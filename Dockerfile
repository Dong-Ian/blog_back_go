FROM golang:alpine3.23 as base

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

    
FROM base as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o backend .


FROM golang:alpine3.23 as release

WORKDIR /home/node

COPY --from=builder /app/backend ./backend

EXPOSE $APP_PORT

CMD [ "./backend" ]