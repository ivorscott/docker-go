package main

import (
	"fmt"
	"net/http"
)

// Go's HTTP server assumes that the effect of any panic is isolated to the go routines serving the atice HTTP request
// every request is handled in it's own go routine. Follwoing a panic our server will log a stack trace to the server
//  error log, unwind the stack and close the underlaying HTTP connection. But it won't terminate the application, so importantly,
// any panic in your handlers won't bring down your server. Setting the Connection: Close header on the repsonse acts as a trigger
// to make Go's HTTP server automatically close the current connection after a response has been sent. It also informs the user
// that the connection will be closed.
//
// It's important to realise that our middleware will only recover panics that happen in the same goroutine that executed the
// recoverPanic() middleware. If, for example you have a handler which spins up another goroutine (e.g., to do some background processing),
//  then any panics that happpen in the second goroutine will not be recovered -- not by the recoverPanic middleware... and not by the
// panic recovery built into Go HTTP server. They will cause your app to exit and bring down the server.
// So, if you are spinning up addtional goroutines from within your web application and there is any chance of apanic, you must make sure
// that you recover any panics from within those too.
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// Set additional security headers
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("XXXX", "XXXX")
		next.ServeHTTP(w, r)
	})
}

// Custom logger for each api request
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}
