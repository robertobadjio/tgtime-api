package router

type Router struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Address     string `json:"address"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	Status      bool   `json:"status"`
	WorkTime    bool   `json:"workTime"`
}

func NewRouter() (*Router, error) {
	return &Router{}, nil
}
