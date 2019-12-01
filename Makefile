IMAGE_NAME = webhook
IMAGE_VERSION = $(shell git log --abbrev-commit --format=%h -s | head -n 1)

bin: patch webhook demo

patch:
	go build -o bin/patch ./cmd/patch

webhook:
	go build -o bin/webhook ./cmd/webhook

demo:
	go build -o bin/demo-readiness-gate ./cmd/demo-readiness-gate
image:
	docker build --no-cache  -t $(IMAGE_NAME):$(IMAGE_VERSION) .

.PHONY: bin
