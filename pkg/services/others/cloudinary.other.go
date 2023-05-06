package others

import (
	"context"
	"glitchz/pkg/utils"
	"log"
	"time"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"
)

func InitCloud() *cloudinary.Cloudinary {
	config, err := utils.LoadConfig("../")
	if err != nil {
		log.Fatal("Error: cannot load config file", err)
	}
	cld, _ := cloudinary.NewFromParams(config.CloudinaryName, config.CloudinaryKey, config.CloudinarySecret)
	return cld
}

func UploadToCloud(ctx *gin.Context) (string, string, error) {
	fileName := ctx.PostForm("glitchz-upload" + time.Now().String())
	fileTags := ctx.PostForm("tags")
	file, _, err := ctx.Request.FormFile("upload")

	if err != nil {
		return "", "", err
	}

	res, err := InitCloud().Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID:    fileName,
		AutoTagging: ctx.GetFloat64(fileTags),
	})

	if err != nil {
		return "", "", err
	}

	return res.SecureURL, res.PublicID, nil
}

func DeleteFromCloud(publicId string, ctx context.Context) error {
	_, err := InitCloud().Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicId,
	})

	if err != nil {
		return err
	}

	return nil
}
