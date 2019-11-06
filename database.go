// CHANGE TO MYSQL
package main

import (
  "fmt"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
  //"os"
)

type Database struct {
  *gorm.DB
}

var DB *gorm.DB

func Init() *gorm.DB {
  db, err := gorm.Open("mysql", "dev:dev-user5@/quikapplicant")
  if err != nil {
    fmt.Println("db err: ", err)
  }

  db.DB().SetMaxIdleConns(10)
  DB = db
  return DB
}

func GetDB() *gorm.DB {
  return DB
}
