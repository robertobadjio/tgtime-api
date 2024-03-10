package transport

import (
	"context"
	"encoding/json"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	"officetime-api/internal/util"
	"officetime-api/pkg/api/endpoints"
)

//var logger log.Logger

/*func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}*/

func NewHTTPHandler(ep endpoints.Set) http.Handler {
	router := mux.NewRouter()

	router.Methods(http.MethodGet).Path("/service/status").Handler(
		httptransport.NewServer(
			ep.ServiceStatusEndpoint,
			decodeHTTPServiceStatusRequest,
			encodeResponse,
		),
	)

	var api = router.PathPrefix("/api").Subrouter()

	var api1 = api.PathPrefix("/v1").Subrouter()
	api1.Methods(http.MethodPost).Path("/login").Handler(
		httptransport.NewServer(
			ep.LoginEndpoint,
			decodeHTTPLoginRequest,
			encodeResponse,
		),
	)
	api1.Methods(http.MethodGet).Path("/router").Handler(
		httptransport.NewServer(
			ep.GetRoutersEndpoint,
			decodeHTTPGetRoutersRequest,
			encodeResponse,
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

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if e, ok := response.(error); ok && e != nil {
		encodeError(ctx, e, w)
		return nil
	}
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case util.ErrUnknown:
		w.WriteHeader(http.StatusNotFound)
	case util.ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
