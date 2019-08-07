package sqlite

import (
	"sort"

	"github.com/illfalcon/avitoTest/internal/database"
	"github.com/illfalcon/avitoTest/pkg/weberrors"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	_ "github.com/mattn/go-sqlite3"
)

func StartService() (*service, func() error, error) {
	db, err := gorm.Open("sqlite3", "database.db")
	if err != nil {
		return nil, nil, errors.Wrap(err, "error when opening connection to sqlite instance")
	}
	db.AutoMigrate(&database.Chat{}, &database.User{}, &database.Message{})
	cl := db.Close
	return &service{db: db}, cl, nil
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
	err := s.db.Where(ids).Find(&uu).Error
	if err != nil {
		//if gorm.IsRecordNotFoundError(err) {
		//	return nil, errors.Wrapf(weberrors.ErrNotFound, "error in getUsers() when querying for ids %+v",
		//		ids)
		//}
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
	if len(uu) == 0 {
		return 0, errors.Wrapf(weberrors.ErrNotFound, "not found any users, will not create chat")
	}
	c := database.Chat{Name: name}
	c.Users = append(c.Users, uu...)
	s.db.Create(&c)
	return c.ID, nil
}

func (s service) getChatByID(chatID int) (database.Chat, error) {
	var c database.Chat
	err := s.db.Where("id = ?", chatID).First(&c).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return c, errors.Wrapf(weberrors.ErrNotFound, "error when querying for chat with id %v", chatID)
		}
		return c, errors.Wrapf(err, "error when querying for chat with id %v", chatID)
	}
	return c, nil
}

func (s service) getUserByID(userID int) (database.User, error) {
	var u database.User
	err := s.db.Where("id = ?", userID).First(&u).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return u, errors.Wrapf(weberrors.ErrNotFound, "error when querying for user with id %v", userID)
		}
		return u, errors.Wrapf(err, "error when querying for user with id %v", userID)
	}
	return u, nil
}

func (s service) chatContainsUser(c database.Chat, u database.User) bool {
	for _, user := range c.Users {
		if user.ID == u.ID {
			return true
		}
	}
	return false
}

func (s *service) SendMessage(chat, author int, text string) (int, error) {
	c, err := s.getChatByID(chat)
	s.db.Model(&c).Related(&c.Users, "Users")
	if err != nil {
		return 0, errors.Wrapf(err, "error when sending message (SendMessage()) from user %v to chat %v",
			author, chat)
	}
	u, err := s.getUserByID(author)
	if err != nil {
		return 0, errors.Wrapf(err, "error when sending message (SendMessage()) from user %v to chat %v",
			author, chat)
	}
	if !s.chatContainsUser(c, u) {
		return 0, errors.Wrapf(weberrors.ErrForbidden, "unable to send to chat from external user, chat %v, user %v",
			chat, author)
	}
	m := database.Message{Chat: c, Author: u, Text: text}
	s.db.Create(&m)
	return m.ID, nil
}

func (s service) GetChatsOfUser(userID int) ([]database.Chat, error) {
	u, err := s.getUserByID(userID)
	if err != nil {
		return nil, errors.Wrapf(err, "error when getting chats of user %v", userID)
	}
	s.db.Model(&u).Preload("Users").Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("messages.created_at desc")
	}).Related(&u.Chats, "Chats")
	sort.Slice(u.Chats, func(i, j int) bool {
		if len(u.Chats[i].Messages) == 0 {
			return len(u.Chats[j].Messages) == 0
		}
		if len(u.Chats[j].Messages) == 0 {
			return true
		}
		return u.Chats[i].Messages[0].CreatedAt.After(u.Chats[j].Messages[0].CreatedAt)
	})
	return u.Chats, nil
}

func (s service) GetMessagesInChat(chatID int) ([]database.Message, error) {
	c, err := s.getChatByID(chatID)
	if err != nil {
		return nil, errors.Wrapf(err, "error when getting messages of chat %v", chatID)
	}
	s.db.Model(&c).Order("messages.created_at asc").Related(&c.Messages)
	return c.Messages, nil
}
