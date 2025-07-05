package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func BackupB2(logger *slog.Logger) {
	// get key pair
	keyPair := os.Getenv("B2_API_KEY")
	keyPairSplit := strings.SplitN(keyPair, ":", 2)
	if len(keyPairSplit) != 2 {
		log.Fatal("incorrect key pair in B2_API_KEY, should be <KEY_ID>:<APP_KEY>")
	}

	accessKeyID := keyPairSplit[0]
	secretAccessKey := keyPairSplit[1]

	// create tmp path
	tmpAlbumPath := path.Join(os.TempDir(), "magomet2.zip")
	logger.Info(fmt.Sprintf("creating %s", tmpAlbumPath))

	tmpAlbum, err := os.Create(tmpAlbumPath)
	if err != nil {
		log.Fatalf("unable to create tmp file: %v", err)
	}
	defer tmpAlbum.Close()

	// download album zip
	logger.Info("downloading album zip")
	httpClient := &http.Client{Timeout: 10 * time.Second}
	albumUrl := "https://imgur.com/a/xUok0eh/zip"
	resp, err := httpClient.Get(albumUrl)
	if err != nil {
		log.Fatalf("unable to create http client: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("unable to access album url: %v", resp)
	}

	// save album
	logger.Info("writing album zip into tmp file")
	_, err = io.Copy(tmpAlbum, resp.Body)
	if err != nil {
		log.Fatalf("unable to write album into tmp file: %v", err)
	}

	// creating s3
	logger.Info("creating s3 client to b2")
	s3Client, err := minio.New("s3.us-west-000.backblazeb2.com", &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		log.Fatalf("unable to init s3 client to b2: %v", err)
	}

	ctx := context.Background()
	putOptions := minio.PutObjectOptions{
		ContentType: "application/zip",
	}

	// upload b2
	logger.Info("uploading to b2")
	uploadInfo, err := s3Client.FPutObject(ctx, "thes1tuation-memes", "magomet2.zip", tmpAlbumPath, putOptions)
	if err != nil {
		log.Fatalf("unable to upload: %v", err)
	}

	sizeMB := float64(uploadInfo.Size) / 1_000_000
	logger.Info(fmt.Sprintf("done, uploaded %.1fMB", sizeMB))
}
