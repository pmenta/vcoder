package domain_test

import (
	"testing"
	"time"

	"github.com/pmenta/goencoder/src/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestValidateIfVideoIsEmpty(t *testing.T) {
	video := domain.NewVideo()
	err := video.Validate()

	require.Error(t, err)
}

func TestVideoIdIsNotAnUuid(t *testing.T) {
	video := domain.NewVideo()
	video.ID = "foo"
	video.ResourceID = "bar"
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	err := video.Validate()
	require.Error(t, err)
}

func TestVideoValidation(t *testing.T) {
	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.ResourceID = "bar"
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	err := video.Validate()
	require.Nil(t, err)
}
