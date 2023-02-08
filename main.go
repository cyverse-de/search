package main

import (
	"context"
	_ "expvar"
	"flag"
	"net/http"

	"github.com/cyverse-de/configurate"
	"github.com/spf13/viper"

	"github.com/cyverse-de/go-mod/otelutils"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"

	"github.com/cyverse-de/search/data"
	"github.com/cyverse-de/search/elasticsearch"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const serviceName = "search"

var log = logrus.WithFields(logrus.Fields{
	"service": serviceName,
	"art-id":  serviceName,
	"group":   "org.cyverse",
})

var (
	cfgPath = flag.String("config", "", "Path to the configuration file.")
	cfg     *viper.Viper
)

func init() {
	flag.Parse()
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func loadConfig(cfgPath string) {
	var err error
	cfg, err = configurate.Init(cfgPath)
	if err != nil {
		log.Fatal(err)
	}
}

//nolint
func GetElasticsearchReadyHandler(e *elasticsearch.Elasticer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if e.Ready {
			w.WriteHeader(200)
			w.Write([]byte("ready"))
		} else {
			w.WriteHeader(500)
			w.Write([]byte("not ready"))
		}
	}
}

func newRouter(e *elasticsearch.Elasticer) *mux.Router {
	r := mux.NewRouter()
	r.Use(otelmux.Middleware(serviceName))
	r.Handle("/debug/vars", http.DefaultServeMux)
	r.HandleFunc("/ready", GetElasticsearchReadyHandler(e))
	data.RegisterRoutes(r.PathPrefix("/data/").Subrouter(), cfg, e, log)

	return r
}

func main() {
	log.Info("Starting up the search service.")

	var tracerCtx, cancel = context.WithCancel(context.Background())
	defer cancel()
	shutdown := otelutils.TracerProviderFromEnv(tracerCtx, serviceName, func(e error) { log.Fatal(e) })
	defer shutdown()

	loadConfig(*cfgPath)
	e := elasticsearch.NewElasticer(cfg.GetString("elasticsearch.base"), cfg.GetString("elasticsearch.user"), cfg.GetString("elasticsearch.password"), cfg.GetString("elasticsearch.index"))
	defer e.Close()

	// Check in a goroutine so the service can start up & respond to health checks sooner
	go func(e *elasticsearch.Elasticer) {
		log.Info("Setting up Elasticsearch connection")
		err := e.Setup()
		if err != nil {
			log.Fatal(err)
		}
		log.Info("Elasticsearch is ready")
	}(e)

	r := newRouter(e)
	listenPortSpec := ":" + "60000"
	log.Infof("Listening on %s", listenPortSpec)
	log.Fatal(http.ListenAndServe(listenPortSpec, r))
}
