package files

import (
	"context"
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util"
	"node-backend/util/auth"
	"node-backend/util/requests"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/files/upload
func uploadFile(c *fiber.Ctx) error {

	if disabled {
		return requests.FailedRequest(c, "file.disabled", nil)
	}

	// Form data
	key := c.FormValue("key", "-")
	name := c.FormValue("name", "-")
	extension := c.FormValue("extension", "-")
	if key == "-" || name == "-" || extension == "-" {
		return requests.InvalidRequest(c)
	}
	file, err := c.FormFile("file")
	if err != nil {
		return requests.InvalidRequest(c)
	}
	accId := util.GetAcc(c)
	fileType := file.Header.Get("Content-Type")
	if fileType == "" {
		return requests.InvalidRequest(c)
	}

	// Check file size
	if file.Size > maxUploadSize {
		return requests.FailedRequest(c, "file.too_large", nil)
	}

	// Check total storage
	totalStorage, err := CountTotalStorage(accId)
	if err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	if totalStorage+file.Size > maxTotalStorage {
		return requests.FailedRequest(c, "file.storage_limit", nil)
	}

	// Generate file name Format: a-[accountId]-[objectIdentifier].[extension]
	fileId := "a-" + accId + "-" + auth.GenerateToken(16) + "." + extension
	if err := database.DBConn.Create(&account.CloudFile{
		Id:       fileId,
		Name:     name,
		Type:     file.Header.Get("Content-Type"),
		Key:      key,
		Account:  accId,
		Favorite: false,
		Size:     file.Size,
	}).Error; err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	f, err := file.Open()
	if err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	// Upload to R2
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileId),
		Body:   f,
		ACL:    "public-read",
	})
	if err != nil {
		return requests.FailedRequest(c, "server.error", err)
	}

	// Update file path
	if err := database.DBConn.Model(&account.CloudFile{}).Where("id = ?", fileId).Update("path", result.Location).Error; err != nil {
		client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(file.Filename),
		})
		return requests.FailedRequest(c, "server.error", err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"id":      fileId,
		"url":     result.Location,
	})
}
