package main

import (
  "context"
  "encoding/json"
  "github.com/ThreeDotsLabs/watermill"
  "github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
  "github.com/ThreeDotsLabs/watermill/message"
  catalogProto "github.com/ckbball/smurfin-catalog/proto/catalog"
  pb "github.com/ckbball/smurfin-checkout/proto/checkout"
  paymentProto "github.com/ckbball/smurfin-payment/proto/payment"
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

  // Construct data for payment requested
  vpEvent := &PaymentRequestedEvent{
    BuyerId:       req.BuyerId,
    Info:          req.Card,
    AmountDollars: cr.Item.PriceDollars,
    AmountCents:   cr.Item.PriceCents,
    AccountId:     cr.Item.Id,
  }
  // Marshal event
  f, err = json.Marshal(vpEvent)
  if err != nil {
    return err
  }
  // create watermill message
  msg := message.NewMessage(watermill.NewUUID(), f)
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
  // go FinishCheckout(s, cr, req.BuyerId)

  return nil
}
