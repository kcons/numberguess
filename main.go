package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
)

func main() {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://963b4d3c413a1d71dfb5bb25f2712659@o4509204366491648.ingest.us.sentry.io/4509204371668992",
		TracesSampleRate: 1.0,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}
	defer sentry.Flush(2 * time.Second)

	// Create an instance of sentryhttp
	sentryHandler := sentryhttp.New(sentryhttp.Options{})

	// Once it's done, you can set up routes and attach the handler as one of your middleware
	http.HandleFunc("/guess/", sentryHandler.HandleFunc(func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		hub := sentry.GetHubFromContext(ctx)
		guessStr := strings.TrimPrefix(r.URL.Path, "/guess/")
		var guess int
		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetTag("phase", "parse")
			time.Sleep(1 * time.Second)
			guessInt, err := strconv.ParseInt(guessStr, 10, 32)
			if err != nil {
				panic(err)
			}
			guess = int(guessInt)
		})
		target := rand.N(10)
		hub.AddBreadcrumb(&sentry.Breadcrumb{
			Message: fmt.Sprintf("Target: %d", target),
		}, nil)
		log.Println(guess, target)
		if guess == target {
			fmt.Fprintf(rw, "You guessed the number! %d", target)
		} else {
			hub.CaptureException(fmt.Errorf("wrong guess"))
			http.Error(rw, "You guessed so wrong!", http.StatusBadRequest)
		}

	}))

	fmt.Println("Listening and serving HTTP on :3000")

	// And run it
	if err := http.ListenAndServe(":3000", nil); err != nil {
		panic(err)
	}
}
