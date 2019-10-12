package transport

import (
	"github.com/google/uuid"
	"gitlab.com/hyperd/titanic"
)

// We have two options to return errors from the business logic.
//
// We could return the error via the endpoint itself. That makes certain things
// a little bit easier, like providing non-200 HTTP responses to the client. But
// Go kit assumes that endpoint errors are (or may be treated as)
// transport-domain errors. For example, an endpoint error will count against a
// circuit breaker error count.
//
// Therefore, it's often better to return service (business logic) errors in the
// response object. This means we have to do a bit more work in the HTTP
// response encoder to detect e.g. a not-found error and provide a proper HTTP
// status code. That work is done with the errorer interface, in transport.go.
// Response types that may contain business-logic errors implement that
// interface.

type PostPeopleRequest struct {
	People titanic.People `json:"people,omitempty"`
}

type PostPeopleResponse struct {
	Err error `json:"err,omitempty"`
}

func (r PostPeopleResponse) error() error { return r.Err }

type GetPeopleRequest struct {
	UUID uuid.UUID `json:"uuid,omitempty"`
}

type GetPeopleResponse struct {
	People titanic.People `json:"people,omitempty"`
	Err    error          `json:"err,omitempty"`
}

func (r GetPeopleResponse) error() error { return r.Err }

type PutPeopleRequest struct {
	UUID   uuid.UUID      `json:"uuid,omitempty"`
	People titanic.People `json:"people,omitempty"`
}

type PutPeopleResponse struct {
	Err error `json:"err,omitempty"`
}

func (r PutPeopleResponse) error() error { return nil }

type PatchPeopleRequest struct {
	UUID   uuid.UUID      `json:"uuid,omitempty"`
	People titanic.People `json:"people,omitempty"`
}

type PatchPeopleResponse struct {
	Err error `json:"err,omitempty"`
}

func (r PatchPeopleResponse) error() error { return r.Err }

type DeletePeopleRequest struct {
	UUID uuid.UUID
}

type DeletePeopleResponse struct {
	Err error `json:"err,omitempty"`
}

func (r DeletePeopleResponse) error() error { return r.Err }

// GetAllPeopleRequest struct
type GetAllPeopleRequest struct {
}

// GetAllPeopleResponse struct
type GetAllPeopleResponse struct {
	AllPeople []titanic.People `json:"people,omitempty"`
	Err       error            `json:"err,omitempty"`
}

func (r GetAllPeopleResponse) error() error { return r.Err }

type GetAPIStatusRequest struct {
}

type GetAPIStatusResponse struct {
	Status string `json:"status,omitempty"`
	Err    error  `json:"err,omitempty"`
}

func (r GetAPIStatusResponse) error() error { return r.Err }
