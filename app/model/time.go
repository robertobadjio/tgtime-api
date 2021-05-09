package model

import 	(
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
		if delta <= 15 * 60 {
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
