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

var accessTokenPublicKey = fmt.Sprintf("%v", os.Getenv("ACCESS_TOKEN_PUBLIC_KEY"))

func SendEmail(receiver string, email string, verificationLink string, deletionLink string) error {
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
		Warning string
	}{
		Name:    receiver,
		Message: "Вашият email адрес трябва да бъде потвърден, за да можете да се регистрирате за Hacktues 9. Моля, натиснете върху линка по-долу, за да потвърдите вашият email адрес : " + verificationLink,
		Warning: "Ако не сте регистрирали се за Hacktues 9, моля натиснете върху линка по-долу, за да отмените регистрацията си : " + deletionLink,
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

func GenerateResetLink(email string, privateKey string, publicKey string, TokenTTL time.Duration) string {
	routeUrl := os.Getenv("ROUTE_URL")
	token, err := jwt.CreateToken(TokenTTL, email, privateKey, publicKey)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return routeUrl + "forgotten-password/reset/" + token
}

func ValidateEmailToken(token string) (string, error) {

	sub, err := jwt.ValidateStringToken(token, accessTokenPublicKey)
	if err != nil {
		fmt.Println("err: ", err)
		return "", err
	}
	return sub, nil
}

func ValidateEmail(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	vars := mux.Vars(r)
	token := vars["token"]
	elsys, err := strconv.ParseBool(vars["elsys"])
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		http.Redirect(w, r, os.Getenv("ROUTE_URL"), http.StatusNotFound)
		return
	}

	email, err := ValidateEmailToken(token)
	if err != nil {
		fmt.Printf("[ ERROR ] [ ValidateEmail ] validate: token: %s", err)
		http.Redirect(w, r, os.Getenv("ROUTE_URL"), http.StatusNotFound)
		return
	}
	security := models.Security{}

	if elsys {
		db.Table("security").Joins("JOIN users ON users.security_id = security.id").Where("users.elsys_email = ?", email).First(&security)

		if security.ElsysEmailVerified {
			w.WriteHeader(http.StatusBadRequest)
			http.Redirect(w, r, os.Getenv("ROUTE_URL"),http.StatusSeeOther)
			return
		}
		security.ElsysEmailVerified = true
		db.Model(models.Security{}).Where("ID = ?", security.ID).Update("elsys_email_verified", security.ElsysEmailVerified)
	} else {
		db.Table("security").Joins("JOIN users ON users.security_id = security.id").Where("users.email = ?", email).First(&security)

		if security.EmailVerified {
			w.WriteHeader(http.StatusBadRequest)
			http.Redirect(w, r, os.Getenv("ROUTE_URL"), http.StatusSeeOther)
			return
		}
		security.EmailVerified = true
		db.Model(models.Security{}).Where("ID = ?", security.ID).Update("email_verified", security.EmailVerified)
	}

	http.Redirect(w, r, os.Getenv("ROUTE_URL"), http.StatusSeeOther)
}

func SendResetLink(reciever string, email string, resetLink string) error {
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
		Message: "Please reset your password by clicking the following link : " + resetLink,
	})

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func ValidateResetLink(token string) (string, error) {

	sub, err := jwt.ValidateStringToken(token, accessTokenPublicKey)
	if err != nil {
		fmt.Println("err: ", err)
		return "", err
	}
	return sub, nil
}

func GenerateDeletionLink(email string, privateKey string, publicKey string, TokenTTL time.Duration) string {
	hostUrl := os.Getenv("HOST_URL")
	token, err := jwt.CreateToken(TokenTTL, email, privateKey, publicKey)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return hostUrl + "api/auth/delete/" + token
}
