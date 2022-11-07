build:
	go build -o main cmd/main.go

image:
	docker build -t vsg .

run: image
	docker run --rm --network host  vsg

gorun:
	go run cmd/main.go

test-image:
	docker build --file test.Dockerfile -t vsg.test .

test: test-image
	docker run --rm --network host vsg.test

gotest:
	go test -count=1 -v ./...