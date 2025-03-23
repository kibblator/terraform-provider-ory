package email_configuration_resource

import (
	"fmt"
	"net/url"

	orytypes "github.com/kibblator/terraform-provider-ory/internal/provider/types"
)

func HttpConfigToApi(tfConfig emailConfigurationResourceModel, httpConfig *orytypes.HTTP) {
	headersMap := make(map[string]string)

	for _, header := range *tfConfig.SMTPHeaders {
		headersMap[header.Key.ValueString()] = header.Value.ValueString()
	}

	authenticationType := tfConfig.HTTPConfig.AuthenticationType.ValueString()

	httpConfig.HttpRequestConfig = &orytypes.HttpRequestConfig{}

	if authenticationType != "none" {
		httpConfig.HttpRequestConfig.HttpAuth = &orytypes.HttpAuth{
			HttpAuthConfig: &orytypes.HttpAuthConfig{},
		}

		httpConfig.HttpRequestConfig.HttpAuth.Type = authenticationType

		if authenticationType == "api_key" {
			httpConfig.HttpRequestConfig.HttpAuth.HttpAuthConfig.In = tfConfig.HTTPConfig.ApiKey.TransportMode.ValueString()
			httpConfig.HttpRequestConfig.HttpAuth.HttpAuthConfig.Name = tfConfig.HTTPConfig.ApiKey.Name.ValueString()
			httpConfig.HttpRequestConfig.HttpAuth.HttpAuthConfig.Value = tfConfig.HTTPConfig.ApiKey.Value.ValueString()
		}

		if authenticationType == "basic_auth" {
			httpConfig.HttpRequestConfig.HttpAuth.HttpAuthConfig.User = tfConfig.HTTPConfig.BasicAuth.Username.ValueString()
			httpConfig.HttpRequestConfig.HttpAuth.HttpAuthConfig.Password = tfConfig.HTTPConfig.BasicAuth.Password.ValueString()
		}
	}

	httpConfig.HttpRequestConfig.Body = "base64://" + tfConfig.HTTPConfig.ActionBody.ValueString()
	httpConfig.HttpRequestConfig.Method = tfConfig.HTTPConfig.RequestMethod.ValueString()
	httpConfig.HttpRequestConfig.Url = tfConfig.HTTPConfig.Url.ValueString()
	httpConfig.HttpRequestConfig.Headers = headersMap
}

func SmtpConfigToApi(tfConfig emailConfigurationResourceModel, smtpConfig *orytypes.SMTP) {
	headersMap := make(map[string]string)

	for _, header := range *tfConfig.SMTPHeaders {
		headersMap[header.Key.ValueString()] = header.Value.ValueString()
	}

	smtpConfig.ConnectionUri = buildSMTPURL(tfConfig.SMTPConfig.Username.ValueString(), tfConfig.SMTPConfig.Password.ValueString(),
		tfConfig.SMTPConfig.Host.ValueString(), tfConfig.SMTPConfig.Port.ValueString(), tfConfig.SMTPConfig.Security.ValueString())
	smtpConfig.FromAddress = tfConfig.SMTPConfig.SenderAddress.ValueString()
	smtpConfig.FromName = tfConfig.SMTPConfig.SenderName.ValueString()
	smtpConfig.Headers = headersMap
}

func buildSMTPURL(username, password, host string, port string, securityType string) string {
	if securityType == "" {
		securityType = "starttls"
	}

	escapedUsername := url.QueryEscape(username)
	escapedPassword := url.QueryEscape(password)

	scheme := "smtp"
	query := ""

	switch securityType {
	case "starttls":
		scheme = "smtp"
		query = ""
	case "starttls_notrust":
		scheme = "smtp"
		query = "skip_ssl_verify=true"
	case "cleartext":
		scheme = "smtp"
		query = "disable_starttls=true"
	case "implicittls":
		scheme = "smtps"
		query = ""
	case "implicittls_notrust":
		scheme = "smtps"
		query = "skip_ssl_verify=true"
	}

	smtpURI := fmt.Sprintf("%s://%s:%s@%s:%s", scheme, escapedUsername, escapedPassword, host, port)

	if query != "" {
		smtpURI = fmt.Sprintf("%s?%s", smtpURI, query)
	}

	return smtpURI
}
