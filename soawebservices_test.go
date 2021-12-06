package soawebservices_test

import (
	"fmt"
	"github.com/diegohordi/soawebservices"
	"net/http"
	"os"
	"testing"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func MustCreateClient(httpClient *http.Client) soawebservices.Client {
	credenciais := soawebservices.Credenciais{Email: "test@test.com", Senha: "test"}
	return soawebservices.NewClient(httpClient, "https://soawebservices.com.br", soawebservices.TestDrive, credenciais)
}

func MustLoadTestDataFile(t *testing.T, fileName string) []byte {
	t.Helper()
	content, err := os.ReadFile(fmt.Sprintf("./test/testdata/%s", fileName))
	if err != nil {
		t.Fatal(err)
	}
	return content
}
