# go-systems
My golang codes

A collection of independent Go tools and experiments — desktop apps, Android utilities, file encryption, image processing, and audio capture. Built on Bazzite Linux using an Ubuntu 24.04 Distrobox container.

---

## Projects

### 🔗 t5code — Map-Based Linked List (CLI prototype)
`t5code.go`

A terminal prototype for a map-based ordered linked list. Each key is a string; each value is a `[prevIndex, currentIndex]` int pair. Supports append and delete with automatic reindexing. The data structure that underpins the image viewer and PDF reader. This concept is later used to build Multi-PDF apk using Dart and flutter which is available as a pre-release.

```
go run t5code.go
Commands: type any string to add | !Del! <key> to remove | !Exit! to quit
```

---

### 📝 t2code + browser — T-pad Text Editor
`t2code.go` · `browser.go`

A Fyne-based desktop text editor with:
- Custom orange/warm theme
- Save / Load with a file dialog (`.txt`, `.md`)
- Debounced auto-save — 1 second after the last keystroke, only if a file path is already set
- Mutex-protected file path state
- "Search the web" button launches a companion `browser` binary (a `webview_go` window) resolved relative to the executable path

```
go build -o tpad t2code.go
go build -o browser browser.go
./tpad
```

> The browser binary must sit alongside the tpad binary. WebKit helper processes can't be fully bundled in an AppImage due to hardcoded system paths.

---

### 💰 t12code — FD / RD Financial Calculator
`t12code.go`

A mobile-friendly financial calculator built with Fyne, targeting Android. All results update live on every keystroke.

**FD (Fixed Deposit)**
- Forward: principal + rate + period → maturity value, interest
- Reverse 1: rate + period + target maturity → required principal
- Reverse 2: principal + rate + target maturity → investment period

**RD (Recurring Deposit)**
- Forward: monthly deposit + rate + period → maturity value, interest
- Reverse 1: target maturity + rate + period → required monthly deposit
- Reverse 2: target maturity + rate + monthly deposit → investment period

Uses quarterly compounding throughout: `M = P × (1 + r/400)^(4t)`

Includes a custom `mobileEntry` widget that overrides `FocusGained` to scroll the tapped field above the Android keyboard after a 250ms delay — something Fyne doesn't handle by default.

```
go run t12code.go                          # desktop
fyne package -os android -appID com.yourname.fincalc   # Android APK
```

---

###  scrambler + unscramble — File Scrambler
`scrambler.go` · `unscramble.go`

A two-file encryption utility that scrambles any file's bytes using a geometric transposition scheme.

**How it works:**
1. File bytes are laid out into a near-square 2D grid (rows ≈ √n)
2. Rows are interleaved top-bottom (upper pointer ↓, lower pointer ↑ alternating)
3. The result is transposed, then columns get the same interleave treatment
4. Final output is re-transposed and serialized

Produces two output files:
- `<filename>scrambled.txt` — scrambled byte stream, prefixed with an 8-byte random session token
- `<filename>key.txt` — same token + original byte count (4 bytes) + position map (4 bytes per position)

Decryption verifies the token match, reads the position map, and uses goroutines (capped at 2000 concurrent) to restore each byte in parallel.

```
go run scrambler.go      # prompts for input file
go run unscramble.go     # prompts for scrambled file + key file
```

---

###  t33code — Text-on-Image Renderer
`t33code.go`

Takes an image and a text/code file, then re-renders the image by replacing each pixel with a 5×5 bitmap glyph drawn in that pixel's original color. The glyph characters cycle through the content of the text file, stripping whitespace first.

- Custom hand-coded 5×5 bitmap font covering full QWERTY: a–z, A–Z, 0–9, and all punctuation
- Each character cell is 6×6 px (5px glyph + 1px gap), so a 100×100 image becomes 600×600 output
- Output is PNG regardless of input format (JPEG and PNG both accepted)
- Loops the text file content if the image has more pixels than characters

```
go run t33code.go
# Prompts: input image → text source file → output filename
```

---

### 🎵 YouTubeAudioD — YouTube Audio Recorder
`YouTubeAudioD.go`

A CLI tool that searches YouTube, opens the selected video in the system browser, then records the system audio output to FLAC using `ffmpeg` and PulseAudio monitor source.

1. Takes a video name as input
2. Scrapes `ytInitialData` JSON from the YouTube search results page
3. Displays up to 10 results for selection
4. Opens the chosen URL via `xdg-open` (Linux) / `open` (macOS) / `start` (Windows)
5. Records `@DEFAULT_SINK@.monitor` via `ffmpeg -f pulse` until interrupted

Output: `<VideoTitle>.flac` in the current directory.

> **Note:** Ads that play before the video will also be recorded. Work in progress.

```
go run YouTubeAudioD.go
# Requires: ffmpeg with PulseAudio support
```

---

## System

- **OS:** Bazzite Linux (`bazzite-gnome-nvidia-open`)
- **Dev container:** Ubuntu 24.04 Distrobox
- **Go:** 1.26.4 (project toolchain) / 1.22 (system)
- **Fyne CLI:** used for Android packaging
- **Android SDK/NDK:** `~/android` inside container

---

## Dependencies

| Project | Dependencies |
|---|---|
| t2code / browser | `fyne.io/fyne/v2`, `github.com/webview/webview_go` |
| t12code | `fyne.io/fyne/v2` |
| t33code | stdlib only |
| scrambler / unscramble | stdlib only |
| YouTubeAudioD | stdlib only + `ffmpeg` (system) |

---

*All projects are independent — each has its own `main` package and can be built or run separately.*
