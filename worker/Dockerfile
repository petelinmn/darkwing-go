# syntax=docker/dockerfile:1
FROM golang:1.21-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY api/ ./api/
COPY worker/ ./worker/
RUN go mod download || true
RUN CGO_ENABLED=0 go build -o /out/worker ./worker

FROM alpine:3.20
RUN adduser -D -H appuser
USER appuser
WORKDIR /app
COPY --from=build /out/worker /app/worker
ENV INTERVAL_SECONDS=5
ENTRYPOINT ["/app/worker"]
