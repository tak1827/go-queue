lint:
	go vet ./...

fmt:
	gofmt -w -l .

test:
	go test ./... -v -race

bench:
	go test ./... -bench=. -benchtime=3s -benchmem
