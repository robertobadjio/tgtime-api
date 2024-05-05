package adapter

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"officetime-api/internal/db"
	"officetime-api/internal/model/user/app/command_query"
	"officetime-api/internal/model/user/domain/user"
	"strings"
	"time"
)

type PgUserRepository struct {
	db *sql.DB
}

func NewPgUserRepository(db *sql.DB) *PgUserRepository {
	if db == nil {
		panic("missing db")
	}

	return &PgUserRepository{db: db}
}

func (prr PgUserRepository) CreateUser(_ context.Context, u *user.User) (*command_query.User, error) {
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	now := time.Now().In(moscowLocation)
	password := randomString(10)
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 14) // TODO: в сервис
	lastInsertId := 0
	err := prr.db.QueryRow("INSERT INTO users (name, email, mac_address, telegram_id, password, created_at, surname, lastname, department_id, position) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id", u.Name, u.Email, u.MacAddress, u.TelegramId, passwordHash, now.Format("2006-01-02 15:04:05"), u.Surname, u.Lastname, u.Department).Scan(&lastInsertId)

	if pgerr, ok := err.(*pq.Error); ok {
		if pgerr.Code == "23505" {
			if strings.Contains(err.Error(), "users_email_unique") {
				return nil, &user.UserAlreadyExists{Email: u.Email}
			} else if strings.Contains(err.Error(), "users_telegram_id_unique") {
				return nil, &user.TelegramAlreadyExists{TgId: u.TelegramId}
			} else if strings.Contains(err.Error(), "users_mac_address_unique") {
				return nil, &user.MacAddressAlreadyExists{MacAddress: u.MacAddress}
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return &command_query.User{
		Password: password,
		User: &user.User{
			Id:         lastInsertId,
			Name:       u.Name,
			Surname:    u.Surname,
			Lastname:   u.Lastname,
			BirthDate:  u.BirthDate,
			Email:      u.Email,
			MacAddress: u.MacAddress,
			TelegramId: u.TelegramId,
			Role:       u.Role,
			Department: u.Department,
			Position:   u.Position,
		},
	}, nil
}

func randomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (prr PgUserRepository) UpdateUser(_ context.Context, u *user.User) (*user.User, error) {
	_, err := prr.db.Exec(
		"UPDATE users SET name = $1, email = $2, mac_address = $3, telegram_id = $4 WHERE id = $5",
		u.Name, u.Email, u.MacAddress, u.TelegramId, u.Id)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (prr PgUserRepository) GetUser(_ context.Context, userId int) (*user.User, error) {
	u := new(user.User)
	row := prr.db.QueryRow("SELECT u.id, u.name, u.email, u.mac_address, u.telegram_id, u.role, u.surname, u.lastname, u.birth_date, u.position FROM users u WHERE u.id = $1", userId)
	err := row.Scan(&u.Id, &u.Name, &u.Email, &u.MacAddress, &u.TelegramId, &u.Role, &u.Surname, &u.Lastname, &u.BirthDate, &u.Position)
	if err != nil {
		return nil, &user.NotFoundUserOfId{UserId: userId}
	}

	return u, nil
}

func (prr PgUserRepository) DeleteUser(_ context.Context, userId int) error {
	_, err := prr.db.Exec("UPDATE users SET deleted = true WHERE id = $1", userId)
	if err != nil {
		return &user.ErrorDeleteUser{UserId: userId}
	}

	return nil
}

func (prr PgUserRepository) GetUsers(_ context.Context, offset, limit int) ([]*user.User, error) {
	var args []interface{}
	statusQuery := ""
	if offset != 0 && limit != 0 {
		statusQuery = " LIMIT $1 OFFSET $2"
		args = append(args, limit)
		args = append(args, offset)
	}

	rows, err := prr.db.Query("SELECT u.id, u.name, u.email, u.mac_address, u.telegram_id, u.role, u.surname, u.lastname, u.birth_date, u.position FROM users u ORDER BY u.name ASC"+statusQuery, args...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	users := make([]*user.User, 0)
	for rows.Next() {
		u := new(user.User)
		err = rows.Scan(&u.Id, &u.Name, &u.Email, &u.MacAddress, &u.TelegramId, &u.Role, &u.Surname, &u.Lastname, &u.BirthDate, &u.Position)
		if err != nil {
			log.Fatal(err)
		}

		users = append(users, u)
	}

	//var usersStruct user.Users
	//usersStruct.Users = users

	return users, nil
}

func (prr PgUserRepository) GetUserByEmail(_ context.Context, email string) (*user.User, error) {
	u := new(user.User)
	row := prr.db.QueryRow("SELECT u.id, u.name, u.email, u.mac_address, u.telegram_id, u.role, u.surname, u.lastname, u.birth_date, u.position FROM users u WHERE u.email = $1", email)
	err := row.Scan(&u.Id, &u.Name, &u.Email, &u.MacAddress, &u.TelegramId, &u.Role, &u.Surname, &u.Lastname, &u.BirthDate, &u.Position)
	if err != nil {
		return nil, &user.NotFoundUser{UserEmail: email}
	}

	return u, nil
}

func (prr PgUserRepository) GetUserByMacAddress(_ context.Context, macAddress string) (*user.User, error) {
	u := new(user.User)
	row := prr.db.QueryRow(
		"SELECT u.id, u.name, u.email, u.mac_address, u.telegram_id, u.role, u.surname, u.lastname, u.birth_date, u.position FROM users u WHERE u.mac_address = $1",
		macAddress,
	)
	err := row.Scan(&u.Id, &u.Name, &u.Email, &u.MacAddress, &u.TelegramId, &u.Role, &u.Surname, &u.Lastname, &u.BirthDate, &u.Position)
	if err != nil {
		return nil, &user.NotFoundUserByMacAddress{MacAddress: macAddress}
	}

	return u, nil
}

// GetUsersByDepartment
// Список сотрудников по отделу
func (prr PgUserRepository) GetUsersByDepartment(_ context.Context, departmentId int) ([]*user.User, error) {
	rows, err := prr.db.Query("SELECT u.id, u.name, u.email, u.mac_address, u.telegram_id, u.role, u.surname, u.lastname, u.birth_date, u.position FROM users u WHERE u.department_id = $1", departmentId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	users := make([]*user.User, 0)
	for rows.Next() {
		u := new(user.User)
		err = rows.Scan(&u.Id, &u.Name, &u.Email, &u.MacAddress, &u.TelegramId, &u.Role, &u.Surname, &u.Lastname, &u.BirthDate, &u.Position)
		if err != nil {
			log.Fatal(err)
		}

		users = append(users, u)
	}

	return users, nil
}

func (prr PgUserRepository) GetUserPasswordHashByEmail(_ context.Context, email string) (string, error) {
	var passwordHash string
	row := prr.db.QueryRow("SELECT u.password FROM users u WHERE u.email = $1", email)
	err := row.Scan(&passwordHash)
	if err != nil {
		panic(err) // TODO: Return error
	}

	return passwordHash, nil
}

func NewPgConnection() (*sql.DB, error) {
	/*config := mysql.NewConfig()

	config.Net = "tcp"
	config.Addr = os.Getenv("MYSQL_ADDR")
	config.User = os.Getenv("MYSQL_USER")
	config.Passwd = os.Getenv("MYSQL_PASSWORD")
	config.DBName = os.Getenv("MYSQL_DATABASE")
	config.ParseTime = true // with that parameter, we can use time.Time in mysqlHour.Hour

	db, err := sqlx.Connect("mysql", config.FormatDSN())
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to MySQL")
	}

	return db, nil*/

	return db.GetDB(), nil
}
