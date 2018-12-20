FROM golang:1.11.3 as builder

EXPOSE 8080

ENV BUILD_DIR /build

RUN mkdir -p ${BUILD_DIR}
WORKDIR ${BUILD_DIR}

COPY go.* ./
RUN go mod download
COPY . .

RUN go test -cover ./...
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o /serve /build/bin/serve

FROM scratch
COPY --from=builder /serve /
ENTRYPOINT ["/serve"]
CMD []
