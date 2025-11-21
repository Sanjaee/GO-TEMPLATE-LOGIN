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
	subject := "Kode Verifikasi OTP Anda"

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
		.email-header { background: linear-gradient(135deg, black 0%%, darkslategray 100%%); padding: 40px 30px; text-align: center; }
		.email-header h1 { color: white; font-size: 28px; font-weight: 700; margin: 0; letter-spacing: -0.5px; }
		.email-body { padding: 40px 30px; }
		.email-body p { margin: 0 0 20px 0; font-size: 16px; color: darkslategray; }
		.otp-container { background-color: whitesmoke; border: 2px dashed silver; border-radius: 12px; padding: 30px; margin: 30px 0; text-align: center; }
		.otp-code { font-size: 36px; font-weight: 700; color: black; letter-spacing: 8px; font-family: 'Courier New', monospace; margin: 10px 0; }
		.otp-label { font-size: 14px; color: gray; text-transform: uppercase; letter-spacing: 1px; margin-bottom: 10px; }
		.warning-box { background-color: snow; border-left: 4px solid darkred; padding: 15px; margin: 25px 0; border-radius: 4px; }
		.warning-box p { margin: 0; font-size: 14px; color: darkred; font-weight: 600; }
		.email-footer { background-color: whitesmoke; padding: 25px 30px; text-align: center; border-top: 1px solid gainsboro; }
		.email-footer p { margin: 5px 0; font-size: 12px; color: gray; }
		@media only screen and (max-width: 600px) {
			.email-wrapper { padding: 20px 10px; }
			.email-header { padding: 30px 20px; }
			.email-header h1 { font-size: 24px; }
			.email-body { padding: 30px 20px; }
			.otp-code { font-size: 28px; letter-spacing: 6px; }
		}
	</style>
</head>
<body>
	<div class="email-wrapper">
		<div class="email-container">
			<div class="email-header">
				<h1>Kode Verifikasi OTP</h1>
			</div>
			<div class="email-body">
				<p>Halo,</p>
				<p>Terima kasih telah menggunakan layanan kami. Gunakan kode verifikasi di bawah ini untuk melanjutkan proses Anda:</p>
				
				<div class="otp-container">
					<div class="otp-label">Kode Verifikasi</div>
					<div class="otp-code">%s</div>
				</div>
				
				<div class="warning-box">
					<p>⚠️ Kode ini berlaku selama 10 menit. Jangan bagikan kode ini kepada siapapun.</p>
				</div>
				
				<p style="font-size: 14px; color: gray;">Jika Anda tidak meminta kode verifikasi ini, silakan abaikan email ini atau hubungi tim support kami.</p>
			</div>
			<div class="email-footer">
				<p>&copy; %d BeRealTime. All rights reserved.</p>
				<p>Email ini dikirim secara otomatis, mohon jangan membalas email ini.</p>
			</div>
		</div>
	</div>
</body>
</html>
`, subject, otpCode, time.Now().Year())

	// Plain text fallback
	textBody := fmt.Sprintf(`
Halo,

Gunakan kode verifikasi berikut untuk memverifikasi email Anda:

Kode OTP: %s

Kode ini akan kedaluwarsa dalam 10 menit.

Demi keamanan, jangan pernah membagikan kode ini kepada siapapun.

Jika Anda tidak meminta kode ini, abaikan email ini.

Terima kasih,
Tim YouApp
`, otpCode)

	return s.sendEmailHTML(to, subject, htmlBody, textBody)
}

func (s *emailService) SendResetPasswordEmail(to, resetLink string) error {
	subject := "Reset Password - Kode OTP Anda"

	// Template HTML Profesional Modern - Tanpa RGB (Untuk Reset Password via OTP)
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
		.email-header { background: linear-gradient(135deg, black 0%%, darkslategray 100%%); padding: 40px 30px; text-align: center; }
		.email-header h1 { color: white; font-size: 28px; font-weight: 700; margin: 0; letter-spacing: -0.5px; }
		.email-body { padding: 40px 30px; }
		.email-body p { margin: 0 0 20px 0; font-size: 16px; color: darkslategray; }
		.info-box { background-color: whitesmoke; border-left: 4px solid black; padding: 20px; margin: 25px 0; border-radius: 4px; }
		.info-box p { margin: 0; font-size: 14px; color: darkslategray; }
		.email-footer { background-color: whitesmoke; padding: 25px 30px; text-align: center; border-top: 1px solid gainsboro; }
		.email-footer p { margin: 5px 0; font-size: 12px; color: gray; }
		@media only screen and (max-width: 600px) {
			.email-wrapper { padding: 20px 10px; }
			.email-header { padding: 30px 20px; }
			.email-header h1 { font-size: 24px; }
			.email-body { padding: 30px 20px; }
		}
	</style>
</head>
<body>
	<div class="email-wrapper">
		<div class="email-container">
			<div class="email-header">
				<h1>Reset Password</h1>
			</div>
			<div class="email-body">
				<p>Halo,</p>
				<p>Kami menerima permintaan untuk mereset password akun Anda. Kode OTP reset password telah dikirim ke email ini.</p>
				
				<div class="info-box">
					<p><strong>Langkah selanjutnya:</strong></p>
					<p style="margin-top: 10px;">1. Buka aplikasi atau website kami</p>
					<p>2. Masukkan kode OTP yang akan dikirim dalam email terpisah</p>
					<p>3. Buat password baru Anda</p>
				</div>
				
				<p style="font-size: 14px; color: gray;">Jika Anda tidak meminta reset password ini, silakan abaikan email ini. Akun Anda tetap aman.</p>
			</div>
			<div class="email-footer">
				<p>&copy; %d BeRealTime. All rights reserved.</p>
				<p>Email ini dikirim secara otomatis, mohon jangan membalas email ini.</p>
			</div>
		</div>
	</div>
</body>
</html>
`, subject, time.Now().Year())

	textBody := fmt.Sprintf(`
Halo,

Kami menerima permintaan untuk mereset password akun Anda.

Kode OTP reset password akan dikirim melalui email terpisah. Silakan cek email Anda dan ikuti instruksi untuk mereset password.

Jika Anda tidak meminta reset password ini, silakan abaikan email ini.

Terima kasih,
Tim BeRealTime
`)

	return s.sendEmailHTML(to, subject, htmlBody, textBody)
}

func (s *emailService) SendVerificationEmail(to, token string) error {
	subject := "Verifikasi Email Anda"
	verificationURL := fmt.Sprintf("%s/auth/verify-email?token=%s", s.config.ClientURL, token)
	
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
		.email-header { background: linear-gradient(135deg, black 0%%, darkslategray 100%%); padding: 40px 30px; text-align: center; }
		.email-header h1 { color: white; font-size: 28px; font-weight: 700; margin: 0; letter-spacing: -0.5px; }
		.email-body { padding: 40px 30px; }
		.email-body p { margin: 0 0 20px 0; font-size: 16px; color: darkslategray; }
		.button-container { text-align: center; margin: 30px 0; }
		.button { display: inline-block; padding: 16px 40px; background-color: black; color: white !important; text-decoration: none; border-radius: 8px; font-weight: 600; font-size: 16px; letter-spacing: 0.5px; transition: all 0.3s ease; }
		.button:hover { background-color: darkslategray; transform: translateY(-2px); box-shadow: 0 4px 12px rgba(0,0,0,0.2); }
		.link-box { background-color: whitesmoke; border: 1px solid gainsboro; border-radius: 8px; padding: 15px; margin: 20px 0; word-break: break-all; font-size: 13px; color: gray; }
		.info-box { background-color: snow; border-left: 4px solid darkorange; padding: 15px; margin: 25px 0; border-radius: 4px; }
		.info-box p { margin: 0; font-size: 14px; color: darkslategray; }
		.email-footer { background-color: whitesmoke; padding: 25px 30px; text-align: center; border-top: 1px solid gainsboro; }
		.email-footer p { margin: 5px 0; font-size: 12px; color: gray; }
		@media only screen and (max-width: 600px) {
			.email-wrapper { padding: 20px 10px; }
			.email-header { padding: 30px 20px; }
			.email-header h1 { font-size: 24px; }
			.email-body { padding: 30px 20px; }
			.button { display: block; width: 90%%; margin: 10px auto; }
		}
	</style>
</head>
<body>
	<div class="email-wrapper">
		<div class="email-container">
			<div class="email-header">
				<h1>Verifikasi Email</h1>
			</div>
			<div class="email-body">
				<p>Halo,</p>
				<p>Terima kasih telah mendaftar di BeRealTime! Untuk mengaktifkan akun Anda, silakan verifikasi alamat email dengan mengklik tombol di bawah ini:</p>
				
				<div class="button-container">
					<a href="%s" class="button">Verifikasi Email Saya</a>
				</div>
				
				<p style="font-size: 14px; color: gray; text-align: center;">Jika tombol di atas tidak berfungsi, salin dan tempel link berikut ke browser Anda:</p>
				
				<div class="link-box">
					%s
				</div>
				
				<div class="info-box">
					<p>⚠️ <strong>Penting:</strong> Link ini berlaku selama 24 jam. Setelah itu, Anda perlu meminta link verifikasi baru.</p>
				</div>
				
				<p style="font-size: 14px; color: gray;">Jika Anda tidak mendaftar untuk akun ini, silakan abaikan email ini.</p>
			</div>
			<div class="email-footer">
				<p>&copy; %d BeRealTime. All rights reserved.</p>
				<p>Email ini dikirim secara otomatis, mohon jangan membalas email ini.</p>
			</div>
		</div>
	</div>
</body>
</html>
	`, subject, verificationURL, verificationURL, time.Now().Year())

	textBody := fmt.Sprintf(`
Halo,

Terima kasih telah mendaftar. Klik link berikut untuk memverifikasi email Anda:

%s

Link ini akan kedaluwarsa dalam 24 jam.

Jika Anda tidak meminta verifikasi ini, abaikan email ini.

Terima kasih,
Tim YouApp
`, verificationURL)

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
