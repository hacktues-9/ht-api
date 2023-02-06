package email

import (
	"bytes"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/hacktues-9/API/pkg/jwt"
	"github.com/hacktues-9/API/pkg/models"
	"gorm.io/gorm"
)

func SendEmail(reciever string, email string, verificationLink string) error {
	from := "hacktues@elsys-bg.org"
	password := os.Getenv("EMAIL_PASSWORD")

	to := []string{
		email,
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	filePrefix, _ := filepath.Abs("./pkg/email/")
	temp, err := template.ParseFiles(filePrefix + "/email.html")
	if err != nil {
		fmt.Println(err)
		return err
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Verify your email for Hacktues 9!\n%s\n\n", mimeHeaders)))

	temp.Execute(&body, struct {
		Name    string
		Message string
	}{
		Name:    reciever,
		Message: "Please verify your email by clicking the following link : " + verificationLink,
	})

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func GenerateVerificationLink(email string, privateKey string, publicKey string, TokenTTL time.Duration) string {
	hostUrl := os.Getenv("HOST_URL")
	elsys := strings.Contains(email, "@elsys-bg.org")
	token, err := jwt.CreateToken(TokenTTL, email, privateKey, publicKey)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return hostUrl + "api/user/verify/" + strconv.FormatBool(elsys) + "/" + token
}

func ValidateEmailToken(token string, publicKey string) (string, error) {
	fmt.Println("token: ", token)
	claims, err := jwt.ValidateToken(token, publicKey)
	if err != nil {
		return "", fmt.Errorf("ValidateEmailToken: %w", err)
	}
	return strconv.FormatUint(uint64(claims), 10), nil
}

func ValidateEmail(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	vars := mux.Vars(r)
	token := vars["token"]
	elsys, err := strconv.ParseBool(vars["elsys"])
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	email, err := ValidateEmailToken(token, os.Getenv("ACCESS_TOKEN_PUBLIC_KEY"))
	if err != nil {
		fmt.Printf("[ ERROR ] [ ValidateEmail ] validate: token: %s", err)
		models.RespHandler(w, r, models.DefaultNegResponse(http.StatusUnauthorized, "validate: token", 0), err, http.StatusUnauthorized, "ValidateEmail")
		return
	}
	security := models.Security{}

	if elsys {
		db.Table("security").Joins("JOIN users ON users.security_id = security.id").Where("users.elsys_email = ?", email).First(&security)

		if security.ElsysEmailVerified {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		security.ElsysEmailVerified = true
		db.Model(models.Security{}).Where("ID = ?", security.ID).Update("elsys_email_verified", security.ElsysEmailVerified)
	} else {
		db.Table("security").Joins("JOIN users ON users.security_id = security.id").Where("users.email = ?", email).First(&security)

		if security.EmailVerified {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		security.EmailVerified = true
		db.Model(models.Security{}).Where("ID = ?", security.ID).Update("email_verified", security.EmailVerified)
	}

	w.WriteHeader(http.StatusOK)
}
