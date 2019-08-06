package validation

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/illfalcon/avitoTest/pkg/weberrors"
	"github.com/pkg/errors"
)

type SelfValidator interface {
	Validate() error
}

func Decode(r *http.Request, v SelfValidator) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return errors.Wrapf(weberrors.ErrBadRequest, "error in Decode(): %v", err)
	}
	return v.Validate()
}

type NewUser struct {
	Username string
}

type NewChat struct {
	Name         string
	Users        []string
	DecodedUsers []int
}

type NewMessage struct {
	Chat          string
	DecodedChat   int
	Author        string
	DecodedAuthor int
	Text          string
}

type GetChats struct {
	User        string
	DecodedUser int
}

type GetMessages struct {
	Chat        string
	DecodedChat int
}

func (nu *NewUser) Validate() error {
	if len(nu.Username) == 0 {
		return errors.Wrap(weberrors.ErrBadRequest, "username cannot be empty")
	}
	return nil
}

func (nc *NewChat) Validate() error {
	if len(nc.Name) == 0 {
		return errors.Wrap(weberrors.ErrBadRequest, "chatname cannot be empty")
	}
	if len(nc.Users) == 0 {
		return errors.Wrap(weberrors.ErrBadRequest, "cannot create chat without users")
	}
	for _, u := range nc.Users {
		du, err := strconv.Atoi(u)
		if err != nil {
			return errors.Wrap(weberrors.ErrBadRequest, "malformed user id")
		}
		log.Println(du)
		nc.DecodedUsers = append(nc.DecodedUsers, du)
	}
	return nil
}

func (nm *NewMessage) Validate() error {
	if len(nm.Text) == 0 {
		return errors.Wrap(weberrors.ErrBadRequest, "message cannot be empty")
	}
	if len(nm.Chat) == 0 {
		return errors.Wrap(weberrors.ErrBadRequest, "cannot create message without chat")
	}
	dc, err := strconv.Atoi(nm.Chat)
	if err != nil {
		return errors.Wrap(weberrors.ErrBadRequest, "malformed chat id")
	}
	nm.DecodedChat = dc
	if len(nm.Author) == 0 {
		return errors.Wrap(weberrors.ErrBadRequest, "cannot create message without author")
	}
	da, err := strconv.Atoi(nm.Author)
	if err != nil {
		return errors.Wrap(weberrors.ErrBadRequest, "malformed author id")
	}
	nm.DecodedAuthor = da
	return nil
}

func (gc *GetChats) Validate() error {
	if len(gc.User) == 0 {
		return errors.Wrap(weberrors.ErrBadRequest, "cannot get chats of no user")
	}
	du, err := strconv.Atoi(gc.User)
	if err != nil {
		return errors.Wrap(weberrors.ErrBadRequest, "malformed user id")
	}
	gc.DecodedUser = du
	return nil
}

func (gm *GetMessages) Validate() error {
	if len(gm.Chat) == 0 {
		return errors.Wrap(weberrors.ErrBadRequest, "cannot get messages of no chat")
	}
	dc, err := strconv.Atoi(gm.Chat)
	if err != nil {
		return errors.Wrap(weberrors.ErrBadRequest, "malformed chat id")
	}
	gm.DecodedChat = dc
	return nil
}
