// Package errx fornece utilitários para padronizar e enriquecer erros
// de aplicação, incluindo códigos, detalhes adicionais e informações
// do local onde o erro ocorreu (caller). Também oferece funções para
// conversão entre códigos da aplicação e status HTTP.
package errx

import (
	"errors"
	"fmt"
	"runtime"
)

type Code string

const (
	INTERNAL     Code = "INTERNAL"
	BAD_REQUEST  Code = "BAD_REQUEST"
	UNAUTHORIZED Code = "UNAUTHORIZED"
	FORBIDDEN    Code = "FORBIDDEN"
	NOT_FOUND    Code = "NOT_FOUND"
	CONFLICT     Code = "CONFLICT"
)

var (
	// separator é o separador usado na representação em string dos erros encadeados.
	separator = "->"
)

// AppError representa um erro da aplicação com metadados adicionais
// que facilitam log, telemetria e respostas HTTP.
type AppError struct {
	// Message é a mensagem principal do erro.
	Message string
	// Code identifica a categoria do erro (ex.: BAD_REQUEST, INTERNAL).
	Code Code
	// Err é o erro interno encadeado (causa raiz ou erro anterior).
	Err error
	// Caller descreve o ponto de origem do erro (arquivo:linha e função).
	Caller string
	// Details contém informações adicionais úteis para diagnóstico.
	Details map[string]interface{}
}

// ShowLogger define a estrutura padronizada para exibição/serialização
// segura do erro em logs e respostas, evitando expor tipos inesperados.
type ShowLogger struct {
	Message string                 `json:"message"`
	Code    Code                   `json:"code"`
	Caller  string                 `json:"caller,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// New cria uma nova AppError com a mensagem fornecida.
func New(message string) *AppError {
	return &AppError{
		Message: message,
	}
}

func SetSeparator(sep string) {
	separator = sep
}

// WithCode define o código do erro.
func (e *AppError) WithCode(code Code) *AppError {
	if e.Code != "" {
		return e
	}

	e.Code = code
	return e
}

// WithError associa um erro interno à AppError. Caso o erro recebido
// já seja uma AppError, herdará metadados úteis (code, caller, details)
// sem sobrescrever valores já definidos neste erro.
func (e *AppError) WithError(err error) *AppError {
	if err == nil {
		return e
	}

	if inner, ok := asAppError(err); ok {

		if e.Code == "" && inner.Code != "" {
			e.Code = inner.Code
		}

		if e.Caller == "" && inner.Caller != "" {
			e.Caller = inner.Caller
		}

		if len(inner.Details) > 0 {
			if e.Details == nil {
				e.Details = make(map[string]interface{}, len(inner.Details))
			}
			// Herdar apenas chaves ausentes para não perder informações já atribuídas
			for k, v := range inner.Details {
				if _, exists := e.Details[k]; !exists {
					e.Details[k] = v
				}
			}
		}
	}

	e.Err = err
	return e
}

// WithCaller adiciona informações do caller à AppError (arquivo, linha e função).
func (e *AppError) WithCaller() *AppError {
	if e.Caller != "" {
		return e
	}

	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	e.Caller = fmt.Sprintf("[%s:%d %s]", file, line, fn.Name())
	return e
}

// WithDetails adiciona detalhes extras à AppError (mesclando com os existentes).
func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	for k, v := range details {
		e.Details[k] = v
	}
	return e
}

// Error implementa a interface error, incluindo a cadeia de erros, quando houver.
func (e *AppError) Error() string {
	if e.Err == nil {
		return e.Message
	}
	return fmt.Sprintf("%s %s %s", e.Message, separator, e.Err)
}

// Unwrap permite integrar com erros do Go (errors.Unwrap/As/Is),
// expondo o erro interno encadeado.
func (e *AppError) Unwrap() error { return e.Err }

// Funções auxiliares
// IsAppError informa se o erro ou algum erro interno é uma AppError.
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// GetAppError retorna a AppError contida em err, se existir; caso contrário, nil.
func GetAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return nil
}

// GetCode extrai o Code de uma AppError; caso contrário, retorna INTERNAL.
func GetCode(err error) Code {
	var appErr *AppError
	if errors.As(err, &appErr) && appErr.Code != "" {
		return appErr.Code
	}
	return INTERNAL
}

// GetStatusCode extrai o status HTTP correspondente ao Code da AppError; caso contrário, retorna 500.
func GetStatusCode(err error) int {
	var appErr *AppError
	if errors.As(err, &appErr) && appErr.Code != "" {
		return ToHTTPCode(appErr.Code)
	}

	return 500
}

// GetCaller extrai a informação do caller da AppError, se houver.
func GetCaller(err error) string {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Caller
	}
	return ""
}

// GetDetails retorna o mapa de detalhes da AppError, se houver.
func GetDetails(err error) map[string]interface{} {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Details
	}
	return nil
}

// PrintLogger devolve uma estrutura pronta para log/serialização do erro.
// ToHTTPCode converte um Code de erro para o status HTTP correspondente.
func ToHTTPCode(code Code) int {
	switch code {
	case BAD_REQUEST:
		return 400
	case UNAUTHORIZED:
		return 401
	case FORBIDDEN:
		return 403
	case NOT_FOUND:
		return 404
	case CONFLICT:
		return 409
	default:
		return 500
	}

}

// StatusToCode converte um status HTTP para o Code de erro correspondente.
func StatusToCode(status int) Code {
	if status >= 500 {
		return INTERNAL
	}

	switch status {
	case 401:
		return UNAUTHORIZED
	case 403:
		return FORBIDDEN
	case 404:
		return NOT_FOUND
	case 409:
		return CONFLICT
	default:
		return BAD_REQUEST
	}

}

// PrintLogger devolve uma estrutura pronta para log/serialização do erro.
func PrintLogger(err error) *ShowLogger {
	if !IsAppError(err) {
		return nil
	}
	appErr := GetAppError(err)
	return &ShowLogger{
		Message: appErr.Error(),
		Code:    appErr.Code,
		Caller:  appErr.Caller,
		Details: appErr.Details,
	}
}

// PrintHttpLogger retorna o status HTTP e o payload serializável do erro.
func PrintHttpLogger(err error) (int, *ShowLogger) {
	if !IsAppError(err) {
		return 500, nil
	}
	appErr := GetAppError(err)
	return ToHTTPCode(appErr.Code), &ShowLogger{
		Code:    appErr.Code,
		Message: appErr.Error(),
		Details: appErr.Details,
	}
}

// asAppError tenta obter *AppError a partir de um error qualquer.
func asAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
