/**
 * Script to delete Google Apps Script projects
 * 
 * Usage:
 * 1. Enable Apps Script API: https://script.google.com/home/usersettings
 * 2. Run: node delete-apps-script-projects.js
 * 
 * Or use the Apps Script API directly via curl/gcloud
 */

// Project IDs to delete (from the image you showed)
const PROJECTS_TO_DELETE = [
  // "stlph-app-scripts-api-468501", // This is a GCP project ID, not Apps Script ID
  // "stlph-prod" // This is also a GCP project ID
];

// Note: The IDs shown in the image are GCP project IDs, not Apps Script script IDs
// To find the actual Apps Script IDs, you need to:
// 1. Open each project in script.google.com
// 2. Go to Project Settings
// 3. Copy the Script ID

console.log(`
To delete Google Apps Script projects:

METHOD 1: Via Web Interface (Easiest)
1. Go to https://script.google.com
2. Find each project in the list
3. Click the three dots menu (⋮) next to each project
4. Select "Delete project"
5. Confirm deletion

METHOD 2: Via Apps Script API
1. Enable Apps Script API: https://script.google.com/home/usersettings
2. Get OAuth token: gcloud auth print-access-token
3. Find the Script ID from each project's settings
4. Use this curl command:

curl -X DELETE \\
  "https://script.googleapis.com/v1/projects/{SCRIPT_ID}" \\
  -H "Authorization: Bearer $(gcloud auth print-access-token)"

METHOD 3: Using clasp (if you have the script IDs)
Unfortunately, clasp doesn't support deletion directly.
You'll need to use the web interface or API.

NOTE: The project names you showed are:
- "STLPH App Scripts API" (stlph-app-scripts-api-468501)
- "STLPH Prod" (stlph-prod)

These appear to be GCP project IDs, not Apps Script script IDs.
To delete Apps Script projects, you need the Script ID, not the GCP project ID.
`);



