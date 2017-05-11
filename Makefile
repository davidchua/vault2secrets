VERSION=$(shell git describe --abbrev=0 --tags)
docker:
	docker build -t vault2secrets .

tag:
	docker tag vault2secrets "cubiclerebels/vault2secrets:$(VERSION)"

push:
	docker push cubiclerebels/vault2secrets:$(VERSION)
	docker push cubiclerebels/vault2secrets:latest
