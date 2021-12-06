package soawebservices

import (
	"context"
	"net/http"
	"time"
)

type Ambiente string

const (
	Producao  Ambiente = "producao"
	TestDrive Ambiente = "test-drive"
)

type CEPService interface {
	ConsultarCEP(ctx context.Context, cep string) (CEP, error)
}

type PessoaFisicaService interface {
	ConsultarCPF(ctx context.Context, cpf string, dataNascimento time.Time) (PessoaFisica, error)
}

type PessoaJuridicaService interface {
	ConsultarCNPJ(ctx context.Context, cnpj string) (PessoaJuridica, error)
}

type Client interface {
	CEPService
	PessoaFisicaService
	PessoaJuridicaService
}

type defaultClient struct {
	httpClient  *http.Client
	baseURL     string
	ambiente    Ambiente
	credenciais Credenciais
}

func NewClient(httpClient *http.Client, baseURL string, ambiente Ambiente, credenciais Credenciais) Client {
	return &defaultClient{
		httpClient:  httpClient,
		baseURL:     baseURL,
		ambiente:    ambiente,
		credenciais: credenciais,
	}
}
