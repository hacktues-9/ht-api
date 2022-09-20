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
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func Migrate(db *gorm.DB) {
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
}

func PopulateDefault(db *gorm.DB) {
	var classValues = []string{"А", "Б", "В", "Г"}
	var roleValues = []string{"student", "member", "captain", "mentor", "tester", "admin"}
	var eatingPreferenceValues = []string{"default", "vegetarian", "vegan", "omnivore", "pescetarian", "other"}
	var allergiesValues = []string{"none", "nuts", "gluten", "dairy", "other"}
	var shirtSizeValues = []string{"XS", "S", "M", "L", "XL"}
	var technologiesValues = []string{"C", "C++", "C#", "Java", "Python", "JavaScript", "HTML", "CSS", "PHP", "Ruby", "Go", "Swift", "Kotlin", "Rust", "Scala", "Haskell", "Elixir", "Erlang", "Clojure", "Lua", "Perl", "R", "SQL", "Bash", "Assembly", "MySQL", "PostgreSQL", "MongoDB", "Redis", "Docker", "Kubernetes"}

	for _, class := range classValues {
		db.Create(&models.Class{Name: class})
	}

	for _, role := range roleValues {
		db.Create(&models.Role{Name: role})
	}

	for _, eatingPreference := range eatingPreferenceValues {
		db.Create(&models.EatingPreference{Name: eatingPreference})
	}

	for _, allergies := range allergiesValues {
		db.Create(&models.Allergies{Name: allergies})
	}

	for _, shirtSize := range shirtSizeValues {
		db.Create(&models.ShirtSize{Name: shirtSize})
	}

	for _, technologies := range technologiesValues {
		db.Create(&models.Technologies{Technology: technologies})
	}
}

func Drop(db *gorm.DB) {
	db.Migrator().DropTable(&models.Log{})
	db.Migrator().DropTable(&models.InvitesTechnologies{})
	db.Migrator().DropTable(&models.InviteTechnologies{})
	db.Migrator().DropTable(&models.ProjectTechnologies{})
	db.Migrator().DropTable(&models.TeamTechnologies{})
	db.Migrator().DropTable(&models.UserTechnologies{})
	db.Migrator().DropTable(&models.InfoAllergies{})
	db.Migrator().DropTable(&models.Invite{})
	db.Migrator().DropTable(&models.User{})
	db.Migrator().DropTable(&models.Team{})
	db.Migrator().DropTable(&models.Invites{})
	db.Migrator().DropTable(&models.Pictures{})
	db.Migrator().DropTable(&models.Project{})
	db.Migrator().DropTable(&models.Security{})
	db.Migrator().DropTable(&models.Info{})
	db.Migrator().DropTable(&models.Socials{})
	db.Migrator().DropTable(&models.Discord{})
	db.Migrator().DropTable(&models.Technologies{})
	db.Migrator().DropTable(&models.ShirtSize{})
	db.Migrator().DropTable(&models.Allergies{})
	db.Migrator().DropTable(&models.EatingPreference{})
	db.Migrator().DropTable(&models.Role{})
	db.Migrator().DropTable(&models.Class{})
}
