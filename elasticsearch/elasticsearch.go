// Package elasticsearch contains simple wrappers for the elastic library for use in this service
package elasticsearch

import (
	"net/http"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"gopkg.in/olivere/elastic.v5"
)

// Elasticer is a type used to interact with Elasticsearch
type Elasticer struct {
	es      *elastic.Client
	baseURL string
	index   string
	ready   bool
}

var httpClient = http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

// NewElasticer returns a pointer to an Elasticer instance
func NewElasticer(elasticsearchBase string, user string, password string, elasticsearchIndex string) (*Elasticer, error) {
	c, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(elasticsearchBase),
		elastic.SetBasicAuth(user, password),
		elastic.SetHealthcheckTimeoutStartup(60*time.Second),
		elastic.SetHttpClient(&httpClient))

	if err != nil {
		return nil, errors.Wrap(err, "Failed to create elastic client")
	}

	return &Elasticer{es: c, baseURL: elasticsearchBase, index: elasticsearchIndex}, nil
}

// Check checks that the connection is good by making a WaitForStatus call to
// the configured Elasticsearch cluster
func (e *Elasticer) Check() error {
	wait := "90s"
	err := e.es.WaitForYellowStatus(wait)
	if err != nil {
		return errors.Wrapf(err, "Cluster did not report yellow or better status within %s", wait)
	}
	return nil
}

// Search returns an *elastic.SearchService set to the right index, for further use
func (e *Elasticer) Search() *elastic.SearchService {
	return e.es.Search().Index(e.index)
}

// Scroll returns an *elastic.ScrollService set to the right index, for further use
func (e *Elasticer) Scroll() *elastic.ScrollService {
	return e.es.Scroll().Index(e.index)
}

// Close calls out to the Stop method of the underlying elastic.Client
func (e *Elasticer) Close() {
	e.es.Stop()
}
