cd cmd/server
go build .
.\server.exe -grpc-port=9091 -http-port=8080 -db-host=blah -db-user=dev -db-password=dev-user5 -db-schema=checkout -catalog-service-address=localhost:9090
