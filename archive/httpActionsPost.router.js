function handlePostRequest(e) {
  const { action } = JSON.parse(e.postData.contents || '{}');
  const version = getVersionFromPath(e);

  switch (version) {
    case 'v1':
      return handlePostV1(action, e);
    case 'v2':
      return handlePostV2(action, e);  // add this stub now
    default:
      return ContentService.createTextOutput(`❌ Unsupported version: ${version}`);
  }
}

function handlePostV1(action, e) {
  switch (action) {
    case 'sendQuote':
      return handle_POST_sendQuote(e);
    case 'createInvoice':
      return handle_POST_createInvoice(e);
    case 'parseLead':
      return handle_POST_parseLead(e);
    // 🔄 Add more actions as needed
    default:
      return ContentService.createTextOutput(`Unknown v1 action: ${action}`);
  }
}

function handlePostV2(action, e) {
  switch (action) {
    case 'createInvoice':
      return handle_POST_createInvoice_V2(e);
    default:
      return ContentService.createTextOutput(`❌ Unknown v2 action: ${action}`);
  }
}
