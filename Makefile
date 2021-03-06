PHONY = check run

check:
	go vet ./...
	go test ./...
	golint ./...

run:
	go run cmd/wildcatting-server/*.go -staticdir=static
