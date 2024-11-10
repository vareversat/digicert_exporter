# syntax=docker/dockerfile:1
FROM golang:1.23.3 as build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/digicert_exporter

FROM alpine:3.20
COPY --from=build /app/digicert_exporter /usr/bin/local/digicert_exporter

CMD [ "/usr/bin/local/digicert_exporter" ]