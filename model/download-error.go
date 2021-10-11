package model

type DownloadErrorState string

const (
	UrlNotValid             DownloadErrorState = "url-not-valid"
	RequestBodyNotValid     DownloadErrorState = "request-body-not-valid"
	InternalError           DownloadErrorState = "internal-error"
	InternalHttpClientError DownloadErrorState = "internal-http-client-error"
	Timeout                 DownloadErrorState = "timeout"
	SanitizerError          DownloadErrorState = "sanitize-error"
	DnsNotResolved                             = "dns-not-resolved"
	ConnectionRefused                          = "connection-refused"
	SysCallGenericError                        = "sys-call-generic-error"
	DnsTimeout                                 = "dns-timeout"
	ClientNotSuccess                           = "client-not-success"
)

type DownloadError struct {
	ErrorState   DownloadErrorState
	ErrorText    string
	ClientStatus int
}
