package message_receiver

type MessageReceiver interface {
	ReceiveMessage()error
}