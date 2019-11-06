package main

import (
  "context"
  "github.com/ThreeDotsLabs/watermill"
  "github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
  "github.com/ThreeDotsLabs/watermill/message"
  catalogProto "github.com/ckbball/smurfin-catalog/proto/catalog"
  pb "github.com/ckbball/smurfin-checkout/proto/checkout"
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
  cr, err := s.catalogClient.FindItems(ctx, &catalogProto.Specification{
    ItemId: req.AccountId,
  })
  log.Printf("Found item with id: %s \n", cr.Item.Id)
  if err != nil {
    return err
  }

  // Construct data for validate-payment
  vpEvent := ValidatePaymentEvent{
    BuyerId:       req.BuyerId,
    Info:          req.Card,
    AmountDollars: cr.Item.PriceDollars,
    AmountCents:   cr.Item.PriceCents,
    AccountId:     cr.Item.Id,
  }
  // Send validate-payment event
  // gob encode vpEvent
  // create watermill message with gob
  // Publish message on checkout topic

  // ****8 Maybe send response here saying process is under way and user will receive email with info when tx complete
  // Listen for PaymentSuccess event and then return response to client when received.

}
