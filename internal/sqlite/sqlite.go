package sqlite

import (
	"github.com/illfalcon/avitoTest/internal/database"
	"github.com/illfalcon/avitoTest/pkg/weberrors"
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

func (s service) isUniqueUsername(username string) (bool, error) {
	var u database.User
	result := s.db.Where("username = ?", username).First(&u)
	if result.RecordNotFound() {
		return true, nil
	}
	if result.Error != nil {
		return false, errors.Wrapf(result.Error,
			"error while querying for user with username %v in isUniqueUsername()", username)
	}
	return false, errors.Wrapf(weberrors.ErrConflict, "user with username %v exists", username)

}

func (s service) isUniqueChatname(chatname string) (bool, error) {
	var c database.Chat
	result := s.db.Where("name = ?", chatname).First(&c)
	if result.RecordNotFound() {
		return true, nil
	}
	if result.Error != nil {
		return false, errors.Wrapf(result.Error,
			"error while querying for chat with chatname %v in isUniqueChatname()", c)
	}
	return false, errors.Wrapf(weberrors.ErrConflict, "chat with chatname %v exists", chatname)
}

func (s *service) AddUser(username string) (int, error) {
	if _, err := s.isUniqueUsername(username); err != nil {
		return 0, errors.Wrap(err, "error in method AddUser()")
	}
	u := database.User{Username: username}
	s.db.Create(&u)
	return u.ID, nil
}

func (s service) getUsers(ids []int) ([]database.User, error) {
	var uu []database.User
	err := s.db.Where("id IN (?)", ids).Find(&uu).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.Wrapf(weberrors.ErrNotFound, "error in getUsers() when querying for ids %+v",
				ids)
		}
		return nil, errors.Wrapf(err, "error in getUsers() when querying for ids %+v", ids)
	}
	return uu, nil
}

func (s *service) CreateChat(name string, users []int) (int, error) {
	if _, err := s.isUniqueChatname(name); err != nil {
		return 0, errors.Wrapf(err, "error in method CreateChat() when creating chat with name %v", name)
	}
	uu, err := s.getUsers(users)
	if err != nil {
		return 0, errors.Wrapf(err, "error in method CreateChat() when creating chat with name %v", name)
	}
	c := database.Chat{Name: name, Users: uu}
	s.db.Create(&c)
	return c.ID, nil
}
