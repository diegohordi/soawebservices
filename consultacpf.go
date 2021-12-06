package soawebservices

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	urlCPF = "cdc/pessoafisicanfe.ashx"
)

const (
	StatusDataNascimentoObrigatoria = "P009M001"
	StatusDataNascimentoInvalida    = "P009M002"
)

const (
	ErrDataNascimentoObrigatoria = Error("data de nascimento obrigatória (P009M001)")
	ErrDataNascimentoInvalida    = Error("data de nascimento inválida (P009M002)")
	ErrCPFInvalido               = Error("o cnpj informado é inválido (G000M003)")
)

func (d *defaultClient) buildConsultaCPFRequestBody(cpf string, dataNascimento time.Time) (io.Reader, error) {
	consultaCPF := newConsultaPessoaFisicaNFe(d.credenciais, cpf, dataNascimento.Format("02/01/2006"))
	buf, err := json.Marshal(consultaCPF)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while build the request body: %w", err)
	}
	return bytes.NewBuffer(buf), err
}

func (d *defaultClient) parseConsultaCPFResponseBody(resp *http.Response) (PessoaFisica, error) {
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	result := pessoaFisicaResult{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return PessoaFisica{}, err
	}
	if result.Transacao.CodigoStatus == StatusDataNascimentoInvalida {
		return PessoaFisica{}, ErrDataNascimentoInvalida
	}
	if result.Transacao.CodigoStatus == StatusDataNascimentoObrigatoria {
		return PessoaFisica{}, ErrDataNascimentoObrigatoria
	}
	if result.Transacao.CodigoStatus == StatusDocumentoInvalido {
		return PessoaFisica{}, ErrCPFInvalido
	}
	if result.Transacao.CodigoStatus == StatusCredenciaisInvalidas {
		return PessoaFisica{}, ErrCredenciaisInvalidas
	}
	if !result.Status {
		return PessoaFisica{}, fmt.Errorf("%s: %s", result.Transacao.CodigoStatus, result.Transacao.CodigoStatusDescricao)
	}
	dataNascimento, _ := time.Parse("02/01/2006", result.DataNascimento)
	return PessoaFisica{
		Documento:      result.Documento,
		Nome:           result.Nome,
		NomeSocial:     result.NomeSocial,
		DataNascimento: dataNascimento,
		Status:         PessoaFisicaStatus(result.CodigoSituacaoCadastral),
	}, nil
}

func (d *defaultClient) ConsultarCPF(ctx context.Context, cpf string, dataNascimento time.Time) (PessoaFisica, error) {
	errChan := make(chan error, 1)
	resultChan := make(chan PessoaFisica, 1)
	requestBody, err := d.buildConsultaCPFRequestBody(cpf, dataNascimento)
	if err != nil {
		return PessoaFisica{}, err
	}
	go func() {
		var resp *http.Response
		serviceURL := fmt.Sprintf("%s/restservices/%s/%s", d.baseURL, d.ambiente, urlCPF)
		resp, err = d.httpClient.Post(serviceURL, "application/json", requestBody)
		if err != nil {
			errChan <- err
			return
		}
		result, err := d.parseConsultaCPFResponseBody(resp)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- result
	}()
	select {
	case err = <-errChan:
		return PessoaFisica{}, err
	case <-ctx.Done():
		return PessoaFisica{}, ctx.Err()
	case result := <-resultChan:
		return result, nil
	}
}
