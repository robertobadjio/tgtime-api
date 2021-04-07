package dao

import (
	"encoding/json"
	"net/http"
	"officetime-api/app/model"
)

type Data struct {
	Content []*model.Router `json:"content"`
}

func GetAllRouters(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := new(Data)
	data.Content = model.GetAllRouters()

	json.NewEncoder(w).Encode(data)
}
