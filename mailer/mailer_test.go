package mailer

import "testing"

var e = "test@test.com.au"
var f = "Rod"
var l = "Lewis"
var rc = "F0001234-0451-4000-B000-000000000000"

// @todo add conditions
func TestSendEmail(t *testing.T) {
	r := NewRequest([]string{f}, "Moonrock password reset")
	r.Send("templates/reset_template.html", map[string]string{
		"reset":    rc,
		"username": e,
	})
}
