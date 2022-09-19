package db

import (
	"log"

	"github.com/hacktues-9/API/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func Init() *gorm.DB {
	dsn := "host=localhost user=postgres password=password dbname=ht92 port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.Class{})
	db.AutoMigrate(&models.Role{})
	db.AutoMigrate(&models.EatingPreference{})
	db.AutoMigrate(&models.Allergies{})
	db.AutoMigrate(&models.ShirtSize{})
	db.AutoMigrate(&models.Technologies{})
	db.AutoMigrate(&models.Discord{})
	db.AutoMigrate(&models.Socials{})
	db.AutoMigrate(&models.Info{})
	db.AutoMigrate(&models.Security{})
	db.AutoMigrate(&models.Project{})
	db.AutoMigrate(&models.Pictures{})
	db.AutoMigrate(&models.Invites{})
	db.AutoMigrate(&models.Team{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Invite{})
	db.AutoMigrate(&models.InfoAllergies{})
	db.AutoMigrate(&models.UserTechnologies{})
	db.AutoMigrate(&models.TeamTechnologies{})
	db.AutoMigrate(&models.ProjectTechnologies{})
	db.AutoMigrate(&models.InviteTechnologies{})
	db.AutoMigrate(&models.InvitesTechnologies{})
	db.AutoMigrate(&models.Log{})

	return db
}
