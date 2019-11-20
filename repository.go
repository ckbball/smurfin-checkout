// CHANGE TO SAVING EVENTS TO JOURNAL

package main

import (
  "context"
  pb "github.com/ckbball/smurfin-checkout/proto/checkout"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
)

type repository interface {
  CreateJournalEntry(event interface{}) error
}

type JournalRepository struct {
  collection *mongo.Collection
}

func (repository *JournalRepository) CreateJournalEntry(event interface{}) error {
  v, ok := event.(*PaymentRequestedEvent)
  if ok {
    work := v
    _, err = repository.collection.InsertOne(context.Background(), work)
    return nil
  }
  w, ok := event.(*AccountPurchasedEvent)
  if ok {
    work := w
    _, err = repository.collection.InsertOne(context.Background(), work)
    return nil
  }
  return errors.New("Event does not match AccountTakenDownEvent or AccountSubmittedEvent in CreatingJournalEntry")
}
