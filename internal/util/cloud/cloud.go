package cloudutil

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"sync"
	"tart-shop-manager/internal/common"
	imagemodel "tart-shop-manager/internal/entity/dtos/sql/image"
)

var (
	sess     *session.Session
	uploader *s3manager.Uploader
	initErr  error
	once     sync.Once
)

type Image struct {
	FileName   string
	FileBuffer []byte
}

type UploadResult struct {
	FileName string
	FileURL  string
	Error    error
}

func initUploader() {
	accessKey := os.Getenv("SUPABASE_ACCESS_KEY")
	secretKey := os.Getenv("SUPABASE_SECRET_KEY")
	endPoint := os.Getenv("SUPABASE_ENDPOINT") // Should be "https://snjwvkwtuboigdykuzcf.supabase.co/storage/v1/s3"
	region := "ap-southeast-1"                 // Update this to match the region from the image

	sess, initErr = session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:         aws.String(endPoint),
		S3ForcePathStyle: aws.Bool(true), // Use path-style addressing
	})

	if initErr != nil {
		initErr = fmt.Errorf("failed to create AWS session: %w", initErr)
		return
	}

	uploader = s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		u.PartSize = 10 * 1024 * 1024 // 10MB per part
		u.Concurrency = 5             // Number of concurrent parts per upload
	})
}

// UploadToS3 uploads one or multiple images to S3 concurrently, using context for cancellation

func UploadToS3(ctx context.Context, images []Image) ([]UploadResult, error) {
	once.Do(initUploader)
	if initErr != nil {
		return nil, common.ErrCloudConnectionFailed(initErr)
	}

	bucketName := os.Getenv("SUPABASE_BUCKET")
	supabaseURL := os.Getenv("SUPABASE_URL")

	// Channel to limit concurrency
	maxConcurrentUploads := 10 // Adjust this number as needed
	semaphore := make(chan struct{}, maxConcurrentUploads)

	// Channel to collect results
	resultsChan := make(chan UploadResult, len(images))

	var wg sync.WaitGroup

	for _, img := range images {
		select {
		case <-ctx.Done():
			// Context was cancelled before starting the upload
			wg.Add(1)
			go func(img Image) {
				defer wg.Done()
				resultsChan <- UploadResult{
					FileName: img.FileName,
					FileURL:  "",
					Error:    ctx.Err(),
				}
			}(img)
			continue
		case semaphore <- struct{}{}:
			// Acquired a slot, proceed with the upload
		}

		wg.Add(1)
		go func(img Image) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release the slot

			var uploadErr error

			// Check for file size limit (example: 100MB)
			if len(img.FileBuffer) > 100*1024*1024 {
				uploadErr = common.ErrCannotUploadFile("image", fmt.Errorf("file size: %d bytes", len(img.FileBuffer)))
			} else {
				if len(img.FileBuffer) < (5 * 1024 * 1024) { // Files smaller than 5MB
					svc := s3.New(sess)
					_, err := svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
						Bucket: aws.String(bucketName),
						Key:    aws.String(img.FileName),
						Body:   bytes.NewReader(img.FileBuffer),
					})
					if err != nil {
						uploadErr = common.ErrCannotUploadFile(img.FileName, err)
					}
				} else {
					_, err := uploader.UploadWithContext(ctx, &s3manager.UploadInput{
						Bucket: aws.String(bucketName),
						Key:    aws.String(img.FileName),
						Body:   bytes.NewReader(img.FileBuffer),
					})
					if err != nil {
						uploadErr = common.ErrCannotUploadFile(img.FileName, err)
					}
				}
			}

			fileURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", supabaseURL, bucketName, img.FileName)
			resultsChan <- UploadResult{
				FileName: img.FileName,
				FileURL:  fileURL,
				Error:    uploadErr,
			}
		}(img)
	}

	// Wait for all uploads to complete
	wg.Wait()
	close(resultsChan)

	// Collect results
	results := make([]UploadResult, 0, len(images))
	var anyError bool
	for result := range resultsChan {
		if result.Error != nil {
			anyError = true
		}
		results = append(results, result)
	}

	if anyError {
		return results, common.ErrCannotUploadFile("images", fmt.Errorf("one or more uploads failed"))
	}

	return results, nil
}

// UploadSingleImageToS3 uploads a single image to S3

func UploadSingleImageToS3(ctx context.Context, fileBuffer []byte, fileName string) (string, error) {
	images := []Image{
		{
			FileName:   fileName,
			FileBuffer: fileBuffer,
		},
	}
	results, err := UploadToS3(ctx, images)
	if err != nil {
		return "", err
	}
	result := results[0]
	return result.FileURL, result.Error
}

func DeleteSingleImageFromS3(ctx context.Context, fileName string) error {
	if sess == nil {
		return common.ErrCloudConnectionFailed(fmt.Errorf("AWS session not initialized"))
	}

	bucketName := os.Getenv("SUPABASE_BUCKET")
	svc := s3.New(sess)

	_, err := svc.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return common.ErrCannotDeleteFile(imagemodel.EntityName, err)
	}

	// Đợi cho đến khi object bị xóa
	err = svc.WaitUntilObjectNotExistsWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return common.ErrCannotDeleteFile(imagemodel.EntityName, err)
	}

	return nil
}

//
//package cloudutil
//
//// Uploader and S3 Client using SDK v2
//import (
//	"bytes"
//	"context"
//	"errors"
//	"fmt"
//	"github.com/aws/aws-sdk-go-v2/aws"
//	"github.com/aws/aws-sdk-go-v2/config"
//	"github.com/aws/aws-sdk-go-v2/credentials"
//	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
//	"github.com/aws/aws-sdk-go-v2/service/s3"
//	"github.com/aws/smithy-go"
//	smithyendpoints "github.com/aws/smithy-go/endpoints"
//	"github.com/joho/godotenv"
//	"log"
//	//"net/url"
//	"os"
//	"sync"
//	"tart-shop-manager/internal/common"
//	imagemodel "tart-shop-manager/internal/entity/dtos/sql/image"
//)
//
//var (
//	uploader *manager.Uploader
//	s3Client *s3.Client
//	initErr  error
//	once     sync.Once
//)
//
//// Image represents an image to upload
//type Image struct {
//	FileName   string
//	FileBuffer []byte
//}
//
//// UploadResult contains the result of the upload process
//type UploadResult struct {
//	FileName string
//	FileURL  string
//	Error    error
//}
//
//// custom resolver for EndpointResolverV2
//type resolverV2 struct{}
//
//func (*resolverV2) ResolveEndpoint(ctx context.Context, params s3.EndpointParameters) (
//	smithyendpoints.Endpoint, error,
//) {
//	// Print the endpoint provided in config
//	fmt.Printf("The endpoint provided in config is %s\n", *params.Endpoint)
//
//	// Fallback to default endpoint resolution
//	return s3.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
//}
//
//// initUploader initializes the S3 client and uploader
//func initUploader(ctx context.Context) error {
//	// Load environment variables from .env file
//	err := godotenv.Load()
//	if err != nil {
//		log.Fatalf("Error loading .env file: %v", err)
//	}
//
//	// Get credentials from environment variables
//	accessKey := os.Getenv("SUPABASE_ACCESS_KEY")
//	secretKey := os.Getenv("SUPABASE_SECRET_KEY")
//	endpoint := os.Getenv("SUPABASE_ENDPOINT")
//	region := os.Getenv("SUPABASE_REGION")
//
//	// Verify credentials and endpoint
//	if accessKey == "" || secretKey == "" || endpoint == "" || region == "" {
//		return fmt.Errorf("missing Supabase configuration in environment variables")
//	}
//
//	// Create AWS config with static credentials
//	cfg, err := config.LoadDefaultConfig(ctx,
//		config.WithCredentialsProvider(
//			credentials.NewStaticCredentialsProvider("c2938582e600cf237b7136dc7370a6b9", "d3a92c46abdcc19061b66442bde2ded8f293029e4fefa73099a4ca1701053837", ""),
//		),
//		//config.WithRegion("ap-southeast-1"), // Use the correct region from .env
//		config.WithClientLogMode(aws.LogSigning),
//		//config.WithEndpointCredentialOptions()
//	)
//	if err != nil {
//		return fmt.Errorf("failed to load AWS configuration: %w", err)
//	}
//
//	// Initialize the S3 client with a custom EndpointResolverV2
//	s3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
//		o.BaseEndpoint = aws.String("https://snjwvkwtuboigdykuzcf.supabase.co/storage/v1/s3")
//		//o.EndpointResolverV2 = &resolverV2{}
//		o.Region = region
//		o.UsePathStyle = true
//	})
//
//	s3Client.ListBuckets(context.Background(), nil)
//
//	// Initialize the uploader
//	uploader = manager.NewUploader(s3Client)
//
//	return nil
//}
//
//// UploadToS3 uploads one or multiple images to S3 concurrently, using context for cancellation
//func UploadToS3(ctx context.Context, images []Image) ([]UploadResult, error) {
//	once.Do(func() {
//		initErr = initUploader(ctx)
//	})
//	if initErr != nil {
//		log.Print(initErr)
//		return nil, common.ErrCloudConnectionFailed(initErr)
//	}
//
//	bucketName := os.Getenv("SUPABASE_BUCKET")
//	supabaseURL := os.Getenv("SUPABASE_URL")
//
//	// Channel to limit concurrency
//	maxConcurrentUploads := 10 // Adjust this number as needed
//	semaphore := make(chan struct{}, maxConcurrentUploads)
//
//	// Channel to collect results
//	resultsChan := make(chan UploadResult, len(images))
//
//	var wg sync.WaitGroup
//
//	for _, img := range images {
//		select {
//		case <-ctx.Done():
//			// Context was cancelled before starting the upload
//			resultsChan <- UploadResult{
//				FileName: img.FileName,
//				FileURL:  "",
//				Error:    ctx.Err(),
//			}
//			continue
//		case semaphore <- struct{}{}:
//			// Acquired a slot, proceed with the upload
//		}
//
//		wg.Add(1)
//		go func(img Image) {
//			defer wg.Done()
//			defer func() { <-semaphore }() // Release the slot
//
//			var uploadErr error
//			var fileURL string
//
//			// Check for file size limit (example: 100MB)
//			if len(img.FileBuffer) > 100*1024*1024 {
//				uploadErr = common.ErrCannotUploadFile("image", fmt.Errorf("file size: %d bytes", len(img.FileBuffer)))
//			} else {
//				// Upload smaller files (<5MB) directly, use uploader for larger files
//				input := &s3.PutObjectInput{
//					Bucket: aws.String(bucketName),
//					Key:    aws.String(img.FileName),
//					Body:   bytes.NewReader(img.FileBuffer),
//					//ContentType: aws.String("application/octet-stream"),
//				}
//				if len(img.FileBuffer) < (5 * 1024 * 1024) {
//					_, uploadErr = s3Client.PutObject(ctx, input)
//				} else {
//					_, uploadErr = uploader.Upload(ctx, input)
//				}
//
//				// If upload succeeds, construct the file URL
//				if uploadErr == nil {
//					fileURL = fmt.Sprintf("%s/storage/v1/object/public/%s/%s", supabaseURL, bucketName, img.FileName)
//				} else {
//					log.Print(uploadErr)
//					uploadErr = common.ErrCannotUploadFile(img.FileName, uploadErr)
//					// Ensure the error is logged with more detail for debugging
//					resultsChan <- UploadResult{
//						FileName: img.FileName,
//						FileURL:  "",
//						Error:    uploadErr,
//					}
//					return // Stop further processing in case of error
//				}
//
//			}
//
//			log.Printf("Attempting to upload file: %s to bucket: %s", img.FileName, bucketName)
//			log.Printf("Using endpoint: %s", *s3Client.Options().BaseEndpoint)
//			resultsChan <- UploadResult{
//				FileName: img.FileName,
//				FileURL:  fileURL,
//				Error:    uploadErr,
//			}
//		}(img)
//	}
//
//	// Wait for all uploads to complete
//	wg.Wait()
//	close(resultsChan)
//
//	// Collect results
//	results := make([]UploadResult, 0, len(images))
//	var anyError bool
//	for result := range resultsChan {
//		if result.Error != nil {
//			anyError = true
//		}
//		results = append(results, result)
//	}
//
//	if anyError {
//		return results, common.ErrCannotUploadFile("images", fmt.Errorf("one or more uploads failed"))
//	}
//
//	return results, nil
//}
//
//// UploadSingleImageToS3 uploads a single image to S3
//func UploadSingleImageToS3(ctx context.Context, fileBuffer []byte, fileName string) (string, error) {
//	images := []Image{
//		{
//			FileName:   fileName,
//			FileBuffer: fileBuffer,
//		},
//	}
//	results, err := UploadToS3(ctx, images)
//	if err != nil {
//		return "", err
//	}
//	result := results[0]
//	return result.FileURL, result.Error
//}
//
//// DeleteSingleImageFromS3 deletes a single image from S3 without using waiters
//func DeleteSingleImageFromS3(ctx context.Context, fileName string) error {
//	once.Do(func() {
//		initErr = initUploader(ctx)
//	})
//	if initErr != nil {
//		return common.ErrCloudConnectionFailed(initErr)
//	}
//
//	bucketName := os.Getenv("SUPABASE_BUCKET")
//
//	// Use the S3 client to delete the object
//	_, err := s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
//		Bucket: aws.String(bucketName),
//		Key:    aws.String(fileName),
//	})
//	if err != nil {
//		var ae smithy.APIError
//		if errors.As(err, &ae) {
//			log.Printf("AWS API Error while deleting %s: %s", fileName, ae.ErrorMessage())
//		} else {
//			log.Printf("Failed to delete %s: %v", fileName, err)
//		}
//		return common.ErrCannotDeleteFile(imagemodel.EntityName, err)
//	}
//
//	log.Printf("Successfully deleted %s from bucket %s", fileName, bucketName)
//	return nil
//}
