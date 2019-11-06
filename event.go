package main

import ()

type PaymentInfo struct {
  Card_num int
  Date_M   int
  Date_Y   int
  Code     int
  First    string
  Last     string
  Zip      int
}

type ValidatePaymentEvent struct {
  BuyerId       string
  Info          PaymentInfo
  AmountDollars int
  AmountCents   int
  AccountId     string
}

type PaymentSuccessEvent struct {
  BuyerId   string
  AccountId string
}

type PaymentValidatedEvent struct {
  BuyerId   string
  AccountId string
}
