package main

import (
	"net/http"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/common/promlog/flag"

	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/version"
	"github.com/vareversat/digicert_exporter/collector"
)

var (
	app         = kingpin.New("digicert_exporter", "🔥 A Prometheus exporter to monitor Digicert certificates.")
	digicertURL = app.Flag(
		"digicert.url",
		"Digicert API URL used to fetch data.",
	).Default("https://www.digicert.com/services/v2/order/certificate").Envar("DIGICERT_URL").String()
	digicertAPIKey = app.Flag("digicert.api-key",
		"Digicert API Key used to authentication.").Envar("DIGICERT_API_KEY").String()
	digicertShowExpiredCertificates = app.Flag(
		"digicert.show-expired-certificates",
		"Show expired certificate.",
	).Default("false").Envar("DIGICERT_SHOW_EXPIRED_CERTIFICATES").Bool()
	digicertMock = app.Flag(
		"digicert.mock",
		"Use mocked data as Digicert API response.",
	).Default("false").Envar("DIGICERT_MOCK").Bool()
	listenAddress = app.Flag("web.listen-port",
		"Port used to run the exporter.").Default(":10005").Envar("EXPORTER_PORT").String()
	metricPath = app.Flag("web.metrics-path",
		"Path under which to expose metrics.").Default("/metrics").Envar("EXPORTER_PATH").String()
	webMetrics = app.Flag(
		"web.exporter-metrics",
		"Show the go and http system metrics for this exporter.",
	).Default("false").Envar("EXPORTER_ENABLE_EXPORTER_METRICS").Bool()
)

func main() {
	promlogConfig := &promlog.Config{}
	flag.AddFlags(app, promlogConfig)

	kingpin.CommandLine.UsageWriter(os.Stdout)
	app.Version(version.Print("digicert_exporter"))
	app.VersionFlag.Short('v')
	app.HelpFlag.Short('h')

	kingpin.MustParse(app.Parse(os.Args[1:]))
	logger := promlog.New(promlogConfig)

	level.Info(logger).
		Log("msg", "Starting digicert_exporter", "port", listenAddress, "path", metricPath, "version", version.Info())
	level.Debug(logger).Log("msg", "Build context", "build_context", version.BuildContext())

	digicertCollector, err := collector.NewDigicertCollector(
		logger,
		*digicertURL,
		*digicertAPIKey,
		*digicertShowExpiredCertificates,
		*digicertMock,
	)
	if err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}

	promRegistry := prometheus.NewRegistry()
	promRegistry.MustRegister(digicertCollector)
	httpHandler := promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{})
	if *webMetrics {
		httpHandler = promhttp.InstrumentMetricHandler(promRegistry, httpHandler)
	}
	http.Handle(*metricPath, httpHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
            <html>
            <head><title>Digicert Exporter Metrics</title></head>
            <body>
            <h1>Digicert Exporter Metrics</h1>
            <p><a href='` + *metricPath + `'>Metrics</a></p>
            </body>
            </html>
        `))
	})

	level.Error(logger).Log("err", http.ListenAndServe(*listenAddress, nil))
}
