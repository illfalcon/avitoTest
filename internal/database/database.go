package database

import "time"

type Chat struct {
	ID        int    `json:"id"`
	Name      string `json:"name" gorm:"unique;not null"`
	Users     []User `json:"users" gorm:"many2many:user_chats;"`
	Messages  []Message
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Chats     []Chat    `gorm:"many2many:user_chats;"`
	Messages  []Message `gorm:"foreignkey:AuthorID"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	ID        int  `json:"id"`
	Chat      Chat `json:"chat"`
	ChatID    int
	Author    User `json:"author"`
	AuthorID  int
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
	AddMessage(chat, author int, text string) (int, error)
	GetChatsOfUser(userID int) ([]Chat, error)
	GetMessagesInChat(chatID int) ([]Message, error)
}
