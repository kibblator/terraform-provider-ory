package orytypes

import "encoding/json"

type Config struct {
	Clients     *Clients     `json:"clients,omitempty"`
	Courier     *Courier     `json:"courier,omitempty"`
	SelfService *SelfService `json:"selfservice,omitempty"`
}

type SelfService struct {
	Flows   *Flows   `json:"flows,omitempty"`
	Methods *Methods `json:"methods,omitempty"`
}

type Methods struct {
	Password PasswordMethod `json:"password,omitempty"`
}

type PasswordMethod struct {
	Config  PasswordMethodConfig `json:"config,omitempty"`
	Enabled bool                 `json:"enabled"`
}

type PasswordMethodConfig struct {
	HaveIBeenPwnedEnabled            bool `json:"haveibeenpwned_enabled,omitempty"`
	IdentifierSimilarityCheckEnabled bool `json:"identifier_similarity_check_enabled,omitempty"`
	IgnoreNetworkErrors              bool `json:"ignore_network_errors,omitempty"`
	MaxBreaches                      bool `json:"max_breaches,omitempty"`
	MinPasswordLength                bool `json:"min_password_length,omitempty"`
}

type Flows struct {
	Registration *Registration `json:"registration,omitempty"`
}

type Registration struct {
	After               After  `json:"after"`
	Before              Before `json:"before"`
	EnableLegacyOneStep bool   `json:"enable_legacy_one_step"`
	Enabled             bool   `json:"enabled"`
	Lifespan            string `json:"lifespan"`
	LoginHints          bool   `json:"login_hints"`
	UIURL               string `json:"ui_url"`
}

type After struct {
	Code     AuthMethod `json:"code"`
	Hooks    []Hook     `json:"hooks"`
	OIDC     AuthMethod `json:"oidc"`
	Passkey  AuthMethod `json:"passkey"`
	Password AuthMethod `json:"password"`
	SAML     AuthMethod `json:"saml"`
	WebAuthn AuthMethod `json:"webauthn"`
}

type Before struct {
	Hooks []Hook `json:"hooks"`
}

type Hook struct {
	Hook   string                 `json:"hook"`
	Config map[string]interface{} `json:"config,omitempty"`
}

type AuthMethod struct {
	Hooks []Hook `json:"hooks"`
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
	ConnectionUri string            `json:"connection_uri,omitempty"`
	FromAddress   string            `json:"from_address,omitempty"`
	FromName      string            `json:"from_name,omitempty"`
	Headers       map[string]string `json:"headers,omitempty"`
}

type HTTP struct {
	HttpRequestConfig *HttpRequestConfig `json:"request_config,omitempty"`
}

type HttpRequestConfig struct {
	HttpAuth *HttpAuth         `json:"auth,omitempty"`
	Body     string            `json:"body,omitempty"`
	Headers  map[string]string `json:"headers,omitempty"`
	Method   string            `json:"method,omitempty"`
	Url      string            `json:"url,omitempty"`
}

type HttpAuth struct {
	HttpAuthConfig *HttpAuthConfig `json:"config,omitempty"`
	Type           string          `json:"type,omitempty"`
}

type HttpAuthConfig struct {
	Password string `json:"password,omitempty"`
	User     string `json:"user,omitempty"`
	In       string `json:"in,omitempty"`
	Name     string `json:"name,omitempty"`
	Value    string `json:"value,omitempty"`
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
