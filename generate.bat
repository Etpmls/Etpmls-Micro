cd proto
protoc --go_out=./empb --go_opt=paths=source_relative --go-grpc_out=./empb --go-grpc_opt=paths=source_relative *.proto
cd ..