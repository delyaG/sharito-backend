package http

import (
	"backend/internal/domain"
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"net/http"
	"strconv"
)

func (a *adapter) JWTAuthMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authMiddleware := func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					r.Header.Set("Authorization", jwtauth.TokenFromHeader(r))

					token, err := jwtauth.VerifyRequest(a.jwtAuth, r, func(r *http.Request) string {
						return r.Header.Get("Authorization")
					})

					if err != nil {
						a.logger.WithError(err).Error(domain.ErrUnauthorized)
						_ = jError(w, domain.ErrUnauthorized)
						return
					}

					ctx := r.Context()
					ctx = jwtauth.NewContext(ctx, token, err)
					next.ServeHTTP(w, r.WithContext(ctx))
				})
			}

			userIDMiddleware := func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					_, claims, err := jwtauth.FromContext(r.Context())
					if err != nil {
						a.logger.WithError(err).Error("Error while extracting JWT token from context!")
						_ = jError(w, domain.ErrUnauthorized)
						return
					}

					sub, ok := claims["sub"]
					if !ok {
						a.logger.Error("Token is without 'sub' field!")
						_ = jError(w, domain.ErrUnauthorized)
						return
					}

					subStr, ok := sub.(string)
					if !ok {
						a.logger.Error("'sub' field is not a string!")
						_ = jError(w, domain.ErrUnauthorized)
						return
					}

					userID, err := strconv.Atoi(subStr)
					if err != nil {
						a.logger.WithError(err).Error("Error while converting string into int!")
						_ = jError(w, domain.ErrUnauthorized)
						return
					}

					r = r.WithContext(context.WithValue(r.Context(), domain.ContextUserID, userID))
					w.Header().Set("User-ID", strconv.Itoa(userID))
					next.ServeHTTP(w, r)
				})
			}
			chi.Chain(authMiddleware, userIDMiddleware).Handler(next).ServeHTTP(w, r)
		})
	}
}
