package services_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/pmenta/goencoder/src/application/repositories"
	"github.com/pmenta/goencoder/src/application/services"
	"github.com/pmenta/goencoder/src/domain"
	"github.com/pmenta/goencoder/src/framework/database"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatalf("could not load env file")
	}
}

func prepare() (*domain.Video, repositories.VideoRepositoryDb) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "test.mp4"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: db}
	repo.Insert(video)

	return video, repo
}

func TestDownload(t *testing.T) {
	video, repo := prepare()
	service := services.NewVideoService()

	service.Video = video
	service.VideoRepository = repo

	err := service.Download(os.Getenv("BUCKET_NAME"))
	require.Nil(t, err)

	err = service.Fragment()
	require.Nil(t, err)

	err = service.Encode()
	require.Nil(t, err)

	err = service.Finish()
	require.Nil(t, err)
}
