package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/project-tilas/svc-auth/handler"

	"github.com/project-tilas/svc-auth/service"

	"github.com/globalsign/mgo"
	"github.com/project-tilas/svc-auth/repository"
)

type (
	health struct {
		ServiceName string `json:"serviceName"`
		Alive       bool   `json:"alive"`
		Version     string `json:"version"`
		PodName     string `json:"podName"`
		NodeName    string `json:"nodeName"`
	}
)

var version string
var addr string
var dbAddr string
var nodeName string
var podName string

func init() {
	fmt.Println("Running SVC_AUTH version: " + version)
	addr = getEnvVar("SVC_AUTH_ADDR", "0.0.0.0:8080")
	dbAddr = getEnvVar("SVC_AUTH_DB_ADDR", "")
	nodeName = getEnvVar("SVC_AUTH_NODE_NAME", "N/A")
	podName = getEnvVar("SVC_AUTH_POD_NAME", "N/A")
}

func main() {
	testGetGolang()

	fmt.Println("Starting SVC_AUTH")

	if dbAddr == "" {
		fmt.Println("No SVC_AUTH_DB_ADDR provided")
		return
	}

	fmt.Println("Connecting to DB")

	dialInfo, err := mgo.ParseURL(dbAddr)
	dialInfo.Timeout = 30 * time.Second
	if err != nil {
		fmt.Println(err)
		return
	}

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		return tls.Dial("tcp", addr.String(), &tls.Config{})
	}

	mongoClient, err := repository.NewMongoClient(dialInfo)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Initialising User Repo")
	userRepo, err := repository.NewMongoUserRespository(mongoClient, "user")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Initialising Token Repo")
	tokenRepo, err := repository.NewMongoTokenRespository(mongoClient, "token")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Intitialising Auth Service")
	authSvc := service.NewRepositoryAuthService(userRepo, tokenRepo)

	fmt.Println("Intialising Routes")
	r := mux.NewRouter()
	apiRouter := r.PathPrefix("/api").Subrouter()
	handler.MakeServerHandler(authSvc, apiRouter)
	printRoutes(r)
	srv := &http.Server{
		Addr: addr,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}
	fmt.Println("Listening on " + addr)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)
	// Block until we receive our signal.
	<-c

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

func printRoutes(r *mux.Router) {
	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			fmt.Println("ROUTE:", pathTemplate)
		}
		pathRegexp, err := route.GetPathRegexp()
		if err == nil {
			fmt.Println("Path regexp:", pathRegexp)
		}
		queriesTemplates, err := route.GetQueriesTemplates()
		if err == nil {
			fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
		}
		queriesRegexps, err := route.GetQueriesRegexp()
		if err == nil {
			fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
		}
		methods, err := route.GetMethods()
		if err == nil {
			fmt.Println("Methods:", strings.Join(methods, ","))
		}
		fmt.Println()
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
}

func getEnvVar(env string, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}

func testGetGolang() {
	response, err := http.Get("http://web-sniffer.net/")
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", string(contents))
	}
}
