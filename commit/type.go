package commit

type Type int8

const (
	TypeFeat     Type = iota
	TypeFix
	TypeDocs
	TypeStyle
	TypeRefactor
	TypeTest
	TypeChore
)

func (t Type) Describe() string {
	switch t {
	case TypeFeat:
		return "feature"
	case TypeFix:
		return "bug fix"
	case TypeDocs:
		return "documentation"
	case TypeStyle:
		return "formatting, missing semi colons, …"
	case TypeRefactor:
		return "refactor"
	case TypeTest:
		return "when adding missing tests"
	case TypeChore:
		return "maintain"
	default:
		return "unknown"
	}
}

func (t Type) Name() string {
	switch t {
	case TypeFeat:
		return "feat"
	case TypeFix:
		return "fix"
	case TypeDocs:
		return "docs"
	case TypeStyle:
		return "style"
	case TypeRefactor:
		return "refactor"
	case TypeTest:
		return "test"
	case TypeChore:
		return "chore"
	default:
		return "unknown"
	}
}
