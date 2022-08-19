package userHandler

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hacktues-9/API/db"
	"github.com/hacktues-9/API/modelUser"
)

func login(w http.ResponseWriter, r *http.Request) {
	var user modelUser.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println(err)
	}
	defer r.Body.Close()

	user, err := modelUser.GetUserByEmailAndPassword(user.Email, user.Password)
	
	if err != nil {
		fmt.Println(err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.Id,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		fmt.Println(err)
	}
	w.Write([]byte(tokenString))
}