package alertsender

type AlertSender interface{
	SendAlert(msg string)error

}