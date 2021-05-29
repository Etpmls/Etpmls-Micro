cd proto_test
protoc --go_out=./pb_test --go_opt=paths=source_relative --go-grpc_out=./pb_test --go-grpc_opt=paths=source_relative *.proto
cd ..\