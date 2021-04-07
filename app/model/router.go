package model

import "log"

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

func GetAllRouters() []*Router {
	rows, err := Db.Query("SELECT id, name, description, address, login, password, status, work_time FROM router")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	routers := make([]*Router, 0)
	for rows.Next() {
		router := new(Router)
		err := rows.Scan(&router.Id, &router.Name, &router.Description, &router.Address, &router.Login, &router.Password, &router.Status, &router.WorkTime)
		if err != nil {
			log.Fatal(err)
		}

		routers = append(routers, router)
	}

	return routers
}
