package domain

type Sms struct {
	PhoneNumber    string `json:"phone_number,omitempty"`
	SignName       string `json:"sign_name,omitempty"`
	TemplateCode   string `json:"template_code,omitempty"`
	TemplateParams any    `json:"template_params,omitempty"`
}
