gen:
	protoc --go_out=. --go-grpc_out=. proto/*.proto

clean:
	rm -rf pb

run: 
	go run main.go

test: 
	go test -cover -race ./...