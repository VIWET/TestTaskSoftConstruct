package server

import (
	"net/http"

	"github.com/gorilla/context"
)

func (s *server) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("UserID")
		if err != nil {
			s.logger.Error(err)
			return
		}

		context.Set(r, "userId", c.Value)

		next(w, r)
	}
}
