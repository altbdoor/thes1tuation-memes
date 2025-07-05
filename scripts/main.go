package main

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

func main() {
	discordFlag := flag.Bool("discord", false, "Parse Discord data")
	imgurFlag := flag.Bool("imgur", false, "Parse imgur data")
	uploadFlag := flag.String("upload", "", "Upload media to Cloudinary")
	backupFlag := flag.Bool("backup", false, "Backup media to Backblaze")
	flag.Parse()

	// ========================================
	// logger
	// ========================================

	baseLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	var logger *slog.Logger

	// ========================================
	// get base dir
	// ========================================

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalln("unable to retrieve script path")
	}

	baseDir := filepath.Join(filepath.Dir(filename), "../")

	// ========================================
	// quick operations
	// ========================================

	if *uploadFlag != "" {
		UploadCloud(*uploadFlag)
		return
	}

	if *backupFlag {
		logger = baseLogger.With("fn", "BackupB2")
		BackupB2(logger)
		return
	}

	// ========================================
	// sync-able operations
	// ========================================
	var wg sync.WaitGroup

	if *discordFlag {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ParseDiscord(baseDir)
		}()
	}

	if *imgurFlag {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ParseImgur(baseDir)
		}()
	}

	wg.Wait()
}
