package services_test

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/pmenta/goencoder/src/application/services"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatalf("could not load env file")
	}
}

func TestVideoServiceUpload(t *testing.T) {
	video, repo := prepare()
	service := services.NewVideoService()
	bucketName := os.Getenv("BUCKET_NAME")

	service.Video = video
	service.VideoRepository = repo

	err := service.Download(bucketName)
	require.Nil(t, err)

	err = service.Fragment()
	require.Nil(t, err)

	err = service.Encode()
	require.Nil(t, err)

	videoUpload := services.NewVideoUpload()
	videoUpload.OutputBucket = bucketName
	videoUpload.VideoPath = os.Getenv("LOCAL_STORAGE_PATH") + "/" + video.ID

	doneUpload := make(chan string)
	go videoUpload.ProcessUpload(50, doneUpload)

	result := <-doneUpload
	require.Equal(t, result, "upload completed")

	err = service.Finish()
	require.Nil(t, err)
}
