package domain_test

import (
	"testing"
	"time"

	"github.com/pmenta/goencoder/src/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestNewJob(t *testing.T) {
	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	job, err := domain.NewJob(video.FilePath, "Converted", video)

	require.Nil(t, err)
	require.NotNil(t, job)
}
