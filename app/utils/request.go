package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
)

func DoRequest(r http.Handler, method, path string, form interface{}, token string) (*httptest.ResponseRecorder, error) {
	return DoRequestWithHeaders(r, method, path, form, token, map[string]string{})
}

func DoRequestWithHeaders(
	r http.Handler,
	method,
	path string,
	form interface{},
	token string,
	headers map[string]string,
) (*httptest.ResponseRecorder, error) {

	w := httptest.NewRecorder()

	var body io.Reader

	if form != nil {

		b, err := json.Marshal(form)
		if err != nil {
			return w, err
		}

		body = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return w, err
	}

	req.Header.Set("Content-Type", "application/json")

	if len(headers) > 0 {
		for header, value := range headers {
			req.Header.Set(header, value)
		}
	}

	if token != "" {
		req.Header.Set("X-ORGANONO-TOKEN", token)
	}

	r.ServeHTTP(w, req)
	return w, nil
}
