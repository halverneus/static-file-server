FROM golang:1.12.3 as builder

ENV VERSION 1.6.4
ENV BUILD_DIR /build

RUN mkdir -p ${BUILD_DIR}
WORKDIR ${BUILD_DIR}

COPY . .
RUN go test -race -cover ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -installsuffix netgo -ldflags "-X github.com/halverneus/static-file-server/cli/version.version=${VERSION}" -o pkg/linux-amd64/serve /build/bin/serve
RUN CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -a -tags netgo -installsuffix netgo -ldflags "-X github.com/halverneus/static-file-server/cli/version.version=${VERSION}" -o pkg/linux-i386/serve /build/bin/serve
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -a -tags netgo -installsuffix netgo -ldflags "-X github.com/halverneus/static-file-server/cli/version.version=${VERSION}" -o pkg/linux-arm6/serve /build/bin/serve
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -a -tags netgo -installsuffix netgo -ldflags "-X github.com/halverneus/static-file-server/cli/version.version=${VERSION}" -o pkg/linux-arm7/serve /build/bin/serve
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -tags netgo -installsuffix netgo -ldflags "-X github.com/halverneus/static-file-server/cli/version.version=${VERSION}" -o pkg/linux-arm64/serve /build/bin/serve
RUN CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -tags netgo -installsuffix netgo -ldflags "-X github.com/halverneus/static-file-server/cli/version.version=${VERSION}" -o pkg/darwin-amd64/serve /build/bin/serve
RUN CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -tags netgo -installsuffix netgo -ldflags "-X github.com/halverneus/static-file-server/cli/version.version=${VERSION}" -o pkg/win-amd64/serve.exe /build/bin/serve

# Metadata
LABEL life.apets.vendor="Halverneus" \
      life.apets.url="https://github.com/halverneus/static-file-server" \
      life.apets.name="Static File Server" \
      life.apets.description="A tiny static file server" \
      life.apets.version="v1.6.4" \
      life.apets.schema-version="1.0"
