package hooks

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/pocketbase/pocketbase/core"
)

func HookGenerateEventoCompleto(app core.App) {
	app.OnRecordAfterCreateSuccess("evento_cabecera").BindFunc(func(e *core.RecordEvent) error {
		collection, err := app.FindCollectionByNameOrId("evento_detalle")
		if err != nil {
			app.Logger().Error(err.Error())
			return err
		}
		cabecera := e.Record
		code, err := gonanoid.Generate("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789", 5)
		cabecera.Set("Code", code)
		detalle := core.NewRecord(collection)
		detalle.Set("Cabecera", cabecera.Id)
		err = app.Save(detalle)
		if err != nil {
			return err
		}
		cabecera.Set("Detalle", detalle.Id)
		err = app.Save(cabecera)
		if err != nil {
			return err
		}
		return e.Next()
	})
}
