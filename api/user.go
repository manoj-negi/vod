package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"mime"
	"net/http"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	db "github.com/tiktok/db/sqlc"
)

type Video struct {
	Name     string `json:"name"`
	ID       int    `json:"id"`
	IsActive bool   `json:"is_active"`
}

type UploadResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func errorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errResponse := ErrorResponse{
		Status:  statusCode,
		Message: message,
	}
	json.NewEncoder(w).Encode(errResponse)
}

func createS3Client(key, secret string) (*session.Session, *s3.S3, error) {
	awsConfig := aws.Config{
		Region:      aws.String("us-west-1"), // Replace with your desired AWS region.
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
	}

	// Create a new AWS session.
	sess, err := session.NewSession(&awsConfig)
	if err != nil {
		return nil, nil, err
	}

	svc := s3.New(sess)
	return sess, svc, nil
}

func (server *Server) listAllVideos(w http.ResponseWriter, r *http.Request) {
	//videos, err := fetchAllVideos()

	ctx := r.Context()
	videos, err := server.store.ListVideos(ctx, server.config.BUCKET_URL)

	if err != nil {
		errorResponse(w, http.StatusNotFound, "No data")
		return
	}
	if len(videos) == 0 {
		// No videos found, construct a custom response.
		emptyResponse := struct {
			Message string `json:"message"`
		}{
			Message: "No videos found.",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(emptyResponse)
		return
	}

	// Return the list of video names as a JSON response.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)
}

func (server *Server) uploadVideoToS3(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		errorResponse(w, http.StatusMethodNotAllowed, "Only POST requests are allowed")
		return
	}
	ctx := r.Context()
	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		errorResponse(w, http.StatusBadRequest, "Failed to get file from request")
		return
	}

	defer file.Close()

	filePath := header.Filename             // Replace this with the actual path to your video file.
	bucketName := server.config.BUCKET_NAME // Replace this with the name of your S3 bucket.
	uniqueFilename := generateUniqueFilename(filePath)
	objectKey := "videos/" + uniqueFilename

	ext := filepath.Ext(filePath)

	sess, svc, err := createS3Client(server.config.AWS_KEY, server.config.AWS_SECRET)
	if err != nil {
		errorResponse(w, http.StatusMethodNotAllowed, "Failed to stream video content")
		return
	}
	defer sess.Config.Credentials.Expire()

	// Prepare the input parameters for the S3 upload.
	input := &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(objectKey),
		Body:          file,
		ContentLength: aws.Int64(header.Size),
		ContentType:   aws.String(mime.TypeByExtension(ext)), // Replace with the appropriate content type for your video file.
		ACL:           aws.String("public-read"),             // Change this as needed, depending on your security requirements.
	}

	// Upload the video to S3.
	_, err = svc.PutObject(input)

	if err != nil {
		slog.Info("=====", err)
		errorResponse(w, http.StatusMethodNotAllowed, "Failed to upload video content")
		return
	}

	arg := db.CreateVideoParams{
		Name:       objectKey,
		UploadedBy: 1,
	}

	video, err := server.store.CreateVideo(ctx, arg)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to insert video name into PostgreSQL")
		return
	}

	response := UploadResponse{
		Status:  "success",
		Message: "Video uploaded successfully to S3.",
		Data:    video,
	}

	// Respond with a JSON response indicating success.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func generateUniqueFilename(originalFilename string) string {
	timestamp := time.Now().Format("20060102150405") // Format: YYYYMMDDHHMMSS
	randomUUID := uuid.New().String()
	randomComponent := rand.Intn(1000000) // Generates a random number between 0 and 999999

	return fmt.Sprintf("%s_%s_%d_%s", timestamp, randomUUID, randomComponent, originalFilename)
}
