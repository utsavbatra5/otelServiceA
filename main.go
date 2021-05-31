package main

import (
	"ServiceA/tracer"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
	"os"
)

const (
	tracingConfigKey = "tracing"
	applicationName  = "ServiceA"
)

func requestHandler(w http.ResponseWriter, r *http.Request, postURL string, client *http.Client) {
	traceID, spanID, _ := tracer.ExtractTraceInfo(r.Context())

	log.Println(fmt.Sprintf("Trace ID for this request in %s is: %s and Span Id is: %s", applicationName, traceID, spanID))

	req, err := http.NewRequestWithContext(r.Context(), "POST", postURL, nil)
	if err != nil {
		log.Printf("Error on request", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error on response", err)
	}
	defer resp.Body.Close()

}

func getViper() *viper.Viper {
	viper.SetConfigName("config")

	// Set the path to look for the configurations file
	viper.AddConfigPath(".")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	return viper.GetViper()

}

func initTracing(v *viper.Viper, appName string) (tracer.Tracing, error) {
	var tracing = tracer.Tracing{
		Propagator:     propagation.TraceContext{},
		TracerProvider: trace.NewNoopTracerProvider(),
	}
	var traceConfig tracer.Config
	err := v.UnmarshalKey(tracingConfigKey, &traceConfig)
	if err != nil {
		return tracer.Tracing{}, err
	}
	traceConfig.ApplicationName = appName
	tracerProvider, err := tracer.ConfigureTracerProvider(traceConfig)
	if err != nil {
		return tracer.Tracing{}, err
	}

	tracing.TracerProvider = tracerProvider
	return tracing, nil
}

func main() {

	v := getViper()

	tracing, err := initTracing(v, applicationName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to build tracing component: %v \n", err)
		return
	}

	r := mux.NewRouter()

	// Auto instrumentation options of mux router.
	otelMuxOptions := []otelmux.Option{
		otelmux.WithPropagators(tracing.Propagator),
		otelmux.WithTracerProvider(tracing.TracerProvider),
	}

	// Auto instrumentation of http client.
	tr := otelhttp.NewTransport(http.DefaultTransport,
		otelhttp.WithPropagators(tracing.Propagator),
		otelhttp.WithTracerProvider(tracing.TracerProvider),
	)

	client := &http.Client{Transport: tr}

	requestHandlerFunc := func(w http.ResponseWriter, r *http.Request) {
		requestHandler(w, r, viper.GetString("postURL"), client)
	}

	r.Use(otelmux.Middleware("ServiceA", otelMuxOptions...), tracer.EchoFirstTraceNodeInfo(tracing.Propagator))
	r.HandleFunc("/", requestHandlerFunc)

	err1 := http.ListenAndServe(":"+v.GetString("port"), r)
	log.Fatal(err1)
}
