package main

import (
	"context"
	"log"
	"os"
	"pocketbaseCustom/internal/api"
	"pocketbaseCustom/internal/crons"
	"pocketbaseCustom/internal/hooks"
	"pocketbaseCustom/internal/utils"
	_ "pocketbaseCustom/migrations"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

func main() {
	app := pocketbase.New()
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: isGoRun,
	})
	api.Register(app)
	hooks.Register(app)
	crons.Register(app)
	if err := utils.InitializeNotificationClient(context.Background()); err != nil {
		app.Logger().Error("failed to initialize Firebase client", err)
	}
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
