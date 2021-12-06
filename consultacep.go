package soawebservices

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

const (
	urlCEP = "cep/cep.asmx"
)

const (
	ErrCEPInvalido            = Error("o cep informado é inválido (P016M001)")
	ErrCEPFalhaTransacao      = Error("ocorreu um erro e não foi possível realizar a consulta (P016M009)")
	ErrCEPServicoIndisponivel = Error("serviço indisponível no momento (P016M010)")
)

const (
	StatusCEPInvalido            = "P016M001"
	StatusCEPFalhaProcessamento  = "P016M009"
	StatusCEPServicoIndisponivel = "P016M010"
	StatusCEPNaoEncontrado       = "P016M002"
)

func (d *defaultClient) buildConsultaCEPRequestBody(cep string) (io.Reader, error) {
	consultaCep := newConsultaCEPEstendida(d.credenciais, cep)
	buf, err := xml.Marshal(newRequestEnvelope(consultaCep))
	if err != nil {
		return nil, fmt.Errorf("an error occurred while build the request body: %w", err)
	}
	return bytes.NewBuffer(buf), err
}

func (d *defaultClient) parseConsultaCEPResponseBody(resp *http.Response) (CEP, error) {
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	var respBody struct {
		XMLName xml.Name
		Body    struct {
			XMLName  xml.Name
			Response struct {
				XMLName xml.Name `xml:"ConsultaCEPEstendidaResponse"`
				Result  consultaCEPEstendidaResult
			}
		}
	}
	if err := xml.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return CEP{}, err
	}
	result := respBody.Body.Response.Result
	if result.Transacao.CodigoStatus == StatusCEPInvalido {
		return CEP{}, ErrCEPInvalido
	}
	if result.Transacao.CodigoStatus == StatusCEPFalhaProcessamento {
		return CEP{}, ErrCEPFalhaTransacao
	}
	if result.Transacao.CodigoStatus == StatusCEPServicoIndisponivel {
		return CEP{}, ErrCEPServicoIndisponivel
	}
	if result.Transacao.CodigoStatus == StatusCredenciaisInvalidas {
		return CEP{}, ErrCredenciaisInvalidas
	}
	if result.Transacao.CodigoStatus == StatusCEPNaoEncontrado {
		return CEP{}, nil
	}
	if !result.Status {
		return CEP{}, fmt.Errorf("%s: %s", result.Transacao.CodigoStatus, result.Transacao.CodigoStatusDescricao)
	}
	return CEP{
		CEP:                   result.Cep,
		UF:                    result.UF,
		TipoLogradouro:        result.TipoLogradouro,
		LogradouroCompleto:    result.LogradouroCompleto,
		LogradouroComplemento: result.LogradouroComplemento,
		Bairro:                result.Bairro,
		Cidade:                result.Cidade,
		CodigoIBGE:            result.CodigoIBGE,
	}, nil
}

func (d *defaultClient) ConsultarCEP(ctx context.Context, cep string) (CEP, error) {
	errChan := make(chan error, 1)
	resultChan := make(chan CEP, 1)
	requestBody, err := d.buildConsultaCEPRequestBody(cep)
	if err != nil {
		return CEP{}, err
	}
	go func() {
		var resp *http.Response
		serviceURL := fmt.Sprintf("%s/webservices/%s/%s", d.baseURL, d.ambiente, urlCEP)
		resp, err = d.httpClient.Post(serviceURL, "text/xml", requestBody)
		if err != nil {
			errChan <- err
			return
		}
		result, err := d.parseConsultaCEPResponseBody(resp)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- result
	}()
	select {
	case err = <-errChan:
		return CEP{}, err
	case <-ctx.Done():
		return CEP{}, ctx.Err()
	case result := <-resultChan:
		return result, nil
	}
}
