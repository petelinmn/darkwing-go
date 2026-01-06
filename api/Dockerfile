# syntax=docker/dockerfile:1
FROM golang:1.21-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY api/ ./api/
COPY worker/ ./worker/
RUN go mod download || true
RUN CGO_ENABLED=0 go build -o /out/api ./api

FROM alpine:3.20
RUN adduser -D -H appuser
USER appuser
WORKDIR /app
COPY --from=build /out/api /app/api
ENV VERSION=dev
EXPOSE 8080
ENTRYPOINT ["/app/api"]
