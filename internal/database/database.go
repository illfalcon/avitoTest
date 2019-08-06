package database

import "time"

type Chat struct {
	ID        int       `json:"id"`
	Name      string    `gorm:"unique;not null" json:"name"`
	Users     []User    `gorm:"many2many:user_chats;" json:"users"`
	Messages  []Message `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID        int       `json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Chats     []Chat    `gorm:"many2many:user_chats;" json:"-"`
	Messages  []Message `gorm:"foreignkey:AuthorID" json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	ID        int       `json:"id"`
	Chat      Chat      `json:"-"`
	ChatID    int       `json:"chat"`
	Author    User      `json:"-"`
	AuthorID  int       `json:"author"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type DBService interface {
	UserService
	ChatService
}

type UserService interface {
	AddUser(username string) (int, error)
}

type ChatService interface {
	CreateChat(name string, users []int) (int, error)
	SendMessage(chat, author int, text string) (int, error)
	GetChatsOfUser(userID int) ([]Chat, error)
	GetMessagesInChat(chatID int) ([]Message, error)
}
