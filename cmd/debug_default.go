//go:build !nodebug

package main

import (
	"net/http"
	"net/http/pprof"
)

// registerDebugHandlers wires the standard /debug/pprof/ endpoints. Build with
// `-tags nodebug` to strip them out of production binaries.
func registerDebugHandlers(mux *http.ServeMux) {
	mux.HandleFunc("GET /debug/pprof/", pprof.Index)
	mux.HandleFunc("GET /debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("GET /debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("GET /debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("GET /debug/pprof/trace", pprof.Trace)
}
