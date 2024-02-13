package db

import (
	"announce-api/entities"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	Driver   string
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
	SSLMode  string
}

type DatabaseClient struct {
	Auth
	AnnouncementActions
}

type Auth interface {
	CreateUser(user *entities.InputUser) error
	AuthorizeUser()
	generateUserAccessToken()
}

type AnnouncementActions interface {
	GetList()
	GetOneById()
	CreateAnnounce()
	UpdateAnnounce()
	HideAnnounce()
	DeleteAnnounce()
}

func newClient(db *sqlx.DB) *DatabaseClient {
	return &DatabaseClient{
		Auth:                &Authenticator{db},
		AnnouncementActions: &AnnouncementManager{db},
	}
}

func Connect(config Config) (*DatabaseClient, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", config.Host, config.Port, config.DBName, config.User, config.Password, config.SSLMode)

	db, err := sqlx.Connect(config.Driver, dataSourceName)
	if err != nil {
		return nil, err
	}

	client := newClient(db)

	db.MustExec(schema)

	fmt.Println("Database connection successfull")
	return client, nil
}

// func (c *DatabaseClient) CreateQueryString() string {
// 	cfg := Config{
// 		Driver: "postgres",
// 		Host:   "localhost",
// 	}

// 	r := reflect.ValueOf(&cfg).Elem()
// 	rt := r.Type()
// 	for i := 0; i < rt.NumField(); i++ {
// 		field := rt.Field(i)
// 		rv := reflect.ValueOf(&cfg)
// 		value := reflect.Indirect(rv).FieldByName(field.Name)
// 		fmt.Println(field.Name, value.String())
// 	}

// 	return ""
// }
