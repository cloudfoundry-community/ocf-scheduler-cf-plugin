FROM golang:1.18 AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /bin/ocf-scheduler-cf-plugin

CMD ["ocf-scheduler-cf-plugin"]