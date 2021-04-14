package dao

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"io/ioutil"
	"net/http"
	"officetime-api/app/model"
	"strconv"
	"time"
)

type Break struct {
	BeginTime int64 `json:"beginTime"`
	EndTime   int64 `json:"endTime"`
}

type RouterResponse struct {
	Name        string `json:"name"`
	Total       int64  `json:"total"`
	Description string `json:"description"`
}

type PeriodUser struct {
	Period int     `json:"period"`
	Time   []*Time `json:"time"`
}

type Time struct {
	Date      string           `json:"date"`
	Total     int64            `json:"total"`
	BeginTime int64            `json:"beginTime"`
	EndTime   int64            `json:"endTime"`
	Break     []*Break         `json:"breaks"`
	Routers   []RouterResponse `json:"routers"`
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
	timeOutput.Total = getDayTotalSecondsByUser(macAddress, date, 0)
	timeOutput.BeginTime = GetDayTimeFromTimeTable(macAddress, date, "ASC")
	timeOutput.EndTime = GetDayTimeFromTimeTable(macAddress, date, "DESC")
	times := GetAllByDate(macAddress, date, 0)
	timeOutput.Break = GetAllBreaksByTimesOld(times)

	err = json.NewEncoder(w).Encode(timeOutput)
	if err != nil {
		panic(err)
	}
}

// TODO: перенести в model.GetTimeByPeriod()
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

	row = Db.QueryRow("SELECT id, name, year, begin_at, ended_at FROM period WHERE id = $1", period)
	periodStruct := new(model.Period)
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

	routers := model.GetAllRouters()

	for curr := begin; curr.Before(end); curr = curr.AddDate(0, 0, 1) {
		timeStruct := new(Time)
		timeStruct.Date = curr.Format("2006-01-02")
		timeStruct.Total = getDayTotalSecondsByUser(macAddress, curr.Format("2006-01-02"), 0)
		timeStruct.BeginTime = GetDayTimeFromTimeTable(macAddress, curr.Format("2006-01-02"), "ASC")
		timeStruct.EndTime = GetDayTimeFromTimeTable(macAddress, curr.Format("2006-01-02"), "DESC")
		times := GetAllByDate(macAddress, curr.Format("2006-01-02"), 0)
		timeStruct.Break = GetAllBreaksByTimesOld(times)

		// Собираем routers
		for _, router := range routers {
			if router.WorkTime {
				continue
			}
			var responseRouters []RouterResponse
			var responseRouter RouterResponse
			responseRouter.Total = getDayTotalSecondsByUser(macAddress, curr.Format("2006-01-02"), router.Id)
			responseRouter.Name = router.Name
			responseRouter.Description = router.Description
			responseRouters = append(responseRouters, responseRouter)
			timeStruct.Routers = responseRouters
			timeStruct.Total -= responseRouter.Total
		}

		response.Time = append(response.Time, timeStruct)
	}

	json.NewEncoder(w).Encode(response)
}

func CreateTime(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the mac address and seconds only in order to update") // TODO: Описание
	}

	var timeUser model.TimeUser
	err = json.Unmarshal(reqBody, &timeUser)
	if err != nil {
		panic(err)
	}

	_, err = Db.Exec("INSERT INTO time (mac_address, second, router_id) VALUES ($1, $2, $3)", timeUser.MacAddress, timeUser.Second, timeUser.RouterId)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(timeUser)
	if err != nil {
		panic(err)
	}
}

func GetAllBreaksByTimes(macAddress, date string) []*Break {
	var breaksJson string
	row := Db.QueryRow("SELECT ts.breaks FROM time_summary ts WHERE ts.mac_address = $1 AND ts.date = $2", macAddress, date)
	err := row.Scan(&breaksJson)
	var s []*Break // TODO: Переименовать переменную
	if err == sql.ErrNoRows {
		return s
	}
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal([]byte(breaksJson), &s)
	if err != nil {
		panic(err)
	}
	return s
}

func GetAllBreaksByTimesOld(times []*model.TimeUser) []*Break {
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

func GetDayTime(macAddress string, date string, sort string) int64 {
	var second int64
	// TODO: Разделить на отдельные методы
	var field string
	if "ASC" == sort {
		field = "ts.seconds_begin"
	} else {
		field = "ts.seconds_end"
	}

	row := Db.QueryRow("SELECT "+field+" FROM time_summary ts WHERE ts.mac_address = $1 AND ts.date = $2", macAddress, date)
	err := row.Scan(&second)

	if err == sql.ErrNoRows {
		return 0
	}

	if err != nil {
		panic(err)
	}

	return second
}

func GetDayTimeFromTimeTable(macAddress string, date string, sort string) int64 {
	var beginSecond int64
	row := Db.QueryRow("SELECT t.second FROM time t WHERE t.mac_address = $1 AND t.second BETWEEN $2 AND $3 ORDER BY t.second "+sort+" LIMIT 1", macAddress, model.GetSecondsByBeginDate(date), model.GetSecondsByEndDate(date))
	err := row.Scan(&beginSecond)

	if err == sql.ErrNoRows {
		return 0
	}

	if err != nil {
		panic(err)
	}

	return beginSecond
}

func getDayTotalSecondsByUser(macAddress, date string, routerId int) int64 {
	var seconds int64

	dateToday, _ := time.Parse("2006-01-02", date)
	now := time.Now()
	if dateToday.Year() == now.Year() && dateToday.Month() == now.Month() && dateToday.Day() == now.Day() {
		times := GetAllByDate(macAddress, date, routerId)
		return model.AggregateDayTotalTime(times)
	}

	row := Db.QueryRow("SELECT ts.seconds FROM time_summary ts WHERE ts.mac_address = $1 AND ts.date = $2", macAddress, date)
	err := row.Scan(&seconds)
	if err == sql.ErrNoRows {
		return 0
	}
	if err != nil {
		panic(err)
	}

	return seconds
}

func GetAllByDate(macAddress string, date string, routerId int) []*model.TimeUser {
	var args []interface{}
	args = append(args, macAddress)
	args = append(args, model.GetSecondsByBeginDate(date))
	args = append(args, model.GetSecondsByEndDate(date))

	var routerQuery string
	if routerId != 0 {
		routerQuery = " AND t.router_id = $4"
		args = append(args, routerId)
	}

	rows, err := Db.Query("SELECT t.mac_address, t.second FROM time t WHERE t.mac_address = $1 AND t.second BETWEEN $2 AND $3"+routerQuery+" ORDER BY t.second", args...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	times := make([]*model.TimeUser, 0)
	for rows.Next() {
		time := new(model.TimeUser)
		err := rows.Scan(&time.MacAddress, &time.Second)
		if err != nil {
			panic(err)
		}

		times = append(times, time)
	}

	return times
}
