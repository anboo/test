package rpc_auth

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	ent "test-question/internal/entity/user"
)

type ctxKey string

const CtxUserID ctxKey = "user_id"

type AuthUseCase interface {
	AuthorizeUser(ctx context.Context, username, password string) (*ent.User, error)
}

func BasicAuthMiddleware(auth AuthUseCase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			h := r.Header.Get("Authorization")
			if !strings.HasPrefix(h, "Basic ") {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			raw, _ := base64.StdEncoding.DecodeString(strings.TrimPrefix(h, "Basic "))
			parts := strings.SplitN(string(raw), ":", 2)
			if len(parts) != 2 {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			username := parts[0]
			password := parts[1]

			user, err := auth.AuthorizeUser(r.Context(), username, password)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), CtxUserID, user.ID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func GetUserID(ctx context.Context) string {
	v := ctx.Value(CtxUserID)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func InjectUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, CtxUserID, userID)
}
