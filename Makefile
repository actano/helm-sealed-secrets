VERSION=0.1.4
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"
.PHONY: build release-bin clean

build:
	go build -o ./build/sealed-secret-template $(LDFLAGS) ./cmd/sealed-secret-template

release-bin:
	for arch in amd64; do \
		for os in linux darwin windows; do \
			CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -o "build/sealed-secret-template_"$$os"_$$arch" $(LDFLAGS) ./cmd/sealed-secret-template; \
		done; \
	done

clean:
	rm -rf build
