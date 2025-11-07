package alertconsumer

import (

	"alert_system/internal/alert_sender"

)


type AlertConsumer interface{
	StartConsumingLog(alerter alertsender.AlertSender) error
	Close()
}