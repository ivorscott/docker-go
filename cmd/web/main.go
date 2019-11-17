package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"time"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	addr := ":" + os.Getenv("ADDR_PORT")
	// A big benefit to logging your messages to the standard streams (stdout and stderr) is that
	// your application and logging are decoupled. Your application itself isn't concerned with the routing or the
	// storage of the logs, and that can make it easier to manage the logs differently depending on the environment
	// For example, we could redirect the stdout and stderr streams to on-disk files when starting the application
	//
	// go run ./cmd/web >>/tmp/info.log 2>>/tmp/error.log
	//
	// Go provides a wide range of methods to write log messages. Avoid using Panic(), Fatal() varaitions outside
	// of your main function. it's good practice to return errors instead and only panic or exit directly from main
	//
	// Custom loggers created by log.New() are concurrently safe. You can share a single logger and use it across
	// multiple goroutines and in your handlers without needing to worry about race conditions
	// If you have multiple loggers writing to the same destination that you need to be careful and ensure that the
	// destination's underlying Write() method is also safe for concurrent use.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// PreferServerCipherSuites controls whether the HTTPS connection should use Go's favored cipher suites
	// or the user's favored cipher suites. By setting this to true, Go's favored cipher suites are given
	// preference and we help increase the likelihood that a strong cipher suite which also supports forward secrecy is used.
	// CurvePreferences lets us specify which elliptic curves should be given preference during TLS handshakes.
	// Go supports fewer elliptic curves, but as of Go 1.11 only tls.CurveP256 and tls.X25519 have assembly implementations.
	// The others are very CPU intensive, so omitting them helps ensure that our server will remain performant under heavy loads
	// It may be desirable to limit the HTTPS server to only support some of these cipher suites. It's important to remember
	// that restricting the supported cipher suites to only include strong, modern ciphers can mean that users with certain older browsers
	// won't be able to use you website.
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	server := &http.Server{
		Addr:      addr,
		ErrorLog:  errorLog,
		Handler:   app.routes(),
		TLSConfig: tlsConfig,
		// IdleTimeout - By default, Go enables keep-alives on all accepted connections. This helps reduce latency (especially for HTTPS connections)
		// because a client can reuse the same connection for multiple requests without having to repeat the handshake.
		// This helps to clear-up connections where the user has unexpectedly disappeared.
		// There is no way to increase this cut-off above 3 minutes (unless you roll your own net.Listener),
		// but you can reduce it via the IdleTimeout setting. In our case, we’ve set it to 1 minute, which means
		// that all keep-alive connections will be automatically closed after 1 minute of inactivity.
		IdleTimeout: time.Minute,
		// ReadTimeout - In our code we’ve also set the ReadTimeout setting to 5 seconds. This means that if the request headers
		// or body are still being read 5 seconds after the request is first accepted, then Go will close the underlying connection.
		// Because this is a ‘hard’ closure on the connection, the user won’t receive any HTTP(S) response.
		// Setting a short ReadTimeout period helps to mitigate the risk from slow-client attacks — such as Slowloris — which could
		// otherwise keep a connection open indefinitely by sending partial, incomplete, HTTP(S) requests.
		// If you set ReadTimeout but don’t set IdleTimeout, then IdleTimeout will default to using the same setting as ReadTimeout.
		// For instance, if you set ReadTimeout to 3 seconds, then there is the side-effect that all keep-alive connections will also
		// be closed after 3 seconds of inactivity.
		ReadTimeout: 5 * time.Second,
		// WriteTimeout - The WriteTimeout setting will close the underlying connection if our server attempts to write
		// to the connection after a given period (in our code, 10 seconds). For HTTPS connections, if some data is written to the
		// connection more than 10 seconds after the request is first accepted
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", addr)
	err := server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}
