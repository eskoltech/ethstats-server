FROM golang:1.11.2 as builder
RUN go get -u -v github.com/eskoltech/ethstats-server
WORKDIR /go/src/github.com/eskoltech/ethstats-server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ethstats-server .

FROM alpine:3.9
WORKDIR /root/
COPY --from=builder /go/src/github.com/eskoltech/ethstats-server .
ENTRYPOINT ["./ethstats-server"]
