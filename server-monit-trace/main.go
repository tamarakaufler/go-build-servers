package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"syscall"
	"time"
)

// middleware definition
//	"middler" because:
// 		middler : one belonging to an intermediate group (Merriam-Webster)
//		middler is someone who will always care for you, be there for you and never betray you (Urban dictionary)

//		for monitoring and tracing middleware etc
type middler func(http.Handler) http.Handler
type middlers []middler

type controller struct {
	logger        *log.Logger
	nextRequestID func(r *http.Request) string
	healthy       int64
}

func main() {

	// customize server address if desired
	listenAddr := ":5000"
	if len(os.Args) == 2 {
		listenAddr = os.Args[1]
	}

	logger := log.New(os.Stdout, "server-monit-trace: ", log.LstdFlags)
	logger.Printf("Server is starting...")

	c := &controller{
		logger: logger,
		nextRequestID: func(r *http.Request) string {
			nri := strconv.FormatInt(time.Now().UnixNano(), 36) + "-" + r.RemoteAddr
			fmt.Printf(">>> nextRequestID = %s <<<\n", nri)
			return nri
		},
	}

	fmt.Println("Setting up routing")
	router := http.NewServeMux()
	router.HandleFunc("/", c.mainHandler)
	router.HandleFunc("/healthz", c.healthzHandler)

	// middlers are executed, so that monitoring is executed close to the route handler
	// and all is wrapped by the trace middler
	server := &http.Server{
		Addr:         listenAddr,
		Handler:      (middlers{c.trace, c.monit}).apply(router),
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	ctx := c.shutdown(context.Background(), server)

	fmt.Println("Starting server ...")

	logger.Printf("Server is ready to handle requests at %q\n", listenAddr)
	atomic.StoreInt64(&c.healthy, time.Now().UnixNano())

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %q: %s\n", listenAddr, err)
	}
	<-ctx.Done()
	logger.Printf("Server stopped\n")
}

// apply applies middlers starting with the closest middler to the handler
// ie the last one provided.
func (m middlers) apply(h http.Handler) http.Handler {
	if len(m) == 0 {
		return h
	}
	l := len(m) - 1
	return m[:l].apply(m[l](h))
}

// graceful server shutdown
func (c *controller) shutdown(ctx context.Context, server *http.Server) context.Context {
	ctx, done := context.WithCancel(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		defer done()

		<-quit
		signal.Stop(quit)
		close(quit)

		atomic.StoreInt64(&c.healthy, 0)
		server.ErrorLog.Printf("Server is shutting down...\n")

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			server.ErrorLog.Fatalf("Could not gracefully shutdown the server: %s\n", err)
		}
	}()

	return ctx
}

// handlers
func (c *controller) mainHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("--> index")

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	time.Sleep(1 * time.Second)
	fmt.Fprintf(w, "Hello, World!\n")
	w.Write([]byte(`You are challenging!`))
}

func (c *controller) healthzHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("--> healthz")

	if a := atomic.LoadInt64(&c.healthy); a == 0 {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		fmt.Fprintf(w, "uptime: %s\n", time.Since(time.Unix(0, a)))
	}
}

// middleware handlers
func (c *controller) monit(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("\tstart in monit: %s", time.Now())
		defer func(start time.Time) {
			requestID := w.Header().Get("X-Request-Id")
			if requestID == "" {
				requestID = "unknown"
			}
			c.logger.Println(requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), time.Since(start))
			log.Printf("\tend in monit: %s", time.Now())
		}(time.Now())
		h.ServeHTTP(w, r)
		log.Printf("\tafter handler in monit finished: %s", time.Now())
	})
}

func (c *controller) trace(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("\tstart in trace: %s", time.Now())

		requestID := r.Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = c.nextRequestID(r)
			log.Printf("\tSetting requestID in trace: %s", requestID)
		}
		w.Header().Set("X-Request-Id", requestID)
		h.ServeHTTP(w, r)
		log.Printf("\tafter handler in trace finished: %s", time.Now())
		log.Printf("\tend in trace: %s", time.Now())
	})
}

// main_test.go
var (
	_ http.Handler = http.HandlerFunc((&controller{}).mainHandler)
	_ http.Handler = http.HandlerFunc((&controller{}).healthzHandler)
	_ middler      = (&controller{}).monit
	_ middler      = (&controller{}).trace
)
