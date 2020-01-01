package v1

import (
  "context"
  "database/sql"
  "log"
  "strconv"
  "time"

  //"github.com/golang/protobuf/ptypes"
  "encoding/json"
  "github.com/ThreeDotsLabs/watermill"
  "github.com/ThreeDotsLabs/watermill/message"
  // "github.com/go-redis/cache/v7"
  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/status"

  catalogProto "github.com/ckbball/smurfin-catalog/pkg/api/v1"
  //paymentProto "github.com/ckbball/smurfin-payment/pkg/api/v1"

  v1 "github.com/ckbball/smurfin-checkout/pkg/api/v1"
)

const (
  apiVersion = "v1"
  eventName  = "account_purchased"
)

type handler struct {
  repo          repository
  catalogClient catalogProto.CatalogServiceClient
  subscriber    message.Subscriber
  publisher     message.Publisher
}

func NewCheckoutServiceServer(repo repository, catalogClient catalogProto.CatalogServiceClient,
  subscriber message.Subscriber, publisher message.Publisher) handler {
  return &handler{
    repo:          repo,
    catalogClient: catalogClient,
    subscriber:    subscriber,
    publisher:     publisher,
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
    return nil, err
  }

  // confirm item info with grpc call
  // will need to change to full get in future
  account_id_int64 := strconv.ParseInt(req.AccountId, 10, 64)
  catalogResponse, err := s.catalogClient.GetById(ctx, &catalogProto.GetByIdRequest{
    Id: account_id_int64,
  })
  if err != nil {
    log.Printf("Error in catalog.GetById()")
    return nil, err
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
    AccountId:            req.AccountId,
    VendorId:             item.VendorId,
    BuyerId:              req.BuyerId,
    BuyerEmail:           req.BuyerEmail,
  }

  f, err = json.Marshal(event)
  if err != nil {
    return nil, err
  }
  // create watermill message
  msg := message.NewMessage(watermill.NewUUID(), f)
  // set message metadate
  msg.Metadata.Set("event_type", eventName)
  // Publish message on checkout topic
  if err = s.publisher.Publish("checkout.topic", msg); err != nil {
    return nil, err
  }

  // return
  return &v1.Response{
    Api:   apiVersion,
    State: "Processing",
    // maybe in future add more data to response about the purchased item.
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
