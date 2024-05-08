FROM golang:1.22

WORKDIR /opt/app
#install utils
RUN apt-get update &&  apt-get install -y protobuf-compiler


#load dependensies
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
COPY . .


# build project
RUN make generate-api
RUN make build-project
# start
CMD ["./main","config/config.example.project.docker.yaml"]

# settings ports
#http
EXPOSE 9000
#https
EXPOSE 9001
#grpc
EXPOSE 9002
#kafka1
EXPOSE 9091
#kafka2
EXPOSE 9092
#kafka3
EXPOSE 9093
#redis
EXPOSE 6379
#postgres
EXPOSE 5432
#prometheus
EXPOSE 10091
