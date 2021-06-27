up: log
	docker-compose up -d web

log:
	mkdir -p log

build:
	docker build \
		-t gitlab-log-reader-mockup \
		--build-arg uid=`id -u` \
		.
