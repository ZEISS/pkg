package smtp

// ReplyCode is the status code for the SMTP server.
type ReplyCode int

const (
	// ReplyCodeSystemStatus is the status code for the system status.
	ReplyCodeSystemStatus = 211
	// ReplyCodeHelpMessage is the status code for the help message.
	ReplyCodeHelpMessage = 214
	// ReplyCodeServiceReady is the status code for the service ready.
	ReplyCodeServiceReady = 220
	// ReplyCodeServiceClosing is the status code for the service closing.
	ReplyCodeServiceClosing = 221
	// ReplyCodeAuthenticationSuccessful is the status code for the authentication successful.
	ReplyCodeMailActionOkay = 250
	// ReplyCodeUserNotLocal is the status code for the user not local.
	ReplyCodeMailActionCompleted = 250
	// ReplyCodeCannotVerifyUser is the status code for the cannot verify user.
	ReplyCodeUserNotLocal = 251
	// ReplyCodeStartMailInput is the status code for the start mail input.
	ReplyCodeCannotVerifyUser               = 252
	ReplyCodeStartMailInput                 = 354
	ReplyCodeServiceNotAvailable            = 421
	ReplyCodeMailboxUnavailable             = 450
	ReplyCodeLocalError                     = 451
	ReplyCodeInsufficientStorage            = 452
	ReplyCodeSyntaxError                    = 500
	ReplyCodeSyntaxErrorInParameters        = 501
	ReplyCodeCommandNotImplemented          = 502
	ReplyCodeCommandBadSequence             = 503
	ReplyCodeCommandParameterNotImplemented = 504
	ReplyCodeRequestActionNotTaken          = 550
	ReplyCodeUserNotLocalForThisHost        = 551
	ReplyCodeRequestedActionAborted         = 552
	ReplyCodeRequestedActionNotTaken        = 553
	ReplyCodeTransactionFailed              = 554
	ReplyCodeMailFromOrRcptToError          = 555
)
