package message_sender

type MessageSender interface {
	SendMessage(msg string)error
}