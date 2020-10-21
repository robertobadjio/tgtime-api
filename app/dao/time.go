package dao

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type TimeUser struct {
	MacAddress string
	Second     int64
}

type Period struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Year      int    `json:"year"`
	BeginDate string `json:"beginDate"`
	EndDate   string `json:"endDate"`
}

type Periods struct {
	Periods []*Period `json:"periods"`
}

type Break struct {
	BeginTime int64 `json:"beginTime"`
	EndTime   int64 `json:"endTime"`
}

type PeriodUser struct {
	Period int     `json:"period"`
	Time   []*Time `json:"time"`
}

type Time struct {
	Date      string   `json:"date"`
	Total     int64    `json:"total"`
	BeginTime int64    `json:"beginTime"`
	EndTime   int64    `json:"endTime"`
	Break     []*Break `json:"breaks"`
}

var Db *sql.DB

func GetTimeDayAll(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["id"]
	date := mux.Vars(r)["date"]

	var macAddress string
	var telegramId int64
	row := Db.QueryRow("SELECT u.mac_address, u.telegram_id FROM users u WHERE u.id = $1", userId)
	err := row.Scan(&macAddress, &telegramId)
	if err != nil {
		panic(err)
	}
	var timeOutput Time
	timeOutput.Date = date
	timeOutput.Total = AggregateDayTotalTime(getDayTimesByUser(macAddress, date))
	timeOutput.BeginTime = GetDayTime(macAddress, date, "ASC")

	err = json.NewEncoder(w).Encode(timeOutput)
	if err != nil {
		panic(err)
	}
}

func GetTimeByPeriod(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["id"]
	period := mux.Vars(r)["period"]

	var macAddress string
	var telegramId int64
	row := Db.QueryRow("SELECT u.mac_address, u.telegram_id FROM users u WHERE u.id = $1", userId)
	err := row.Scan(&macAddress, &telegramId)
	if err != nil {
		panic(err)
	}

	row = Db.QueryRow("SELECT p.id, p.name, p.year, p.begin_at, p.ended_at FROM period p WHERE p.id = $1", period)
	periodStruct := new(Period)
	err = row.Scan(&periodStruct.Id, &periodStruct.Name, &periodStruct.Year, &periodStruct.BeginDate, &periodStruct.EndDate)
	if err != nil {
		panic(err)
	}
	var response PeriodUser
	response.Period, _ = strconv.Atoi(period)

	begin, err := time.Parse(time.RFC3339, periodStruct.BeginDate)
	end, err := time.Parse(time.RFC3339, periodStruct.EndDate)
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	now := time.Now().In(moscowLocation)
	if end.After(now) {
		end = now
	}

	for curr := begin; curr.Before(end); curr = curr.AddDate(0, 0, 1) {
		timeStruct := new(Time)
		timeStruct.Date = curr.Format("2006-01-02")
		times := getDayTimesByUser(macAddress, curr.Format("2006-01-02"))
		timeStruct.Total = AggregateDayTotalTime(times)
		timeStruct.BeginTime = GetDayTime(macAddress, curr.Format("2006-01-02"), "ASC")
		timeStruct.EndTime = GetDayTime(macAddress, curr.Format("2006-01-02"), "DESC")
		timeStruct.Break = GetAllBreaksByTimes(times)

		response.Time = append(response.Time, timeStruct)
	}

	json.NewEncoder(w).Encode(response)
}

func CreateTime(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the mac address and seconds only in order to update")
	}

	var timeUser TimeUser
	err = json.Unmarshal(reqBody, &timeUser)
	if err != nil {
		panic(err)
	}

	_, err = Db.Exec("INSERT INTO time (mac_address, second) VALUES ($1, $2)", timeUser.MacAddress, timeUser.Second)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(timeUser)
	if err != nil {
		panic(err)
	}
}

func GetAllBreaksByTimes(times []*TimeUser) []*Break {
	breaks := make([]*Break, 0)
	for i, time := range times {
		if i == 0 {
			continue
		}

		breakStruct := new(Break)

		delta := time.Second - times[i-1].Second
		if delta <= 33 {
			continue
		} else if delta <= (10 * 60) { // TODO: в параметры
			continue
		} else {
			breakStruct.BeginTime = times[i-1].Second
			breakStruct.EndTime = time.Second
			breaks = append(breaks, breakStruct)
		}
	}

	return breaks
}

func AggregateDayTotalTime(times []*TimeUser) int64 {
	num := 1
	for i, time := range times {
		if i == 0 {
			continue
		}

		delta := time.Second - times[i-1].Second
		if delta <= 33 { // TODO: в параметры
			num++
		}
	}

	return int64(num * 30) // TODO: в параметры
}

func GetSecondsByBeginDate(date string) int64 {
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	t, _ := time.ParseInLocation("2006-01-02", date, moscowLocation)

	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, moscowLocation).Unix()
}

func GetSecondsByEndDate(date string) int64 {
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	t, _ := time.ParseInLocation("2006-01-02", date, moscowLocation)

	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location()).Unix()
}

func GetDayTime(macAddress string, date string, sort string) int64 {
	var beginSecond int64
	row := Db.QueryRow("SELECT t.second FROM time t WHERE t.mac_address = $1 AND t.second BETWEEN $2 AND $3 ORDER BY t.second "+sort+" LIMIT 1", macAddress, GetSecondsByBeginDate(date), GetSecondsByEndDate(date))
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
	rows, err := Db.Query("SELECT t.mac_address, t.second FROM time t WHERE t.mac_address = $1 AND t.second BETWEEN $2 AND $3 ORDER BY t.second", macAddress, GetSecondsByBeginDate(date), GetSecondsByEndDate(date))
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	times := make([]*TimeUser, 0)
	for rows.Next() {
		timeUser := new(TimeUser)
		err := rows.Scan(&timeUser.MacAddress, &timeUser.Second)
		if err != nil {
			log.Fatal(err)
		}

		times = append(times, timeUser)
	}

	return times
}
