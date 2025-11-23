package utils

import (
	"context"
	"fmt"
	"sync"

	"firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/pocketbase/pocketbase/core"
	"google.golang.org/api/option"
)

var (
	fcmOnce       sync.Once
	fcmClient     *messaging.Client
	fcmInitErr    error
	credFilePath  = "internal/utils/cred.json"
	initLogPrefix = "[firebase-notification-service]"
)

func SendNotificationWithUser(app core.App, userId string, title string, message string) error {
	app.Logger().Info(fmt.Sprintf("%s dispatch requested userId=%s title=%q body=%q", initLogPrefix, userId, title, message))

	record, err := app.FindRecordById("users", userId)
	if err != nil {
		return err
	}
	var tokens []string
	err = record.UnmarshalJSONField("sns", &tokens)
	if err != nil {
		app.Logger().Info("Error unmarshalling tokens")
		return err
	}

	userName := record.GetString("name")
	userEmail := record.GetString("email")

	app.Logger().Info(fmt.Sprintf("%s user context userId=%s userName=%s email=%s tokens=%d tokenList=%v", initLogPrefix, userId, userName, userEmail, len(tokens), tokens))
	if len(tokens) == 0 {
		app.Logger().Info(fmt.Sprintf("%s no push tokens for userId=%s userName=%s", initLogPrefix, userId, userName))
		return nil
	}

	sent := 0
	for _, token := range tokens {
		app.Logger().Info(fmt.Sprintf("%s sending userId=%s userName=%s token=%s title=%q body=%q", initLogPrefix, userId, userName, token, title, message))
		if err := SendNotification(token, title, message); err != nil {
			app.Logger().Error(fmt.Sprintf("%s failed userId=%s userName=%s token=%s title=%q body=%q", initLogPrefix, userId, userName, token, title, message), err)
			continue
		}
		sent++
	}
	app.Logger().Info(fmt.Sprintf("%s dispatch completed userId=%s userName=%s success=%d/%d", initLogPrefix, userId, userName, sent, len(tokens)))
	return nil
}

func getFCMClient(ctx context.Context) (*messaging.Client, error) {
	fcmOnce.Do(func() {
		firebaseApp, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(credFilePath))
		if err != nil {
			fmt.Println(initLogPrefix, "firebase.NewApp failed:", err)
			fcmInitErr = err
			return
		}
		fmt.Println(initLogPrefix, "firebase.NewApp initialized successfully")

		client, err := firebaseApp.Messaging(ctx)
		if err != nil {
			fmt.Println(initLogPrefix, "firebase.Messaging client init failed:", err)
			fcmInitErr = err
			return
		}
		fmt.Println(initLogPrefix, "firebase.Messaging client initialized successfully")
		fcmClient = client
	})
	return fcmClient, fcmInitErr
}

func InitializeNotificationClient(ctx context.Context) error {
	_, err := getFCMClient(ctx)
	return err
}

func SendNotification(token string, title string, message string) error {
	ctx := context.Background()
	client, err := getFCMClient(ctx)
	if err != nil {
		fmt.Println(initLogPrefix, "messaging client unavailable:", err)
		return err
	}

	// 🔹 Log de envío
	fmt.Println(initLogPrefix, "send request token=", token, "title=", title, "body=", message)

	// 🔹 Mensaje híbrido (Notification + Data)
	msg := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  message,
		},
		Data: map[string]string{
			"title":   title,
			"body":    message,
			"mensaje": message,
			"type":    "general",
		},
		Android: &messaging.AndroidConfig{
			Priority: "high",
			Notification: &messaging.AndroidNotification{
				ChannelID: "default_channel",
				Sound:     "default",
				Tag:       "general",
			},
		},
	}

	_, err = client.Send(ctx, msg)
	if err != nil {
		fmt.Println(initLogPrefix, "send failed token=", token, "title=", title, "error:", err)
		return err
	}

	fmt.Println(initLogPrefix, "send succeeded token=", token, "title=", title)
	return nil
}
