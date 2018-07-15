FROM golang:1.10.3 as builder

ENV BUILD_DIR /go/src/github.com/halverneus/static-file-server
ENV MAIN github.com/halverneus/static-file-server/bin/serve

RUN curl -fsSL -o /usr/local/bin/dep \
    https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && \
    chmod +x /usr/local/bin/dep

RUN mkdir -p ${BUILD_DIR}
WORKDIR ${BUILD_DIR}
COPY . .

RUN dep ensure -vendor-only
RUN go test ./...
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o /serve ${MAIN}

FROM scratch
COPY --from=builder /serve /
CMD ["/serve"]
