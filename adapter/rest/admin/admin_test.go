package admin_test

import (
	"testing"

	. "github.com/ryutah/virtual-ec/adapter/rest/admin"
	"github.com/stretchr/testify/assert"
)

func TestNewHandler(t *testing.T) {
	handler := NewHandler(
		NewProductEndpoint(new(mockProductSearcher), new(mockProductFinder)),
	)
	assert.NotNil(t, handler)
}
