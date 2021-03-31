run-example:
	go run ./example/simple.go
race:
	go run -race ./example/...
test:
	go test -v -count=1 ./...
bench:
	go test -v -count=1 -bench=. -benchmem -benchtime=100x -run=none ./...
tidy:
	go mod tidy