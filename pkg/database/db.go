package db

import (
	"log"
	"os"

	"github.com/hacktues-9/API/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func Init() *gorm.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("SSL_MODE")
	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=" + sslmode
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(&models.Class{})
	if err != nil {
		log.Fatal(err)
		return
	}
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
	db.AutoMigrate(&models.Users{})
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
	var roleValues = []string{"MEMBER", "CAPTAIN", "STUDENT", "MENTOR", "ADMIN"}
	var eatingPreferenceValues = []string{"VEGETARIAN", "VEGAN", "NONE"}
	var allergiesValues = []string{"EGGS", "NUTS", "MILK", "GLUTEN"}
	var shirtSizeValues = []string{"XS", "S", "M", "L", "XL", "XXL"}
	var technologiesValues = []string{"React",
		"Angular",
		"Vue",
		"Node",
		"Express",
		"MongoDB",
		"MySQL",
		"PostgreSQL",
		"Docker",
		"Kubernetes",
		"AWS",
		"Azure",
		"Google Cloud",
		"Git",
		"C",
		"C++",
		"C#",
		"Java",
		"Python",
		"PHP",
		"HTML",
		"CSS",
		"JavaScript",
		"TypeScript",
		"NoSQL",
		"REST",
		"Ruby",
		"Ruby on Rails",
		"Swift",
		"Kotlin",
		"Go",
		"Rust",
		"Elixir",
		"Svelte",
		"Bash",
		"Vim",
		"Linux",
		"TensorFlow",
		"Apache",
		"Nginx",
		"Deno",
		"Next",
		"Jest",
		"Dart",
		"Flutter",
		"Firebase",
		"Vercel",
		"Netlify",
		"Material UI",
		"Tailwind",
		"SASS",
		"Webpack",
		"Babel",
		"JQuery",
		"Bootstrap",
	}

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
	db.Migrator().DropTable(&models.Users{})
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
