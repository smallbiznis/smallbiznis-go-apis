NAME := 127.0.0.1:5001/oauth2
TAG := $$(git log -1 --pretty=%h)
IMG := ${NAME}:${TAG}
LATEST := ${NAME}:latest

buildimage:
	@docker build --platform=linux/amd64 -t ${IMG} .
	@docker tag ${IMG} ${LATEST}

pushimage: buildimage
	@docker push ${LATEST}

e2e:
	go test -timeout 1000s ottogift-product/test/e2e -count=1 -v