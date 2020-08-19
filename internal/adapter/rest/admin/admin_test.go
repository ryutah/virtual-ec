package admin_test

import (
	"testing"

	. "github.com/ryutah/virtual-ec/internal/adapter/rest/admin"
	"github.com/stretchr/testify/assert"
)

func TestNewHandler(t *testing.T) {
	handler := NewHandler(
		NewProductEndpoint(new(mockProductSearcher), new(mockProductFinder), new(mockProductCreator)),
	)
	assert.NotNil(t, handler)
}
