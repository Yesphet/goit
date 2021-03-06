package commit

import (
	"fmt"
	"os"
	"strings"

	"github.com/Yesphet/goit/config"
)

type Type struct {
	Name        string
	Description string
}

var Types []Type
var TypeUnknown = Type{Name: "unknown", Description: "Unknown"}

func init() {
	if !config.Global.Commit.DisableDefaultTypes {
		addDefaultTypes()
	}

	for _, s := range config.Global.Commit.Types {
		AddCustomType(s)
	}

	if len(Types) == 0 {
		fmt.Println("There has no types defined.")
		os.Exit(1)
	}
}

func addDefaultTypes() {
	AddCustomType("feat: A new feature")
	AddCustomType("fix: A bug fix")
	AddCustomType("modify: A code change that neither fixes a bug nor adds a feature")
	AddCustomType("test: Adding missing tests or correcting existing tests")
	AddCustomType("docs: Documentation only changes")
	AddCustomType("style: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)")
	AddCustomType("build: Changes that affect the build system or external dependencies")
	AddCustomType("chore: Tool changes, configuration changes, version releases, etc")
}

func NewType(s string) Type {
	part := strings.Split(s, ":")
	name := strings.TrimSpace(part[0])
	desc := ""
	if len(part) > 1 {
		desc = strings.TrimSpace(strings.Join(part[1:], ":"))
	}
	return Type{
		Name:        name,
		Description: desc,
	}
}

func (t Type) String() string {
	return fmt.Sprintf("%-12s %s", t.Name+":", t.Description)
}

func AddCustomType(s string) {
	Types = append(Types, NewType(s))
}

func TypeOf(name string) Type {
	for i := range Types {
		if Types[i].Name == name {
			return Types[i]
		}
	}
	return TypeUnknown
}
