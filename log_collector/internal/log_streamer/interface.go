package log_streamer

type LogProducer interface {
	SendLog(msg string)error
	Close()error
}