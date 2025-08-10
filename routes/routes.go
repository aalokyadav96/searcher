package routes

import (
	"naevis/ratelim"
	"naevis/search"
	"net/http"
	_ "net/http/pprof"

	"github.com/julienschmidt/httprouter"
)

func AddStaticRoutes(router *httprouter.Router) {
	router.ServeFiles("/static/uploads/*filepath", http.Dir("static/uploads"))
}

func AddSearchRoutes(router *httprouter.Router, rateLimiter *ratelim.RateLimiter) {
	router.GET("/api/v1/ac", search.Autocompleter)
	router.GET("/api/v1/search/:entityType", rateLimiter.Limit(search.SearchHandler))
	router.POST("/api/v1/emitted", search.EventHandler)
}
