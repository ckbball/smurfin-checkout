package v1

import (
  "context"
  "database/sql"
)

type repository interface {
  CreateJournalEntry(event interface{}) error
}

type JournalRepository struct {
  db *sql.DB
}

func (repository *JournalRepository) CreateJournalEntry(ctx context.Context, event interface{}) error {

  w, ok := event.(*AccountPurchased)
  if ok {
    work := w
    // get SQL connection
    c, err := repository.connect(ctx)
    if err != nil {
      return err
    }
    defer c.Close()

    res, err := c.ExecContext(ctx, `INSERT INTO checkout_journal (purchase_date, account_login_name, account_login_password, account_email, account_email_password, account_id, vendor_id, buyer_id, buyer_email) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`,
      work.PurchaseDate, work.AccountLoginName, work.AccountLoginPassword, work.AccountEmail, work.AccountEmailPassword, work.AccountId, work.VendorId, work.BuyerId, work.BuyerEmail)
    if err != nil {
      return errors.New("Error inserting journal to checkout: ", err)
    }
    return nil
  }
  return errors.New("Event does not match AccountPurchased in CreatingJournalEntry")
}

func (s *JournalRepository) connect(ctx context.Context) (*sql.Conn, error) {
  c, err := s.db.Conn(ctx)
  if err != nil {
    return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
  }
  return c, nil
}
