# BOS Folder and Single Spreadsheet Setup

## Overview

The script automatically manages:
- **BOS Folder**: A dedicated folder in Google Drive for all analytics
- **Single Google Sheet**: One spreadsheet with multiple tabs (pages) instead of separate files

## Automatic Setup

### First Time Setup

1. **Run `initializeBOS()`** - This is the only setup function you need!

   ```javascript
   initializeBOS()
   ```

2. **What it does:**
   - Creates "BOS" folder in your Google Drive root
   - Creates "Email Revenue Analytics" spreadsheet inside BOS folder
   - Creates all 8 tabs (pages) in the spreadsheet:
     - Raw Data
     - Clients
     - Monthly Revenue
     - Yearly Summary
     - Pattern Discovery
     - Client Matching
     - Processing Log
     - Sample Analysis
   - Sets up headers and formatting for each tab
   - Stores folder and spreadsheet IDs for future use

3. **Check the execution log** - It will show:
   - BOS Folder ID
   - Spreadsheet ID
   - Spreadsheet URL (click to open)

### After Initial Setup

Once `initializeBOS()` has run:
- The script remembers the folder and spreadsheet
- All functions automatically use the correct spreadsheet
- No need to manually configure IDs
- The spreadsheet is stored in the BOS folder (organized!)

## Folder Structure

```
Google Drive/
└── BOS/
    └── Email Revenue Analytics (Google Sheet)
        ├── Raw Data (tab)
        ├── Clients (tab)
        ├── Monthly Revenue (tab)
        ├── Yearly Summary (tab)
        ├── Pattern Discovery (tab)
        ├── Client Matching (tab)
        ├── Processing Log (tab)
        └── Sample Analysis (tab)
```

## Benefits of Single Spreadsheet

✅ **Easier Navigation**: All data in one place  
✅ **Better Organization**: All tabs in one file  
✅ **Easier Charting**: Reference multiple tabs in charts  
✅ **Simpler Sharing**: Share one file instead of many  
✅ **Better Performance**: Faster to work with one file  

## Accessing Your Spreadsheet

After running `initializeBOS()`, you can:

1. **From the execution log**: Click the spreadsheet URL
2. **From Google Drive**: Navigate to BOS folder → Email Revenue Analytics
3. **From Apps Script**: The script automatically uses it

## Re-initializing

If you need to start fresh:

1. **Delete the existing spreadsheet** (optional)
2. **Run `initializeBOS()` again**
3. It will create a new spreadsheet in the BOS folder

## Customization

### Change Folder Name

Edit this line in the script:
```javascript
const FOLDER_NAME = 'BOS'; // Change to your preferred name
```

### Change Spreadsheet Name

Edit this line:
```javascript
const SPREADSHEET_NAME = 'Email Revenue Analytics'; // Change to your preferred name
```

### Add More Tabs

1. Add tab name to `SHEETS` object:
   ```javascript
   const SHEETS = {
     // ... existing tabs
     NEW_TAB: 'New Tab Name'
   };
   ```

2. Add setup function:
   ```javascript
   function setupNewTabSheet(sheet) {
     if (sheet.getLastRow() === 0) {
       sheet.getRange(1, 1, 1, 5).setValues([[
         'Column 1', 'Column 2', 'Column 3', 'Column 4', 'Column 5'
       ]]);
       sheet.getRange(1, 1, 1, 5).setFontWeight('bold');
     }
   }
   ```

3. Call it in `setupAllSheets()`:
   ```javascript
   setupNewTabSheet(ss.getSheetByName(SHEETS.NEW_TAB));
   ```

## Troubleshooting

### "Cannot find spreadsheet"
- Run `initializeBOS()` first
- Check execution log for errors

### "Folder not found"
- The script will create it automatically
- Make sure you have Google Drive access

### "Permission denied"
- Authorize the script when prompted
- Grant Drive and Sheets permissions

### Want to use existing spreadsheet?
- You can manually move an existing spreadsheet into the BOS folder
- Update the script properties (advanced - not recommended)

## Notes

- The BOS folder is created in your Google Drive root
- The spreadsheet is removed from root and only exists in BOS folder
- All IDs are stored in script properties (persists across runs)
- You can have multiple BOS folders if needed (change FOLDER_NAME)
