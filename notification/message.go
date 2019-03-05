package notification


type Message struct {
	from string
	to []string
	cc []string
	bcc []string
	subject string
	body string
}

func (m *Message) SetFrom(from string) {
	m.from = from
}
func (m *Message) SetTo(to []string) {
	m.to = to
}
func (m *Message) SetCc(cc []string)  {
	m.cc = cc
}
func (m *Message) SetBcc(bcc []string) {
	m.bcc = bcc
}
func (m *Message) SetSubject(subject string) {
	m.subject = subject
}
func (m *Message) SetBody(body string)  {
	m.body = body
}

func (m *Message) From() string {
	return m.from
}
func (m *Message) To() []string {
	return m.to
}
func (m *Message) Cc() []string {
	return m.cc
}
func (m *Message) Bcc() []string {
	return m.bcc
}
func (m *Message) Subject() string {
	return m.subject
}
func (m *Message) Body() string {
	return m.body
}