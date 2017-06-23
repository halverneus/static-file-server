FROM golang:latest as builder
COPY serve.go /
WORKDIR /
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o serve .

FROM scratch
COPY --from=builder /serve /
CMD ["/serve"]
