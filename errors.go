package soawebservices

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	StatusCredenciaisInvalidas = "G000M000"
	StatusDocumentoInvalido    = "G000M003"
)

const (
	ErrCredenciaisInvalidas = Error("credenciais inválidas (G000M000)")
)
