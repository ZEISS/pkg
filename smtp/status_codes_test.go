package smtp_test

import (
	"testing"

	"github.com/katallaxie/pkg/smtp"
	"github.com/stretchr/testify/assert"
)

func TestNewStatusCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		message      string
		replyCode    smtp.ReplyCode
		enhancedCode smtp.EnhancedMailSystemStatusCode
	}{
		{
			name:         "ok",
			message:      "ok",
			replyCode:    smtp.ReplyCodeServiceReady,
			enhancedCode: smtp.EnhancedStatusCodeUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := smtp.NewStatusCode(tt.replyCode, tt.enhancedCode, tt.message)
			assert.Equal(t, tt.message, s.Message())
			assert.Equal(t, int(tt.replyCode), s.ReplyCode())
			assert.Equal(t, tt.enhancedCode, s.EnhancedStatusCode())
		})
	}
}
