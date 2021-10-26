package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/namsral/flag"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

const (
	defaultMetricsInterval = 30
	defaultAddr            = ":9295"
	defaultReashScheme     = "http"
	defaultRedashHost      = "localhost"
	defaultRedashPort      = "5000"
)

const rootDoc = `<html>
<head><title>Redash Exporter</title></head>
<body>
<h1>Redash Exporter</h1>
<p><a href="/metrics">Metrics</a></p>
</body>
</html>
`

var (
	addr            = flag.String("listen_address", defaultAddr, "The address to listen HTTP requests.")
	metricsInterval = flag.Int("metrics_interval", defaultMetricsInterval, "Interval to scrape status.")
	redashScheme    = flag.String("redash_scheme", defaultReashScheme, "target Redash scheme.")
	redashHost      = flag.String("redash_host", defaultRedashHost, "target Redash host.")
	redashPort      = flag.String("redash_port", defaultRedashPort, "target Redash port.")
)

var apiKey = os.Getenv("REDASH_API_KEY")

type redashStatus struct {
	DashboardsCount float64 `json:"dashboards_count"`
	DatabaseMetrics struct {
		Metrics [][]interface{} `json:"metrics"`
	} `json:"database_metrics"`
	Manager struct {
		OutdatedQueriesCount float64 `json:"outdated_queries_count,string"`
		Queues               struct {
			Default struct {
				Size float64 `json:"size"`
			} `json:"default"`
			Periodic struct {
				Size float64 `json:"size"`
			} `json:"periodic"`
			Queries struct {
				Size float64 `json:"size"`
			} `json:"queries"`
			ScheduledQueries struct {
				Size float64 `json:"size"`
			} `json:"scheduled_queries"`
			Schemas struct {
				Size float64 `json:"size"`
			} `json:"schemas"`
		} `json:"queues"`
	} `json:"manager"`
	QueriesCount            float64 `json:"queries_count"`
	QueryResultsCount       float64 `json:"query_results_count"`
	RedisUsedMemory         float64 `json:"redis_used_memory"`
	UnusedQueryResultsCount float64 `json:"unused_query_results_count"`
	RedashVersion           string  `json:"version"`
	WidgetsCount            float64 `json:"widgets_count"`
}

var infoLabels = []string{
	"redash_version",
}

var labels = []string{}

var (
	info = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "redash_info",
			Help: "Information of Redash.",
		},
		infoLabels,
	)

	dashboardsCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "redash_dashboards_count",
			Help: "Number of dashboards in Redash.",
		},
		labels,
	)

	queryResultsSize = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "redash_query_results_size_bytes",
			Help: "Size of Redash query results.",
		},
		labels,
	)

	dbSize = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "redash_db_size_bytes",
			Help: "Size of Redash database.",
		},
		labels,
	)

	outdatedQueriesCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "redash_outdated_queries_count",
			Help: "Number of outdated queries.",
		},
		labels,
	)

	queuesDefault = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "redash_queues_default",
			Help: "Number of default queues.",
		},
		labels,
	)

	queuesPeriodic = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "redash_queues_periodic",
			Help: "Number of periodic queues.",
		},
		labels,
	)

	queuesQueries = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "redash_queues_queries",
			Help: "Number of query queues.",
		},
		labels,
	)

	queuesScheduledQueries = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "redash_queues_scheduled_queries",
			Help: "Number of scheduled query queues.",
		},
		labels,
	)

	queuesSchemas = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "redash_queues_schemas",
			Help: "Number of schemas queues.",
		},
		labels,
	)

	queriesCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "redash_queries_count",
			Help: "Number of queries stored in redash.",
		},
		labels,
	)

	queryResultsCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "redash_query_results_count",
			Help: "Number of query results.",
		},
		labels,
	)

	redisUsedMemory = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "redash_redis_used_memory_bytes",
			Help: "Memory size used by redis in Redash.",
		},
		labels,
	)

	unusedQueryResultsCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "redash_unused_query_results_count",
			Help: "Number of unused query results.",
		},
		labels,
	)

	widgetsCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "redash_wigets_count",
			Help: "Number of widgets.",
		},
		labels,
	)
)

func getRedashStatus() (redashStatus, error) {
	url := *redashScheme + "://" + *redashHost + ":" + *redashPort
	endpoint := "/status.json"
	resp, e := http.Get(url + endpoint + "?api_key=" + apiKey)
	if e != nil {
		return redashStatus{}, fmt.Errorf("httpGet error : %v", e)
	}
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return redashStatus{}, fmt.Errorf("io read error : %v", e)
	}
	var jsonBody redashStatus
	e = json.Unmarshal(body, &jsonBody)
	if e != nil {
		return redashStatus{}, fmt.Errorf("json parse error : %v. Is api key correct?", e)
	}
	return jsonBody, nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(rootDoc))
}

func main() {
	flag.Parse()
	log.Info("start Redash exporter.")
	go func() {
		for {
			status, err := getRedashStatus()
			if err != nil {
				log.Error(err)
				time.Sleep(time.Duration(*metricsInterval) * time.Second)
				continue
			}
			label := prometheus.Labels{}
			infoLabel := prometheus.Labels{
				"redash_version": status.RedashVersion,
			}
			metrics := map[string]float64{}
			for _, metric := range status.DatabaseMetrics.Metrics {
				var key string
				var val float64
				for _, m := range metric {
					switch m.(type) {
					case string:
						key = m.(string)
					case float64:
						val = m.(float64)
					}
				}
				metrics[key] = val
			}
			info.With(infoLabel).Set(float64(1))
			dashboardsCount.With(label).Set(status.DashboardsCount)
			queryResultsSize.With(label).Set(metrics["Query Results Size"])
			dbSize.With(label).Set(metrics["Redash DB Size"])
			outdatedQueriesCount.With(label).Set(float64(status.Manager.OutdatedQueriesCount))
			queuesDefault.With(label).Set(status.Manager.Queues.Default.Size)
			queuesPeriodic.With(label).Set(status.Manager.Queues.Periodic.Size)
			queuesQueries.With(label).Set(status.Manager.Queues.Queries.Size)
			queuesScheduledQueries.With(label).Set(status.Manager.Queues.ScheduledQueries.Size)
			queuesSchemas.With(label).Set(status.Manager.Queues.Schemas.Size)
			queriesCount.With(label).Set(status.QueriesCount)
			queryResultsCount.With(label).Set(status.QueryResultsCount)
			redisUsedMemory.With(label).Set(status.RedisUsedMemory)
			unusedQueryResultsCount.With(label).Set(status.UnusedQueryResultsCount)
			widgetsCount.With(label).Set(status.WidgetsCount)

			time.Sleep(time.Duration(*metricsInterval) * time.Second)
		}
	}()
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
