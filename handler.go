package main

import (
  "bytes"
  "context"
  "encoding/gob"
  "github.com/ThreeDotsLabs/watermill"
  "github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
  "github.com/ThreeDotsLabs/watermill/message"
  catalogProto "github.com/ckbball/smurfin-catalog/proto/catalog"
  pb "github.com/ckbball/smurfin-checkout/proto/checkout"
  "time"
)

type handler struct {
  repo          repository
  catalogClient catalogProto.CatalogServiceClient
  subscriber    message.Subscriber
  publisher     message.Publisher
}

// Checkout - Receives transaction info and begins checkout process
// takes a proto.specification as input
// returns - a list of item(s)
func (s *handler) Checkout(ctx context.Context, req *pb.Request, res *pb.Response) error {
  // Make api call to catalog to validate item info
  // returns a catalogProto.Item object
  cr, err := s.catalogClient.FindItems(ctx, &catalogProto.Specification{
    ItemId: req.AccountId,
  })
  log.Printf("Found item with id: %s \n", cr.Item.Id)
  if err != nil {
    return err
  }

  // Construct data for validate-payment
  vpEvent := &ValidatePaymentEvent{
    BuyerId:       req.BuyerId,
    Info:          req.Card,
    AmountDollars: cr.Item.PriceDollars,
    AmountCents:   cr.Item.PriceCents,
    AccountId:     cr.Item.Id,
  }
  // Send validate-payment event
  // gob encode vpEvent
  var network bytes.Buffer
  enc := gob.NewEncoder(&network)
  err = enc.Encode(vpEvent)
  if err != nil {
    return err
  }
  byteSlice := network.Bytes()
  // create watermill message with gob
  msg := message.NewMessage(watermill.NewUUID(), byteSlice)
  // Publish message on checkout topic
  if err = s.publisher.Publish("checkout.topic", msg); err != nil {
    return err
  }

  // Add event to journal
  err = s.repo.CreateJournalEntry(vpEvent)
  if err != nil {
    return err
  }
  res.State = "Processing"
  // ****8 Maybe send response here saying process is under way and user will receive email with info when tx complete
  // Listen for PaymentSuccess event and then return response to client when received.

  // Finishes all the background processing needed to complete checkout
  // called as goroutine so the response can be sent back to client in a timely manner
  go FinishCheckout(s, cr, req.BuyerId)
}

func FinishCheckout(s *handler, account *catalogProto.Item, buyer_id string) {
  messages, err := s.subscriber.Subscribe(context.Background(), "payment.succes")
  if err != nil {
    log.Printf(err)
  }
  _ := process(messages, buyer_id, account.Id) //****** PROCESS ISNT COMPLETE NEED TO FIGURE OUT HOW TO PUT TIMEL LIMIT ON ITS EXECUTION

  // create emailaccountevent
  ea := &EmailAccountEvent{
    BuyerId:              buyer_id,
    AccountLogin:         account.Login,
    AccountPassword:      account.LoginPassword,
    AccountEmail:         account.Email,
    AccountEmailPassword: account.Password,
  }

  // send emailaccountevent
  var network bytes.Buffer
  enc := gob.NewEncoder(&network)
  dec := gob.NewDecoder(&network)
  err = enc.Encode(ea)
  if err != nil {
    log.Printf("error encoding EmailAccountEvent: ", err)
  }
  byteSlice := network.Bytes()
  // create watermill message with gob
  msg := message.NewMessage(watermill.NewUUID(), byteSlice)
  // Publish message on email topic
  if err = s.publisher.Publish("email.topic", msg); err != nil {
    log.Printf("error publishing emailaccountevent: ", err)
  }

  // Add event to journal
  err = s.repo.CreateJournalEntry(ea)
  if err != nil {
    log.Printf("error creating journal entry of EmailAccountEvent: ", err)
  }

  // Clear bytes buffer
  _ = dec.Decode(&EmailAccountEvent{})
  // create removeitemevent
  ri := &RemoveItemEvent{
    BuyerId:   buyer_id,
    AccountId: account.Id,
  }
  // send remove-account event to cart and catalog
  err = enc.Encode(ri)
  if err != nil {
    log.Printf("error encoding RemoveItemEvent: ", err)
  }
  byteSlice = network.Bytes()
  // create watermill message with gob
  msg = message.NewMessage(watermill.NewUUID(), byteSlice)
  // Publish message on remove topic
  if err = s.publisher.Publish("remove.topic", msg); err != nil {
    log.Printf("error publishing RemoveItemEvent: ", err)
  }

  // Add event to journal
  err = s.repo.CreateJournalEntry(ri)
  if err != nil {
    log.Printf("error creating journal entry of RemoveItemEvent: ", err)
  }

  // Clear bytes buffer
  _ = dec.Decode(&RemoveItemEvent{})
  // some other stuff i think
}

///  NOT COMPLETE - NEED TO FIGURE OUT HOW TO PUT TIME LIMIT ON EACH CALL TO PROCESS -------
func process(messages <-chan *message.Message, buyer_id string, accountId string) bool {
  timer := time.NewTimer(10 * time.Second)
  for msg := range messages {
    // decode msg payload back into struct
    var network bytes.Buffer
    var ps PaymentSuccessEvent
    network.Write(msg.payload)
    dec := gob.NewDecoder(&network)
    err = dec.Decode(&ps)
    if err != nil {
      log.Fatal("decode error: ", err)
    }
    log.Printf("received message: %s, payload buyer_id: %s", msg.UUID, ps.BuyerId)
    log.Printf("Checking if correct payload received. buyer id: %s || account: %s", buyer_id, accountId)
    if ps.BuyerId == buyer_id && ps.AccountId == accountId {
      return true
    }
  }
  return false
}
