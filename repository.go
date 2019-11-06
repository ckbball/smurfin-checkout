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
  // Somehow determine which event it is, of the two
  // Store it in db
  // return error
}
