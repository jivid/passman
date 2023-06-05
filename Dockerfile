FROM golang:1.20 AS build

WORKDIR /app
COPY . .
VOLUME ["/var/passman"]
RUN make server
EXPOSE 8080
CMD ["./bin/server", "-dir", "/var/passman"]
