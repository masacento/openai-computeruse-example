package computeruse

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
)

// Browser represents a browser instance for automation
type Browser struct {
	browser *rod.Browser
	page    *rod.Page
	width   int
	height  int
}

// NewBrowser creates a new browser instance with the specified dimensions
func NewBrowser(width, height int) *Browser {
	browser := rod.New().MustConnect()
	return &Browser{browser, nil, width, height}
}

// Close closes the browser instance
func (b *Browser) Close() {
	b.browser.MustClose()
}

// Open opens a URL in the browser
func (b *Browser) Open(url string) error {
	page, err := b.browser.Page(proto.TargetCreateTarget{URL: url})
	if err != nil {
		return fmt.Errorf("error opening page: %w", err)
	}
	page.MustSetViewport(b.width, b.height, 1, false)
	page.MustWaitStable()
	b.page = page
	return nil
}

// Screenshot takes a screenshot of the current page
func (b *Browser) Screenshot() ([]byte, error) {
	screenshot, err := b.page.Screenshot(false, nil)
	if err != nil {
		return nil, fmt.Errorf("error taking screenshot: %w", err)
	}
	return screenshot, nil
}

// GetCurrentUrl returns the current URL of the page
func (b *Browser) GetCurrentUrl() string {
	return b.page.MustInfo().URL
}

// Keypress simulates pressing keys on the keyboard
func (b *Browser) Keypress(keys []string) {
	keyb := b.page.Keyboard
	for _, key := range keys {
		switch strings.ToLower(key) {
		case "enter":
			keyb.Press(input.Enter)
		case "return":
			keyb.Press(input.Enter)
		case "delete":
			keyb.Press(input.Delete)
		case "tab":
			keyb.Press(input.Tab)
		case "escape":
			keyb.Press(input.Escape)
		case "left":
			keyb.Press(input.ArrowLeft)
		case "right":
			keyb.Press(input.ArrowRight)
		case "up":
			keyb.Press(input.ArrowUp)
		case "down":
			keyb.Press(input.ArrowDown)
		case "page_up":
			keyb.Press(input.PageUp)
		case "page_down":
			keyb.Press(input.PageDown)
		default:
			fmt.Printf("key: %v is not implemented", key)
		}
	}
	b.page.MustWaitStable()
}

// Type types text into the active element
func (b *Browser) Type(text string) {
	page := b.page
	page.InsertText(text)
}

// Move moves the mouse to the specified coordinates
func (b *Browser) Move(x, y int) {
	mouse := b.page.Mouse
	mouse.MustMoveTo(float64(x), float64(y))
}

// Click clicks at the specified coordinates with the specified button
func (b *Browser) Click(x, y int, button string) {
	mouse := b.page.Mouse
	mouse.MustMoveTo(float64(x), float64(y))

	switch button {
	case "right":
		mouse.MustDown("right")
		mouse.MustUp("right")
	default: // "left" is default
		mouse.MustDown("left")
		mouse.MustUp("left")
	}
	b.page.MustWaitStable()
}

// DoubleClick double-clicks at the specified coordinates
func (b *Browser) DoubleClick(x, y int) {
	mouse := b.page.Mouse
	mouse.MustMoveTo(float64(x), float64(y))
	mouse.MustClick("left")
	mouse.MustClick("left")
	b.page.MustWaitStable()
}

// Scroll scrolls the page at the specified coordinates
func (b *Browser) Scroll(x, y, scrollX, scrollY int) {
	mouse := b.page.Mouse
	mouse.MustMoveTo(float64(x), float64(y))
	b.page.Mouse.MustScroll(float64(scrollX), float64(scrollY))
	b.page.MustWaitStable()
}

// Wait waits for the specified number of milliseconds
func (b *Browser) Wait(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

// Drag performs a drag operation along the specified path
func (b *Browser) Drag(path []map[string]int) {
	fmt.Println("Drag not implemented")
}
