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

func newRouter(e *elasticsearch.Elasticer) *mux.Router {
	r := mux.NewRouter()
	r.Use(otelmux.Middleware(serviceName))
	r.Handle("/debug/vars", http.DefaultServeMux)
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
	e, err := elasticsearch.NewElasticer(cfg.GetString("elasticsearch.base"), cfg.GetString("elasticsearch.user"), cfg.GetString("elasticsearch.password"), cfg.GetString("elasticsearch.index"))
	if err != nil {
		log.Fatal(err)
	}
	defer e.Close()

	r := newRouter(e)
	listenPortSpec := ":" + "60000"
	log.Infof("Listening on %s", listenPortSpec)
	log.Fatal(http.ListenAndServe(listenPortSpec, r))
}
