const ESTIMATE_SENT_CALENDAR_ID =
  "c_f8c0098141f20b9bcb25d5e3c05d54c450301eb4f21bff9c75a04b1612138b54@group.calendar.google.com";
const FOLDER_ID = "1qH-4Vq6aLHhnM7jpUD3SbMLX8SW89oDl"; // Folder where leads spreadsheet is stored

const LEAD_ID_AUTO = "lead/auto-id-v1";
const LEAD_ID_BACKFILL = "lead/backfill";

const TEST_SRE_EMAIL = "test-sre@stlpartyhelpers.com";

const DataSource = Object.freeze({
  ZAPIER: "zappier",
  EMAIL_LABEL_EXTRACTING: "email_label_extracing"
});

// Constants: Booking Deposit options mapped by price
const BOOKING_DEPOSITS = [
  { value: 5000, id: "price_1RXLiuIzH4MDwV7sXM4GCCl1" },
  { value: 4500, id: "price_1RXLiUIzH4MDwV7sJLT2V4n2" },
  { value: 4000, id: "price_1RXLi3IzH4MDwV7sp60jdyp6" },
  { value: 3000, id: "price_1RXLg6IzH4MDwV7sXCz5vBeU" },
  { value: 2500, id: "price_1RXLf1IzH4MDwV7sJayxv7ho" },
  { value: 2000, id: "price_1RXLeiIzH4MDwV7scjIAiDnD" },
  { value: 1250, id: "price_1RXLePIzH4MDwV7sQrDa6Y4S" },
  { value: 1000, id: "price_1RXLgrIzH4MDwV7szxFYM3yU" },
  { value: 750,  id: "price_1RXLe0IzH4MDwV7sze3Ym7tg" },
  { value: 500,  id: "price_1RXLYOIzH4MDwV7sWTRoBMDZ" },
  { value: 400,  id: "price_1RXLaOIzH4MDwV7sPMmgHwO5" },
  { value: 350,  id: "price_1RXLa3IzH4MDwV7sJaLxBkqH" },
  { value: 300,  id: "price_1RXLZgIzH4MDwV7sgLr1hvGG" },
  { value: 250,  id: "price_1RXLYzIzH4MDwV7sB3cKw5ug" },
  { value: 200,  id: "price_1RbpRfIzH4MDwV7swxFBma8P" },
  { value: 150,  id: "price_1RXLU5IzH4MDwV7sdR9It89Z" },
  { value: 100,  id: "price_1ReLk0IzH4MDwV7smKRJTGkj" }, 
  { value: 50,   id: "price_1RXLtPIzH4MDwV7s6JKkKRT7" },
];

const PublicGetActions = {
  TEST_HEALTHCHECK: 'TEST_YAGOOD',
  TEST_HEALTHCHECK_V1: 'TEST_YAGOOD_v1',
  TEST_EMAIL_SENDING: 'TEST_EMAIL_SENDING',
  TEST_EMAIL_SENDING_QUOTE_EMAIL: 'TEST_EMAIL_SENDING_QUOTE_EMAIL',
  TEST_EMAIL_QUOTESENDING_WITHFORWARDING: 'TEST_EMAIL_QUOTESENDING_WITHFORWARDING',
  GET_LAT_LNG: 'GET_LAT_LNG',
  CALCULATE_ESTIMATE: 'CALCULATE_ESTIMATE',
  STRIPE_GET_BOOKING_DEPOSIT_AMOUNT: 'STRIPE_GET_BOOKING_DEPOSIT_AMOUNT',
};

const TEST_HTMLBODY = "<table style='width:100%; border:1px solid black;'><tr><td><b>Event Helper</b></td><td>3 hours</td></tr></table>";



function test_internal_doGet_healthcheckV1() {
  const mockEvent = {
    parameter: {
      action: PublicGetActions.TEST_HEALTHCHECK_V1,
      "req-initiator": "internal_test"
    }
  };

  const response = doGet(mockEvent);
  Logger.log(response.getContent());  // should show: "YEZZIR!" or whatever your response is
}


// Publicly Available
function getLatLng(address) {
  const apiKey = 'AIzaSyCBLoZHDCSKmYarkvNJht4-qAHnQtA7GBQ';
  const url = `https://maps.googleapis.com/maps/api/geocode/json?address=${encodeURIComponent(address)}&key=${apiKey}`;
  
  const response = UrlFetchApp.fetch(url);
  const json = JSON.parse(response.getContentText());
  
  if (json.status === 'OK') {
    const location = json.results[0].geometry.location;
    return {
      lat: location.lat,
      lng: location.lng,
      fullAddress: json.results[0].formatted_address
    };
  } else {
    throw new Error(`Geocoding failed: ${json.status}`);
  }
}

function getUpdateFrequency(daysUntilEvent) {
  if (daysUntilEvent >= 548) return 90; // 1.5+ Years (Every 3 months)
  if (daysUntilEvent >= 365) return 60; // 1 Year (Every 2 months)
  if (daysUntilEvent >= 180) return 30; // 6 Months (Every 1 month)
  if (daysUntilEvent >= 90) return 14; // 3 Months (Every 2 weeks)
  if (daysUntilEvent >= 30) return 7; // 1 Month (Every week)
  if (daysUntilEvent >= 14) return 3; // 2 Weeks (Every 3 days)
  if (daysUntilEvent >= 7) return 2; // 1 Week (Every 2 days)
  if (daysUntilEvent >= 3) return 1; // 3 Days (Every day)
  return 0; // 24 Hours (Immediate update)
}

function testEmailQuotePreview() {
  const clientName = "Alex Shev";
  const eventDate = "August 28, 2025";
  const eventTime = "1:00 PM";
  const eventLocation = "4220 Duncan Ave., St. Louis, MO 63110";
  const occasion = "Birthday Party";
  const guestCount = 50;
  const helpers = 2;
  const hours = 5;
  const rate = 45;
  const total = "$490";

  const htmlBody = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8" />
    <title>STL Party Helpers - Quote</title>
  </head>
  <body style="margin:0; padding:0; font-family: Arial, sans-serif; background-color: #ffffff; color: #333;">
    <table width="100%; padding: 10px;" cellpadding="0" cellspacing="2" border="0" style="background-color: #ffffff;">
      <tr>
        <td align="center" style="padding: 0 16px;">
          <table width="100%" cellpadding="0" cellspacing="0" border="0" style="max-width: 650px; border: 1px solid #ccc; padding: 20px;">
            <!-- Header -->
            <tr>
              <td align="center" style="padding: 5px;">
                <img src="https://stlpartyhelpers.com/wp-content/uploads/2025/08/stlph-logo-2.jpg" alt="STL Party Helpers Logo" />
              </td>
            </tr>
            <tr>
              <td align="center" style="font-size: 22px; font-weight: bold;">Hi ${clientName}!</td>
            </tr>
            <tr>
              <td align="center" style="padding-bottom: 10px;">
                Thank you for reaching out!<br />
                Below is your event quote, along with important details and next steps.
              </td>
            </tr>

             <!-- Pricing -->
            <tr>
              <td style="font-size: 14px; font-weight: bold; padding-top: 15px;">Our Rates & Pricing</td>
            </tr>
            <tr>
              <td>
                <table width="100%" cellpadding="0" cellspacing="0" border="0" style="background-color: #f9f9f9; margin-top: 10px;">
                  <tr>
                    <td style="padding: 8px 10px; font-weight: bold;">Base Rate:</td>
                    <td style="padding: 8px 10px;">$200 / helper (first 4 hours)</td>
                  </tr>
                  <tr>
                    <td style="padding: 8px 10px; font-weight: bold;">Additional Hours:</td>
                    <td style="padding: 8px 10px;">$${rate} per additional hour per helper</td>
                  </tr>
                  <tr>
                    <td style="padding: 8px 10px; font-weight: bold;">Estimated Total:</td>
                    <td style="padding: 8px 10px;">${total}</td>
                  </tr>
                </table>
                <p style="font-size: 12px; color: #666; padding-top: 5px;">
                  Final total may adjust based on our call. Gratuity is not included but always appreciated!
                </p>
              </td>
            </tr>

            <!-- Event Details -->
            <tr>
              <td style="font-size: 14px; font-weight: bold; padding-top: 20px;">Event Details</td>
            </tr>
            <tr>
              <td>
                <table width="100%" cellpadding="0" cellspacing="0" border="0" style="background-color: #f9f9f9; margin-top: 10px;">
                  <tr>
                    <td style="padding: 8px 10px; font-weight: bold;">&#128197; When:</td>
                    <td style="padding: 8px 10px;">${eventDate} ${eventTime}</td>
                  </tr>
                  <tr>
                    <td style="padding: 8px 10px; font-weight: bold;">&#128205; Where:</td>
                    <td style="padding: 8px 10px;">${eventLocation}</td>
                  </tr>
                  <tr>
                    <td style="padding: 8px 10px; font-weight: bold;">&#127760; Occasion:</td>
                    <td style="padding: 8px 10px;">${occasion}</td>
                  </tr>
                  <tr>
                    <td style="padding: 8px 10px; font-weight: bold;">&#128101; Guest Count:</td>
                    <td style="padding: 8px 10px;">${guestCount}</td>
                  </tr>
                  <tr>
                    <td style="padding: 8px 10px; font-weight: bold;">&#129491; Helpers Needed:</td>
                    <td style="padding: 8px 10px;">${helpers}</td>
                  </tr>
                  <tr>
                    <td style="padding: 8px 10px; font-weight: bold;">&#9201; For How Long:</td>
                    <td style="padding: 8px 10px;">${hours} Hours</td>
                  </tr>
                </table>
              </td>
            </tr>

           

            <!-- Services -->
            <tr>
              <td style="font-size: 14px; font-weight: bold; padding-top: 15px;">Services Included</td>
            </tr>
            <tr>
              <td style="padding: 10px 0;">
                <ul style="padding-left: 20px; margin: 0;">
                  <li><strong>Setup & Presentation</strong>
                    <ul>
                      <li>Arranging tables, chairs, and decorations</li>
                      <li>Buffet setup & live buffet service</li>
                      <li>Butler-passed appetizers & cocktails</li>
                    </ul>
                  </li>
                  <li><strong>Dining & Guest Assistance</strong>
                    <ul>
                      <li>Multi-course plated dinners</li>
                      <li>General bussing (plates, silverware, glassware)</li>
                      <li>Beverage service (water, wine, champagne, coffee, etc.)</li>
                      <li>Special services (cake cutting, dessert plating, etc.)</li>
                    </ul>
                  </li>
                  <li><strong>Cleanup & End-of-Event Support</strong>
                    <ul>
                      <li>Washing dishes, managing trash, and keeping the event space tidy</li>
                      <li>Kitchen cleanup & end-of-event breakdown</li>
                      <li>Assisting with food storage & leftovers</li>
                    </ul>
                  </li>
                </ul>
                <p>Need something specific? Let us know! We’ll do our best to accommodate your request.</p>
              </td>
            </tr>

            <!-- Payment Options -->
            <tr>
              <td style="font-size: 14px; font-weight: bold; padding-top: 15px;">Payment Options</td>
            </tr>
            <tr>
              <td style="background-color: #f9f9f9; padding: 10px;">
                Check, Debit / Credit Cards (via Stripe), Venmo, Zelle, Apple Pay
              </td>
            </tr>

            <!-- Next Steps -->
            <tr>
              <td style="font-size: 14px; font-weight: bold; padding-top: 15px;">What Happens Next</td>
            </tr>
           <tr>

  </tr>
  <tr>
    <td style="background-color: #f9f9f9; padding: 10px;">
      <span style="text-align:center; font-size: 14px; font-weight: bold;">Booked already?</span><br />
      <table cellpadding="0" cellspacing="0" border="0" style="font-size: 14px;">
        <tr>
          <td valign="top" style="padding-right: 8px;">1.</td>
          <td>We’ll call you at your scheduled time to go over details.</td>
        </tr>
        <tr>
          <td valign="top" style="padding-right: 8px;">2.</td>
          <td>If all looks good after our call, we’ll send a Stripe deposit link to proceed.</td>
        </tr>
        <tr>
          <td valign="top" style="padding-right: 8px;">3.</td>
          <td>Once the deposit is in, your reservation is locked in.</td>
        </tr>
      </table>
      <p style="font-size: 13px; text-align: center; color: #666; margin-top: 8px;">
        Deposit is 40–50% of the estimate rounded for simplicity. 
      </p> <p style="font-size: 13px; text-align: center; color: #666; margin-top: 5px;">
      ❌ Required to confirm your reservation.
      </p>
    </td>
  </tr>
              <tr>
                <td style="background-color: #fff4e5; text-align: center; padding: 10px; margin-top: 5px; border: 1px solid #fddfb4; ">
                  <strong>Haven’t scheduled a call yet?</strong><br />
                  <strong>Book now to get started</strong><br />
                  <span style="font-size: 0.9em; padding-bottom: 20px;">(to confirm helpers, tasks, and setup)</span><br />
                  <a href="https://calendly.com/stlpartyhelpers/quote-intake" style="display:inline-block; background-color:#0047ab; color:#fff; padding:8px 14px; margin-top: 12px; text-decoration:none; font-weight:bold; border-radius:4px;">Click Here to Schedule Appointment</a><br/>
                
                </td> 
              </tr>
  <!-- Footer -->
  <tr>
    <td align="center" style="font-size: 12px; padding-top: 20px; color: #666;">
    
      4220 Duncan Ave., Ste. 201, St. Louis, MO 63110<br />
      <a href="tel:+13147145514" style="display:inline-block;background-color:#ffffff;color:#000000;padding: 9px 10px;text-decoration:none;border-radius:4px;margin-top:8px;border: 1px solid gray;margin-top: 12px;margin-bottom: 12px;" target="_blank">Tap to Call Us: (314) 714-5514</a><br />
      <a href="https://stlpartyhelpers.com" style="color:#0047ab; display: inline-block; margin-bottom: 8px;">stlpartyhelpers.com</a>
      <br />
      &copy; 2025 STL Party Helpers<br />
      <span style="font-size: 0.55em;">v1.1</span>
    </td>
  </tr>
          </table>
        </td>
      </tr>
    </table>
  </body>
</html>`;

  GmailApp.sendEmail("alexey@shevelyov.com", "[Preview] Full HTML Quote Template", "Fallback", {
    htmlBody,
    from: "test-sre@stlpartyhelpers.com",
    headers: {
    "Content-Type": "text/html; charset=UTF-8"
  }
  });

  Logger.log("✅ Test email sent.");
}


function testEmailSending(fromAddress, toAddress, htmlBody) {
  const allowedFroms = GmailApp.getAliases();

  if (!allowedFroms.includes(fromAddress)) {
    throw new Error("Alias not allowed: " + fromAddress + ". Allowed aliases: " + allowedFroms.join(", "));
  }

  const subject = "[test-sre] testEmailSending(fromAddress, toAddress, htmlBody)";
  
  GmailApp.sendEmail(toAddress, subject, "Plain text fallback", {
    name: `Site Reliability Team`,
    htmlBody: htmlBody,
    from: fromAddress,
      headers: { "Content-Type": "text/html; charset=UTF-8" }
  });

  return {
    status: 'sent',
    to: toAddress,
    from: fromAddress,
    timestamp: new Date().toISOString()
  };
}

/*
function testSendQuoteEmail(
  fromAddress,
  toAddress,
  helpers = 1,
  hours = 3,
  rate = 45,
  total = null,
  notes = `This is a test quote sent on ${new Date().toLocaleString()}`,
  clientName = "Test Client",
  eventDate = "December 13, 2025",
  eventTime = "02:00 PM",
  eventLocation = "942 Guelbreth Ln, St. Louis, MO 63146",
  occasion = "Holiday Party",
  guestCount = 100
) {
  const subject = "[test-sre] Test Quote Email - Original";

  // Calculate total if not explicitly provided
  if (total === null || isNaN(total)) {
    const result = calculateTotalCost_v1(helpers, hours, {
      baseRate: 200,
      hourlyRate: rate,
      baseHours: 4,
      minHours: 4
    });

    if (result.error) {
      throw new Error(`Failed to calculate total: ${result.message}`);
    }

    total = result.totalCost;
  }

  // Send structured object
  const htmlBody = generateQuoteEmail_v1({
    clientName,
    eventDate,
    eventTime,
    eventLocation,
    occasion,
    guestCount,
    helpers,
    hours,
    rate,
    total,
    notes
  });

  const htmlBody1 = `<!-- Full Quote Email HTML Refactored -->
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Quote Email</title>
    <style>
      body {
        font-family: Arial, sans-serif;
        margin: 0;
        padding: 0;
        background: #ffffff;
        color: #333;
      }
      .container {
        width: 100%;
        max-width: 650px;
        margin: auto;
        border: 1px solid #ccc;
        padding: 20px;
      }
      h2 {
        color: #222;
      }
      .section {
        margin-bottom: 25px;
      }
      .highlight {
        background-color: #f9f9f9;
        padding: 10px;
        border-radius: 5px;
      }
      .header {
        text-align: center;
      }
      .footer {
        text-align: center;
        margin-top: 40px;
        font-size: 0.9em;
        color: #666;
      }
      .button {
        display: inline-block;
        background-color: #0047ab;
        color: #fff;
        padding: 12px 20px;
        margin-top: 20px;
        border-radius: 5px;
        text-decoration: none;
        font-weight: bold;
      }
      .tag {
        background: #fef4e5;
        padding: 10px;
        border-radius: 4px;
        margin-top: 15px;
        border: 1px solid #fddfb4;
      }
      .row {
        display: flex;
        justify-content: space-between;
        margin: 5px 0;
      }
      .label {
        font-weight: bold;
      }
    </style>
  </head>
  <body>
    <div class="container">
      <div class="header">
        <img src="https://stlpartyhelpers.com/wp-content/uploads/2024/02/STL-Party-Helpers-Logo.webp" width="200" />
        <h2>Hi {{clientName}}!</h2>
        <p>Thank you for reaching out!</p>
        <p>Below is your event quote, along with important details and next steps.</p>
      </div>

      <div class="section">
        <h3>Event Details</h3>
        <div class="highlight">
          <div class="row"><span class="label">When</span> {{eventDate}} {{eventTime}}</div>
          <div class="row"><span class="label">📍 Where</span> {{eventLocation}}</div>
          <div class="row"><span class="label">🎉 Occasion</span> {{occasion}}</div>
          <div class="row"><span class="label">👥 Guest Count</span> {{guestCount}}</div>
          <div class="row"><span class="label">🧍‍♀️ Helpers Needed</span> {{helpers}}</div>
          <div class="row"><span class="label">⏱️ For How Long</span> {{hours}} Hours</div>
        </div>
      </div>

      <div class="section">
        <h3>Our Rates & Pricing</h3>
        <div class="highlight">
          <div class="row"><span class="label">💲 Base Rate:</span> $200 / helper (first 4 hours)</div>
          <div class="row"><span class="label">⏳ Additional Hours:</span> ${{rate}} per additional hour per helper</div>
          <div class="row"><span class="label">📊 Estimated Total:</span> ${{total}}</div>
          <p style="font-size: 0.9em; color: #666;">Final total may adjust based on our call. Gratuity is not included but always appreciated!</p>
        </div>
      </div>

      <div class="section">
        <h3>Services Included</h3>
        <div class="highlight">
          <strong>Setup & Presentation</strong>
          <ul>
            <li>Arranging tables, chairs, and decorations</li>
            <li>Buffet setup & live buffet service</li>
            <li>Butler-passed appetizers & cocktails</li>
          </ul>

          <strong>Dining & Guest Assistance</strong>
          <ul>
            <li>Multi-course plated dinners</li>
            <li>General bussing (plates, silverware, glassware)</li>
            <li>Beverage service (water, wine, champagne, coffee, etc.)</li>
            <li>Special services (cake cutting, dessert plating, etc.)</li>
          </ul>

          <strong>Cleanup & End-of-Event Support</strong>
          <ul>
            <li>Washing dishes, managing trash, and keeping the event space tidy</li>
            <li>Kitchen cleanup & end-of-event breakdown</li>
            <li>Assisting with food storage & leftovers</li>
          </ul>
          <p>Need something specific? Let us know! We’ll do our best to accommodate your request.</p>
        </div>
      </div>

      <div class="section">
        <h3>Payment Options</h3>
        <p class="highlight">Check, Debit / Credit Cards (via Stripe), Venmo, Zelle, Apple Pay</p>
      </div>

      <div class="section">
        <h3>What Happens Next</h3>
        <div class="highlight">
          <strong>📅 Booked already?</strong>
          <p>
            We’ll call you at your scheduled time to go over details.<br />
            If all looks good after our call, we’ll send a Stripe deposit link to proceed.<br />
            Once the deposit is in, your reservation is locked in.
          </p>
          <p>Deposit is 40 – 50% of the estimate rounded for simplicity.<br />❌ (required to confirm your reservation)</p>
        </div>

        <div class="tag">
          <strong>📞 Haven’t scheduled a call yet?</strong>
          <br /><strong>📅 Book now to get started</strong><br />
          <span style="font-size: 0.9em;">(to confirm helpers, tasks, and setup)</span>
          <br /><br />
          <a href="https://calendly.com/stlpartyhelpers/quote-intake" class="button">Click Here to Schedule Appointment</a>
        </div>
      </div>

      <div class="footer">
        <p>© 2025 STL Party Helpers<br />
        <a href="https://stlpartyhelpers.com">stlpartyhelpers.com</a><br />
        4220 Duncan Ave., Ste. 201<br />
        St. Louis, MO 63110</p>
        <p>
          <a class="button" href="tel:+13147145514">Tap to Call Us: (314) 714-5514</a>
        </p>
        <p style="margin-top: 20px; font-size: 0.7em;">v1.00</p>
      </div>
    </div>
  </body>
</html>
`;

  GmailApp.sendEmail(toAddress, subject, "Plain text fallback", {
    name: `Site Reliability Team`,
    htmlBody1,
    from: fromAddress,
    headers: { "Content-Type": "text/html; charset=UTF-8" }
  });

  return {
    status: 'sent',
    to: toAddress,
    total,
    clientName,
    eventDate,
    eventTime,
    notes
  };
}
*/

function testSendAndForwardAndDiff(toAddress) {
  testSendQuoteEmail(toAddress);
  
  const threads = GmailApp.search(`to:(${toAddress}) subject:"Test Quote Email - Original"`);
  if (threads.length === 0) throw new Error("Original email not found");

  const originalMessage = threads[0].getMessages()[0];
  const originalBody = originalMessage.getBody();

  // Forward to self
  originalMessage.forward(Session.getActiveUser().getEmail(), {
    subject: "FWD: Test Quote Email - Forwarded"
  });

  Utilities.sleep(5000); // wait for Gmail to process forwarding

  const fwdThreads = GmailApp.search(`subject:"FWD: Test Quote Email - Forwarded"`);
  if (fwdThreads.length === 0) throw new Error("Forwarded email not found");

  const forwardedMessage = fwdThreads[0].getMessages()[0];
  const forwardedBody = forwardedMessage.getBody();

  return {
    originalLength: originalBody.length,
    forwardedLength: forwardedBody.length,
    diffLength: forwardedBody.length - originalBody.length,
    originalSnippet: originalBody.substring(0, 200),
    forwardedSnippet: forwardedBody.substring(0, 200)
  };
}

function testFindBookingDepositId() {
  let passed = 0;
  let failed = 0;

  BOOKING_DEPOSITS.forEach(entry => {
    const { value, id } = entry;

    // Create a fake estimate such that this value is 45% of it
    const estimate = value / 0.45;

    // Run the function
    const resultObj = findBookingDeposit(estimate);
    const returnedId = resultObj?.id;
    const returnedValue = resultObj?.value;

    const result = returnedId === id ? "✅ PASS" : "❌ FAIL";

    Logger.log(`${result}: Estimate = ${estimate.toFixed(2)}, Expected = { id: ${id}, value: ${value} }, Got = { id: ${returnedId}, value: ${returnedValue} }`);

    if (returnedId === id) {
      passed++;
    } else {
      failed++;
    }
  });

  Logger.log(`\nTest complete. Passed: ${passed}, Failed: ${failed}`);
}

// Utility function to find the best matching deposit
function findBookingDeposit(estimate) {
  const max = estimate * 0.5;
  const min = estimate * 0.4;

  // Filter all options under or equal to 50%
  const validOptions = BOOKING_DEPOSITS.filter(opt => opt.value <= max);

  // Sort by closeness to 50% of the estimate
  validOptions.sort((a, b) => {
    const diffA = Math.abs(estimate * 0.5 - a.value);
    const diffB = Math.abs(estimate * 0.5 - b.value);
    return diffA - diffB;
  });

  // Return the closest valid match
  return validOptions[0] || null;
}

// Example usage
const deposit = findBookingDeposit(445);
console.log(deposit);
// Output: { value: 200, id: 'price_1RbpRfIzH4MDwV7swxFBma8P' }


function backfillLeads() {
  processLeadsByLabel(
    LEAD_ID_BACKFILL,
    (options = {
      sendEmail: false,
      markProcessed: false,
      markBackfilled: true,
      excludeLabels: [],
    })
  );
}

function markLeadAsBackfilled(thread) {
  try {
    var backfilledLabel =
      GmailApp.getUserLabelByName("backfilled-ok") ||
      GmailApp.createLabel("backfilled-ok");
    var backfillLabel = GmailApp.getUserLabelByName("lead/backfill");

    // ✅ Add "backfilled-ok" label
    thread.addLabel(backfilledLabel);
    console.log(
      `✅ Marked lead as 'backfilled-ok': ${thread.getFirstMessageSubject()}`
    );

    // ✅ Remove "lead/backfill" label if it exists
    if (backfillLabel) {
      thread.removeLabel(backfillLabel);
      console.log(
        `🚀 Removed 'lead/backfill' label from: ${thread.getFirstMessageSubject()}`
      );
    }
  } catch (error) {
    console.error(`❌ ERROR marking lead as backfilled: ${error.message}`);
  }
}

/**
 * Processes leads dynamically by fetching labels and excluding certain labels.
 */
function processLeadsDynamic(
  options = {
    sendEmail: true,
    markProcessed: true,
    markBackfilled: false,
    includeLabels: [],
    excludeLabels: [],
  }
) {
  console.log(`📌 Fetching labels and processing leads dynamically...`);

  var allLabels = GmailApp.getUserLabels().map((label) => label.getName());
  console.log(`📌 Found ${allLabels.length} labels.`);

  // Apply include/exclude filters
  var filteredLabels =
    options.includeLabels.length > 0
      ? allLabels.filter((label) => options.includeLabels.includes(label))
      : allLabels.filter((label) => !options.excludeLabels.includes(label));

  if (filteredLabels.length === 0) {
    console.log(`✅ No labels matched for processing.`);
    return;
  }

  console.log(
    `📬 Processing emails from ${filteredLabels.length
    } labels: ${filteredLabels.join(", ")}`
  );

  filteredLabels.forEach((labelName) => {
    try {
      processLeadsByLabel(labelName, options);
    } catch (error) {
      console.error(
        `❌ ERROR processing label '${labelName}': ${error.message}`
      );
    }
  });

  console.log(`✅ Processing complete for all applicable labels.`);
}

/**
 * Processes leads from a given label (supports normal processing & backfilling).
 */
function processLeadsByLabel(labelName, options) {
  console.log(`📌 Fetching unprocessed leads from label: ${labelName}`);

  var label = GmailApp.getUserLabelByName(labelName);
  if (!label) {
    console.error(`⚠️ ERROR: Label '${labelName}' not found.`);
    return;
  }

  var query = `label:${labelName}`;
  if (options.excludeLabels.length > 0) {
    options.excludeLabels.forEach((excludeLabel) => {
      query += ` -label:${excludeLabel}`;
    });
  }

  var threads = GmailApp.search(query);
  if (threads.length === 0) {
    console.log(`✅ No new leads found under '${labelName}'.`);
    return;
  }

  console.log(`📬 Checking ${threads.length} thread(s) for leads...`);

  threads.forEach((thread) => {
    try {
      processSingleLead(thread, options);

      if (options.markBackfilled) {
        markLeadAsBackfilled(thread);
      }

      if (options.markProcessed) {
        markLeadAsProcessed(thread);
      }
    } catch (error) {
      console.error(
        `❌ ERROR processing lead in thread '${thread.getId()}': ${error.message
        }`
      );
      markLeadAsFailed(thread, "auto-quote-sending-failed");
    }
  });

  console.log(`✅ Processing complete for label: ${labelName}`);
}

/**
 * Processes a single lead, logging it in Google Sheets and optionally sending a quote.
 */
function processSingleLead(thread, options) {
  var messages = thread.getMessages();
  var lastMessage = messages[messages.length - 1];

  var subject = lastMessage.getSubject();
  var body = lastMessage.getPlainBody();
  var sender = lastMessage.getFrom();

  console.log(`📩 Processing lead: ${subject} from ${sender}`);
  var parsedData = parseLeadData(body);

  if (!validateLeadData(parsedData, parsedData.clientName)) {
    console.error(
      `🚨 Missing required fields. Marking lead as "auto-quote-sending-failed".`
    );
    //markLeadAsFailed(thread, "auto-quote-sending-failed");
    return;
  }

  // Extract relevant data
  var {
    clientName,
    eventDate,
    eventTime,
    phone,
    location,
    numHelpers,
    duration,
    email,
    occasion,
    guestCount,
  } = parsedData;
  var emailId = email || "Not Provided";
  var threadIdForCalendar = thread ? thread.getId() : "";

  // Calculate Total Cost
  var { totalCost, baseCost, additionalCost } = calculateTotalCost(
    numHelpers,
    duration
  );
  console.log(
    `📊 Pricing Breakdown: Base Cost = $${baseCost}, Additional Cost = $${additionalCost}, Total Cost = $${totalCost}`
  );

  // ✅ Create Calendar Event
  var eventResponse = calendareventcreator.createEvent({
    calendarId: ESTIMATE_SENT_CALENDAR_ID,
    clientName,
    occasion,
    guestCount,
    eventDate,
    eventTime,
    phone,
    location,
    numHelpers,
    duration,
    totalCost,
    emailId,
    threadIdForCalendar,
  });

  if (!eventResponse || eventResponse.error) {
    console.error(
      `❌ ERROR: Failed to create calendar event. Reason: ${eventResponse.error}`
    );
    markLeadAsFailed(thread, "auto-quote-sending-failed");
    return;
  }

  var eventId = eventResponse.eventId;
  logQuoteToSheet(
    clientName,
    eventDate,
    eventTime,
    location,
    numHelpers,
    duration,
    totalCost,
    eventId,
    options
  );

  if (typeof thread !== "object" || typeof thread.getMessages !== "function") {
    throw new Error(
      "🛑 CRITICAL: thread is not a valid GmailThread before sending email."
    );
  }
  // Only send email if allowed
  if (options.sendEmail) {
    sendQuoteEmail(
      clientName,
      email,
      eventDate,
      eventTime,
      location,
      numHelpers,
      duration,
      totalCost,
      occasion,
      guestCount,
      thread
    );
  }

  // Only mark as processed if required
  if (options.markProcessed) {
    var processedLabel =
      GmailApp.getUserLabelByName("auto-quote-sent") ||
      GmailApp.createLabel("auto-quote-sent");
    thread.addLabel(processedLabel);
  }
}

/**
 * Logs the quote details to a Google Sheet stored inside a specified Google Drive folder.
 */
function logQuoteToSheet(
  clientName,
  eventDate,
  eventTime,
  location,
  numHelpers,
  duration,
  totalCost,
  eventId,
  options
) {
  var folderId = "1qH-4Vq6aLHhnM7jpUD3SbMLX8SW89oDl";
  var folder = DriveApp.getFolderById(folderId);
  var ss = getOrCreateSpreadsheet(folder);
  var sheet = ss.getSheetByName("Leads") || ss.insertSheet("Leads");

  var logType = options.markProcessed ? "New Lead Processed" : "Backfill Lead";

  sheet.appendRow([
    clientName,
    eventDate,
    eventTime,
    location,
    numHelpers,
    duration,
    totalCost,
    eventId,
    logType,
  ]);
  console.log(`📊 Logged ${logType} for ${clientName}.`);
}

/**
 * Gets or creates the Leads spreadsheet inside the specified Drive folder.
 */
function getOrCreateSpreadsheet() {
  var folder = DriveApp.getFolderById(FOLDER_ID);
  var files = folder.getFilesByType(MimeType.GOOGLE_SHEETS);

  if (files.hasNext()) {
    return SpreadsheetApp.openById(files.next().getId());
  }

  var ss = SpreadsheetApp.create("Leads");
  folder.addFile(DriveApp.getFileById(ss.getId()));
  return ss;
}

function handleStripeBookingDepositGet(e) {
  const estimateParam = e.parameter.estimate;
  const estimate = parseFloat(estimateParam);

  if (isNaN(estimate)) {
    return ContentService
      .createTextOutput(JSON.stringify({ error: "Missing or invalid 'estimate'" }))
      .setMimeType(ContentService.MimeType.JSON);
  }

  const match = findBookingDeposit(estimate);

  if (!match) {
    return ContentService
      .createTextOutput(JSON.stringify({ error: "No valid deposit found in 40–50% range" }))
      .setMimeType(ContentService.MimeType.JSON);
  }

  return ContentService
    .createTextOutput(JSON.stringify({
      estimate,
      deposit: match
    }))
    .setMimeType(ContentService.MimeType.JSON);
}

/* PUBLIC POST REQUESTS */
  const PublicPostActions = {
    SEND_ESTIMATE_ADD_TO_CALENDAR: 'send_est_add_to_cal', // we look in our gmail box
    SEND_ESTIMATE_ADD_TO_CALENDAR_FROM_ZAPPIER: 'send_est_add_to_cal_from_zappier', // zappier inputs data
    //SEND_ESTIMATE_ADD_TO_CALENDAR_FROM_ZAPPIER_V1: 'send_est_add_to_cal_from_zappier_v1', // zappier inputs data
    SEND_PAID_BOOKING_INV_SEND_EMAIL: 'send_paid_booking_inv_send_email',
    CLIENT_PAID_BOOKING_DEPOSIT: 'client_paid_booking_deposit',
    STRIPE_GET_BOOKING_DEPOSIT_AMOUNT: 'stripe_get_booking_deposit_amount',
    SEND_FINAL_INVOICE: 'send_final_invoice',
    CLIENT_PAID_FINAL_INVOICE: 'client_paid_final_invoice',
    TEST_PING: 'test_ping',
    // V1
    // SINGLE_RESP functions - accessible via POST (for external testability)
    SINGLE_RESP_CREATE_CALENDAR_EVENT: '', // create calendar to track event in Google Calendar (monday crm sucks)
    SINGLE_RESP_MOVE_CALENDAR_EVENT: '', // moving calendar between different calendars / stages
    SINGLE_RESP_DELETE_CALENDAR_EVENT: '', // event represents real event. different calendars represent stages of the event [prelim, confirmed, cancelled]
    SINGLE_RESP_GET_ESTIMATE: '', // supply event reservation  date and time, number of helpers, number of hours

    // each email will have a template
    // Gmail UI will serve as Accurate and Reliable E-Mail Visual Editor
    // E-Mail Template Keeper - all in one place. Accessable as code if needed
    SINGLE_RESP_SEND_EMAIL_QUOTE: '', // fully automated process.
    SINGLE_RESP_SEND_EMAIL_REVIEW: '', // manual process
    SINGLE_RESP_SEND_EMAIL_BOOKING_DEPOSIT_RECEIVED: '', // we can automatically send event reservation confirmation email within 5 seconds
    SINGLE_RESP_SEND_EMAIL_EVENT_RESERVATION_CONFIRMATION: '',
    SINGLE_RESP_SEND_EMAIL_EVENT_PASSED_BUSINESS_THANK_YOU_RECEIPT: '', // thank you for your business

    /* IN ACTIVE DEVELOPMENT */
    STRIPE_GENERATE_BOOKING_DEPOSIT_INVOICE: 'STRIPE_GENERATE_BOOKING_DEPOSIT_INVOICE', // In Postman
    SINGLE_RESP_SEND_EMAIL_BOOKING_DEPOSIT: 'SINGLE_RESP_SEND_EMAIL_BOOKING_DEPOSIT', // first time clients or offenders only // Not In Postman

    CONTROLLER_GENERATE_STRIPE_INVOICE_AND_SEND_BOOKING_DEPOSIT_EMAIL: 'CONTROLLER_GENERATE_STRIPE_INVOICE_AND_SEND_BOOKING_DEPOSIT_EMAIL', // aggregate functions like this will orchestrate to achieve a higher level business objective / goal / fullfill request.
//   CONTROLLER
//   NAME: CONTROLLER_MONDAY_CRM_LEAD_SEND_BOOKING_DEPOSIT_EMAIL
//   DESC: Send booking-deposit email and update CRM after Stripe invoice generation.
//
//   #   FROM  VERB    -> TO    WHAT
//   --  ----  ------- -- ----  ---------------------------------------------------
/* 0)  [CRM]  click       —    BTN: SEND_BOOKING_DEPOSIT_EMAIL                  */
/* 1)  [CRM]  emit     -> [SLK] Message posted to a specific channel            */
/* 2)  [SLK]  emit     -> [ZPR] Zapier trigger fires                            */
/* 3)  [ZPR]  call     -> [STR] Generate booking-deposit invoice                */
/* 4)  [ZPR]  call     -> [GAS] Send deposit email (uses Stripe invoice data)   */
/* 5)  [ZPR]  update   -> [CRM] Update Payment board                            */
/* 6)  [ZPR]  update   -> [CRM] Write Stripe invoice URL, due date, PDF         */
//   --  ----  ------- -- ----  ---------------------------------------------------

    CONTROLLER_COMBINE_2_OR_MORE_SINGLE_RESP_EXAMPLE: '', // aggregate functions like this will orchestrate to achieve a higher level business objective / goal / fullfill request.
  };
/**
 * Unified doPost entrypoint for WebApp
 */
function doPost(e) {
  var action = null; // define outside so catch can use it safely
  try {
    // 1) Parse incoming JSON safely
    if (!e || !e.postData || !e.postData.contents) {
      return jsonOut({ ok: false, error: "Missing postData.contents" });
    }
    var payload = JSON.parse(e.postData.contents);

    // 2) Extract action
    action = payload.action || null;
    if (!action) {
      return jsonOut({ ok: false, error: "Missing action" });
    }

    // 3) Route by action
    switch (action) {
      case PublicPostActions.SEND_ESTIMATE_ADD_TO_CALENDAR:
        return sendEstimateAndAddToCalendar(payload);

      case PublicPostActions.STRIPE_GET_BOOKING_DEPOSIT_AMOUNT:
        return getStripePaymentIdByEstimateAmount(payload);
      // Note: this one appears to be the only one in use.
      case PublicPostActions.SEND_ESTIMATE_ADD_TO_CALENDAR_FROM_ZAPPIER:
        return sendEstimateAndAddToCalendarFromZappier(payload);

      case PublicPostActions.SEND_ESTIMATE_ADD_TO_CALENDAR_FROM_ZAPPIER_V1:
        var cleaned = transformZapierPayload(payload, {
          yesNoFields: ["schedule_call"] // optional
        });
        return sendEstimateAndAddToCalendarFromZappier_v1(cleaned);

      case PublicPostActions.SEND_PAID_BOOKING_INV_SEND_EMAIL:
        return sendPaidBookingInvoiceAndEmail(payload);

      case PublicPostActions.CLIENT_PAID_BOOKING_DEPOSIT:
        return handleBookingDepositPaid(payload);

      case PublicPostActions.STRIPE_GENERATE_BOOKING_DEPOSIT_INVOICE:
        return STRIPE_GENERATE_BOOKING_DEPOSIT_INVOICE(payload);

      case PublicPostActions.SINGLE_RESP_SEND_EMAIL_BOOKING_DEPOSIT:
        return jsonOut(SINGLE_RESP_SEND_EMAIL_BOOKING_DEPOSIT(payload));

      case PublicPostActions.CONTROLLER_GENERATE_STRIPE_INVOICE_AND_SEND_BOOKING_DEPOSIT_EMAIL:
        return CONTROLLER_GENERATE_STRIPE_INVOICE_AND_SEND_BOOKING_DEPOSIT_EMAIL(payload);

      case PublicPostActions.SEND_FINAL_INVOICE:
        return sendFinalInvoice(payload);
   
      case PublicPostActions.CLIENT_PAID_FINAL_INVOICE:
        return handleFinalInvoicePaid(payload);

      case PublicPostActions.TEST_PING:
        return jsonOut({ ok: true, pong: true });

      default:
        return jsonOut({ ok: false, error: "Unknown action", action: action });
    }
  } catch (err) {
    // 4) Guaranteed JSON error response
    return jsonOut({
      ok: false,
      error: (err && err.message) ? err.message : String(err),
      action: action,
      stack: err && err.stack ? err.stack : undefined
    });
  }
}


function jsonOut(obj, pretty) {
  // Pass-through if it's already a TextOutput
  if (obj && typeof obj.getContent === 'function' && typeof obj.setMimeType === 'function') {
    return obj;
  }
  var body = pretty ? JSON.stringify(obj, null, 2) : JSON.stringify(obj);
  return ContentService.createTextOutput(body).setMimeType(ContentService.MimeType.JSON);
}

/** Получить ключ Stripe: тестовый или продовый */
function getStripeKey_(useTest, overrideKey) {
  if (overrideKey) return overrideKey; // вручную
  var props = PropertiesService.getScriptProperties();
  // Заведи в Script Properties два ключа:
  // STRIPE_SECRET_KEY (prod) и STRIPE_SECRET_KEY_TEST (test)
  return useTest
    ? props.getProperty('STRIPE_SECRET_KEY_TEST')
    : props.getProperty('STRIPE_SECRET_KEY');
}

/**
 * Find a template by a single Gmail label.
 * @param {string} labelName  e.g. "email-templates-request-deposit"
 * @returns {{subject:string, html:string}}
 */
function getEmailTemplateByLabel_(labelName) {
  var query = 'in:anywhere label:"' + labelName + '"';
  Logger.log('🔍 Gmail search query: ' + query);

  var threads = GmailApp.search(query, 0, 5);
  Logger.log('📌 Found ' + threads.length + ' threads for label: ' + labelName);

  if (!threads.length) {
    throw new Error('Template not found for label: ' + labelName);
  }

  var msg = threads[0].getMessages()[0]; // take first message
  return { subject: msg.getSubject(), html: msg.getBody() };
}

/**
 * Sends booking-deposit email using a Gmail template under one label.
 * Placeholders supported: {firstName}, {lastNameAbv}, {stripe_deposit_link}, {deposit_amount}, plus any extraVars.
 */


function SINGLE_RESP_SEND_EMAIL_BOOKING_DEPOSIT(opt) {
  opt = opt || {};
  var toEmail = String(opt.toEmail || '').trim();
  if (!toEmail) throw new Error('toEmail is required');
  if (!opt.firstName) throw new Error('firstName is required');
  if (!opt.stripeDepositLink) throw new Error('stripeDepositLink is required');

  // 1) Label (single)
  var label = (opt.label && String(opt.label).trim()) || 'email-templates-request-deposit';
  Logger.log('🔎 Looking for template with label: ' + label);
  var tpl = getEmailTemplateByLabel_(label);

  // 2) Link text + optional button style
  var linkUrl  = String(opt.stripeDepositLink);
  var linkText = opt.linkText || 'Pay Your Deposit Securely with Stripe';
  var linkHtml = '<a href="' + linkUrl + '" target="_blank">' + linkText + '</a>';

  // Optional: button-styled version you can place with {pay-deposit-button}
  var buttonHtml =
    '<a href="' + linkUrl + '" target="_blank" ' +
    'style="background:#635bff;color:#fff;padding:10px 18px;border-radius:6px;' +
    'text-decoration:none;display:inline-block;font-weight:bold;font-family:Arial,sans-serif;">' +
    linkText + '</a>';

  // 3) Vars
  var vars = {
    firstName: opt.firstName,
    lastNameAbv: (opt.lastName || '').trim().charAt(0),
    stripe_deposit_link: linkUrl,          // keeps backward compatibility
    'pay-deposit-via-stripe': linkHtml,    // ← new: text + link
    'pay-deposit-button': buttonHtml,      // ← optional: button CTA
    deposit_amount: opt.depositAmount ? ('$' + opt.depositAmount) : ''
  };

  // Merge extra vars
  if (opt.extraVars && typeof opt.extraVars === 'object') {
    Object.keys(opt.extraVars).forEach(function (k) { vars[k] = opt.extraVars[k]; });
  }

  Logger.log('📨 Vars merged: ' + JSON.stringify(vars));

  // 4) Render
  var subj = renderTemplate_(tpl.subject, vars);
  var html = renderTemplate_(tpl.html, vars);
  var plain = subj + '\n\n' + linkText + ': ' + linkUrl + '\n'; // nicer plaintext fallback

  Logger.log('✉️ Sending to ' + toEmail + ' with subject: "' + subj + '"');
  GmailApp.sendEmail(toEmail, subj, plain, { 
  htmlBody: html,
  bcc: 'qa-booking-deposit@stlpartyhelpers.com'
});


  return { ok: true, sentTo: toEmail, subject: subj, labelUsed: label, depositAmount: vars.deposit_amount };
}

/** Safely unwrap ContentService TextOutput -> JSON object */
function unwrapJsonResponse_(res) {
  try {
    // If it's already an object, just return
    if (res && typeof res === 'object' && !res.getContent) return res;

    // Apps Script TextOutput has getContent(); fall back to String(res)
    var body = (res && typeof res.getContent === 'function') ? res.getContent() : String(res || '');
    return JSON.parse(body);
  } catch (e) {
    throw new Error('Failed to parse JSON from nested response: ' + (e && e.message ? e.message : e));
  }
}

function CONTROLLER_GENERATE_STRIPE_INVOICE_AND_SEND_BOOKING_DEPOSIT_EMAIL(payload) {
  var t0 = Date.now();
  var traceId = Utilities.getUuid();
  var steps = [];
  var debug = !!(payload && payload.debug);

  try {
    payload = payload || {};
    steps.push({ name: 'init', t: Date.now() });

    // Normalize inputs
    var customerEmail = String((payload.customerEmail || payload.email || payload.toEmail || '')).trim();
    var first = payload.first_name || payload.firstName || '';
    var last  = payload.last_name  || payload.lastName  || '';
    if (!customerEmail) throw new Error('customerEmail (or email/toEmail) is required');
    if (!first) throw new Error('first_name is required');

    // Resolve depositValue
    var depositValue = Number(payload.depositValue || 0);
    var depositPickedBy = 'given';

    if (!(isFinite(depositValue) && depositValue > 0)) {
      var est = Number(payload.estimatedTotal || payload.estimate || 0);
      if (isFinite(est) && est > 0) {
        var pick = findBookingDeposit(est);
        depositValue = pick.value;
        depositPickedBy = pick.pickedBy;  // 'estimate' or 'fallback'
        // propagate for the SR below
        payload.depositValue = depositValue;
        steps.push({ name: 'auto_pick_deposit', value: depositValue, pickedBy: depositPickedBy, estimate: est });
      } else {
        throw new Error('Provide depositValue OR estimatedTotal to compute a deposit.');
      }
    } else {
      steps.push({ name: 'use_given_deposit', value: depositValue });
      payload.depositValue = depositValue; // ensure downstream sees it
    }

    // === 1) Generate Stripe invoice (SR) ===
    var s1 = Date.now();
    var invoiceTextOut = STRIPE_GENERATE_BOOKING_DEPOSIT_INVOICE(Object.assign({}, payload, {
      email: customerEmail,
      useTest: (payload.useTest !== false) // default: test mode
    }));
    var invoiceJson = unwrapJsonResponse_(invoiceTextOut);
    steps.push({ name: 'STRIPE_GENERATE_BOOKING_DEPOSIT_INVOICE', ms: Date.now() - s1, ok: !!invoiceJson && invoiceJson.ok });

    if (!invoiceJson || invoiceJson.ok === false) {
      throw new Error('Invoice creation failed: ' + (invoiceJson && invoiceJson.error ? invoiceJson.error : 'unknown'));
    }
    if (!invoiceJson.hosted_invoice_url) throw new Error('Invoice created but no hosted_invoice_url returned.');

    // === 2) Send booking-deposit email ===
    var s2 = Date.now();
    var sendRes = SINGLE_RESP_SEND_EMAIL_BOOKING_DEPOSIT({
      toEmail: customerEmail,
      firstName: first,
      lastName: last,
      stripeDepositLink: invoiceJson.hosted_invoice_url,
      label: payload.label || 'email-templates-request-deposit',
      depositAmount: depositValue,
      extraVars: {
        eventDateTime:  payload.eventDateTimeLocal || '',
        helpersCount:   payload.helpersCount != null ? String(payload.helpersCount) : '',
        hours:          payload.hours != null ? String(payload.hours) : '',
        estimatedTotal: payload.estimatedTotal != null ? String(payload.estimatedTotal) : ''
      },
      qaEmail: payload.qaEmail
    });
    steps.push({ name: 'SINGLE_RESP_SEND_EMAIL_BOOKING_DEPOSIT', ms: Date.now() - s2, ok: !!sendRes && sendRes.ok });

    // Response
    var out = {
      ok: true,
      traceId: traceId,
      elapsed_ms: Date.now() - t0,
      invoice: {
        id:                 invoiceJson.invoiceId,
        status:             invoiceJson.status,
        hosted_invoice_url: invoiceJson.hosted_invoice_url,
        invoice_pdf:        invoiceJson.invoice_pdf,
        amount_due:         invoiceJson.amount_due,
        currency:           invoiceJson.currency,
        due_date:           invoiceJson.due_date
      },
      email: sendRes,
      deposit: { value: depositValue, pickedBy: depositPickedBy }
    };

    if (debug) {
      out.debug = { inputs: payload, steps: steps, raw_invoice: invoiceJson };
    }
    return jsonOut(out, debug);

  } catch (err) {
    return jsonOut({
      ok: false,
      traceId: traceId,
      elapsed_ms: Date.now() - t0,
      error: (err && err.message) ? err.message : String(err),
      where: 'CONTROLLER_GENERATE_STRIPE_INVOICE_AND_SEND_BOOKING_DEPOSIT_EMAIL',
      steps: steps
    }, !!(payload && payload.debug));
  }
}



/** Simple HTML placeholder replacer: {key} or {{key}} */
function renderTemplate_(text, vars) {
  if (!text) return '';
  return Object.keys(vars || {}).reduce(function(out, k) {
    var v = (vars[k] == null) ? '' : String(vars[k]);
    // {key}
    out = out.replace(new RegExp('\\{' + k + '\\}', 'g'), v);
    // {{key}}
    out = out.replace(new RegExp('\\{\\{' + k + '\\}\\}', 'g'), v);
    return out;
  }, String(text));
}

/** Build a Gmail search query from labels. */
function buildLabelQuery_(labels, opts) {
  opts = opts || {};
  var base = 'in:anywhere ' + (opts.onlyDraft ? 'has:draft ' : '');
  if (!Array.isArray(labels) || labels.length === 0) return base.trim();
  var parts = labels.map(function(l) {
    l = String(l || '').trim();
    if (!l) return '';
    // allow users to pass either raw name or already quoted
    return 'label:"' + l.replace(/"/g, '\\"') + '"';
  }).filter(Boolean);
  return (base + parts.join(' ')).trim();
}

/**
 * Find the first message that matches labels.
 * @param {string[]} labels  e.g. ["email-templates/automation","email-templates/request-deposit"]
 * @param {{onlyDraft?:boolean, maxThreads?:number}} [opts]
 * @returns {{subject:string, html:string}}
 */
function getEmailTemplateByLabels_(labels, opts) {
  opts = opts || {};
  var q = buildLabelQuery_(labels, { onlyDraft: true }); // drafts by default
  var threads = GmailApp.search(q, 0, opts.maxThreads || 10);

  // Prefer the newest draft in the newest thread
  for (var i = 0; i < threads.length; i++) {
    var msgs = threads[i].getMessages();
    for (var m = msgs.length - 1; m >= 0; m--) {
      var msg = msgs[m];
      // Some accounts return drafts via isDraft(); some store templates as draft-like. Check both.
      if (msg.isDraft && msg.isDraft()) {
        return { subject: msg.getSubject(), html: msg.getBody() };
      }
    }
  }

  // If no explicit drafts found, allow any message under those labels as a fallback
  if (!threads.length) {
    // try without has:draft
    q = buildLabelQuery_(labels, { onlyDraft: false });
    threads = GmailApp.search(q, 0, opts.maxThreads || 10);
  }
  if (threads.length) {
    var last = threads[0].getMessages().slice(-1)[0];
    return { subject: last.getSubject(), html: last.getBody() };
  }

  throw new Error('Template not found. Labels tried: [' + (labels || []).join(', ') + '] Query="' + q + '"');
}


/** tiny utility: FirstName + first letter of LastName (with no dot or with dot as you prefer) */
function firstLetter_(s) { return (s || '').trim().charAt(0) || ''; }

function normalizeEmailName_(payload) {
  var rawEmail = payload && payload.email;
  var rawName  = payload && payload.name;

  // Если пришёл объект (например {email:"...", name:"..."})
  if (rawEmail && typeof rawEmail === 'object') {
    if (rawEmail.email) rawEmail = rawEmail.email;
    if (!rawName && rawEmail.name) rawName = rawEmail.name; // иногда имя внутри того же объекта
  }
  if (rawName && typeof rawName === 'object' && rawName.name) {
    rawName = rawName.name;
  }

  // В строки
  var email = rawEmail ? String(rawEmail).trim() : '';
  var name  = rawName  ? String(rawName).trim()  : '';

  return { email: email, name: name };
}

function STRIPE_GENERATE_BOOKING_DEPOSIT_INVOICE(payload) {
  try {
    payload = payload || {};
    var idn   = normalizeEmailName_(payload);
    var email = idn.email || 'saint.louis.mail@gmail.com';
    var name  = idn.name  || 'Client';

    // depositValue or derive from estimate/estimatedTotal
    var value = Number(payload.depositValue || 0);
    var pickedBy = 'given';

    if (!(isFinite(value) && value > 0)) {
      var est = Number(payload.estimatedTotal || payload.estimate || 0);
      if (!(isFinite(est) && est > 0)) {
        throw new Error('Provide depositValue OR estimatedTotal to compute a deposit.');
      }
      var pick = findBookingDeposit(est);
      value    = pick.value;
      pickedBy = pick.pickedBy; // 'estimate' or 'fallback'
    }

    var memo  = (payload.memo) || 'Booking deposit invoice';

    var price = BOOKING_DEPOSITS.find(function (p) { return p.value === value; });
    if (!price) throw new Error('Unknown deposit value (no matching price): ' + value);
    Logger.log('Using price: ' + JSON.stringify(price));

    var customFields = sanitizeCustomFields_([
      { name: 'Event Date & Time', value: payload.eventDateTimeLocal },
      { name: 'Helpers Count',     value: payload.helpersCount },
      { name: 'Hours',             value: payload.hours },
      { name: 'Estimated Total',   value: payload.estimatedTotal }
    ]);

    var metadata = sanitizeMetadata_({
      source: 'apps-script',
      kind:   'booking_deposit',
      pickedBy: pickedBy,
      eventDateTimeLocal: payload.eventDateTimeLocal
    });

    var res = createStripeInvoice_({
      useTest: (payload.useTest !== false),  // default TEST
      email: email,
      name:  name,
      priceId: price.id,
      quantity: 1,
      collectionMethod: 'send_invoice',
      daysUntilDue: payload.daysUntilDue != null ? payload.daysUntilDue : 7,
      descriptionInvoice: memo,
      footer:
        'This payment serves as a booking deposit for your upcoming event.\n' +
        'It confirms the reservation of helper(s) for the scheduled date.\n\n' +
        'The remaining balance (total event cost minus this deposit) will be invoiced ' +
        'separately on the next business day following the event, with a 7-day due date for payment.\n\n' +
        'Please retain this invoice as confirmation of your booking. We appreciate your business!',
      customFields: customFields,
      metadata: metadata,
      idempotencyKey: ['deposit', price.id, email, (payload.eventDateTimeLocal || '')].join(':')
    });

    Logger.log('Invoice created: ' + JSON.stringify(res, null, 2));
    // include how the deposit was chosen
    res.deposit = { value: value, pickedBy: pickedBy };
    return ContentService.createTextOutput(JSON.stringify(res, null, 2))
      .setMimeType(ContentService.MimeType.JSON);

  } catch (err) {
    Logger.log('ERROR: ' + err.stack);
    return ContentService.createTextOutput(JSON.stringify({ ok:false, error: err.message }, null, 2))
      .setMimeType(ContentService.MimeType.JSON);
  }
}


function createStripeInvoice_(opt) {
  opt = opt || {};
  var apiKey = getStripeKey_(opt.useTest !== false, opt.apiKey);
  if (!apiKey) throw new Error('Stripe API key not set. Put STRIPE_SECRET_KEY_TEST / STRIPE_SECRET_KEY in Script Properties.');

  // 1) Customer
  var customerId = getOrCreateStripeCustomer_(apiKey, {
    customerId: opt.customerId,
    email: opt.email,
    name:  opt.name
  });
  if (!customerId) throw new Error('No customerId could be resolved');

  // 🔹 NEW: make sure we don't accumulate old pending lines
  clearPendingInvoiceItems_(apiKey, customerId);

  // 2) Exactly ONE line item (priceId OR amount+currency)
  if (opt.priceId) {
    var iiByPrice = addInvoiceItemByPrice_(apiKey, customerId, opt.priceId, opt.quantity || 1);
    Logger.log('InvoiceItem by price created: ' + JSON.stringify(iiByPrice));
  } else if (isFinite(opt.amount) && opt.currency) {
    var iiByAmt = addInvoiceItemByAmount_(apiKey, customerId, opt.amount, opt.currency, opt.description || '');
    Logger.log('InvoiceItem by amount created: ' + JSON.stringify(iiByAmt));
  } else {
    throw new Error('Provide either priceId OR amount+currency for line item.');
  }

  // 3) Sanitize metadata/custom_fields
  var md = {};
  if (opt.metadata) {
    Object.keys(opt.metadata).forEach(function(k){
      var v = opt.metadata[k];
      if (v == null) return;
      var sv = String(v).trim();
      if (sv !== '') md[k] = sv;
    });
  }
  var cf = [];
  if (Array.isArray(opt.customFields)) {
    opt.customFields.forEach(function(item){
      if (!item) return;
      var name  = String(item.name  || '').trim();
      var value = String(item.value || '').trim();
      if (name && value) cf.push({ name: name, value: value });
    });
  }

  // 4) Create invoice (draft) — include the one we just added
  var invoicePayload = {
    customer: customerId,
    collection_method: opt.collectionMethod || 'send_invoice',
    auto_advance: true,
    pending_invoice_items_behavior: 'include', // ← include pending items (only the fresh one now)
    expand: ['lines.data.price', 'customer']
  };
  // Memo + footer
  if (opt.descriptionInvoice || opt.memo) invoicePayload.description = opt.descriptionInvoice || opt.memo; // "Memo"
  if (opt.footer)                           invoicePayload.footer      = opt.footer;

  if (invoicePayload.collection_method === 'send_invoice') {
    invoicePayload.days_until_due = (opt.daysUntilDue != null) ? opt.daysUntilDue : 7;
  }
  if (Object.keys(md).length) invoicePayload.metadata      = md;
  if (cf.length)              invoicePayload.custom_fields = cf;

  var invoiceDraft = stripeRequest_('POST', '/v1/invoices', invoicePayload, apiKey, opt.idempotencyKey);
  Logger.log('Draft invoice: ' + JSON.stringify(invoiceDraft));

  // 5) Finalize (no email)
  var finalized = stripeRequest_(
    'POST',
    '/v1/invoices/' + encodeURIComponent(invoiceDraft.id) + '/finalize',
    {},
    apiKey
  );
  Logger.log('Finalized invoice: ' + JSON.stringify(finalized));

  if (!finalized || !finalized.id || !finalized.hosted_invoice_url) {
    throw new Error('Invoice finalization did not return expected fields: ' + JSON.stringify(finalized));
  }

  return {
    ok: true,
    invoiceId: finalized.id,
    status: finalized.status,
    hosted_invoice_url: finalized.hosted_invoice_url,
    invoice_pdf: finalized.invoice_pdf,
    amount_due: finalized.amount_due,
    currency: finalized.currency,
    due_date: finalized.due_date,
    customer: (typeof finalized.customer === 'object' && finalized.customer) ? finalized.customer.id : finalized.customer,
    lines_count: finalized.lines && finalized.lines.data ? finalized.lines.data.length : 0
  };
}

/** Delete all PENDING invoice items (invoice == null) for this customer. */
function clearPendingInvoiceItems_(apiKey, customerId) {
  // list pending invoice items
  var list = stripeRequest_(
    'GET',
    '/v1/invoiceitems',
    { customer: customerId, limit: 100 },
    apiKey
  );

  if (!list || !list.data) return 0;

  var deleted = 0;
  list.data.forEach(function(ii) {
    if (!ii.invoice) { // pending (not attached to any invoice yet)
      try {
        stripeRequest_('DELETE', '/v1/invoiceitems/' + encodeURIComponent(ii.id), {}, apiKey);
        deleted++;
      } catch (e) {
        Logger.log('Failed to delete invoiceitem ' + ii.id + ': ' + e.message);
      }
    }
  });
  Logger.log('Cleared ' + deleted + ' pending invoice items for ' + customerId);
  return deleted;
}

function sanitizeCustomFields_(arr) {
  // Удаляем элементы без name или с пустым value, приводим к строкам
  return (arr || []).reduce(function(out, item) {
    if (!item || item.value == null) return out;
    var name  = String(item.name || '').trim();
    var value = String(item.value).trim();
    if (!name || !value) return out; // пустые — выкидываем
    out.push({ name: name, value: value });
    return out;
  }, []);
}

function debugStripeProps() {
  var p = PropertiesService.getScriptProperties().getProperties();
  Logger.log(JSON.stringify(p, null, 2));
}


function sanitizeMetadata_(obj) {
  // Превращаем значения в строки и убираем пустые
  var out = {};
  if (!obj) return out;
  Object.keys(obj).forEach(function(k) {
    var v = obj[k];
    if (v == null) return;
    var sv = String(v).trim();
    if (sv !== '') out[k] = sv;
  });
  return out;
}

/** Добавить invoice item по Price ID */
function addInvoiceItemByPrice_(apiKey, customerId, priceId, quantity) {
  return stripeRequest_('POST', '/v1/invoiceitems', {
    customer: customerId,
    price: priceId,
    quantity: quantity || 1
  }, apiKey);
}

/** Добавить invoice item по сумме (в минимальных единицах) */
function addInvoiceItemByAmount_(apiKey, customerId, amount, currency, description) {
  return stripeRequest_('POST', '/v1/invoiceitems', {
    customer: customerId,
    amount: amount,
    currency: currency,
    description: description || ''
  }, apiKey);
}

function getOrCreateStripeCustomer_(apiKey, opts) {
  var email = opts && opts.email;
  var name  = opts && opts.name;

  if (email && typeof email === 'object') {
    if (email.email) email = email.email;
    if (!name && email.name) name = email.name;
  }
  if (name && typeof name === 'object' && name.name) {
    name = name.name;
  }
  email = email ? String(email).trim() : '';
  name  = name  ? String(name).trim()  : '';

  if (opts.customerId) return opts.customerId;
  if (!email && !name) throw new Error('Need at least email or name to create/find customer.');

  var list = stripeRequest_('GET', '/v1/customers', { email: email || undefined, limit: 1 }, apiKey);
  if (list && list.data && list.data.length) return list.data[0].id;

  var created = stripeRequest_('POST', '/v1/customers', { email: email || undefined, name: name || undefined }, apiKey);
  return created.id;
}

/** Add invoice item by Price ID (price_...). */
function createStripeInvoiceItemByPrice_(apiKey, customerId, priceId, quantity) {
  var body = {
    customer: customerId,
    price: priceId,
    quantity: quantity || 1
  };
  return stripeRequest_('post', '/v1/invoiceitems', body, apiKey);
}

/** Add invoice item by raw amount (in minor units) + currency. */
function createStripeInvoiceItemByAmount_(apiKey, customerId, amount, currency, description) {
  var body = {
    customer: customerId,
    amount: amount,       // cents
    currency: currency,   // "usd"
    description: description || ''
  };
  return stripeRequest_('post', '/v1/invoiceitems', body, apiKey);
}

function demo_CreateInvoiceByAmount() {
  var res = createStripeInvoice_({
    useTest: true,
    email: 'client@example.com',
    name: 'Client Name',
    amount: 28000,              // $280.00
    currency: 'usd',
    description: 'Event setup (2 helpers × 4h)',
    collectionMethod: 'send_invoice',
    daysUntilDue: 7,
    // Memo (aka description on the invoice)
  memo: 'Have billing questions? Call us at 314-350-4400',

  customFields: [
    { name: 'Event Date & Time', value: '2025-09-18 10:00 PM' },
    { name: 'Helpers Count',     value: '1' },
    { name: 'Hours',             value: '5' },
    { name: 'Estimated Total',   value: '300' }
  ],
  metadata: { source: 'apps-script', kind: 'booking_deposit' },
  // or descriptionInvoice: '...'

  // Footer text (multi-line supported)
  footer:
    'This payment serves as a booking deposit...\n' +
    'Please retain this invoice as confirmation of your booking.',
    send: true
  });
  Logger.log(JSON.stringify(res, null, 2));
}

/** Generic Stripe request via UrlFetchApp. */
function stripeRequest_(method, path, params, apiKey) {
  var url = 'https://api.stripe.com' + path;

  var options = {
    method: method.toUpperCase(),
    headers: {
      Authorization: 'Bearer ' + apiKey
    },
    muteHttpExceptions: true
  };

  if (options.method === 'GET') {
    // Encode params to querystring
    if (params && Object.keys(params).length) {
      url += '?' + toFormUrlEncoded_(params);
    }
  } else {
    // Stripe expects form-encoded body
    options.contentType = 'application/x-www-form-urlencoded';
    options.payload = toFormUrlEncoded_(params || {});
  }

  var resp = UrlFetchApp.fetch(url, options);
  var code = resp.getResponseCode();
  var text = resp.getContentText();
  if (code < 200 || code >= 300) {
    throw new Error('Stripe API error ' + code + ': ' + text);
  }
  return JSON.parse(text);
}

/** Convert nested object to application/x-www-form-urlencoded (Stripe-friendly). */
function toFormUrlEncoded_(obj, prefix) {
  var str = [];
  for (var p in obj) {
    if (!obj.hasOwnProperty(p)) continue;
    var k = prefix ? prefix + '[' + p + ']' : p;
    var v = obj[p];
    if (v === null || v === undefined) continue;
    if (typeof v === 'object' && !Array.isArray(v)) {
      str.push(toFormUrlEncoded_(v, k));
    } else if (Array.isArray(v)) {
      v.forEach(function (item, idx) {
        if (typeof item === 'object') {
          str.push(toFormUrlEncoded_(item, k + '[' + idx + ']'));
        } else {
          str.push(encodeURIComponent(k + '[]') + '=' + encodeURIComponent(item));
        }
      });
    } else {
      str.push(encodeURIComponent(k) + '=' + encodeURIComponent(v));
    }
  }
  return str.join('&');
}

/*
function test_sendEstimateAndAddToCalendarFromZappier() {
  const payload = {
    action: "SEND_ESTIMATE_ADD_TO_CALENDAR_FROM_ZAPPIER",
    first_name: "Alex",
    last_name: "Shevelyov",
    email_address: "stlph-crm-test@shevelyov.com",
    phone_number: "3145555555",
    event_date: "July 10, 2025 4:00 PM",
    event_time: "4:00 PM",
    event_location: "2300 Hitzert Ct, Fenton, MO 63026, USA",
    guests_expected: "300 - 400 Guests",
    helpers_requested: "I Need 2 Helpers",
    for_how_many_hours: "for 5 Hours",
    occasion: "Christmas"
  };

  const result = sendEstimateAndAddToCalendarFromZappier(payload);
  Logger.log(result.getContent());
}
*/

// ==== 🔁 PUBLIC POST ACTION HANDLERS ====
// These functions are invoked via `doPost` based on `payload.action`.
// Each handles a single externally-triggered business process (e.g., WPForms, Monday.com, Stripe webhooks).

// Our first communication with clients
function getStripePaymentIdByEstimateAmount(payload) {
  try {
    
    // Use the passed-in payload directly (already parsed in doPost)
    findBookingDeposit(payload); // Optionally use payload to pass data
    return ContentService.createTextOutput("Success - sendEstimateAndAddToCalendar");
  } catch (err) {
    return ContentService.createTextOutput("Error: " + err.message);
  }
}

// Our first communication with clients
function sendEstimateAndAddToCalendar(payload) {
  try {
    // Use the passed-in payload directly (already parsed in doPost)
    processNewLead(payload); // Optionally use payload to pass data
    return ContentService.createTextOutput("Success - sendEstimateAndAddToCalendar");
  } catch (err) {
    return ContentService.createTextOutput("Error: " + err.message);
  }
}

function sendEstimateAndAddToCalendarFromZappier(payload) {
  try {
    // Capture and return the actual result
    return processNewLeadFromZapier(payload);
  } catch (err) {
    return ContentService.createTextOutput(
      JSON.stringify({ success: false, error: err.message })
    ).setMimeType(ContentService.MimeType.JSON);
  }
}

function sendEstimateAndAddToCalendarFromZappier_v1(cleanedPayload) {
  try {
    return processNewLeadFromZappier_v1(cleanedPayload);
  } catch (err) {
    return ContentService.createTextOutput(
      JSON.stringify({ success: false, error: err.message })
    ).setMimeType(ContentService.MimeType.JSON);
  }
}

function sendPaidBookingInvoiceAndEmail(payload) {
  try {
    // Use the passed-in payload directly (already parsed in doPost)
    processSendingPaidBookingInvoiceAndEmail(payload); // Optionally use payload to pass data
    return ContentService.createTextOutput("Success - sendPaidBookingInvoiceAndEmail");
  } catch (err) {
    return ContentService.createTextOutput("Error: " + err.message);
  }
}

/*
// DELETE?
function sendTestToZapier() {
  const url = 'https://hooks.zapier.com/hooks/catch/21931276/uo7w4yw/ '; // replace this

  const payload = {
    itemId: 123456789,
    columnId: "status",
    newValue: "Confirmed"
  };

  const options = {
    method: "post",
    contentType: "application/json",
    payload: JSON.stringify(payload)
  };

  UrlFetchApp.fetch(url, options);
}
*/

function processSendingPaidBookingInvoiceAndEmail(payload) {
  const invoiceUrl = createStripeInvoice(payload);
  const htmlBody = generateBookingInvoiceEmail(payload, invoiceUrl);

  /* Payload Params */
  // original invoice amount 

  MailApp.sendEmail({
    to: payload.emailAddress,
    subject: "Booking Deposit Invoice – STL Party Helpers",
    htmlBody: htmlBody
  });
}

function createStripeInvoice(payload) {
  const STRIPE_SECRET_KEY = 'sk_test_...' // Replace with your actual secret key
  const PRODUCT_LOOKUP_KEY = 'booking_deposit'; // Must match product set in Stripe
  const INVOICE_TEMPLATE_ID = ''; // Optional — if you use a saved invoice template

  const headers = {
    Authorization: 'Bearer ' + STRIPE_SECRET_KEY
  };

  try {
    // Step 1: Create or find customer
    const customerRes = UrlFetchApp.fetch('https://api.stripe.com/v1/customers', {
      method: 'post',
      headers,
      payload: {
        email: payload.emailAddress,
        name: payload.name
      }
    });

    const customer = JSON.parse(customerRes.getContentText());

    // Step 2: Create invoice item using price via lookup_key
    const invoiceItemRes = UrlFetchApp.fetch('https://api.stripe.com/v1/invoiceitems', {
      method: 'post',
      headers,
      payload: {
        customer: customer.id,
        price_data: {
          currency: 'usd',
          product_data: {
            name: "Booking Deposit"
          },
          unit_amount: 15000 // $150.00 in cents
        }
      }
    });

    // Step 3: Create invoice
    const invoiceRes = UrlFetchApp.fetch('https://api.stripe.com/v1/invoices', {
      method: 'post',
      headers,
      payload: {
        customer: customer.id,
        auto_advance: true,
        collection_method: 'send_invoice',
        days_until_due: 3
      }
    });

    const invoice = JSON.parse(invoiceRes.getContentText());
    return invoice.hosted_invoice_url;

  } catch (err) {
    throw new Error("Stripe error: " + err.message);
  }
}

function handleEstimateFromZapierJson(payload) {
  try {
    const {
      first_name,
      last_name,
      email_address,
      phone_number,
      event_date,
      event_time,
      event_location,
      guests_expected,
      helpers_requested,
      for_how_many_hours,
      occasion
    } = payload;

    const clientName = `${first_name} ${last_name}`;
    const guestCount = parseInt(guests_expected || 0, 10);
    const numHelpers = parseInt(helpers_requested || 0, 10);
    const duration = parseFloat(for_how_many_hours || 0);

    const baseRate = 250;
    const overtimeRate = 45;
    const baseHours = 4;

    const totalBase = numHelpers * baseRate;
    const extraHours = Math.max(0, duration - baseHours);
    const overtimeCost = numHelpers * overtimeRate * extraHours;
    const totalCost = totalBase + overtimeCost;

    const subject = "Your Event Quote from STL Party Helpers";
    const referenceNumber = new Date().getTime();
    const quoteHTML = generateQuoteEmail(
      clientName,
      event_date,
      event_time,
      event_location,
      numHelpers,
      duration,
      totalCost,
      occasion,
      guestCount
    );

    return ContentService.createTextOutput(JSON.stringify({
      success: true,
      estimate: totalCost,
      quoteHTML,
      referenceNumber,
      breakdown: {
        clientName,
        guestCount,
        numHelpers,
        duration,
        totalBase,
        overtimeCost,
        event_date,
        event_time,
        event_location
      }
    })).setMimeType(ContentService.MimeType.JSON);

  } catch (err) {
    return ContentService.createTextOutput(JSON.stringify({
      success: false,
      error: err.message
    })).setMimeType(ContentService.MimeType.JSON);
  }
} // END

// Old way - does everything - will retire soon
function processNewLead() {
  var folderId = "1qH-4Vq6aLHhnM7jpUD3SbMLX8SW89oDl";
  var folder = DriveApp.getFolderById(folderId);
  var ss = getOrCreateSpreadsheet(folder);
  var sheet = ss.getSheetByName("Leads") || ss.insertSheet("Leads");

  console.log("📌 Fetching unprocessed leads...");
  var newLeads = getUnprocessedLeads();

  if (newLeads.length === 0) {
    console.log("✅ No new leads to process.");
    return;
  }

  newLeads.forEach((lead) => {
    try {
      var { subject, body, sender, thread } = lead;

      console.log(`📩 Processing lead: ${subject} from ${sender}`);
      var parsedData = parseLeadData(body);

      // ✅ Validate parsed data before proceeding
      if (!validateLeadData(parsedData, parsedData.clientName)) {
        console.error(
          `🚨 Missing required fields. Marking lead as "auto-quote-sending-failed".`
        );
        markLeadAsFailed(thread, "auto-quote-sending-failed");
        return;
      }

      // ✅ Process the lead
      processLead(parsedData, ss, thread);
    } catch (error) {
      console.error(`❌ ERROR processing lead: ${error.message}`);
      markLeadAsFailed(lead.thread, "auto-quote-sending-failed");
    }
  });
}


function processNewLeadFromZapier(parsedData) {
  if (!parsedData) {
    console.error("❌ ERROR: No parsed data available.");
    return ContentService.createTextOutput("Error: No parsed data");
  }

  var {
    first_name,
    last_name,
    event_date,
    event_time,
    phone_number,
    event_location,
    helpers_requested,
    for_how_many_hours,
    email_address,
    occasion,
    guests_expected,
    dryRun
  } = parsedData;

  const clientName = `${first_name} ${last_name}`;
  const guestCount = parseInt(guests_expected || 0, 10);
  const numHelpers = parseInt((helpers_requested.match(/\d+/) || [0])[0], 10);
  const duration = parseFloat((for_how_many_hours.match(/\d+(\.\d+)?/) || [0])[0]);
  const phone = phone_number;
  const location = event_location;
  const eventDate = event_date;
  const eventTime = event_time;
  const email = email_address;
  const dryRunFlag = dryRun;

  const { totalCost, baseCost, additionalCost, baseRate, hourlyRate, rateLabel } = calculateTotalCost_v2(numHelpers, duration, eventDate);
  const referenceNumber = generateShortQuoteID(email, eventDate);

  console.log(`📊 Pricing Breakdown: Base = $${baseCost}, Extra = $${additionalCost}, Total = $${totalCost}`);

  // Try calendar creation
  const eventIdResponse = calendareventcreator.createEvent({
    calendarId: ESTIMATE_SENT_CALENDAR_ID,
    clientName,
    occasion,
    guestCount,
    eventDate,
    eventTime,
    phone,
    location,
    numHelpers,
    duration,
    totalCost,
    emailId: email,
    threadId: "",
    dataSource: DataSource.ZAPIER
  });

  let calendarError = null;
  let eventId = null;

  if (!eventIdResponse || eventIdResponse.error) {
    calendarError = eventIdResponse?.error || "Unknown error";
    console.error(`❌ ERROR: Failed to create calendar event. Reason: ${calendarError}`);
  } else {
    eventId = eventIdResponse.eventId;
    console.log(`✅ Calendar event created: ${eventId}`);
  }


  // ✅ Always send the quote email
  try {
    // ✅ Send email quote
sendQuoteEmailOnly_v2(
   clientName,
   email,
   eventDate,
   eventTime,
   location,
   numHelpers,
   duration,
   baseRate,
   hourlyRate,
   totalCost,
   occasion,
   guestCount,
   rateLabel,
   dryRunFlag
  );
    console.log(`📧 Quote email sent to ${email}`);
  } catch (err) {
    console.error(`❌ ERROR: Failed to send quote email to ${email}. ${err.message}`);
    return ContentService.createTextOutput(
      JSON.stringify({
        success: false,
        error: "Quote email failed",
        reason: err.message
      })
    ).setMimeType(ContentService.MimeType.JSON);
  }

   const geoData = getLatLng(location); // Replace with your address field


  // ✅ Respond back with success + errors if any
  return ContentService.createTextOutput(
    JSON.stringify({
      referenceNumber,
      success: true,
      emailSent: true,
      lat: geoData?.lat || null,
      long: geoData?.lng || null,
      fullAddress: geoData?.fullAddress || null,
      eventId: eventId || null,
      estimate: totalCost,
      calendarCreated: !calendarError,
      calendarError: calendarError || null
    })
  ).setMimeType(ContentService.MimeType.JSON);
}


// Old way - does everything - will retire soon
function processNewLeadFromZappier_v1(cleanedData) {
  var folderId = "1qH-4Vq6aLHhnM7jpUD3SbMLX8SW89oDl";
  var folder = DriveApp.getFolderById(folderId);
  var ss = getOrCreateSpreadsheet(folder);
  var sheet = ss.getSheetByName("Leads") || ss.insertSheet("Leads");

  try {
    var { subject, sender, thread } = lead;

    console.log(`📩 Processing lead from Zappier: ${subject} from ${sender}`);


    // ✅ Validate parsed data before proceeding
    if (!validateLeadData(cleanedData, cleanedData.clientName)) {
      console.error(
        `🚨 Missing required fields. Auto quote for lead from Zappier not sent".`
      );

      return;
    }

    // ✅ Process the lead
    processLead(cleanedData, ss, thread);
  } catch (error) {
    console.error(`❌ ERROR processing lead from Zappier: ${error.message}`);

  }
}

function processNewLeadFromZapier_v1(clean) {
  // --- Guard: ensure the cleaned payload is present and has critical fields
  if (!clean) {
    return ContentService.createTextOutput(
      JSON.stringify({ success: false, error: "No data" })
    ).setMimeType(ContentService.MimeType.JSON);
  }

  const {
    clientName,
    email,
    phone,
    eventDate,
    eventTime,
    location,
    numHelpers,
    duration,
    guestCount,
    occasion
  } = clean;

  // Minimal validations (you can assume transform already validated—
  // these are just safety nets)
  if (!clientName)  return _jsonError("Missing clientName");
  if (!email)       return _jsonError("Missing email");
  if (!eventDate)   return _jsonError("Missing eventDate");
  if (!eventTime)   return _jsonError("Missing eventTime");

  // --- Pricing + reference
  const { totalCost, baseCost, additionalCost } = calculateTotalCost(numHelpers, duration);
  const referenceNumber = generateShortQuoteID(email, eventDate);

  console.log(`📊 Pricing Breakdown: Base = $${baseCost}, Extra = $${additionalCost}, Total = $${totalCost}`);

  // --- Calendar event
  const eventIdResponse = calendareventcreator.createEvent({
    calendarId: ESTIMATE_SENT_CALENDAR_ID,
    clientName,
    occasion,
    guestCount,
    eventDate,
    eventTime,
    phone,
    location,
    numHelpers,
    duration,
    totalCost,
    emailId: email,
    threadId: "",
    dataSource: DataSource.ZAPIER
  });

  let calendarError = null;
  let eventId = null;
  if (!eventIdResponse || eventIdResponse.error) {
    calendarError = eventIdResponse?.error || "Unknown error";
    console.error(`❌ ERROR: Failed to create calendar event. Reason: ${calendarError}`);
  } else {
    eventId = eventIdResponse.eventId;
    console.log(`✅ Calendar event created: ${eventId}`);
  }

  // --- Send quote email (always)
  try {
    sendQuoteEmailOnly(
      clientName,
      email,
      eventDate,
      eventTime,
      location,
      numHelpers,
      duration,
      totalCost,
      occasion,
      guestCount
    );
    console.log(`📧 Quote email sent to ${email}`);
  } catch (err) {
    console.error(`❌ ERROR: Failed to send quote email to ${email}. ${err.message}`);
    return ContentService.createTextOutput(
      JSON.stringify({ success: false, error: "Quote email failed", reason: err.message })
    ).setMimeType(ContentService.MimeType.JSON);
  }

  // --- Geocode (unchanged)
  const geoData = getLatLng(location);

  // --- Final JSON response
  return ContentService.createTextOutput(
    JSON.stringify({
      referenceNumber,
      success: true,
      emailSent: true,
      lat: geoData?.lat || null,
      long: geoData?.lng || null,
      fullAddress: geoData?.fullAddress || null,
      eventId: eventId || null,
      estimate: totalCost,
      calendarCreated: !calendarError,
      calendarError: calendarError || null
    })
  ).setMimeType(ContentService.MimeType.JSON);
}

/** small shared helper */
function _jsonError(msg) {
  return ContentService.createTextOutput(
    JSON.stringify({ success: false, error: msg })
  ).setMimeType(ContentService.MimeType.JSON);
}


/** "I Need 2 Helpers" -> 2 */
function parseHelpers(text) {
  const n = numberFromText(text);
  return Number.isFinite(n) ? n : 0;
}

/** "for 4 Hours (minimum)" -> 4 */
function parseDurationHours(text) {
  const n = numberFromText(text);
  return Number.isFinite(n) ? n : 0;
}

/** Lowercase includes "yes" -> "Checked"/"Unchecked" (if you need checkbox mapping) */
function yesToCheckbox(text) {
  const s = String(text || "").trim().toLowerCase();

  // explicit yes patterns
  const hasYes = /\b(yes|y|true|1|yeah|yep|sure|please|call\s*me)\b/.test(s);

  // explicit no patterns
  const hasNo = /\b(no|n|false|0|nope|dont|don't|do not|nah)\b/.test(s);

  if (hasYes && !hasNo) return "Checked";     // definitely yes
  if (hasNo && !hasYes) return "Unchecked";   // definitely no

  // fallback: simple contains "yes"
  return s.includes("yes") ? "Checked" : "Unchecked";
}



/***********************
 * Low-level helpers
 ***********************/
function _firstNumber(value) {
  const m = String(value || "").match(/-?\d+(\.\d+)?/);
  return m ? Number(m[0]) : 0;
}


function consolidateOther(main, alt) {
  const m = String(main || "");
  if (!m) return "Unspecified";
  return m.toLowerCase().includes("other")
    ? (alt && String(alt).trim() ? alt : "Unspecified")
    : m;
}

/***********************
 * Per-field transformers
 ***********************/
function getClientName(first_name, last_name) {
  return `${first_name || ""} ${last_name || ""}`.trim();
}

function getEmail(email_address) {
  return (email_address || "").trim();
}

function getPhone(phone_number) {
  // Keep as-is for now; normalize here if you need E.164 later.
  return (phone_number || "").trim();
}

function getEventDate(event_date) {
  // Expect already formatted by Zapier; leave as-is.
  return event_date;
}

function getEventTime(event_time) {
  return event_time;
}

function getLocation(event_location) {
  return (event_location || "").trim();
}

function getNumHelpers(helpers_requested) {
  // e.g. "I Need 2 Helpers" -> 2
  return parseInt(_firstNumber(helpers_requested), 10) || 0;
}

function getDurationHours(for_how_many_hours) {
  // e.g. "for 4 Hours (minimum)" -> 4
  return Number(_firstNumber(for_how_many_hours)) || 0;
}

function getGuestCountFirst(guests_expected) {
  // Preserve your current behavior: first number only ("10 - 25 Guests" -> 10)
  return parseInt(_firstNumber(guests_expected), 10) || 0;
}

function getOccasion(occasion, occasion_as_you_see_it) {
  // If "Other", use alt text; else keep original
  return consolidateOther(occasion, occasion_as_you_see_it);
}

/***********************
 * Main: transform only
 * - Calls individual functions above
 * - Returns a CLEAN SUBSET (not the whole payload)
 * - Optionally converts yes/no fields to "Checked"/"Unchecked"
 ***********************/
function transformZapierPayload(payload, options) {
  options = options || {};
  const yesNoFields = options.yesNoFields || []; // e.g. ['alcohol_service', 'need_cleanup']

  const {
    first_name,
    last_name,
    email_address,
    phone_number,
    event_date,
    event_time,
    event_location,
    helpers_requested,
    for_how_many_hours,
    guests_expected,
    occasion,
    occasion_as_you_see_it
  } = payload || {};

  const transformed = {
    clientName: getClientName(first_name, last_name),
    first_name: first_name || "",
    last_name: last_name || "",
    email: getEmail(email_address),
    phone: getPhone(phone_number),
    eventDate: getEventDate(event_date),
    eventTime: getEventTime(event_time),
    location: getLocation(event_location),
    numHelpers: getNumHelpers(helpers_requested),
    duration: getDurationHours(for_how_many_hours),
    guestCount: getGuestCountFirst(guests_expected),
    occasion: getOccasion(occasion, occasion_as_you_see_it)
  };

  // Map any specified yes/no fields → "Checked"/"Unchecked"
  yesNoFields.forEach(function (key) {
    if (Object.prototype.hasOwnProperty.call(payload || {}, key)) {
      transformed[key] = yesToCheckbox(payload[key]);
    }
  });

  return transformed;
}

/*

function processLead(parsedData, ss, thread) {
  if (!parsedData) {
    console.error("❌ ERROR: No parsed data available.");
    markLeadAsFailed(thread, "auto-quote-sending-failed");
    return;
  }

  // ✅ Extract relevant data from parsedData
  var {
    clientName,
    eventDate,
    eventTime,
    phone,
    location,
    numHelpers,
    duration,
    email,
    occasion,
    guestCount,
  } = parsedData;

  var emailId = email || "Not Provided"; // ✅ Ensure email is set
  var threadIdForCalendar = thread ? thread.getId() : ""; // 🛠️ Avoid clobbering 'thread'

  // ✅ Validate required data
  if (!validateLeadData(parsedData, clientName)) {
    console.error(
      "🚨 Validation failed. Marking lead as 'auto-quote-sending-failed'."
    );
    markLeadAsFailed(thread, "auto-quote-sending-failed");
    return;
  }

  // ✅ Calculate Total Cost
  var { totalCost, baseCost, additionalCost } = calculateTotalCost(
    numHelpers,
    duration
  );
  console.log(
    `📊 Pricing Breakdown: Base Cost = $${baseCost}, Additional Cost = $${additionalCost}, Total Cost = $${totalCost}`
  );

  // ✅ DEBUG: Log Data Before Sending to `createCalendarEvent`
  console.log("📌 DEBUG: Passing to calendareventcreator.createEvent:", {
    calendarId: ESTIMATE_SENT_CALENDAR_ID,
    clientName,
    occasion,
    guestCount,
    eventDate,
    eventTime,
    phone,
    location,
    numHelpers,
    duration,
    totalCost,
    emailId,
    threadIdForCalendar,
  });

  // ✅ Call External Calendar Event Creator Function
  var eventIdResponse = calendareventcreator.createEvent({
    calendarId: ESTIMATE_SENT_CALENDAR_ID,
    clientName,
    occasion,
    guestCount,
    eventDate,
    eventTime,
    phone,
    location,
    numHelpers,
    duration,
    totalCost,
    emailId,
    threadId: threadIdForCalendar, // ✅ renamed var that still passes string ID,
    dataSource: DataSource.EMAIL_LABEL_EXTRACTING
  });

  // ✅ Handle Missing Event ID Response
  if (!eventIdResponse || eventIdResponse.error) {
    console.error(
      `❌ ERROR: Failed to create calendar event. Reason: ${eventIdResponse.error}`
    );
    markLeadAsFailed(thread, "auto-quote-sending-failed");
    return;
  }

  var eventId = eventIdResponse.eventId; // ✅ Extract eventId from response

  // ✅ Log Quote in Google Sheets
  logQuoteDetails(
    ss,
    referenceNumber,
    clientName,
    eventDate,
    eventTime,
    phone,
    location,
    numHelpers,
    duration,
    totalCost,
    eventId
  );


  sendQuoteAndLog(
    referenceNumber,
    clientName,
    email,
    eventDate,
    eventTime,
    location,
    numHelpers,
    duration,
    totalCost,
    occasion,
    guestCount,
    thread
  );

  // ✅ Schedule Follow-Up
  addFollowUpTask(clientName, eventDate, totalCost);

  // ✅ Mark Lead as Processed
  var processedLabel =
    GmailApp.getUserLabelByName("auto-quote-sent") ||
    GmailApp.createLabel("auto-quote-sent");
  thread.addLabel(processedLabel);
}
*/

/*
function processLead(parsedData, ss, thread) {
  if (!parsedData) {
    console.error("❌ ERROR: No parsed data available.");
    markLeadAsFailed(thread, "auto-quote-sending-failed");
    return;
  }

  // ✅ Extract relevant data from parsedData
  var {
    clientName,
    eventDate,
    eventTime,
    phone,
    location,
    numHelpers,
    duration,
    email,
    occasion,
    guestCount,
  } = parsedData;

  var emailId = email || "Not Provided"; // ✅ Ensure email is set
  var threadIdForCalendar = thread ? thread.getId() : ""; // 🛠️ Avoid clobbering 'thread'

  // ✅ Validate required data
  if (!validateLeadData(parsedData, clientName)) {
    console.error(
      "🚨 Validation failed. Marking lead as 'auto-quote-sending-failed'."
    );
    markLeadAsFailed(thread, "auto-quote-sending-failed");
    return;
  }

  // ✅ Calculate Total Cost
  var { totalCost, baseCost, additionalCost } = calculateTotalCost(
    numHelpers,
    duration
  );
  console.log(
    `📊 Pricing Breakdown: Base Cost = $${baseCost}, Additional Cost = $${additionalCost}, Total Cost = $${totalCost}`
  );

  // ✅ DEBUG: Log Data Before Sending to `createCalendarEvent`
  console.log("📌 DEBUG: Passing to calendareventcreator.createEvent:", {
    calendarId: ESTIMATE_SENT_CALENDAR_ID,
    clientName,
    occasion,
    guestCount,
    eventDate,
    eventTime,
    phone,
    location,
    numHelpers,
    duration,
    totalCost,
    emailId,
    threadIdForCalendar,
  });

  // ✅ Call External Calendar Event Creator Function
  var eventIdResponse = calendareventcreator.createEvent({
    calendarId: ESTIMATE_SENT_CALENDAR_ID,
    clientName,
    occasion,
    guestCount,
    eventDate,
    eventTime,
    phone,
    location,
    numHelpers,
    duration,
    totalCost,
    emailId,
    threadId: threadIdForCalendar, // ✅ renamed var that still passes string ID
  });

  // ✅ Handle Missing Event ID Response
  if (!eventIdResponse || eventIdResponse.error) {
    console.error(
      `❌ ERROR: Failed to create calendar event. Reason: ${eventIdResponse.error}`
    );
    markLeadAsFailed(thread, "auto-quote-sending-failed");
    return;
  }

  var eventId = eventIdResponse.eventId; // ✅ Extract eventId from response

  let referenceNumber = generateShortQuoteID(email, eventDate);
  // ✅ Log Quote in Google Sheets
  logQuoteDetails(
    ss,
    referenceNumber,
    clientName,
    eventDate,
    eventTime,
    phone,
    location,
    numHelpers,
    duration,
    totalCost,
    eventId
  );


  sendQuoteAndLog(
    referenceNumber,
    clientName,
    email,
    eventDate,
    eventTime,
    location,
    numHelpers,
    duration,
    totalCost,
    occasion,
    guestCount,
    thread
  );

  // ✅ Schedule Follow-Up
  addFollowUpTask(clientName, eventDate, totalCost);

  // ✅ Mark Lead as Processed
  var processedLabel =
    GmailApp.getUserLabelByName("auto-quote-sent") ||
    GmailApp.createLabel("auto-quote-sent");
  thread.addLabel(processedLabel);
}
*/

function sendQuoteEmailOnly_v2(
  clientName,
  email,
  eventDate,
  eventTime,
  eventLocation,
  numHelpers,
  duration,
  baseRate,
  hourlyRate,
  totalCost,
  occasion,
  guestCount,
  rateLabel,
  dryRun
) {
  const subject = `Party Helpers for ${occasion} - ${eventDate} - Estimate & Details for ${clientName}`;
  const body = generateQuoteEmail(
    clientName,
    eventDate,
    eventTime,
    eventLocation,
    numHelpers,
    duration,
    baseRate,
    hourlyRate,
    totalCost,
    occasion,
    guestCount,
    rateLabel
  );


if(!dryRun){
 MailApp.sendEmail({
    to: email,
    name: `STL Party Helpers Team`,
    subject: subject,
    htmlBody: body,
    cc: `qa-quote@stlpartyhelpers.com`,
    headers: { "Content-Type": "text/html; charset=UTF-8" }
  });
  console.log(`📧 Sent quote email to ${email} (no thread logging)`);
} else {
  MailApp.sendEmail({
    to: email,
    name: `STL Party Helpers Team`,
    subject: 'Dry Run - ' + subject,
    htmlBody: body,
    headers: { "Content-Type": "text/html; charset=UTF-8" }
  });
  console.log(`📧 Sent quote email to ${email} (no thread logging)`);
}
  
  
}

/*
function sendQuoteEmailOnly(
  clientName,
  email,
  eventDate,
  eventTime,
  eventLocation,
  numHelpers,
  duration,
  totalCost,
  occasion,
  guestCount
) {
  const subject = `Party Helpers for ${occasion} - ${eventDate} - Estimate & Details for ${clientName}`;
  const body = generateQuoteEmail(
    clientName,
    eventDate,
    eventTime,
    eventLocation,
    numHelpers,
    duration,
    totalCost,
    occasion,
    guestCount
  );

  MailApp.sendEmail({
    to: email,
    name: `STL Party Helpers Team`,
    subject: subject,
    htmlBody: body,
    cc: `qa-quote@stlpartyhelpers.com`,
    headers: { "Content-Type": "text/html; charset=UTF-8" }
  });
  console.log(`📧 Sent quote email to ${email} (no thread logging)`);
  
}
*/
function getUpdateFrequency(daysUntilEvent) {
  if (daysUntilEvent >= 548) return 90; // 1.5+ Years (Every 3 months)
  if (daysUntilEvent >= 365) return 60; // 1 Year (Every 2 months)
  if (daysUntilEvent >= 180) return 30; // 6 Months (Every 1 month)
  if (daysUntilEvent >= 90) return 14; // 3 Months (Every 2 weeks)
  if (daysUntilEvent >= 30) return 7; // 1 Month (Every week)
  if (daysUntilEvent >= 14) return 3; // 2 Weeks (Every 3 days)
  if (daysUntilEvent >= 7) return 2; // 1 Week (Every 2 days)
  if (daysUntilEvent >= 3) return 1; // 3 Days (Every day)
  return 0; // 24 Hours (Immediate update)
}

function generateShortQuoteID(email, eventDate) {
  if (!eventDate || !email) return "XXXX"; // Fallback if missing data

  let datePart = eventDate.replace(/-/g, ""); // Remove dashes (YYYYMMDD)

  // Create a simple hash from date + email
  let rawString = datePart + email;
  let hash = rawString
    .split("")
    .reduce((acc, char) => acc + char.charCodeAt(0), 0);

  // Convert hash to a 4-character alphanumeric code
  let shortID = (hash % 1679616).toString(36).toUpperCase().padStart(4, "X");

  return shortID;
}

// TODO: make it so every yeah hourly rate 5 is added to hourly rate and baserate grows by 50 each year
function calculateTotalCost(
  numHelpers,
  duration,
  config = { baseRate: 250, hourlyRate: 45, baseHours: 4, minHours: 4 }
) {
  numHelpers = Number(numHelpers) || 0;
  duration = Number(duration) || 0;

  if (numHelpers < 1) {
    console.error(
      `❌ ERROR: Invalid input - numHelpers: ${numHelpers}, duration: ${duration}`
    );
    return { error: true, message: "Number of helpers must be at least 1." };
  }

  if (duration < config.minHours) {
    console.error(
      `❌ ERROR: Invalid input - duration: ${duration} (Minimum is ${config.minHours} hours)`
    );
    return {
      error: true,
      message: `Minimum duration is ${config.minHours} hours.`,
    };
  }

  let { baseRate, hourlyRate, baseHours } = config;

  let baseCost = numHelpers * baseRate;
  let extraHours = Math.max(0, duration - baseHours);
  let additionalCost = numHelpers * extraHours * hourlyRate;
  let totalCost = baseCost + additionalCost;

  console.log(
    `📊 Total Cost Calculation: Base Cost ($${baseCost}) + Additional Cost ($${additionalCost}) = $${totalCost}`
  );

  return { totalCost, baseCost, additionalCost, numHelpers, duration };
}

function testCalculateTotalCostV1(){
  calculateTotalCost_v1(1, 6)
}

// TODO: add double rate for certain rates 
// TODO: add date to it. So we can start basing it of the date. 
// First step will be to just go off it extracting the year
// That way if the year is 2026 - we can already adjust accordingly 
// Second step could be adding a flat where we check availability for the date
// We can dynamically adjust rates for certain busy dates
function calculateTotalCost_v1(
  numHelpers,
  duration,
  config = { baseRate: 250, hourlyRate: 45, baseHours: 4, minHours: 4 }
) {
  numHelpers = Number(numHelpers) || 0;
  duration = Number(duration) || 0;

  if (numHelpers < 1) {
    console.error(
      `❌ ERROR: Invalid input - numHelpers: ${numHelpers}, duration: ${duration}`
    );
    return { error: true, message: "Number of helpers must be at least 1." };
  }

  if (duration < config.minHours) {
    console.error(
      `❌ ERROR: Invalid input - duration: ${duration} (Minimum is ${config.minHours} hours)`
    );
    return {
      error: true,
      message: `Minimum duration is ${config.minHours} hours.`,
    };
  }

  // 🎯 Apply year-based rate increases
  const baseYear = 2025;
  const currentYear = new Date().getFullYear();
  const yearDiff = Math.max(0, currentYear - baseYear);

  let { baseRate, hourlyRate, baseHours } = config;
  baseRate += yearDiff * 50;
  hourlyRate += yearDiff * 5;

  let baseCost = numHelpers * baseRate;
  let extraHours = Math.max(0, duration - baseHours);
  let additionalCost = numHelpers * extraHours * hourlyRate;
  let totalCost = baseCost + additionalCost;

  console.log(
    `📊 Total Cost Calculation (${currentYear}): Base Cost ($${baseCost}) + Additional Cost ($${additionalCost}) = $${totalCost}`
  );

  return {
    totalCost,
    baseCost,
    additionalCost,
    numHelpers,
    duration,
    yearAdjusted: currentYear,
    baseRate,
    hourlyRate
  };
}


/*******************************
 * Test cases for calculateTotalCost
 *******************************/
/*******************************
 * Test harness for calculateTotalCost
 *******************************/
function assertEqual(actual, expected, msg) {
  if (actual !== expected) {
    throw new Error(`AssertEqual failed: ${msg}\nExpected: ${expected}\nActual:   ${actual}`);
  }
}
function assertClose(actual, expected, epsilon, msg) {
  if (Math.abs(actual - expected) > (epsilon || 1e-9)) {
    throw new Error(`AssertClose failed: ${msg}\nExpected: ${expected}\nActual:   ${actual}`);
  }
}
function assertTrue(value, msg) {
  if (!value) {
    throw new Error(`AssertTrue failed: ${msg}`);
  }
}
function assertError(result, msg) {
  assertTrue(result && result.error === true, msg + " (expected error:true)");
}
function logPass(name) {
  console.log("✅ " + name + " passed");
}



// Run all tests
function run_CalculateTotalCost_AllTests() {
  console.log("🧪 Running all calculateTotalCost_v2 tests...\n");
  
  try {
    test_CalculateTotalCost_v2_InvalidInputs();
    test_CalculateTotalCost_v2_2024_baseline_noExtras();
    test_CalculateTotalCost_v2_2025_escalator();
    test_CalculateTotalCost_v2_2026_escalator();
    test_CalculateTotalCost_v2_busyDate_dec24();
    test_CalculateTotalCost_v2_busyDate_jan1_multiHelpers();
    test_CalculateTotalCost_v2_busyDate_nov27();
    test_CalculateTotalCost_v2_busyDate_dec31();
    test_CalculateTotalCost_v2_dateObjectInput();
    test_CalculateTotalCost_v2_baseHours_effect();
    test_CalculateTotalCost_v2_multiplier_after_escalator_ordering();
    
    console.log("\n🎉 All tests passed successfully!");
  } catch (error) {
    console.error("\n❌ Test failed:", error.message);
    throw error;
  }
}

/*******************************
 * Test cases for calculateTotalCost
 *******************************/

// Alternative test function name (for compatibility)
function test_CalculateTotalCost_v2_InvalidInputs() {
  const cfg = { baseRate: 200, hourlyRate: 45, baseHours: 4, minHours: 4 };
  assertError(calculateTotalCost_v2(0, 4, "2025-01-02", cfg), "numHelpers < 1 should error");
  assertError(calculateTotalCost_v2(1, 3, "2025-01-02", cfg), "duration < minHours should error");
  assertError(calculateTotalCost_v2(1, 4, "not-a-date", cfg), "invalid date should error");
  logPass("test_CalculateTotalCost_v2_InvalidInputs");
}

// 2024 baseline (no escalator), no busy date, no extras
function test_CalculateTotalCost_v2_2024_baseline_noExtras() {
  const r = calculateTotalCost_v2(1, 4, "2024-06-01");
  assertEqual(r.baseRate, 200, "baseRate 2024");
  assertEqual(r.hourlyRate, 45, "hourlyRate 2024");
  assertEqual(r.additionalCost, 0, "no extras at baseHours");
  assertEqual(r.totalCost, 200, "total cost at baseline");
  assertEqual(r.multiplierApplied, 1, "no multiplier");
  logPass("test_CalculateTotalCost_v2_2024_baseline_noExtras");
}

// 2025 escalator test (fixed function name and date)
function test_CalculateTotalCost_v2_2025_escalator() {
  console.log("Testing: 2025 escalator applies (+$50 base, +$5 hourly), no busy date");
  const result = calculateTotalCost_v2(1, 6, "2025-06-15"); // Use June 15th (not a holiday)
  console.log("Result:", result);
  assertEqual(result.baseRate, 250, "2025 baseRate");
  assertEqual(result.hourlyRate, 50, "2025 hourlyRate");
  assertEqual(result.baseCost, 250, "2025 baseCost");
  assertEqual(result.additionalCost, 100, "2h * $50");
  assertEqual(result.totalCost, 350, "total 2025");
  assertEqual(result.multiplierApplied, 1, "no multiplier");
  logPass("test_CalculateTotalCost_v2_2025_escalator");
}

// 2026 escalator applies (+$100 base, +$10 hourly vs 2024), no busy date
function test_CalculateTotalCost_v2_2026_escalator() {
  const r = calculateTotalCost_v2(2, 5, "2026-03-15");
  // 2026 => base 300, hourly 55; extras = 1h
  assertEqual(r.baseRate, 300, "2026 baseRate");
  assertEqual(r.hourlyRate, 55, "2026 hourlyRate");
  assertEqual(r.baseCost, 600, "2 helpers * 300");
  assertEqual(r.additionalCost, 110, "2 * 1h * 55");
  assertEqual(r.totalCost, 710, "total 2026");
  assertEqual(r.multiplierApplied, 1, "no multiplier");
  logPass("test_CalculateTotalCost_v2_2026_escalator");
}

// Busy date doubles (Dec 24)
function test_CalculateTotalCost_v2_busyDate_dec24() {
  const r = calculateTotalCost_v2(1, 6, "2025-12-24T12:00:00Z");
  // 2025 escalator then double: base 250->500, hourly 50->100; extras=2h
  assertEqual(r.baseRate, 500, "Dec 24 baseRate doubled");
  assertEqual(r.hourlyRate, 100, "Dec 24 hourlyRate doubled");
  assertEqual(r.baseCost, 500, "base cost");
  assertEqual(r.additionalCost, 200, "2h * 100");
  assertEqual(r.totalCost, 700, "total on Dec 24");
  assertEqual(r.multiplierApplied, 2, "multiplier 2");
  logPass("test_CalculateTotalCost_v2_busyDate_dec24");
}

// Busy date doubles (Jan 1) with multiple helpers
function test_CalculateTotalCost_v2_busyDate_jan1_multiHelpers() {
  const r = calculateTotalCost_v2(2, 5, "2026-01-01T12:00:00Z");
  // 2026 escalator base 300/hourly 55, doubled => base 600/hourly 110, extras=1h
  assertEqual(r.baseRate, 600, "Jan 1 baseRate doubled");
  assertEqual(r.hourlyRate, 110, "Jan 1 hourlyRate doubled");
  assertEqual(r.baseCost, 1200, "2 * 600");
  assertEqual(r.additionalCost, 220, "2 * 1h * 110");
  assertEqual(r.totalCost, 1420, "total on Jan 1");
  assertEqual(r.multiplierApplied, 2, "multiplier 2");
  logPass("test_CalculateTotalCost_v2_busyDate_jan1_multiHelpers");
}

// Busy date doubles (Nov 27)
function test_CalculateTotalCost_v2_busyDate_nov27() {
  const r = calculateTotalCost_v2(1, 4, "2027-11-27T12:00:00Z");
  // 2027 escalator: base 350, hourly 60; doubled => base 700, hourly 120; extras 0
  assertEqual(r.baseRate, 700, "Nov 27 baseRate doubled");
  assertEqual(r.hourlyRate, 120, "Nov 27 hourlyRate doubled");
  assertEqual(r.totalCost, 700, "total on Nov 27");
  assertEqual(r.multiplierApplied, 2, "multiplier 2");
  logPass("test_CalculateTotalCost_v2_busyDate_nov27");
}

// Dec 31 doubles
function test_CalculateTotalCost_v2_busyDate_dec31() {
  const r = calculateTotalCost_v2(1, 4, "2025-12-31T12:00:00Z");
  // escalated 2025 -> base 250 doubled => 500
  assertEqual(r.baseRate, 500, "Dec 31 baseRate doubled");
  assertEqual(r.hourlyRate, 100, "Dec 31 hourlyRate doubled");
  assertEqual(r.totalCost, 500, "Dec 31 total no extras");
  assertEqual(r.multiplierApplied, 2, "multiplier 2");
  logPass("test_CalculateTotalCost_v2_busyDate_dec31");
}

// Date object input is accepted
function test_CalculateTotalCost_v2_dateObjectInput() {
  const r = calculateTotalCost_v2(1, 4, new Date("2025-12-31T10:00:00Z"));
  assertEqual(r.multiplierApplied, 2, "Date object input OK");
  logPass("test_CalculateTotalCost_v2_dateObjectInput");
}

// Base hours effect (no extras vs extras)
function test_CalculateTotalCost_v2_baseHours_effect() {
  const cfg = { baseRate: 200, hourlyRate: 45, baseHours: 4, minHours: 4 };
  const r1 = calculateTotalCost_v2(1, 4, "2024-07-07", cfg);
  assertEqual(r1.additionalCost, 0, "no extras at baseHours");
  const r2 = calculateTotalCost_v2(1, 6, "2024-07-07", cfg);
  assertEqual(r2.additionalCost, 90, "2h * $45 extras in 2024");
  logPass("test_CalculateTotalCost_v2_baseHours_effect");
}

// Sanity: multiplier applied AFTER escalator (i.e., double the escalated rates)
function test_CalculateTotalCost_v2_multiplier_after_escalator_ordering() {
  const r = calculateTotalCost_v2(1, 5, "2026-12-24T12:00:00Z"); // busy date
  // 2026 escalator: base 300, hourly 55 → doubled: base 600, hourly 110
  assertEqual(r.baseRate, 600, "order of operations base");
  assertEqual(r.hourlyRate, 110, "order of operations hourly");
  // baseCost 600; extras 1h → additional 110; total 710
  assertEqual(r.totalCost, 710, "order of operations total");
  logPass("test_CalculateTotalCost_v2_multiplier_after_escalator_ordering");
}


/*******************************
 * Test harness for calculateTotalCost
 *******************************/
function assertEqual(actual, expected, msg) {
  if (actual !== expected) {
    throw new Error(`AssertEqual failed: ${msg}\nExpected: ${expected}\nActual:   ${actual}`);
  }
}
function assertClose(actual, expected, epsilon, msg) {
  if (Math.abs(actual - expected) > (epsilon || 1e-9)) {
    throw new Error(`AssertClose failed: ${msg}\nExpected: ${expected}\nActual:   ${actual}`);
  }
}
function assertTrue(value, msg) {
  if (!value) {
    throw new Error(`AssertTrue failed: ${msg}`);
  }
}
function assertError(result, msg) {
  assertTrue(result && result.error === true, msg + " (expected error:true)");
}
function logPass(name) {
  console.log("✅ " + name + " passed");
}


/*******************************
 * Test runner
 *******************************/
function test_RunCalculateTotalCost_v2_AllTests() {
  const tests = [
    test_CalculateTotalCost_v2_InvalidInputs,
    test_CalculateTotalCost_v2_2024_baseline_noExtras,
    test_CalculateTotalCost_v2_2025_escalator,
    test_CalculateTotalCost_v2_2026_escalator,
    test_CalculateTotalCost_v2_busyDate_dec24,
    test_CalculateTotalCost_v2_busyDate_jan1_multiHelpers,
    test_CalculateTotalCost_v2_busyDate_nov27,
    test_CalculateTotalCost_v2_busyDate_dec31,
    test_CalculateTotalCost_v2_dateObjectInput,
    test_CalculateTotalCost_v2_baseHours_effect,
    test_CalculateTotalCost_v2_multiplier_after_escalator_ordering
  ];
  var passed = 0;
  for (var i = 0; i < tests.length; i++) {
    var name = tests[i].name;
    try {
      tests[i]();
      passed++;
    } catch (e) {
      console.log("❌ " + name + " FAILED: " + e.message + "\n" + (e.stack || ""));
    }
  }
  console.log("———");
  console.log("Tests passed: " + passed + " / " + tests.length);
}

function try_calc_v2(){
  calculateTotalCost_v2(1, 4, "2026-12-27T12:00:00Z");
}


/**
 * Busy-date multiplier:
 *  - Nov 27
 *  - Dec 24, Dec 25, Dec 31
 *  - Jan 1
 * Yearly escalators:
 *  - Starting in 2025, add +$50 to base and +$5 to hourly each year.
 *    (i.e., 2025 = +1 step vs 2024 baseline)
 */
function calculateTotalCost_v2(
  numHelpers,
  duration,
  dateInput, // string | Date | number (timestamp). If omitted, uses "now".
  config = { baseRate: 250, hourlyRate: 45, baseHours: 4, minHours: 4 }
) {
  // ---- Input normalization
  numHelpers = Number(numHelpers) || 0;
  duration  = Number(duration)  || 0;
  const date = dateInput ? new Date(dateInput) : new Date();
  if (Number.isNaN(date.getTime())) {
    console.error(`❌ ERROR: Invalid date input: ${dateInput}`);
    return { error: true, message: "Invalid date." };
  }

  if (numHelpers < 1) {
    console.error(`❌ ERROR: Invalid input - numHelpers: ${numHelpers}, duration: ${duration}`);
    return { error: true, message: "Number of helpers must be at least 1." };
  }

  if (duration < config.minHours) {
    console.error(`❌ ERROR: Invalid input - duration: ${duration} (Minimum is ${config.minHours} hours)`);
    return { error: true, message: `Minimum duration is ${config.minHours} hours.` };
  }

  // ---- Helpers
  const y = date.getFullYear();
  const m = date.getMonth() + 1;
  const d = date.getDate();

// US Thanksgiving day-of-month for the next 5 years
const THANKSGIVING_BY_YEAR = {
  2025: 27,
  2026: 26,
  2027: 25,
  2028: 23,
  2029: 22,
};

 const isDoubleRateDate = () => {
  const tDay = THANKSGIVING_BY_YEAR[y] ?? 27; // fallback to Nov 27 if not listed
  return (
    (m === 1  && d === 1) ||                 // Jan 1
    (m === 11 && d === tDay) ||              // Thanksgiving (hardcoded per year)
    (m === 12 && (d === 24 || d === 25 || d === 31)) // Dec 24, 25, 31
  );
};

  // ---- Year-based rate increases
  const stepsSince2025 = Math.max(0, y - 2025);

  let { baseRate, hourlyRate, baseHours } = config;
  baseRate  += stepsSince2025 * 50;
  hourlyRate += stepsSince2025 * 5;

  // ---- Busy-date multiplier
  const isHolidayOrBusy = isDoubleRateDate();
  const multiplier = isHolidayOrBusy ? 2 : 1;
  baseRate  *= multiplier;
  hourlyRate *= multiplier;

  // ---- Label for clarity
  let rateLabel = "Base Rate";
  if (isHolidayOrBusy) {
    rateLabel = "Holiday Rate";
  }

  // ---- Math
  const baseCost       = numHelpers * baseRate;
  const extraHours     = Math.max(0, duration - baseHours);
  const additionalCost = numHelpers * extraHours * hourlyRate;
  const totalCost      = baseCost + additionalCost;

  console.log(
    `💰 Total Cost (${date.toISOString().slice(0,10)}): Base $${baseCost} + Additional $${additionalCost} = $${totalCost} [${rateLabel}]`
  );

  return {
    totalCost,
    baseCost,
    additionalCost,
    numHelpers,
    duration,
    dateISO: date.toISOString(),
    yearAdjusted: y,
    baseRate,
    hourlyRate,
    multiplierApplied: multiplier,
    rateLabel // 👈 added
  };
}




function testCalculateTotalCost() {
  var tests = [
    { numHelpers: 4, duration: 4, expected: 800 }, // (4 x $200) + (0 x $45)
    { numHelpers: 4, duration: 5, expected: 980 }, // (4 x $200) + (4 x $45)
    { numHelpers: 2, duration: 6, expected: 580 }, // (2 x $200) + (2 x 2 x $45)
    { numHelpers: 1, duration: 10, expected: 470 }, // (1 x $200) + (1 x 6 x $45)
    { numHelpers: 3, duration: 7, expected: 1005 }, // (3 x $200) + (3 x 3 x $45)

    // ✅ New Test Cases:
    { numHelpers: 0, duration: 5, expected: "ERROR" }, // No helpers
    { numHelpers: 2, duration: 0, expected: "ERROR" }, // Zero duration
    { numHelpers: 1, duration: 2, expected: "ERROR" }, // 🔴 Less than 4 hours (NEW)
    { numHelpers: 2, duration: 3, expected: "ERROR" }, // 🔴 Less than 4 hours (NEW)
  ];

  tests.forEach((test, index) => {
    try {
      var result = calculateTotalCost(test.numHelpers, test.duration);

      // Handling expected errors
      var passed =
        (test.expected === "ERROR" && result.error) ||
        (test.expected !== "ERROR" && result.totalCost === test.expected);

      Logger.log(`Test ${index + 1}: ${passed ? "✅ PASS" : "❌ FAIL"}`);
      Logger.log(
        `    ➡️ Input: numHelpers=${test.numHelpers}, duration=${test.duration}`
      );
      Logger.log(`    ➡️ Expected: ${test.expected}, Got: ${result.totalCost}`);
    } catch (e) {
      if (test.expected === "ERROR") {
        Logger.log(`✅ PASS: Test ${index + 1} correctly threw an error.`);
      } else {
        Logger.log(
          `❌ FAIL: Test ${index + 1} threw an unexpected error: ${e.message}`
        );
      }
    }
  });
}
/*
function tests_runCalculateTotalCost() {
  const tests = [
    {
      description: "✅ Calculates total cost for 2 helpers over 5 hours using default config in the current year (2025).",
      input: { numHelpers: 2, duration: 5 },
      expectError: false
    },
    {
      description: "❌ Should throw an error when number of helpers is zero.",
      input: { numHelpers: 0, duration: 5 },
      expectError: true
    },
    {
      description: "❌ Should reject duration shorter than the minimum required (default min 4 hours).",
      input: { numHelpers: 2, duration: 2 },
      expectError: true
    },
    {
      description: "✅ Confirms no additional charges are applied when duration equals baseHours.",
      input: { numHelpers: 1, duration: 4 },
      expectError: false
    },
    {
      description: "✅ Calculates cost using custom baseRate and hourlyRate from config override.",
      input: {
        numHelpers: 2,
        duration: 6,
        config: { baseRate: 300, hourlyRate: 60, baseHours: 3, minHours: 2 }
      },
      expectError: false
    },
    {
      description: "✅ Applies year-based pricing increase by simulating the year 2026.",
      input: {
        numHelpers: 1,
        duration: 5,
        config: { baseRate: 200, hourlyRate: 45, baseHours: 4, minHours: 2 },
        mockYear: 2026
      },
      expectError: false
    },
    {
      description: "✅ Confirms no pricing adjustment is applied when simulating the base year (2024).",
      input: {
        numHelpers: 1,
        duration: 5,
        config: { baseRate: 200, hourlyRate: 45, baseHours: 4, minHours: 2 },
        mockYear: 2024
      },
      expectError: false
    },
    {
      description: "❌ Should throw an error when number of helpers is negative.",
      input: { numHelpers: -1, duration: 5 },
      expectError: true
    },
    {
      description: "❌ Should fail gracefully on non-numeric input for helpers and duration.",
      input: { numHelpers: "abc", duration: "xyz" },
      expectError: true
    },
    {
      description: "✅ Calculates correct cost for large bookings (3 helpers, 12-hour duration).",
      input: { numHelpers: 3, duration: 12 },
      expectError: false
    }
  ];

  tests.forEach((test, index) => {
    const { numHelpers, duration, config, mockYear } = test.input;

    // Temporarily override system year if needed
    const realDate = Date;
    if (mockYear) {
      globalThis.Date = class extends realDate {
        static now() {
          return new realDate(mockYear, 0, 1).getTime();
        }
        constructor() {
          return new realDate(mockYear, 0, 1);
        }
        static getFullYear() {
          return mockYear;
        }
      };
    }

    let result;
    try {
      result = calculateTotalCost_v1(numHelpers, duration, config);
      const passed = (test.expectError && result.error) || (!test.expectError && !result.error);
      console.log(
        `Test #${index + 1}: ${test.description} → ${passed ? "✅ PASSED" : "❌ FAILED"}`
      );
      if (!passed) console.log("  Result:", JSON.stringify(result, null, 2));
    } catch (err) {
      console.log(`Test #${index + 1}: ${test.description} → ❌ EXCEPTION`);
      console.error(err);
    } finally {
      // Restore real Date object
      if (mockYear) globalThis.Date = realDate;
    }
  });
}
*/

/*
function testCalculateCostMatrix_v1(writeToSheet = true) {
  const maxHelpers = 5;
  const maxHours = 12;
  const yearOverride = 2025;

  const baseConfig = {
    baseRate: 200,
    hourlyRate: 45,
    baseHours: 4,
    minHours: 4
  };

  const headers = [
    "Helpers", "Hours",
    "BaseRate", "HourlyRate", "BaseHours",
    "BaseCost", "AdditionalCost", "ExpectedTotal", "ReturnedTotal",
    "Formula Check", "Mismatch"
  ];

  const data = [headers];
  let mismatchCount = 0;

  const realDate = Date;
  globalThis.Date = class extends realDate {
    constructor() { return new realDate(yearOverride, 0, 1); }
    static now() { return new realDate(yearOverride, 0, 1).getTime(); }
    static getFullYear() { return yearOverride; }
  };

  for (let helpers = 1; helpers <= maxHelpers; helpers++) {
    for (let hours = baseConfig.minHours; hours <= maxHours; hours++) {
      const result = calculateTotalCost_v1(helpers, hours, baseConfig);

      if (!result.error) {
        const expected = result.baseCost + result.additionalCost;
        const mismatch = expected !== result.totalCost ? "❌" : "";
        if (mismatch) mismatchCount++;

        data.push([
          helpers,
          hours,
          result.baseRate,
          result.hourlyRate,
          baseConfig.baseHours,
          result.baseCost,
          result.additionalCost,
          expected,
          result.totalCost,
          `=IF(H${data.length + 1}=I${data.length + 1}, "✅", "❌")`,
          mismatch
        ]);
      } else {
        mismatchCount++;
        data.push([helpers, hours, "ERROR", "", "", "", "", "", "", "", "❌"]);
      }
    }
  }

  globalThis.Date = realDate;

  if (!writeToSheet) {
    const csvOutput = data.map(row => row.join(",")).join("\n");
    Logger.log(csvOutput);
    return ContentService.createTextOutput(csvOutput).setMimeType(ContentService.MimeType.CSV);
  }

const name = `Cost Breakdown Matrix - ${yearOverride} - base${baseConfig.baseRate}_hr${baseConfig.hourlyRate} - ${new Date().toLocaleString().replace(/[/:]/g, "-")}`;
const ss = SpreadsheetApp.create(name);

  const sheet = ss.getActiveSheet();
  sheet.getRange(1, 1, data.length, data[0].length).setValues(data);
  sheet.autoResizeColumns(1, data[0].length);

  // 👇 Enforce minimum column widths
  for (let col = 1; col <= data[0].length; col++) {
    const currentWidth = sheet.getColumnWidth(col);
    if (currentWidth < 100) {
      sheet.setColumnWidth(col, 100);
    }
  }

  // Bold header
  const headerRange = sheet.getRange("1:1");
  headerRange.setFontWeight("bold").setBackground("#e6f0ff").setHorizontalAlignment("center");

  // Center all data
  sheet.getDataRange().setHorizontalAlignment("center");

  // Format dollar columns (BaseCost, Add'l Cost, ExpectedTotal, ReturnedTotal)
  const dollarCols = [6, 7, 8, 9];
  dollarCols.forEach(col => {
    sheet.getRange(2, col, data.length - 1).setNumberFormat('$#,##0.00');
  });

  // Conditional formatting: green H&I when match
  const matchGreen = SpreadsheetApp.newConditionalFormatRule()
    .whenFormulaSatisfied(`=$J2="✅"`)
    .setBackground("#e7ffe7")
    .setRanges([
      sheet.getRange(2, 8, data.length - 1),
      sheet.getRange(2, 9, data.length - 1)
    ])
    .build();

  // Conditional formatting: red J (Formula Check) if ❌
  const formulaFail = SpreadsheetApp.newConditionalFormatRule()
    .whenTextEqualTo("❌")
    .setBackground("#ffe5e5")
    .setFontColor("#b30000")
    .setRanges([sheet.getRange(2, 10, data.length - 1)])
    .build();

  // Conditional formatting: red K (Mismatch column) if ❌
  const mismatchFail = SpreadsheetApp.newConditionalFormatRule()
    .whenTextEqualTo("❌")
    .setBackground("#ffe5e5")
    .setFontColor("#b30000")
    .setRanges([sheet.getRange(2, 11, data.length - 1)])
    .build();

  sheet.setConditionalFormatRules([matchGreen, formulaFail, mismatchFail]);

  // Add summary section
  const totalTests = (maxHelpers) * (maxHours - baseConfig.minHours + 1);
  const summaryRow = data.length + 2;
  sheet.getRange(`A${summaryRow}`).setValue("Total Tests:");
  sheet.getRange(`B${summaryRow}`).setValue(totalTests);
  sheet.getRange(`A${summaryRow + 1}`).setValue("Mismatches:");
  sheet.getRange(`B${summaryRow + 1}`).setValue(mismatchCount);
  sheet.getRange(`A${summaryRow + 2}`).setValue("Mismatch %:");
  sheet.getRange(`B${summaryRow + 2}`).setFormula(`=B${summaryRow + 1}/B${summaryRow}`);
  sheet.getRange(`A${summaryRow}:B${summaryRow + 2}`).setFontWeight("bold");

  Logger.log(`✅ Sheet created: ${ss.getUrl()}`);
  return ContentService.createTextOutput(ss.getUrl());
}
*/

/**
 * Creates a comprehensive cost matrix for calculateTotalCost_v2
 * Tests all combinations of helpers, hours, and different dates/years
 * Saves results to Google Sheets with detailed formatting
 */
function testCalculateCostMatrix_v2(writeToSheet = true) {
  const maxHelpers = 5;
  const maxHours = 12;
  const yearOverride = 2025;

  const baseConfig = {
    baseRate: 200,
    hourlyRate: 45,
    baseHours: 4,
    minHours: 4
  };

  const headers = [
    "Helpers", "Hours", "Date",
    "BaseRate", "HourlyRate", "BaseHours", "Multiplier",
    "BaseCost", "AdditionalCost", "ExpectedTotal", "ReturnedTotal",
    "Formula Check", "Mismatch", "Year", "IsHoliday"
  ];

  const data = [headers];
  let mismatchCount = 0;

  // Test different years and dates
  const testYears = [2025, 2026, 2027];
  const testDates = [
    { year: 2025, month: 6, day: 15, desc: "2025 Regular" },
    { year: 2025, month: 12, day: 24, desc: "2025 Christmas Eve" },
    { year: 2026, month: 1, day: 1, desc: "2026 New Year" },
    { year: 2026, month: 6, day: 15, desc: "2026 Regular" },
    { year: 2027, month: 11, day: 27, desc: "2027 Thanksgiving" },
    { year: 2027, month: 6, day: 15, desc: "2027 Regular" }
  ];

  for (const testDate of testDates) {
    const realDate = Date;
    globalThis.Date = class extends realDate {
      constructor() { return new realDate(testDate.year, testDate.month - 1, testDate.day); }
      static now() { return new realDate(testDate.year, testDate.month - 1, testDate.day).getTime(); }
      static getFullYear() { return testDate.year; }
    };

    for (let helpers = 1; helpers <= maxHelpers; helpers++) {
      for (let hours = baseConfig.minHours; hours <= maxHours; hours++) {
        const result = calculateTotalCost_v2(helpers, hours, `${testDate.year}-${testDate.month.toString().padStart(2, '0')}-${testDate.day.toString().padStart(2, '0')}`, baseConfig);

        if (!result.error) {
          const expected = result.baseCost + result.additionalCost;
          const mismatch = expected !== result.totalCost ? "❌" : "";
          if (mismatch) mismatchCount++;

          // Calculate if it's a holiday
          const isHoliday = (testDate.month === 1 && testDate.day === 1) || 
                           (testDate.month === 11 && testDate.day === 27) || 
                           (testDate.month === 12 && (testDate.day === 24 || testDate.day === 25 || testDate.day === 31));

          data.push([
            helpers,
            hours,
            testDate.desc,
            result.baseRate,
            result.hourlyRate,
            baseConfig.baseHours,
            result.multiplierApplied,
            result.baseCost,
            result.additionalCost,
            expected,
            result.totalCost,
            `=IF(J${data.length + 1}=K${data.length + 1}, "✅", "❌")`,
            mismatch,
            result.yearAdjusted,
            isHoliday ? "🎉" : "📅"
          ]);
        } else {
          mismatchCount++;
          data.push([helpers, hours, testDate.desc, "ERROR", "", "", "", "", "", "", "", "", "❌", "", ""]);
        }
      }
    }

    globalThis.Date = realDate;
  }

  if (!writeToSheet) {
    const csvOutput = data.map(row => row.join(",")).join("\n");
    Logger.log(csvOutput);
    return ContentService.createTextOutput(csvOutput).setMimeType(ContentService.MimeType.CSV);
  }

  const name = `Cost Matrix v2 - ${yearOverride} Baseline - ${new Date().toLocaleString().replace(/[/:]/g, "-")}`;
  const ss = SpreadsheetApp.create(name);

  const sheet = ss.getActiveSheet();
  sheet.getRange(1, 1, data.length, data[0].length).setValues(data);
  sheet.autoResizeColumns(1, data[0].length);

  // Enforce minimum column widths
  for (let col = 1; col <= data[0].length; col++) {
    const currentWidth = sheet.getColumnWidth(col);
    if (currentWidth < 100) {
      sheet.setColumnWidth(col, 100);
    }
  }

  // Bold header
  const headerRange = sheet.getRange("1:1");
  headerRange.setFontWeight("bold").setBackground("#e6f0ff").setHorizontalAlignment("center");

  // Center all data
  sheet.getDataRange().setHorizontalAlignment("center");

  // Format dollar columns
  const dollarCols = [8, 9, 10, 11]; // BaseCost, AdditionalCost, ExpectedTotal, ReturnedTotal
  dollarCols.forEach(col => {
    sheet.getRange(2, col, data.length - 1).setNumberFormat('$#,##0.00');
  });

  // Conditional formatting: green when match
  const matchGreen = SpreadsheetApp.newConditionalFormatRule()
    .whenFormulaSatisfied(`=$L2="✅"`)
    .setBackground("#e7ffe7")
    .setRanges([
      sheet.getRange(2, 10, data.length - 1), // ExpectedTotal
      sheet.getRange(2, 11, data.length - 1)  // ReturnedTotal
    ])
    .build();

  // Conditional formatting: red for mismatches
  const formulaFail = SpreadsheetApp.newConditionalFormatRule()
    .whenTextEqualTo("❌")
    .setBackground("#ffe5e5")
    .setFontColor("#b30000")
    .setRanges([
      sheet.getRange(2, 12, data.length - 1), // Formula Check
      sheet.getRange(2, 13, data.length - 1)  // Mismatch
    ])
    .build();

  // Conditional formatting: highlight holidays
  const holidayHighlight = SpreadsheetApp.newConditionalFormatRule()
    .whenTextEqualTo("🎉")
    .setBackground("#fff2cc")
    .setFontColor("#d6b656")
    .setRanges([sheet.getRange(2, 15, data.length - 1)]) // IsHoliday column
    .build();

  sheet.setConditionalFormatRules([matchGreen, formulaFail, holidayHighlight]);

  // Add summary section
  const totalTests = testDates.length * maxHelpers * (maxHours - baseConfig.minHours + 1);
  const summaryRow = data.length + 2;
  
  sheet.getRange(`A${summaryRow}`).setValue("📊 SUMMARY");
  sheet.getRange(`A${summaryRow}`).setFontWeight("bold").setFontSize(12);
  
  sheet.getRange(`A${summaryRow + 1}`).setValue("Total Tests:");
  sheet.getRange(`B${summaryRow + 1}`).setValue(totalTests);
  
  sheet.getRange(`A${summaryRow + 2}`).setValue("Mismatches:");
  sheet.getRange(`B${summaryRow + 2}`).setValue(mismatchCount);
  
  sheet.getRange(`A${summaryRow + 3}`).setValue("Success Rate:");
  sheet.getRange(`B${summaryRow + 3}`).setFormula(`=1-B${summaryRow + 2}/B${summaryRow + 1}`);
  sheet.getRange(`B${summaryRow + 3}`).setNumberFormat('0.00%');
  
  sheet.getRange(`A${summaryRow + 4}`).setValue("Test Years:");
  sheet.getRange(`B${summaryRow + 4}`).setValue(testYears.join(", "));
  
  sheet.getRange(`A${summaryRow + 5}`).setValue("Holiday Dates:");
  sheet.getRange(`B${summaryRow + 5}`).setValue("Jan 1, Nov 27, Dec 24/25/31");
  
  sheet.getRange(`A${summaryRow + 1}:B${summaryRow + 5}`).setFontWeight("bold");

  // Add legend
  const legendRow = summaryRow + 7;
  sheet.getRange(`A${legendRow}`).setValue("📋 LEGEND");
  sheet.getRange(`A${legendRow}`).setFontWeight("bold").setFontSize(12);
  
  sheet.getRange(`A${legendRow + 1}`).setValue("🎉 = Holiday (2x rates)");
  sheet.getRange(`A${legendRow + 2}`).setValue("📅 = Regular day");
  sheet.getRange(`A${legendRow + 3}`).setValue("✅ = Test passed");
  sheet.getRange(`A${legendRow + 4}`).setValue("❌ = Test failed");

  Logger.log(`✅ Sheet created: ${ss.getUrl()}`);
  return ContentService.createTextOutput(ss.getUrl());
}


function sendQuoteAndLog(
  clientName,
  email,
  eventDate,
  eventTime,
  eventLocation,
  numHelpers,
  duration,
  totalCost,
  occasion,
  guestCount,
  thread
) {
  try {
    console.log(`🧚 DEBUG: typeof thread = ${typeof thread}`);
    console.log(
      `🧚 DEBUG: thread.getMessages? ${typeof thread.getMessages === "function"
      }`
    );

    if (!thread || typeof thread.getMessages !== "function") {
      throw new Error("Invalid thread object passed to sendQuoteAndLog.");
    }

    const subject = `Party Helpers for ${occasion} - ${eventDate} - Estimate & Details for ${clientName}`;
    const body = generateQuoteEmail(
      clientName,
      eventDate,
      eventTime,
      eventLocation,
      numHelpers,
      duration,
      totalCost,
      occasion,
      guestCount
    );
    const cleanedBody = cleanEmailBody(body);

    const messages = thread.getMessages();
    const lastMessage = messages[messages.length - 1];


    // ✅ Send email to client
    MailApp.sendEmail({
      to: email,
      name: `STL Party Helpers Team`,
      subject: subject,
      htmlBody: body,
      cc: `qa-quote@stlpartyhelpers.com`,
      headers: { "Content-Type": "text/html; charset=UTF-8" }
    });
    console.log(`📧 Sent quote email to ${email}`);


    // ✅ Add internal log reply to thread
    lastMessage.getThread().reply(cleanedBody, {
      htmlBody: cleanedBody,
      subject: `Re: ` + subject,
      name: `STL Party Helpers Team`,
      headers: { "Content-Type": "text/html; charset=UTF-8" },
    });
    console.log(`📩 Internal log reply added to lead thread`);
  } catch (error) {
    console.error("❌ ERROR in sendQuoteAndLog:", error.message);
    throw error;
  }
}

function cleanEmailBody(body) {
  return body.normalize("NFKD").replace(/[^\x00-\x7F]/g, ""); // Removes non-ASCII characters
}

function getOrCreateSpreadsheet(folder) {
  var files = folder.getFilesByType(MimeType.GOOGLE_SHEETS);
  if (files.hasNext()) {
    return SpreadsheetApp.openById(files.next().getId());
  }
  var ss = SpreadsheetApp.create("Leads");
  folder.addFile(DriveApp.getFileById(ss.getId()));
  return ss;
}

function getUnprocessedLeads() {
  console.log("📌 Fetching unprocessed leads...");

  var teamEmail = "team@stlpartyhelpers.com"; // Change this if needed
  var leadLabel = getGmailLabel(teamEmail, LEAD_ID_AUTO);
  var failedLabel = getGmailLabel(teamEmail, "auto-quote-sending-failed");
  var sentLabel = getGmailLabel(teamEmail, "auto-quote-sent");
  var manualSentLabel = getGmailLabel(teamEmail, "@-quote-manually-sent");

  if (!leadLabel) {
    console.error("⚠️ ERROR: Label" + LEAD_ID_AUTO + " not found.");
    return [];
  }

  var threads = leadLabel.getThreads();
  if (threads.length === 0) {
    console.log("✅ No new leads found under " + LEAD_ID_AUTO + ".");
    return [];
  }

  var leads = [];
  console.log(`📬 Checking ${threads.length} thread(s) for new leads...`);

  threads.forEach((thread) => {
    try {
      var threadLabels = thread.getLabels().map((label) => label.getName());

      // ✅ Skip if thread has "auto-quote-sending-failed"
      if (failedLabel && threadLabels.includes(failedLabel.getName())) {
        console.log(
          `⏩ SKIPPED: Thread "${thread.getFirstMessageSubject()}" already marked as 'auto-quote-sending-failed'.`
        );
        return;
      }

      // ✅ Skip if thread already has "auto-quote-sent"
      if (sentLabel && threadLabels.includes(sentLabel.getName())) {
        console.log(
          `⏩ SKIPPED: Thread "${thread.getFirstMessageSubject()}" already marked as 'auto-quote-sent'.`
        );
        return;
      }

      // ✅ Skip if thread has "@-quote-manually-sent"
      if (manualSentLabel && threadLabels.includes(manualSentLabel.getName())) {
        console.log(
          `⏩ SKIPPED: Thread "${thread.getFirstMessageSubject()}" marked as '@-quote-manually-sent'.`
        );
        return;
      }

      var messages = thread.getMessages();
      if (messages.length === 0) {
        console.warn(
          `⚠️ SKIPPED: Thread "${thread.getFirstMessageSubject()}" has no messages.`
        );
        return;
      }

      messages.forEach((message) => {
        var subject = message.getSubject();
        var body = message.getPlainBody();
        var sender = message.getFrom();

        console.log(`📩 Found lead: ${subject} from ${sender}`);
        console.log(`✅ Adding lead to processing list: ${subject}`);

        leads.push({ subject, body, sender, thread });
      });
    } catch (error) {
      console.error(`❌ ERROR processing thread: ${error.message}`);
    }
  });

  console.log(
    `🔍 FINAL COUNT: ${leads.length} new lead(s) ready for processing.`
  );
  return leads;
}

/**
 * Fetches a Gmail label from a specified account.
 * Uses delegation to access another Gmail account's labels.
 * @param {string} account - The email account to retrieve labels from.
 * @param {string} labelName - The label to fetch.
 * @returns {GmailLabel|null}
 */
function getGmailLabel(account, labelName) {
  try {
    var label = GmailApp.getUserLabelByName(labelName);
    if (!label) {
      console.warn(`⚠️ WARNING: Label '${labelName}' not found in ${account}`);
      return null;
    }
    return label;
  } catch (error) {
    console.error(
      `❌ ERROR: Could not access label '${labelName}' for ${account}: ${error.message}`
    );
    return null;
  }
}

function cleanParsedData(parsedData) {
  Object.keys(parsedData).forEach((key) => {
    if (typeof parsedData[key] === "string") {
      parsedData[key] = parsedData[key]
        .normalize("NFKD")
        .replace(/[^\x00-\x7F]/g, ""); // Removes non-ASCII characters
      parsedData[key] = parsedData[key].replace(/^[?\s]+/, ""); // Remove leading ? and spaces
    }
  });
}

function markLeadAsFailed(thread, labelName) {
  try {
    var failedLabel =
      GmailApp.getUserLabelByName(labelName) || GmailApp.createLabel(labelName);
    //thread.addLabel(failedLabel);
    console.log(`⚠️ Marked lead as failed: ${thread.getFirstMessageSubject()}`);
  } catch (error) {
    console.error(`❌ ERROR: Failed to mark lead as "${labelName}".`, error);
  }
}

function parseLeadData(body) {
  // Remove WPForms junk formatting and trim extra spaces
  body = body
    .replace(/\[image:.*?\]/g, "")
    .replace(/\*/g, "")
    .trim();

  var lines = body
    .split("\n")
    .map((line) => line.trim())
    .filter((line) => line.length > 0);
  var parsedData = {};

  var fieldMappings = {
    "What Is Your Role In This Event?": "role",
    "Occasion?": "occasion",
    "When?": "eventDateTime",
    "Event Location Information": "location",
    "Guests Expected": "guestCount",
    "How Many Party Helpers Needed?": "numHelpers",
    "For How Many Hours?": "duration",
    "First Name": "firstName",
    "Last Name": "lastName",
    "Contact Phone Number": "phone",
    "Contact Email Address": "email",
  };

  // Extract values using regex
  lines.forEach((line) => {
    Object.keys(fieldMappings).forEach((field) => {
      let regex = new RegExp(`^\\s*${field}\\s*:?\\s*(.*)`, "i");
      let match = line.match(regex);
      if (match) {
        parsedData[fieldMappings[field]] = match[1]
          .trim()
          .normalize("NFKD")
          .replace(/[^\x00-\x7F]/g, ""); // Normalize & Remove non-ASCII
      }
    });
  });

  // Clean up extracted values
  Object.keys(parsedData).forEach((key) => {
    if (typeof parsedData[key] === "string") {
      parsedData[key] = parsedData[key].replace(/^[?\s]+/, "").trim(); // Remove leading ? and spaces
    }
  });

  // 🛠 Extract phone number safely
  let phoneRegex =
    /(\+?\d{1,3}[-.\s]?(\(\d{1,3}\)|\d{1,3})[-.\s]?\d{3}[-.\s]?\d{4,6})/;
  let phoneMatch = parsedData.phone ? parsedData.phone.match(phoneRegex) : null;
  parsedData.phone = phoneMatch ? phoneMatch[0].trim() : "Not Provided";

  // 🛠 Extract numbers correctly from "I Need 4 Helpers" and "for 5 Hours"
  parsedData.numHelpers = parseInt(
    (parsedData.numHelpers || "").match(/\d+/)?.[0] || "0",
    10
  );
  parsedData.duration = parseInt(
    (parsedData.duration || "").match(/\d+/)?.[0] || "0",
    10
  );

  // 🛠 Extract event date & time safely
  let eventDateTimeStr = parsedData.eventDateTime || "";
  if (eventDateTimeStr) {
    let dateTimeParts = eventDateTimeStr.match(
      /(.+?)\s+(\d{1,2}:\d{2}\s*(AM|PM))/i
    );
    if (dateTimeParts) {
      parsedData.eventDate = dateTimeParts[1].trim(); // Extracts "April 11, 2025"
      parsedData.eventTime = dateTimeParts[2].trim(); // Extracts "8:30 AM"
    } else {
      console.error(
        "❌ ERROR: Unable to parse eventDateTime:",
        eventDateTimeStr
      );
      parsedData.eventDate = "Invalid";
      parsedData.eventTime = "Invalid";
    }
  } else {
    console.error("❌ ERROR: Missing eventDateTime.");
    parsedData.eventDate = "Invalid";
    parsedData.eventTime = "Invalid";
  }

  // 🛠 Ensure client name is always defined
  parsedData.clientName = `${parsedData.firstName || ""} ${parsedData.lastName || ""
    }`.trim();
  if (!parsedData.clientName || parsedData.clientName === " ") {
    parsedData.clientName = "Unknown Client";
  }

  // ✅ Clean parsed data to remove any encoding issues
  cleanParsedData(parsedData);

  // 🛠 Log results to catch parsing issues
  console.log("📋 Parsed Lead Data:", JSON.stringify(parsedData, null, 2));

  return parsedData;
}

function validateLeadData(parsedData, clientName) {
  let requiredFields = [
    "clientName",
    "eventDate",
    "eventTime",
    "location",
    "numHelpers",
    "duration",
  ];
  let missingFields = [];

  requiredFields.forEach((field) => {
    if (
      !parsedData[field] ||
      parsedData[field] === "Invalid" ||
      parsedData[field] === 0
    ) {
      missingFields.push(`${field}: "${parsedData[field]}"`);
    }
  });

  if (missingFields.length > 0) {
    console.error(
      `❌ ERROR: Missing required fields for lead: New Lead: ${clientName}, ${parsedData.eventDate} ${parsedData.eventTime}`
    );
    console.error(`⛔ Missing Fields: ${missingFields.join(" | ")}`);
    return false;
  }

  return true;
}

function getLastNameInitial(lastName) {
  if (typeof lastName !== "string" || lastName.trim() === "") {
    return ""; // Return empty string if lastName is invalid or empty
  }
  return lastName.trim().charAt(0).toUpperCase() + ".";
}

function logQuoteDetails(
  ss,
  referenceNumber,
  clientName,
  eventDate,
  eventTime,
  location,
  numHelpers,
  duration,
  totalCost,
  eventId
) {
  var logSheet =
    ss.getSheetByName("Quote Logs") || ss.insertSheet("Quote Logs");
  logSheet.appendRow([
    referenceNumber,
    clientName,
    eventDate,
    eventTime,
    location,
    numHelpers,
    duration,
    totalCost,
    eventId,
  ]);
  console.log(`📊 Logged quote for ${clientName}.`);
}

function addFollowUpTask(clientName, eventDate, totalCost) {
  var daysUntilEvent =
    (new Date(eventDate) - new Date()) / (1000 * 60 * 60 * 24);
  var followUpDate = new Date();
  followUpDate.setDate(followUpDate.getDate() + (daysUntilEvent > 7 ? 7 : 2));
  console.log(
    `📅 Follow-up task scheduled for ${clientName} on ${followUpDate}`
  );
}

function wrapEmailHTML(content) {
  return `
    <!DOCTYPE html>
    <html>
      <head>
        <meta charset="UTF-8">
        <style>
          body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 0;
            color: #000;
          }
          table {
            width: 100%;
            border-collapse: collapse;
            margin: 20px 0;
          }
          td, th {
            padding: 12px;
            border: 1px solid #ccc;
            text-align: left;
          }
        </style>
      </head>
      <body>
        ${content}
      </body>
    </html>
  `;
}

function emailHeaderSection() {
  return `
    <div style="text-align: center; margin-bottom: 20px;">
      <img src="https://stlpartyhelpers.com/wp-content/uploads/2024/05/FullLogo_Transparent_NoBuffer-1-e1716631491571.png" alt="STL Party Helpers" style="max-height: 60px;" />
      <h2 style="margin-top: 8px; color: #333;">Your Party Helpers Estimate</h2>
    </div>
  `;
}

function eventDetailsSection(data) {
  return `
    <h3>Event Details</h3>
    <table>
      <tr><td><strong>Client Name:</strong></td><td>${data.clientName}</td></tr>
      <tr><td><strong>Occasion:</strong></td><td>${data.occasion}</td></tr>
      <tr><td><strong>Date:</strong></td><td>${data.eventDate}</td></tr>
      <tr><td><strong>Time:</strong></td><td>${data.eventTime}</td></tr>
      <tr><td><strong>Location:</strong></td><td>${data.eventLocation}</td></tr>
      <tr><td><strong>Guest Count:</strong></td><td>${data.guestCount || "Not Provided"}</td></tr>
    </table>
  `;
}

function staffingAndPricingSection(data) {
  return `
    <h3>Staffing & Pricing</h3>
    <table>
      <tr><td><strong>Helpers Needed:</strong></td><td>${data.numHelpers}</td></tr>
      <tr><td><strong>Duration:</strong></td><td>${data.duration} hours</td></tr>
      <tr><td><strong>Hourly Rate:</strong></td><td>$${data.rate}</td></tr>
      <tr><td><strong>Total Estimate:</strong></td><td><strong>$${data.total}</strong></td></tr>
    </table>
  `;
}

function notesSection(notes) {
  if (!notes) return "";
  const safeNotes = notes
    .replace(/\u2013|\u2014/g, "-")
    .replace(/\u2018|\u2019/g, "'")
    .replace(/\u201C|\u201D/g, '"')
    .replace(/\u2026/g, "...");
  return `
    <h3>Notes</h3>
    <p>${safeNotes}</p>
  `;
}

function footerSection() {
  return `
    <hr style="margin: 30px 0;">
    <p style="font-size: 12px; color: #666; text-align: center;">
      STL Party Helpers · <a href="https://stlpartyhelpers.com">stlpartyhelpers.com</a><br>
      Questions? Just reply to this email or call us.
    </p>
  `;
}


function generateQuoteEmail_v1(data) {
  return `
    <!DOCTYPE html>
    <html>
      <head>
        <meta charset="UTF-8" />
        <style>
          body { font-family: Arial, sans-serif; color: #222; }
          table { width: 100%; border-collapse: collapse; margin-bottom: 20px; }
          td { border: 1px solid #ccc; padding: 10px; }
          h3 { margin-top: 30px; color: #444; }
        </style>
      </head>
      <body>
        ${emailHeaderSection()}
        ${eventDetailsSection(data)}
        ${staffingAndPricingSection(data)}
        ${notesSection(data.notes)}
        ${footerSection()}
      </body>
    </html>
  `;
}

function generateHtmlQuoteV1({ clientName, eventDate, eventTime, eventLocation, occasion, guestCount, helpers, hours, rate, total }) {
  return `<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8" />
    <title>STL Party Helpers - Quote</title>
  </head>
  <body style="margin:0; padding:0; font-family: Arial, sans-serif; background-color: #ffffff; color: #333;">
    <table width="100%" cellpadding="0" cellspacing="2" border="0" style="background-color: #ffffff;">
      <tr>
        <td align="center" style="padding: 0 16px;">
          <table width="100%" cellpadding="0" cellspacing="0" border="0" style="max-width: 650px; border: 1px solid #ccc; padding: 20px;">
            <tr>
              <td align="center" style="padding: 5px;">
                <img src="https://stlpartyhelpers.com/wp-content/uploads/2025/08/stlph-logo-2.jpg" alt="STL Party Helpers Logo" />
              </td>
            </tr>
            <tr>
              <td align="center" style="font-size: 22px; font-weight: bold;">Hi ${clientName}!</td>
            </tr>
            <tr>
              <td align="center" style="padding-bottom: 10px;">
                Thank you for reaching out!<br />
                Below is your event quote, along with important details and next steps.
              </td>
            </tr>

            <!-- Pricing -->
            <tr>
              <td style="font-size: 14px; font-weight: bold; padding-top: 15px;">Our Rates & Pricing</td>
            </tr>
            <tr>
              <td>
                <table width="100%" cellpadding="0" cellspacing="0" border="0" style="background-color: #f9f9f9; margin-top: 10px;">
                  <tr>
                    <td style="padding: 8px 10px; font-weight: bold;">Base Rate:</td>
                    <td style="padding: 8px 10px;">$200 / helper (first 4 hours)</td>
                  </tr>
                  <tr>
                    <td style="padding: 8px 10px; font-weight: bold;">Additional Hours:</td>
                    <td style="padding: 8px 10px;">$${rate} per additional hour per helper</td>
                  </tr>
                  <tr>
                    <td style="padding: 8px 10px; font-weight: bold;">Estimated Total:</td>
                    <td style="padding: 8px 10px;">${total}</td>
                  </tr>
                </table>
                <p style="font-size: 12px; color: #666; padding-top: 5px;">
                  Final total may adjust based on our call. Gratuity is not included but always appreciated!
                </p>
              </td>
            </tr>

            <!-- Event Details -->
            <tr>
              <td style="font-size: 14px; font-weight: bold; padding-top: 20px;">Event Details</td>
            </tr>
            <tr>
              <td>
                <table width="100%" cellpadding="0" cellspacing="0" border="0" style="background-color: #f9f9f9; margin-top: 10px;">
                  <tr><td style="padding: 8px 10px; font-weight: bold;">&#128197; When:</td><td style="padding: 8px 10px;">${eventDate} ${eventTime}</td></tr>
                  <tr><td style="padding: 8px 10px; font-weight: bold;">&#128205; Where:</td><td style="padding: 8px 10px;">${eventLocation}</td></tr>
                  <tr><td style="padding: 8px 10px; font-weight: bold;">&#127760; Occasion:</td><td style="padding: 8px 10px;">${occasion}</td></tr>
                  <tr><td style="padding: 8px 10px; font-weight: bold;">&#128101; Guest Count:</td><td style="padding: 8px 10px;">${guestCount}</td></tr>
                  <tr><td style="padding: 8px 10px; font-weight: bold;">&#129491; Helpers Needed:</td><td style="padding: 8px 10px;">${helpers}</td></tr>
                  <tr><td style="padding: 8px 10px; font-weight: bold;">&#9201; For How Long:</td><td style="padding: 8px 10px;">${hours} Hours</td></tr>
                </table>
              </td>
            </tr>

            <!-- Services -->
            <tr>
              <td style="font-size: 14px; font-weight: bold; padding-top: 15px;">Services Included</td>
            </tr>
            <tr>
              <td style="padding: 10px 0;">
                <ul style="padding-left: 20px; margin: 0;">
                  <li><strong>Setup & Presentation</strong>
                    <ul>
                      <li>Arranging tables, chairs, and decorations</li>
                      <li>Buffet setup & live buffet service</li>
                      <li>Butler-passed appetizers & cocktails</li>
                    </ul>
                  </li>
                  <li><strong>Dining & Guest Assistance</strong>
                    <ul>
                      <li>Multi-course plated dinners</li>
                      <li>General bussing (plates, silverware, glassware)</li>
                      <li>Beverage service (water, wine, champagne, coffee, etc.)</li>
                      <li>Special services (cake cutting, dessert plating, etc.)</li>
                    </ul>
                  </li>
                  <li><strong>Cleanup & End-of-Event Support</strong>
                    <ul>
                      <li>Washing dishes, managing trash, and keeping the event space tidy</li>
                      <li>Kitchen cleanup & end-of-event breakdown</li>
                      <li>Assisting with food storage & leftovers</li>
                    </ul>
                  </li>
                </ul>
                <p>Need something specific? Let us know! We’ll do our best to accommodate your request.</p>
              </td>
            </tr>

            <!-- Payment -->
            <tr><td style="font-size: 14px; font-weight: bold; padding-top: 15px;">Payment Options</td></tr>
            <tr><td style="background-color: #f9f9f9; padding: 10px;">Check, Debit / Credit Cards (via Stripe), Venmo, Zelle, Apple Pay</td></tr>

            <!-- Next Steps -->
            <tr><td style="font-size: 14px; font-weight: bold; padding-top: 15px;">What Happens Next</td></tr>
            <tr>
              <td style="background-color: #f9f9f9; padding: 10px;">
                <span style="text-align:center; font-size: 14px; font-weight: bold;">Booked already?</span><br />
                <table cellpadding="0" cellspacing="0" border="0" style="font-size: 14px;">
                  <tr><td valign="top" style="padding-right: 8px;">1.</td><td>We’ll call you at your scheduled time to go over details.</td></tr>
                  <tr><td valign="top" style="padding-right: 8px;">2.</td><td>If all looks good after our call, we’ll send a Stripe deposit link to proceed.</td></tr>
                  <tr><td valign="top" style="padding-right: 8px;">3.</td><td>Once the deposit is in, your reservation is locked in.</td></tr>
                </table>
                <p style="font-size: 13px; text-align: center; color: #666; margin-top: 8px;">Deposit is 40–50% of the estimate rounded for simplicity.</p>
                <p style="font-size: 13px; text-align: center; color: #666; margin-top: 5px;">❌ Required to confirm your reservation.</p>
              </td>
            </tr>
            <tr>
              <td style="background-color: #fff4e5; text-align: center; padding: 10px; margin-top: 5px; border: 1px solid #fddfb4;">
                <strong>Haven’t scheduled a call yet?</strong><br />
                <strong>Book now to get started</strong><br />
                <span style="font-size: 0.9em;">(to confirm helpers, tasks, and setup)</span><br />
                <a href="https://calendly.com/stlpartyhelpers/quote-intake" style="display:inline-block; background-color:#0047ab; color:#fff; padding:8px 14px; margin-top: 12px; text-decoration:none; font-weight:bold; border-radius:4px;">Click Here to Schedule Appointment</a>
              </td>
            </tr>

            <!-- Footer -->
            <tr>
              <td align="center" style="font-size: 12px; padding-top: 20px; color: #666;">
                4220 Duncan Ave., Ste. 201, St. Louis, MO 63110<br />
                <a href="tel:+13147145514" style="display:inline-block;background-color:#ffffff;color:#000000;padding: 9px 10px;text-decoration:none;border-radius:4px;margin-top:12px;margin-bottom: 12px;border: 1px solid gray;" target="_blank">Tap to Call Us: (314) 714-5514</a><br />
                <a href="https://stlpartyhelpers.com" style="color:#0047ab; display: inline-block; margin-bottom: 8px;">stlpartyhelpers.com</a><br />
                &copy; 2025 STL Party Helpers<br />
                <span style="font-size: 0.55em;">v1.1</span>
              </td>
            </tr>
          </table>
        </td>
      </tr>
    </table>
  </body>
</html>`;
}


function generateQuoteEmail(
  clientName,
  eventDate,
  eventTime,
  eventLocation,
  numHelpers,
  duration,
  baseRate,
  hourlyRate,
  totalCost,
  occasion,
  guestCount,
  rateLabel
) 
{
  return `<head>
  <meta charset="UTF-8"> <!-- ✅ Forces UTF-8 encoding -->
  <style>
    body {
      font-family: Arial, sans-serif;
      font-size: 14px;
    }
    
    .container {
      max-width: 700px;
      margin: auto;
      padding: 0px;
      background-color: #f9f9f9;
      border-radius: 8px;
      box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
    }

    .header {
      background-color: #262578;
      padding: 5px;
      border-radius: 8px 8px 0 0;
      color: white;
      text-align: center;
    }

    .header h3 {
      margin: 0;
      display: inline-block;
      background-color: white;
      color: black;
      padding: 2px;
    }
body, p, .table td {
  font-size: 14px;
  line-height: 1.2;
}
    .content {
      padding: 5px;
      background-color: white;
      border-radius: 0 0 8px 8px;
    }

    .section-title {
      color: #333;
      padding: 0px;
      padding-bottom: 8px;
      margin: 0px;
      padding-top: 12px;
    }

    .section {
      
      margin: 0px 0;
      padding-top: 0px;
      padding-bottom: 10px;
    }
    tr {
      padding-bottom: 7px;
    }
    .table {
      width: 100%;
      font-family: Arial, sans-serif;
      line-height: 1.5;
      border-collapse: collapse;
      padding-bottom: 20px;
    }

    .table td {
      padding: 5px;
    }

    .highlight-box {
      display: inline-block;
      background-color: #f8ff94;
      padding: 10px;
      border-radius: 6px;
      font-size: 15px;
      font-weight: bold;
      border: 1px solid #d4c200;
    }

    .button {
      color: #ffffff;
      text-decoration: none;
      display: inline-block;
      background-color: #673AB7;
      padding: 5px;
      border-radius: 6px;
      font-size: 15px;
      font-weight: bold;
    }

    .footer {
      margin-top: 5px;
      padding: 5px;
      
      color: white;
      text-align: center;
      border-radius: 8px;
    }

    .footer a {
      color: white;
      text-decoration: none;
    }

    .italic-text {
      color: #555;
      font-style: italic;
      font-size: 12px;
    }
    .list p{
      font-size: 13px;
      line-height: 1;
    }

  </style>
</head>
<body>
  <div class="container">
    <table class="email_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0">
    <tbody>
      <tr>
        <td class="email_bg bg_light px pt_md" style="font-size: 0;text-align: center;line-height: 100%;background-color: #F7F7FA;padding-left: 8px;padding-right: 8px;padding-top: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
          <!--[if (mso)|(IE)]>
          <table role="presentation" width="800" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
            <tbody>
              <tr>
                <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
          <![endif]-->
          <table class="content_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="max-width: 800px;margin: 0 auto;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
            <tbody>
              <tr>
                <td class="content_cell bg_white brounded_top bt_primary px py" style="font-size: 0;text-align: center;background-color: #ffffff;border-top: 4px solid #262578;border-radius: 4px 4px 0 0;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                  <div class="column_row" style="font-size: 0;text-align: center;max-width: 624px;margin: 0 auto;">
                    <!--[if (mso)|(IE)]>
                    <table role="presentation" width="624" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
                      <tbody>
                        <tr>
                          <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
                    <![endif]-->
                      <table class="column" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="vertical-align: top;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                        <tbody>
                          <tr>
                            <td class="column_cell px py_xs text_primary text_center" style="vertical-align: top;color: #2376dc;text-align: center;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                              <p class="img_inline" style="color: inherit;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;font-size: 16px;line-height: 100%;clear: both;"><a href="#" style="-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;text-decoration: none;color: #2376dc;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;"><img  src="https://stlpartyhelpers.com/wp-content/uploads/2024/05/FullLogo_Transparent_NoBuffer-1-e1716631491571.png" width="110" height="" alt="STL Party Helpers is always here for you!" style="max-width: 110px;-ms-interpolation-mode: bicubic;border: 0;height: auto;line-height: 100%; display: block; margin: 0 auto;outline: none;text-decoration: none;"></a></p>
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    <!--[if (mso)|(IE)]>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                    <![endif]-->
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          <!--[if (mso)|(IE)]>
                </td>
              </tr>
            </tbody>
          </table>
          <![endif]-->
        </td>
      </tr>
    </tbody>
  </table>
  <table class="email_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0">
    <tbody>
      <tr>
        <td class="email_bg bg_light px" style="font-size: 0;text-align: center;line-height: 100%;background-color: #F7F7FA;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
          <!--[if (mso)|(IE)]>
          <table role="presentation" width="800" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
            <tbody>
              <tr>
                <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
          <![endif]-->
          <table class="content_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="max-width: 800px;margin: 0 auto;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
            <tbody>
              <tr>
                <td class="content_cell bg_white px py_md" style="font-size: 0;text-align: center;background-color: #ffffff;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                  <div class="column_row" style="font-size: 0;text-align: center;max-width: 624px;margin: 0 auto;">
                    <!--[if (mso)|(IE)]>
                    <table role="presentation" width="624" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
                      <tbody>
                        <tr>
                          <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
                    <![endif]-->
                      <table class="column" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="vertical-align: top;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                        <tbody>
                          <tr>
                            <td class="column_cell px text_dark text_link text_center" style="vertical-align: top;color: #333333;text-align: center;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                              <h5 class="mb_xs" style="color: inherit;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 8px;word-break: break-word;font-size: 19px;line-height: 21px;font-weight: bold;">Hi ${clientName}!</h5>
                              <p style="color: inherit;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;font-size: 16px;line-height: 26px;"></p>
                              <p>Thank you for reaching out!</p>
                              <p>Below is your event quote, along with important details and next steps.</p>
                            </td>
                           
                          </tr>
                        </tbody>
                      </table>
                    <!--[if (mso)|(IE)]>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                    <![endif]-->
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          <!--[if (mso)|(IE)]>
                </td>
              </tr>
            </tbody>
          </table>
          <![endif]-->
        </td>
      </tr>
    </tbody>
  </table>
    
  
  <table class="email_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0">
    <tbody>
      <tr>
        <td class="email_bg bg_light px" style="font-size: 0;text-align: center;line-height: 100%;background-color: #f7f7fa;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
          <!--[if (mso)|(IE)]>
          <table role="presentation" width="800" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
            <tbody>
              <tr>
                <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
          <![endif]-->
          <table class="content_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="max-width: 800px;margin: 0 auto;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
            <tbody>
              <tr>
                <td class="content_cell bg_white px" style="font-size: 0;text-align: center;background-color: #ffffff;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                  <div class="column_row" style="font-size: 0;text-align: center;max-width: 624px;margin: 0 auto;">
                    <!--[if (mso)|(IE)]>
                    <table role="presentation" width="312" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
                      <tbody>
                        <tr>
                          <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
                    <![endif]-->
                      <div class="col_2" style="vertical-align: top;display: inline-block;width: 100%;max-width: 312px;">
                        <table class="column" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="vertical-align: top;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                          <tbody>
                            <tr>
                              <td class="column_cell px py text_dark text_center" style="vertical-align: top;color: #333333;text-align: center;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                                <h3 class="mb_xs" style="color: inherit;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;word-break: break-word;font-size: 18px;padding-top: 10px;font-weight: bold;">Event Details</h3>
                               
                              </td>
                            </tr>
                          </tbody>
                        </table>
                      </div>
                    <!--[if (mso)|(IE)]>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                    <![endif]-->
                  
                </td>
              </tr>
            </tbody>
          </table>
          <!--[if (mso)|(IE)]>
                </td>
              </tr>
            </tbody>
          </table>
          <![endif]-->
        </td>
      </tr>
    </tbody>
  </table>

<table class="email_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0">
    <tbody>
      <tr>
        <td class="email_bg bg_light px" style="font-size: 0;text-align: center;line-height: 100%;background-color: #F7F7FA;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
          <!--[if (mso)|(IE)]>
          <table role="presentation" width="800" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
            <tbody>
              <tr>
                <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
          <![endif]-->
          <table class="content_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="max-width: 800px;margin: 0 auto;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
            <tbody>
              <tr>
                <td class="content_cell bg_white px py" style="font-size: 0;text-align: center;background-color: #ffffff;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                  <div class="column_row" style="font-size: 0;text-align: center;max-width: 624px;margin: 0 auto;">
                    <!--[if (mso)|(IE)]>
                    <table role="presentation" width="416" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
                      <tbody>
                        <tr>
                          <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
                    <![endif]-->
                      <div class="col_3" style="vertical-align: top;display: inline-block;width: 100%;max-width: 516px;">
                        <table class="column" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="vertical-align: top;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                          <tbody>
                            <tr>
                              <td class="column_cell bg_secondary brounded px py_xs text_dark text_left mobile_center" style="vertical-align: top;background-color: #f7f7fa;color: #333333;border-radius: 4px;text-align: left;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                                 <table class="table">
          <tr><td  style="text-align:center; font-size: 14px;"><strong>📅 When</strong></td></tr>
          <tr><td  style=" width:100%; text-align: center; padding-bottom: 12px; font-size: 14px;">${eventDate}</td></tr>
          <tr><td style="text-align:center; font-size: 14px;"><strong>📌 Where</strong></td></tr>
          <tr><td  style="padding-bottom: 12px; text-align:center; font-size: 14px;">${eventLocation}</td></tr>
         
        </table>

         <table style="padding-top: 8px;" class="table">
        
          <tr><td style="font-size: 14px;" colspan="2"><strong>🎉 Occasion</strong></td><td style="padding-bottom: 12px; text-align: right; font-size: 14px;">${occasion}</td></tr>
          <tr><td style="font-size: 14px;" colspan="2"><strong>👥 Guest Count</strong></td><td style="padding-bottom: 12px; text-align: right; font-size: 14px;">${guestCount}</td></tr>
          <tr><td style="font-size: 14px;" colspan="2"><strong>👨‍🍳 Helpers Needed</strong></td><td style="padding-bottom: 12px; font-size: 14px; text-align: right;">${numHelpers}</td></tr>
          <tr><td style="font-size: 14px;" colspan="2"><strong>🕒 For How Long</strong></td><td style="padding-bottom: 12px;text-align: right; font-size: 14px;">${duration} Hours</td></tr>
        </table>
                      </div>
                    <!--[if (mso)|(IE)]>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                    <![endif]-->
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          <!--[if (mso)|(IE)]>
                </td>
              </tr>
            </tbody>
          </table>
          <![endif]-->
        </td>
      </tr>
    </tbody>
</table>

<table class="email_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0">
    <tbody>
      <tr>
        <td class="email_bg bg_light px" style="font-size: 0;text-align: center;line-height: 100%;background-color: #f7f7fa;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
          <!--[if (mso)|(IE)]>
          <table role="presentation" width="800" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
            <tbody>
              <tr>
                <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
          <![endif]-->
          <table class="content_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="max-width: 800px;margin: 0 auto;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
            <tbody>
              <tr>
                <td class="content_cell bg_white px" style="font-size: 0;text-align: center;background-color: #ffffff;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                  <div class="column_row" style="font-size: 0;text-align: center;max-width: 624px;margin: 0 auto;">
                    <!--[if (mso)|(IE)]>
                    <table role="presentation" width="312" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
                      <tbody>
                        <tr>
                          <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
                    <![endif]-->
                      <div class="col_2" style="vertical-align: top;display: inline-block;width: 100%;max-width: 312px;">
                        <table class="column" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="vertical-align: top;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                          <tbody>
                            <tr>
                              <td class="column_cell px py text_dark text_center" style="vertical-align: top;color: #333333;text-align: center;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                               <h3 style="color:inherit;font-family:Arial,Helvetica,sans-serif;margin-top:8px;word-break:break-word;font-size:18px;padding-top:10px;font-weight:bold">Our Rates & Pricing</h3>
                               
                              </td>
                            </tr>
                          </tbody>
                        </table>
                      </div>
                    <!--[if (mso)|(IE)]>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                    <![endif]-->
                  
                </td>
              </tr>
            </tbody>
          </table>
          <!--[if (mso)|(IE)]>
                </td>
              </tr>
            </tbody>
          </table>
          <![endif]-->
        </td>
      </tr>
    </tbody>
  </table>

<table class="email_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0">
    <tbody>
      <tr>
        <td class="email_bg bg_light px" style="font-size: 0;text-align: center;line-height: 100%;background-color: #F7F7FA;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
          <!--[if (mso)|(IE)]>
          <table role="presentation" width="800" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
            <tbody>
              <tr>
                <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
          <![endif]-->
          <table class="content_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="max-width: 800px;margin: 0 auto;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
            <tbody>
              <tr>
                <td class="content_cell bg_white px py" style="font-size: 0;text-align: center;background-color: #ffffff;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                  <div class="column_row" style="font-size: 0;text-align: center;max-width: 624px;margin: 0 auto;">
                    <!--[if (mso)|(IE)]>
                    <table role="presentation" width="416" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
                      <tbody>
                        <tr>
                          <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
                    <![endif]-->
                      <div class="col_3" style="vertical-align: top;display: inline-block;width: 100%;max-width: 516px;">
                        <table class="column" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="vertical-align: top;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                          <tbody>
                           <tr>
                              <td class="column_cell bg_secondary brounded px py_xs text_dark text_left mobile_center" style="vertical-align: top;background-color: #f7f7fa;color: #333333;border-radius: 4px;text-align: left;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                                <table class="table">
        <tr>
            <td style="font-size: 14px;"><strong>💲 ${rateLabel}:</strong></td>
            
          </tr>
          <tr>
            <td style="text-align: right;  width: 100%; font-size: 14px;">$${baseRate} / helper (first 4 hours)</td>
          </tr>
          <tr>
            <td style="font-size: 14px;"><strong>⏳ ${rateLabel} Additional Hours:</strong></td>
            
          </tr>
                    <tr>

            <td style="text-align: right; width: 100%; font-size: 14px;">$${hourlyRate} per additional hour per helper</td>
          </tr>

          <tr>

            <td style="font-size: 14px;"><strong>💰 ${rateLabel} Estimated Total:</strong></td>
            
          </tr>
          <tr>

            
            <td style="text-align: right; font-size: 14px;">$${totalCost}</td>
          </tr>
          
          <tr>
            <td  class="italic-text"  style="font-size: 12px; text-align: center; padding-top: 7px;">Final total may adjust based on our call.</td>
          </tr>
          <tr>
            <td  class="italic-text"  style="font-size: 12px; text-align: center; padding-top: 7px;">Gratuity is not included but always appreciated!</td>
          </tr>
        
        </table>
                             
                              </td>
                            </tr>

                         
                          </tbody>
                        </table>
                      </div>
                    <!--[if (mso)|(IE)]>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                    <![endif]-->
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          <!--[if (mso)|(IE)]>
                </td>
              </tr>
            </tbody>
          </table>
          <![endif]-->
        </td>
      </tr>
    </tbody>
</table>







  <!-- Event Details Section Ends --> 

<table class="email_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0">
    <tbody>
      <tr>
        <td class="email_bg bg_light px" style="font-size: 0;text-align: center;line-height: 100%;background-color: #f7f7fa;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
          <!--[if (mso)|(IE)]>
          <table role="presentation" width="800" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
            <tbody>
              <tr>
                <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
          <![endif]-->
          <table class="content_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="max-width: 800px;margin: 0 auto;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
            <tbody>
              <tr>
                <td class="content_cell bg_white px" style="font-size: 0;text-align: center;background-color: #ffffff;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                  <div class="column_row" style="font-size: 0;text-align: center;max-width: 624px;margin: 0 auto;">
                    <!--[if (mso)|(IE)]>
                    <table role="presentation" width="312" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
                      <tbody>
                        <tr>
                          <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
                    <![endif]-->
                      <div class="col_2" style="vertical-align: top;display: inline-block;width: 100%;max-width: 312px;">
                        <table class="column" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="vertical-align: top;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                          <tbody>
                            <tr>
                              <td class="column_cell px py text_dark text_center" style="vertical-align: top;color: #333333;text-align: center;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                                <h3 style="color:inherit;font-family:Arial,Helvetica,sans-serif;margin-top:8px;word-break:break-word;font-size:18px;padding-top:10px;font-weight:bold">Services Included</h3>
                               
                              </td>
                            </tr>
                          </tbody>
                        </table>
                      </div>
                    <!--[if (mso)|(IE)]>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                    <![endif]-->
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          <!--[if (mso)|(IE)]>
                </td>
              </tr>
            </tbody>
          </table>
          <![endif]-->
        </td>
      </tr>
    </tbody>
  </table>

  <!-- Event Details Section 2 Ends --> 

<table class="email_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0">
    <tbody>
      <tr>
        <td class="email_bg bg_light px" style="font-size: 0;text-align: center;line-height: 100%;background-color: #F7F7FA;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
          <!--[if (mso)|(IE)]>
          <table role="presentation" width="800" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
            <tbody>
              <tr>
                <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
          <![endif]-->
          <table class="content_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="max-width: 800px;margin: 0 auto;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
            <tbody>
              <tr>
                <td class="content_cell bg_white px py" style="font-size: 0;text-align: center;background-color: #ffffff;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                  <div class="column_row" style="font-size: 0;text-align: center;max-width: 624px;margin: 0 auto;">
                    <!--[if (mso)|(IE)]>
                    <table role="presentation" width="416" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
                      <tbody>
                        <tr>
                          <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
                    <![endif]-->
                      <div class="col_3" style="vertical-align: top;display: inline-block;width: 100%;max-width: 516px;">
                        <table class="column" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="vertical-align: top;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                          <tbody>
                            <tr>
                              <td class="column_cell bg_secondary brounded px py_xs text_dark text_left mobile_center" style="vertical-align: top;background-color: #f7f7fa;color: #333333;border-radius: 4px;text-align: left;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                                
                                
                                
                                <table class="column" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="vertical-align: top;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                                  <tbody>
                                    <tr>
                                      <td class="column_cell pl py bb_light text_dark text_left" style="vertical-align: top;color: #333333;border-bottom: 1px solid #dee0e1;text-align: left;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                                        <h5 class="mb_xs" style="color: inherit;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 8px;word-break: break-word;font-size: 16px;line-height: 21px;font-weight: bold;">Setup & Presentation</h5>
                                        <p class="text_xs text_secondary" style="color: #000000;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;font-size: 14px;line-height: 22px;">- Arranging tables, chairs, and decorations</p>
                                        <p class="text_xs text_secondary" style="color: #000000;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;font-size: 14px;line-height: 22px;">- Buffet setup & live buffet service</p>
                                        <p class="text_xs text_secondary" style="color: #000000;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;font-size: 14px;line-height: 22px;">- Butler-passed appetizers & cocktails</p>
                                      </td>
                                    </tr>
                                                                    <tr>
                                      <td class="column_cell pl py bb_light text_dark text_left" style="vertical-align: top;color: #333333;border-bottom: 1px solid #dee0e1;text-align: left;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
  <h5 class="mb_xs" style="color: inherit;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 8px;word-break: break-word;font-size: 16px;line-height: 21px;font-weight: bold;">Dining & Guest Assistance</h5>
  <p class="text_xs text_secondary" style="color: #000000;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;font-size: 14px;line-height: 22px;">- Multi-course plated dinners</p>
  <p class="text_xs text_secondary" style="color: #000000;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;font-size: 14px;line-height: 22px;">- General bussing </p>
  <p class="text_xs text_secondary" style="color: #000000;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;font-size: 13px;line-height: 22px;">(plates, silverware, glassware)</p>
  <p class="text_xs text_secondary" style="color: #000000;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;font-size: 14px;line-height: 22px;">- Beverage service</p>
  <p class="text_xs text_secondary" style="color: #000000;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;font-size: 13px;line-height: 22px;">(water, wine, champagne, coffee, etc.)</p>
   
  <p class="text_xs text_secondary" style="color: #000000;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;font-size: 14px;line-height: 22px;">- Special services</p>
  <p class="text_xs text_secondary" style="color: #000000;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;font-size: 13px;line-height: 22px;"> (cake cutting, dessert plating, etc.)</p>
</td>
</tr>
<tr>
<td class="column_cell pl py bb_light text_dark text_left" style="vertical-align: top;color: #333333;border-bottom: 1px solid #dee0e1;text-align: left;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
  <h5 class="mb_xs" style="color: inherit;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 8px;word-break: break-word;font-size: 16px;line-height: 21px;font-weight: bold;">Cleanup & End-of-Event Support</h5>
  <p class="text_xs text_secondary" style="color: #000000;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;font-size: 14px;line-height: 22px;">- Washing dishes, managing trash, and keeping the event space tidy</p>
  <p class="text_xs text_secondary" style="color: #000000;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;font-size: 14px;line-height: 22px;">- Kitchen cleanup & end-of-event breakdown</p>
  <p class="text_xs text_secondary" style="color: #000000;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;font-size: 14px;line-height: 22px;">- Assisting with food storage & leftovers</p>
</td>                   
                                    </tr>
                                   <tr>
  <td class="column_cell pl py bb_light text_dark text_left" style="vertical-align: top;color: #333333;text-align: left;padding-left: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
    <table style="padding-top: 8px;" class="table">
        
          <tr>
          <td style="padding-bottom: 12px; text-align: center; font-size: 14px;">Need something specific?</td>
          </tr>
 <tr>
          <td style="padding-bottom: 12px; text-align: center; font-size: 14px;"> Let us know!</td>
          </tr>
           <tr>
          <td style="padding-bottom: 12px; text-align: center; font-size: 14px;">We’ll do our best to accommodate your request.</td>
          </tr>
        </table>


  </td>
</tr>


                                  </tbody>
                                </table>
                              </td>
                            </tr>
                          </tbody>
                        </table>
                      </div>
                    <!--[if (mso)|(IE)]>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                    <![endif]-->
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          <!--[if (mso)|(IE)]>
                </td>
              </tr>
            </tbody>
          </table>
          <![endif]-->
        </td>
      </tr>
    </tbody>
  </table>











<table class="email_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0">
    <tbody>
      <tr>
        <td class="email_bg bg_light px" style="font-size: 0;text-align: center;line-height: 100%;background-color: #f7f7fa;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
          <!--[if (mso)|(IE)]>
          <table role="presentation" width="800" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
            <tbody>
              <tr>
                <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
          <![endif]-->
          <table class="content_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="max-width: 800px;margin: 0 auto;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
            <tbody>
              <tr>
                <td class="content_cell bg_white px" style="font-size: 0;text-align: center;background-color: #ffffff;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                  <div class="column_row" style="font-size: 0;text-align: center;max-width: 624px;margin: 0 auto;">
                    <!--[if (mso)|(IE)]>
                    <table role="presentation" width="312" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
                      <tbody>
                        <tr>
                          <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
                    <![endif]-->
                      <div class="col_2" style="vertical-align: top;display: inline-block;width: 100%;max-width: 312px;">
                        <table class="column" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="vertical-align: top;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                          <tbody>
                            <tr>
                              <td class="column_cell px py text_dark text_center" style="vertical-align: top;color: #333333;text-align: center;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                               <h3 style="color:inherit;font-family:Arial,Helvetica,sans-serif;margin-top:14px;word-break:break-word;font-size:18px;padding-top:10px;font-weight:bold">Payment Options</h3>
                               
                              </td>
                            </tr>
                          </tbody>
                        </table>
                      </div>
                    <!--[if (mso)|(IE)]>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                    <![endif]-->
                  
                </td>
              </tr>
            </tbody>
          </table>
          <!--[if (mso)|(IE)]>
                </td>
              </tr>
            </tbody>
          </table>
          <![endif]-->
        </td>
      </tr>
    </tbody>
  </table>

<table class="email_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0">
    <tbody>
      <tr>
        <td class="email_bg bg_light px" style="font-size: 0;text-align: center;line-height: 100%;background-color: #F7F7FA;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
          <!--[if (mso)|(IE)]>
          <table role="presentation" width="800" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
            <tbody>
              <tr>
                <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
          <![endif]-->
          <table class="content_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="max-width: 800px;margin: 0 auto;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
            <tbody>
              <tr>
                <td class="content_cell bg_white px py" style="font-size: 0;text-align: center;background-color: #ffffff;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                  <div class="column_row" style="font-size: 0;text-align: center;max-width: 624px;margin: 0 auto;">
                    <!--[if (mso)|(IE)]>
                    <table role="presentation" width="416" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
                      <tbody>
                        <tr>
                          <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
                    <![endif]-->
                      <div class="col_3" style="vertical-align: top;display: inline-block;width: 100%;max-width: 516px;">
                        <table class="column" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="vertical-align: top;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                          <tbody>
                           <tr>
                              <td class="column_cell bg_secondary brounded px py_xs text_dark text_left mobile_center" style="vertical-align: top;background-color: #f7f7fa;color: #333333;border-radius: 4px;text-align: left;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                                 <table class="table">
      
          <tr>
            
            <td style="text-align:center; font-size: 14px;">Check, Debit / Credit Cards (via Stripe), Venmo, Zelle, Apple Pay</td>
          </tr>
        
        </table>
                             
                              </td>
                            </tr>

                         
                          </tbody>
                        </table>
                      </div>
                    <!--[if (mso)|(IE)]>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                    <![endif]-->
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          <!--[if (mso)|(IE)]>
                </td>
              </tr>
            </tbody>
          </table>
          <![endif]-->
        </td>
      </tr>
    </tbody>
</table>




<table class="email_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0">
    <tbody>
      <tr>
        <td class="email_bg bg_light px" style="font-size: 0;text-align: center;line-height: 100%;background-color: #f7f7fa;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
          <!--[if (mso)|(IE)]>
          <table role="presentation" width="800" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
            <tbody>
              <tr>
                <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
          <![endif]-->
          <table class="content_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="max-width: 800px;margin: 0 auto;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
            <tbody>
              <tr>
                <td class="content_cell bg_white px" style="font-size: 0;text-align: center;background-color: #ffffff;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                  <div class="column_row" style="font-size: 0;text-align: center;max-width: 624px;margin: 0 auto;">
                    <!--[if (mso)|(IE)]>
                    <table role="presentation" width="312" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
                      <tbody>
                        <tr>
                          <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
                    <![endif]-->
                      <div class="col_2" style="vertical-align: top;display: inline-block;width: 100%;max-width: 312px;">
                        <table class="column" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="vertical-align: top;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                          <tbody>
                            <tr>
                              <td class="column_cell px py text_dark text_center" style="vertical-align: top;color: #333333;text-align: center;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                               <h3 style="color:inherit;font-family:Arial,Helvetica,sans-serif;margin-top:14px;word-break:break-word;font-size:18px;padding-top:10px;font-weight:bold">What Happens Next</h3>
                               
                              </td>
                            </tr>
                          </tbody>
                        </table>
                      </div>
                    <!--[if (mso)|(IE)]>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                    <![endif]-->
                  
                </td>
              </tr>
            </tbody>
          </table>
          <!--[if (mso)|(IE)]>
                </td>
              </tr>
            </tbody>
          </table>
          <![endif]-->
        </td>
      </tr>
    </tbody>
  </table>

<table class="email_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0">
    <tbody>
      <tr>
        <td class="email_bg bg_light px" style="font-size: 0;text-align: center;line-height: 100%;background-color: #F7F7FA;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
          <!--[if (mso)|(IE)]>
          <table role="presentation" width="800" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
            <tbody>
              <tr>
                <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
          <![endif]-->
          <table class="content_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="max-width: 800px;margin: 0 auto;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
            <tbody>
              <tr>
                <td class="content_cell bg_white px py" style="font-size: 0;text-align: center;background-color: #ffffff;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                  <div class="column_row" style="font-size: 0;text-align: center;margin: 0 auto;">
                    <!--[if (mso)|(IE)]>
                    <table role="presentation" width="650" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
                      <tbody>
                        <tr>
                          <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
                    <![endif]-->
                      <div class="col_3" style="vertical-align: top;display: inline-block;width: 100%;">
                        <table class="column" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="vertical-align: top;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                          <tbody>
                           <tr>
                              <td class="column_cell bg_secondary brounded px py_xs text_dark text_left mobile_center" style="vertical-align: top;background-color: #f7f7fa;color: #333333;border-radius: 4px;text-align: left;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                               
          <table  class="table">
        
        <tr>
            <td style="padding: 10px; text-align: center; font-weight: bold; font-size: 15px;" >Booked already?</td>
             
          </tr>
          <tr>
            <td style="padding: 10px; text-align: center; font-size: 14px;" > We’ll call you at your scheduled time to go over details.</td>
             
          </tr>
          <tr>
            
             <td style="padding: 10px; text-align: center; font-size: 14px;" >If all looks good after our call, we’ll send a Stripe deposit link to proceed.</td>
          </tr>
          
         
         
         <tr>
            <td colspan="2" style="padding: 10px; text-align: center; font-size: 14px;">Once the deposit is in, your reservation is locked in.</td>
          </tr>

           <tr>
            <td colspan="2" style="padding: 10px; background: white; font-size: 14px; text-align: center;" >Deposit is 40 – 50% of the estimate rounded for simplicity.</td>
          </tr>
          <tr>
            <td colspan="2" style="background: white; text-align: center; font-weight: 400; font-size: 12px;" >🚫 (required to confirm your reservation)</td>
          </tr>

        </table>
                       <table style="margin-top: 24px;" class="table">
        <tr>
            <td colspan="2" style="padding: 10px; font-size: 14px; background: #ffede0; text-align: center; font-weight: bold;">📞 Haven’t scheduled a call yet?</td>
          </tr>
          <tr>
            <td colspan="2" style="padding: 10px; font-size: 14px; background: #ffede0; text-align: center; font-weight: bold;">📅 Book now to get started</td>
          </tr>
           <tr>
            <td colspan="2" style="padding: 10px; background: #ffede0; text-align: center; font-weight: 400; font-size: 12px;">(to confirm helpers, tasks, and setup)</td>
          </tr>



          </table>
          <table style="margin-top: 8px;" class="table">
          <td class="bg_primary" style="font-size: 14px;padding: 10px 20px;border-radius: 4px;line-height: normal;text-align: center;font-weight: bold;">
                    <a href="https://stlpartyhelpers.com/book-appointment" style="text-decoration: none;font-family: Arial, Helvetica, sans-serif;background: white;padding: 12px;border: 1px solid;">
                      Click Here to Schedule Appointment
                    </a>
                  </td></table>        
                              </td>
                            </tr>

                         
                          </tbody>
                        </table>
                      </div>
                    <!--[if (mso)|(IE)]>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                    <![endif]-->
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          <!--[if (mso)|(IE)]>
                </td>
              </tr>
            </tbody>
          </table>
          <![endif]-->
        </td>
      </tr>
    </tbody>
</table>






  <!-- Event Details Section Ends --> 



  <!-- Event Details Section 2 Ends --> 


        <table class="column" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="vertical-align: top;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                        <tbody>
                          <tr>
                            <td class="column_cell px py_md text_secondary text_center" style="vertical-align: top;color: #959ba0;text-align: center;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                              <p class="img_inline mb_md" style="display: none; color: inherit;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 16px;word-break: break-word;font-size: 16px;line-height: 100%;clear: both;">
                                <a href="#" style="-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;text-decoration: none;color: #959ba0;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;"><img src="images/facebook.png" width="24" height="24" alt="Facebook" style="max-width: 24px;-ms-interpolation-mode: bicubic;border: 0;height: auto;line-height: 100%;outline: none;text-decoration: none;"></a> &nbsp;&nbsp;
                                <a href="#" style="-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;text-decoration: none;color: #959ba0;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;"><img src="images/twitter.png" width="24" height="24" alt="Twitter" style="max-width: 24px;-ms-interpolation-mode: bicubic;border: 0;height: auto;line-height: 100%;outline: none;text-decoration: none;"></a> &nbsp;&nbsp;
                                <a href="#" style="-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;text-decoration: none;color: #959ba0;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;"><img src="images/instagram.png" width="24" height="24" alt="Instagram" style="max-width: 24px;-ms-interpolation-mode: bicubic;border: 0;height: auto;line-height: 100%;outline: none;text-decoration: none;"></a> &nbsp;&nbsp;
                                <a href="#" style="-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;text-decoration: none;color: #959ba0;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;"><img src="images/pinterest.png" width="24" height="24" alt="Pinterest" style="max-width: 24px;-ms-interpolation-mode: bicubic;border: 0;height: auto;line-height: 100%;outline: none;text-decoration: none;"></a>
                              </p>
                              <p class="mb" style="color: inherit;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 16px;word-break: break-word;font-size: 15px;line-height: 26px;">
                              	© 2025 STL Party Helpers<br> <a href="https://stlpartyhelpers.com">stlpartyhelpers.com</a><br/>
                              	<span class="text_adr" href="#" style="-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;text-decoration: none;color: #959ba0;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;">
                                <span style="color: #959ba0;text-decoration: none;">4220 Duncan Ave., Ste. 201</span><br/>
                                <span style="color: #959ba0;text-decoration: none;">St. Louis, MO 63110</span>
                                </span>
                                <br/>
                              <a href="tel:13147145514" style="background-color:#ffffff;color: #262578;padding:8px 8px;text-decoration:none;border-radius:4px;font-family:Arial,sans-serif;display:inline-block;margin-top:10px;border: 3px solid #262578;font-weight: bold;" target="_blank">
 Tap to Call Us: (314) 714-5514
</a>

                              </p>
                            <span id="stlth-version" style="color: #959ba0;text-decoration: none; font-size: 8px;">v 1.00</span>  
                            </td>
                            
                          </tr>
                        </tbody>
                      </table>                             
     
                             
                              </td>
                            </tr>
                          </tbody>
                        </table>
                      </div>
                    <!--[if (mso)|(IE)]>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                    <![endif]-->
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          <!--[if (mso)|(IE)]>
                </td>
              </tr>
            </tbody>
          </table>
          <![endif]-->
        </td>
      </tr>
    </tbody>
</table>



   
    </div>
  </div>
</body>

  `;
}

function generateBookingDepositEmail(
  subject,
  referenceNumber,
  clientName,
  eventDate,
  eventTime,
  eventLocation,
  numHelpers,
  duration,
  totalCost,
  occasion,
  guestCount
) {
  return `<head>
  <meta charset="UTF-8"> <!-- ✅ Forces UTF-8 encoding -->
  <style>
    body {
      font-family: Arial, sans-serif;
      font-size: 14px;
    }
    
    .container {
      max-width: 700px;
      margin: auto;
      padding: 0px;
      background-color: #f9f9f9;
      border-radius: 8px;
      box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
    }

    .header {
      background-color: #262578;
      padding: 5px;
      border-radius: 8px 8px 0 0;
      color: white;
      text-align: center;
    }

    .header h3 {
      margin: 0;
      display: inline-block;
      background-color: white;
      color: black;
      padding: 2px;
    }
body, p, .table td {
  font-size: 14px;
  line-height: 1.2;
}
    .content {
      padding: 5px;
      background-color: white;
      border-radius: 0 0 8px 8px;
    }

    .section-title {
      color: #333;
      padding: 0px;
      padding-bottom: 8px;
      margin: 0px;
      padding-top: 12px;
    }

    .section {
      
      margin: 0px 0;
      padding-top: 0px;
      padding-bottom: 10px;
    }
tr {
  padding-bottom: 7px;
}
    .table {
      width: 100%;
      font-family: Arial, sans-serif;
      line-height: 1.5;
      border-collapse: collapse;
      padding-bottom: 20px;
    }

    .table td {
      padding: 5px;
    }

    .highlight-box {
      display: inline-block;
      background-color: #f8ff94;
      padding: 10px;
      border-radius: 6px;
      font-size: 15px;
      font-weight: bold;
      border: 1px solid #d4c200;
    }

    .button {
      color: #ffffff;
      text-decoration: none;
      display: inline-block;
      background-color: #673AB7;
      padding: 5px;
      border-radius: 6px;
      font-size: 15px;
      font-weight: bold;
    }

    .footer {
      margin-top: 5px;
      padding: 5px;
      
      color: white;
      text-align: center;
      border-radius: 8px;
    }

    .footer a {
      color: white;
      text-decoration: none;
    }

    .italic-text {
      color: #555;
      font-style: italic;
      font-size: 12px;
    }
    .list p{
      font-size: 13px;
      line-height: 1;
    }

  </style>
</head>
<body>
  <div class="container">
    <table class="email_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0">
    <tbody>
      <tr>
        <td class="email_bg bg_light px pt_md" style="font-size: 0;text-align: center;line-height: 100%;background-color: #F7F7FA;padding-left: 8px;padding-right: 8px;padding-top: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
          <!--[if (mso)|(IE)]>
          <table role="presentation" width="800" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
            <tbody>
              <tr>
                <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
          <![endif]-->
          <table class="content_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="max-width: 800px;margin: 0 auto;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
            <tbody>
              <tr>
                <td class="content_cell bg_white brounded_top bt_primary px py" style="font-size: 0;text-align: center;background-color: #ffffff;border-top: 4px solid #262578;border-radius: 4px 4px 0 0;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                  <div class="column_row" style="font-size: 0;text-align: center;max-width: 624px;margin: 0 auto;">
                    <!--[if (mso)|(IE)]>
                    <table role="presentation" width="624" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
                      <tbody>
                        <tr>
                          <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
                    <![endif]-->
                      <table class="column" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="vertical-align: top;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                        <tbody>
                          <tr>
                            <td class="column_cell px py_xs text_primary text_center" style="vertical-align: top;color: #2376dc;text-align: center;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                              <p class="img_inline" style="color: inherit;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;font-size: 16px;line-height: 100%;clear: both;"><a href="#" style="-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;text-decoration: none;color: #2376dc;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;"><img src="https://stlpartyhelpers.com/wp-content/uploads/2024/05/FullLogo_Transparent_NoBuffer-1-e1716631491571.png" width="110" height="" alt="STL Party Helpers is always here for you!" style="max-width: 110px;-ms-interpolation-mode: bicubic;border: 0;height: auto;line-height: 100%;outline: none;text-decoration: none; display: block; margin: 0 auto;"></a></p>
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    <!--[if (mso)|(IE)]>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                    <![endif]-->
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          <!--[if (mso)|(IE)]>
                </td>
              </tr>
            </tbody>
          </table>
          <![endif]-->
        </td>
      </tr>
    </tbody>
  </table>
  <table class="email_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0">
    <tbody>
      <tr>
        <td class="email_bg bg_light px" style="font-size: 0;text-align: center;line-height: 100%;background-color: #F7F7FA;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
          <!--[if (mso)|(IE)]>
          <table role="presentation" width="800" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
            <tbody>
              <tr>
                <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
          <![endif]-->
          <table class="content_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="max-width: 800px;margin: 0 auto;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
            <tbody>
              <tr>
                <td class="content_cell bg_white px py_md" style="font-size: 0;text-align: center;background-color: #ffffff;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                  <div class="column_row" style="font-size: 0;text-align: center;max-width: 624px;margin: 0 auto;">
                    <!--[if (mso)|(IE)]>
                    <table role="presentation" width="624" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
                      <tbody>
                        <tr>
                          <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
                    <![endif]-->
                      <table class="column" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="vertical-align: top;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                        <tbody>
                          <tr>
                            <td class="column_cell px text_dark text_link text_center" style="vertical-align: top;color: #333333;text-align: center;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                              <h5 class="mb_xs" style="color: inherit;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 8px;word-break: break-word;font-size: 19px;line-height: 21px;font-weight: bold;">Hi ${clientName}!</h5>
                              <p style="color: inherit;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;font-size: 16px;line-height: 26px;"></p>
                              <p>Thank you for reaching out!</p>
                              <p>Below is your event quote, along with important details and next steps.</p>
                            </td>
                           
                          </tr>
                        </tbody>
                      </table>
                    <!--[if (mso)|(IE)]>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                    <![endif]-->
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          <!--[if (mso)|(IE)]>
                </td>
              </tr>
            </tbody>
          </table>
          <![endif]-->
        </td>
      </tr>
    </tbody>
  </table>
    
  
  <table class="email_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0">
    <tbody>
      <tr>
        <td class="email_bg bg_light px" style="font-size: 0;text-align: center;line-height: 100%;background-color: #f7f7fa;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
          <!--[if (mso)|(IE)]>
          <table role="presentation" width="800" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
            <tbody>
              <tr>
                <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
          <![endif]-->
          <table class="content_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="max-width: 800px;margin: 0 auto;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
            <tbody>
              <tr>
                <td class="content_cell bg_white px" style="font-size: 0;text-align: center;background-color: #ffffff;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                  <div class="column_row" style="font-size: 0;text-align: center;max-width: 624px;margin: 0 auto;">
                    <!--[if (mso)|(IE)]>
                    <table role="presentation" width="312" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
                      <tbody>
                        <tr>
                          <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
                    <![endif]-->
                      <div class="col_2" style="vertical-align: top;display: inline-block;width: 100%;max-width: 312px;">
                        <table class="column" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="vertical-align: top;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                          <tbody>
                            <tr>
                              <td class="column_cell px py text_dark text_center" style="vertical-align: top;color: #333333;text-align: center;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                                <h3 class="mb_xs" style="color: inherit;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;word-break: break-word;font-size: 18px;padding-top: 10px;font-weight: bold;">Event Details (BDE)</h3>
                               
                              </td>
                            </tr>
                          </tbody>
                        </table>
                      </div>
                    <!--[if (mso)|(IE)]>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                    <![endif]-->
                  
                </td>
              </tr>
            </tbody>
          </table>
          <!--[if (mso)|(IE)]>
                </td>
              </tr>
            </tbody>
          </table>
          <![endif]-->
        </td>
      </tr>
    </tbody>
  </table>

<table class="email_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0">
    <tbody>
      <tr>
        <td class="email_bg bg_light px" style="font-size: 0;text-align: center;line-height: 100%;background-color: #F7F7FA;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
          <!--[if (mso)|(IE)]>
          <table role="presentation" width="800" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
            <tbody>
              <tr>
                <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
          <![endif]-->
          <table class="content_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="max-width: 800px;margin: 0 auto;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
            <tbody>
              <tr>
                <td class="content_cell bg_white px py" style="font-size: 0;text-align: center;background-color: #ffffff;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                  <div class="column_row" style="font-size: 0;text-align: center;max-width: 624px;margin: 0 auto;">
                    <!--[if (mso)|(IE)]>
                    <table role="presentation" width="416" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
                      <tbody>
                        <tr>
                          <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
                    <![endif]-->
                      <div class="col_3" style="vertical-align: top;display: inline-block;width: 100%;max-width: 516px;">
                        <table class="column" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="vertical-align: top;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                          <tbody>
                            <tr>
                              <td class="column_cell bg_secondary brounded px py_xs text_dark text_left mobile_center" style="vertical-align: top;background-color: #f7f7fa;color: #333333;border-radius: 4px;text-align: left;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                                 <table class="table">
          <tr><td  style="text-align:center; font-size: 14px;"><strong>📅 When</strong></td></tr>
          <tr><td  style=" width:100%; text-align: center; padding-bottom: 12px; font-size: 14px;">${eventDate}</td></tr>
          <tr><td style="text-align:center; font-size: 14px;"><strong>📌 Where</strong></td></tr>
          <tr><td style="padding-bottom: 12px; text-align:center; font-size: 14px;">${eventLocation}</td></tr>
         
        </table>

         <table style="padding-top: 8px;" class="table">
        
          <tr><td style="font-size: 14px;" colspan="2"><strong>🎉 Occasion</strong></td><td style="padding-bottom: 12px; text-align: right; font-size: 14px;">${occasion}</td></tr>
          <tr><td style="font-size: 14px;" colspan="2"><strong>👥 Guest Count</strong></td><td style="padding-bottom: 12px; text-align: right; font-size: 14px;">${guestCount}</td></tr>
          <tr><td style="font-size: 14px;" colspan="2"><strong>👨‍🍳 Helpers Needed</strong></td><td style="padding-bottom: 12px; font-size: 14px; text-align: right;">${numHelpers}</td></tr>
          <tr><td style="font-size: 14px;" colspan="2"><strong>🕒 For How Long</strong></td><td style="padding-bottom: 12px;text-align: right; font-size: 14px;">${duration} Hours</td></tr>
        </table>
                             
                              </td>
                            </tr>
                          </tbody>
                        </table>
                      </div>
                    <!--[if (mso)|(IE)]>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                    <![endif]-->
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          <!--[if (mso)|(IE)]>
                </td>
              </tr>
            </tbody>
          </table>
          <![endif]-->
        </td>
      </tr>
    </tbody>
</table>

















<table class="email_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0">
    <tbody>
      <tr>
        <td class="email_bg bg_light px" style="font-size: 0;text-align: center;line-height: 100%;background-color: #f7f7fa;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
          <!--[if (mso)|(IE)]>
          <table role="presentation" width="800" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
            <tbody>
              <tr>
                <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
          <![endif]-->
          <table class="content_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="max-width: 800px;margin: 0 auto;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
            <tbody>
              <tr>
                <td class="content_cell bg_white px" style="font-size: 0;text-align: center;background-color: #ffffff;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                  <div class="column_row" style="font-size: 0;text-align: center;max-width: 624px;margin: 0 auto;">
                    <!--[if (mso)|(IE)]>
                    <table role="presentation" width="312" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
                      <tbody>
                        <tr>
                          <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
                    <![endif]-->
                      <div class="col_2" style="vertical-align: top;display: inline-block;width: 100%;max-width: 312px;">
                        <table class="column" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="vertical-align: top;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                          <tbody>
                            <tr>
                              <td class="column_cell px py text_dark text_center" style="vertical-align: top;color: #333333;text-align: center;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                               <h3 style="color:inherit;font-family:Arial,Helvetica,sans-serif;margin-top:8px;word-break:break-word;font-size:18px;padding-top:10px;font-weight:bold">Our Rates & Pricing</h3>
                               
                              </td>
                            </tr>
                          </tbody>
                        </table>
                      </div>
                    <!--[if (mso)|(IE)]>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                    <![endif]-->
                  
                </td>
              </tr>
            </tbody>
          </table>
          <!--[if (mso)|(IE)]>
                </td>
              </tr>
            </tbody>
          </table>
          <![endif]-->
        </td>
      </tr>
    </tbody>
  </table>

<table class="email_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0">
    <tbody>
      <tr>
        <td class="email_bg bg_light px" style="font-size: 0;text-align: center;line-height: 100%;background-color: #F7F7FA;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
          <!--[if (mso)|(IE)]>
          <table role="presentation" width="800" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
            <tbody>
              <tr>
                <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
          <![endif]-->
          <table class="content_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="max-width: 800px;margin: 0 auto;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
            <tbody>
              <tr>
                <td class="content_cell bg_white px py" style="font-size: 0;text-align: center;background-color: #ffffff;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                  <div class="column_row" style="font-size: 0;text-align: center;max-width: 624px;margin: 0 auto;">
                    <!--[if (mso)|(IE)]>
                    <table role="presentation" width="416" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
                      <tbody>
                        <tr>
                          <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
                    <![endif]-->
                      <div class="col_3" style="vertical-align: top;display: inline-block;width: 100%;max-width: 516px;">
                        <table class="column" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="vertical-align: top;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                          <tbody>
                           <tr>
                              <td class="column_cell bg_secondary brounded px py_xs text_dark text_left mobile_center" style="vertical-align: top;background-color: #f7f7fa;color: #333333;border-radius: 4px;text-align: left;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                                 <table class="table">
        <tr>
            <td style="font-size: 14px;"><strong>💲 Base Rate:</strong></td>
            
          </tr>
          <tr>
            <td style="text-align: right;  width: 100%; font-size: 14px;">$200 / helper (first 4 hours)</td>
          </tr>
          <tr>
            <td style="font-size: 14px;"><strong>⏳ Additional Hours:</strong></td>
            
          </tr>
                    <tr>

            <td style="text-align: right; width: 100%; font-size: 14px;">$45 per additional hour per helper</td>
          </tr>

          <tr>

            <td style="font-size: 14px;"><strong>💰 Estimated Total:</strong></td>
            
          </tr>
          <tr>

            
            <td style="text-align: right; font-size: 14px;">$${totalCost}</td>
          </tr>
          
          <tr>
            <td  class="italic-text"  style="font-size: 12px; text-align: center; padding-top: 7px;">Final total may adjust based on our call.</td>
          </tr>
          <tr>
            <td  class="italic-text"  style="font-size: 12px; text-align: center; padding-top: 7px;">Gratuity is not included but always appreciated!</td>
          </tr>
        
        </table>
                             
                              </td>
                            </tr>

                         
                          </tbody>
                        </table>
                      </div>
                    <!--[if (mso)|(IE)]>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                    <![endif]-->
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          <!--[if (mso)|(IE)]>
                </td>
              </tr>
            </tbody>
          </table>
          <![endif]-->
        </td>
      </tr>
    </tbody>
</table>







  <!-- Event Details Section Ends --> 

<table class="email_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0">
    <tbody>
      <tr>
        <td class="email_bg bg_light px" style="font-size: 0;text-align: center;line-height: 100%;background-color: #f7f7fa;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
          <!--[if (mso)|(IE)]>
          <table role="presentation" width="800" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
            <tbody>
              <tr>
                <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
          <![endif]-->
          <table class="content_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="max-width: 800px;margin: 0 auto;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
            <tbody>
              <tr>
                <td class="content_cell bg_white px" style="font-size: 0;text-align: center;background-color: #ffffff;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                  <div class="column_row" style="font-size: 0;text-align: center;max-width: 624px;margin: 0 auto;">
                    <!--[if (mso)|(IE)]>
                    <table role="presentation" width="312" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
                      <tbody>
                        <tr>
                          <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
                    <![endif]-->
                      <div class="col_2" style="vertical-align: top;display: inline-block;width: 100%;max-width: 312px;">
                        <table class="column" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="vertical-align: top;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                          <tbody>
                            <tr>
                              <td class="column_cell px py text_dark text_center" style="vertical-align: top;color: #333333;text-align: center;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                                <h3 style="color:inherit;font-family:Arial,Helvetica,sans-serif;margin-top:8px;word-break:break-word;font-size:18px;padding-top:10px;font-weight:bold">Services Included</h3>
                               
                              </td>
                            </tr>
                          </tbody>
                        </table>
                      </div>
                    <!--[if (mso)|(IE)]>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                    <![endif]-->
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          <!--[if (mso)|(IE)]>
                </td>
              </tr>
            </tbody>
          </table>
          <![endif]-->
        </td>
      </tr>
    </tbody>
  </table>

  <!-- Event Details Section 2 Ends --> 

<table class="email_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0">
    <tbody>
      <tr>
        <td class="email_bg bg_light px" style="font-size: 0;text-align: center;line-height: 100%;background-color: #F7F7FA;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
          <!--[if (mso)|(IE)]>
          <table role="presentation" width="800" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
            <tbody>
              <tr>
                <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
          <![endif]-->
          <table class="content_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="max-width: 800px;margin: 0 auto;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
            <tbody>
              <tr>
                <td class="content_cell bg_white px py" style="font-size: 0;text-align: center;background-color: #ffffff;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                  <div class="column_row" style="font-size: 0;text-align: center;max-width: 624px;margin: 0 auto;">
                    <!--[if (mso)|(IE)]>
                    <table role="presentation" width="416" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
                      <tbody>
                        <tr>
                          <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
                    <![endif]-->
                      <div class="col_3" style="vertical-align: top;display: inline-block;width: 100%;max-width: 516px;">
                        <table class="column" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="vertical-align: top;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                          <tbody>
                            <tr>
                              <td class="column_cell bg_secondary brounded px py_xs text_dark text_left mobile_center" style="vertical-align: top;background-color: #f7f7fa;color: #333333;border-radius: 4px;text-align: left;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                                
                                
                                
                                <table class="column" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="vertical-align: top;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                                  <tbody>
                                    <tr>
                                      <td class="column_cell pl py bb_light text_dark text_left" style="vertical-align: top;color: #333333;border-bottom: 1px solid #dee0e1;text-align: left;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                                        <h5 class="mb_xs" style="color: inherit;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 8px;word-break: break-word;font-size: 16px;line-height: 21px;font-weight: bold;">Setup & Presentation</h5>
                                        <p class="text_xs text_secondary" style="color: #000000;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;font-size: 14px;line-height: 22px;">- Arranging tables, chairs, and decorations</p>
                                        <p class="text_xs text_secondary" style="color: #000000;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;font-size: 14px;line-height: 22px;">- Buffet setup & live buffet service</p>
                                        <p class="text_xs text_secondary" style="color: #000000;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;font-size: 14px;line-height: 22px;">- Butler-passed appetizers & cocktails</p>
                                      </td>
                                    </tr>
                                                                    <tr>
                                      <td class="column_cell pl py bb_light text_dark text_left" style="vertical-align: top;color: #333333;border-bottom: 1px solid #dee0e1;text-align: left;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
  <h5 class="mb_xs" style="color: inherit;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 8px;word-break: break-word;font-size: 16px;line-height: 21px;font-weight: bold;">Dining & Guest Assistance</h5>
  <p class="text_xs text_secondary" style="color: #000000;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;font-size: 14px;line-height: 22px;">- Multi-course plated dinners</p>
  <p class="text_xs text_secondary" style="color: #000000;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;font-size: 14px;line-height: 22px;">- General bussing </p>
  <p class="text_xs text_secondary" style="color: #000000;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;font-size: 13px;line-height: 22px;">(plates, silverware, glassware)</p>
  <p class="text_xs text_secondary" style="color: #000000;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;font-size: 14px;line-height: 22px;">- Beverage service</p>
  <p class="text_xs text_secondary" style="color: #000000;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;font-size: 13px;line-height: 22px;">(water, wine, champagne, coffee, etc.)</p>
   
  <p class="text_xs text_secondary" style="color: #000000;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;font-size: 14px;line-height: 22px;">- Special services</p>
  <p class="text_xs text_secondary" style="color: #000000;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;font-size: 13px;line-height: 22px;"> (cake cutting, dessert plating, etc.)</p>
</td>
</tr>
<tr>
<td class="column_cell pl py bb_light text_dark text_left" style="vertical-align: top;color: #333333;border-bottom: 1px solid #dee0e1;text-align: left;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
  <h5 class="mb_xs" style="color: inherit;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 8px;word-break: break-word;font-size: 16px;line-height: 21px;font-weight: bold;">Cleanup & End-of-Event Support</h5>
  <p class="text_xs text_secondary" style="color: #000000;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;font-size: 14px;line-height: 22px;">- Washing dishes, managing trash, and keeping the event space tidy</p>
  <p class="text_xs text_secondary" style="color: #000000;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;font-size: 14px;line-height: 22px;">- Kitchen cleanup & end-of-event breakdown</p>
  <p class="text_xs text_secondary" style="color: #000000;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;font-size: 14px;line-height: 22px;">- Assisting with food storage & leftovers</p>
</td>                   
                                    </tr>
                                   <tr>
  <td class="column_cell pl py bb_light text_dark text_left" style="vertical-align: top;color: #333333;text-align: left;padding-left: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
     <table style="padding-top: 8px;" class="table">
        
          <tr>
          <td style="padding-bottom: 12px; text-align: center; font-size: 14px;">Need something specific?</td>
          </tr>
 <tr>
          <td style="padding-bottom: 12px; text-align: center; font-size: 14px;"> Let us know!</td>
          </tr>
           <tr>
          <td style="padding-bottom: 12px; text-align: center; font-size: 14px;">We’ll do our best to accommodate your request.</td>
          </tr>
        </table>

  </td>
</tr>


                                  </tbody>
                                </table>
                              </td>
                            </tr>
                          </tbody>
                        </table>
                      </div>
                    <!--[if (mso)|(IE)]>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                    <![endif]-->
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          <!--[if (mso)|(IE)]>
                </td>
              </tr>
            </tbody>
          </table>
          <![endif]-->
        </td>
      </tr>
    </tbody>
  </table>











<table class="email_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0">
    <tbody>
      <tr>
        <td class="email_bg bg_light px" style="font-size: 0;text-align: center;line-height: 100%;background-color: #f7f7fa;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
          <!--[if (mso)|(IE)]>
          <table role="presentation" width="800" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
            <tbody>
              <tr>
                <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
          <![endif]-->
          <table class="content_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="max-width: 800px;margin: 0 auto;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
            <tbody>
              <tr>
                <td class="content_cell bg_white px" style="font-size: 0;text-align: center;background-color: #ffffff;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                  <div class="column_row" style="font-size: 0;text-align: center;max-width: 624px;margin: 0 auto;">
                    <!--[if (mso)|(IE)]>
                    <table role="presentation" width="312" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
                      <tbody>
                        <tr>
                          <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
                    <![endif]-->
                      <div class="col_2" style="vertical-align: top;display: inline-block;width: 100%;max-width: 312px;">
                        <table class="column" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="vertical-align: top;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                          <tbody>
                            <tr>
                              <td class="column_cell px py text_dark text_center" style="vertical-align: top;color: #333333;text-align: center;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                               <h3 style="color:inherit;font-family:Arial,Helvetica,sans-serif;margin-top:14px;word-break:break-word;font-size:18px;padding-top:10px;font-weight:bold">Payment Options</h3>
                               
                              </td>
                            </tr>
                          </tbody>
                        </table>
                      </div>
                    <!--[if (mso)|(IE)]>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                    <![endif]-->
                  
                </td>
              </tr>
            </tbody>
          </table>
          <!--[if (mso)|(IE)]>
                </td>
              </tr>
            </tbody>
          </table>
          <![endif]-->
        </td>
      </tr>
    </tbody>
  </table>

<table class="email_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0">
    <tbody>
      <tr>
        <td class="email_bg bg_light px" style="font-size: 0;text-align: center;line-height: 100%;background-color: #F7F7FA;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
          <!--[if (mso)|(IE)]>
          <table role="presentation" width="800" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
            <tbody>
              <tr>
                <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
          <![endif]-->
          <table class="content_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="max-width: 800px;margin: 0 auto;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
            <tbody>
              <tr>
                <td class="content_cell bg_white px py" style="font-size: 0;text-align: center;background-color: #ffffff;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                  <div class="column_row" style="font-size: 0;text-align: center;max-width: 624px;margin: 0 auto;">
                    <!--[if (mso)|(IE)]>
                    <table role="presentation" width="416" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
                      <tbody>
                        <tr>
                          <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
                    <![endif]-->
                      <div class="col_3" style="vertical-align: top;display: inline-block;width: 100%;max-width: 516px;">
                        <table class="column" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="vertical-align: top;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                          <tbody>
                           <tr>
                              <td class="column_cell bg_secondary brounded px py_xs text_dark text_left mobile_center" style="vertical-align: top;background-color: #f7f7fa;color: #333333;border-radius: 4px;text-align: left;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                                 <table class="table">
      
          <tr>
            
            <td style="text-align:center; font-size: 14px;">Check, Debit / Credit Cards (via Stripe), Venmo, Zelle, Apple Pay</td>
          </tr>
        
        </table>
                             
                              </td>
                            </tr>

                         
                          </tbody>
                        </table>
                      </div>
                    <!--[if (mso)|(IE)]>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                    <![endif]-->
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          <!--[if (mso)|(IE)]>
                </td>
              </tr>
            </tbody>
          </table>
          <![endif]-->
        </td>
      </tr>
    </tbody>
</table>




<table class="email_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0">
    <tbody>
      <tr>
        <td class="email_bg bg_light px" style="font-size: 0;text-align: center;line-height: 100%;background-color: #f7f7fa;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
          <!--[if (mso)|(IE)]>
          <table role="presentation" width="800" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
            <tbody>
              <tr>
                <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
          <![endif]-->
          <table class="content_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="max-width: 800px;margin: 0 auto;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
            <tbody>
              <tr>
                <td class="content_cell bg_white px" style="font-size: 0;text-align: center;background-color: #ffffff;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                  <div class="column_row" style="font-size: 0;text-align: center;max-width: 624px;margin: 0 auto;">
                    <!--[if (mso)|(IE)]>
                    <table role="presentation" width="312" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
                      <tbody>
                        <tr>
                          <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
                    <![endif]-->
                      <div class="col_2" style="vertical-align: top;display: inline-block;width: 100%;max-width: 312px;">
                        <table class="column" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="vertical-align: top;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                          <tbody>
                            <tr>
                              <td class="column_cell px py text_dark text_center" style="vertical-align: top;color: #333333;text-align: center;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                               <h3 style="color:inherit;font-family:Arial,Helvetica,sans-serif;margin-top:14px;word-break:break-word;font-size:18px;padding-top:10px;font-weight:bold">What Happens Next</h3>
                               
                              </td>
                            </tr>
                          </tbody>
                        </table>
                      </div>
                    <!--[if (mso)|(IE)]>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                    <![endif]-->
                  
                </td>
              </tr>
            </tbody>
          </table>
          <!--[if (mso)|(IE)]>
                </td>
              </tr>
            </tbody>
          </table>
          <![endif]-->
        </td>
      </tr>
    </tbody>
  </table>

<table class="email_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0">
    <tbody>
      <tr>
        <td class="email_bg bg_light px" style="font-size: 0;text-align: center;line-height: 100%;background-color: #F7F7FA;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
          <!--[if (mso)|(IE)]>
          <table role="presentation" width="800" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
            <tbody>
              <tr>
                <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
          <![endif]-->
          <table class="content_section" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="max-width: 800px;margin: 0 auto;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
            <tbody>
              <tr>
                <td class="content_cell bg_white px py" style="font-size: 0;text-align: center;background-color: #ffffff;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                  <div class="column_row" style="font-size: 0;text-align: center;margin: 0 auto;">
                    <!--[if (mso)|(IE)]>
                    <table role="presentation" width="650" border="0" cellspacing="0" cellpadding="0" align="center" style="vertical-align:top;Margin:0 auto;">
                      <tbody>
                        <tr>
                          <td align="center" style="line-height:0px;font-size:0px;mso-line-height-rule:exactly;vertical-align:top;">
                    <![endif]-->
                      <div class="col_3" style="vertical-align: top;display: inline-block;width: 100%;">
                        <table class="column" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="vertical-align: top;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                          <tbody>
                           <tr>
                              <td class="column_cell bg_secondary brounded px py_xs text_dark text_left mobile_center" style="vertical-align: top;background-color: #f7f7fa;color: #333333;border-radius: 4px;text-align: left;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                               
          <table  class="table">
        
        <tr>
            <td style="padding: 10px; text-align: center; font-weight: bold; font-size: 15px;" >Booked already?</td>
             
          </tr>
          <tr>
            <td style="padding: 10px; text-align: center; font-size: 14px;" > We’ll call you at your scheduled time to go over details.</td>
             
          </tr>
          <tr>
            
             <td style="padding: 10px; text-align: center; font-size: 14px;" >If all looks good after our call, we’ll send a Stripe deposit link to proceed.</td>
          </tr>
          
         
         
         <tr>
            <td colspan="2" style="padding: 10px; text-align: center; font-size: 14px;">Once the deposit is in, your reservation is locked in.</td>
          </tr>

           <tr>
            <td colspan="2" style="padding: 10px; font-size: 14px; background: white; text-align: center;" >Deposit is 40 – 50% of the estimate rounded for simplicity.</td>
          </tr>
          <tr>
            <td colspan="2" style="background: white; text-align: center; font-weight: 400; font-size: 12px;" >🚫 (required to confirm your reservation)</td>
          </tr>

        </table>
                       <table style="margin-top: 24px;" class="table">
        <tr>
            <td colspan="2" style="padding: 10px; background: #ffede0; text-align: center; font-size: 14px; font-weight: bold;">📞 Haven’t scheduled a call yet?</td>
          </tr>
          <tr>
            <td colspan="2" style="padding: 10px; background: #ffede0; font-size: 14px; text-align: center; font-weight: bold;">📅 Book now to get started</td>
          </tr>
           <tr>
            <td colspan="2" style="padding: 10px; background: #ffede0; text-align: center; font-weight: 400; font-size: 12px;">(to confirm helpers, tasks, and setup)</td>
          </tr>



          </table>
          <table style="margin-top: 8px;" class="table">
          <td class="bg_primary" style="font-size: 14px;padding: 10px 20px;border-radius: 4px;line-height: normal;text-align: center;font-weight: bold;">
                   <a href="https://stlpartyhelpers.com/book-appointment" style="text-decoration: none;font-family: Arial, Helvetica, sans-serif;background: white;padding: 12px;border: 1px solid;">
                      Click Here to Schedule Appointment
                    </a>
                  </td></table>        
                              </td>
                            </tr>

                         
                          </tbody>
                        </table>
                      </div>
                    <!--[if (mso)|(IE)]>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                    <![endif]-->
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          <!--[if (mso)|(IE)]>
                </td>
              </tr>
            </tbody>
          </table>
          <![endif]-->
        </td>
      </tr>
    </tbody>
</table>






  <!-- Event Details Section Ends --> 



  <!-- Event Details Section 2 Ends --> 


        <table class="column" role="presentation" align="center" width="100%" cellspacing="0" cellpadding="0" border="0" style="vertical-align: top;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                        <tbody>
                          <tr>
                            <td class="column_cell px py_md text_secondary text_center" style="vertical-align: top;color: #959ba0;text-align: center;padding-top: 8px;padding-bottom: 8px;padding-left: 8px;padding-right: 8px;-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;mso-table-lspace: 0pt;mso-table-rspace: 0pt;">
                              <p class="img_inline mb_md" style="display: none; color: inherit;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 16px;word-break: break-word;font-size: 16px;line-height: 100%;clear: both;">
                                <a href="#" style="-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;text-decoration: none;color: #959ba0;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;"><img src="images/facebook.png" width="24" height="24" alt="Facebook" style="max-width: 24px;-ms-interpolation-mode: bicubic;border: 0;height: auto;line-height: 100%;outline: none;text-decoration: none;"></a> &nbsp;&nbsp;
                                <a href="#" style="-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;text-decoration: none;color: #959ba0;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;"><img src="images/twitter.png" width="24" height="24" alt="Twitter" style="max-width: 24px;-ms-interpolation-mode: bicubic;border: 0;height: auto;line-height: 100%;outline: none;text-decoration: none;"></a> &nbsp;&nbsp;
                                <a href="#" style="-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;text-decoration: none;color: #959ba0;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;"><img src="images/instagram.png" width="24" height="24" alt="Instagram" style="max-width: 24px;-ms-interpolation-mode: bicubic;border: 0;height: auto;line-height: 100%;outline: none;text-decoration: none;"></a> &nbsp;&nbsp;
                                <a href="#" style="-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;text-decoration: none;color: #959ba0;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;"><img src="images/pinterest.png" width="24" height="24" alt="Pinterest" style="max-width: 24px;-ms-interpolation-mode: bicubic;border: 0;height: auto;line-height: 100%;outline: none;text-decoration: none;"></a>
                              </p>
                              <p class="mb" style="color: inherit;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 16px;word-break: break-word;font-size: 15px;line-height: 26px;">
                              	© 2025 STL Party Helpers<br> <a href="https://stlpartyhelpers.com">stlpartyhelpers.com</a><br/>
                              	<span class="text_adr" href="#" style="-webkit-text-size-adjust: 100%;-ms-text-size-adjust: 100%;text-decoration: none;color: #959ba0;font-family: Arial, Helvetica, sans-serif;margin-top: 0px;margin-bottom: 0px;word-break: break-word;">
                                <span style="color: #959ba0;text-decoration: none;">4220 Duncan Ave., Ste. 201</span><br/>
                                <span style="color: #959ba0;text-decoration: none;">St. Louis, MO 63110</span>
                                
                                </span>
                                <br/>
                              <a href="tel:13147145514" style="background-color:#ffffff;color: #262578;padding:8px 8px;text-decoration:none;border-radius:4px;font-family:Arial,sans-serif;display:inline-block;margin-top:10px;border: 3px solid #262578;font-weight: bold;" target="_blank">
 Tap to Call Us: (314) 714-5514
</a>

                              </p>
                               <span id="stlth-version" style="color: #959ba0;text-decoration: none; font-size: 8px;">v 1.00</span>
                            </td>
                          </tr>
                        </tbody>
                      </table>                             
     
                             
                              </td>
                            </tr>
                          </tbody>
                        </table>
                      </div>
                    <!--[if (mso)|(IE)]>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                    <![endif]-->
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          <!--[if (mso)|(IE)]>
                </td>
              </tr>
            </tbody>
          </table>
          <![endif]-->
        </td>
      </tr>
    </tbody>
</table>



   
    </div>
  </div>
</body>

  `;
}

