package hooks

import (
	"pocketbaseCustom/internal/utils"
	_ "pocketbaseCustom/internal/utils"

	_ "fmt"
	"strconv"

	"github.com/pocketbase/pocketbase/core"
)

func HookEventDetalleUpdate(app core.App) {
	app.OnRecordAfterUpdateSuccess("evento_cabecera").BindFunc(func(e *core.RecordEvent) error {
		old := e.Record.Original()
		if old == nil {
			return e.Next()
		}
		newRecord := e.Record
		app.Logger().Info("Comparando estado de evento")
		if old.GetBool("Activo") == false && newRecord.GetBool("Activo") == true {
			utils.SendEventPeapleNotification(app,
				newRecord.Id,
				"! Hora de la fiesta "+newRecord.GetString("Titulo")+"!",
				"Es el momento, puedes empezar a tomar tus fotos!!!")
		}
		if old.GetBool("Activo") == true && newRecord.GetBool("Activo") == false {
			utils.SendEventPeapleNotification(app,
				newRecord.Id,
				"! Se ha terminado la fiesta "+newRecord.GetString("Titulo")+"!",
				"La fiesta a terminado, pero puedes recordar todos los momentos con las fotos !!!")
		}
		return e.Next()
	})
	app.OnRecordAfterUpdateSuccess("evento_detalle").BindFunc(func(e *core.RecordEvent) error {
		if e == nil || e.Record == nil {
			return e.Next() // o return nil
		}

		old := e.Record.Original()
		if old == nil {
			return e.Next()
		}
		oldSlice := old.GetStringSlice("Fotos")
		oldLength := len(oldSlice)
		newSlice := e.Record.GetStringSlice("Fotos")
		newLength := len(newSlice)
		if oldLength < newLength {
			var diff = int64(newLength - oldLength)
			strDiff := strconv.FormatInt(diff, 10)
			eventoCabecera, _ := app.FindRecordById("evento_cabecera", e.Record.GetString("Cabecera"))
			utils.SendEventPeapleNotification(app, e.Record.GetString("Cabecera"),
				"Nueva foto!!!",
				"Alguien a agregado "+strDiff+" foto(s) en "+eventoCabecera.GetString("Titulo"))
		}
		return e.Next()
	})
}
