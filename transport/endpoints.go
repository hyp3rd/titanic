package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"gitlab.com/hyperd/titanic"
)

// Endpoints collects all of the endpoints that compose a People titanic.People.
type Endpoints struct {
	PostPeopleEndpoint   endpoint.Endpoint
	GetPeopleEndpoint    endpoint.Endpoint
	PutPeopleEndpoint    endpoint.Endpoint
	PatchPeopleEndpoint  endpoint.Endpoint
	DeletePeopleEndpoint endpoint.Endpoint
	GetAllPeopleEndpoint endpoint.Endpoint
	GetAPIStatusEndpoint endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the provided titanic.Service. Useful in a cs titanic.Service
// server.
func MakeServerEndpoints(s titanic.Service) Endpoints {
	return Endpoints{
		PostPeopleEndpoint:   MakePostPeopleEndpoint(s),
		GetPeopleEndpoint:    MakeGetPeopleEndpoint(s),
		PutPeopleEndpoint:    MakePutPeopleEndpoint(s),
		PatchPeopleEndpoint:  MakePatchPeopleEndpoint(s),
		DeletePeopleEndpoint: MakeDeletePeopleEndpoint(s),
		GetAllPeopleEndpoint: MakeGetAllPeopleEndpoint(s),
		GetAPIStatusEndpoint: MakeGetAPIStatusEndpoint(s),
	}
}

// MakePostPeopleEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakePostPeopleEndpoint(s titanic.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(PostPeopleRequest)
		e := s.PostPeople(ctx, req.People)
		return PostPeopleResponse{Err: e}, nil
	}
}

// MakeGetPeopleEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeGetPeopleEndpoint(s titanic.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetPeopleRequest)
		p, e := s.GetPeople(ctx, req.UUID)
		return GetPeopleResponse{People: p, Err: e}, nil
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
		e := s.DeletePeople(ctx, req.UUID)
		return DeletePeopleResponse{Err: e}, nil
	}
}

// MakeGetAllPeopleEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeGetAllPeopleEndpoint(s titanic.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		// req := request.(getAllPeopleRequest)
		a, e := s.GetAllPeople(ctx)
		return GetAllPeopleResponse{AllPeople: a, Err: e}, nil
	}
}

// MakeGetAPIStatusEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeGetAPIStatusEndpoint(s titanic.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		// req := request.(getAllPeopleRequest)
		s, e := s.GetAPIStatus(ctx)
		return GetAPIStatusResponse{Status: s, Err: e}, nil
	}
}
