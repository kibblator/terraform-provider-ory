package orytypes

import "encoding/json"

type Config struct {
	Clients *Clients `json:"clients,omitempty"`
	Courier *Courier `json:"courier,omitempty"`
}

type Clients struct {
	WebHook *WebHook `json:"web_hook,omitempty"`
}

type WebHook struct {
	HeaderAllowlist []string `json:"header_allowlist,omitempty"`
}

type Courier struct {
	SMTP             *SMTP      `json:"smtp,omitempty"`
	HTTP             *HTTP      `json:"http,omitempty"`
	DeliveryStrategy *string    `json:"delivery_strategy,omitempty"`
	Templates        *Templates `json:"templates,omitempty"`
}

type SMTP struct {
	ConnectionUri string `json:"connection_uri,omitempty"`
	FromAddress   string `json:"from_address,omitempty"`
	FromName      string `json:"from_name,omitempty"`
	Headers       string `json:"headers,omitempty"`
}

type HTTP struct {
	FromName string `json:"from_name,omitempty"`
}

type Templates struct {
	LoginCode *TemplateBody `json:"login_code,omitempty"`
	Recovery  *TemplateBody `json:"recovery,omitempty"`
}

type TemplateBody struct {
	Valid   *TemplateDetail `json:"valid,omitempty"`
	Invalid *TemplateDetail `json:"invalid,omitempty"`
}

type TemplateDetail struct {
	Email *EmailBody `json:"email,omitempty"`
	SMS   *SMSBody   `json:"sms,omitempty"`
}

type EmailBody struct {
	Body interface{} `json:"body,omitempty"`
}

type SMSBody struct {
	Body interface{} `json:"body,omitempty"`
}

func TransformToConfig(data map[string]interface{}, config *Config) error {
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(raw, &config)
	if err != nil {
		return err
	}

	return nil
}
