run:
	go run cmd/main.go

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o warehouse-service cmd/main.go