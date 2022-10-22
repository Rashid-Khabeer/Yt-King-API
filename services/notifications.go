package services

import (
	"context"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

type Notification interface {
	SendNotificationToTopic(title string, body string, topic string)
}

func NewNotificationService() *notificationService {
	path, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}
	cred := option.WithCredentialsFile(path + "/admin-firebase.json")
	app, err := firebase.NewApp(context.Background(), nil, cred)
	if err != nil {
		panic(err)
	}

	msgClient, err := app.Messaging(context.TODO())
	if err != nil {
		panic(err)
	}
	return &notificationService{messagingClient: msgClient}
}

func (n *notificationService) SendNotificationToTopic(title string, body string, topic string) {
	message := prepareMessage(title, body)
	message.Topic = topic
	send, err := n.messagingClient.Send(context.TODO(), message)
	if err != nil {
		panic(err)
	}
	fmt.Println(send)
}

func prepareMessage(title string, body string) *messaging.Message {
	message := &messaging.Message{}
	message.APNS = nil
	message.Android = &messaging.AndroidConfig{
		Notification: &messaging.AndroidNotification{
			Title:        title,
			Body:         body,
			Priority:     messaging.PriorityHigh,
			DefaultSound: true,
		},
	}

	return message
}

type notificationService struct {
	messagingClient *messaging.Client
}
