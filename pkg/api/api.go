// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Base API that all services implement and extend
// ----------------------------------------------------------------------------

package api

import (
	"github.com/gorilla/mux"
)

// Base holds a standard set of values for all services & APIs
type Base struct {
	ServiceName string
	Healthy     bool
	Version     string
	BuildInfo   string
}

//
// New creates and returns a new Base API instance
//
func NewBase(name, ver, info string, healthy bool, router *mux.Router) *Base {
	b := &Base{
		ServiceName: name,
		Healthy:     healthy,
		Version:     ver,
		BuildInfo:   info,
	}

	return b
}
