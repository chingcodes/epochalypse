package main

import (
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Clock struct {
	Form *tview.Form

	unixField,
	timeField,
	utcField *tview.InputField

	updating bool
}

func NewClock() *Clock {
	c := &Clock{
		timeField: tview.NewInputField(),
		unixField: tview.NewInputField(),
		utcField:  tview.NewInputField(),
	}

	c.timeField.SetLabel("Local").
		SetFieldWidth(35).
		SetDoneFunc(func(tcell.Key) { c.timeUpdated() }).
		SetAcceptanceFunc(tview.InputFieldMaxLength(35))

	c.utcField.SetLabel("UTC").
		SetFieldWidth(35).
		SetDoneFunc(func(tcell.Key) { c.utcUpdated() }).
		SetAcceptanceFunc(tview.InputFieldMaxLength(35))

	c.unixField.SetLabel("Unix").
		SetFieldWidth(10).
		SetDoneFunc(func(tcell.Key) { c.unixUpdated() }).
		SetAcceptanceFunc(func(s string, r rune) bool {
			return len(s) <= 10 && tview.InputFieldInteger(s, r)
		})

	c.Form = tview.NewForm().
		AddFormItem(c.unixField).
		AddFormItem(c.timeField).
		AddFormItem(c.utcField).
		AddButton("Now", func() { c.SetTime(time.Now()) })

	c.SetTime(time.Now())
	return c
}

// Simple way to keep from Changes don't keep triggering other changes
func (c *Clock) update(f func()) {
	if c.updating == true {
		return
	}
	c.updating = true
	f()
	c.updating = false
}

func (c *Clock) unixUpdated() {
	t, err := strconv.ParseInt(c.unixField.GetText(), 10, 64)
	if err != nil {
		return
	}
	c.SetTime(time.Unix(t, 0))
}

func (c *Clock) timeUpdated() {
	t, err := time.Parse(time.UnixDate, c.timeField.GetText())
	if err != nil {
		return
	}
	c.SetTime(t)
}

func (c *Clock) utcUpdated() {
	t, err := time.Parse(time.UnixDate, c.utcField.GetText())
	if err != nil {
		return
	}
	c.SetTime(t.UTC())
}

func (c *Clock) SetTime(t time.Time) {
	c.update(func() {
		c.timeField.SetText(t.Local().Format(time.UnixDate))
		c.utcField.SetText(t.UTC().Format(time.UnixDate))
		c.unixField.SetText(strconv.FormatInt(t.Unix(), 10))
	})
}

func main() {
	app := tview.NewApplication()

	c := NewClock()

	c.Form.SetBorder(true).SetTitle(" epochalypse ").SetTitleAlign(tview.AlignLeft)

	if err := app.SetRoot(c.Form, true).SetFocus(c.Form).Run(); err != nil {
		panic(err)
	}
}
