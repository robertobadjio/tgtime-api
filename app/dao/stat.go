package dao

import (
	"encoding/json"
	"net/http"
	"officetime-api/app/model"
)

func GetStatByPeriodsAndRouters(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(model.GetAllTimesByPeriodsAndRouters())
}
