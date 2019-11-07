package main

import (
  "context"
  "github.com/ThreeDotsLabs/watermill"
  "github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
  "github.com/ThreeDotsLabs/watermill/message"
  "log"
  "time"
)

type PaymentInfo struct {
  Card_num int
  Date_M   int
  Date_Y   int
  Code     int
  First    string
  Last     string
  Zip      int
}

type ValidatePaymentEvent struct {
  BuyerId       string
  Info          PaymentInfo
  AmountDollars int
  AmountCents   int
  AccountId     string
}

type EmailAccountEvent struct {
  BuyerId              string
  AccountLogin         string
  AccountPassword      string
  AccountEmail         string
  AccountEmailPassword string
}

type PaymentValidatedEvent struct {
  BuyerId   string
  AccountId string
}

type PaymentSuccessEvent struct {
  BuyerId   string
  AccountId string
}

type RemoveItemEvent struct {
  BuyerId   string
  AccountId string
}

func InitSubscriber(config kafka.SubscriberConfig) *kafka.Subscriber {
  subscriber, err := kafka.NewSubscriber(
    kafka.SubscriberConfig{
      Brokers:               []string{"kafka:9092"},
      Unmarshaler:           kafka.DefaultMarshaler{},
      OverwriteSaramaConfig: config,
      ConsumerGroup:         "test_consumer_group",
    },
    watermill.NewStdLogger(false, false),
  )
  if err != nil {
    panic(err)
  }
  return subscriber
}

func InitPublisher() *kafka.Publisher {
  publisher, err := kafka.NewPublisher(
    kafka.PublisherConfig{
      Brokers:   []string{"kafka:9092"},
      Marshaler: kafka.DefaultMarshaler{},
    },
    watermill.NewStdLogger(false, false),
  )
  if err != nil {
    panic(err)
  }
  return publisher
}
