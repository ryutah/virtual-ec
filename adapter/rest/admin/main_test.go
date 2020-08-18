//go:generate oapi-codegen -generate client,types -package admin_test -o client_gen_test.go ../../../documents/admin/openapi.yaml

package admin_test

import (
	"fmt"
	"net/http/httptest"

	. "github.com/ryutah/virtual-ec/adapter/rest/admin"
	"github.com/ryutah/virtual-ec/adapter/rest/admin/internal"
)

func strPtr(s string) *string {
	return &s
}

type testServerConfig struct {
	productFinder   internal.ProductFinder
	productSearcher internal.ProductSearcher
	productCreator  internal.ProductCreator
}

func newTestServerConfig(opts ...testServerOption) testServerConfig {
	conf := testServerConfig{
		productFinder:   new(mockProductFinder),
		productSearcher: new(mockProductSearcher),
		productCreator:  new(mockProductCreator),
	}
	for _, opt := range opts {
		opt(&conf)
	}
	return conf
}

type testServerOption func(*testServerConfig)

func withProductFinder(finder internal.ProductFinder) testServerOption {
	return func(c *testServerConfig) {
		c.productFinder = finder
	}
}

func withProductSearcher(seacher internal.ProductSearcher) testServerOption {
	return func(c *testServerConfig) {
		c.productSearcher = seacher
	}
}

func withProductCreator(creator internal.ProductCreator) testServerOption {
	return func(c *testServerConfig) {
		c.productCreator = creator
	}
}

func startTestServerAndNewClient(opts ...testServerOption) (client *Client, finish func()) {
	conf := newTestServerConfig(opts...)
	handler := NewHandler(
		NewProductEndpoint(conf.productSearcher, conf.productFinder, conf.productCreator),
	)
	s := httptest.NewServer(handler)

	client, err := NewClient(fmt.Sprintf("%s/api", s.URL))
	if err != nil {
		panic(err)
	}
	return client, s.Close
}
