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

  w, ok := event.(*AccountPurchasedEvent)
  if ok {
    work := w
    // get SQL connection
    c, err := s.connect(ctx)
    if err != nil {
      return nil, err
    }
    defer c.Close()

    res, err := c.ExecContext(ctx, `INSERT INTO items (VendorId, BlueEssence, RiotPoints, Solo, Flex, PriceDollars, PriceCents, Level, Email, Password, Login, LoginPassword) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
      req.Item.VendorId, req.Item.BlueEssence, req.Item.RiotPoints, req.Item.Solo, req.Item.Flex, req.Item.PriceDollars, req.Item.PriceCents, req.Item.Level, req.Item.Email, req.Item.EmailPassword, req.Item.LoginName, req.Item.LoginPassword)
    if err != nil {
      return nil, status.Error(codes.Unknown, "failed to insert into item-> "+err.Error())
    }
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
