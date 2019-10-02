FROM golang:1.12.6 as build
WORKDIR ${GOPATH}/src/github/quintoandar
RUN apt-get update && apt-get install make git curl  -y && git clone https://github.com/quintoandar/mysqld_exporter.git
WORKDIR ${GOPATH}/src/github/quintoandar/mysqld_exporter
RUN go get -u github.com/prometheus/promu
RUN go get -u github.com/golang/dep/cmd/dep && dep ensure
RUN make build
RUN chmod +x mysqld_exporter

FROM  quay.io/prometheus/busybox:latest
COPY --from=build ["${GOPATH}/src/github/quintoandar/mysqld_exporter/mysqld_exporter", "/bin/" ]
EXPOSE      9104
ENTRYPOINT  [ "/bin/mysqld_exporter" ]
