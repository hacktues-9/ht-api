package teams

import (
	"github.com/hacktues-9/API/pkg/models"
	"gorm.io/gorm"
)

// var (
// 	accessTokenTTL         = time.Hour * 24
// 	refreshTokenTTL        = time.Hour * 24 * 7
// 	accessTokenPrivateKey  = os.Getenv("ACCESS_TOKEN_PRIVATE_KEY")
// 	accessTokenPublicKey   = os.Getenv("ACCESS_TOKEN_PUBLIC_KEY")
// 	refreshTokenPrivateKey = os.Getenv("REFRESH_TOKEN_PRIVATE_KEY")
// 	refreshTokenPublicKey  = os.Getenv("REFRESH_TOKEN_PUBLIC_KEY")
// )

func AddTeamMember(user *models.Users, team *models.Team, db *gorm.DB) error {
	//check if team is full max is 5 members
	return nil
}
