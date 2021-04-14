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

	for _, user := range model.GetAllUsers(0, 0).Users {
		times := dao.GetAllByDate(user.MacAddress, date, 0)
		seconds := model.AggregateDayTotalTime(times)
		breaks := dao.GetAllBreaksByTimesOld(times)
		breaksJson, err := json.Marshal(breaks)
		begin := dao.GetDayTimeFromTimeTable(user.MacAddress, date, "ASC")
		end := dao.GetDayTimeFromTimeTable(user.MacAddress, date, "DESC")

		_, err = Db.Exec("INSERT INTO time_summary (mac_address, date, seconds, breaks, seconds_begin, seconds_end) VALUES ($1, $2, $3, $4, $5, $6)", user.MacAddress, date, seconds, breaksJson, begin, end)
		if err != nil {
			panic(err)
		}
	}
}
