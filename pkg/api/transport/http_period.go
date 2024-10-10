package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/robertobadjio/tgtime-api/pkg/api/endpoints"
	"net/http"
	"strconv"
)

func decodeHTTPGetPeriodsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.GetPeriodsRequest
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

func decodeHTTPGetPeriodRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.GetPeriodRequest
	params := mux.Vars(r)

	periodId, err := strconv.Atoi(params["periodId"])
	if err != nil {
		return nil, fmt.Errorf("invalid path param period identifier")
	}

	req.PeriodId = periodId

	return req, nil
}

func decodeHTTPGetPeriodCurrentRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	var req endpoints.GetPeriodCurrentRequest

	return req, nil
}

func decodeHTTPCreatePeriodRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.CreatePeriodRequest

	err := json.NewDecoder(r.Body).Decode(&req.Period)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("error decode json body: %v", err))
	}

	return req, nil
}

func decodeHTTPUpdatePeriodRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.UpdatePeriodRequest
	params := mux.Vars(r)

	periodId, err := strconv.Atoi(params["periodId"])
	if err != nil {
		return nil, fmt.Errorf("invalid path param period identifier")
	}
	req.PeriodId = periodId

	err = json.NewDecoder(r.Body).Decode(&req.Period)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("error decode json body: %v", err))
	}

	return req, nil
}

func decodeHTTPDeletePeriodRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.DeletePeriodRequest
	params := mux.Vars(r)

	periodId, err := strconv.Atoi(params["periodId"])
	if err != nil {
		return nil, fmt.Errorf("invalid path param period identifier")
	}

	req.PeriodId = periodId

	return req, nil
}
