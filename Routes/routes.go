package Routes

import (
	"github.com/bmizerany/pat"
	"linkShortener/Handlers"
	"net/http"
)

func SetupRoutes() *pat.PatternServeMux {
	mux := pat.New()

	mux.Post("/link", http.HandlerFunc(Handlers.InsertLink))
	mux.Get("/link", http.HandlerFunc(Handlers.GetGeneralLink))

	return mux
}
