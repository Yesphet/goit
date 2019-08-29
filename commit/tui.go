package commit

import (
	"github.com/rivo/tview"
	"github.com/gdamore/tcell"
	"fmt"
	"strings"
)

const retPrefix = "  "

func CreateMessageFromTUI() (*Message, error) {
	app := TUIApplication{
		msg: &Message{},
		style: TUIStyleConfig{
			TipTextColor:          tcell.ColorWhite,
			HintTextColor:         tcell.ColorWhite,
			SelectedHintTextColor: tcell.ColorGreen,
			RetTextColor:          tcell.ColorBlue,
		},
	}
	app.Start()
	return nil, nil
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
			retTexView.SetText(fmt.Sprintf("  %-12s %s", t.Name()+":", t.Describe()))
			flex.RemoveItem(typeList).
				AddItem(retTexView, 0, 1, false)
			app.msg.Type = t
			app.foldCurrentAndDrawNextFlex(flex, nextFlex)
		}
	}

	for _, t := range Types {
		typeList.AddItem(fmt.Sprintf("  %-12s %s", t.Name()+":", t.Describe()), "", 0, doneFunc(t))
	}
	return flex
}

func (app *TUIApplication) denoteScope(nextFlex *tview.Flex) *tview.Flex {
	hintText := "Denote the scope of this change, use commas to separate multiple scopes:"
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
		app.msg.Scope = strings.Split(input.GetText(), ",")
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
	return app.commonWriteFlex(hintText, nextFlex, func(inputText string) {
		app.msg.Subject = inputText
	})
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
		fmt.Println(app.msg.Format())
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
