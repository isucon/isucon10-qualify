package main_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

client := new(http.Client)

func TestGetChairDetail(t *testing.T) {
	router := NewRouter()

	req := httptest.NewRequest("GET", "/api/chair/10", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestBuyChair(t *testing.T) {
	t.Fail()
}

func TestResponseChairRange(t *testing.T) {
	t.Fail()
}
func TestSearchChairs(t *testing.T) {
	t.Fail()
}
