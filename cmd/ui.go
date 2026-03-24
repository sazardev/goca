package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// UIRenderer provides styled terminal output for the CLI.
// It centralizes all user-facing output with consistent styling.
type UIRenderer struct {
	writer      io.Writer
	noColor     bool
	interactive bool
	verbosity   int // 0=quiet, 1=normal, 2=verbose
}

// Global UI instance used by all commands
var ui *UIRenderer

// Color palette
var (
	colorGreen   = lipgloss.Color("#00d787")
	colorRed     = lipgloss.Color("#ff5f87")
	colorYellow  = lipgloss.Color("#ffd700")
	colorCyan    = lipgloss.Color("#00d7ff")
	colorMagenta = lipgloss.Color("#d787ff")
	colorDim     = lipgloss.Color("#6c6c6c")
	colorWhite   = lipgloss.Color("#ffffff")
)

// NewUIRenderer creates a new UIRenderer.
// If noColor is true, all styling is disabled.
// verbosity: 0=quiet (only Success/Error), 1=normal, 2=verbose (adds Debug/Trace)
func NewUIRenderer(writer io.Writer, noColor bool, verbosity int) *UIRenderer {
	if writer == nil {
		writer = os.Stdout
	}
	if noColor || os.Getenv("NO_COLOR") != "" {
		lipgloss.SetColorProfile(termenv.Ascii)
	}
	return &UIRenderer{
		writer:      writer,
		noColor:     noColor || os.Getenv("NO_COLOR") != "",
		interactive: true,
		verbosity:   verbosity,
	}
}

// SetInteractive controls whether interactive prompts are enabled
func (u *UIRenderer) SetInteractive(interactive bool) {
	u.interactive = interactive
}

// IsInteractive returns whether interactive prompts are enabled
func (u *UIRenderer) IsInteractive() bool {
	return u.interactive
}

// Header prints a bold header line
func (u *UIRenderer) Header(text string) {
	if u.verbosity < 1 {
		return
	}
	style := lipgloss.NewStyle().Bold(true)
	fmt.Fprintln(u.writer, style.Render(text))
}

// Step prints a numbered step indicator
func (u *UIRenderer) Step(number int, text string) {
	if u.verbosity < 1 {
		return
	}
	numStyle := lipgloss.NewStyle().Foreground(colorCyan).Bold(true)
	fmt.Fprintf(u.writer, "%s %s\n", numStyle.Render(fmt.Sprintf("%d.", number)), text)
}

// Success prints a success message with a green checkmark
func (u *UIRenderer) Success(text string) {
	prefix := lipgloss.NewStyle().Foreground(colorGreen).Bold(true).Render("✓")
	fmt.Fprintf(u.writer, "%s %s\n", prefix, text)
}

// Error prints an error message with a red cross
func (u *UIRenderer) Error(text string) {
	prefix := lipgloss.NewStyle().Foreground(colorRed).Bold(true).Render("✗")
	fmt.Fprintf(u.writer, "%s %s\n", prefix, text)
}

// Warning prints a warning message with a yellow indicator
func (u *UIRenderer) Warning(text string) {
	if u.verbosity < 1 {
		return
	}
	prefix := lipgloss.NewStyle().Foreground(colorYellow).Bold(true).Render("⚠")
	fmt.Fprintf(u.writer, "%s %s\n", prefix, text)
}

// Info prints an informational message with a cyan indicator
func (u *UIRenderer) Info(text string) {
	if u.verbosity < 1 {
		return
	}
	prefix := lipgloss.NewStyle().Foreground(colorCyan).Render("ℹ")
	fmt.Fprintf(u.writer, "%s %s\n", prefix, text)
}

// DryRun prints a dry-run prefixed message
func (u *UIRenderer) DryRun(text string) {
	if u.verbosity < 1 {
		return
	}
	tag := lipgloss.NewStyle().
		Foreground(colorMagenta).
		Bold(true).
		Render("[DRY-RUN]")
	fmt.Fprintf(u.writer, "%s %s\n", tag, text)
}

// FileCreated prints a file creation message with a styled path
func (u *UIRenderer) FileCreated(path string) {
	check := lipgloss.NewStyle().Foreground(colorGreen).Render("✓")
	dimPath := lipgloss.NewStyle().Foreground(colorDim).Render(path)
	fmt.Fprintf(u.writer, "  %s Created: %s\n", check, dimPath)
}

// FileBackedUp prints a file backup message
func (u *UIRenderer) FileBackedUp(from, to string) {
	arrow := lipgloss.NewStyle().Foreground(colorCyan).Render("→")
	fmt.Fprintf(u.writer, "  Backed up: %s %s %s\n", from, arrow, to)
}

// KeyValue prints a key-value pair with styled key
func (u *UIRenderer) KeyValue(key, value string) {
	if u.verbosity < 1 {
		return
	}
	k := lipgloss.NewStyle().Foreground(colorCyan).Render(key + ":")
	fmt.Fprintf(u.writer, "%s %s\n", k, value)
}

// KeyValueFromConfig prints a key-value pair with a "from config" annotation
func (u *UIRenderer) KeyValueFromConfig(key, value string) {
	if u.verbosity < 1 {
		return
	}
	k := lipgloss.NewStyle().Foreground(colorCyan).Render(key + ":")
	dimTag := lipgloss.NewStyle().Foreground(colorDim).Render("(from config)")
	fmt.Fprintf(u.writer, "%s %s %s\n", k, value, dimTag)
}

// Feature prints a feature toggle line (e.g., "✓ Including validation")
func (u *UIRenderer) Feature(text string, fromConfig bool) {
	if u.verbosity < 1 {
		return
	}
	check := lipgloss.NewStyle().Foreground(colorGreen).Render("✓")
	if fromConfig {
		dimTag := lipgloss.NewStyle().Foreground(colorDim).Render("(from config)")
		fmt.Fprintf(u.writer, "%s %s %s\n", check, text, dimTag)
	} else {
		fmt.Fprintf(u.writer, "%s %s\n", check, text)
	}
}

// Table prints a formatted table with headers and rows
func (u *UIRenderer) Table(headers []string, rows [][]string) {
	if len(headers) == 0 {
		return
	}

	// Calculate column widths
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Build separator
	sepParts := make([]string, len(widths))
	for i, w := range widths {
		sepParts[i] = strings.Repeat("─", w+2)
	}
	separator := "├" + strings.Join(sepParts, "┼") + "┤"
	topBorder := "┌" + strings.Join(sepParts, "┬") + "┐"
	bottomBorder := "└" + strings.Join(sepParts, "┴") + "┘"

	// Print top border
	borderStyle := lipgloss.NewStyle().Foreground(colorDim)
	fmt.Fprintln(u.writer, borderStyle.Render(topBorder))

	// Print header
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(colorCyan)
	headerCells := make([]string, len(headers))
	for i, h := range headers {
		padded := h + strings.Repeat(" ", widths[i]-len(h))
		headerCells[i] = " " + headerStyle.Render(padded) + " "
	}
	fmt.Fprintln(u.writer, borderStyle.Render("│")+strings.Join(headerCells, borderStyle.Render("│"))+borderStyle.Render("│"))

	// Print separator
	fmt.Fprintln(u.writer, borderStyle.Render(separator))

	// Print rows
	for _, row := range rows {
		cells := make([]string, len(headers))
		for i := range headers {
			val := ""
			if i < len(row) {
				val = row[i]
			}
			padded := val + strings.Repeat(" ", widths[i]-len(val))
			cells[i] = " " + padded + " "
		}
		fmt.Fprintln(u.writer, borderStyle.Render("│")+strings.Join(cells, borderStyle.Render("│"))+borderStyle.Render("│"))
	}

	// Print bottom border
	fmt.Fprintln(u.writer, borderStyle.Render(bottomBorder))
}

// Println prints a plain line
func (u *UIRenderer) Println(text string) {
	fmt.Fprintln(u.writer, text)
}

// Printf prints a formatted string
func (u *UIRenderer) Printf(format string, args ...any) {
	fmt.Fprintf(u.writer, format, args...)
}

// Blank prints an empty line
func (u *UIRenderer) Blank() {
	if u.verbosity < 1 {
		return
	}
	fmt.Fprintln(u.writer)
}

// Dim prints dimmed text
func (u *UIRenderer) Dim(text string) {
	if u.verbosity < 1 {
		return
	}
	style := lipgloss.NewStyle().Foreground(colorDim)
	fmt.Fprintln(u.writer, style.Render(text))
}

// Section prints a section with a title and indented content
func (u *UIRenderer) Section(title string) {
	if u.verbosity < 1 {
		return
	}
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(colorWhite).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(colorDim)
	fmt.Fprintln(u.writer, style.Render(title))
}

// NextSteps prints a formatted "Next steps" block
func (u *UIRenderer) NextSteps(steps []string) {
	if u.verbosity < 1 {
		return
	}
	title := lipgloss.NewStyle().Bold(true).Foreground(colorCyan).Render("Next steps:")
	fmt.Fprintln(u.writer, title)
	for _, step := range steps {
		fmt.Fprintf(u.writer, "  %s\n", step)
	}
}

// initUI initializes the global UI renderer.
// verbosity: 0=quiet, 1=normal, 2=verbose
func initUI(noColor bool, verbosity int) {
	ui = NewUIRenderer(os.Stdout, noColor, verbosity)
}

// Debug prints a debug message (verbosity >= 2)
func (u *UIRenderer) Debug(text string) {
	if u.verbosity < 2 {
		return
	}
	prefix := lipgloss.NewStyle().Foreground(colorDim).Render("[debug]")
	fmt.Fprintf(u.writer, "%s %s\n", prefix, text)
}

// Trace prints a trace message (verbosity >= 2)
func (u *UIRenderer) Trace(text string) {
	if u.verbosity < 2 {
		return
	}
	prefix := lipgloss.NewStyle().Foreground(colorDim).Render("[trace]")
	fmt.Fprintf(u.writer, "%s %s\n", prefix, text)
}

// Spinner starts a spinner animation with the given text.
// Returns a stop function that should be called when the operation is complete.
// The stop function prints a completion message with a green checkmark.
func (u *UIRenderer) Spinner(text string) func() {
	if u.noColor {
		fmt.Fprintf(u.writer, "%s... ", text)
		return func() {
			fmt.Fprintln(u.writer, "done")
		}
	}

	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	style := lipgloss.NewStyle().Foreground(colorCyan)

	var once sync.Once
	done := make(chan struct{})

	go func() {
		i := 0
		for {
			select {
			case <-done:
				return
			default:
				fmt.Fprintf(u.writer, "\r%s %s", style.Render(frames[i%len(frames)]), text)
				i++
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()

	return func() {
		once.Do(func() {
			close(done)
			check := lipgloss.NewStyle().Foreground(colorGreen).Bold(true).Render("✓")
			fmt.Fprintf(u.writer, "\r%s %s\n", check, text)
		})
	}
}
