package main

import (
	"math"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// mobileEntry is a normal Entry with one extra thing:
// it knows which scroll container it lives in,
// so it can scroll itself into view when tapped (needed on Android).
type mobileEntry struct {
	widget.Entry
	scroll *container.Scroll
}

// newMobileEntry creates a mobileEntry with a placeholder and a scroll reference.
func newMobileEntry(placeholder string, scroll *container.Scroll) *mobileEntry {
	e := &mobileEntry{scroll: scroll}
	e.ExtendBaseWidget(e) // tells Fyne: use THIS type's methods, not the inner Entry's
	e.PlaceHolder = placeholder
	return e
}

// FocusGained is called automatically by Fyne when the user taps this field.
// We override it to also scroll the field into view after the keyboard appears.
func (e *mobileEntry) FocusGained() {
	e.Entry.FocusGained() // do the normal focus behaviour first
	if e.scroll != nil {
		go func() {
			time.Sleep(250 * time.Millisecond) // wait for keyboard animation to finish
			fyne.Do(func() {                   // fyne.Do runs this on the main UI thread safely
				e.scroll.ScrollToOffset(fyne.NewPos(0, e.Position().Y-40))
			})
		}()
	}
}

func main() {
	a := app.New()
	w := a.NewWindow("Financial Calculator")

	// --- Scroll containers created first so entries can hold a reference to them ---
	scroll1 := container.NewVScroll(nil)
	scroll2 := container.NewVScroll(nil)

	// --- FD Section ---
	l1 := widget.NewLabel("FD")
	p := newMobileEntry("Initial Amount", scroll1)
	ari := newMobileEntry("Annual rate of interest", scroll1)
	tm := newMobileEntry("Period of deposit in Years", scroll1)
	r := widget.NewLabel("Maturity value: ...")
	interest := widget.NewLabel("Interest: ...")
	var mv float64
	fd := func(string) {
		f1, e1 := strconv.ParseFloat(p.Text, 64)
		f2, e2 := strconv.ParseFloat(ari.Text, 64)
		f3, e3 := strconv.ParseFloat(tm.Text, 64)
		if e1 == nil && e2 == nil && e3 == nil {
			mv = f1 * math.Pow((1+(f2/400)), (f3*4))
			r.SetText("Maturity value: " + strconv.FormatFloat(mv, 'f', 2, 64))
			interest.SetText("Interest: " + strconv.FormatFloat((mv-f1), 'f', 2, 64))
		} else {
			r.SetText("Maturity value: ...")
			interest.SetText("Interest: ...")
		}
	}
	p.OnChanged = fd
	ari.OnChanged = fd
	tm.OnChanged = fd

	// --- Reverse FD 1 ---
	l2 := widget.NewLabel("Reverse FD - Compute Initial Amount")
	ari2 := newMobileEntry("Annual rate of interest", scroll1)
	tm2 := newMobileEntry("Period of deposit in Years", scroll1)
	mv2 := newMobileEntry("Maturity Value", scroll1)
	p2 := widget.NewLabel("Initial amount: ...")
	i2 := widget.NewLabel("Interest: ...")
	var pri float64
	fd2 := func(string) {
		f4, e4 := strconv.ParseFloat(ari2.Text, 64)
		f5, e5 := strconv.ParseFloat(tm2.Text, 64)
		f6, e6 := strconv.ParseFloat(mv2.Text, 64)
		if e4 == nil && e5 == nil && e6 == nil {
			pri = f6 / math.Pow((1+(f4/400)), (f5*4))
			p2.SetText("Initial amount: " + strconv.FormatFloat(pri, 'f', 2, 64))
			i2.SetText("Interest: " + strconv.FormatFloat((f6-pri), 'f', 2, 64))
		} else {
			p2.SetText("Initial amount: ...")
			i2.SetText("Interest: ...")
		}
	}
	ari2.OnChanged = fd2
	tm2.OnChanged = fd2
	mv2.OnChanged = fd2

	// --- Reverse FD 2 ---
	l3 := widget.NewLabel("Reverse FD - Compute period of investment")
	ari3 := newMobileEntry("Annual rate of interest", scroll1)
	p3 := newMobileEntry("Initial amount", scroll1)
	mv3 := newMobileEntry("Maturity Value", scroll1)
	tm3 := widget.NewLabel("Period of investment: ...")
	i3 := widget.NewLabel("Interest: ...")
	var t3 float64
	fd3 := func(string) {
		f11, e7 := strconv.ParseFloat(ari3.Text, 64)
		f12, e8 := strconv.ParseFloat(p3.Text, 64)
		f13, e9 := strconv.ParseFloat(mv3.Text, 64)
		if e7 == nil && e8 == nil && e9 == nil {
			t3 = (math.Log(f13/f12) / math.Log(1+f11/400)) / 4
			tm3.SetText("Period of investment: " + strconv.FormatFloat(t3, 'f', 2, 64))
			i3.SetText("Interest: " + strconv.FormatFloat((f13-f12), 'f', 2, 64))
		} else {
			tm3.SetText("Period of investment: ...")
			i3.SetText("Interest: ...")
		}
	}
	ari3.OnChanged = fd3
	p3.OnChanged = fd3
	mv3.OnChanged = fd3

	// --- RD Section ---
	winl2 := widget.NewLabel("RD")
	mi := newMobileEntry("Monthly deposit", scroll2)
	ari4 := newMobileEntry("Annual rate of interest", scroll2)
	p4 := newMobileEntry("Period of investment in years", scroll2)
	mv4 := widget.NewLabel("Maturity value: ...")
	i4 := widget.NewLabel("Interest: ...")
	var m4 float64
	fd4 := func(string) {
		g1, h1 := strconv.ParseFloat(mi.Text, 64)
		g2, h2 := strconv.ParseFloat(ari4.Text, 64)
		g3, h3 := strconv.ParseFloat(p4.Text, 64)
		if h1 == nil && h2 == nil && h3 == nil {
			m4 = g1 * (math.Pow((1+(g2/400)), (4*g3)) - 1) / (1 - (1 / math.Pow((1+(g2/400)), (1.0/3.0))))
			mv4.SetText("Maturity value: " + strconv.FormatFloat(m4, 'f', 2, 64))
			i4.SetText("Interest: " + strconv.FormatFloat((m4-(g3*g1*12)), 'f', 2, 64))
		} else {
			mv4.SetText("Maturity value: ...")
			i4.SetText("Interest: ...")
		}
	}
	mi.OnChanged = fd4
	ari4.OnChanged = fd4
	p4.OnChanged = fd4

	// --- Reverse RD 1 ---
	winl3 := widget.NewLabel("Reverse RD - Compute monthly deposit")
	mv5 := newMobileEntry("Maturity Value", scroll2)
	ari5 := newMobileEntry("Annual rate of interest", scroll2)
	p5 := newMobileEntry("Period of investment in years", scroll2)
	mi2 := widget.NewLabel("Monthly investment: ...")
	i5 := widget.NewLabel("Interest: ...")
	var minv float64
	fd5 := func(string) {
		g4, h4 := strconv.ParseFloat(mv5.Text, 64)
		g5, h5 := strconv.ParseFloat(ari5.Text, 64)
		g6, h6 := strconv.ParseFloat(p5.Text, 64)
		if h4 == nil && h5 == nil && h6 == nil {
			minv = g4 * (1 - (1 / math.Pow((1+(g5/400)), (1.0/3.0)))) / (math.Pow((1+(g5/400)), (4*g6)) - 1)
			mi2.SetText("Monthly investment: " + strconv.FormatFloat(minv, 'f', 2, 64))
			i5.SetText("Interest: " + strconv.FormatFloat((g4-(g6*minv*12)), 'f', 2, 64))
		} else {
			mi2.SetText("Monthly investment: ...")
			i5.SetText("Interest: ...")
		}
	}
	mv5.OnChanged = fd5
	ari5.OnChanged = fd5
	p5.OnChanged = fd5

	// --- Reverse RD 2 ---
	winl4 := widget.NewLabel("Reverse RD - Compute period of investment")
	mv6 := newMobileEntry("Maturity Value", scroll2)
	ari6 := newMobileEntry("Annual rate of interest", scroll2)
	mi3 := newMobileEntry("Monthly investment", scroll2)
	tm4 := widget.NewLabel("Period of investment in years: ...")
	i6 := widget.NewLabel("Interest: ...")
	var t float64
	fd6 := func(string) {
		g7, h7 := strconv.ParseFloat(mv6.Text, 64)
		g8, h8 := strconv.ParseFloat(ari6.Text, 64)
		g9, h9 := strconv.ParseFloat(mi3.Text, 64)
		if h7 == nil && h8 == nil && h9 == nil {
			t = (math.Log(((g7/g9)*(1-(1/math.Pow((1+(g8/400)), (1.0/3.0)))))+1) / math.Log(1+(g8/400))) / 4
			tm4.SetText("Period of investment in years: " + strconv.FormatFloat(t, 'f', 2, 64))
			i6.SetText("Interest: " + strconv.FormatFloat((g7-(g9*t*12)), 'f', 2, 64))
		} else {
			tm4.SetText("Period of investment in years: ...")
			i6.SetText("Interest: ...")
		}
	}
	mv6.OnChanged = fd6
	ari6.OnChanged = fd6
	mi3.OnChanged = fd6

	// --- Build the two pages and assign directly into scroll containers ---
	scroll1.Content = container.NewVBox(
		widget.NewButton("Switch to RD", func() {
			w.SetContent(scroll2)
			w.Resize(fyne.NewSize(400, 700))
		}),
		container.New(layout.NewCenterLayout(), l1), p, ari, tm, r, interest,
		container.New(layout.NewCenterLayout(), l2), mv2, ari2, tm2, p2, i2,
		container.New(layout.NewCenterLayout(), l3), mv3, p3, ari3, tm3, i3,
	)

	scroll2.Content = container.NewVBox(
		widget.NewButton("Switch to FD", func() {
			w.SetContent(scroll1)
			w.Resize(fyne.NewSize(400, 700))
		}),
		container.New(layout.NewCenterLayout(), winl2), mi, ari4, p4, mv4, i4,
		container.New(layout.NewCenterLayout(), winl3), mv5, ari5, p5, mi2, i5,
		container.New(layout.NewCenterLayout(), winl4), mv6, ari6, mi3, tm4, i6,
	)

	w.SetContent(scroll1)
	w.Resize(fyne.NewSize(400, 700))
	w.ShowAndRun()
}
