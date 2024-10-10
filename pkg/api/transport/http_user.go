package transport

import (
	"context"
	"encoding/json"
	"github.com/robertobadjio/tgtime-api/pkg/api/endpoints"
	"net/http"
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
