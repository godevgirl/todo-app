FROM golang:1.23.2 AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o todo-app .

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/todo-app . 

EXPOSE 8080

CMD [ "./todo-app" ]

