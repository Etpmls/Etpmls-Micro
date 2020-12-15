:protoc --go_out=plugins=grpc:protobuf/ protobuf/proto/*.proto

cd src\application\protobuf\proto
protoc -I . --go_out ../ --go_opt paths=source_relative --go-grpc_out ../ --go-grpc_opt paths=source_relative *.proto
protoc -I . --grpc-gateway_out ../ --grpc-gateway_opt logtostderr=true,allow_delete_body=true --grpc-gateway_opt paths=source_relative *.proto
cd ..\..\..\..