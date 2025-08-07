/**
 * Sends an HTML email from an authorized sender with validation.
 *
 * @param {string} fromAddress - Must be defined in ALLOWED_SENDERS
 * @param {string} toAddress - Email address of the recipient
 * @param {string} htmlBody - HTML content of the email
 * @param {string} [subject] - Optional custom subject
 * @returns {Object} - Status, sender info, and timestamp
 */
function sendEmail(fromAddress, toAddress, htmlBody, subject) {
  const senderMeta = ALLOWED_SENDERS[fromAddress];
  if (!senderMeta) {
    throw new Error(`Unauthorized sender: ${fromAddress}`);
  }

  if (!toAddress || !/^[^@]+@[^@]+\.[^@]+$/.test(toAddress)) {
    throw new Error(`Invalid recipient email: ${toAddress}`);
  }

  if (!htmlBody || htmlBody.trim() === "") {
    throw new Error("htmlBody must not be empty");
  }

  const finalSubject = subject || `[${senderMeta.prefix}] Default Subject`;

  GmailApp.sendEmail(toAddress, finalSubject, "Plain text fallback", {
    name: senderMeta.name,
    htmlBody: htmlBody,
    from: fromAddress,
    headers: { "Content-Type": "text/html; charset=UTF-8" }
  });

  return {
    status: "sent",
    to: toAddress,
    from: fromAddress,
    subject: finalSubject,
    senderName: senderMeta.name,
    timestamp: new Date().toISOString()
  };
}
