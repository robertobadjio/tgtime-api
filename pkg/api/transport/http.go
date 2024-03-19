package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	util "officetime-api/internal/util"
	"officetime-api/pkg/api/endpoints"
	"strconv"
)

//var logger log.Logger

/*func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}*/

type errorer interface {
	Error() error
}

func NewHTTPHandler(ep endpoints.Set) http.Handler {
	router := mux.NewRouter()

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	router.Methods(http.MethodGet).Path("/service/status").Handler(
		httptransport.NewServer(
			ep.ServiceStatusEndpoint,
			decodeHTTPServiceStatusRequest,
			encodeResponse,
			opts...,
		),
	)

	var api = router.PathPrefix("/api").Subrouter()

	var api1 = api.
		PathPrefix("/v1").
		Subrouter()

	api1.Methods(http.MethodPost).
		Path("/login").
		Handler(
			httptransport.NewServer(
				ep.LoginEndpoint,
				decodeHTTPLoginRequest,
				encodeResponse,
				opts...,
			),
		)
	api1.Methods(http.MethodGet).
		Path("/router").
		Handler(
			httptransport.NewServer(
				ep.GetRoutersEndpoint,
				decodeHTTPGetRoutersRequest,
				encodeResponse,
				opts...,
			),
		)

	api1.Methods(http.MethodGet).
		Path("/router/{routerId}").
		Handler(
			httptransport.NewServer(
				ep.GetRouterEndpoint,
				decodeHTTPGetRouterRequest,
				encodeResponse,
				opts...,
			),
		)

	api1.Methods(http.MethodPost).
		Path("/router").
		Handler(
			httptransport.NewServer(
				ep.CreateRouterEndpoint,
				decodeHTTPCreateRouterRequest,
				encodeResponse,
				opts...,
			),
		)

	api1.Methods(http.MethodPut).
		Path("/router/{routerId}").
		Handler(
			httptransport.NewServer(
				ep.UpdateRouterEndpoint,
				decodeHTTPUpdateRouterRequest,
				encodeResponse,
				opts...,
			),
		)

	api1.Methods(http.MethodDelete).
		Path("/router/{routerId}").
		Handler(
			httptransport.NewServer(
				ep.DeleteRouterEndpoint,
				decodeHTTPDeleteRouterRequest,
				encodeResponse,
				opts...,
			),
		)

	return router
}

func decodeHTTPLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.LoginRequest
	if r.ContentLength == 0 {
		//logger.Log("Get request with no body")
		return req, nil
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeHTTPServiceStatusRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.ServiceStatusRequest
	return req, nil
}

func decodeHTTPGetRoutersRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.GetRoutersRequest
	if r.ContentLength == 0 {
		//logger.Log("Get request with no body")
		return req, nil
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeHTTPGetRouterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.GetRouterRequest
	params := mux.Vars(r)

	routerId, err := strconv.Atoi(params["routerId"])
	if err != nil {
		return nil, fmt.Errorf("invalid path param router identifier")
	}

	req.RouterId = routerId

	return req, nil
}

func decodeHTTPCreateRouterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.CreateRouterRequest

	err := json.NewDecoder(r.Body).Decode(&req.Router)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("error decode json body: %v", err))
	}

	return req, nil
}

func decodeHTTPUpdateRouterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.UpdateRouterRequest
	params := mux.Vars(r)

	routerId, err := strconv.Atoi(params["routerId"])
	if err != nil {
		return nil, fmt.Errorf("invalid path param router identifier")
	}
	req.RouterId = routerId

	err = json.NewDecoder(r.Body).Decode(&req.Router)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("error decode json body: %v", err))
	}

	return req, nil
}

func decodeHTTPDeleteRouterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.DeleteRouterRequest
	params := mux.Vars(r)

	routerId, err := strconv.Atoi(params["routerId"])
	if err != nil {
		return nil, fmt.Errorf("invalid path param router identifier")
	}

	req.RouterId = routerId

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if e, ok := response.(errorer); ok && e.Error() != nil {
		encodeError(ctx, e.Error(), w)
		return nil
	}

	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch {
	/*case errors.Is(err, util.ErrInvalidArgument):
	w.WriteHeader(http.StatusBadRequest)*/
	case errors.Is(err, util.ErrRouterNotFound):
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
