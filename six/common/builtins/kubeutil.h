// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2019 Datadog, Inc.
#ifndef DATADOG_AGENT_SIX_KUBEUTIL_H
#define DATADOG_AGENT_SIX_KUBEUTIL_H

/*! \file kubeutil.h
    \brief Six Kubeutil builtin header file.

    The prototypes here defined provide functions to initialize the python kubeutil
    builtin module, and set its relevant callbacks for the six caller.
*/
/*! \fn void _set_get_connection_info_cb(cb_get_connection_info_t)
    \brief Sets a callback to be used by six for kubernetes connection information
    retrieval.
    \param object A function pointer with cb_get_connection_info_t prototype to the
    callback function.

    The callback is expected to be provided by the six caller - in go-context: CGO.
*/

#include <Python.h>
#include <six_types.h>

#define KUBEUTIL_MODULE_NAME "kubeutil"

#ifdef __cplusplus
extern "C" {
#endif

#ifdef DATADOG_AGENT_THREE
//PyMODINIT_FUNC macro already specifies extern "C", nesting these is legal
PyMODINIT_FUNC PyInit_kubeutil(void);
#elif defined(DATADOG_AGENT_TWO)
void Py2_init_kubeutil();
#endif

void _set_get_connection_info_cb(cb_get_connection_info_t);

#ifdef __cplusplus
}
#endif

#endif // DATADOG_AGENT_SIX_KUBEUTIL_H
