//Package restiful is a fork that serves as facilitator for Middleware implementation
package restiful

import (
	"net/http"
)

// Handler function is a wrapper for the way HttpRouter deals with httpHandlers
type Handler func(w http.ResponseWriter, r *http.Request) error

// Handle calls each of the Handlers for the API endpoint in order.
func Handle(handlers ...Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, handler := range handlers {
			err := handler(w, r)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
		}
	})
}
