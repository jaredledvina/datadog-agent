// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.

package security

import (
	"fmt"
	"io/ioutil"
	"os/user"
	"syscall"

	"github.com/DataDog/datadog-agent/pkg/util/log"
	acl "github.com/hectane/go-acl"
	"golang.org/x/sys/windows"
)

var (
	wellKnownSidStrings = map[string]string{
		"Administrators": "S-1-5-32-544",
		"System":         "S-1-5-18",
		"Users":          "S-1-5-32-545",
	}
	wellKnownSids = make(map[string]*windows.SID)
)

func init() {
	for key, val := range wellKnownSidStrings {
		sid, err := windows.StringToSid(val)
		if err == nil {
			wellKnownSids[key] = sid
		}
	}
}

// lookupUsernameAndDomain obtains the username and domain for usid.
func lookupUsernameAndDomain(usid *syscall.SID) (username, domain string, e error) {
	username, domain, t, e := usid.LookupAccount("")
	if e != nil {
		return "", "", e
	}
	if t != syscall.SidTypeUser {
		return "", "", fmt.Errorf("user: should be user account type, not %d", t)
	}
	return username, domain, nil
}

// writes auth token(s) to a file with the same permissions as datadog.yaml
func saveAuthToken(token, tokenPath string) error {
	// get the current user
	var sidString string
	currUser, err := user.Current()
	if err != nil {
		log.Warnf("Unable to get current user %v", err)
		log.Infof("Attempting to get current user information directly")
		tok, e := syscall.OpenCurrentProcessToken()
		if e != nil {
			log.Warnf("Couldn't get process token %v", e)
			return e
		}
		defer tok.Close()
		user, e := tok.GetTokenUser()
		if e != nil {
			log.Warnf("Couldn't get  token user %v", e)
			return e
		}
		sidString, e = user.User.Sid.String()
		if e != nil {
			log.Warnf("Couldn't get  user sid string %v", e)
			return e
		}

		log.Infof("Got sidstring from token user")

		// now just do some debugging, see what we weren't able to get.
		pg, e := tok.GetTokenPrimaryGroup()
		if e != nil {
			log.Warnf("Would have failed getting token PG %v", e)
		}
		_, e = pg.PrimaryGroup.String()
		if e != nil {
			log.Warnf("Would have failed getting  PG  string %v", e)
		}
		dir, e := tok.GetUserProfileDirectory()
		if e != nil {
			log.Warnf("Would have failed getting  primary directory %v", e)
		} else {
			log.Infof("Profile directory is %v", dir)
		}
		username, domain, e := lookupUsernameAndDomain(user.User.Sid)
		if e != nil {
			log.Warnf("Would have failed getting username and domain %v", e)
		} else {
			log.Infof("Username/domain is %v %v", username, domain)
		}

	} else {
		log.Infof("Getting sidstring from current user")
		sidString = currUser.Uid
	}
	currUserSid, err := windows.StringToSid(sidString)
	if err != nil {
		log.Warnf("Unable to get current user sid %v", err)
		return err
	}
	err = ioutil.WriteFile(tokenPath, []byte(token), 0755)
	if err == nil {
		err = acl.Apply(
			tokenPath,
			true,  // replace the file permissions
			false, // don't inherit
			acl.GrantSid(windows.GENERIC_ALL, wellKnownSids["Administrators"]),
			acl.GrantSid(windows.GENERIC_ALL, wellKnownSids["System"]),
			acl.GrantSid(windows.GENERIC_ALL, currUserSid))
		log.Infof("Wrote auth token acl %v", err)
	}
	return err
}
