import (
  "fmt"
  "net/http"

  "github.com/getsentry/sentry-go"
  sentryhttp "github.com/getsentry/sentry-go/http"
)

func main() {

// To initialize Sentry's handler, you need to initialize Sentry itself beforehand
if err := sentry.Init(sentry.ClientOptions{
  Dsn: "https://963b4d3c413a1d71dfb5bb25f2712659@o4509204366491648.ingest.us.sentry.io/4509204371668992",
  // Set TracesSampleRate to 1.0 to capture 100%
  // of transactions for tracing.
  // We recommend adjusting this value in production,
  TracesSampleRate: 1.0,
}); err != nil {
  fmt.Printf("Sentry initialization failed: %v\n", err)
}

// Create an instance of sentryhttp
sentryHandler := sentryhttp.New(sentryhttp.Options{})

// Once it's done, you can set up routes and attach the handler as one of your middleware
http.Handle("/", sentryHandler.Handle(&handler{}))
http.HandleFunc("/foo", sentryHandler.HandleFunc(func(rw http.ResponseWriter, r *http.Request) {
  panic("y tho")
}))

fmt.Println("Listening and serving HTTP on :3000")

// And run it
if err := http.ListenAndServe(":3000", nil); err != nil {
  panic(err)
}
}
