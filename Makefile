PHONY = check run

check:
	go test ./...
	go lint ./...

run:
	go run cmd/wildcatting-server/*.go -staticdir=static
