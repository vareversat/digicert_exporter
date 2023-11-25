# Digicert Certificates Exporter

Export [Digicert](https://www.digicert.com/) certificates information to [Prometheus](https://prometheus.io).

Metrics are computed by retrieving data using
the [Digicert REST API](https://dev.digicert.com/en/certcentral-apis/services-api/orders/list-orders.html).
Currently, the exporter use the **v2** version of the **/order/certificate** REST endpoint

## Prerequisites

In order to run this exporter, you need :

1) A valid Digicert account **and** create an API
   key [here](https://www.digicert.com/secure/automation/api-keys/) with the **View Only** permission
2) Docker and/or Go v1.21 installed on your computer

## To run it

- With go installed :

```shell
export DIGICERT_API_KEY=my-key && go run .
```

- With Docker installed :

```shell
docker build -t digicert_exporter . && docker run -e DIGICERT_API_KEY=my-key digicert_exporter
```

## Exporter's Metrics

| Metric                                          | Description                                 | Labels                                                          | Mandatory ? |
|-------------------------------------------------|---------------------------------------------|-----------------------------------------------------------------|-------------|
| `digicert_api_up`                               | Was the last Digicert API query successful  |                                                                 | ✅           |
| `digicert_certificate_expire_timestamp_seconds` | Certificate expiration date                 | certificate_common_name, certificate_id, order_id, organization | ✅           |
| `digicert_scrape_duration_seconds`              | Exporter scrape duration in seconds         |                                                                 | ✅           |
| `promhttp_metric_handler_requests_in_flight`    | Current number of scrapes being served      |                                                                 | ❌           |
| `promhttp_metric_handler_requests_total`        | Total number of scrapes by HTTP status code | code                                                            | ❌           |

## Flags

```shell
./digicert_exporter --help
```

| Flag                                      | Description                             | Default                                                  | Env vars                           |
|-------------------------------------------|-----------------------------------------|----------------------------------------------------------|------------------------------------|
| --log.level                               | Logging level                           | `info`                                                   | ❌                                  |
| --log.format                              | Logging format                          | `logfmt`                                                 | ❌                                  |
| --web.listen-port                         | Port used to run the exporter           | `:10005`                                                 | EXPORTER_PORT                      |
| --web.metrics-path                        | Path under which to expose metrics      | `/metrics`                                               | EXPORTER_PATH                      |
| --digicert.url                            | Digicert API URL used to fetch data     | `https://www.digicert.com/services/v2/order/certificate` | DIGICERT_URL                       |
| --digicert.api-key                        | Digicert API Key used to authentication | `""`                                                     | DIGICERT_API_KEY                   |
| --[no-]digicert.show-expired-certificates | Show expired certificate                | `false`                                                  | DIGICERT_SHOW_EXPIRED_CERTIFICATES |
