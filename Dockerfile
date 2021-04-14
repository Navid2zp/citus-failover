FROM ubuntu:latest

WORKDIR ~

COPY citus-failover .

CMD ["./citus-failover", "-f", "config/config.yml"]