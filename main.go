package main

import (
  "context"
  "fmt"
  pb "github.com/ckbball/smurfin-checkout/proto/checkout"
  "github.com/micro/go-micro"
  "github.com/Shopify/sarama"
  "log"
  "os"
)

const (
  defaultHost = "datastore:27017"
)

func main() {
  srv := micro.NewService(
    micro.Name("smurfin.checkout")
  )

  srv.Init()

  uri := os.Getenv("DB_HOST")
  if uri == "" {
    uri = defaultHost
  }

  // Will need to change for MySQL
  client, err := CreateClient(uri)
  if err != nil {
    log.Panic(err)
  }
  defer client.Disconnect(context.TODO())

  journalCollection := client.Database("smurfin-checkout").Collection("journal")
  repository := &JournalRepository{
    journalCollection,
  }
  catalogClient := catalogProto.NewCatalogServiceClient("smurfin.catalog.client", srv.Client())

  // Make subscriber config here
  saramaSubscriberConfig := kafka.DefaultSaramaSubscriberConfig()
  saramaSubscriberConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

  // Make subscriber pointer here
  subscriber := InitSubscriber(saramaSubscriberConfig)

  // Make publisher pointer here
  publisher := InitPublisher()

  // Make handler here with stuff
  h := &handler{repository, catalogClient, subscriber, publisher}

  // Register handler and server
  pb.RegisterCheckoutServiceHandler(srv.Server(), h)

  r := InitWaterRouter(publisher, subscriber, repository)

  fmt.Println("Spinning up router")
  go r.Run()
  <- r.Running()
  fmt.Println("Router is running")
  // to check router is running

  // Run Server
  if err := srv.Run(); err != nil {
    fmt.Println(err)
  }
}