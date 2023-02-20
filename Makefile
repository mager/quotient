
dev:
	go mod tidy && go run main.go

test:
	go test ./...

build:
	gcloud builds submit --tag gcr.io/quotient-378412/quotient

deploy:
	gcloud run deploy quotient \
		--image gcr.io/quotient-378412/quotient \
		--platform managed

ship:
	make test && make build && make deploy