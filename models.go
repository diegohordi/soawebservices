package soawebservices

import (
	"encoding/xml"
	"time"
)

const (
	defaultXSI       = "http://www.w3.org/2001/XMLSchema-instance"
	defaultXSD       = "http://www.w3.org/2001/XMLSchema"
	defaultSOAP      = "http://schemas.xmlsoap.org/soap/envelope/"
	defaultNamespace = "SOAWebServices"
)

type body struct {
	XMLName xml.Name
	Body    interface{}
}

type requestEnvelope struct {
	XMLName xml.Name `xml:"soap:Envelope"`
	XSI     string   `xml:"xmlns:xsi,attr"`
	XSD     string   `xml:"xmlns:xsd,attr"`
	Soap    string   `xml:"xmlns:soap,attr"`
	Body    *body    `xml:"soap:Body"`
}

func newRequestEnvelope(bodyRef interface{}) *requestEnvelope {
	return &requestEnvelope{
		XSI:  defaultXSI,
		XSD:  defaultXSD,
		Soap: defaultSOAP,
		Body: &body{Body: bodyRef},
	}
}

type transacao struct {
	Status                bool   `xml:"Status" json:"status"`
	CodigoStatus          string `xml:"CodigoStatus" json:"CodigoStatus"`
	CodigoStatusDescricao string `xml:"CodigoStatusDescricao" json:"CodigoStatusDescricao"`
}

type Credenciais struct {
	Email string `xml:"Email"`
	Senha string `xml:"Senha"`
}

type consultaCEPEstendida struct {
	XMLName     xml.Name    `xml:"ConsultaCEPEstendida"`
	Namespace   string      `xml:"xmlns,attr"`
	Credenciais Credenciais `xml:"Credenciais"`
	Cep         string      `xml:"CEP"`
}

func newConsultaCEPEstendida(credenciais Credenciais, cep string) *consultaCEPEstendida {
	return &consultaCEPEstendida{
		Namespace:   defaultNamespace,
		Credenciais: credenciais,
		Cep:         cep,
	}
}

type consultaCEPEstendidaResult struct {
	XMLName               xml.Name  `xml:"ConsultaCEPEstendidaResult"`
	Cep                   string    `xml:"CEP"`
	UF                    string    `xml:"UF"`
	TipoLogradouro        string    `xml:"TipoLogradouro"`
	LogradouroCompleto    string    `xml:"LogradouroCompleto"`
	LogradouroComplemento string    `xml:"LogradouroComplemento"`
	Bairro                string    `xml:"Bairro"`
	Cidade                string    `xml:"Cidade"`
	CodigoIBGE            string    `xml:"CodigoIBGE"`
	Mensagem              string    `xml:"Mensagem"`
	Status                bool      `xml:"Status"`
	Transacao             transacao `xml:"Transacao"`
}

type CEP struct {
	CEP                   string
	UF                    string
	TipoLogradouro        string
	LogradouroCompleto    string
	LogradouroComplemento string
	Bairro                string
	Cidade                string
	CodigoIBGE            string
}

type consultaPessoaFisicaNFe struct {
	Credenciais    Credenciais `json:"Credenciais"`
	Documento      string      `json:"Documento"`
	DataNascimento string      `json:"DataNascimento"`
}

func newConsultaPessoaFisicaNFe(credenciais Credenciais, documento string, dataNascimento string) *consultaPessoaFisicaNFe {
	return &consultaPessoaFisicaNFe{
		Credenciais:    credenciais,
		Documento:      documento,
		DataNascimento: dataNascimento,
	}
}

type pessoaFisicaResult struct {
	Documento               string    `json:"Documento"`
	Nome                    string    `json:"Nome"`
	NomeSocial              string    `json:"NomeSocial"`
	DataNascimento          string    `json:"DataNascimento"`
	DataInscricao           string    `json:"DataInscricao"`
	AnoObito                string    `json:"AnoObito"`
	MensagemObito           string    `json:"MensagemObito"`
	CodigoSituacaoCadastral string    `json:"CodigoSituacaoCadastral"`
	SituacaoRFB             string    `json:"SituacaoRFB"`
	DataConsultaRFB         time.Time `json:"DataConsultaRFB"`
	ProtocoloRFB            string    `json:"ProtocoloRFB"`
	DigitoVerificador       string    `json:"DigitoVerificador"`
	DIRPF                   string    `json:"DIRPF"`
	Mensagem                string    `json:"Mensagem"`
	Status                  bool      `json:"Status"`
	Transacao               transacao `json:"Transacao"`
}

type PessoaFisicaStatus string

const (
	Regular                      PessoaFisicaStatus = "1"
	Suspensa                     PessoaFisicaStatus = "2"
	TitularFalecido              PessoaFisicaStatus = "3"
	CanceladaPorMultiplicidade   PessoaFisicaStatus = "4"
	PendenteRegularizacao        PessoaFisicaStatus = "5"
	CanceladaOficio              PessoaFisicaStatus = "6"
	CanceladaEncerramento        PessoaFisicaStatus = "7"
	Cancelada                    PessoaFisicaStatus = "8"
	Nula                         PessoaFisicaStatus = "9"
	SituacaoCadastralInexistente PessoaFisicaStatus = "12"
	DadosIncompletos             PessoaFisicaStatus = "13"
)

type PessoaFisica struct {
	Documento      string
	Nome           string
	NomeSocial     string
	DataNascimento time.Time
	Status         PessoaFisicaStatus
}

type consultaPessoaJuridicaNFe struct {
	Credenciais Credenciais `json:"Credenciais"`
	Documento   string      `json:"Documento"`
}

func newConsultaPessoaJuridicaNFe(credenciais Credenciais, documento string) *consultaPessoaJuridicaNFe {
	return &consultaPessoaJuridicaNFe{Credenciais: credenciais, Documento: documento}
}

type pessoaJuridicaResult struct {
	Documento                         string    `json:"Documento"`
	RazaoSocial                       string    `json:"RazaoSocial"`
	NomeFantasia                      string    `json:"NomeFantasia"`
	DataFundacao                      string    `json:"DataFundacao"`
	MatrizFilial                      string    `json:"MatrizFilial"`
	Capital                           string    `json:"Capital"`
	CodigoAtividadeEconomica          string    `json:"CodigoAtividadeEconomica"`
	CodigoAtividadeEconomicaDescricao string    `json:"CodigoAtividadeEconomicaDescricao"`
	CodigoNaturezaJuridica            string    `json:"CodigoNaturezaJuridica"`
	CodigoNaturezaJuridicaDescricao   string    `json:"CodigoNaturezaJuridicaDescricao"`
	SituacaoRFB                       string    `json:"SituacaoRFB"`
	DataSituacaoRFB                   string    `json:"DataSituacaoRFB"`
	DataConsultaRFB                   string    `json:"DataConsultaRFB"`
	MotivoSituacaoRFB                 string    `json:"MotivoSituacaoRFB"`
	DataMotivoEspecialSituacaoRFB     string    `json:"DataMotivoEspecialSituacaoRFB"`
	Email                             string    `json:"Email"`
	Telefone                          string    `json:"Telefone"`
	Mensagem                          string    `json:"Mensagem"`
	Status                            bool      `json:"Status"`
	Transacao                         transacao `json:"Transacao"`
}

type CNAE struct {
	Codigo    string
	Descricao string
}

type NaturezaJuridica struct {
	Codigo    string
	Descricao string
}

type PessoaJuridica struct {
	Documento        string
	RazaoSocial      string
	NomeFantasia     string
	DataFundacao     time.Time
	Matriz           bool
	CNAE             CNAE
	NaturezaJuridica NaturezaJuridica
	Email            string
	Telefone         string
}
