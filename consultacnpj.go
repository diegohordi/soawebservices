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
	urlCNPJ = "cdc/pessoajuridicanfe.ashx"
)

const (
	ErrCNPJInvalido = Error("o cnpj informado é inválido (G000M003)")
)

func (d *defaultClient) buildConsultaCNPJRequestBody(cnpj string) (io.Reader, error) {
	consultaCNPJ := newConsultaPessoaJuridicaNFe(d.credenciais, cnpj)
	buf, err := json.Marshal(consultaCNPJ)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while build the request body: %w", err)
	}
	return bytes.NewBuffer(buf), err
}

func (d *defaultClient) parseConsultaCNPJResponseBody(resp *http.Response) (PessoaJuridica, error) {
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	result := pessoaJuridicaResult{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return PessoaJuridica{}, err
	}
	if result.Transacao.CodigoStatus == StatusDocumentoInvalido {
		return PessoaJuridica{}, ErrCNPJInvalido
	}
	if result.Transacao.CodigoStatus == StatusCredenciaisInvalidas {
		return PessoaJuridica{}, ErrCredenciaisInvalidas
	}
	if !result.Status {
		return PessoaJuridica{}, fmt.Errorf("%s: %s", result.Transacao.CodigoStatus, result.Transacao.CodigoStatusDescricao)
	}
	dataFundacao, _ := time.Parse("02/01/2006", result.DataFundacao)
	return PessoaJuridica{
		Documento:    result.Documento,
		RazaoSocial:  result.RazaoSocial,
		NomeFantasia: result.NomeFantasia,
		DataFundacao: dataFundacao,
		Matriz:       result.MatrizFilial == "MATRIZ",
		CNAE: CNAE{
			Codigo:    result.CodigoAtividadeEconomica,
			Descricao: result.CodigoAtividadeEconomicaDescricao,
		},
		NaturezaJuridica: NaturezaJuridica{
			Codigo:    result.CodigoNaturezaJuridica,
			Descricao: result.CodigoNaturezaJuridicaDescricao,
		},
		Email:    result.Email,
		Telefone: result.Telefone,
	}, nil
}

func (d *defaultClient) ConsultarCNPJ(ctx context.Context, cnpj string) (PessoaJuridica, error) {
	errChan := make(chan error, 1)
	resultChan := make(chan PessoaJuridica, 1)
	requestBody, err := d.buildConsultaCNPJRequestBody(cnpj)
	if err != nil {
		return PessoaJuridica{}, err
	}
	go func() {
		var resp *http.Response
		serviceURL := fmt.Sprintf("%s/restservices/%s/%s", d.baseURL, d.ambiente, urlCNPJ)
		resp, err = d.httpClient.Post(serviceURL, "application/json", requestBody)
		if err != nil {
			errChan <- err
			return
		}
		result, err := d.parseConsultaCNPJResponseBody(resp)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- result
	}()
	select {
	case err = <-errChan:
		return PessoaJuridica{}, err
	case <-ctx.Done():
		return PessoaJuridica{}, ctx.Err()
	case result := <-resultChan:
		return result, nil
	}
}
