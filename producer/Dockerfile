FROM golang

WORKDIR /

COPY . .

RUN go mod download

EXPOSE 5500
CMD ["$ docker run --net=host --rm confluentinc/cp-kafka:5.0.0 kafka-topics --create --topic messages --partitions 1 --replication-factor 1 --if-not-exists --zookeeper localhost:22181"]
CMD ["go", "run", "producer.go"]