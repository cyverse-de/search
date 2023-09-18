// Package elasticsearch contains simple wrappers for the elastic library for use in this service
package elasticsearch

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Elasticer is a type used to interact with Elasticsearch
type Elasticer struct {
	es       *elastic.Client
	baseURL  string
	index    string
	user     string
	password string

	Ready bool
}

var httpClient = http.Client{Transport: otelhttp.NewTransport(&http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}})}

// NewElasticer returns a pointer to an Elasticer instance that needs to be set up with Setup()
func NewElasticer(elasticsearchBase string, user string, password string, elasticsearchIndex string) *Elasticer {
	return &Elasticer{es: nil, baseURL: elasticsearchBase, index: elasticsearchIndex, user: user, password: password, Ready: false}
}

// Setup sets up the elasticsearch client and checks that the connection is
// good by making a WaitForStatus call to the configured Elasticsearch cluster
func (e *Elasticer) Setup() error {
	c, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(e.baseURL),
		elastic.SetBasicAuth(e.user, e.password),
		elastic.SetHealthcheckTimeoutStartup(60*time.Second),
		elastic.SetHttpClient(&httpClient))

	if err != nil {
		return errors.Wrap(err, "Failed to create elastic client")
	}
	e.es = c

	wait := "90s"
	err = e.es.WaitForYellowStatus(wait)
	if err != nil {
		return errors.Wrapf(err, "Cluster did not report yellow or better status within %s", wait)
	}
	e.Ready = true
	return nil
}

// Search returns an *elastic.SearchService set to the right index, for further use
func (e *Elasticer) Search() *elastic.SearchService {
	if e.Ready {
		return e.es.Search().Index(e.index)
	} else {
		return nil
	}
}

// Scroll returns an *elastic.ScrollService set to the right index, for further use
func (e *Elasticer) Scroll() *elastic.ScrollService {
	if e.Ready {
		return e.es.Scroll().Index(e.index)
	} else {
		return nil
	}
}

// Close calls out to the Stop method of the underlying elastic.Client
func (e *Elasticer) Close() {
	if e.es != nil {
		e.es.Stop()
	}
}
