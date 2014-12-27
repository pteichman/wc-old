PHONY = check run

check:
	go test ./...
	golint ./...

run:
	go run cmd/wildcatting-server/*.go -staticdir=static
