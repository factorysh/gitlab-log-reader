build: back

back: bin
	go build -o bin/log-reader-back ./cmd/back/

bin:
	mkdir -p bin

test:
	go test -cover \
		github.com/factorysh/gitlab-log-reader/rpc \
		github.com/factorysh/gitlab-log-reader/rg \
		github.com/factorysh/gitlab-log-reader/state


docker-build:
	docker run --rm \
	-v `pwd`:/src \
	-w /src \
	bearstech/golang-dev \
	make build

docker-image:
	docker build \
		--build-arg uid=`id -u` \
		-t gitlab-log-reader \
		.

clean:
	rm -rf bin
