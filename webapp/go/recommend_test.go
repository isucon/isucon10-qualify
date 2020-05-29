package main_test
import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

client := new(http.Client)

func TestSearchRecommendChair(t *testing.T) {
	router := NewRouter()

	req := httptest.NewRequest("GET", "/api/chair/10", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}