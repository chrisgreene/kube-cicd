package main

import (
  "net/http"
  "net/http/httptest"
  "testing"
)

func TestHandleGet(t *testing.T) {
  mux := http.NewServeMux()
  mux.HandleFunc("/", homeHandler)

  writer := httptest.NewRecorder()
  request, _ := http.NewRequest("GET", "/", nil)
  mux.ServeHTTP(writer, request)

  if writer.Code != 200 {
    t.Errorf("Response code is %v", writer.Code)
  }
}
