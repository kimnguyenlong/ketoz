FROM golang:1.24 AS build-stage

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./main.go

FROM oryd/keto:v0.14.0 AS build-release-stage

WORKDIR /ketoz

COPY --from=build-stage /app/main .
COPY ./keto/config.yml /home/ory/config.yml
COPY ./keto/namespaces.ts /home/ory/namespaces.ts
COPY start.sh .

ENTRYPOINT ["/bin/sh", "start.sh"]