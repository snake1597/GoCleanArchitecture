FROM golang:alpine AS builder

WORKDIR /app

COPY go.sum .

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o GoCleanArchitecture

FROM alpine

WORKDIR /app

COPY --from=builder /app/wait-for-it.sh /app/

COPY --from=builder /app/config.yml /app/

COPY --from=builder /app/database /app/database

COPY --from=builder /app/GoCleanArchitecture /app/

RUN apk add --no-cache bash

ENV PORT 8080 

CMD ["./wait-for-it.sh", "db:3306", "--", "./GoCleanArchitecture"]