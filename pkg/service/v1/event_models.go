package v1

import (
// v1 "github.com/ckbball/smurfin-checkout/pkg/api/v1"
)

type AccountPurchased struct {
  PurchaseDate         int64  `json:"purchase_date"`
  AccountLoginName     string `json:"account_login_name"`
  AccountLoginPassword string `json:"account_login_password"`
  AccountEmail         string `json:"account_email"`
  AccountEmailPassword string `json:"account_email_password"`
  AccountId            string `json:"account_id"`
  VendorId             string `json:"vendor_id"`
  BuyerId              string `json:"buyer_id"`
  BuyerEmail           string `json:"buyer_email"`
}
