package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// uploads the file to cloudinary, used for the migration of discord media.
// usage example:
//
//	CLOUDINARY_USER=api-key:api-secret go run *.go -upload 'path/to/file.mp4'
func UploadCloud(pathToFile string) {
	cloudKeys := os.Getenv("CLOUDINARY_USER")
	if cloudKeys == "" {
		log.Fatalln("(!) cloudinary: invalid key secret pair")
	}

	cloudKeySecretPair := strings.Split(cloudKeys, ":")
	cloudClient, err := cloudinary.NewFromParams("dsakciquw", cloudKeySecretPair[0], cloudKeySecretPair[1])
	if err != nil {
		log.Fatalln("(!) cloudinary: unable to init API")
	}

	cleanPath, err := filepath.Abs(pathToFile)
	if err != nil {
		log.Fatalf("(!) cloudinary: unable to standardize %s\n", pathToFile)
	}

	log.Println("(i) cloudinary: uploading", cleanPath)
	mediaFile, err := os.Open(cleanPath)

	if err != nil {
		log.Fatalln("(!) cloudinary: unable to open file")
	}

	var (
		trueVal  = true
		falseVal = false
	)

	uploadResult, err := cloudClient.Upload.Upload(
		context.Background(),
		mediaFile,
		uploader.UploadParams{
			AssetFolder:              "magomet_media",
			UseFilename:              &trueVal,
			UseFilenameAsDisplayName: &trueVal,
			UniqueFilename:           &falseVal,
			Overwrite:                &falseVal,
		},
	)

	if err != nil {
		log.Fatalln("(!) cloudinary:", err)
	}

	log.Println("(i) cloudinary:", uploadResult.URL)
}
