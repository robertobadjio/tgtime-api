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
