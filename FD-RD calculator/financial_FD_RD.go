package main

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type mobileEntry struct {
	widget.Entry
	scroll *container.Scroll
}

func newMobileEntry(placeholder string, scroll *container.Scroll) *mobileEntry {
	e := &mobileEntry{scroll: scroll}
	e.ExtendBaseWidget(e)
	e.PlaceHolder = placeholder
	return e
}

func (e *mobileEntry) FocusGained() {
	e.Entry.FocusGained()
	if e.scroll != nil {
		go func() {
			time.Sleep(250 * time.Millisecond)
			fyne.Do(func() {
				e.scroll.ScrollToOffset(fyne.NewPos(0, e.Position().Y-40))
			})
		}()
	}
}

// calcBlock encapsulates the creation of a calculation section
func calcBlock(
	title string,
	placeholders [3]string,
	scroll *container.Scroll,
	outLabel1Prefix, outLabel2Prefix string,
	calcFunc func(v1, v2, v3 float64) (float64, float64),
) fyne.CanvasObject {
	e1 := newMobileEntry(placeholders[0], scroll)
	e2 := newMobileEntry(placeholders[1], scroll)
	e3 := newMobileEntry(placeholders[2], scroll)

	out1 := widget.NewLabel(outLabel1Prefix + ": ...")
	out2 := widget.NewLabel(outLabel2Prefix + ": ...")

	recalc := func(string) {
		f1, err1 := strconv.ParseFloat(e1.Text, 64)
		f2, err2 := strconv.ParseFloat(e2.Text, 64)
		f3, err3 := strconv.ParseFloat(e3.Text, 64)

		if err1 == nil && err2 == nil && err3 == nil {
			r1, r2 := calcFunc(f1, f2, f3)
			out1.SetText(fmt.Sprintf("%s: %.2f", outLabel1Prefix, r1))
			out2.SetText(fmt.Sprintf("%s: %.2f", outLabel2Prefix, r2))
		} else {
			out1.SetText(outLabel1Prefix + ": ...")
			out2.SetText(outLabel2Prefix + ": ...")
		}
	}

	e1.OnChanged = recalc
	e2.OnChanged = recalc
	e3.OnChanged = recalc

	return container.NewVBox(
		container.New(layout.NewCenterLayout(), widget.NewLabel(title)),
		e1, e2, e3, out1, out2,
	)
}

func main() {
	a := app.New()
	w := a.NewWindow("Financial Calculator")

	scroll1 := container.NewVScroll(nil)
	scroll2 := container.NewVScroll(nil)

	// --- FD Section Blocks ---
	fdBlock := calcBlock("FD", [3]string{"Initial Amount", "Annual rate of interest", "Period of deposit in Years"}, scroll1, "Maturity value", "Interest",
		func(p, r, t float64) (float64, float64) {
			mv := p * math.Pow(1+(r/400), t*4)
			return mv, mv - p
		})

	revFd1 := calcBlock("Reverse FD - Compute Initial Amount", [3]string{"Maturity Value", "Annual rate of interest", "Period of deposit in Years"}, scroll1, "Initial amount", "Interest",
		func(mv, r, t float64) (float64, float64) {
			p := mv / math.Pow(1+(r/400), t*4)
			return p, mv - p
		})

	revFd2 := calcBlock("Reverse FD - Compute period of investment", [3]string{"Maturity Value", "Initial amount", "Annual rate of interest"}, scroll1, "Period of investment", "Interest",
		func(mv, p, r float64) (float64, float64) {
			t := (math.Log(mv/p) / math.Log(1+r/400)) / 4
			return t, mv - p
		})

	// --- RD Section Blocks ---
	rdBlock := calcBlock("RD", [3]string{"Monthly deposit", "Annual rate of interest", "Period of investment in years"}, scroll2, "Maturity value", "Interest",
		func(mi, r, t float64) (float64, float64) {
			mv := mi * (math.Pow(1+(r/400), 4*t) - 1) / (1 - (1 / math.Pow(1+(r/400), 1.0/3.0)))
			return mv, mv - (t * mi * 12)
		})

	revRd1 := calcBlock("Reverse RD - Compute monthly deposit", [3]string{"Maturity Value", "Annual rate of interest", "Period of investment in years"}, scroll2, "Monthly investment", "Interest",
		func(mv, r, t float64) (float64, float64) {
			mi := mv * (1 - (1 / math.Pow(1+(r/400), 1.0/3.0))) / (math.Pow(1+(r/400), 4*t) - 1)
			return mi, mv - (t * mi * 12)
		})

	revRd2 := calcBlock("Reverse RD - Compute period of investment", [3]string{"Maturity Value", "Annual rate of interest", "Monthly investment"}, scroll2, "Period of investment in years", "Interest",
		func(mv, r, mi float64) (float64, float64) {
			t := (math.Log(((mv/mi)*(1-(1/math.Pow(1+(r/400), 1.0/3.0))))+1) / math.Log(1+(r/400))) / 4
			return t, mv - (mi * t * 12)
		})

	// --- Construct Pages ---
	scroll1.Content = container.NewVBox(
		widget.NewButton("Switch to RD", func() { w.SetContent(scroll2); w.Resize(fyne.NewSize(400, 700)) }),
		fdBlock, revFd1, revFd2,
	)

	scroll2.Content = container.NewVBox(
		widget.NewButton("Switch to FD", func() { w.SetContent(scroll1); w.Resize(fyne.NewSize(400, 700)) }),
		rdBlock, revRd1, revRd2,
	)

	w.SetContent(scroll1)
	w.Resize(fyne.NewSize(400, 700))
	w.ShowAndRun()
}
