FROM golang:latest
RUN mkdir /app
RUN mkdir /app/bin
ADD . /app/
WORKDIR /app

RUN go build -mod=vendor -o /bin/mocked_service mocked_service/main.go
RUN go build -mod=vendor -o /bin/service service/main.go
