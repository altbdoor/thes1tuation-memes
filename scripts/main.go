package main

import (
	"flag"
	"log"
	"path/filepath"
	"runtime"
	"sync"
)

func main() {
	discordFlag := flag.Bool("discord", false, "Parse Discord data")
	imgurFlag := flag.Bool("imgur", false, "Parse imgur data")
	flag.Parse()

	// ========================================
	// get base dir
	// ========================================

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalln("unable to retrieve script path")
	}

	baseDir := filepath.Join(filepath.Dir(filename), "../")
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
