package media

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"strings"
)

func filenameWithoutExt(filename string) string {
	return filename[:len(filename)-len(filepath.Ext(filename))]
}

// extractPublicId return extracted media publicId from url
// extractPublicId("https://res.cloudinary.com/example_cloud/image/upload/v12345678/folder/example.png") -> "folder/example"
func extractPublicId(url string) string {
	separatedUrl := strings.Split(url, "/")
	separatedUrl = separatedUrl[len(separatedUrl)-2:]
	separatedUrl[1] = filenameWithoutExt(separatedUrl[1])
	return strings.Join(separatedUrl, "/")
}

func ValidateImage(image *multipart.FileHeader) error {
	// validate image extension
	var imageType, _ = regexp.Compile(`^.*\.(jpeg|JPEG|jpg|JPG|gif|GIF|png|PNG|svg|SVG|webp|WebP|WEBP)$`)
	if isImage := imageType.MatchString(image.Filename); !isImage {
		return errors.New("invalid file type")
	}

	// validate image size, max 1 MB (1048576 bytes)
	if image.Size > 1048576 {
		return errors.New("invalid file size")
	}

	return nil
}
