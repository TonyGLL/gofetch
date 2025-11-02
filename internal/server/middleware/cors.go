package middleware

import "net/http"

// CORS is a middleware that adds Cross-Origin Resource Sharing headers.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the headers that allow cross-origin requests.
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow any origin
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// If it's a preflight request (OPTIONS), we handle it here and stop the chain.
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Otherwise, pass the request to the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}
