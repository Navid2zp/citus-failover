FROM ubuntu:latest

WORKDIR ~

COPY citus-failover .

RUN apt-get update && apt-get install -y ca-certificates && apt-get clean

CMD ["./citus-failover", "-f", "config/config.yml"]