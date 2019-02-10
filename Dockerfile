FROM golang:1.11.2 as builder
RUN go get -u -v github.com/eskoltech/ethstats
WORKDIR /go/src/github.com/eskoltech/ethstats
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ethstats .

FROM alpine:3.9
WORKDIR /root/
COPY --from=builder /go/src/github.com/eskoltech/ethstats .
ENTRYPOINT ["./ethstats"]
