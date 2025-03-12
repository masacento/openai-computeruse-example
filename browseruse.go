package computeruse

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"
)

// BrowserUse automates browser interactions using OpenAI's computer-use model
// Parameters:
// - url: The URL to open in the browser
// - instruction: The instruction to send to the AI model
// Returns an error if any operation fails
func BrowserUse(ctx context.Context, url, instruction string, maxTurns int) error {
	model := "computer-use-preview-2025-03-11"

	browser := NewBrowser(1024, 768)
	err := browser.Open(url)
	if err != nil {
		return fmt.Errorf("error opening browser: %w", err)
	}
	defer browser.Close()

	var responseID string
	var callID string
	var callResp *ComputerOutput

	for i := 0; i < maxTurns; i++ {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context canceled")
		default:
		}

		messages := []Input{}
		if responseID == "" {
			messages = append(messages, Input{
				Role:    "user",
				Content: instruction,
			})
		} else {
			messages = append(messages, Input{
				Type:   "computer_call_output",
				CallID: callID,
				Output: callResp,
			})
		}

		debugInput(messages)
		response, err := Responses(model, responseID, messages)
		if err != nil {
			return fmt.Errorf("error calling OpenAI API: %w", err)
		}
		debugResponse(response)

		responseID = response.ID

		finalOutput := ""
		for _, o := range response.Output {
			if o.Action != nil {
				var err error
				callResp, err = computerCall(browser, o.Action)
				if err != nil {
					return fmt.Errorf("error executing browser action: %w", err)
				}
				callID = o.CallID
				if len(o.PendingSafetyChecks) > 0 {
					fmt.Println("pending safety checks:", o.PendingSafetyChecks)
				}
				debugComputerOutput(callResp)
			}
			if o.Content != nil {
				if o.Role == "assistant" {
					finalOutput = fmt.Sprint(o.Content[0])
					break
				}
			}
		}

		if finalOutput != "" {
			fmt.Println("Final output:", finalOutput)
			break
		}
		time.Sleep(1 * time.Second)
	}

	return nil
}

// computerCall executes a browser action and returns the resulting output
func computerCall(b *Browser, action *Action) (*ComputerOutput, error) {
	switch action.Type {
	case "screenshot":
		// Just take a screenshot, no additional action needed
	case "type":
		b.Type(action.Text)
	case "click":
		b.Click(action.X, action.Y, action.Button)
	case "scroll":
		b.Scroll(action.X, action.Y, action.ScrollX, action.ScrollY)
	case "keypress":
		b.Keypress(action.Keys)
	case "wait":
		time.Sleep(3 * time.Second)
	}

	screenshot, err := b.Screenshot()
	if err != nil {
		return nil, fmt.Errorf("error taking screenshot: %w", err)
	}
	return &ComputerOutput{
		Type:       "input_image",
		ImageURL:   dataURL(screenshot),
		CurrentURL: b.GetCurrentUrl(),
	}, nil
}

// dataURL converts binary data to a base64-encoded data URL
func dataURL(data []byte) string {
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(data)
}

// debugResponse formats and displays Response details
func debugResponse(response *Response) {
	fmt.Println("\n📩 ----- RESPONSE DETAILS -----")
	fmt.Printf("🆔 Response ID: %s\n", response.ID)
	fmt.Printf("📊 Status: %s\n", response.Status)

	if len(response.Output) > 0 {
		fmt.Printf("📦 Output items: %d\n", len(response.Output))

		for i, o := range response.Output {
			fmt.Printf("\n📦 Output item #%d:\n", i+1)

			if o.Action != nil {
				fmt.Println("🎮 ----- BROWSER ACTION -----")
				fmt.Printf("  Type: %s\n", o.Action.Type)

				if o.Action.Text != "" {
					textPreview := o.Action.Text
					if len(textPreview) > 50 {
						textPreview = textPreview[:47] + "..."
					}
					fmt.Printf("  Text: %s\n", textPreview)
				}

				if o.Action.Button != "" {
					fmt.Printf("  Button: %s\n", o.Action.Button)
				}

				if len(o.Action.Keys) > 0 {
					fmt.Printf("  Keys: %v\n", o.Action.Keys)
				}

				if o.Action.X != 0 || o.Action.Y != 0 {
					fmt.Printf("  Position: (%d, %d)\n", o.Action.X, o.Action.Y)
				}

				if o.Action.ScrollX != 0 || o.Action.ScrollY != 0 {
					fmt.Printf("  Scroll: (%d, %d)\n", o.Action.ScrollX, o.Action.ScrollY)
				}

				fmt.Println("  --------------------------")
			}

			if o.Content != nil && o.Role == "assistant" {
				fmt.Println("🤖 ----- ASSISTANT RESPONSE -----")
				for j, content := range o.Content {
					fmt.Printf("  Content #%d: %s\n", j+1, content)
				}
				fmt.Println("  ------------------------------")
			}

			if len(o.PendingSafetyChecks) > 0 {
				fmt.Println("⚠️ ----- PENDING SAFETY CHECKS -----")
				for _, check := range o.PendingSafetyChecks {
					fmt.Printf("  %s: %s\n", check.Code, check.Message)
				}
				fmt.Println("  ---------------------------------")
			}
		}
	}

	fmt.Println("📩 ----- END OF RESPONSE DETAILS -----\n")
}

// debugComputerOutput saves the screenshot from ComputerOutput to a file
func debugComputerOutput(out *ComputerOutput) {
	dataurl := out.ImageURL
	if dataurl == "" {
		fmt.Println("📷 No screenshot available")
		return
	}

	database64 := strings.Split(dataurl, ",")[1]
	data, err := base64.StdEncoding.DecodeString(database64)
	if err != nil {
		fmt.Printf("❌ Error decoding screenshot: %v\n", err)
		return
	}

	// Create filename with timestamp
	os.MkdirAll("screenshots", 0755)
	filename := fmt.Sprintf("screenshots/%s.png", time.Now().Format("20060102150405"))

	// Save the file
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Printf("❌ Error saving screenshot: %v\n", err)
		return
	}

	fmt.Printf("📷 Screenshot saved: %s\n", filename)

	// Log browser state if available
	if out.CurrentURL != "" {
		fmt.Printf("🌐 Current URL: %s\n", out.CurrentURL)
	}
	if out.Type != "" {
		fmt.Printf("📊 Output type: %s\n", out.Type)
	}
}

// debugInput prints input message details for debugging
func debugInput(input []Input) {
	fmt.Println("\n📥 ----- INPUT MESSAGE DETAILS -----")

	for i, v := range input {
		fmt.Printf("📌 Message #%d:\n", i+1)

		if v.Role != "" {
			fmt.Printf("  🔹 Role: %s\n", v.Role)
		}

		if v.Type != "" {
			fmt.Printf("  🔹 Type: %s\n", v.Type)
		}

		if v.CallID != "" {
			fmt.Printf("  🔹 Call ID: %s\n", v.CallID)
		}

		if v.Content != "" {
			contentPreview := v.Content
			if len(contentPreview) > 100 {
				contentPreview = contentPreview[:97] + "..."
			}
			fmt.Printf("  🔹 Content: %s\n", contentPreview)
		}

		if v.Output != nil {
			fmt.Println("  🔹 Output details:")
			if v.Output.CurrentURL != "" {
				fmt.Printf("    - URL: %s\n", v.Output.CurrentURL)
			}
			if v.Output.Type != "" {
				fmt.Printf("    - Type: %s\n", v.Output.Type)
			}
		}

		fmt.Println("  ------------------------------")
	}

	fmt.Println("📥 ----- END OF INPUT DETAILS -----\n")
}
