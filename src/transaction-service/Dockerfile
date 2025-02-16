#start from golang base image
FROM golang:1.23-alpine

#install dependencies
RUN apk update
RUN apk add --no-cache git build-base

#setup workdir
RUN mkdir -p /app
ENV WORKDIR /app
WORKDIR $WORKDIR

RUN go install github.com/air-verse/air@latest


COPY ./src/transaction-service/ ./src/transaction-service/
COPY grpc/ grpc/
COPY transaction.air.toml .
COPY .env .
COPY go.mod .
COPY go.sum .
COPY makefile .
COPY tmp .

# install dependencies project
WORKDIR $WORKDIR/src/transaction-service/
RUN go mod download

# run the Go app.
WORKDIR $WORKDIR
CMD ["air", "-c", "transaction.air.toml"]
