FROM golang:1.23 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux CGO_ENABLED=0 go build -o /app/auction cmd/auction/main.go

FROM alpine

WORKDIR /app

COPY ./cmd/auction/.env /app/cmd/auction/.env 
COPY --from=build /app/auction /app/auction

EXPOSE 8080

ENTRYPOINT ["/app/auction"]