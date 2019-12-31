package v1

import (
  "context"
  "database/sql"
  "fmt"
  "strconv"
  "time"
  "log"

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
  eventName = "account_purchased"
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

/* Checkout handles api calls to grpc method Checkout and REST endpoint: /v1/checkout
  Checkout is the process by which a user purchases an account; sending account and card info
  Input:
  v1.Request{
    fill in later
  }
  Output:
  - Makes payment api call for the account
  - Queues AccountPurchased event to be later published for other services
  - returns any error generated or nil if no errors.
*/
func (s *handler) Checkout(ctx context.Context, req *v1.Request) (*v1.Response, error) {
  // check api version
  if err := s.checkAPI(req.Api); err != nil {
    return err
  }

  // confirm item info with grpc call
  // will need to change to full get in future
  catalogResponse, err := s.catalogClient.GetById(ctx, &catalogProto.GetByIdRequest{
    Id: req.AccountId,
    })
  if err != nil {
    return err
  }
  log.Printf("Received response from catalog: response:%s\n", catalogResponse)

  item := catalogResponse.Item


  // send payment info to payment api
  // boop boop
  /*
  card struct
  user_id
  item_id
  ? maybe more
  */

  // queue account purchased event info, full item - buyer_id - ? maybe more
  //private item info, buyerid, buyer email, item_id, vendor_id
  event := &AccountPurchased{
    PurchaseDate:         time.Now().Unix(),
    AccountLoginName:     item.LoginName,
    AccountLoginPassword: item.LoginPassword,
    AccountEmail:         item.Email,
    AccountEmailPassword: item.EmailPassword,
    AccountId:            item.Id,
    VendorId:             item.VendorId,
    BuyerId:              req.BuyerId,
    BuyerEmail:           req.BuyerEmail, 
  }

  f, err = json.Marshal(event)
  if err != nil {
    return err
  }
  // create watermill message
  msg := message.NewMessage(watermill.NewUUID(), f)
  // set message metadate
  msg.Metadata.Set("event_type", eventName)
  // Publish message on checkout topic
  if err = s.publisher.Publish("checkout.topic", msg); err != nil {
    return err
  }
  

  // return
  return &v1.Response{
    Api:  apiVersion,
    State: "Processing",
  }, nil
}

/* Steps for publisher worker
1. Grab data from kafka
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