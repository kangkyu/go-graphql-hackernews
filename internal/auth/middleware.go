package auth

import (
	"context"
	"net/http"
	"strconv"

	"github.com/kangkyu/hackernews/internal/users"
	"github.com/kangkyu/hackernews/pkg/jwt"
)

type contextKey struct {
	name string
}

var userCtxKey = &contextKey{"user"}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			// Allow unauthentiated users in
			if header == "" {
				next.ServeHTTP(w, r)
				return
			}
			tokenStr := header
			username, err := jwt.ParseToken(tokenStr)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			user := users.User{Username: username}
			id, err := users.GetUserIdByUsername(username)

			// Allow users do not exist in db
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			user.ID = strconv.Itoa(id)
			ctx := context.WithValue(r.Context(), userCtxKey, &user)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// ForContext finds the user from ctx (requires middleware)
func ForContext(ctx context.Context) *users.User {
	raw, _ := ctx.Value(userCtxKey).(*users.User)
	return raw
}
