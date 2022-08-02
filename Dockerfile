################################################################################
## GO BUILDER
################################################################################
FROM golang:1.19.0 as builder

ENV VERSION 1.8.7
ENV CGO_ENABLED 0
ENV BUILD_DIR /build

RUN mkdir -p ${BUILD_DIR}
WORKDIR ${BUILD_DIR}

COPY go.* ./
RUN go mod download
COPY . .

RUN go test -cover ./...
RUN go build -a -tags netgo -installsuffix netgo -ldflags "-X github.com/halverneus/static-file-server/cli/version.version=${VERSION}" -o /serve /build/bin/serve

RUN adduser --system --no-create-home --uid 1000 --shell /usr/sbin/nologin static

################################################################################
## DEPLOYMENT CONTAINER
################################################################################
FROM scratch

EXPOSE 8080
COPY --from=builder /serve /
COPY --from=builder /etc/passwd /etc/passwd

USER static
ENTRYPOINT ["/serve"]
CMD []

# Metadata
LABEL life.apets.vendor="Halverneus" \
    life.apets.url="https://github.com/halverneus/static-file-server" \
    life.apets.name="Static File Server" \
    life.apets.description="A tiny static file server" \
    life.apets.version="v1.8.7" \
    life.apets.schema-version="1.0"
