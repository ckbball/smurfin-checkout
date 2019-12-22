package v1

import (
  v1 "github.com/ckbball/smurfin-checkout/pkg/api/v1"
  "time"
)

type AccountPurchased struct {
  PurchaseDate         time.Time
  AccountName          string
  AccountPassword      string
  AccountEmail         string
  AccountEmailPassword string
  AccountId            string
  VendorId             string
  BuyerId              string
  BuyerEmail           string
}
