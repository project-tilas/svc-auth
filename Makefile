.PHONY: install test build serve clean pack serve-container deploy ship

TAG?=$(shell git rev-list HEAD --max-count=1 --abbrev-commit)

export TAG
-include .env
export $(shell sed 's/=.*//' .env)


install:
	dep ensure -vendor-only

test:
	go test ./...

build: install
	go build -ldflags "-X main.version=$(TAG)" -o svc-auth .


clean:
	rm ./svc-auth

pack:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 make build
	docker build -t gcr.io/project-tilas/svc-auth:$(TAG) .

serve: build
	./svc-auth

serve-container: pack
	docker run -d -it -p 8080:8080 --name=svc-auth gcr.io/project-tilas/svc-auth:$(TAG)


upload:
	docker push gcr.io/project-tilas/svc-auth:$(TAG)

deploy:
	envsubst < k8s/deployment.yml | kubectl apply -f -

ship: test pack upload deploy clean