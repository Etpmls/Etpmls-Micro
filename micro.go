package em

import em_library "github.com/Etpmls/Etpmls-Micro/library"

var Micro center

type center struct {
	Config *em_library.Configuration
	Response *response
	Request *request
	Auth *auth
	Client *client
	Middleware *middleware
}