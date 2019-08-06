package interactions

import (
	"encoding/json"

	"github.com/illfalcon/avitoTest/internal/database"
	"github.com/pkg/errors"
)

type MyInteractor struct {
	db database.DBService
}

func NewMyInteractor(service database.DBService) MyInteractor {
	return MyInteractor{db: service}
}

func (mi *MyInteractor) AddUser(username string) ([]byte, error) {
	id, err := mi.db.AddUser(username)
	if err != nil {
		return nil, errors.Wrapf(err, "error when adding user %v in interactions.AddUser()", username)
	}
	r := map[string]int{
		"id": id,
	}
	resp, err := json.Marshal(r)
	if err != nil {
		return nil, errors.Wrap(err, "error when marsalling response in interactions.AddUser()")
	}
	return resp, nil
}

func (mi *MyInteractor) CreateChat(name string, users []int) ([]byte, error) {
	id, err := mi.db.CreateChat(name, users)
	if err != nil {
		return nil, errors.Wrapf(err, "error when adding chat %v in interactions.CreateChat()", name)
	}
	r := map[string]int{
		"id": id,
	}
	resp, err := json.Marshal(r)
	if err != nil {
		return nil, errors.Wrap(err, "error when marsalling response in interactions.CreateChat()")
	}
	return resp, nil
}

func (mi *MyInteractor) SendMessage(chat, author int, text string) ([]byte, error) {
	id, err := mi.db.SendMessage(chat, author, text)
	if err != nil {
		return nil, errors.Wrapf(err, "error when sending message to chat %v from user %v in interactions.SendMessage()",
			chat, author)
	}
	r := map[string]int{
		"id": id,
	}
	resp, err := json.Marshal(r)
	if err != nil {
		return nil, errors.Wrap(err, "error when marsalling response in interactions.SendMessage()")
	}
	return resp, nil
}

func (mi MyInteractor) GetChatsOfUser(userID int) ([]byte, error) {
	cc, err := mi.db.GetChatsOfUser(userID)
	if err != nil {
		return nil, errors.Wrapf(err, "error when getting chats of user %v in interactions", userID)
	}
	resp, err := json.Marshal(cc)
	if err != nil {
		return nil, errors.Wrap(err, "error when marsalling response in interactions.GetChatsOfUser()")
	}
	return resp, nil
}

func (mi MyInteractor) GetMessagesInChat(chatID int) ([]byte, error) {
	mm, err := mi.db.GetMessagesInChat(chatID)
	if err != nil {
		return nil, errors.Wrapf(err, "error when getting messages of chat %v in interactions", chatID)
	}
	resp, err := json.Marshal(mm)
	if err != nil {
		return nil, errors.Wrap(err, "error when marsalling response in interactions.GetMessagesInChat()")
	}
	return resp, nil
}
