package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"pocketbaseCustom/internal/utils"
)

func main() {
	token := flag.String("token", "", "FCM registration token to target")
	title := flag.String("title", "Test notification", "Notification title")
	body := flag.String("body", "Hola desde el tester de notificaciones", "Notification body")
	flag.Parse()

	if *token == "" {
		log.Fatal("flag --token es obligatoria")
	}

	ctx := context.Background()
	if err := utils.InitializeNotificationClient(ctx); err != nil {
		log.Fatalf("no se pudo inicializar el cliente FCM: %v", err)
	}

	if err := utils.SendNotification(*token, *title, *body); err != nil {
		log.Fatalf("fallo el envio de la notificacion: %v", err)
	}

	fmt.Printf("Notificacion enviada correctamente al token %s\n", *token)
}
