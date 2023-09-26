package centrifugo

import (
	"context"
	"fmt"
	"log"
	"notification/configs"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	centrifuge "github.com/centrifugal/centrifuge-go"
)

var (
	centrifugoAddr = os.Getenv("CENTRIFUGO_ADDR")
	centrifugoKey  = os.Getenv("CENTRIFUGO_KEY")
	channel        string
	action         string
	userID         int
	tokenString    string
)

type service struct {
	Client     *centrifuge.Client
	Subscribed *centrifuge.Subscription
}

type Service interface {
	Broadcast(ch string)
}

func NewService(ch, act string, user_id int) *service {
	jwtKey := configs.Jwt().SecretKey()
	channel = ch
	action = act
	userID = user_id
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": fmt.Sprintf("%d", userID),
		"exp": time.Now().Add(time.Minute * 3).Unix(),
		// "channel": fmt.Sprintf("%s:%s#%d", channel, action, userID),
	})
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		log.Println("error generating token. Error: ", err.Error())
		return nil
	}

	client := centrifuge.NewJsonClient(
		fmt.Sprintf("ws://%s/connection/websocket", centrifugoAddr),
		centrifuge.Config{
			EnableCompression: true,
			Token:             tokenString,
			// Uncomment to make it work with Centrifugo JWT auth.
			// Token: "apikey my_api_key",
			//Token: connToken("49", 0),
		},
	)
	err = client.Connect()
	if err != nil {
		panic(err)
	}
	client.OnConnecting(func(e centrifuge.ConnectingEvent) {
		log.Printf("Connecting - %d (%s)", e.Code, e.Reason)
	})
	client.OnConnected(func(e centrifuge.ConnectedEvent) {
		log.Printf("Connected with ID %s", e.ClientID)
		log.Printf("Data: %b", e.Data)
	})
	client.OnDisconnected(func(e centrifuge.DisconnectedEvent) {
		log.Printf("Disconnected: %d (%s)", e.Code, e.Reason)
	})

	client.OnError(func(e centrifuge.ErrorEvent) {
		log.Printf("Error: %s", e.Error.Error())
	})

	client.OnMessage(func(e centrifuge.MessageEvent) {
		log.Printf("Message from server: %s", string(e.Data))
	})

	client.OnSubscribed(func(e centrifuge.ServerSubscribedEvent) {
		log.Printf("Subscribed to server-side channel %s: (was recovering: %v, recovered: %v)", e.Channel, e.WasRecovering, e.Recovered)
	})
	client.OnSubscribing(func(e centrifuge.ServerSubscribingEvent) {
		log.Printf("Subscribing to server-side channel %s", e.Channel)
	})
	client.OnUnsubscribed(func(e centrifuge.ServerUnsubscribedEvent) {
		log.Printf("Unsubscribed from server-side channel %s", e.Channel)
	})

	client.OnPublication(func(e centrifuge.ServerPublicationEvent) {
		log.Printf("Publication from server-side channel %s: %s (offset %d)", e.Channel, e.Data, e.Offset)
	})
	client.OnJoin(func(e centrifuge.ServerJoinEvent) {
		log.Printf("Join to server-side channel %s: %s (%s)", e.Channel, e.User, e.Client)
	})
	client.OnLeave(func(e centrifuge.ServerLeaveEvent) {
		log.Printf("Leave from server-side channel %s: %s (%s)", e.Channel, e.User, e.Client)
	})
	return &service{
		Client: client,
	}
}

func (s *service) Subs() *service {
	jwtKey := configs.Jwt().SecretKey()
	log.Println(jwtKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":     fmt.Sprintf("%d", userID),
		"channel": fmt.Sprintf("%s:%s#%d", channel, action, userID),
		"exp":     time.Now().Add(time.Minute * 3).Unix(),
	})
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		log.Println("error generating token. Error: ", err.Error())
		return nil
	}

	sub, err := s.Client.NewSubscription(fmt.Sprintf("%s:%s#%d", channel, action, userID), centrifuge.SubscriptionConfig{
		Recoverable: true,
		JoinLeave:   true,
		Token:       tokenString,
	})
	if err != nil {
		log.Println(err)
		return nil
	}
	sub.OnSubscribing(func(e centrifuge.SubscribingEvent) {
		log.Printf("Subscribing on channel %s - %d (%s)", sub.Channel, e.Code, e.Reason)
	})
	sub.OnSubscribed(func(e centrifuge.SubscribedEvent) {
		log.Printf("Subscribed on channel %s, (was recovering: %v, recovered: %v)", sub.Channel, e.WasRecovering, e.Recovered)
	})
	sub.OnUnsubscribed(func(e centrifuge.UnsubscribedEvent) {
		log.Printf("Unsubscribed from channel %s - %d (%s)", sub.Channel, e.Code, e.Reason)
	})

	sub.OnError(func(e centrifuge.SubscriptionErrorEvent) {
		log.Printf("Subscription error %s: %s", sub.Channel, e.Error)
	})

	sub.OnPublication(func(e centrifuge.PublicationEvent) {
		log.Printf("Someone says via channel %s: %s (offset %d)", sub.Channel, e.Data, e.Offset)
	})
	sub.OnJoin(func(e centrifuge.JoinEvent) {
		log.Printf("Someone joined %s: user id %s, client id %s", sub.Channel, e.User, e.Client)
	})
	sub.OnLeave(func(e centrifuge.LeaveEvent) {
		log.Printf("Someone left %s: user id %s, client id %s", sub.Channel, e.User, e.Client)
	})
	return &service{
		Client:     s.Client,
		Subscribed: sub,
	}
}

func (s *service) BroadcastMessage(data []byte) {
	defer s.Client.Close()
	if err := s.Subscribed.Subscribe(); err != nil {
		log.Println(err)
		return
	}
	_, err := s.Subscribed.Publish(context.Background(), data)
	if err != nil {
		log.Println(err)
		return
	}
	return
}
