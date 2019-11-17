package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Create middleware chains that can be assigned to variables, appended to, and reused
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	// TODO: add authentication middleware
	dynamicMiddleware := alice.New()

	// Go's default server mux doesn't support method based routing or semantic URLS with variables in them
	// bmizerany/pat is a lightweight, provides method-based routing and support for semantic URLS.
	// A downside is the package isn't really maintained anymore but its has a clear and well-written api
	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/products", dynamicMiddleware.ThenFunc(app.products))
	return standardMiddleware.Then(mux)
}
