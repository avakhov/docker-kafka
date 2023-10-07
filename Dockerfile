FROM golang:1.19 as builder
WORKDIR /app
COPY wait-brokers .
RUN go build -o kafka-wait-brokers .

FROM confluentinc/cp-kafka:7.5.1
COPY --from=builder /app/kafka-wait-brokers /usr/local/bin/kafka-wait-brokers
