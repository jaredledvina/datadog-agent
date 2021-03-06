package testtagger

import (
	"fmt"
	"os"
	"regexp"
	"testing"
)

func TestMain(m *testing.M) {
	err := setUp()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up tests: %v", err)
		os.Exit(-1)
	}

	ret := m.Run()
	tearDown()

	os.Exit(ret)
}

func TestGetTags(t *testing.T) {
	code := fmt.Sprintf(`
	import json
	with open(r'%s', 'w') as f:
		f.write(json.dumps(tagger.get_tags("base", False)))
	`, tmpfile.Name())
	out, err := run(code)
	if err != nil {
		t.Fatal(err)
	}
	if out != "[\"a\", \"b\", \"c\"]" {
		t.Errorf("Unexpected printed value: '%s'", out)
	}
}

func TestGetTagsHighCard(t *testing.T) {
	code := fmt.Sprintf(`
	import json
	with open(r'%s', 'w') as f:
		f.write(json.dumps(tagger.get_tags("base", True)))
	`, tmpfile.Name())
	out, err := run(code)
	if err != nil {
		t.Fatal(err)
	}
	if out != "[\"A\", \"B\", \"C\"]" {
		t.Errorf("Unexpected printed value: '%s'", out)
	}
}

func TestGetTagsUnknown(t *testing.T) {
	code := fmt.Sprintf(`
	import json
	with open(r'%s', 'w') as f:
		f.write(json.dumps(tagger.get_tags("default_switch", True)))
	`, tmpfile.Name())
	out, err := run(code)
	if err != nil {
		t.Fatal(err)
	}
	if out != "[]" {
		t.Errorf("Unexpected printed value: '%s'", out)
	}
}

func TestGetTagsErrorType(t *testing.T) {
	code := fmt.Sprintf(`tagger.get_tags(1234, True)`)
	out, err := run(code)
	if err != nil {
		t.Fatal(err)
	}
	if matched, err := regexp.Match("TypeError: argument 1 must be (str|string), not int", []byte(out)); err != nil && !matched {
		t.Errorf("Unexpected printed value: '%s'", out)
	}
}

func TestTagsLow(t *testing.T) {
	code := fmt.Sprintf(`
	import json
	with open(r'%s', 'w') as f:
		f.write(json.dumps(tagger.get_tags("base", tagger.LOW)))
	`, tmpfile.Name())
	out, err := run(code)
	if err != nil {
		t.Fatal(err)
	}
	if out != "[\"a\", \"b\", \"c\"]" {
		t.Errorf("Unexpected printed value: '%s'", out)
	}
}

func TestTagsHigh(t *testing.T) {
	code := fmt.Sprintf(`
	import json
	with open(r'%s', 'w') as f:
		f.write(json.dumps(tagger.tag("base", tagger.HIGH)))
	`, tmpfile.Name())
	out, err := run(code)
	if err != nil {
		t.Fatal(err)
	}
	if out != "[\"A\", \"B\", \"C\"]" {
		t.Errorf("Unexpected printed value: '%s'", out)
	}
}

func TestTagsOrchestrator(t *testing.T) {
	code := fmt.Sprintf(`
	import json
	with open(r'%s', 'w') as f:
		f.write(json.dumps(tagger.tag("base", tagger.ORCHESTRATOR)))
	`, tmpfile.Name())
	out, err := run(code)
	if err != nil {
		t.Fatal(err)
	}
	if out != "[\"1\", \"2\", \"3\"]" {
		t.Errorf("Unexpected printed value: '%s'", out)
	}
}

func TestTagsInvalidCardinality(t *testing.T) {
	code := fmt.Sprintf(`
	import json
	with open(r'%s', 'w') as f:
		f.write(json.dumps(tagger.tag("default_switch", 123456)))
	`, tmpfile.Name())
	out, err := run(code)
	if err != nil {
		t.Fatal(err)
	}
	if out != "TypeError: Invalid cardinality" {
		t.Errorf("Unexpected printed value: '%s'", out)
	}
}

func TestTagsUnknown(t *testing.T) {
	code := fmt.Sprintf(`
	import json
	with open(r'%s', 'w') as f:
		f.write(json.dumps(tagger.tag("default_switch", tagger.LOW)))
	`, tmpfile.Name())
	out, err := run(code)
	if err != nil {
		t.Fatal(err)
	}
	if out != "[]" {
		t.Errorf("Unexpected printed value: '%s'", out)
	}
}

func TestTagsErrorType(t *testing.T) {
	code := fmt.Sprintf(`tagger.tag(1234, tagger.LOW)`)
	out, err := run(code)
	if err != nil {
		t.Fatal(err)
	}
	if matched, err := regexp.Match("TypeError: argument 1 must be (str|string), not int", []byte(out)); err != nil && !matched {
		t.Errorf("Unexpected printed value: '%s'", out)
	}
}
