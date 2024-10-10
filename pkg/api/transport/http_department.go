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

func decodeHTTPGetDepartmentsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.GetDepartmentsRequest
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

func decodeHTTPGetDepartmentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.GetDepartmentRequest
	params := mux.Vars(r)

	departmentId, err := strconv.Atoi(params["departmentId"])
	if err != nil {
		return nil, fmt.Errorf("invalid path param department identifier")
	}

	req.DepartmentId = departmentId

	return req, nil
}

func decodeHTTPCreateDepartmentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.CreateDepartmentRequest

	err := json.NewDecoder(r.Body).Decode(&req.Department)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("error decode json body: %v", err))
	}

	return req, nil
}

func decodeHTTPUpdateDepartmentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.UpdateDepartmentRequest
	params := mux.Vars(r)

	departmentId, err := strconv.Atoi(params["departmentId"])
	if err != nil {
		return nil, fmt.Errorf("invalid path param department identifier")
	}
	req.DepartmentId = departmentId

	err = json.NewDecoder(r.Body).Decode(&req.Department)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("error decode json body: %v", err))
	}

	return req, nil
}

func decodeHTTPDeleteDepartmentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.DeleteDepartmentRequest
	params := mux.Vars(r)

	departmentId, err := strconv.Atoi(params["departmentId"])
	if err != nil {
		return nil, fmt.Errorf("invalid path param department identifier")
	}

	req.DepartmentId = departmentId

	return req, nil
}
