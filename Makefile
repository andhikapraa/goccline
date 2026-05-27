# goccline — build + release.
#
# `make`           — build for the host platform into ./bin/
# `make install`   — install the host binary to /opt/homebrew/bin (macOS)
#                    or /usr/local/bin (Linux)
# `make release`   — cross-compile + tarball all platforms into ./dist/
# `make publish V=v1.1.0`
#                  — release tarballs + upload to a GitHub release tagged V

NAME    := goccline
PKG     := github.com/andhikapraa/goccline
VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo dev)
LDFLAGS := -ldflags "-s -w -X $(PKG)/internal/components.Version=$(VERSION)"

PLATFORMS := \
	darwin/amd64 \
	darwin/arm64 \
	linux/amd64 \
	linux/arm64

.PHONY: all build install release publish clean

all: build

build:
	@mkdir -p bin
	go build $(LDFLAGS) -o bin/$(NAME) .

install: build
	@if [ -d /opt/homebrew/bin ]; then \
		install -m 0755 bin/$(NAME) /opt/homebrew/bin/$(NAME); \
		echo "→ /opt/homebrew/bin/$(NAME)"; \
	else \
		install -m 0755 bin/$(NAME) /usr/local/bin/$(NAME); \
		echo "→ /usr/local/bin/$(NAME)"; \
	fi

release: clean
	@mkdir -p dist
	@for p in $(PLATFORMS); do \
		os=$$(echo $$p | cut -d/ -f1); \
		arch=$$(echo $$p | cut -d/ -f2); \
		out="$(NAME)_$${os}_$${arch}"; \
		echo "→ $$out"; \
		GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build $(LDFLAGS) -o dist/$$out/$(NAME) .; \
		tar -C dist/$$out -czf dist/$$out.tar.gz $(NAME); \
		(cd dist && shasum -a 256 $$out.tar.gz >> SHA256SUMS); \
	done
	@ls -la dist/*.tar.gz dist/SHA256SUMS

publish: release
	@if [ -z "$(V)" ]; then echo "usage: make publish V=v1.1.0" >&2; exit 1; fi
	gh release create $(V) --title $(V) --generate-notes \
		dist/*.tar.gz dist/SHA256SUMS

clean:
	rm -rf bin dist
