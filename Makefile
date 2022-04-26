VERSION=0.17.4ß
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"
.PHONY: build release-bin clean

build:
	go build -o ./build/helm-sealed-secrets $(LDFLAGS) ./cmd/helm-sealed-secrets

release-bin:
	for arch in amd64; do \
		for os in linux darwin windows; do \
			CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -o "build/helm-sealed-secrets_"$$os"_$$arch" $(LDFLAGS) ./cmd/helm-sealed-secrets; \
		done; \
	done

clean:
	rm -rf build
