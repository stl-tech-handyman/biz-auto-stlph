# Email Template Compatibility Guide

## Overview

This document explains the email compatibility requirements and testing approach for all email templates in the system. Email templates must render consistently across all major email clients, including:

- **Desktop**: Gmail, Outlook (Windows/Mac), Apple Mail, Yahoo Mail
- **Mobile**: iPhone Mail, Gmail App (iOS/Android), Outlook Mobile
- **Web**: Gmail Web, Outlook Web, Yahoo Web

## Email Compatibility Principles

### 1. **Table-Based Layout**
- ✅ **Use**: HTML tables for layout structure
- ❌ **Avoid**: `<div>`, CSS Grid, Flexbox
- **Why**: Tables are the most reliable layout method across all email clients

### 2. **Inline Styles Only**
- ✅ **Use**: Inline `style=""` attributes on every element
- ❌ **Avoid**: `<style>` blocks, external stylesheets, CSS classes
- **Why**: Many email clients strip out `<style>` blocks and external CSS

### 3. **Simple CSS Properties**
- ✅ **Use**: Basic properties: `font-size`, `color`, `padding`, `margin`, `background-color`, `border`, `text-align`
- ❌ **Avoid**: Advanced CSS: `transform`, `position: absolute/fixed`, `z-index`, `box-shadow`, complex selectors
- **Why**: Limited CSS support varies widely across email clients

### 4. **Font Fallbacks**
- ✅ **Use**: `font-family: Arial, Helvetica, sans-serif`
- ❌ **Avoid**: Custom fonts, web fonts, `@font-face`
- **Why**: System fonts ensure consistent rendering

### 5. **Absolute URLs**
- ✅ **Use**: `https://example.com/image.jpg`
- ❌ **Avoid**: Relative URLs, protocol-relative URLs (`//example.com`)
- **Why**: Email clients may not resolve relative URLs correctly

### 6. **No JavaScript**
- ❌ **Never use**: `<script>`, `onclick`, `javascript:`, event handlers
- **Why**: All email clients block JavaScript for security

### 7. **Mobile Viewport**
- ✅ **Use**: `<meta name="viewport" content="width=device-width, initial-scale=1.0" />`
- **Why**: Ensures proper mobile rendering

### 8. **Readable Font Sizes**
- ✅ **Use**: Minimum 12px for body text, 14px+ for important content
- ❌ **Avoid**: Font sizes below 11px
- **Why**: Small text is unreadable on mobile devices

### 9. **Width Constraints**
- ✅ **Use**: `max-width: 600px` for main container
- ✅ **Use**: `width: 100%` for responsive tables
- **Why**: Prevents horizontal scrolling on mobile

### 10. **Color Contrast**
- ✅ **Use**: High contrast colors (black text on white background)
- ❌ **Avoid**: Light gray text on white backgrounds
- **Why**: Ensures readability across all devices and email clients

## Testing Requirements

### Automated Tests

All email templates must pass the compatibility test suite located in:
- `go/internal/util/quote_email_template_test.go` (example)
- Similar tests should be created for all email templates

#### Test Coverage

1. **Structure Tests**
   - ✅ Uses table-based layout
   - ✅ No `<div>` tags
   - ✅ Inline styles only
   - ✅ No CSS classes

2. **Content Tests**
   - ✅ All required content present
   - ✅ No format errors (`%!d`, `%!s`, etc.)
   - ✅ Proper data substitution

3. **Compatibility Tests**
   - ✅ UTF-8 encoding specified
   - ✅ Viewport meta tag present
   - ✅ Font fallbacks specified
   - ✅ Max-width constraint set
   - ✅ No JavaScript

4. **Mobile Tests**
   - ✅ Viewport meta tag
   - ✅ Responsive table widths
   - ✅ Readable font sizes (12px+)

### Manual Testing Checklist

Before deploying any email template changes, manually test in:

#### Desktop Clients
- [ ] Gmail (Chrome)
- [ ] Gmail (Firefox)
- [ ] Outlook 2016/2019 (Windows)
- [ ] Outlook (Mac)
- [ ] Apple Mail (Mac)
- [ ] Yahoo Mail (Web)

#### Mobile Clients
- [ ] Gmail App (iPhone)
- [ ] Gmail App (Android)
- [ ] iPhone Mail App
- [ ] Outlook Mobile (iPhone)
- [ ] Outlook Mobile (Android)

#### Test Scenarios
- [ ] Email renders correctly on first open
- [ ] All images load (if applicable)
- [ ] All links work correctly
- [ ] Text is readable (not too small)
- [ ] Layout doesn't break on mobile
- [ ] Colors display correctly
- [ ] No horizontal scrolling on mobile

## Common Issues and Solutions

### Issue: Different rendering on iPhone vs Desktop
**Solution**: 
- Use simpler table structure
- Remove complex CSS
- Use `cellpadding` and `cellspacing` instead of CSS padding/margin where possible
- Test on actual devices

### Issue: Fonts look different
**Solution**:
- Use system fonts only (Arial, Helvetica, sans-serif)
- Specify font-size in pixels (px), not em or rem
- Avoid custom fonts

### Issue: Images not loading
**Solution**:
- Use absolute URLs (https://)
- Host images on reliable CDN
- Provide alt text for all images
- Consider embedding small images as base64 (with caution)

### Issue: Links not working
**Solution**:
- Use absolute URLs
- Test all links before sending
- Avoid JavaScript in links
- Use simple `<a href="...">` tags

### Issue: Colors look different
**Solution**:
- Use hex colors (#000000, not black)
- Test in multiple clients
- Avoid transparency/opacity
- Use high contrast

## Template Structure Best Practices

### Recommended HTML Structure

```html
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Email Title</title>
  </head>
  <body style="margin:0; padding:0; font-family: Arial, Helvetica, sans-serif; background-color: #ffffff; color: #333333;">
    <table width="100%" cellpadding="0" cellspacing="0" border="0" style="background-color: #ffffff; width: 100%;">
      <tr>
        <td align="center" style="padding: 10px;">
          <table width="100%" cellpadding="0" cellspacing="0" border="0" style="max-width: 600px; border: 1px solid #cccccc; padding: 15px; background-color: #ffffff;">
            <!-- Content here -->
          </table>
        </td>
      </tr>
    </table>
  </body>
</html>
```

### Key Points
- Outer table: Full width container
- Inner table: Max-width 600px for desktop, 100% width for mobile
- All styles inline
- Simple, nested table structure

## Tools for Testing

### Recommended Tools
1. **Litmus** - Comprehensive email testing across clients
2. **Email on Acid** - Email preview and testing
3. **Mailtrap** - Safe email testing environment
4. **Browser DevTools** - For initial testing
5. **Actual Devices** - Final verification on real phones

### Free Alternatives
- Gmail "Show Original" - View raw HTML
- Browser responsive mode - Simulate mobile
- Multiple email accounts - Test across providers

## Running Compatibility Tests

```bash
# Run all email template compatibility tests
cd go
go test ./internal/util/... -v -run TestQuoteEmailTemplate

# Run specific test
go test ./internal/util/... -v -run TestQuoteEmailTemplateCompatibility
```

## Template Version History

- **v1.0**: Initial template
- **v1.1**: AIDA structure implementation
- **v1.2**: Email compatibility improvements
  - Removed emojis from data fields
  - Simplified table structure
  - Standardized inline styles
  - Added viewport meta tag
  - Fixed parameter order
  - Improved mobile compatibility

## Checklist for New Templates

When creating a new email template:

- [ ] Uses table-based layout
- [ ] All styles are inline
- [ ] No CSS classes
- [ ] No JavaScript
- [ ] UTF-8 encoding specified
- [ ] Viewport meta tag included
- [ ] Font fallbacks specified
- [ ] Max-width 600px set
- [ ] All URLs are absolute
- [ ] Font sizes 12px or larger
- [ ] High contrast colors
- [ ] Compatibility tests written
- [ ] Manual testing completed
- [ ] Documentation updated

---

**Last Updated**: 2025-12-24  
**Maintained By**: Development Team  
**Related Docs**: 
- `QUOTE_EMAIL_TEMPLATE_STRUCTURE.md` - Marketing rationale
- Email template test files in `go/internal/util/`

