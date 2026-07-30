# my programs

A collection of small, self-contained Go programs covering file utilities, media tools, and desktop apps.

---

## Tools

### `scrambler.go` / `unscrambler.go` — File Scrambler

Scrambles any file's byte content using a matrix transposition and interleave algorithm, producing a scrambled output and a separate key file. The key is required to restore the original.

**How it works:**
- Lays file bytes into a 2D matrix (roughly square)
- Applies two rounds of row interleaving and transposition
- Prepends an 8-byte random header to both output files for pair validation
- The key file stores the original byte positions; without it, the scrambled file is unreadable

**Usage:**
```bash
go run scrambler.go
# Enter filename when prompted → produces <file>scrambled.txt and <file>key.txt

go run unscrambler.go
# Enter scrambled file and key file when prompted → produces unscrambled_<original>
```

---

### `text-embed.go` — Text-in-Image Renderer

Takes an image and a text file, then renders the text content onto the image pixel-by-pixel using a custom 5×5 bitmap font. Each pixel of the source image becomes a character cell, coloured to match the original pixel. Output is a PNG scaled 6× the source resolution.

Supports the full ASCII printable character set (a–z, A–Z, 0–9, punctuation).

**Usage:**
```bash
go run text-embed.go
# Prompts for: image file (PNG/JPEG), text file (.txt/.md), output filename
# Output: <name>.png
```

**Notes:**
- If the text is shorter than the pixel count, it wraps
- If the text is longer than the pixel count, you are warned and can choose to continue or abort
- Exit at any prompt with `~Exit~`

---

### `YouTubeAudioD.go` — YouTube Audio Downloader

Searches YouTube by name via `yt-dlp`, presents the top 10 results, then either downloads the audio directly or opens the video in your default browser and records system audio — your choice per video.

**Dependencies:** `yt-dlp` (recommended), `ffmpeg` (for legacy record mode)

**Usage:**
```bash
go run YouTubeAudioD.go
# Type a search query → pick from results → choose:
#   1. Direct download via yt-dlp (FLAC, no ads)
#   2. Browser + system audio record (legacy, Linux/PulseAudio only)
```

Exit with `!Exi!T`.

**Notes:**
- Direct download mode (option 1) is clean — no ads, no recording noise
- Legacy record mode captures `@DEFAULT_SINK@.monitor` (PulseAudio, Linux only); advertisements will also be captured
- Falls back gracefully if `yt-dlp` is not installed

---

### `financial_FD_RD.go` — FD/RD Financial Calculator

A mobile-friendly desktop calculator for Indian Fixed Deposit (FD) and Recurring Deposit (RD) instruments, built with [Fyne](https://fyne.io/). Supports both forward calculations (given principal/deposit → maturity) and reverse calculations (given maturity → derive the unknown variable).

**Calculations supported:**

| Mode | Inputs | Outputs |
|---|---|---|
| FD | Principal, rate, years | Maturity value, interest earned |
| Reverse FD | Maturity value, rate, years | Required principal, interest |
| Reverse FD | Maturity value, principal, rate | Investment period, interest |
| RD | Monthly deposit, rate, years | Maturity value, interest earned |
| Reverse RD | Maturity value, rate, years | Required monthly deposit, interest |
| Reverse RD | Maturity value, rate, monthly deposit | Investment period, interest |

Results update live as you type. The FD and RD pages are separate screens with a toggle button between them.

**Dependencies:** `fyne.io/fyne/v2`

**Usage:**
```bash
go run financial_FD_RD.go
```

---

### `browser.go` — Minimal WebView Browser

A bare-bones desktop browser using [`webview/webview_go`](https://github.com/webview/webview_go). Opens a URL in an 800×600 native webview window. Used internally by `tpad` for its web search button.

**Dependencies:** `github.com/webview/webview_go` and its native WebKit dependencies

**Usage:**
```bash
go run browser.go [url]
# Defaults to https://google.com if no URL is given
```

Or use the pre-built binary:
```bash
./browser https://example.com
```

---

### `tpad.go` — Minimal Text Editor (T-pad)

A lightweight desktop text editor built with [Fyne](https://fyne.io/). Features a warm custom colour theme, debounced auto-save, file open/save dialogs, and a "Search the web" button that launches the bundled `browser` binary.

**Dependencies:** `fyne.io/fyne/v2`, the compiled `browser` binary in the same directory as the `tpad` binary

**Features:**
- Auto-saves 1 second after the last keystroke (once a file path is known)
- Supports `.txt` and `.md` files
- Custom warm theme (terracotta/sand palette)
- Integrated web search via the `browser` binary

**Usage:**
```bash
go run tpad.go
```

Or use the pre-built binary (place `browser` in the same directory):
```bash
./tpad
```

---

## Building

Each tool is its own `package main`. Build individually:

```bash
go build -o scrambler scrambler.go
go build -o unscrambler unscrambler.go
go build -o text-embed text-embed.go
go build -o ytdl YouTubeAudioD.go
go build -o browser browser.go
go build -o tpad tpad.go
```

For `tpad`, the `browser` binary must be in the same directory as the `tpad` binary at runtime.

---

## Requirements

| Tool | External deps |
|---|---|
| `scrambler` / `unscrambler` | None |
| `text-embed` | None |
| `YouTubeAudioD` | `yt-dlp`, `ffmpeg` (legacy mode), PulseAudio (legacy mode) |
| `financial_FD_RD` | Fyne native deps |
| `browser` | WebKit/GTK (Linux), WebView2 (Windows), WKWebView (macOS) |
| `tpad` | Fyne native deps, `browser` binary |

Go 1.18+ recommended.
