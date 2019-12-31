package main

import (
  "encoding/json"
  "flag"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "strings"
  "time"
)

func main() {
  // get configuration
  address := flag.String("server", "http://localhost:8080", "HTTP gateway url, e.g. http://localhost:8080")
  flag.Parse()

  t := time.Now().In(time.UTC)
  pfx := t.Format(time.RFC3339Nano)

  var body string

  // Call Checkout
  resp, err := http.Post(*address+"/v1/checkout", "application/json", strings.NewReader(fmt.Sprintf(`
    {
      "api":"v1",
      "buyer_id": "2",
      "account_id": "3",
      "card": {
        "card_num":1234567890123456,
        "date_m":"03",
        "date_y":"22",
        "code":333,
        "first": "bobby",
        "last": "McFlannagan",
        "zip": 88333
      },
      "buyer_email": "bobby@gmail.com"
    }
  `, pfx, pfx, pfx)))
  if err != nil {
    log.Fatalf("failed to call Checkout method: %v", err)
  }
  bodyBytes, err := ioutil.ReadAll(resp.Body)
  resp.Body.Close()
  if err != nil {
    body = fmt.Sprintf("failed read Checkout response body: %v", err)
  } else {
    body = string(bodyBytes)
  }
  log.Printf("Checkout response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

  // parse status of checkout
  var created struct {
    API    string `json:"api"`
    Status string `json:"state"`
  }
  err = json.Unmarshal(bodyBytes, &created)
  if err != nil {
    log.Fatalf("failed to unmarshal JSON response of Checkout method: %v", err)
    fmt.Println("error:", err)
  }

}
