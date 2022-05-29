package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	Port             string
	ServerAddress    string
	SecretKey        string
	CloudinaryName   string
	CloudinaryKey    string
	CloudinarySecret string
	DBSource         string
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Couldn't loading .env file")
	}

	Port = os.Getenv("PORT")
	ServerAddress = os.Getenv("SERVER_ADDRESS")
	SecretKey = os.Getenv("SECRET_KEY")
	CloudinaryName = os.Getenv("CLOUDINARY_NAME")
	CloudinaryKey = os.Getenv("CLOUDINARY_KEY")
	CloudinarySecret = os.Getenv("CLOUDINARY_SECRET")
	DBSource = os.Getenv("DB_SOURCE")
}
