package errors

type State string

const (
	// NoError no error found
	NoError State = ""

	// DomainNotValid whois
	DomainNotValid State = "domain-not-valid"
	WhoisError     State = "whois-error"

	// UrlNotValid download
	UrlNotValid             State = "url-not-valid"
	RequestBodyNotValid     State = "request-body-not-valid"
	InternalError           State = "internal-error"
	InternalHttpClientError State = "internal-http-client-error"
	DialTimeout             State = "dial-timeout"
	Timeout                 State = "timeout"
	TlsTimeout              State = "tls-timeout"
	SanitizerError          State = "sanitize-error"
	DnsNotResolved          State = "dns-not-resolved"
	ConnectionRefused       State = "connection-refused"
	SysCallGenericError     State = "sys-call-generic-error"
	DnsTimeout              State = "dns-timeout"
	ClientNotSuccess        State = "client-not-success"
)
