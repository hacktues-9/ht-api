package email

import (
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/hacktues-9/API/pkg/jwt"
	"github.com/hacktues-9/API/pkg/models"
	"gorm.io/gorm"
)

func SendEmail(reciever string, email string, verificationLink string) error {
	auth := smtp.PlainAuth(
		"",
		"hacktues@elsys-bg.org",
		os.Getenv("EMAIL_PASSWORD"),
		"smtp.gmail.com",
	)

	//verification email
	msg := "Subject: " + verificationLink

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"hacktues@elsys-bg.org",
		[]string{email},
		[]byte(msg),
	)
	if err != nil {
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
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := models.Users{}
	if result := db.Preload("Security").Where("email = ?", email).First(&user); result.Error != nil {
		if result := db.Preload("Security").Where("elsys_email = ?", email).First(&user); result.Error != nil {
			fmt.Println("login: user: find:", result.Error)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("login: user: find: " + result.Error.Error()))
			return
		}
	}
	if elsys {
		if user.Security.ElsysEmailVerified {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		user.Security.ElsysEmailVerified = true
		db.Model(models.Security{}).Where("ID = ?", user.SecurityID).Update("elsys_email_verified", user.Security.ElsysEmailVerified)
	} else {
		if user.Security.EmailVerified {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		user.Security.EmailVerified = true
		db.Model(models.Security{}).Where("ID = ?", user.SecurityID).Update("email_verified", user.Security.EmailVerified)
	}

	w.WriteHeader(http.StatusOK)
}
