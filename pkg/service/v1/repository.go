package v1

import (
  "context"
  "database/sql"
  "fmt"
)

type repository interface {
  CreateJournalEntry(event interface{}) error
}

type JournalRepository struct {
  db *sql.DB
}

func (repository *JournalRepository) CreateJournalEntry(ctx context.Context, event interface{}) error {
  // Check if payment requested event was passed
  v, ok := event.(*PaymentRequestedEvent)
  if ok {
    work := v

    _, err = repository.db.InsertOne(context.Background(), work)
    return nil
  }

  w, ok := event.(*AccountPurchasedEvent)
  if ok {
    work := w
    _, err = repository.db.InsertOne(context.Background(), work)
    return nil
  }
  return errors.New("Event does not match AccountTakenDownEvent or AccountSubmittedEvent in CreatingJournalEntry")
}

func (s *JournalRepository) connect(ctx context.Context) (*sql.Conn, error) {
  c, err := s.db.Conn(ctx)
  if err != nil {
    return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
  }
  return c, nil
}
