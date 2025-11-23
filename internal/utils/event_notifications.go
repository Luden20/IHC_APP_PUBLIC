package utils

import (
	_ "fmt"

	"github.com/pocketbase/pocketbase/core"
)

func SendEventPeapleNotification(app core.App, id_cabecera string, title string, message string) {
	eventoCabecera, _ := app.FindRecordById("evento_cabecera", id_cabecera)
	asistentes := eventoCabecera.GetStringSlice("Asistentes")
	for _, asistente := range asistentes {
		_ = SendNotificationWithUser(app, asistente, title, message)
	}
	_ = SendNotificationWithUser(app, eventoCabecera.GetString("Creador"), title, message)
}
