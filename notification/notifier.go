package notification

type Notifier interface {
	Prepare (prepareFunc func() (Messager)) Driver
}

type Driver interface {
	Fire() bool
	SetMessager(message Messager)
	Messager
}

type Messager interface {
	SetFrom(from string)
	SetTo(to []string)
	SetBcc(bcc []string)
	SetCc(cc []string)
	SetSubject(subject string)
	SetBody(body string)

	From() string
	To() []string
	Bcc() []string
	Cc() []string
	Subject() string
	Body() string
}