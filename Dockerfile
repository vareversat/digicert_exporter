# syntax=docker/dockerfile:1
FROM golang:1.21.4 as build
LABEL org.opencontainers.image.authors="dev@vareversat.fr"
LABEL org.opencontainers.image.source="https://github.com/vareversat/digicert_exporter"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/digicert_exporter

FROM alpine:3.18
COPY --from=build /app/digicert_exporter /usr/bin/local/digicert_exporter
EXPOSE 8080

CMD [ "/usr/bin/local/digicert_exporter" ]