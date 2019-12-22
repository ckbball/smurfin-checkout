protoc --proto_path=proto/checkout/v1 --proto_path=third_party --go_out=plugins=grpc:pkg/api/v1 checkout.proto
protoc --proto_path=proto/checkout/v1 --proto_path=third_party --grpc-gateway_out=logtostderr=true:pkg/api/v1 checkout.proto
protoc --proto_path=proto/checkout/v1 --proto_path=third_party --swagger_out=logtostderr=true:api/swagger/v1 checkout.proto