package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"gitlab.com/hyperd/titanic"
	"gitlab.com/hyperd/titanic/transport"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

// MakeHTTPHandler mounts all of the service endpoints into an http.Handler.
func MakeHTTPHandler(s titanic.Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := transport.MakeServerEndpoints(s)
	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	// POST    /people/                       	   adds another passenger to the people collection
	// GET     /people/:uuid                       retrieves the given passenger by uuid from the people collection
	// PUT     /people/:uuid                       post updated information about a passenger (uuid)
	// PATCH   /people/:uuid                       partial update of the passenger information
	// DELETE  /people/:uuid                       removes the given passenger
	// GET     /people/           				   retrieves all the passengers from the people collection
	// GET     /           						   returns the API status

	r.Methods("POST").Path("/people/").Handler(kithttp.NewServer(
		e.PostPeopleEndpoint,
		decodePostPeopleRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/people/{uuid}").Handler(kithttp.NewServer(
		e.GetPeopleByIDEndpoint,
		decodeGetPeopleByIDRequest,
		encodeResponse,
		options...,
	))
	r.Methods("PUT").Path("/people/{uuid}").Handler(kithttp.NewServer(
		e.PutPeopleEndpoint,
		decodePutPeopleRequest,
		encodeResponse,
		options...,
	))
	r.Methods("PATCH").Path("/people/{uuid}").Handler(kithttp.NewServer(
		e.PatchPeopleEndpoint,
		decodePatchPeopleRequest,
		encodeResponse,
		options...,
	))
	r.Methods("DELETE").Path("/people/{uuid}").Handler(kithttp.NewServer(
		e.DeletePeopleEndpoint,
		decodeDeletePeopleRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/people/").Handler(kithttp.NewServer(
		e.GetPeopleEndpoint,
		decodeGetPeopleRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/").Handler(kithttp.NewServer(
		e.GetAPIStatusEndpoint,
		decodeGetAPIStatusRequest,
		encodeStatusResponse,
		options...,
	))
	return r
}

func decodePostPeopleRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req transport.PostPeopleRequest
	if e := json.NewDecoder(r.Body).Decode(&req.People); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeGetPeopleByIDRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["uuid"])

	if err != nil {
		return nil, ErrBadRouting
	}

	return transport.GetPeopleByIDRequest{ID: id}, nil
}

func decodePutPeopleRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["uuid"])

	if err != nil {
		return nil, ErrBadRouting
	}

	var people titanic.People
	if err := json.NewDecoder(r.Body).Decode(&people); err != nil {
		return nil, err
	}
	return transport.PutPeopleRequest{
		ID:     id,
		People: people,
	}, nil
}

func decodePatchPeopleRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["uuid"])

	if err != nil {
		return nil, ErrBadRouting
	}

	var people titanic.People
	if err := json.NewDecoder(r.Body).Decode(&people); err != nil {
		return nil, err
	}
	return transport.PatchPeopleRequest{
		ID:     id,
		People: people,
	}, nil
}

func decodeDeletePeopleRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["uuid"])

	if err != nil {
		return nil, ErrBadRouting
	}

	return transport.DeletePeopleRequest{ID: id}, nil
}

func decodeGetPeopleRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return transport.GetPeopleRequest{}, nil
}

func decodeGetAPIStatusRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return transport.GetAPIStatusRequest{}, nil
}

func encodePostPeopleRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("POST").Path("/people/")
	req.URL.Path = "/people/"
	return encodeRequest(ctx, req, request)
}

func encodeGetPeopleByIDRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("GET").Path("/people/{uuid}")
	r := request.(transport.GetPeopleByIDRequest)
	peopleID := url.QueryEscape(r.ID.String())
	req.URL.Path = "/people/" + peopleID
	return encodeRequest(ctx, req, request)
}

func encodePutPeopleRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("PUT").Path("/people/{uuid}")
	r := request.(transport.PutPeopleRequest)
	peopleID := url.QueryEscape(r.ID.String())
	req.URL.Path = "/people/" + peopleID
	return encodeRequest(ctx, req, request)
}

func encodePatchPeopleRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("PATCH").Path("/people/{uuid}")
	r := request.(transport.PatchPeopleRequest)
	peopleID := url.QueryEscape(r.ID.String())
	req.URL.Path = "/people/" + peopleID
	return encodeRequest(ctx, req, request)
}

func encodeDeletePeopleRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("DELETE").Path("/people/{uuid}")
	r := request.(transport.DeletePeopleRequest)
	peopleID := url.QueryEscape(r.ID.String())
	req.URL.Path = "/people/" + peopleID
	return encodeRequest(ctx, req, request)
}

func encodeGetPeopleRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("GET").Path("/people/")
	req.URL.Path = "/people/"
	return encodeRequest(ctx, req, request)
}

func encodeGetAPIStatusRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("GET").Path("/")
	req.URL.Path = "/"
	return encodeRequest(ctx, req, request)
}

func decodePostPeopleResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response transport.PostPeopleResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeGetPeopleByIDResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response transport.GetPeopleByIDRequest
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodePutPeopleResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response transport.PutPeopleResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodePatchPeopleResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response transport.PatchPeopleResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeDeletePeopleResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response transport.DeletePeopleResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeGetPeopleResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response transport.GetPeopleResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeGetAPIStatusResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response transport.GetAPIStatusResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error. For more information, read the
// big comment in endpoints.go.
type errorer interface {
	error() error
}

// encodeResponse is the common method to encode all response types to the
// client. I chose to do it this way because, since we're using JSON, there's no
// reason to provide anything more specific. It's certainly possible to
// specialize on a per-response (per-method) basis.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}

	// Configure the http security headers
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// w.Header().Set("X-XSS-Protection", "1; mode=block")
	// w.Header().Set("X-Content-Type-Options", "nosniff")
	// w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	// w.Header().Set("Content-Security-Policy", "upgrade-insecure-requests;")

	return json.NewEncoder(w).Encode(response)
}

func encodeStatusResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	started := time.Now()

	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	duration := time.Since(started)

	if duration.Seconds() > 10 {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("error: %v", duration.Seconds())))
	} else {
		w.WriteHeader(200)
		// w.Write([]byte("ok"))
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// encodeRequest likewise JSON-encodes the request to the HTTP request body.
// Don't use it directly as a transport/http.Client EncodeRequestFunc:
// titanic endpoints require mutating the HTTP method and request path.
func encodeRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case titanic.ErrNotFound:
		return http.StatusNotFound
	case titanic.ErrAlreadyExists, titanic.ErrInconsistentIDs:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
