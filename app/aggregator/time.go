package aggregator

import (
	"context"
	"database/sql"
	"encoding/json"
	"officetime-api/app/model"
	"officetime-api/internal/db"
	userAdapter "officetime-api/internal/model/user/adapter"
	userApp "officetime-api/internal/model/user/app"
	userCommand "officetime-api/internal/model/user/app/command"
	userCommandQuery "officetime-api/internal/model/user/app/command_query"
	userQuery "officetime-api/internal/model/user/app/query"
	"time"
)

var Db *sql.DB

func AggregateTime() {
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	date := time.Now().AddDate(0, 0, -1).In(moscowLocation)

	uApp := buildUserApp()
	qr := userQuery.GetUsers{Offset: 0, Limit: 0}
	ctx := context.TODO()
	users, _ := uApp.Queries.GetUsers.Handle(ctx, qr) // TODO: Handle error

	for _, user := range users {
		times := model.GetAllByDate(user.MacAddress, date, 0)
		seconds := model.AggregateDayTotalTime(times)
		breaks := model.GetAllBreaksByTimesOld(times)
		breaksJson, err := json.Marshal(breaks)
		begin := model.GetDayTimeFromTimeTable(user.MacAddress, date, "ASC")
		end := model.GetDayTimeFromTimeTable(user.MacAddress, date, "DESC")

		_, err = Db.Exec("INSERT INTO time_summary (mac_address, date, seconds, breaks, seconds_begin, seconds_end) VALUES ($1, $2, $3, $4, $5, $6)", user.MacAddress, date, seconds, breaksJson, begin, end)
		if err != nil {
			panic(err)
		}
	}
}

func buildUserApp() userApp.Application {
	userRepository := userAdapter.NewPgUserRepository(db.GetDB())
	return userApp.Application{
		Commands: userApp.Commands{
			UpdateUser: userCommand.NewUpdateUserHandler(userRepository),
			DeleteUser: userCommand.NewDeleteUserHandler(userRepository),
		},
		Queries: userApp.Queries{
			GetUser:  userQuery.NewGetUserHandler(userRepository),
			GetUsers: userQuery.NewGetUsersHandler(userRepository),
		},
		CommandsQueries: userApp.CommandsQueries{
			CreateUser: userCommandQuery.NewCreateUserHandler(userRepository),
		},
	}
}
