package main

import (
  "crypto/sha1"
  "strings"

  catalogProto "github.com/ckbball/smurfin-catalog/proto/catalog"
)

// this file will hold the data structure that holds the relvant info for processing future events.
// it will be a map where the key is buyer_id.account_id hashed?

type Queue struct {
  queue map[int]*Entry
}

var que *Queue

func NewQueue() error {
  que.queue = make(map[int]*Entry)
  return nil
}

func GetQueue() *Queue {
  return &que
}

func (q *Queue) Add(b string, a *catalogProto.Item) error {
  key := Hash(a.Id, b)
  e := &Entry{
    Buyer:   b,
    Account: a,
  }
  q.queue[key] = e
  return nil
}

func (q *Queue) Find(b string, a string) *Entry {
  key := Hash(a, b)
  e := q.queue[key]
  return e
}

func (q *Queue) Remove(b string, a string) {
  key := Hash(a, b)
  delete(q.queue, key)
}

type Entry struct {
  Buyer   string
  Account *catalogProto.Item
}

// takes acc_id, buyer_id. hashes into key for queue
func Hash(a string, b string) string {
  var sb strings.Builder
  sb.WriteString(b)
  sb.WriteString("-.-")
  sb.WriteString(a)
  key := sb.String()
  hash := sha1.New()
  hash.Write([]byte{key})
  ins := string(hash)
  return ins
}
