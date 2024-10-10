package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/robertobadjio/tgtime-api/pkg/api/endpoints"
	"net/http"
	"strconv"
	"strings"
)

func decodeHTTPGetUsersRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.GetUsersRequest
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

func decodeHTTPGetUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.GetUserRequest
	params := mux.Vars(r)

	userId, err := strconv.Atoi(params["userId"])
	if err != nil {
		return nil, fmt.Errorf("invalid path param user identifier")
	}

	req.UserId = userId

	return req, nil
}

func decodeHTTPGetUserByMacAddressRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.GetUserByMacAddressRequest
	params := mux.Vars(r)

	req.MacAddress = strings.TrimSpace(params["macAddress"])

	return req, nil
}
