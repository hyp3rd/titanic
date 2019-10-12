package middleware

import (
	"gitlab.com/hyperd/titanic"
)

// Middleware describes the titanic service (as opposed to endpoint) middleware.
type Middleware func(titanic.Service) titanic.Service
