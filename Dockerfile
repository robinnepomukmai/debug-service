FROM golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN apk add git
RUN go get github.com/prometheus/client_golang/prometheus
RUN go build -o debug-service .

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/debug-service /app/
WORKDIR /app
CMD ["./debug-service"]