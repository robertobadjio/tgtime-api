package model

import (
	"database/sql"
	"officetime-api/app/dao"
)

var Db *sql.DB

func GetAllByDate(macAddress string, date string) []*dao.TimeUser {
	rows, err := Db.Query("SELECT t.mac_address, t.second FROM time t WHERE t.mac_address = $1 AND t.second BETWEEN $2 AND $3 ORDER BY t.second", macAddress, dao.GetSecondsByBeginDate(date), dao.GetSecondsByEndDate(date))
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	times := make([]*dao.TimeUser, 0)
	for rows.Next() {
		time := new(dao.TimeUser)
		err := rows.Scan(&time.MacAddress, &time.Second)
		if err != nil {
			panic(err)
		}

		times = append(times, time)
	}

	return times
}
