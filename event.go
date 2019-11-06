package main

import (
"context"
  "log"
  "time"
"github.com/ThreeDotsLabs/watermill"
  "github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
  "github.com/ThreeDotsLabs/watermill/message"
)

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

func 