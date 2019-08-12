package auth

import (
	"bytes"
	"io/ioutil"
	"context"
	"net/http"
	"encoding/base64"
	"compress/gzip"
)

type CookieSetHandler struct {
	Handler http.Handler
}

func (handler CookieSetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(AuthCookieName)
	if err == nil {
		ctx := context.WithValue(r.Context(), CSRFRequiredKey, handler.isCSRFRequired(r))
		r = r.WithContext(ctx)

		if r.Header.Get("Authorization") == "" {
			r.Header.Set("Authorization", decompress(cookie.Value))
		}
	}

	handler.Handler.ServeHTTP(w, r)
}

func decompress(str string) string {
	data, _ := base64.StdEncoding.DecodeString(str)
	gz, err := gzip.NewReader(bytes.NewBuffer([]byte(data)))
	if err != nil {
		panic(err)
	}
	decompressed, err := ioutil.ReadAll(gz)
	if err != nil {
		panic(err)
	}
	return string(decompressed)
}

// We don't validate CSRF token for GET requests
// since they are not changing the state
func (handler CookieSetHandler) isCSRFRequired(r *http.Request) bool {
	return (r.Method != http.MethodGet && r.Method != http.MethodHead && r.Method != http.MethodOptions)
}

func IsCSRFRequired(r *http.Request) bool {
	required, ok := r.Context().Value(CSRFRequiredKey).(bool)
	return ok && required
}
