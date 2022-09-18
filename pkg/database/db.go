package db

import (
	"log"

	"github.com/hacktues-9/API/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func Init() *gorm.DB {
	dsn := "host=localhost user=postgres password=password dbname=ht9 port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	UserDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{TablePrefix: "users.", SingularTable: true}})
	if err != nil {
		log.Fatal(err)
	}
	TeamDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{TablePrefix: "teams.", SingularTable: true}})
	if err != nil {
		log.Fatal(err)
	}
	LogDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{TablePrefix: "logs.", SingularTable: true}})
	if err != nil {
		log.Fatal(err)
	}

	UserDB.Migrator().CreateTable(&models.Class{})
	UserDB.Migrator().CreateTable(&models.Role{})
	UserDB.Migrator().CreateTable(&models.EatingPreference{})
	UserDB.Migrator().CreateTable(&models.Allergies{})
	UserDB.Migrator().CreateTable(&models.ShirtSize{})
	db.Migrator().CreateTable(&models.Technologies{})
	UserDB.Migrator().CreateTable(&models.Discord{})
	UserDB.Migrator().CreateTable(&models.Socials{})
	UserDB.Migrator().CreateTable(&models.Info{})
	UserDB.Migrator().CreateTable(&models.Security{})
	TeamDB.Migrator().CreateTable(&models.Project{})
	TeamDB.Migrator().CreateTable(&models.Pictures{})
	TeamDB.Migrator().CreateTable(&models.Invites{})
	TeamDB.Migrator().CreateTable(&models.Team{})
	UserDB.Migrator().CreateTable(&models.User{})
	TeamDB.Migrator().CreateTable(&models.Invite{})
	UserDB.Migrator().CreateTable(&models.InfoAllergies{})
	TeamDB.Migrator().CreateTable(&models.TeamTechnologies{})
	TeamDB.Migrator().CreateTable(&models.ProjectTechnologies{})
	TeamDB.Migrator().CreateTable(&models.InviteTechnologies{})
	TeamDB.Migrator().CreateTable(&models.InvitesTechnologies{})
	LogDB.Migrator().CreateTable(&models.Log{})

	// LogDB.Migrator().DropTable(&models.Log{})
	// TeamDB.Migrator().DropTable(&models.InvitesTechnologies{})
	// TeamDB.Migrator().DropTable(&models.InviteTechnologies{})
	// TeamDB.Migrator().DropTable(&models.ProjectTechnologies{})
	// TeamDB.Migrator().DropTable(&models.TeamTechnologies{})
	// UserDB.Migrator().DropTable(&models.InfoAllergies{})
	// TeamDB.Migrator().DropTable(&models.Invite{})
	// UserDB.Migrator().DropTable(&models.User{})
	// TeamDB.Migrator().DropTable(&models.Team{})
	// TeamDB.Migrator().DropTable(&models.Invites{})
	// TeamDB.Migrator().DropTable(&models.Pictures{})
	// TeamDB.Migrator().DropTable(&models.Project{})
	// UserDB.Migrator().DropTable(&models.Security{})
	// UserDB.Migrator().DropTable(&models.Info{})
	// UserDB.Migrator().DropTable(&models.Socials{})
	// UserDB.Migrator().DropTable(&models.Discord{})
	// db.Migrator().DropTable(&models.Technologies{})
	// UserDB.Migrator().DropTable(&models.ShirtSize{})
	// UserDB.Migrator().DropTable(&models.Allergies{})
	// UserDB.Migrator().DropTable(&models.EatingPreference{})
	// UserDB.Migrator().DropTable(&models.Role{})
	// UserDB.Migrator().DropTable(&models.Class{})

	return db
}
