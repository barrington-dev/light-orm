# syntax=docker/dockerfile:1
FROM golang:1.22-alpine as base

# Setup dev environment
FROM base AS dev

RUN apk add --no-cache make \
    && go install github.com/cosmtrek/air@latest \
	&& go install github.com/pressly/goose/v3/cmd/goose@latest \
	&& go install github.com/volatiletech/sqlboiler/v4@latest \
	&& go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest

WORKDIR /opt/app/api

COPY . .

RUN go get github.com/volatiletech/sqlboiler/v4 \
    && go get github.com/stretchr/testify

RUN go mod tidy

RUN export CGO_ENABLED=0 GOOS=linux \
    && go build -o main cmd/api/main.go

RUN go mod tidy

CMD ["air"]

# Setup prod environment
FROM alpine:edge AS prod

WORKDIR /app

COPY --from=dev /opt/app/api/main .

RUN apk --no-cache add ca-certificates tzdata

ENTRYPOINT ["/app/main"]