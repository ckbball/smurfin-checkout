package main

import (
  "context"
  "fmt"
  "log"
  "time"

  "github.com/ThreeDotsLabs/watermill"
  "github.com/ThreeDotsLabs/watermill/message"
  "github.com/ThreeDotsLabs/watermill/message/router/middleware"
  "github.com/ThreeDotsLabs/watermill/message/router/plugin"
)

var (
  logger = watermill.NewStdLogger(false, false)
)

func InitWaterRouter(pub *kafka.Publisher, sub *kafka.Subscriber, repo repository) *message.Router {
  router, err := message.NewRouter(message.RouterConfig{}, logger)
  if err != nil {
    panic(err)
  }

  router.AddPlugin(plugin.SignalsHandler)

  router.AddMiddleware(
    middleware.CorrelationID,

    middleware.Retry{
      MaxRetries:      3,
      InitialInterval: time.Millisecond * 100,
      Logger:          logger,
    }.Middleware,

    middleware.Recoverer,
  )

  handler := router.AddHandler(
    "checkout_handler",
    "function.payments.result",
    sub,
    "function.accounts.result",
    pub,
    checkoutHandler{repo}.Handler,
  )

  return router
}

type checkoutHandler struct {
  repo repository
}

// Input is the message consumed by subscriber on the subscribed topic
func (c *checkoutHandler) Handler(msg *message.Message) ([]*message.Message, error) {
  log.Println("checkout handler received message: ", msg.UUID)
  q := GetQueue()

  // extract payload
  var ps PaymentProcessedEvent
  if err := json.Unmarshal(msg.Payload, &ps); err != nil {
    log.Printf("Decode error process(): ", err)
  }
  log.Printf("received message: %s, payload buyer_id: %s", msg.UUID, ps.BuyerId)
  log.Printf("Checking payment status and proceeding accordingly")

  // Grab data from queue
  e := q.Find(ps.BuyerId, ps.AccountId)
  if ps.Status == "failed" {
    out := message.NewMessage(watermill.NewUUID(), []byte("payment failed"))
    return message.Messages{msg}, nil
  } else {
    // create removeitemevent
    ri := &AccountPurchasedEvent{
      BuyerId:              e.Buyer,
      AccountLogin:         e.Account.Login,
      AccountPassword:      e.Account.LoginPassword,
      AccountEmail:         e.Account.Email,
      AccountEmailPassword: e.Account.Password,
      AccountId:            e.Account.Id,
    }
    // send AccountPurchasedEvent  to cart and catalog user, email, vendor
    b, err = json.Marshal(ri)
    if err != nil {
      log.Printf("error encoding AccountPurchasedEvent: ", err)
    }
    // create watermill message
    msg = message.NewMessage(watermill.NewUUID(), b)
  }

  err = c.repo.CreateJournalEntry(ri)
  if err != nil {
    log.Printf("Error creating journal entry for AccountPurchasedEvent")
  }

  // here we grab extra data from queue.go struct and create new message for cart, catalog, user, vendor etc.
  return message.Messages{msg}, nil
}
