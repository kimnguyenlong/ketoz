FROM golang:1.24 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./main.go

FROM alpine:3.22.0 AS build-release-stage

WORKDIR /ketoz

COPY --from=build-stage /app/main .

ENTRYPOINT ["./main"]