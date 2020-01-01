package v1

import (
  "context"
  "database/sql"
  "errors"
)

type repository interface {
  CreateJournalEntry(event interface{}) error
}

type JournalRepository struct {
  Db *sql.DB
}

func (repository *JournalRepository) CreateJournalEntry(event interface{}) error {
  ctx := context.Background()

  w, ok := event.(*AccountPurchased)
  if ok {
    work := w
    // get SQL connection
    c, err := repository.connect(ctx)
    if err != nil {
      return err
    }
    defer c.Close()

    _, err = c.ExecContext(ctx, `INSERT INTO checkout_journal (purchase_date, account_login_name, account_login_password, account_email, account_email_password, account_id, vendor_id, buyer_id, buyer_email) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`,
      work.PurchaseDate, work.AccountLoginName, work.AccountLoginPassword, work.AccountEmail, work.AccountEmailPassword, work.AccountId, work.VendorId, work.BuyerId, work.BuyerEmail)
    if err != nil {
      return err
    }
    return nil
  }
  return errors.New("Event does not match AccountPurchased in CreatingJournalEntry")
}

func (s *JournalRepository) connect(ctx context.Context) (*sql.Conn, error) {
  c, err := s.Db.Conn(ctx)
  if err != nil {
    return nil, err
  }
  return c, nil
}
