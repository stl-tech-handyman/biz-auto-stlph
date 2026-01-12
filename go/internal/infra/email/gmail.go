
// GetMessage retrieves an email message from Gmail by message ID
func (g *GmailSender) GetMessage(ctx context.Context, messageID string) (string, error) {
	userEmail := g.from
	
	// Get the message with full format to get HTML body
	msg, err := g.service.Users.Messages.Get(userEmail, messageID).Format("full").Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to get message: %w", err)
	}
	
	// Extract HTML body from message parts
	var htmlBody string
	if msg.Payload != nil {
		htmlBody = g.extractHTMLFromPayload(msg.Payload)
	}
	
	if htmlBody == "" {
		return "", fmt.Errorf("no HTML body found in message")
	}
	
	return htmlBody, nil
}

// extractHTMLFromPayload recursively extracts HTML content from message payload
func (g *GmailSender) extractHTMLFromPayload(payload *gmail.MessagePart) string {
	if payload == nil {
		return ""
	}
	
	// Check if this part has HTML content
	if payload.MimeType == "text/html" && payload.Body != nil && payload.Body.Data != "" {
		data, err := base64.URLEncoding.DecodeString(payload.Body.Data)
		if err == nil {
			return string(data)
		}
	}
	
	// Recursively check parts
	if payload.Parts != nil {
		for _, part := range payload.Parts {
			if html := g.extractHTMLFromPayload(part); html != "" {
				return html
			}
		}
	}
	
	return ""
}
