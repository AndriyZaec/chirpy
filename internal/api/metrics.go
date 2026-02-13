package api

import (
	"fmt"
	"net/http"
)

func (cfg *ApiConfig) MetricHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(200)
	hits := cfg.FileserverHits.Load()
	fmt.Fprintf(w, `
    <html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
    `, hits)
}
