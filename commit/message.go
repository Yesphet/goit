package commit

import (
	"strings"
	"fmt"
)

type Message struct {
	Type    Type
	Scope   []string
	Subject string
	Body    string
	Footer  string
}

func (msg *Message) Format() string {
	scope := ""
	if msg.Scope != nil && len(msg.Scope) != 0 && strings.Join(msg.Scope, "") != "" {
		scope = "(" + strings.Join(msg.Scope, ",") + ")"
	}

	header := fmt.Sprintf("%s%s: %s\n", msg.Type.Name(), scope, msg.Subject)

	body := ""
	if msg.Body != "" {
		body = "\n" + msg.Body + "\n"
	}

	footer := ""
	if msg.Footer != "" {
		footer = "\n" + msg.Footer
	}

	return header + body + footer
}
