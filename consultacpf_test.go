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

func Test_defaultClient_ConsultarCPF(t *testing.T) {
	type args struct {
		httpClient     func() *http.Client
		ctx            func() (context.Context, context.CancelFunc)
		cpf            string
		dataNascimento time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    soawebservices.PessoaFisica
		wantErr bool
	}{
		{
			name: "should return a Pessoa Física",
			args: args{
				httpClient: func() *http.Client {
					return &http.Client{
						Transport: RoundTripFunc(func(req *http.Request) *http.Response {
							resp := httptest.NewRecorder()
							resp.Body.Write(MustLoadTestDataFile(t, "consultacpf_success.json"))
							return resp.Result()
						}),
						Timeout: 5 * time.Second,
					}
				},
				ctx: func() (context.Context, context.CancelFunc) {
					return context.TODO(), nil
				},
				cpf:            "999.999.999-99",
				dataNascimento: time.Now(),
			},
			want: soawebservices.PessoaFisica{
				Documento:      "99999999999",
				Nome:           "DOCUMENTO CPF DE TESTE",
				NomeSocial:     "",
				DataNascimento: time.Time{},
			},
		},
		{
			name: "should fail due to the given invalid data de nascimento",
			args: args{
				httpClient: func() *http.Client {
					return &http.Client{
						Transport: RoundTripFunc(func(req *http.Request) *http.Response {
							resp := httptest.NewRecorder()
							resp.Body.Write(MustLoadTestDataFile(t, "consultacpf_invalid_data_nascimento.json"))
							return resp.Result()
						}),
						Timeout: 5 * time.Second,
					}
				},
				ctx: func() (context.Context, context.CancelFunc) {
					return context.TODO(), nil
				},
				cpf:            "999.999.999-99",
				dataNascimento: time.Now(),
			},
			want:    soawebservices.PessoaFisica{},
			wantErr: true,
		},
		{
			name: "should fail due to the given empty data de nascimento",
			args: args{
				httpClient: func() *http.Client {
					return &http.Client{
						Transport: RoundTripFunc(func(req *http.Request) *http.Response {
							resp := httptest.NewRecorder()
							resp.Body.Write(MustLoadTestDataFile(t, "consultacpf_required_data_nascimento.json"))
							return resp.Result()
						}),
						Timeout: 5 * time.Second,
					}
				},
				ctx: func() (context.Context, context.CancelFunc) {
					return context.TODO(), nil
				},
				cpf:            "999.999.999-99",
				dataNascimento: time.Now(),
			},
			want:    soawebservices.PessoaFisica{},
			wantErr: true,
		},
		{
			name: "should fail due to a invalid CPF",
			args: args{
				httpClient: func() *http.Client {
					return &http.Client{
						Transport: RoundTripFunc(func(req *http.Request) *http.Response {
							resp := httptest.NewRecorder()
							resp.Body.Write(MustLoadTestDataFile(t, "consultacpf_invalid_cpf.json"))
							return resp.Result()
						}),
						Timeout: 5 * time.Second,
					}
				},
				ctx: func() (context.Context, context.CancelFunc) {
					return context.TODO(), nil
				},
				cpf:            "999.999.999-99",
				dataNascimento: time.Now(),
			},
			want:    soawebservices.PessoaFisica{},
			wantErr: true,
		},
		{
			name: "should fail due to the wrong credentials",
			args: args{
				httpClient: func() *http.Client {
					return &http.Client{
						Transport: RoundTripFunc(func(req *http.Request) *http.Response {
							resp := httptest.NewRecorder()
							resp.Body.Write(MustLoadTestDataFile(t, "consultacpf_wrong_credentials.json"))
							return resp.Result()
						}),
						Timeout: 5 * time.Second,
					}
				},
				ctx: func() (context.Context, context.CancelFunc) {
					return context.TODO(), nil
				},
				cpf:            "999.999.999-99",
				dataNascimento: time.Now(),
			},
			want:    soawebservices.PessoaFisica{},
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
				cpf:            "999.999.999-99",
				dataNascimento: time.Now(),
			},
			want:    soawebservices.PessoaFisica{},
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
				cpf:            "999.999.999-99",
				dataNascimento: time.Now(),
			},
			want:    soawebservices.PessoaFisica{},
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
			result, err := client.ConsultarCPF(ctx, tt.args.cpf, tt.args.dataNascimento)
			if err != nil && !tt.wantErr {
				t.Error(err)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("ConsultarCPF() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(result, tt.want) {
				t.Error("want ", tt.want, " but got ", result)
			}
		})
	}
}
