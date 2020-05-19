VERSION=v0.1.3
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"
.PHONY: build release-bin clean

build:
	go build -o ./build/helm-vault-template $(LDFLAGS) ./cmd/helm-vault-template

release-bin:
	for arch in amd64; do \
		for os in linux darwin windows; do \
			CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -o "build/helm-vault-template_"$$os"_$$arch" $(LDFLAGS) ./cmd/helm-vault-template; \
		done; \
	done

clean:
	rm -rf build
