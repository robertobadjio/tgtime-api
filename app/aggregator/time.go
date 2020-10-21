package aggregator

import (
	"database/sql"
	"encoding/json"
	"officetime-api/app/dao"
	"officetime-api/app/model"
)

var Db *sql.DB

func AggregateTime() {
	macAddress := "7c:a1:ae:a9:c0:d8"
	date := "2020-10-21"
	times := model.GetAllByDate(macAddress, date)
	seconds := dao.AggregateDayTotalTime(times)
	breaks := dao.GetAllBreaksByTimes(times)
	breaksJson, err := json.Marshal(breaks)
	begin := dao.GetDayTime(macAddress, date, "ASC")
	end := dao.GetDayTime(macAddress, date, "DESC")

	_, err = Db.Exec("INSERT INTO time_summary (mac_address, date, seconds, breaks, seconds_begin, seconds_end) VALUES ($1, $2, $3, $4, $5, $6)", macAddress, date, seconds, breaksJson, begin, end)
	if err != nil {
		panic(err)
	}
}
