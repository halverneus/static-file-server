FROM golang:1.11.0 as builder

ENV BUILD_DIR /go/src/github.com/halverneus/static-file-server
ENV MAIN github.com/halverneus/static-file-server/bin/serve
ENV DEP_VERSION v0.5.0

RUN curl -fsSL -o /usr/local/bin/dep \
    https://github.com/golang/dep/releases/download/$DEP_VERSION/dep-linux-amd64 && \
    chmod +x /usr/local/bin/dep

RUN mkdir -p ${BUILD_DIR}
WORKDIR ${BUILD_DIR}
COPY . .

RUN dep ensure -vendor-only
RUN go test -race -cover ./...
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o /serve ${MAIN}

FROM scratch
COPY --from=builder /serve /
CMD ["/serve"]
