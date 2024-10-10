package transport

import (
	"context"
	"encoding/json"
	"errors"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	util "github.com/robertobadjio/tgtime-api/internal/util"
	"github.com/robertobadjio/tgtime-api/pkg/api/endpoints"
	"net/http"
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

	api1.Methods(http.MethodGet).
		Path("/period").
		Handler(
			httptransport.NewServer(
				ep.GetPeriodsEndpoint,
				decodeHTTPGetPeriodsRequest,
				encodeResponse,
				opts...,
			),
		)
	api1.Methods(http.MethodGet).
		Path("/period/{periodId}").
		Handler(
			httptransport.NewServer(
				ep.GetPeriodEndpoint,
				decodeHTTPGetPeriodRequest,
				encodeResponse,
				opts...,
			),
		)
	api1.Methods(http.MethodGet).
		Path("/period/current").
		Handler(
			httptransport.NewServer(
				ep.GetPeriodCurrentEndpoint,
				decodeHTTPGetPeriodCurrentRequest,
				encodeResponse,
				opts...,
			),
		)
	api1.Methods(http.MethodPost).
		Path("/period").
		Handler(
			httptransport.NewServer(
				ep.CreatePeriodEndpoint,
				decodeHTTPCreatePeriodRequest,
				encodeResponse,
				opts...,
			),
		)
	api1.Methods(http.MethodPut).
		Path("/period/{periodId}").
		Handler(
			httptransport.NewServer(
				ep.UpdatePeriodEndpoint,
				decodeHTTPUpdatePeriodRequest,
				encodeResponse,
				opts...,
			),
		)
	api1.Methods(http.MethodDelete).
		Path("/period/{periodId}").
		Handler(
			httptransport.NewServer(
				ep.DeletePeriodEndpoint,
				decodeHTTPDeletePeriodRequest,
				encodeResponse,
				opts...,
			),
		)

	api1.Methods(http.MethodGet).
		Path("/department").
		Handler(
			httptransport.NewServer(
				ep.GetDepartmentsEndpoint,
				decodeHTTPGetDepartmentsRequest,
				encodeResponse,
				opts...,
			),
		)
	api1.Methods(http.MethodGet).
		Path("/department/{departmentId}").
		Handler(
			httptransport.NewServer(
				ep.GetDepartmentEndpoint,
				decodeHTTPGetDepartmentRequest,
				encodeResponse,
				opts...,
			),
		)
	api1.Methods(http.MethodPost).
		Path("/department").
		Handler(
			httptransport.NewServer(
				ep.CreateDepartmentEndpoint,
				decodeHTTPCreateDepartmentRequest,
				encodeResponse,
				opts...,
			),
		)
	api1.Methods(http.MethodPut).
		Path("/department/{departmentId}").
		Handler(
			httptransport.NewServer(
				ep.UpdateDepartmentEndpoint,
				decodeHTTPUpdateDepartmentRequest,
				encodeResponse,
				opts...,
			),
		)
	api1.Methods(http.MethodDelete).
		Path("/department/{departmentId}").
		Handler(
			httptransport.NewServer(
				ep.DeleteDepartmentEndpoint,
				decodeHTTPDeleteDepartmentRequest,
				encodeResponse,
				opts...,
			),
		)

	api1.Methods(http.MethodGet).
		Path("/user").
		Handler(
			httptransport.NewServer(
				ep.GetUsersEndpoint,
				decodeHTTPGetUsersRequest,
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
