FROM golang:1.18-alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

#runnung stage 
RUN CGO_ENABLE=0 go build -o bookingsApp ./cmd/web

RUN chmod +x /app/bookingsApp

FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/bookingsApp /app

CMD [ "/app/bookingsApp" ]