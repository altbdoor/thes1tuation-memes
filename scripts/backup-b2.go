package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func BackupB2(logger *slog.Logger, path string) {
	// get key pair
	keyPair := os.Getenv("B2_API_KEY")
	keyPairSplit := strings.SplitN(keyPair, ":", 2)
	if len(keyPairSplit) != 2 {
		log.Fatal("incorrect key pair in B2_API_KEY, should be <KEY_ID>:<APP_KEY>")
	}

	accessKeyID := keyPairSplit[0]
	secretAccessKey := keyPairSplit[1]

	// check zip file
	tmpAlbumPath := path
	if strings.HasPrefix(tmpAlbumPath, "~") {
		homedir, _ := os.UserHomeDir()
		tmpAlbumPath = homedir + tmpAlbumPath[1:]
	}
	tmpAlbumPath = filepath.Clean(tmpAlbumPath)
	logger.Info(fmt.Sprintf("using zip at %s", tmpAlbumPath))

	_, err := os.Stat(tmpAlbumPath)
	if err != nil {
		log.Fatalf(`unable to open file at "%s": %v`, path, err)
	}

	// creating s3 client
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
