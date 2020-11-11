package aggregator

import (
	"database/sql"
	"encoding/json"
	"officetime-api/app/dao"
	"officetime-api/app/model"
	"time"
)

var Db *sql.DB

func AggregateTime() {
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	date := time.Now().AddDate(0, 0, -1).In(moscowLocation).Format("2006-01-02")

	for _, user := range model.GetAllUsers().Users {
		times := dao.GetAllByDate(user.MacAddress, date)
		seconds := dao.AggregateDayTotalTime(times)
		breaks := dao.GetAllBreaksByTimes(user.MacAddress, date)
		breaksJson, err := json.Marshal(breaks)
		begin := dao.GetDayTime(user.MacAddress, date, "ASC")
		end := dao.GetDayTime(user.MacAddress, date, "DESC")

		_, err = Db.Exec("INSERT INTO time_summary (mac_address, date, seconds, breaks, seconds_begin, seconds_end) VALUES ($1, $2, $3, $4, $5, $6)", user.MacAddress, date, seconds, breaksJson, begin, end)
		if err != nil {
			panic(err)
		}
	}
}
