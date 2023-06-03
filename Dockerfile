# Builder image
FROM golang:1.17-buster AS builder
# Git is required for fetching the dependencies.
RUN apt-get update && apt-get install -y git ca-certificates && update-ca-certificates
WORKDIR ./DirtyAPI
COPY . .
RUN go get -d -v
# RUN go build -o /go/bin/dusty
RUN go build -o /DirtyAPI
# Smallish image for actual running
FROM debian:buster-slim

# Install required packages
RUN apt-get update && apt-get install -y libaio1 unzip

# Copy the Oracle Instant Client files
COPY ./instantclient_21_9 /opt/oracle/instantclient_21_9

# Set environment variables
ENV ORACLE_HOME /opt/oracle/instantclient_21_9
ENV LD_LIBRARY_PATH $LD_LIBRARY_PATH:$ORACLE_HOME

# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /DirtyAPI /DirtyAPI
# COPY ./data.sqlite /go/bin
COPY ./Nexus.Sqlite .
COPY ./tnsnames.ora .
EXPOSE 8080
ENTRYPOINT ["./DirtyAPI"]