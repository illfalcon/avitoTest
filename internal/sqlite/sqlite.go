package sqlite

import (
	"github.com/illfalcon/avitoTest/internal/database"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	_ "github.com/mattn/go-sqlite3"
)

func StartService() (*service, error) {
	db, err := gorm.Open("sqlite3", "database.db")
	if err != nil {
		return nil, errors.Wrap(err, "error when opening connection to sqlite instance")
	}
	db.AutoMigrate(&database.Chat{}, &database.User{}, &database.Message{})
	return &service{db: db}, nil
}

type service struct {
	db *gorm.DB
}

//func (s service) checkUniqueUsername(username string) (int, error) {
//
//}
//
//func (s *service) AddUser(username string) (int, error) {
//
//}
