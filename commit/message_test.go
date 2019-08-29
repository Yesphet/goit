package commit

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestMessage_Format(t *testing.T) {
	msg := &Message{
		Type:    TypeFeat,
		Scope:   []string{"plugins"},
		Subject: "Add new plugin xxx",
		Body:    "Add new plugin xxx, it ...........",
		Footer:  "Close #234",
	}

	expect := `feat(plugins): Add new plugin xxx

Add new plugin xxx, it ...........

Close #234`
	assert.EqualValues(t, expect, msg.Format())
}

func TestMessage_Format_WithoutBody(t *testing.T) {
	msg := &Message{
		Type:    TypeFeat,
		Scope:   []string{"plugins"},
		Subject: "Add new plugin xxx",
		Footer:  "Close #234",
	}

	expect := `feat(plugins): Add new plugin xxx

Add new plugin xxx, it ...........

Close #234`
	fmt.Println(msg.Format())
	assert.EqualValues(t, expect, msg.Format())
}