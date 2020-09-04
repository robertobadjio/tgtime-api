package dao

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

type TimeUser struct {
	macAddress string
	Second     int64
}

type Time struct {
	Hours     int   `json:"hours"`
	Minutes   int   `json:"minutes"`
	BeginTime int64 `json:"beginTime"`
}

func GetTimeDayAll(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["id"]
	date := mux.Vars(r)["date"]

	database := GetDB()
	var macAddress string
	var telegramId int64
	row := database.QueryRow("SELECT u.mac_address, u.telegram_id FROM users u WHERE u.id = $1", userId)
	err := row.Scan(&macAddress, &telegramId)
	if err != nil {
		panic(err)
	}

	var timeOutput Time

	times := getDayTimesByUser(macAddress, date)
	totalDayTime := aggregateDayTotalTime(times)
	timeOutput.Hours = totalDayTime / 3600
	timeOutput.Minutes = (totalDayTime / 60) - (timeOutput.Hours * 60)

	beginTimeSeconds := beginDayTime(macAddress, date)
	timeOutput.BeginTime = beginTimeSeconds // TODO: Если дата меньше сегодняшней, пишем "Вы сегодня не были на работе"

	json.NewEncoder(w).Encode(timeOutput)
}

func aggregateDayTotalTime(times []*TimeUser) int {
	num := 0
	for i, time := range times {
		if i == 0 {
			continue
		}

		delta := time.Second - times[i-1].Second
		if delta <= 33 {
			num++
		}
	}

	return num * 30
}

func getSecondsByDate(date string) int64 {
	t, _ := time.Parse("2006-01-02", date)
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location()).Unix()
}

func beginDayTime(macAddress string, date string) int64 {
	database := GetDB()
	var beginSecond int64
	row := database.QueryRow("SELECT t.second FROM time t WHERE t.mac_address = $1 AND t.second > $2 ORDER BY t.second ASC LIMIT 1", macAddress, getSecondsByDate(date))
	err := row.Scan(&beginSecond)

	if err == sql.ErrNoRows {
		return 0
	}

	if err != nil {
		panic(err)
	}

	return beginSecond
}

func getDayTimesByUser(macAddress string, date string) []*TimeUser {
	database := GetDB()
	rows, err := database.Query("SELECT t.mac_address, t.second FROM time t WHERE t.mac_address = $1 AND t.second >= $2 ORDER BY t.second", macAddress, getSecondsByDate(date))
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	times := make([]*TimeUser, 0)
	for rows.Next() {
		timeUser := new(TimeUser)
		err := rows.Scan(&timeUser.macAddress, &timeUser.Second)
		if err != nil {
			log.Fatal(err)
		}

		times = append(times, timeUser)
	}

	return times
}