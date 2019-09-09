package commit

import (
	"strings"

	"github.com/Yesphet/goit/config"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

const (
	retPrefix      = "  "
	scopeSeparator = ","
)

func Do() error {
	gh, err := newGitHelper()
	if err != nil {
		return err
	}

	if err := gh.gitCommitPreCheck(); err != nil {
		return err
	}
	msg, err := CreateMessageFromTUI()
	if err != nil {
		return err
	}

	return gh.gitCommit(msg.Format())
}

func CreateMessageFromTUI() (*Message, error) {
	app := TUIApplication{
		msg: &Message{},
		style: TUIStyleConfig{
			TipTextColor:          tcell.ColorWhite,
			HintTextColor:         tcell.ColorWhite,
			SelectedHintTextColor: tcell.ColorGreen,
			RetTextColor:          tcell.ColorBlue,
		},
		autoScopeList: config.Global.Commit.Scopes,
	}
	err := app.Start()
	if err != nil {
		return nil, err
	}
	return app.msg, nil
}

type TUIStyleConfig struct {
	TipTextColor          tcell.Color
	HintTextColor         tcell.Color
	SelectedHintTextColor tcell.Color
	RetTextColor          tcell.Color
}

type TUIApplication struct {
	tviewApp *tview.Application
	layout   *tview.Flex
	msg      *Message
	style    TUIStyleConfig

	autoScopeList []string
}

func (app *TUIApplication) Start() error {
	app.tviewApp = tview.NewApplication()
	app.layout = tview.NewFlex().
		SetDirection(tview.FlexRow)

	writeFooterFlex := app.writeFooter()
	writeBodyFlex := app.writeBody(writeFooterFlex)
	writeSubjectFlex := app.writeSubject(writeBodyFlex)
	denoteScopeFlex := app.denoteScope(writeSubjectFlex)
	selectTypeFlex := app.selectType(denoteScopeFlex)

	app.layout.AddItem(selectTypeFlex, 0, 1, true)
	if err := app.tviewApp.SetRoot(app.layout, true).Run(); err != nil {
		return err
	}
	return nil
}

func (app *TUIApplication) tipsView() *tview.TextView {
	textView := tview.NewTextView().
		SetText("Line 1 will be cropped at 100 characters. All other lines will be wrapped after 100 characters.").
		SetTextColor(app.style.TipTextColor)
	return textView
}

func (app *TUIApplication) foldCurrentAndDrawNextFlex(current, next *tview.Flex) {
	app.layout.
		RemoveItem(current).
		AddItem(current, 2, 0, false)
	if next != nil {
		app.layout.AddItem(next, 0, 5, false)
		app.tviewApp.SetFocus(next)
	}
}

func (app *TUIApplication) selectType(nextFlex *tview.Flex) *tview.Flex {
	hintText := "Select the type of change that you're committing: "
	hintTextView := tview.NewTextView().
		SetTextColor(app.style.HintTextColor).
		SetText("? " + hintText)

	typeList := tview.NewList().
		ShowSecondaryText(false)

	retTexView := tview.NewTextView().
		SetTextColor(app.style.RetTextColor)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(hintTextView, 1, 1, false).
		AddItem(typeList, 0, len(Types), true)

	doneFunc := func(t Type) func() {
		return func() {
			hintTextView.SetText("√ " + hintText)
			hintTextView.SetTextColor(app.style.SelectedHintTextColor)
			retTexView.SetText("  " + t.String())
			flex.RemoveItem(typeList).
				AddItem(retTexView, 0, 1, false)
			app.msg.Type = t
			app.foldCurrentAndDrawNextFlex(flex, nextFlex)
		}
	}

	for _, t := range Types {
		typeList.AddItem("  "+t.String(), "", 0, doneFunc(t))
	}
	return flex
}

func (app *TUIApplication) scopeAutocompleteFunc(currentText string) []string {
	scopes := strings.Split(currentText, scopeSeparator)
	semiWord := scopes[len(scopes)-1]
	entries := make([]string, 0)
loop:
	for _, defined := range app.autoScopeList {
		if strings.HasPrefix(defined, semiWord) {
			for i := 0; i < len(scopes)-1; i++ {
				if scopes[i] == defined {
					continue loop
				}
			}
			prefix := ""
			if len(scopes) > 1 {
				prefix = strings.Join(scopes[:len(scopes)-1], ",") + ","
			}
			entries = append(entries, prefix+defined)
		}
	}
	return entries
}

func (app *TUIApplication) denoteScope(nextFlex *tview.Flex) *tview.Flex {
	hintText := "Denote the scope of this change, use commas(,) to separate multiple scopes:"
	hintTextView := tview.NewTextView().
		SetTextColor(app.style.HintTextColor).
		SetText("? " + hintText)

	retTexView := tview.NewTextView().
		SetTextColor(app.style.RetTextColor)

	input := tview.NewInputField().
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorWhite).
		SetLabel(retPrefix).
		SetAutocompleteFunc(app.scopeAutocompleteFunc)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(hintTextView, 1, 1, false).
		AddItem(input, 0, 2, true)

	input.SetFinishedFunc(func(key tcell.Key) {
		hintTextView.SetText("√ " + hintText)
		hintTextView.SetTextColor(app.style.SelectedHintTextColor)
		retTexView.SetText(retPrefix + input.GetText())
		flex.RemoveItem(input).
			AddItem(retTexView, 0, 1, false)
		if input.GetText() != "" {
			app.msg.Scope = strings.Split(input.GetText(), scopeSeparator)
		}
		app.foldCurrentAndDrawNextFlex(flex, nextFlex)
	})
	return flex
}

func (app *TUIApplication) commonWriteFlex(hintText string, nextFlex *tview.Flex, finishInputFunc func(inputText string)) *tview.Flex {
	hintTextView := tview.NewTextView().
		SetTextColor(app.style.HintTextColor).
		SetText("? " + hintText)

	retTexView := tview.NewTextView().
		SetTextColor(app.style.RetTextColor)

	input := tview.NewInputField().
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorWhite).
		SetLabel(retPrefix)
	input.SetInputCapture(ignoreUpAndDownInputCapture)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(hintTextView, 1, 1, false).
		AddItem(input, 0, 2, true)

	input.SetFinishedFunc(func(key tcell.Key) {
		hintTextView.SetText("√ " + hintText)
		hintTextView.SetTextColor(app.style.SelectedHintTextColor)
		retTexView.SetText(retPrefix + input.GetText())
		flex.RemoveItem(input).
			AddItem(retTexView, 0, 1, false)
		finishInputFunc(input.GetText())
		app.foldCurrentAndDrawNextFlex(flex, nextFlex)
	})
	return flex
}

func (app *TUIApplication) writeSubject(nextFlex *tview.Flex) *tview.Flex {
	hintText := "Write a short, imperative and present tense description of the change:"
	hintTextView := tview.NewTextView().
		SetTextColor(app.style.HintTextColor).
		SetText("? " + hintText)

	retTexView := tview.NewTextView().
		SetTextColor(app.style.RetTextColor)

	input := tview.NewInputField().
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorWhite).
		SetLabel(retPrefix)
	emptyInputValidateInputCapture := func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter && input.GetText() == "" {
			hintTextView.SetText("? Short description can't be empty").SetTextColor(tcell.ColorRed)
			return nil
		}
		return event
	}
	input.SetInputCapture(multiInputCapture(ignoreUpAndDownInputCapture, emptyInputValidateInputCapture))

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(hintTextView, 1, 1, false).
		AddItem(input, 0, 2, true)

	input.SetFinishedFunc(func(key tcell.Key) {
		hintTextView.SetText("√ " + hintText)
		hintTextView.SetTextColor(app.style.SelectedHintTextColor)
		retTexView.SetText(retPrefix + input.GetText())
		flex.RemoveItem(input).
			AddItem(retTexView, 0, 1, false)
		app.msg.Subject = input.GetText()
		app.foldCurrentAndDrawNextFlex(flex, nextFlex)
	})
	return flex
}

func (app *TUIApplication) writeBody(nextFlex *tview.Flex) *tview.Flex {
	hintText := "Provide a longer description of the change:"
	return app.commonWriteFlex(hintText, nextFlex, func(inputText string) {
		app.msg.Body = inputText
	})
}

func (app *TUIApplication) writeFooter() *tview.Flex {
	hintText := "List any breaking changes or issues closed by this change:"
	return app.commonWriteFlex(hintText, nil, func(inputText string) {
		app.msg.Footer = inputText
		app.tviewApp.Stop()
	})
}

func ignoreUpAndDownInputCapture(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyUp:
		return nil
	case tcell.KeyDown:
		return nil
	default:
		return event
	}
}

func multiInputCapture(captures ...func(event *tcell.EventKey) *tcell.EventKey) func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		input := event
		for _, capture := range captures {
			input = capture(input)
		}
		return input
	}
}
