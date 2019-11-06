package main

import (
  "context"
  "fmt"
  pb "github.com/ckbball/smurfin-checkout/proto/checkout"
  "github.com/micro/go-micro"
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

  client, err := CreateClient(uri)
  if err != nil {
    log.Panic(err)
  }
  defer client.Disconnect(context.TODO())

  journalCollection := client.Database("smurfin").Collection("journal")
  repository := &JournalRepository{
    journalCollection,
  }
  catalogClient := catalogProto.NewCatalogServiceClient("smurfin.catalog.client", srv.Client())
  h := &handler{repository, catalogClient, }
  pb.RegisterCheckoutServiceHandler(srv.Server(), &handler{repository})

  if err := srv.Run(); err != nil {
    fmt.Println(err)
  }
}