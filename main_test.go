package main

import (
	"ServiceA/tracer"
	"net/http"
	"reflect"
	"testing"

	"github.com/spf13/viper"
)

func Test_requestHandler(t *testing.T) {
	type args struct {
		w       http.ResponseWriter
		r       *http.Request
		postURL string
		tracing tracer.Tracing
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestHandler(tt.args.w, tt.args.r, tt.args.postURL, tt.args.tracing)
		})
	}
}

func Test_getViper(t *testing.T) {
	tests := []struct {
		name string
		want *viper.Viper
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getViper(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getViper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_initTracing(t *testing.T) {
	type args struct {
		v       *viper.Viper
		appName string
	}
	tests := []struct {
		name    string
		args    args
		want    tracer.Tracing
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := initTracing(tt.args.v, tt.args.appName)
			if (err != nil) != tt.wantErr {
				t.Errorf("initTracing() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initTracing() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
