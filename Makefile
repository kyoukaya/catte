build:
	go build .

test:
	go test .

lint:
	golangci-lint run .

format:
	gofmt -w .
	goimports -w -local github.com/kyoukaya/catte .
