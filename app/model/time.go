package model

import (
	"database/sql"
	"log"
	"math"
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
	Id           int                      `json:"routerId"`
	Name         string                   `json:"routerName"`
	NumEmployees int                      `json:"numEmployees"`
	Periods      []map[string]interface{} `json:"periods"`
}

type StatDepartments struct {
	Departments []*StatDepartment `json:"departments"`
}

type StatDepartment struct {
	Id           int    `json:"departmentId"`
	Name         string `json:"departmentName"`
	NumEmployees int    `json:"numEmployees"`
	Total        int64  `json:"total"`
	TotalDay     int64  `json:"totalDay"`
}

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
	Weekend   bool             `json:"weekend"`
	Total     int64            `json:"total"`
	BeginTime int64            `json:"beginTime"`
	EndTime   int64            `json:"endTime"`
	Break     []*Break         `json:"breaks"`
	Routers   []RouterResponse `json:"routers"`
}

const TotalWorkingDayInSeconds = 8 * 60 * 60 // TODO: Сколько должно быть отработано в день

// GetAllTimesByPeriodsAndRouters
// Стаститика. Общее время по периоду по всем сотрудникам
// TODO: Ограничить 7 последними периодами
func GetAllTimesByPeriodsAndRouters() *StatByPeriodsAndRouters {
	routers := GetAllRouters()
	periods := GetAllPeriods()

	stat := new(StatByPeriodsAndRouters)
	for _, router := range routers {
		statRouter := new(StatRouter)
		statRouter.Id = router.Id
		statRouter.Name = router.Name
		statRouter.NumEmployees = len(GetAllUsers(0, 0).Users) // TODO: Количество сотрудников у которых есть доступ к роутеру
		for _, period := range periods {
			tempPeriod := make(map[string]interface{})
			tempPeriod["id"] = period.Id
			tempPeriod["name"] = period.Name
			tempPeriod["total"] = getTotalTimeByPeriodAndRouter(period, router.Id)
			tempPeriod["totalWorkTime"] = 20 * 8 * 60 * 60 // TODO: Общее рабочее время по периоду + обед, взять из производственного календаря

			statRouter.Periods = append(statRouter.Periods, tempPeriod)
		}
		stat.Routers = append(stat.Routers, statRouter)
	}

	return stat
}

// GetAllTimesDepartmentsByDate
// Стаститика. Общее время за день по отделам
func GetAllTimesDepartmentsByDate(date string) *StatDepartments {
	departments := GetAllDepartments()
	routers := GetAllRouters()
	data := new(StatDepartments)
	for _, department := range departments {
		item := new(StatDepartment)
		item.Id = department.Id
		item.Name = department.Name
		employees, _ := GetUsersByDepartment(department.Id)
		item.NumEmployees = len(employees)
		item.Total = 0
		item.TotalDay = TotalWorkingDayInSeconds
		for _, router := range routers {
			for _, employee := range employees {
				times := GetAllByDate(employee.MacAddress, date, router.Id)
				item.Total += AggregateDayTotalTime(times)
			}
		}
		data.Departments = append(data.Departments, item)
	}

	return data
}

func getTotalTimeByPeriodAndRouter(period *Period, routerId int) int64 {
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
			err = rows.Scan(&timeUser.MacAddress, &timeUser.Second, &timeUser.RouterId)
			if err != nil {
				log.Fatal(err)
			}

			timesUser = append(timesUser, timeUser)
		}
		totalSeconds += AggregateDayTotalTime(timesUser)
	}

	return totalSeconds
}

/*func AggregateDayTotalTime(times []*TimeUser) int64 {
	num := 1
	for i, time := range times {
		if i == 0 {
			continue
		}

		delta := time.Second - times[i-1].Second
		if delta <= 34 { // TODO: в параметры
			num++
		}
	}

	if 1 == num {
		return 0
	}

	return int64(num * 30) // TODO: в параметры
}*/

// AggregateDayTotalTime Подсчет общего количества секунд
func AggregateDayTotalTime(times []*TimeUser) int64 {
	var sum int64
	for i, time := range times {
		if i == 0 {
			continue
		}
		delta := time.Second - times[i-1].Second
		// Не учитываем перерывы меньше 15 минут
		if delta <= 15*60 {
			sum += delta
		}
	}

	return sum
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

func GetAllByDate(macAddress string, date string, routerId int) []*TimeUser {
	var args []interface{}
	args = append(args, macAddress)
	args = append(args, GetSecondsByBeginDate(date))
	args = append(args, GetSecondsByEndDate(date))

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

	times := make([]*TimeUser, 0)
	for rows.Next() {
		time := new(TimeUser)
		err = rows.Scan(&time.MacAddress, &time.Second)
		if err != nil {
			panic(err)
		}

		times = append(times, time)
	}

	return times
}

func GetTimeByPeriod(userId, period int) PeriodUser {
	var macAddress string
	var telegramId int64
	row := Db.QueryRow("SELECT u.mac_address, u.telegram_id FROM users u WHERE u.id = $1", userId)
	err := row.Scan(&macAddress, &telegramId)
	if err != nil {
		panic(err)
	}

	row = Db.QueryRow("SELECT id, name, year, begin_at, ended_at FROM period WHERE id = $1", period)
	periodStruct := new(Period)
	err = row.Scan(&periodStruct.Id, &periodStruct.Name, &periodStruct.Year, &periodStruct.BeginDate, &periodStruct.EndDate)
	if err != nil {
		panic(err)
	}
	var response PeriodUser
	response.Period = period

	begin, err := time.Parse(time.RFC3339, periodStruct.BeginDate)
	end, err := time.Parse(time.RFC3339, periodStruct.EndDate)
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	now := time.Now().In(moscowLocation)
	if end.After(now) {
		end = now
	}

	routers := GetAllRouters()
	weekend := GetWeekendByPeriod(begin, end)
	for curr := begin; curr.Before(end); curr = curr.AddDate(0, 0, 1) {
		timeStruct := new(Time)
		timeStruct.Date = curr.Format("2006-01-02")
		// TODO: В метод
		if _, ok := weekend[curr.Format("2006-01-02")]; ok {
			timeStruct.Weekend = true
		} else {
			timeStruct.Weekend = false
		}

		timeStruct.Total = GetDayTotalSecondsByUser(macAddress, curr.Format("2006-01-02"), 0)
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
			responseRouter.Total = GetDayTotalSecondsByUser(macAddress, curr.Format("2006-01-02"), router.Id)
			responseRouter.Name = router.Name
			responseRouter.Description = router.Description
			responseRouters = append(responseRouters, responseRouter)
			timeStruct.Routers = responseRouters
			timeStruct.Total -= responseRouter.Total
		}

		response.Time = append(response.Time, timeStruct)
	}

	return response
}

func GetTimeDayAll(userId int, date string) Time {
	var macAddress string
	var telegramId int64
	row := Db.QueryRow("SELECT u.mac_address, u.telegram_id FROM users u WHERE u.id = $1", userId)
	err := row.Scan(&macAddress, &telegramId)
	if err != nil {
		panic(err)
	}
	var timeOutput Time
	timeOutput.Date = date
	timeOutput.Total = GetDayTotalSecondsByUser(macAddress, date, 0)
	timeOutput.BeginTime = GetDayTimeFromTimeTable(macAddress, date, "ASC")
	timeOutput.EndTime = GetDayTimeFromTimeTable(macAddress, date, "DESC")
	times := GetAllByDate(macAddress, date, 0)
	timeOutput.Break = GetAllBreaksByTimesOld(times)

	return timeOutput
}

func GetDayTotalSecondsByUser(macAddress, date string, routerId int) int64 {
	var seconds int64

	dateToday, _ := time.Parse("2006-01-02", date)
	now := time.Now()
	if dateToday.Year() == now.Year() && dateToday.Month() == now.Month() && dateToday.Day() == now.Day() {
		times := GetAllByDate(macAddress, date, routerId)
		return AggregateDayTotalTime(times)
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

func GetDayTimeFromTimeTable(macAddress string, date string, sort string) int64 {
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

func GetAllBreaksByTimesOld(times []*TimeUser) []*Break {
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

type statWorkingPeriod struct {
	StartWorkingDate  string `json:"start_working_date"`
	EndWorkingDate    string `json:"end_working_date"`
	WorkingHours      int64  `json:"working_hours"`
	TotalWorkingHours int    `json:"total_working_hours"`
}

func GetStatWorkingPeriod(userId, periodId int) *statWorkingPeriod {
	period, _ := GetPeriodById(periodId)

	periodUser := GetTimeByPeriod(userId, periodId)
	var totalMonthWorkingTime int64
	for _, timeResponse := range periodUser.Time {
		totalMonthWorkingTime += timeResponse.Total
	}

	start, err := time.Parse(time.RFC3339, period.BeginDate)
	if err != nil {
		panic(err)
	}
	end, err := time.Parse(time.RFC3339, period.EndDate)
	if err != nil {
		panic(err)
	}

	return &statWorkingPeriod{
		StartWorkingDate:  start.Format("02.01.2006"),
		EndWorkingDate:    end.Format("02.01.2006"),
		WorkingHours:      totalMonthWorkingTime / 3600,
		TotalWorkingHours: GetWorkHoursBetween(start, GetNow()),
	}
}

func GetNow() time.Time {
	return time.Now().In(GetMoscowLocation())
}

func GetMoscowLocation() *time.Location {
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	return moscowLocation
}

// getWeekdaysBetween
// https://switch-case.ru/61590709
func getWeekdaysBetween(start, end time.Time) int {
	offset := -int(start.Weekday())
	start = start.AddDate(0, 0, -int(start.Weekday()))

	offset += int(end.Weekday())
	if end.Weekday() == time.Sunday {
		offset++
	}
	end = end.AddDate(0, 0, -int(end.Weekday()))

	dif := end.Sub(start).Truncate(time.Hour * 24)
	weeks := (dif.Hours() / 24) / 7

	return int(math.Round(weeks)*5) + offset
}

func GetWorkHoursBetween(start, end time.Time) int {
	return getWeekdaysBetween(start, end) * 8
}
