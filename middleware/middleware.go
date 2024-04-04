package middleware

import (
	"context"
	"encoding/base64"
	"log"
	"net/http"
	"strings"
	"time"
)

const HeaderContentType = "Content-Type"
const MimeApplicationJSON = "application/json"

const AuthUserID = "middleware.auth.UserID"
const bearerPrefix = "Bearer "

type Middleware func(http.Handler) http.Handler

type wrapResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrapResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func NewStack(ms ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(ms) - 1; i >= 0; i-- {
			m := ms[i]
			next = m(next)
		}
		return next
	}
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := wrapResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(&wrapped, r)
		log.Println(wrapped.statusCode, r.Method, r.URL.Path, time.Since(start))
	})
}

func writeUnauthed(w http.ResponseWriter) {
	w.Header().Set(HeaderContentType, MimeApplicationJSON)
	w.WriteHeader(http.StatusUnauthorized)
}

func IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if !strings.HasPrefix(authorization, bearerPrefix) {
			writeUnauthed(w)
			return
		}

		encodedToken := strings.TrimPrefix(authorization, bearerPrefix)
		token, err := base64.StdEncoding.DecodeString(encodedToken)
		if err != nil {
			writeUnauthed(w)
			return
		}

		userID := string(token)

		ctx := context.WithValue(r.Context(), AuthUserID, userID)
		req := r.WithContext(ctx)

		log.Println("UserID:", userID)

		next.ServeHTTP(w, req)
	})
}
