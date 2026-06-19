package main

import (
	"image/color"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func openBrowser(url string) {
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	browserpath := filepath.Join(exeDir, "browser")
	exec.Command(browserpath, url).Start()
}

// --- Auto-save state ---
var (
	currentFilePath string     // empty means file has not been saved yet
	saveMu          sync.Mutex // protects currentFilePath
	saveTimer       *time.Timer
)

// saveToPath writes content to the given file path.
func saveToPath(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// scheduleSave triggers a debounced auto-save 1 second after the last keystroke.
// Only runs if a file path is already known.
func scheduleSave(content string) {
	saveMu.Lock()
	path := currentFilePath
	saveMu.Unlock()

	if path == "" {
		return // no file path yet, wait for manual save
	}

	if saveTimer != nil {
		saveTimer.Stop()
	}
	saveTimer = time.AfterFunc(1*time.Second, func() {
		saveToPath(path, content)
	})
}

// --- Theme ---
type CTheme struct{}

func (c CTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNamePrimary:
		return color.RGBA{R: 200, G: 120, B: 80, A: 255}
	case theme.ColorNameSelection:
		return color.RGBA{R: 200, G: 100, B: 100, A: 100}
	case theme.ColorNamePlaceHolder:
		return color.RGBA{R: 0, G: 0, B: 0, A: 155}
	case theme.ColorNameBackground:
		return color.RGBA{R: 224, G: 172, B: 105, A: 255}
	case theme.ColorNameInputBackground:
		return color.RGBA{R: 247, G: 200, B: 187, A: 255}
	case theme.ColorNameForeground:
		return color.RGBA{R: 20, G: 20, B: 20, A: 255}
	case theme.ColorNameFocus:
		return color.RGBA{R: 0, G: 0, B: 0, A: 255}
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (t CTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t CTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t CTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func main() {
	a := app.New()
	w := a.NewWindow("T-pad!")
	a.Settings().SetTheme(CTheme{})

	textBox := widget.NewMultiLineEntry()
	textBox.SetPlaceHolder("Start typing here")
	textBox.TextStyle.Italic = false

	// --- Save button ---
	saveButton := widget.NewButton("Save", func() {
		saveMu.Lock()
		path := currentFilePath
		saveMu.Unlock()

		if path != "" {
			// File already has a path — just save silently
			go saveToPath(path, textBox.Text)
			return
		}

		// No path yet — open save dialog
		fd := dialog.NewFileSave(func(uc fyne.URIWriteCloser, err error) {
			if err != nil || uc == nil {
				return
			}
			uc.Close()
			newPath := uc.URI().Path()
			saveMu.Lock()
			currentFilePath = newPath
			saveMu.Unlock()
			go saveToPath(newPath, textBox.Text)
			w.SetTitle("T-pad! — " + filepath.Base(newPath))
		}, w)
		fd.SetFileName("untitled.txt")
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".txt", ".md"}))
		fd.Show()
	})

	// --- Load button ---
	loadButton := widget.NewButton("Load", func() {
		fd := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
			if err != nil || uc == nil {
				return
			}
			defer uc.Close()
			data, err := os.ReadFile(uc.URI().Path())
			if err != nil {
				return
			}
			newPath := uc.URI().Path()
			saveMu.Lock()
			currentFilePath = newPath
			saveMu.Unlock()
			textBox.SetText(string(data))
			w.SetTitle("T-pad! — " + filepath.Base(newPath))
		}, w)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".txt", ".md"}))
		fd.Show()
	})

	// --- Search button ---
	searchButton := widget.NewButton("Search the web", func() {
		go openBrowser("https://google.com")
	})

	// Auto-save on every keystroke (only if file path is known)
	textBox.OnChanged = func(content string) {
		go scheduleSave(content)
	}

	buttons := container.NewHBox(saveButton, loadButton, searchButton)

	w.Resize(fyne.NewSize(800, 500))
	w.SetContent(container.NewBorder(buttons, nil, nil, nil, textBox))
	w.ShowAndRun()
}
