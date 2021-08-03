build: fat

fat: bin
	CGO_ENABLED=0 go build -o bin/log-reader-ip ./cmd/fat/

bin:
	mkdir -p bin

test:
	go test -cover \
		github.com/factorysh/gitlab-log-reader/rpc \
		github.com/factorysh/gitlab-log-reader/web \
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
