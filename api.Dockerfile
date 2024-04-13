# syntax=docker/dockerfile:1
FROM golang:1.22-alpine as base

FROM base as dev

RUN apk add --no-cache make \
    && go install github.com/cosmtrek/air@latest \
	&& go install github.com/pressly/goose/v3/cmd/goose@latest \
	&& go install github.com/volatiletech/sqlboiler/v4@latest \
	&& go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest

WORKDIR /opt/app/api

COPY . .

RUN go get github.com/volatiletech/sqlboiler/v4

RUN go mod tidy

CMD ["air"]
