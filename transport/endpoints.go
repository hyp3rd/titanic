package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
	"gitlab.com/hyperd/titanic"
)

// Endpoints collects all of the endpoints that compose a People titanic.People.
type Endpoints struct {
	PostPeopleEndpoint    endpoint.Endpoint
	GetPeopleByIDEndpoint endpoint.Endpoint
	PutPeopleEndpoint     endpoint.Endpoint
	PatchPeopleEndpoint   endpoint.Endpoint
	DeletePeopleEndpoint  endpoint.Endpoint
	GetPeopleEndpoint     endpoint.Endpoint
	GetAPIStatusEndpoint  endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the provided titanic.Service. Useful in a cs titanic.Service
// server.
func MakeServerEndpoints(s titanic.Service) Endpoints {
	return Endpoints{
		PostPeopleEndpoint:    MakePostPeopleEndpoint(s),
		GetPeopleByIDEndpoint: MakeGetPeopleByIDEndpoint(s),
		PutPeopleEndpoint:     MakePutPeopleEndpoint(s),
		PatchPeopleEndpoint:   MakePatchPeopleEndpoint(s),
		DeletePeopleEndpoint:  MakeDeletePeopleEndpoint(s),
		GetPeopleEndpoint:     MakeGetPeopleEndpoint(s),
		GetAPIStatusEndpoint:  MakeGetAPIStatusEndpoint(),
	}
}

// MakePostPeopleEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakePostPeopleEndpoint(s titanic.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(PostPeopleRequest)
		id, e := s.PostPeople(ctx, req.People)
		return PostPeopleResponse{ID: id, Err: e}, nil
	}
}

// MakeGetPeopleByIDEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeGetPeopleByIDEndpoint(s titanic.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetPeopleByIDRequest)
		p, e := s.GetPeopleByID(ctx, req.UUID)
		return GetPeopleByIDResponse{People: p, Err: e}, nil
	}
}

// MakePutPeopleEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakePutPeopleEndpoint(s titanic.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(PutPeopleRequest)
		e := s.PutPeople(ctx, req.UUID, req.People)
		return PutPeopleResponse{Err: e}, nil
	}
}

// MakePatchPeopleEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakePatchPeopleEndpoint(s titanic.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(PatchPeopleRequest)
		e := s.PatchPeople(ctx, req.UUID, req.People)
		return PatchPeopleResponse{Err: e}, nil
	}
}

// MakeDeletePeopleEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeDeletePeopleEndpoint(s titanic.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeletePeopleRequest)
		id, e := s.DeletePeople(ctx, req.UUID)
		return DeletePeopleResponse{ID: id, Err: e}, nil
	}
}

// MakeGetPeopleEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeGetPeopleEndpoint(s titanic.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		// req := request.(GetPeopleResponse)
		a, e := s.GetPeople(ctx)
		return GetPeopleResponse{People: a, Err: e}, nil
	}
}

// MakeGetAPIStatusEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeGetAPIStatusEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		return GetAPIStatusResponse{Status: "Healthy"}, nil
	}
}

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

// PostPeopleRequest request object
type PostPeopleRequest struct {
	People titanic.People `json:"people,omitempty"`
}

// PostPeopleResponse response object
type PostPeopleResponse struct {
	ID  string `json:"id,omitempty"`
	Err error  `json:"err,omitempty"`
}

func (r PostPeopleResponse) error() error { return r.Err }

// GetPeopleByIDRequest request object
type GetPeopleByIDRequest struct {
	UUID uuid.UUID `json:"uuid,omitempty"`
}

// GetPeopleByIDResponse response object
type GetPeopleByIDResponse struct {
	People titanic.People `json:"people,omitempty"`
	Err    error          `json:"err,omitempty"`
}

func (r GetPeopleByIDResponse) error() error { return r.Err }

// PutPeopleRequest request object
type PutPeopleRequest struct {
	UUID   uuid.UUID      `json:"uuid,omitempty"`
	People titanic.People `json:"people,omitempty"`
}

// PutPeopleResponse response object
type PutPeopleResponse struct {
	Err error `json:"err,omitempty"`
}

func (r PutPeopleResponse) error() error { return nil }

// PatchPeopleRequest request object
type PatchPeopleRequest struct {
	UUID   uuid.UUID      `json:"uuid,omitempty"`
	People titanic.People `json:"people,omitempty"`
}

// PatchPeopleResponse response object
type PatchPeopleResponse struct {
	Err error `json:"err,omitempty"`
}

func (r PatchPeopleResponse) error() error { return r.Err }

// DeletePeopleRequest request object
type DeletePeopleRequest struct {
	UUID uuid.UUID
}

// DeletePeopleResponse response object
type DeletePeopleResponse struct {
	ID  string `json:"id,omitempty"`
	Err error  `json:"err,omitempty"`
}

func (r DeletePeopleResponse) error() error { return r.Err }

// GetPeopleRequest struct
type GetPeopleRequest struct{}

// GetPeopleResponse response object
type GetPeopleResponse struct {
	People []titanic.People `json:"people,omitempty"`
	Err    error            `json:"err,omitempty"`
}

func (r GetPeopleResponse) error() error { return r.Err }

// GetAPIStatusRequest request object
type GetAPIStatusRequest struct{}

// GetAPIStatusResponse response object
type GetAPIStatusResponse struct {
	Status string `json:"status,omitempty"`
	Err    error  `json:"err,omitempty"`
}

func (r GetAPIStatusResponse) error() error { return r.Err }
