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

func Test_defaultClient_ConsultarCNPJ(t *testing.T) {
	type args struct {
		httpClient func() *http.Client
		ctx        func() (context.Context, context.CancelFunc)
		cnpj       string
	}
	tests := []struct {
		name    string
		args    args
		want    soawebservices.PessoaJuridica
		wantErr bool
	}{
		{
			name: "should return a Pessoa Jurídica",
			args: args{
				httpClient: func() *http.Client {
					return &http.Client{
						Transport: RoundTripFunc(func(req *http.Request) *http.Response {
							resp := httptest.NewRecorder()
							resp.Body.Write(MustLoadTestDataFile(t, "consultacnpj_success.json"))
							return resp.Result()
						}),
						Timeout: 5 * time.Second,
					}
				},
				ctx: func() (context.Context, context.CancelFunc) {
					return context.TODO(), nil
				},
				cnpj: "99.999.999/9999-62",
			},
			want: soawebservices.PessoaJuridica{
				Documento:    "99999999999962",
				RazaoSocial:  "DOCUMENTO CNPJ DE TESTES",
				NomeFantasia: "EMPRESA DE TESTES",
				DataFundacao: time.Date(2007, 05, 02, 0, 0, 0, 0, time.UTC),
				Matriz:       true,
				CNAE: soawebservices.CNAE{
					Codigo:    "82.91-1-00",
					Descricao: "Atividades de cobranças e informações cadastrais",
				},
				NaturezaJuridica: soawebservices.NaturezaJuridica{
					Codigo:    "206-2",
					Descricao: "SOCIEDADE EMPRESARIA LIMITADA",
				},
				Email:    "email@email.com",
				Telefone: "1199999999",
			},
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
				cnpj: "99.999.999/9999-62",
			},
			want:    soawebservices.PessoaJuridica{},
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
				cnpj: "99.999.999/9999-62",
			},
			want:    soawebservices.PessoaJuridica{},
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
				cnpj: "99.999.999/9999-62",
			},
			want:    soawebservices.PessoaJuridica{},
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
				cnpj: "99.999.999/9999-62",
			},
			want:    soawebservices.PessoaJuridica{},
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
			result, err := client.ConsultarCNPJ(ctx, tt.args.cnpj)
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
