package v1

import (
  //"context"
  "github.com/Shopify/sarama"
  "github.com/ThreeDotsLabs/watermill"
  "github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
  //"github.com/ThreeDotsLabs/watermill/message"
  //"log"
  //"time"
)

func InitSubscriber(config *sarama.Config) *kafka.Subscriber {
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
