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

	target := os.Getenv("LOCAL_STORAGE_PATH") + "/"
	f, err := os.Create(target + "/" + v.Video.ID + ".mp4")
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
	source := os.Getenv("LOCAL_STORAGE_PATH") + "/" + v.Video.ID + ".mp4"
	target := os.Getenv("LOCAL_STORAGE_PATH") + "/" + v.Video.ID + ".frag"

	cmd := exec.Command("mp4fragment", source, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func (v *VideoService) Encode() error {
	target := os.Getenv("LOCAL_STORAGE_PATH") + "/" + v.Video.ID
	err := os.Mkdir(target, os.ModePerm)
	if err != nil {
		return err
	}

	cmdArgs := []string{}
	cmdArgs = append(cmdArgs, os.Getenv("LOCAL_STORAGE_PATH")+"/"+v.Video.ID+".frag")
	cmdArgs = append(cmdArgs, "--use-segment-timeline")
	cmdArgs = append(cmdArgs, "-o")
	cmdArgs = append(cmdArgs, os.Getenv("LOCAL_STORAGE_PATH")+"/"+v.Video.ID)
	cmdArgs = append(cmdArgs, "-f")
	cmdArgs = append(cmdArgs, "--exec-dir")
	cmdArgs = append(cmdArgs, "/opt/homebrew/bin/mp4dash")

	cmd := exec.Command("mp4dash", cmdArgs...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func (v *VideoService) Finish() error {
	err := os.RemoveAll(os.Getenv("LOCAL_STORAGE_PATH") + "/" + v.Video.ID)
	if err != nil {
		log.Fatalf("error removing: %v", v.Video.ID)
		return err
	}

	err = os.Remove(os.Getenv("LOCAL_STORAGE_PATH") + "/" + v.Video.ID + ".mp4")
	if err != nil {
		log.Fatalf("error removing: %v", v.Video.ID)
		return err
	}

	err = os.Remove(os.Getenv("LOCAL_STORAGE_PATH") + "/" + v.Video.ID + ".frag")
	if err != nil {
		log.Fatalf("error removing: %v", v.Video.ID)
		return err
	}

	log.Printf("files have been removed: %v", v.Video.ID)

	return nil
}

func printOutput(out []byte) {
	if len(out) > 0 {
		log.Printf("===========> Output: %s\n", string(out))
	}
}
