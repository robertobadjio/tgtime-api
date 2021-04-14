package model

import (
	"database/sql"
	"log"
	"time"
)

var Db *sql.DB

type TimeUser struct {
	MacAddress string
	Second     int64
	RouterId   int8
}

type TotalTimeByPeriod struct {
	PeriodId int   `json:"periodId"`
	RouterId int   `json:"routerId"`
	Total    int64 `json:"total"`
}

type StatByPeriodsAndRouters struct {
	Routers []*StatRouter `json:"routers"`
}

type StatRouter struct {
	Id int `json:"routerId"`
	Name string `json:"routerName"`
	Periods []map[string]interface{} `json:"periods"`
}

// GetAllTimesByPeriodsAndRouters
// Общее время по периоду по всем сотрудникам
// TODO: Ограничить 7 последними периодами
func GetAllTimesByPeriodsAndRouters() *StatByPeriodsAndRouters {
	routers := GetAllRouters()
	periods := GetAllPeriods()

	stat := new(StatByPeriodsAndRouters)
	for _, router := range routers {
		statRouter := new(StatRouter)
		statRouter.Id = router.Id
		statRouter.Name = router.Name
		for _, period := range periods {
			tempPeriod := make(map[string]interface{})
			tempPeriod["id"] = period.Id
			tempPeriod["name"] = period.Name
			tempPeriod["total"] = getTotalTimeByPeriodAndRouter(period, router.Id)

			statRouter.Periods = append(statRouter.Periods, tempPeriod)
		}
		stat.Routers = append(stat.Routers, statRouter)
	}

	return stat
}

func getTotalTimeByPeriodAndRouter(period *Period, routerId int) int64  {
	var totalSeconds int64
	users := GetAllUsers(0, 0)
	for _, user := range users.Users {
		rows, err := Db.Query("SELECT t.mac_address, t.second, t.router_id FROM time t WHERE t.router_id = $1 AND t.second BETWEEN $2 AND $3 AND t.mac_address = $4", routerId, GetSecondsByBeginDate(period.BeginDate), GetSecondsByEndDate(period.EndDate), user.MacAddress)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		var timesUser []*TimeUser
		for rows.Next() {
			timeUser := new(TimeUser)
			err := rows.Scan(&timeUser.MacAddress, &timeUser.Second, &timeUser.RouterId)
			if err != nil {
				log.Fatal(err)
			}

			timesUser = append(timesUser, timeUser)
		}
		totalSeconds += AggregateDayTotalTime(timesUser)
	}

	return totalSeconds
}

func AggregateDayTotalTime(times []*TimeUser) int64 {
	num := 1
	for i, time := range times {
		if i == 0 {
			continue
		}

		delta := time.Second - times[i-1].Second
		if delta <= 31 { // TODO: в параметры
			num++
		}
	}

	if 1 == num {
		return 0
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
