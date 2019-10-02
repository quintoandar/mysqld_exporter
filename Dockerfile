FROM golang:1.12.6-alpine as build
WORKDIR ${GOPATH}/src/github.com/quintoandar
RUN apk update && apk add make git curl && git clone https://github.com/quintoandar/mysqld_exporter.git
WORKDIR ${GOPATH}/src/github.com/quintoandar/mysqld_exporter
RUN go get -u github.com/prometheus/promu
RUN go get -u github.com/golang/dep/cmd/dep && dep ensure
RUN make build
RUN chmod +x mysqld_exporter && mv mysqld_exporter /tmp/mysqld_exporter

FROM  quay.io/prometheus/busybox:latest
COPY --from=build ["/tmp/mysqld_exporter", "/bin/" ]
EXPOSE      9104
ENTRYPOINT  [ "/bin/mysqld_exporter" ]
