package interactions

type Interactor interface {
	AddUser(username string) ([]byte, error)
	CreateChat(name string, users []int) ([]byte, error)
	SendMessage(chat, author int, text string) ([]byte, error)
	GetChatsOfUser(userID int) ([]byte, error)
	GetMessagesInChat(chatID int) ([]byte, error)
}
