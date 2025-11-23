package api

import (
	"net/http"
	"pocketbaseCustom/internal/dto"
	_ "pocketbaseCustom/internal/dto"
	"pocketbaseCustom/internal/utils"
	"slices"

	_ "pocketbaseCustom/internal/utils"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

func VerificarEvento(e *core.RequestEvent, app core.App) (*core.Record, *core.Record, error) {
	code := e.Request.PathValue("id")
	evento_cabecera, err := app.FindFirstRecordByFilter(
		"evento_cabecera",
		"Code = {:code}",
		dbx.Params{"code": code},
	)
	if err != nil {
		return nil, nil, e.UnauthorizedError("No existe el evento", nil)
	}
	detalle, err := app.FindFirstRecordByFilter(
		"evento_detalle",
		"Cabecera.Code = {:code}",
		dbx.Params{"code": code},
	)
	if err != nil {
		return nil, nil, e.UnauthorizedError("No existe el detalle del evento", nil)
	}
	return evento_cabecera, detalle, nil
}

func EventosRoutes(se *core.ServeEvent, app core.App) {
	grp := se.Router.Group("/api/eventos")
	grp.PUT("/{id}/invite", func(e *core.RequestEvent) error {
		userId := e.Auth.Id
		cabecera, _, err := VerificarEvento(e, app)
		if err != nil {
			res := dto.FromErrorResult(err)
			return e.JSON(http.StatusOK, dto.ToMap(&res))
		}
		asistentes := cabecera.GetStringSlice("Asistentes")
		if slices.Contains(asistentes, userId) {
			res := dto.ErrorResult("Ya esta inscrito en este evento")
			return e.JSON(http.StatusOK, dto.ToMap(&res))
		}
		if cabecera.GetString("Creador") == userId {
			res := dto.ErrorResult("Ya esta inscrito en este evento, eres el organizador!!!")
			return e.JSON(http.StatusOK, dto.ToMap(&res))
		}
		nuevoParticipante, _ := app.FindRecordById("users", userId)

		//antes de agregar el asistente deberia enviar las notificaciones
		utils.SendEventPeapleNotification(app,
			cabecera.Id,
			"Nueva participante!!",
			"Preparate para la llegada de "+nuevoParticipante.GetString("name"))
		asistentes = append(asistentes, userId)
		cabecera.Set("Asistentes", asistentes)
		if err := app.Save(cabecera); err != nil {
			res := dto.ErrorResult("No se pudo inscribir")
			return e.JSON(http.StatusOK, dto.ToMap(&res))
		}
		err = utils.SendNotificationWithUser(app, userId, "Evento", "Te has inscrito en el evento "+cabecera.GetString("nombre"))
		res := dto.SucessResult("Inscrito correctamente", cabecera.Id)
		return e.JSON(http.StatusOK, dto.ToMap(&res))
	})
}
