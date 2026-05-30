package emails

// All templates use Go's html/template syntax.
// Template data maps are defined in subscriptions.go per handler.

// baseOpen is the shared HTML header / wrapper (open tag).
const baseOpen = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{.Subject}}</title>
</head>
<body style="margin:0;padding:0;background-color:#f4f4f4;font-family:Arial,Helvetica,sans-serif;">
  <table width="100%" cellpadding="0" cellspacing="0" style="background-color:#f4f4f4;padding:30px 0;">
    <tr>
      <td align="center">
        <table width="600" cellpadding="0" cellspacing="0" style="max-width:600px;width:100%;background-color:#ffffff;border-radius:8px;overflow:hidden;box-shadow:0 2px 8px rgba(0,0,0,0.08);">
          <!-- Header -->
          <tr>
            <td style="background-color:#1a1a2e;padding:24px 32px;">
              <p style="margin:0;color:#ffffff;font-size:22px;font-weight:bold;letter-spacing:1px;">{{.StoreName}}</p>
            </td>
          </tr>
          <!-- Body -->
          <tr>
            <td style="padding:32px;">`

// baseClose is the shared HTML footer / wrapper (close tag).
const baseClose = `
            </td>
          </tr>
          <!-- Footer -->
          <tr>
            <td style="background-color:#f9f9f9;padding:20px 32px;border-top:1px solid #eeeeee;">
              <p style="margin:0;color:#999999;font-size:12px;text-align:center;">
                This is an automated message from {{.StoreName}}. Please do not reply to this email.
              </p>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>`

// ─── Order Confirmation ───────────────────────────────────────────────────────

const orderConfirmationTmpl = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Thank you for your order</title>
</head>
<body style="margin:0;padding:0;background-color:#f4f4f4;font-family:Arial,Helvetica,sans-serif;">
  <table width="100%" cellpadding="0" cellspacing="0" style="background-color:#f4f4f4;padding:30px 0;">
    <tr><td align="center">
      <table width="600" cellpadding="0" cellspacing="0" style="max-width:600px;width:100%;background-color:#ffffff;border-radius:8px;overflow:hidden;box-shadow:0 2px 8px rgba(0,0,0,0.08);">
        <tr>
          <td style="background-color:#1a1a2e;padding:24px 32px;">
            <p style="margin:0;color:#ffffff;font-size:22px;font-weight:bold;">{{.StoreName}}</p>
          </td>
        </tr>
        <tr>
          <td style="padding:32px;">
            <h1 style="margin:0 0 8px;color:#1a1a2e;font-size:24px;">Thank you for your order!</h1>
            <p style="margin:0 0 24px;color:#555555;font-size:16px;">We&#39;ve received your order and are getting it ready.</p>
            <table width="100%" cellpadding="0" cellspacing="0" style="background-color:#f9f9f9;border-radius:6px;padding:20px;margin-bottom:24px;">
              <tr>
                <td style="padding:8px 0;">
                  <p style="margin:0;color:#777777;font-size:13px;text-transform:uppercase;letter-spacing:0.5px;">Order ID</p>
                  <p style="margin:4px 0 0;color:#1a1a2e;font-size:16px;font-weight:bold;">#{{.OrderID}}</p>
                </td>
              </tr>
              <tr>
                <td style="padding:8px 0;border-top:1px solid #eeeeee;">
                  <p style="margin:0;color:#777777;font-size:13px;text-transform:uppercase;letter-spacing:0.5px;">Items</p>
                  <p style="margin:4px 0 0;color:#1a1a2e;font-size:16px;">{{.ItemCount}}</p>
                </td>
              </tr>
              <tr>
                <td style="padding:8px 0;border-top:1px solid #eeeeee;">
                  <p style="margin:0;color:#777777;font-size:13px;text-transform:uppercase;letter-spacing:0.5px;">Total</p>
                  <p style="margin:4px 0 0;color:#1a1a2e;font-size:16px;font-weight:bold;">{{.Total}}</p>
                </td>
              </tr>
              <tr>
                <td style="padding:8px 0;border-top:1px solid #eeeeee;">
                  <p style="margin:0;color:#777777;font-size:13px;text-transform:uppercase;letter-spacing:0.5px;">Status</p>
                  <p style="margin:4px 0 0;color:#1a1a2e;font-size:16px;">{{.Status}}</p>
                </td>
              </tr>
            </table>
            <p style="margin:0;color:#555555;font-size:14px;">We&#39;ll send you another email when your order ships.</p>
          </td>
        </tr>
        <tr>
          <td style="background-color:#f9f9f9;padding:20px 32px;border-top:1px solid #eeeeee;">
            <p style="margin:0;color:#999999;font-size:12px;text-align:center;">This is an automated message from {{.StoreName}}. Please do not reply to this email.</p>
          </td>
        </tr>
      </table>
    </td></tr>
  </table>
</body>
</html>`

// ─── Order Shipped ────────────────────────────────────────────────────────────

const orderShippedTmpl = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Your order has been shipped</title>
</head>
<body style="margin:0;padding:0;background-color:#f4f4f4;font-family:Arial,Helvetica,sans-serif;">
  <table width="100%" cellpadding="0" cellspacing="0" style="background-color:#f4f4f4;padding:30px 0;">
    <tr><td align="center">
      <table width="600" cellpadding="0" cellspacing="0" style="max-width:600px;width:100%;background-color:#ffffff;border-radius:8px;overflow:hidden;box-shadow:0 2px 8px rgba(0,0,0,0.08);">
        <tr>
          <td style="background-color:#1a1a2e;padding:24px 32px;">
            <p style="margin:0;color:#ffffff;font-size:22px;font-weight:bold;">{{.StoreName}}</p>
          </td>
        </tr>
        <tr>
          <td style="padding:32px;">
            <h1 style="margin:0 0 8px;color:#1a1a2e;font-size:24px;">Your order is on its way!</h1>
            <p style="margin:0 0 24px;color:#555555;font-size:16px;">Great news — your order has been shipped and is heading your way.</p>
            <table width="100%" cellpadding="0" cellspacing="0" style="background-color:#f9f9f9;border-radius:6px;padding:20px;margin-bottom:24px;">
              <tr>
                <td style="padding:8px 0;">
                  <p style="margin:0;color:#777777;font-size:13px;text-transform:uppercase;letter-spacing:0.5px;">Order ID</p>
                  <p style="margin:4px 0 0;color:#1a1a2e;font-size:16px;font-weight:bold;">#{{.OrderID}}</p>
                </td>
              </tr>
              <tr>
                <td style="padding:8px 0;border-top:1px solid #eeeeee;">
                  <p style="margin:0;color:#777777;font-size:13px;text-transform:uppercase;letter-spacing:0.5px;">Status</p>
                  <p style="margin:4px 0 0;color:#1a1a2e;font-size:16px;">{{.Status}}</p>
                </td>
              </tr>
            </table>
          </td>
        </tr>
        <tr>
          <td style="background-color:#f9f9f9;padding:20px 32px;border-top:1px solid #eeeeee;">
            <p style="margin:0;color:#999999;font-size:12px;text-align:center;">This is an automated message from {{.StoreName}}. Please do not reply to this email.</p>
          </td>
        </tr>
      </table>
    </td></tr>
  </table>
</body>
</html>`

// ─── Order Delivered ──────────────────────────────────────────────────────────

const orderDeliveredTmpl = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Your order has been delivered</title>
</head>
<body style="margin:0;padding:0;background-color:#f4f4f4;font-family:Arial,Helvetica,sans-serif;">
  <table width="100%" cellpadding="0" cellspacing="0" style="background-color:#f4f4f4;padding:30px 0;">
    <tr><td align="center">
      <table width="600" cellpadding="0" cellspacing="0" style="max-width:600px;width:100%;background-color:#ffffff;border-radius:8px;overflow:hidden;box-shadow:0 2px 8px rgba(0,0,0,0.08);">
        <tr>
          <td style="background-color:#1a1a2e;padding:24px 32px;">
            <p style="margin:0;color:#ffffff;font-size:22px;font-weight:bold;">{{.StoreName}}</p>
          </td>
        </tr>
        <tr>
          <td style="padding:32px;">
            <h1 style="margin:0 0 8px;color:#1a1a2e;font-size:24px;">Your order has been delivered!</h1>
            <p style="margin:0 0 24px;color:#555555;font-size:16px;">We hope you enjoy your purchase. Thank you for shopping with us!</p>
            <table width="100%" cellpadding="0" cellspacing="0" style="background-color:#f9f9f9;border-radius:6px;padding:20px;margin-bottom:24px;">
              <tr>
                <td style="padding:8px 0;">
                  <p style="margin:0;color:#777777;font-size:13px;text-transform:uppercase;letter-spacing:0.5px;">Order ID</p>
                  <p style="margin:4px 0 0;color:#1a1a2e;font-size:16px;font-weight:bold;">#{{.OrderID}}</p>
                </td>
              </tr>
            </table>
          </td>
        </tr>
        <tr>
          <td style="background-color:#f9f9f9;padding:20px 32px;border-top:1px solid #eeeeee;">
            <p style="margin:0;color:#999999;font-size:12px;text-align:center;">This is an automated message from {{.StoreName}}. Please do not reply to this email.</p>
          </td>
        </tr>
      </table>
    </td></tr>
  </table>
</body>
</html>`

// ─── Customer Welcome ─────────────────────────────────────────────────────────

const customerWelcomeTmpl = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Welcome!</title>
</head>
<body style="margin:0;padding:0;background-color:#f4f4f4;font-family:Arial,Helvetica,sans-serif;">
  <table width="100%" cellpadding="0" cellspacing="0" style="background-color:#f4f4f4;padding:30px 0;">
    <tr><td align="center">
      <table width="600" cellpadding="0" cellspacing="0" style="max-width:600px;width:100%;background-color:#ffffff;border-radius:8px;overflow:hidden;box-shadow:0 2px 8px rgba(0,0,0,0.08);">
        <tr>
          <td style="background-color:#1a1a2e;padding:24px 32px;">
            <p style="margin:0;color:#ffffff;font-size:22px;font-weight:bold;">{{.StoreName}}</p>
          </td>
        </tr>
        <tr>
          <td style="padding:32px;">
            <h1 style="margin:0 0 8px;color:#1a1a2e;font-size:24px;">Welcome to {{.StoreName}}!</h1>
            <p style="margin:0 0 16px;color:#555555;font-size:16px;">Hi {{.Name}}, we&#39;re excited to have you on board.</p>
            <p style="margin:0 0 24px;color:#555555;font-size:16px;">Your account has been created with the email address: <strong>{{.Email}}</strong></p>
            <p style="margin:0;color:#555555;font-size:14px;">Start exploring our store and enjoy your shopping experience!</p>
          </td>
        </tr>
        <tr>
          <td style="background-color:#f9f9f9;padding:20px 32px;border-top:1px solid #eeeeee;">
            <p style="margin:0;color:#999999;font-size:12px;text-align:center;">This is an automated message from {{.StoreName}}. Please do not reply to this email.</p>
          </td>
        </tr>
      </table>
    </td></tr>
  </table>
</body>
</html>`

// ─── Payment Completed ────────────────────────────────────────────────────────

const paymentCompletedTmpl = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Payment confirmed</title>
</head>
<body style="margin:0;padding:0;background-color:#f4f4f4;font-family:Arial,Helvetica,sans-serif;">
  <table width="100%" cellpadding="0" cellspacing="0" style="background-color:#f4f4f4;padding:30px 0;">
    <tr><td align="center">
      <table width="600" cellpadding="0" cellspacing="0" style="max-width:600px;width:100%;background-color:#ffffff;border-radius:8px;overflow:hidden;box-shadow:0 2px 8px rgba(0,0,0,0.08);">
        <tr>
          <td style="background-color:#1a1a2e;padding:24px 32px;">
            <p style="margin:0;color:#ffffff;font-size:22px;font-weight:bold;">{{.StoreName}}</p>
          </td>
        </tr>
        <tr>
          <td style="padding:32px;">
            <h1 style="margin:0 0 8px;color:#1a1a2e;font-size:24px;">Payment confirmed</h1>
            <p style="margin:0 0 24px;color:#555555;font-size:16px;">Your payment has been successfully processed.</p>
            <table width="100%" cellpadding="0" cellspacing="0" style="background-color:#f9f9f9;border-radius:6px;padding:20px;margin-bottom:24px;">
              <tr>
                <td style="padding:8px 0;">
                  <p style="margin:0;color:#777777;font-size:13px;text-transform:uppercase;letter-spacing:0.5px;">Order ID</p>
                  <p style="margin:4px 0 0;color:#1a1a2e;font-size:16px;font-weight:bold;">#{{.OrderID}}</p>
                </td>
              </tr>
              <tr>
                <td style="padding:8px 0;border-top:1px solid #eeeeee;">
                  <p style="margin:0;color:#777777;font-size:13px;text-transform:uppercase;letter-spacing:0.5px;">Amount Charged</p>
                  <p style="margin:4px 0 0;color:#1a1a2e;font-size:16px;font-weight:bold;">{{.Amount}}</p>
                </td>
              </tr>
            </table>
          </td>
        </tr>
        <tr>
          <td style="background-color:#f9f9f9;padding:20px 32px;border-top:1px solid #eeeeee;">
            <p style="margin:0;color:#999999;font-size:12px;text-align:center;">This is an automated message from {{.StoreName}}. Please do not reply to this email.</p>
          </td>
        </tr>
      </table>
    </td></tr>
  </table>
</body>
</html>`

// ─── Refund Completed ─────────────────────────────────────────────────────────

const refundCompletedTmpl = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Refund processed</title>
</head>
<body style="margin:0;padding:0;background-color:#f4f4f4;font-family:Arial,Helvetica,sans-serif;">
  <table width="100%" cellpadding="0" cellspacing="0" style="background-color:#f4f4f4;padding:30px 0;">
    <tr><td align="center">
      <table width="600" cellpadding="0" cellspacing="0" style="max-width:600px;width:100%;background-color:#ffffff;border-radius:8px;overflow:hidden;box-shadow:0 2px 8px rgba(0,0,0,0.08);">
        <tr>
          <td style="background-color:#1a1a2e;padding:24px 32px;">
            <p style="margin:0;color:#ffffff;font-size:22px;font-weight:bold;">{{.StoreName}}</p>
          </td>
        </tr>
        <tr>
          <td style="padding:32px;">
            <h1 style="margin:0 0 8px;color:#1a1a2e;font-size:24px;">Refund processed</h1>
            <p style="margin:0 0 24px;color:#555555;font-size:16px;">Your refund has been successfully processed. Please allow a few business days for it to appear in your account.</p>
            <table width="100%" cellpadding="0" cellspacing="0" style="background-color:#f9f9f9;border-radius:6px;padding:20px;margin-bottom:24px;">
              <tr>
                <td style="padding:8px 0;">
                  <p style="margin:0;color:#777777;font-size:13px;text-transform:uppercase;letter-spacing:0.5px;">Order ID</p>
                  <p style="margin:4px 0 0;color:#1a1a2e;font-size:16px;font-weight:bold;">#{{.OrderID}}</p>
                </td>
              </tr>
              <tr>
                <td style="padding:8px 0;border-top:1px solid #eeeeee;">
                  <p style="margin:0;color:#777777;font-size:13px;text-transform:uppercase;letter-spacing:0.5px;">Refund Amount</p>
                  <p style="margin:4px 0 0;color:#1a1a2e;font-size:16px;font-weight:bold;">{{.Amount}}</p>
                </td>
              </tr>
            </table>
          </td>
        </tr>
        <tr>
          <td style="background-color:#f9f9f9;padding:20px 32px;border-top:1px solid #eeeeee;">
            <p style="margin:0;color:#999999;font-size:12px;text-align:center;">This is an automated message from {{.StoreName}}. Please do not reply to this email.</p>
          </td>
        </tr>
      </table>
    </td></tr>
  </table>
</body>
</html>`

// baseOpen and baseClose are defined but not used directly — they serve as
// documentation of the shared layout structure used across the templates above.
var _ = baseOpen
var _ = baseClose
