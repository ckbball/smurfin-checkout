package cmd

import (
  "context"
  "database/sql"
  "flag"
  "fmt"

  // mysql driver
  "github.com/go-redis/cache/v7"
  "github.com/go-redis/redis/v7"
  _ "github.com/go-sql-driver/mysql"
  "github.com/vmihailenco/msgpack/v4"

  "github.com/ckbball/smurfin-catalog/pkg/protocol/grpc"
  "github.com/ckbball/smurfin-catalog/pkg/protocol/rest"
  "github.com/ckbball/smurfin-catalog/pkg/service/v1"
)

// Config is configuration for Server
type Config struct {
  // gRPC server start parameters section
  // gRPC is TCP port to listen by gRPC server
  GRPCPort string

  // the port to listen for http calls
  HTTPPort string

  // DB Datastore parameters section
  // DatastoreDBHost is host of database
  DatastoreDBHost string
  // DatastoreDBUser is username to connect to database
  DatastoreDBUser string
  // DatastoreDBPassword password to connect to database
  DatastoreDBPassword string
  // DatastoreDBSchema is schema of database
  DatastoreDBSchema string
  // address for single redis node
  RedisAddress string
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
  ctx := context.Background()

  // get configuration
  var cfg Config
  flag.StringVar(&cfg.GRPCPort, "grpc-port", "", "gRPC port to bind")
  flag.StringVar(&cfg.HTTPPort, "http-port", "", "http port to bind")
  flag.StringVar(&cfg.DatastoreDBHost, "db-host", "", "Database host")
  flag.StringVar(&cfg.DatastoreDBUser, "db-user", "", "Database user")
  flag.StringVar(&cfg.DatastoreDBPassword, "db-password", "", "Database password")
  flag.StringVar(&cfg.DatastoreDBSchema, "db-schema", "", "Database schema")
  flag.StringVar(&cfg.RedisAddress, "redis-address", "", "Redis address")
  flag.Parse()

  if len(cfg.GRPCPort) == 0 {
    return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
  }

  if len(cfg.HTTPPort) == 0 {
    return fmt.Errorf("invalid TCP port for http server: '%s'", cfg.HTTPPort)
  }

  // add MySQL driver specific parameter to parse date/time
  // Drop it for another database
  param := "parseTime=true"

  // for non localhost db %s:%s@tcp(%s)/%s?%s
  // currently set for localhost
  dsn := fmt.Sprintf("%s:%s@/%s?%s",
    cfg.DatastoreDBUser,
    cfg.DatastoreDBPassword,
    // cfg.DatastoreDBHost,
    cfg.DatastoreDBSchema,
    param)
  db, err := sql.Open("mysql", dsn)
  if err != nil {
    return fmt.Errorf("failed to open database: %v", err)
  }
  defer db.Close()

  redisPool := initRedis(cfg.RedisAddress)

  v1API := v1.NewCatalogServiceServer(db, redisPool)

  // run http gateway
  go func() {
    _ = rest.RunServer(ctx, cfg.GRPCPort, cfg.HTTPPort)
  }()

  return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}

func initRedis(address string) *cache.Codec {
  ring := redis.NewRing(&redis.RingOptions{
    Addrs: map[string]string{
      "server1": ":" + address,
    },
  })

  codec := &cache.Codec{
    Redis: ring,

    Marshal: func(v interface{}) ([]byte, error) {
      return msgpack.Marshal(v)
    },
    Unmarshal: func(b []byte, v interface{}) error {
      return msgpack.Unmarshal(b, v)
    },
  }

  return codec
}
