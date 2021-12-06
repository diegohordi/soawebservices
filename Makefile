format:
	docker run --rm -v $(shell pwd):/data cytopia/gofmt -l -w .

lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.42.1 golangci-lint run -v ./...

test_race:
	docker run --rm -v $(shell pwd):/app -w /app golang:1.17 go test -count=1 -race -short ./...

test_short:
	docker run --rm -e "CGO_ENABLED=0" -v $(shell pwd):/app -w /app golang:1.17 go test -count=1 -short -cover ./...

