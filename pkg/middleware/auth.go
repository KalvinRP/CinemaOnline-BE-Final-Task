package middleware

import (
	"context"
	"encoding/json"
	dto "finaltask/dto/result"
	jwtToken "finaltask/pkg/jwt"
	"net/http"
	"strings"
)

type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func AuthAdmin(point http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			token := r.Header.Get("Authorization")

			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "You are not allowed to be here."}
				json.NewEncoder(w).Encode(response)
				return
			}

			token = strings.Split(token, " ")[1]
			claims, err := jwtToken.DecodeToken(token)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				response := Result{Code: http.StatusUnauthorized, Message: "unauthorized"}
				json.NewEncoder(w).Encode(response)
				return
			}

			roling := claims["role"].(string)
			if roling != "admin" {
				w.WriteHeader(http.StatusUnauthorized)
				response := Result{Code: http.StatusUnauthorized, Message: "unauthorized"}
				json.NewEncoder(w).Encode(response)
				return
			}

			ctx := context.WithValue(r.Context(), "userInfo", claims)
			r = r.WithContext(ctx)
			point.ServeHTTP(w, r.WithContext(ctx))
		})
}

func Auth(point http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			token := r.Header.Get("Authorization")

			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "You are not allowed to be here."}
				json.NewEncoder(w).Encode(response)
				return
			}

			token = strings.Split(token, " ")[1]
			claims, err := jwtToken.DecodeToken(token)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				response := Result{Code: http.StatusUnauthorized, Message: "unauthorized"}
				json.NewEncoder(w).Encode(response)
				return
			}

			ctx := context.WithValue(r.Context(), "userInfo", claims)
			r = r.WithContext(ctx)
			point.ServeHTTP(w, r.WithContext(ctx))
		})
}
