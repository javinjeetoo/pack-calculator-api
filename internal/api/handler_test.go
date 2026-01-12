package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalculate_ItemsTooLarge(t *testing.T) {
	h := NewHandler([]int{250, 500, 1000, 2000, 5000})

	req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewBufferString(`{"items":1000001}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.Calculate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
}
