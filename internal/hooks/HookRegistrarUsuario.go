package hooks

import (
	"github.com/pocketbase/pocketbase/core"
)

func HookRegistrarUsuaario(app core.App) {
	app.OnRecordAfterCreateSuccess("users").BindFunc(func(e *core.RecordEvent) error {
		usuario := e.Record
		usuario.Set("public", false)
		err := app.Save(usuario)
		if err != nil {
			return err
		}
		return e.Next()
	})
}
