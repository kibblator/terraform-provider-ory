package email_configuration_resource

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kibblator/terraform-provider-ory/internal/provider/helpers"
	orytypes "github.com/kibblator/terraform-provider-ory/internal/provider/types"
)

func ApiToHttpConfig(httpConfig *orytypes.HTTP, tfConfig emailConfigurationResourceModel) {
	httpAuthType := "none"

	if httpConfig.HttpRequestConfig.HttpAuth != nil {
		httpAuthType = httpConfig.HttpRequestConfig.HttpAuth.Type
	}

	username, password, in, name, value := parseHTTPAuthParams(httpAuthType, httpConfig.HttpRequestConfig.HttpAuth.HttpAuthConfig)

	tfConfig.HTTPConfig = &HTTPConfig{
		Url:                helpers.StringOrNil(httpConfig.HttpRequestConfig.Url),
		RequestMethod:      helpers.StringOrNil(httpConfig.HttpRequestConfig.Method),
		AuthenticationType: helpers.StringOrNil(httpAuthType),
		ActionBody:         helpers.StringOrNil(httpConfig.HttpRequestConfig.Body),
	}

	if username != "" && password != "" {
		tfConfig.HTTPConfig.BasicAuth = &BasicAuth{
			Username: helpers.StringOrNil(username),
			Password: helpers.StringOrNil(password),
		}
	}

	if in != "" && name != "" && value != "" {
		tfConfig.HTTPConfig.ApiKey = &APIKey{
			TransportMode: helpers.StringOrNil(in),
			Name:          helpers.StringOrNil(name),
			Value:         helpers.StringOrNil(value),
		}
	}

	if httpConfig.HttpRequestConfig.Headers != nil {
		headers := []SMTPHeader{}

		for key, value := range httpConfig.HttpRequestConfig.Headers {
			headers = append(headers, SMTPHeader{
				Key:   types.StringValue(key),
				Value: types.StringValue(value),
			})
		}

		tfConfig.SMTPHeaders = &headers
	}
}

func ApiToSmtpConfig(smtpConfig *orytypes.SMTP, tfConfig emailConfigurationResourceModel) error {
	argUsername, argPassword, argHost, argPort, argSecurity, argErr := parseSMTPURL(smtpConfig.ConnectionUri)

	if argErr != nil {
		return argErr
	}

	tfConfig.SMTPConfig = &SMTPConfig{
		SenderName:    helpers.StringOrNil(smtpConfig.FromName),
		SenderAddress: helpers.StringOrNil(smtpConfig.FromAddress),
		Host:          helpers.StringOrNil(argHost),
		Port:          helpers.StringOrNil(argPort),
		Security:      helpers.StringOrNil(argSecurity),
		Username:      helpers.StringOrNil(argUsername),
		Password:      helpers.StringOrNil(argPassword),
	}

	if smtpConfig.Headers != nil {
		headers := []SMTPHeader{}

		for key, value := range smtpConfig.Headers {
			headers = append(headers, SMTPHeader{
				Key:   types.StringValue(key),
				Value: types.StringValue(value),
			})
		}

		tfConfig.SMTPHeaders = &headers
	}

	return nil
}

func parseSMTPURL(smtpURL string) (username, password, host, port, security string, err error) {
	parsedURL, err := url.Parse(smtpURL)
	if err != nil {
		return "", "", "", "", "", err
	}

	// Extract username and password
	if parsedURL.User != nil {
		username = parsedURL.User.Username()
		password, _ = parsedURL.User.Password() // Password might be empty
	}

	// Extract host and port
	host = parsedURL.Hostname()
	port = parsedURL.Port()

	// Default security type
	if host == "" {
		security = ""
	} else {
		security = "starttls"
	}

	// Check for security parameter
	queryParams := parsedURL.Query()
	param := ""

	for key, values := range queryParams {
		param = fmt.Sprintf("%s=%s", key, values[0])
		break
	}

	switch param {
	case "skip_ssl_verify=true":
		if parsedURL.Scheme == "smtps" {
			security = "implicittls_notrust"
		} else {
			security = "starttls_notrust"
		}
	case "disable_starttls=true":
		security = "cleartext"
	case "":
		if parsedURL.Scheme == "smtps" {
			security = "implicittls"
		} else {
			security = "starttls"
		}
	}

	return username, password, host, port, security, nil
}

func parseHTTPAuthParams(hhttpAuthType string, httpAuthConfig *orytypes.HttpAuthConfig) (username, password, in, name, value string) {
	if hhttpAuthType == "api_key" {
		return "", "", httpAuthConfig.In, httpAuthConfig.Name, httpAuthConfig.Value
	}

	if hhttpAuthType == "basic_auth" {
		return httpAuthConfig.User, httpAuthConfig.Password, "", "", ""
	}

	return "", "", "", "", ""
}
