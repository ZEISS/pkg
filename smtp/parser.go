package smtp

// Verb is a command verb of SMTP.
type Verb string

const (
	// HELO is the SMTP verb for HELO.
	HELO Verb = "HELO"
	// EHLO is the SMTP verb for EHLO.
	EHLO Verb = "EHLO"
	// MAIL is the SMTP verb for MAIL.
	MAIL Verb = "MAIL"
	// RCPT is the SMTP verb for RCPT.
	RCPT Verb = "RCPT"
	// DATA is the SMTP verb for DATA.
	DATA Verb = "DATA"
	// QUI is the SMTP verb for QUIT.
	QUIT Verb = "QUIT"
	// RSET is the SMTP verb for RSET.
	RSET Verb = "RSET"
	// NOOP is the SMTP verb for NOOP.
	NOOP Verb = "NOOP"
	// HELP is the SMTP verb for HELP.
	HELP Verb = "HELP"
	// STARTTLS is the SMTP verb for STARTTLS.
	STARTTLS Verb = "STARTTLS"
	// AUTH is the SMTP verb for AUTH.
	AUTH Verb = "AUTH"
)
