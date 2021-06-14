build: back

back: bin
	go build -o bin/log-reader-back ./cmd/back/

bin:
	mkdir -p bin
