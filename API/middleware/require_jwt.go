package middleware

import (
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/samueldaviddelacruz/go-job-board/API/models"
	"net/http"
)

type RequireJWT struct {
	Secret string
}

// Apply assumes that User middleware has already been run
// otherwise it will not work correctly.
func (mw *RequireJWT) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

// ApplyFn assumes that User middleware has already been run
// otherwise it will not work correctly.
func (mw *RequireJWT) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	var hs = jwt.NewHS256([]byte(mw.Secret))
	return func(w http.ResponseWriter, r *http.Request) {
		var pl models.CustomPayload
		token := r.Header.Get("Authorization")

		_, err := jwt.Verify([]byte(token), hs, &pl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
