FROM golang:1.11 as builder

EXPOSE 8080

RUN mkdir -p /build
WORKDIR /build
COPY . .

RUN go test -race -cover ./...
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o /serve /build/bin/serve

FROM scratch
COPY --from=builder /serve /
ENTRYPOINT ["/serve"]
CMD []
