# syntax=docker/dockerfile:1
FROM golang:1.22.1 as build

LABEL org.opencontainers.image.title="Digicert exporter Docker image"
LABEL org.opencontainers.image.authors="dev@vareversat.fr"
LABEL org.opencontainers.image.source="https://github.com/vareversat/digicert_exporter"
LABEL org.opencontainers.image.created="2023-11-24T22:00:00.000+01:00"
LABEL org.opencontainers.image.license="MIT"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/digicert_exporter

FROM alpine:3.19
COPY --from=build /app/digicert_exporter /usr/bin/local/digicert_exporter
EXPOSE 8080

CMD [ "/usr/bin/local/digicert_exporter" ]