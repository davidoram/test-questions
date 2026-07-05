test:
	go test ./...

clean: 
	rm -rf ./bin

race:
	go test --race ./...

benchmark:
	go test -bench=. ./ring

build: build-merge build-depend

build-merge:
	go build -o bin/merge ./merge/cmd

build-depend:
	go build -o bin/depend ./dependency/cmd

test-depend: build-depend
	bin/depend dependency/data app-dev
