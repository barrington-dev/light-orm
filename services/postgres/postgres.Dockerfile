# syntax=docker/dockerfile:1
FROM postgres:latest

COPY ./scripts/. /docker-entrypoint-initdb.d/

EXPOSE 5432
