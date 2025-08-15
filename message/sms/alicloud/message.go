package alicloud

type Message struct {
	TemplateParam string   `json:"templateParam"`
	TemplateCode  string   `json:"templateCode"`
	PhoneNumbers  []string `json:"phoneNumbers"`
}
