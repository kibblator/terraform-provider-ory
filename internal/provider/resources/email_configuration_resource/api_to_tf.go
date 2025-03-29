package email_configuration_resource

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kibblator/terraform-provider-ory/internal/provider/helpers"
	orytypes "github.com/kibblator/terraform-provider-ory/internal/provider/types"
)

func ApiToHttpConfig(httpConfig *orytypes.HTTP, tfConfig *emailConfigurationResourceModel) error {
	httpAuthType := "none"

	if httpConfig.HttpRequestConfig.HttpAuth != nil {
		httpAuthType = httpConfig.HttpRequestConfig.HttpAuth.Type
	}

	username, password, in, name, value := parseHTTPAuthParams(httpAuthType, httpConfig.HttpRequestConfig.HttpAuth.HttpAuthConfig)

	tfConfig.HTTPConfig = &HTTPConfig{
		Url:                helpers.StringOrNil(httpConfig.HttpRequestConfig.Url),
		RequestMethod:      helpers.StringOrNil(httpConfig.HttpRequestConfig.Method),
		AuthenticationType: helpers.StringOrNil(httpAuthType),
	}

	if httpConfig.HttpRequestConfig.Body != "" {
		body, err := getBodyContents(httpConfig.HttpRequestConfig.Body)

		if err != nil {
			return err
		}

		tfConfig.HTTPConfig.ActionBody = helpers.StringOrNil(body)
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

	return nil
}

func ApiToSmtpConfig(smtpConfig *orytypes.SMTP, tfConfig *emailConfigurationResourceModel) error {
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

	if parsedURL.User != nil {
		username = parsedURL.User.Username()
		password, _ = parsedURL.User.Password() // Password might be empty
	}

	host = parsedURL.Hostname()
	port = parsedURL.Port()

	if host == "" {
		security = ""
	} else {
		security = "starttls"
	}

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

func parseHTTPAuthParams(httpAuthType string, httpAuthConfig *orytypes.HttpAuthConfig) (username, password, in, name, value string) {
	if httpAuthType == "api_key" {
		return "", "", httpAuthConfig.In, httpAuthConfig.Name, httpAuthConfig.Value
	}

	if httpAuthType == "basic_auth" {
		return httpAuthConfig.User, httpAuthConfig.Password, "", "", ""
	}

	return "", "", "", "", ""
}

func getBodyContents(input string) (string, error) {
	parsedURL, err := url.ParseRequestURI(input)
	if err != nil || !strings.HasPrefix(parsedURL.Scheme, "http") {
		return input, nil
	}

	resp, err := http.Get(input)
	if err != nil {
		return "", fmt.Errorf("error fetching file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error: received non-200 response code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return base64.RawStdEncoding.WithPadding(base64.StdPadding).EncodeToString(data), nil
}
