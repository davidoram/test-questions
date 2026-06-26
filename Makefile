test:
	go test ./...

race:
	go test --race ./...

benchmark:
	go test -bench=. ./ring

build: build-merge

build-merge:
	go build -o bin/merge ./merge/cmd
