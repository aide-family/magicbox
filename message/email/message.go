package email

type Attachment struct {
	Filename string
	Data     []byte
}

type Message struct {
	To          []string
	Cc          []string
	Subject     string
	Body        string
	ContentType string
	Attachments []*Attachment
	Headers     map[string][]string
}
