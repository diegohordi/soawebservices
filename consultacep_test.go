package soawebservices_test

import (
	"context"
	"github.com/diegohordi/soawebservices"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func Test_defaultClient_ConsultarCEP(t *testing.T) {
	type args struct {
		httpClient func() *http.Client
		ctx        func() (context.Context, context.CancelFunc)
		cep        string
	}
	tests := []struct {
		name    string
		args    args
		want    soawebservices.CEP
		wantErr bool
	}{
		{
			name: "should return a CEP",
			args: args{
				httpClient: func() *http.Client {
					return &http.Client{
						Transport: RoundTripFunc(func(req *http.Request) *http.Response {
							resp := httptest.NewRecorder()
							resp.Body.Write(MustLoadTestDataFile(t, "consultacep_success.xml"))
							return resp.Result()
						}),
						Timeout: 5 * time.Second,
					}
				},
				ctx: func() (context.Context, context.CancelFunc) {
					return context.TODO(), nil
				},
				cep: "99.999999",
			},
			want: soawebservices.CEP{
				CEP:                   "99999999",
				UF:                    "XX",
				TipoLogradouro:        "RUA",
				LogradouroCompleto:    "RUA LOGRADOURO DE TESTES",
				LogradouroComplemento: "LOGRADOURO COMPLEMENTO TESTES",
				Bairro:                "BAIRRO DE TESTES",
				Cidade:                "CIDADE DE TESTES",
				CodigoIBGE:            "99999",
			},
		},
		{
			name: "should return an empty CEP",
			args: args{
				httpClient: func() *http.Client {
					return &http.Client{
						Transport: RoundTripFunc(func(req *http.Request) *http.Response {
							resp := httptest.NewRecorder()
							resp.Body.Write(MustLoadTestDataFile(t, "consultacep_not_found.xml"))
							return resp.Result()
						}),
						Timeout: 5 * time.Second,
					}
				},
				ctx: func() (context.Context, context.CancelFunc) {
					return context.TODO(), nil
				},
				cep: "12345123",
			},
			want:    soawebservices.CEP{},
			wantErr: false,
		},
		{
			name: "should fail due to the given invalid CEP",
			args: args{
				httpClient: func() *http.Client {
					return &http.Client{
						Transport: RoundTripFunc(func(req *http.Request) *http.Response {
							resp := httptest.NewRecorder()
							resp.Body.Write(MustLoadTestDataFile(t, "consultacep_invalid_cep.xml"))
							return resp.Result()
						}),
						Timeout: 5 * time.Second,
					}
				},
				ctx: func() (context.Context, context.CancelFunc) {
					return context.TODO(), nil
				},
				cep: "99.999999",
			},
			want:    soawebservices.CEP{},
			wantErr: true,
		},
		{
			name: "should fail due to service unavailability",
			args: args{
				httpClient: func() *http.Client {
					return &http.Client{
						Transport: RoundTripFunc(func(req *http.Request) *http.Response {
							resp := httptest.NewRecorder()
							resp.Body.Write(MustLoadTestDataFile(t, "consultacep_service_unavailable.xml"))
							return resp.Result()
						}),
						Timeout: 5 * time.Second,
					}
				},
				ctx: func() (context.Context, context.CancelFunc) {
					return context.TODO(), nil
				},
				cep: "99.999999",
			},
			want:    soawebservices.CEP{},
			wantErr: true,
		},
		{
			name: "should fail due to an unknown server error",
			args: args{
				httpClient: func() *http.Client {
					return &http.Client{
						Transport: RoundTripFunc(func(req *http.Request) *http.Response {
							resp := httptest.NewRecorder()
							resp.Body.Write(MustLoadTestDataFile(t, "consultacep_unknown_server_error.xml"))
							return resp.Result()
						}),
						Timeout: 5 * time.Second,
					}
				},
				ctx: func() (context.Context, context.CancelFunc) {
					return context.TODO(), nil
				},
				cep: "99.999999",
			},
			want:    soawebservices.CEP{},
			wantErr: true,
		},
		{
			name: "should fail due to the given wrong credentials",
			args: args{
				httpClient: func() *http.Client {
					return &http.Client{
						Transport: RoundTripFunc(func(req *http.Request) *http.Response {
							resp := httptest.NewRecorder()
							resp.Body.Write(MustLoadTestDataFile(t, "consultacep_wrong_credentials.xml"))
							return resp.Result()
						}),
						Timeout: 5 * time.Second,
					}
				},
				ctx: func() (context.Context, context.CancelFunc) {
					return context.TODO(), nil
				},
				cep: "99.999999",
			},
			want:    soawebservices.CEP{},
			wantErr: true,
		},
		{
			name: "should fail due to context timeout",
			args: args{
				httpClient: func() *http.Client {
					return &http.Client{
						Transport: RoundTripFunc(func(req *http.Request) *http.Response {
							time.Sleep(10 * time.Second)
							return &http.Response{}
						}),
					}
				},
				ctx: func() (context.Context, context.CancelFunc) {
					return context.WithTimeout(context.TODO(), 1*time.Millisecond)
				},
				cep: "99.999999",
			},
			want:    soawebservices.CEP{},
			wantErr: true,
		},
		{
			name: "should fail due to client timeout",
			args: args{
				httpClient: func() *http.Client {
					return &http.Client{
						Timeout: 1 * time.Millisecond,
					}
				},
				ctx: func() (context.Context, context.CancelFunc) {
					return context.TODO(), nil
				},
				cep: "12345-123",
			},
			want:    soawebservices.CEP{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client := MustCreateClient(tt.args.httpClient())
			ctx, cancel := tt.args.ctx()
			if cancel != nil {
				defer cancel()
			}
			result, err := client.ConsultarCEP(ctx, tt.args.cep)
			if err != nil && !tt.wantErr {
				t.Error(err)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("ConsultarCEP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(result, tt.want) {
				t.Error("want ", tt.want, " but got ", result)
			}
		})
	}
}
