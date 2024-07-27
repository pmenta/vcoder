package services

import (
	"context"
	"io"
	"log"
	"os"
	"os/exec"

	"cloud.google.com/go/storage"
	"github.com/pmenta/goencoder/src/application/repositories"
	"github.com/pmenta/goencoder/src/domain"
)

type VideoService struct {
	Video           *domain.Video
	VideoRepository repositories.VideoRepository
}

func NewVideoService() VideoService {
	return VideoService{}
}

func (v *VideoService) Download(bucket_name string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	bkt := client.Bucket(bucket_name)
	obj := bkt.Object(v.Video.FilePath)
	r, err := obj.NewReader(ctx)

	if err != nil {
		return err
	}
	defer r.Close()

	body, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	f, err := os.Create(os.Getenv("LOCAL_STORAGE_PATH") + "/" + v.Video.ID + "mp4")
	if err != nil {
		return err
	}

	_, err = f.Write(body)
	if err != nil {
		return err
	}
	defer f.Close()

	log.Printf("video %v has been store", v.Video.ID)

	return nil
}

func (v *VideoService) Fragment() error {
	err := os.Mkdir(os.Getenv("LOCAL_STORAGE_PATH")+"/"+v.Video.ID, os.ModePerm)
	if err != nil {
		return err
	}

	source := os.Getenv("LOCAL_STORAGE_PATH" + "/" + v.Video.ID + ".mp4")
	target := os.Getenv("LOCAL_STORAGE_PATH" + "/" + v.Video.ID + ".frag")

	cmd := exec.Command("mp4fragment", source, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func printOutput(out []byte) {
	if len(out) > 0 {
		log.Printf("===========> Output: %s\n", string(out))
	}
}
