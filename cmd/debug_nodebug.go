//go:build nodebug

package main

import "net/http"

func registerDebugHandlers(_ *http.ServeMux) {}
