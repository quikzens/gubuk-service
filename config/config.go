package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
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
		log.Fatalf("Error loading .env file")
	}

	ServerAddress = os.Getenv("SERVER_ADDRESS")
	SecretKey = os.Getenv("SECRET_KEY")
	CloudinaryName = os.Getenv("CLOUDINARY_NAME")
	SecretKey = os.Getenv("CLOUDINARY_KEY")
	SecretKey = os.Getenv("CLOUDINARY_SECRET")
	DBSource = os.Getenv("DB_SOURCE")
}
