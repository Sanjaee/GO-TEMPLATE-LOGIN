package service

import (
	"fmt"
	"net/smtp"
	"time"
	"yourapp/internal/config"
)

type EmailService interface {
	SendOTPEmail(to, otpCode string) error
	SendResetPasswordEmail(to, otpCode string) error
	SendVerificationEmail(to, token string) error
	SendWelcomeEmail(to, name string) error
}

type emailService struct {
	config *config.Config
}

func NewEmailService(cfg *config.Config) EmailService {
	return &emailService{
		config: cfg,
	}
}

func (s *emailService) sendEmail(to, subject, body string) error {
	return s.sendEmailHTML(to, subject, body, body)
}

func (s *emailService) sendEmailHTML(to, subject, htmlBody, textBody string) error {
	if s.config.SMTPUsername == "" || s.config.SMTPPassword == "" {
		// In development, just log the email
		fmt.Printf("[EMAIL] To: %s, Subject: %s\nBody: %s\n", to, subject, textBody)
		return nil
	}

	auth := smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPHost)

	from := s.config.EmailFrom
	if from == "" {
		from = s.config.SMTPUsername
	}

	// Create multipart message with HTML and plain text
	boundary := "----=_NextPart_" + fmt.Sprintf("%d", time.Now().UnixNano())

	headers := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: multipart/alternative; boundary=\"%s\"\r\n\r\n",
		from, to, subject, boundary)

	// Plain text part
	textPart := fmt.Sprintf("--%s\r\nContent-Type: text/plain; charset=UTF-8\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\n%s\r\n",
		boundary, textBody)

	// HTML part
	htmlPart := fmt.Sprintf("--%s\r\nContent-Type: text/html; charset=UTF-8\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\n%s\r\n",
		boundary, htmlBody)

	// End boundary
	endBoundary := fmt.Sprintf("--%s--\r\n", boundary)

	msg := []byte(headers + textPart + htmlPart + endBoundary)

	addr := fmt.Sprintf("%s:%s", s.config.SMTPHost, s.config.SMTPPort)
	err := smtp.SendMail(addr, auth, from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (s *emailService) SendOTPEmail(to, otpCode string) error {
	subject := "Kode Verifikasi Email Anda"

	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>%s</title>
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: #f5f5f5;">
	<table role="presentation" style="width: 100%%; border-collapse: collapse; background-color: #f5f5f5;">
		<tr>
			<td style="padding: 40px 20px;">
				<table role="presentation" style="max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
					<tr>
						<td style="padding: 40px 30px; text-align: center; background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); border-radius: 8px 8px 0 0;">
							<h1 style="margin: 0; color: #ffffff; font-size: 28px; font-weight: 600;">Kode Verifikasi</h1>
						</td>
					</tr>
					<tr>
						<td style="padding: 40px 30px;">
							<p style="margin: 0 0 20px 0; color: #333333; font-size: 16px; line-height: 1.6;">Halo,</p>
							<p style="margin: 0 0 30px 0; color: #333333; font-size: 16px; line-height: 1.6;">Terima kasih telah mendaftar. Gunakan kode verifikasi berikut untuk memverifikasi email Anda:</p>
							
							<div style="background-color: #f8f9fa; border: 2px dashed #667eea; border-radius: 8px; padding: 30px; text-align: center; margin: 30px 0;">
								<div style="font-size: 36px; font-weight: 700; color: #667eea; letter-spacing: 8px; font-family: 'Courier New', monospace;">%s</div>
							</div>
							
							<p style="margin: 20px 0; color: #666666; font-size: 14px; line-height: 1.6; text-align: center;">
								<strong style="color: #e74c3c;">⏰ Kode ini akan kedaluwarsa dalam 10 menit.</strong>
							</p>
							
							<div style="margin-top: 40px; padding-top: 30px; border-top: 1px solid #e0e0e0;">
								<p style="margin: 0 0 10px 0; color: #999999; font-size: 12px; line-height: 1.6;">Jika Anda tidak meminta kode ini, abaikan email ini.</p>
								<p style="margin: 0; color: #999999; font-size: 12px; line-height: 1.6;">Terima kasih,<br><strong style="color: #667eea;">Tim YouApp</strong></p>
							</div>
						</td>
					</tr>
				</table>
			</td>
		</tr>
	</table>
</body>
</html>
`, subject, otpCode)

	// Plain text fallback
	textBody := fmt.Sprintf(`
Halo,

Terima kasih telah mendaftar. Gunakan kode verifikasi berikut untuk memverifikasi email Anda:

Kode OTP: %s

Kode ini akan kedaluwarsa dalam 10 menit.

Jika Anda tidak meminta kode ini, abaikan email ini.

Terima kasih,
Tim YouApp
`, otpCode)

	return s.sendEmailHTML(to, subject, htmlBody, textBody)
}

func (s *emailService) SendResetPasswordEmail(to, resetLink string) error {
	subject := "Reset Password"
	
	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>%s</title>
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: #f5f5f5;">
	<table role="presentation" style="width: 100%%; border-collapse: collapse; background-color: #f5f5f5;">
		<tr>
			<td style="padding: 40px 20px;">
				<table role="presentation" style="max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
					<tr>
						<td style="padding: 40px 30px; text-align: center; background: linear-gradient(135deg, #f093fb 0%%, #f5576c 100%%); border-radius: 8px 8px 0 0;">
							<h1 style="margin: 0; color: #ffffff; font-size: 28px; font-weight: 600;">Reset Password</h1>
						</td>
					</tr>
					<tr>
						<td style="padding: 40px 30px;">
							<p style="margin: 0 0 20px 0; color: #333333; font-size: 16px; line-height: 1.6;">Halo,</p>
							<p style="margin: 0 0 30px 0; color: #333333; font-size: 16px; line-height: 1.6;">Anda telah meminta untuk mereset password. Klik tombol di bawah ini untuk melanjutkan:</p>
							
							<div style="text-align: center; margin: 40px 0;">
								<a href="%s" style="display: inline-block; padding: 14px 32px; background: linear-gradient(135deg, #f093fb 0%%, #f5576c 100%%); color: #ffffff; text-decoration: none; border-radius: 6px; font-weight: 600; font-size: 16px; box-shadow: 0 4px 6px rgba(245, 87, 108, 0.3);">Reset Password</a>
							</div>
							
							<p style="margin: 20px 0; color: #666666; font-size: 14px; line-height: 1.6; text-align: center;">
								Atau salin dan buka link berikut di browser Anda:
							</p>
							<div style="background-color: #f8f9fa; border: 1px solid #e0e0e0; border-radius: 6px; padding: 12px; margin: 20px 0; word-break: break-all;">
								<p style="margin: 0; color: #667eea; font-size: 12px; font-family: 'Courier New', monospace;">%s</p>
							</div>
							
							<p style="margin: 20px 0; color: #666666; font-size: 14px; line-height: 1.6; text-align: center;">
								<strong style="color: #e74c3c;">⏰ Link ini akan kedaluwarsa dalam 1 jam.</strong>
							</p>
							
							<div style="margin-top: 40px; padding-top: 30px; border-top: 1px solid #e0e0e0;">
								<p style="margin: 0 0 10px 0; color: #999999; font-size: 12px; line-height: 1.6;">Jika Anda tidak meminta reset password, abaikan email ini.</p>
								<p style="margin: 0; color: #999999; font-size: 12px; line-height: 1.6;">Terima kasih,<br><strong style="color: #f5576c;">Tim YouApp</strong></p>
							</div>
						</td>
					</tr>
				</table>
			</td>
		</tr>
	</table>
</body>
</html>
`, subject, resetLink, resetLink)

	textBody := fmt.Sprintf(`
Halo,

Anda telah meminta untuk mereset password. Klik link berikut untuk melanjutkan:

%s

Link ini akan kedaluwarsa dalam 1 jam.

Jika Anda tidak meminta reset password, abaikan email ini.

Terima kasih,
Tim YouApp
`, resetLink)

	return s.sendEmailHTML(to, subject, htmlBody, textBody)
}

func (s *emailService) SendVerificationEmail(to, token string) error {
	subject := "Verifikasi Email Anda"
	verificationURL := fmt.Sprintf("%s/auth/verify-email?token=%s", s.config.ClientURL, token)
	body := fmt.Sprintf(`
Halo,

Terima kasih telah mendaftar. Klik link berikut untuk memverifikasi email Anda:

%s

Link ini akan kedaluwarsa dalam 24 jam.

Jika Anda tidak meminta verifikasi ini, abaikan email ini.

Terima kasih,
Tim YouApp
`, verificationURL)

	return s.sendEmail(to, subject, body)
}

func (s *emailService) SendWelcomeEmail(to, name string) error {
	subject := "Selamat Datang di YouApp"
	body := fmt.Sprintf(`
Halo %s,

Selamat datang di YouApp! Kami senang Anda bergabung dengan kami.

Jika Anda memiliki pertanyaan, jangan ragu untuk menghubungi kami.

Terima kasih,
Tim YouApp
`, name)

	return s.sendEmail(to, subject, body)
}
