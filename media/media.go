package media

import (
	"context"
	"gubuk-service/config"
	"log"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

var cld *cloudinary.Cloudinary

func init() {
	cloudinaryClient, err := cloudinary.NewFromParams(config.CloudinaryName, config.CloudinaryKey, config.CloudinarySecret)
	if err != nil {
		log.Fatal(err)
	}

	cld = cloudinaryClient
}

func UploadMedia(folder string, media *multipart.FileHeader) (string, error) {
	file, err := media.Open()
	if err != nil {
		return "", err
	}

	uploadResult, err := cld.Upload.Upload(context.TODO(), file, uploader.UploadParams{
		Folder: folder,
	})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}

func DestroyMedia(destroyedMediaUrl string) error {
	publicId := extractPublicId(destroyedMediaUrl)

	_, err := cld.Upload.Destroy(context.TODO(), uploader.DestroyParams{
		PublicID: publicId,
	})
	if err != nil {
		return err
	}

	return nil
}

func UpdateMedia(folder string, destroyedMediaUrl string, media *multipart.FileHeader) (string, error) {
	newMediaUrl, err := UploadMedia(folder, media)
	if err != nil {
		return "", err
	}

	err = DestroyMedia(destroyedMediaUrl)
	if err != nil {
		return "", err
	}

	return newMediaUrl, nil
}
