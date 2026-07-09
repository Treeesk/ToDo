FROM golang:1.25 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY backend ./backend
RUN go build -o /app/server ./backend/cmd

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=build /app/server ./server

EXPOSE 8080

CMD ["./server"]