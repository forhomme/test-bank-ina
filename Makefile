build-deployment:
	go build -tags musl -v -o test-bank-ina main.go

run-deployment:
	./test-bank-ina

run:
	go run main.go