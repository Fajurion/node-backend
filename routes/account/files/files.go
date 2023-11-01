package files

import (
	"context"
	"log"
	"node-backend/database"
	"node-backend/entities/account"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
)

var bucketName string
var client *s3.Client
var uploader *manager.Uploader
var disabled = false

// Configuration
const minUploadSize = 1_000            // 1 KB
const maxUploadSize = 10_000_000       // 10 MB
const maxFavoriteStorage = 500_000_000 // 500 MB
const maxTotalStorage = 1_000_000_000  // 1 GB

func Authorized(router fiber.Router) {
	url := os.Getenv("R2_URL")
	bucketName = os.Getenv("R2_BUCKET")
	accessKeyId := os.Getenv("R2_CLIENT_ID")
	accessKeySecret := os.Getenv("R2_CLIENT_SECRET")

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: url,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
	)
	if err != nil {
		disabled = true
		log.Println("Failed to connect to R2. File integration disabled.")
		log.Fatal(err)
	}

	// Setup uploader
	client = s3.NewFromConfig(cfg)
	uploader = manager.NewUploader(client)
	log.Println("Successfully connected to R2.")

	// Setup file routes
	router.Post("/upload", uploadFile)
	router.Post("/delete", deleteFile)
	router.Post("/list", listFiles)
	router.Post("/favorite", favoriteFile)
	router.Post("/unfavorite", unfavoriteFile)
}

func CountTotalStorage(accId string) (int64, error) {

	// Get total storage
	var totalStorage int64
	unix := time.Now().Add(-time.Hour * 24 * 30).UnixMilli()
	if err := database.DBConn.Model(&account.CloudFile{}).Where("account = ? AND (created_at > ? OR favorite = ?)", accId, unix, true).Select("sum(size)").Scan(&totalStorage).Error; err != nil {
		return 0, err
	}

	return totalStorage, nil
}

func CountFavoriteStorage(accId string) (int64, error) {

	// Get total storage
	var favoriteStorage int64
	if err := database.DBConn.Model(&account.CloudFile{}).Where("account = ? AND favorite = ?", accId, true).Select("sum(size)").Scan(&favoriteStorage).Error; err != nil {
		return 0, err
	}

	return favoriteStorage, nil
}
