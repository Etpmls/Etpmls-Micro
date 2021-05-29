package em

import em_library "github.com/Etpmls/Etpmls-Micro/v3/library"

var Micro center

type center struct {
	Config *em_library.Configuration
	Request *request
	Auth *auth
	Middleware *middleware
}