package model

type ErrorState string

const (
	// DomainNotValid whois
	DomainNotValid ErrorState = "domain-not-valid"
	WhoisError     ErrorState = "whois-error"

	// UrlNotValid download
	UrlNotValid             ErrorState = "url-not-valid"
	RequestBodyNotValid     ErrorState = "request-body-not-valid"
	InternalError           ErrorState = "internal-error"
	InternalHttpClientError ErrorState = "internal-http-client-error"
	DialTimeout             ErrorState = "dial-timeout"
	Timeout                 ErrorState = "timeout"
	TlsTimeout              ErrorState = "tls-timeout"
	SanitizerError          ErrorState = "sanitize-error"
	DnsNotResolved          ErrorState = "dns-not-resolved"
	ConnectionRefused       ErrorState = "connection-refused"
	SysCallGenericError     ErrorState = "sys-call-generic-error"
	DnsTimeout              ErrorState = "dns-timeout"
	ClientNotSuccess        ErrorState = "client-not-success"
)

func (e ErrorState) Error() *Error {
	return &Error{
		ErrorState: e,
	}
}

func (e ErrorState) ErrorWithMsg(msg string) *Error {
	return &Error{
		ErrorState: e,
	}
}

func (e ErrorState) ErrorWithError(err error) *Error {
	if err == nil {
		return nil
	}

	return &Error{
		ErrorState: e,
	}
}

type Error struct {
	ErrorState ErrorState
}
