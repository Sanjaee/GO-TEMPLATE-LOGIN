package service

import (
	"fmt"
	"net/smtp"
	"time"
	"yourapp/internal/config"
)

// EmailService mendefinisikan antarmuka untuk layanan pengiriman email.
type EmailService interface {
	SendOTPEmail(to, otpCode string) error
	SendResetPasswordEmail(to, resetLink string) error
	SendVerificationEmail(to, token string) error
	SendWelcomeEmail(to, name string) error
}

type emailService struct {
	config *config.Config
}

// NewEmailService membuat instance baru dari EmailService.
func NewEmailService(cfg *config.Config) EmailService {
	return &emailService{
		config: cfg,
	}
}

// sendEmail adalah helper untuk mengirim email tanpa HTML (text-only fallback).
func (s *emailService) sendEmail(to, subject, body string) error {
	return s.sendEmailHTML(to, subject, body, body)
}

// sendEmailHTML mengirim email multipart dengan versi HTML dan plain text.
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

	// Use custom email name if available
	fromHeader := from
	if s.config.EmailName != "" {
		fromHeader = fmt.Sprintf("%s <%s>", s.config.EmailName, from)
	}

	// Create multipart message with HTML and plain text
	boundary := "----=_NextPart_" + fmt.Sprintf("%d", time.Now().UnixNano())

	headers := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: multipart/alternative; boundary=\"%s\"\r\n\r\n",
		fromHeader, to, subject, boundary)

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
	subject := "Verifikasi Email - Kode OTP Anda"

	// Template HTML Modern dengan Warna dan OTP Besar - Responsive
	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>%s</title>
	<style>
		* { margin: 0; padding: 0; box-sizing: border-box; }
		body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: #333; line-height: 1.6; }
		.email-wrapper { width: 100%%; background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); padding: 40px 20px; min-height: 100vh; }
		.email-container { max-width: 600px; margin: 0 auto; background-color: white; border-radius: 16px; box-shadow: 0 10px 40px rgba(0,0,0,0.2); overflow: hidden; }
		.email-header { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); padding: 50px 30px; text-align: center; }
		.email-header h1 { color: white; font-size: 32px; font-weight: 700; margin: 0; letter-spacing: -0.5px; }
		.email-header .subtitle { color: rgba(255,255,255,0.9); font-size: 16px; margin-top: 10px; }
		.email-body { padding: 50px 30px; background-color: white; }
		.email-body p { margin: 0 0 20px 0; font-size: 16px; color: #555; }
		.otp-container { background: linear-gradient(135deg, #f5f7fa 0%%, #c3cfe2 100%%); border: 3px solid #667eea; border-radius: 16px; padding: 40px 30px; margin: 40px 0; text-align: center; box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3); }
		.otp-code { font-size: 48px; font-weight: 800; color: #667eea; letter-spacing: 12px; font-family: 'Courier New', monospace; margin: 20px 0; text-shadow: 2px 2px 4px rgba(0,0,0,0.1); }
		.otp-label { font-size: 14px; color: #667eea; text-transform: uppercase; letter-spacing: 2px; margin-bottom: 15px; font-weight: 600; }
		.warning-box { background: linear-gradient(135deg, #fff5f5 0%%, #ffe5e5 100%%); border-left: 5px solid #ef4444; padding: 20px; margin: 30px 0; border-radius: 8px; }
		.warning-box p { margin: 0; font-size: 14px; color: #dc2626; font-weight: 600; }
		.email-footer { background-color: #f8f9fa; padding: 30px; text-align: center; border-top: 1px solid #e5e7eb; }
		.email-footer p { margin: 5px 0; font-size: 12px; color: #6b7280; }
		.brand-name { color: #667eea; font-weight: 700; }
		@media only screen and (max-width: 600px) {
			.email-wrapper { padding: 20px 10px; }
			.email-header { padding: 40px 20px; }
			.email-header h1 { font-size: 26px; }
			.email-body { padding: 40px 20px; }
			.otp-code { font-size: 36px; letter-spacing: 8px; }
			.otp-container { padding: 30px 20px; }
		}
	</style>
</head>
<body>
	<div class="email-wrapper">
		<div class="email-container">
			<div class="email-header">
				<h1>🔐 Verifikasi Email</h1>
				<p class="subtitle">Kode OTP Anda</p>
			</div>
			<div class="email-body">
				<p>Halo,</p>
				<p>Terima kasih telah mendaftar di <span class="brand-name">%s</span>! Gunakan kode verifikasi di bawah ini untuk memverifikasi email Anda:</p>
				
				<div class="otp-container">
					<div class="otp-label">Kode Verifikasi</div>
					<div class="otp-code">%s</div>
				</div>
				
				<div class="warning-box">
					<p>⚠️ Kode ini berlaku selama 10 menit. Jangan bagikan kode ini kepada siapapun.</p>
				</div>
				
				<p style="font-size: 14px; color: #6b7280;">Jika Anda tidak meminta kode verifikasi ini, silakan abaikan email ini atau hubungi tim support kami.</p>
			</div>
			<div class="email-footer">
				<p>&copy; %d <span class="brand-name">%s</span>. All rights reserved.</p>
				<p>Email ini dikirim secara otomatis, mohon jangan membalas email ini.</p>
			</div>
		</div>
	</div>
</body>
</html>
`, subject, s.config.EmailName, otpCode, time.Now().Year(), s.config.EmailName)

	// Plain text fallback
	textBody := fmt.Sprintf(`
Halo,

Terima kasih telah mendaftar di %s! Gunakan kode verifikasi berikut untuk memverifikasi email Anda:

Kode OTP: %s

Kode ini akan kedaluwarsa dalam 10 menit.

Demi keamanan, jangan pernah membagikan kode ini kepada siapapun.

Jika Anda tidak meminta kode ini, abaikan email ini.

Terima kasih,
Tim %s
`, s.config.EmailName, otpCode, s.config.EmailName)

	return s.sendEmailHTML(to, subject, htmlBody, textBody)
}

func (s *emailService) SendResetPasswordEmail(to, otpCode string) error {
	subject := "Reset Password - Kode OTP Anda"

	// Template HTML Modern Reset Password dengan Warna - Responsive
	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>%s</title>
	<style>
		* { margin: 0; padding: 0; box-sizing: border-box; }
		body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background: linear-gradient(135deg, #f093fb 0%%, #f5576c 100%%); color: #333; line-height: 1.6; }
		.email-wrapper { width: 100%%; background: linear-gradient(135deg, #f093fb 0%%, #f5576c 100%%); padding: 40px 20px; min-height: 100vh; }
		.email-container { max-width: 600px; margin: 0 auto; background-color: white; border-radius: 16px; box-shadow: 0 10px 40px rgba(0,0,0,0.2); overflow: hidden; }
		.email-header { background: linear-gradient(135deg, #f093fb 0%%, #f5576c 100%%); padding: 50px 30px; text-align: center; }
		.email-header h1 { color: white; font-size: 32px; font-weight: 700; margin: 0; letter-spacing: -0.5px; }
		.email-header .subtitle { color: rgba(255,255,255,0.9); font-size: 16px; margin-top: 10px; }
		.email-body { padding: 50px 30px; background-color: white; }
		.email-body p { margin: 0 0 20px 0; font-size: 16px; color: #555; }
		.otp-container { background: linear-gradient(135deg, #fff5f7 0%%, #ffeef2 100%%); border: 3px solid #f5576c; border-radius: 16px; padding: 40px 30px; margin: 40px 0; text-align: center; box-shadow: 0 4px 15px rgba(245, 87, 108, 0.3); }
		.otp-code { font-size: 48px; font-weight: 800; color: #f5576c; letter-spacing: 12px; font-family: 'Courier New', monospace; margin: 20px 0; text-shadow: 2px 2px 4px rgba(0,0,0,0.1); }
		.otp-label { font-size: 14px; color: #f5576c; text-transform: uppercase; letter-spacing: 2px; margin-bottom: 15px; font-weight: 600; }
		.info-box { background: linear-gradient(135deg, #fff9e6 0%%, #fff5d6 100%%); border-left: 5px solid #f59e0b; padding: 20px; margin: 30px 0; border-radius: 8px; }
		.info-box p { margin: 0 0 8px 0; font-size: 14px; color: #92400e; font-weight: 500; }
		.info-box strong { color: #78350f; }
		.warning-box { background: linear-gradient(135deg, #fff5f5 0%%, #ffe5e5 100%%); border-left: 5px solid #ef4444; padding: 20px; margin: 30px 0; border-radius: 8px; }
		.warning-box p { margin: 0; font-size: 14px; color: #dc2626; font-weight: 600; }
		.email-footer { background-color: #f8f9fa; padding: 30px; text-align: center; border-top: 1px solid #e5e7eb; }
		.email-footer p { margin: 5px 0; font-size: 12px; color: #6b7280; }
		.brand-name { color: #f5576c; font-weight: 700; }
		@media only screen and (max-width: 600px) {
			.email-wrapper { padding: 20px 10px; }
			.email-header { padding: 40px 20px; }
			.email-header h1 { font-size: 26px; }
			.email-body { padding: 40px 20px; }
			.otp-code { font-size: 36px; letter-spacing: 8px; }
			.otp-container { padding: 30px 20px; }
		}
	</style>
</head>
<body>
	<div class="email-wrapper">
		<div class="email-container">
			<div class="email-header">
				<h1>🔑 Reset Password</h1>
				<p class="subtitle">Kode OTP Anda</p>
			</div>
			<div class="email-body">
				<p>Halo,</p>
				<p>Kami menerima permintaan untuk mereset password akun <span class="brand-name">%s</span> Anda. Gunakan kode OTP di bawah ini:</p>
				
				<div class="otp-container">
					<div class="otp-label">Kode Reset Password</div>
					<div class="otp-code">%s</div>
				</div>
				
				<div class="info-box">
					<p><strong>Langkah selanjutnya:</strong></p>
					<p>1. Masukkan kode OTP di atas</p>
					<p>2. Buat password baru yang kuat</p>
					<p>3. Login dengan password baru Anda</p>
				</div>
				
				<div class="warning-box">
					<p>⚠️ Kode ini berlaku selama 10 menit. Jangan bagikan kode ini kepada siapapun.</p>
				</div>
				
				<p style="font-size: 14px; color: #6b7280;">Jika Anda tidak meminta reset password ini, silakan abaikan email ini. Akun Anda tetap aman.</p>
			</div>
			<div class="email-footer">
				<p>&copy; %d <span class="brand-name">%s</span>. All rights reserved.</p>
				<p>Email ini dikirim secara otomatis, mohon jangan membalas email ini.</p>
			</div>
		</div>
	</div>
</body>
</html>
`, subject, s.config.EmailName, otpCode, time.Now().Year(), s.config.EmailName)

	textBody := fmt.Sprintf(`
Halo,

Kami menerima permintaan untuk mereset password akun %s Anda.

Kode OTP Reset Password: %s

Langkah selanjutnya:
1. Masukkan kode OTP di atas
2. Buat password baru yang kuat
3. Login dengan password baru Anda

Kode ini berlaku selama 10 menit. Jangan bagikan kode ini kepada siapapun.

Jika Anda tidak meminta reset password ini, silakan abaikan email ini.

Terima kasih,
Tim %s
`, s.config.EmailName, otpCode, s.config.EmailName)

	return s.sendEmailHTML(to, subject, htmlBody, textBody)
}

func (s *emailService) SendVerificationEmail(to, token string) error {
	subject := "Selamat Datang - Verifikasi Email Anda"
	verificationURL := fmt.Sprintf("%s/auth/verify-email?token=%s", s.config.ClientURL, token)

	// Template HTML Modern Verifikasi Email dengan Warna - Responsive
	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>%s</title>
	<style>
		* { margin: 0; padding: 0; box-sizing: border-box; }
		body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background: linear-gradient(135deg, #4facfe 0%%, #00f2fe 100%%); color: #333; line-height: 1.6; }
		.email-wrapper { width: 100%%; background: linear-gradient(135deg, #4facfe 0%%, #00f2fe 100%%); padding: 40px 20px; min-height: 100vh; }
		.email-container { max-width: 600px; margin: 0 auto; background-color: white; border-radius: 16px; box-shadow: 0 10px 40px rgba(0,0,0,0.2); overflow: hidden; }
		.email-header { background: linear-gradient(135deg, #4facfe 0%%, #00f2fe 100%%); padding: 50px 30px; text-align: center; }
		.email-header h1 { color: white; font-size: 32px; font-weight: 700; margin: 0; letter-spacing: -0.5px; }
		.email-header .subtitle { color: rgba(255,255,255,0.9); font-size: 16px; margin-top: 10px; }
		.email-body { padding: 50px 30px; background-color: white; }
		.email-body p { margin: 0 0 20px 0; font-size: 16px; color: #555; }
		.button-container { text-align: center; margin: 40px 0; }
		.button { display: inline-block; padding: 18px 45px; background: linear-gradient(135deg, #4facfe 0%%, #00f2fe 100%%); color: white !important; text-decoration: none; border-radius: 12px; font-weight: 700; font-size: 18px; letter-spacing: 0.5px; box-shadow: 0 4px 15px rgba(79, 172, 254, 0.4); transition: all 0.3s ease; }
		.button:hover { transform: translateY(-2px); box-shadow: 0 6px 20px rgba(79, 172, 254, 0.5); }
		.link-box { background: linear-gradient(135deg, #f0f9ff 0%%, #e0f2fe 100%%); border: 2px solid #4facfe; border-radius: 12px; padding: 20px; margin: 30px 0; word-break: break-all; font-size: 13px; color: #0369a1; }
		.info-box { background: linear-gradient(135deg, #fff9e6 0%%, #fff5d6 100%%); border-left: 5px solid #f59e0b; padding: 20px; margin: 30px 0; border-radius: 8px; }
		.info-box p { margin: 0; font-size: 14px; color: #92400e; font-weight: 600; }
		.email-footer { background-color: #f8f9fa; padding: 30px; text-align: center; border-top: 1px solid #e5e7eb; }
		.email-footer p { margin: 5px 0; font-size: 12px; color: #6b7280; }
		.brand-name { color: #4facfe; font-weight: 700; }
		@media only screen and (max-width: 600px) {
			.email-wrapper { padding: 20px 10px; }
			.email-header { padding: 40px 20px; }
			.email-header h1 { font-size: 26px; }
			.email-body { padding: 40px 20px; }
			.button { display: block; width: 90%%; margin: 10px auto; padding: 16px 30px; }
		}
	</style>
</head>
<body>
	<div class="email-wrapper">
		<div class="email-container">
			<div class="email-header">
				<h1>🎉 Selamat Datang!</h1>
				<p class="subtitle">Verifikasi Email Anda</p>
			</div>
			<div class="email-body">
				<p>Halo,</p>
				<p>Terima kasih telah mendaftar di <span class="brand-name">%s</span>! Untuk mengaktifkan akun Anda, silakan verifikasi alamat email dengan mengklik tombol di bawah ini:</p>
				
				<div class="button-container">
					<a href="%s" class="button">✨ Verifikasi Email Saya</a>
				</div>
				
				<p style="font-size: 14px; color: #6b7280; text-align: center;">Jika tombol di atas tidak berfungsi, salin dan tempel link berikut ke browser Anda:</p>
				
				<div class="link-box">
					%s
				</div>
				
				<div class="info-box">
					<p>⚠️ <strong>Penting:</strong> Link ini berlaku selama 24 jam. Setelah itu, Anda perlu meminta link verifikasi baru.</p>
				</div>
				
				<p style="font-size: 14px; color: #6b7280;">Jika Anda tidak mendaftar untuk akun ini, silakan abaikan email ini.</p>
			</div>
			<div class="email-footer">
				<p>&copy; %d <span class="brand-name">%s</span>. All rights reserved.</p>
				<p>Email ini dikirim secara otomatis, mohon jangan membalas email ini.</p>
			</div>
		</div>
	</div>
</body>
</html>
	`, subject, s.config.EmailName, verificationURL, verificationURL, time.Now().Year(), s.config.EmailName)

	textBody := fmt.Sprintf(`
Halo,

Terima kasih telah mendaftar di %s! Klik link berikut untuk memverifikasi email Anda:

%s

Link ini akan kedaluwarsa dalam 24 jam.

Jika Anda tidak meminta verifikasi ini, abaikan email ini.

Terima kasih,
Tim %s
`, s.config.EmailName, verificationURL, s.config.EmailName)

	return s.sendEmailHTML(to, subject, htmlBody, textBody)
}

func (s *emailService) SendWelcomeEmail(to, name string) error {
	subject := "Selamat Datang di BeRealTime"

	// Template HTML Profesional Modern - Tanpa RGB
	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>%s</title>
	<style>
		* { margin: 0; padding: 0; box-sizing: border-box; }
		body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: whitesmoke; color: darkslategray; line-height: 1.6; }
		.email-wrapper { width: 100%%; background-color: whitesmoke; padding: 40px 20px; }
		.email-container { max-width: 600px; margin: 0 auto; background-color: white; border-radius: 12px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); overflow: hidden; }
		.email-header { background: linear-gradient(135deg, black 0%%, darkslategray 100%%); padding: 50px 30px; text-align: center; }
		.email-header h1 { color: white; font-size: 32px; font-weight: 700; margin: 0; letter-spacing: -0.5px; }
		.email-header p { color: white; font-size: 16px; margin: 10px 0 0 0; opacity: 0.9; }
		.email-body { padding: 40px 30px; }
		.email-body p { margin: 0 0 20px 0; font-size: 16px; color: darkslategray; }
		.highlight-box { background: linear-gradient(135deg, whitesmoke 0%%, white 100%%); border: 2px solid gainsboro; border-radius: 12px; padding: 25px; margin: 30px 0; }
		.highlight-box h2 { color: black; font-size: 20px; margin: 0 0 15px 0; }
		.highlight-box ul { margin: 10px 0 0 20px; padding: 0; }
		.highlight-box li { margin: 8px 0; color: darkslategray; font-size: 15px; }
		.email-footer { background-color: whitesmoke; padding: 25px 30px; text-align: center; border-top: 1px solid gainsboro; }
		.email-footer p { margin: 5px 0; font-size: 12px; color: gray; }
		@media only screen and (max-width: 600px) {
			.email-wrapper { padding: 20px 10px; }
			.email-header { padding: 40px 20px; }
			.email-header h1 { font-size: 26px; }
			.email-body { padding: 30px 20px; }
		}
	</style>
</head>
<body>
	<div class="email-wrapper">
		<div class="email-container">
			<div class="email-header">
				<h1>🎉 Selamat Datang!</h1>
				<p>Di BeRealTime</p>
			</div>
			<div class="email-body">
				<p>Halo <strong>%s</strong>,</p>
				<p>Terima kasih telah bergabung dengan BeRealTime! Kami sangat senang menyambut Anda sebagai bagian dari komunitas kami.</p>
				
				<div class="highlight-box">
					<h2>Apa yang bisa Anda lakukan:</h2>
					<ul>
						<li>Nikmati semua fitur yang tersedia</li>
						<li>Jelajahi pengalaman yang menyenangkan</li>
						<li>Hubungi tim support jika ada pertanyaan</li>
					</ul>
				</div>
				
				<p>Jika Anda memiliki pertanyaan atau memerlukan bantuan, jangan ragu untuk menghubungi tim dukungan kami. Kami selalu siap membantu!</p>
				
				<p style="margin-top: 30px;">Hormat kami,<br><strong style="color: black; font-size: 16px;">Tim BeRealTime</strong></p>
			</div>
			<div class="email-footer">
				<p>&copy; %d BeRealTime. All rights reserved.</p>
				<p>Email ini dikirim secara otomatis, mohon jangan membalas email ini.</p>
			</div>
		</div>
	</div>
</body>
</html>
	`, subject, name, time.Now().Year())

	body := fmt.Sprintf(`
Halo %s,

Selamat datang di BeRealTime! Kami sangat senang Anda bergabung dengan komunitas kami.

Anda sekarang siap untuk mulai menjelajahi semua fitur yang kami tawarkan. Jika Anda memiliki pertanyaan atau memerlukan bantuan, jangan ragu untuk menghubungi tim dukungan kami.

Terima kasih,
Tim BeRealTime
`, name)

	return s.sendEmailHTML(to, subject, htmlBody, body)
}
