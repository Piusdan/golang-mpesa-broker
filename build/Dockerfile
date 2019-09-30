FROM golang:latest as builder
WORKDIR /go/src/github.com/Piusdan/mpesa-gateway
RUN go get -d -v github.com/gorilla/mux
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mpesa-gateway .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/Piusdan/mpesa-gateway/mpesa-gateway .
ENTRYPOINT ["./mpesa-gateway"]