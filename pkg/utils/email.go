package utils

  import (
        "bytes"
        "crypto/tls"
        "errors"
        "fmt"
        "log"
        "net"
        "net/smtp"
        "strings"
        "text/template"
        "time"

        "github.com/zyj/my-blog/pkg/constant"
  )

  func RenderVerificationEmail(code string) (string, error) {
        tmpl, err := template.ParseFiles(
                "internal/static/assets/email/verification-code.tmpl",
        )
        if err != nil {
                return "", err
        }

        var body bytes.Buffer
        err = tmpl.Execute(
                &body,
                struct {
                        Code string
                }{
                        Code: code,
                },
        )
        if err != nil {
                return "", err
        }

        return body.String(), nil
  }

  func SendEmail(to, subject, body string) error {
        if !GetAsBool(constant.EnvKeyEmailEnable, false) {
                mode := constant.Mode(
                        Get(
                                constant.EnvKeyMode,
                                string(constant.ModeDev),
                        ),
                )

                if mode == constant.ModeDev {
                        log.Printf(
                                "email disabled: to=%s subject=%s body=%s",
                                to,
                                subject,
                                body,
                        )
                        return nil
                }

                return errors.New("email is disabled")
        }

        host := Get(constant.EnvKeyEmailHost)
        port := GetAsInt(constant.EnvKeyEmailPort, 465)
        username := Get(constant.EnvKeyEmailUsername)
        password := Get(constant.EnvKeyEmailPassword)
        from := Get(constant.EnvKeyEmailAddress)

        if from == "" {
                from = username
        }

        if host == "" || from == "" {
                return errors.New("email configuration is incomplete")
        }

        address := net.JoinHostPort(
                host,
                fmt.Sprintf("%d", port),
        )

        message := buildEmailMessage(
                from,
                to,
                subject,
                body,
        )

        var auth smtp.Auth
        if username != "" {
                auth = smtp.PlainAuth(
                        "",
                        username,
                        password,
                        host,
                )
        }

        if GetAsBool(constant.EnvKeyEmailSSL, true) {
                return sendEmailWithSSL(
                        address,
                        host,
                        auth,
                        from,
                        to,
                        message,
                )
        }

        return smtp.SendMail(
                address,
                auth,
                from,
                []string{to},
                message,
        )
  }

   func buildEmailMessage(
        from,
        to,
        subject,
        body string,
  ) []byte {
	return fmt.Appendf(
		nil,
		"From: %s\r\n"+
                        "To: %s\r\n"+
                        "Subject: %s\r\n"+
                        "MIME-Version: 1.0\r\n"+
                        "Content-Type: text/html; charset=UTF-8\r\n"+
                        "\r\n%s",
                from,
                to,
                sanitizeEmailHeader(subject),
		body,
	)
}

  func sanitizeEmailHeader(value string) string {
        return strings.NewReplacer(
                "\r", "",
                "\n", "",
        ).Replace(value)
  }

  func sendEmailWithSSL(
        address string,
        host string,
        auth smtp.Auth,
        from string,
        to string,
        message []byte,
  ) error {
        connection, err := tls.DialWithDialer(
                &net.Dialer{
                        Timeout: 10 * time.Second,
                },
                "tcp",
                address,
                &tls.Config{
                        ServerName: host,
                        MinVersion: tls.VersionTLS12,
                },
        )
        if err != nil {
                return err
        }

        client, err := smtp.NewClient(connection, host)
        if err != nil {
                _ = connection.Close()
                return err
        }
        defer client.Close()

        if auth != nil {
                if err := client.Auth(auth); err != nil {
                        return err
                }
        }

        if err := client.Mail(from); err != nil {
                return err
        }

        if err := client.Rcpt(to); err != nil {
                return err
        }

        writer, err := client.Data()
        if err != nil {
                return err
        }

        if _, err := writer.Write(message); err != nil {
                _ = writer.Close()
                return err
        }

        if err := writer.Close(); err != nil {
                return err
        }

        return client.Quit()
  }
