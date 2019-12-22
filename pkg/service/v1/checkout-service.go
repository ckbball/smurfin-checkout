package v1

import (
  "context"
  "database/sql"
  "fmt"
  "strconv"
  "time"

  //"github.com/golang/protobuf/ptypes"
  "github.com/go-redis/cache/v7"
  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/status"
  "context"
  "encoding/json"
  "github.com/ThreeDotsLabs/watermill"
  "github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
  "github.com/ThreeDotsLabs/watermill/message"

  catalogProto "github.com/ckbball/smurfin-catalog/pkg/api/v1"
  paymentProto "github.com/ckbball/smurfin-payment/pkg/api/v1"
  "time"

  v1 "github.com/ckbball/smurfin-checkout/pkg/api/v1"
)

const (
  apiVersion = "v1"
)

type handler struct {
  repo          repository
  catalogClient catalogProto.CatalogServiceClient
  subscriber    message.Subscriber
  publisher     message.Publisher
}

func NewCheckoutServiceServer(repo repository, catalogClient catalogProto.CatalogServiceClient,
 subscriber message.Subscriber, publisher message.Publisher) {
  return &handler{
    repo: repo, 
    catalogClient: catalogClient,
    subscriber: subscriber,
    publisher: publisher
  }
}

func (s *handler) checkAPI(api string) error {
  if len(api) > 0 {
    if apiVersion != api {
      return status.Errorf(codes.Unimplemented,
        "unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
    }
  }
  return nil
}

func (s *handler) connect(ctx context.Context) (*sql.Conn, error) {
  c, err := s.repo.Conn(ctx)
  if err != nil {
    return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
  }
  return c, nil
}

func (s *handler) Checkout(ctx context.Context, req *v1.Request) error {
  // check api version
  if err := s.checkAPI(req.Api); err != nil {
    return nil, err
  }

  // confirm item info with grpc call

  // send payment info to payment api
  /*
  card struct
  user_id
  item_id
  ? maybe more
  */

  // queue account purchased event info, full item - buyer_id - ? maybe more

  // in queue functionality write each entry to file? or some sort of persistent store

  // return
}

/* Steps for publisher worker
1. Grab data from queue
2. Build watermill message
3. Publish message
4. Create journal entry
*/

/* Queue functionality, FIFO

*/

/* Using Kafka As queue
topic: checkout.queue
in server.go start worker pool of consumers like one of watermill's examples
somethin like this https://github.com/ThreeDotsLabs/watermill/blob/master/_examples/basic/2-realtime-feed/consumer/main.go
*/

/*
Linked list 
routine grabs head
locks head
gets value from head and moves pointer to next value in queue
unlocks head
routine operates on data it grabbed.
*/