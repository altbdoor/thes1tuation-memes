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

func UploadCloud(pathToFile string) {
	cloudKeys := os.Getenv("CLOUDINARY_USER")
	if cloudKeys == "" {
		log.Fatalln("(!) cloudinary: invalid key secret pair")
	}

	cloudKeySecretPair := strings.Split(cloudKeys, ":")
	cloudClient, _ := cloudinary.NewFromParams("dsakciquw", cloudKeySecretPair[0], cloudKeySecretPair[1])

	cleanPath, _ := filepath.Abs(pathToFile)
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
