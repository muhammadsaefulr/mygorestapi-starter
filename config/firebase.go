package config

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

func InitFirebaseAuthClient() *auth.Client {
	credentialsPath := GoogleAppCredentials
	if credentialsPath == "" {
		log.Fatal("Environment variable GOOGLE_APPLICATION_CREDENTIALS is not set.")
	}

	opt := option.WithCredentialsFile(credentialsPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing Firebase app: %v\n", err)
	}

	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Firebase Auth client: %v\n", err)
	}

	log.Println("Firebase Auth client initialized successfully")
	return authClient
}

func InitFirebaseMessagingClient() *messaging.Client {
	credentialsPath := GoogleAppCredentials
	if credentialsPath == "" {
		log.Fatal("Environment variable GOOGLE_APPLICATION_CREDENTIALS is not set.")
	}

	opt := option.WithCredentialsFile(credentialsPath)

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing Firebase app for Messaging: %v\n", err)
	}

	messagingClient, err := app.Messaging(context.Background())
	if err != nil {
		log.Fatalf("error getting Firebase Messaging client: %v\n", err)
	}

	log.Println("Firebase Messaging client initialized successfully")
	return messagingClient
}
