package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

var albumId = [2]string{
	"yzKq60n",
	"xUok0eh",
}

type ImgurResponseMedia struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	MimeType  string `json:"mime_type"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Size      int    `json:"size"`
	URL       string `json:"url"`
}

type ImgurResponse struct {
	Media []ImgurResponseMedia `json:"media"`
}

type ParsedImage struct {
	Index       int      `json:"index"`
	ID          string   `json:"id"`
	Datetime    int64    `json:"datetime"`
	Type        string   `json:"type"`
	Width       int      `json:"width"`
	Height      int      `json:"height"`
	Size        int      `json:"size"`
	Link        string   `json:"link"`
	Thumbnail   string   `json:"thumbnail"`
	TimeDisplay string   `json:"timeDisplay"`
	GroupBy     string   `json:"groupBy"`
	Year        string   `json:"year"`
	Tags        []string `json:"tags"`
}

type GroupParsedImage struct {
	Name  string        `json:"name"`
	Items []ParsedImage `json:"items"`
}

func main() {
	// ========================================
	// get base dir
	// ========================================

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("unable to retrieve script path")
		os.Exit(1)
	}

	baseDir := filepath.Join(filepath.Dir(filename), "../../")

	// ========================================
	// call imgur api
	// ========================================

	clientId := os.Getenv("IMGUR_CLIENT_ID")
	if clientId == "" {
		fmt.Println("please provide IMGUR_CLIENT_ID")
		os.Exit(1)
	}

	rawImages := []ImgurResponseMedia{}
	client := &http.Client{}

	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		fmt.Sprintf("Chrome/79.0.3945.%d Safari/537.%d", rand.Intn(9999), rand.Intn(99))

	for _, id := range albumId {
		fmt.Printf("(i) calling somewhat official imgur API for album ID %s\n", id)

		albumUrl := fmt.Sprintf("https://api.imgur.com/post/v1/albums/%s?include=media,tags,account", id)
		req, _ := http.NewRequest("GET", albumUrl, nil)

		req.Header.Set("User-Agent", userAgent)
		req.Header.Set("Authorization", "Client-ID "+clientId)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("error sending request", err)
			os.Exit(1)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			fmt.Println("error reading body", err)
			os.Exit(1)
		}

		var responseData ImgurResponse
		err = json.Unmarshal(body, &responseData)
		if err != nil {
			fmt.Println("error parsing json", err)
			os.Exit(1)
		}

		rawImages = append(rawImages, responseData.Media...)
	}

	imagesCount := len(rawImages)

	// ========================================
	// prepare the existing tags
	// ========================================

	tagsDataPath := filepath.Join(baseDir, "./_data/imgur-tags.yml")
	tagsFile, _ := os.Open(tagsDataPath)
	oldTagsMap := make(map[string][]string)

	scanner := bufio.NewScanner(tagsFile)
	for scanner.Scan() {
		lineText := scanner.Text()
		if lineText == "" {
			continue
		}

		lineData := strings.Split(lineText, ": ")
		tagsArray := strings.Trim(lineData[1], "[] ")
		oldTagsMap[lineData[0]] = strings.Split(tagsArray, ", ")
	}

	tagsFile.Close()
	fmt.Println("(i) loaded all tags")

	// ========================================
	// parsing images into custom format
	// ========================================

	// make slice, so it can be sorted
	parsedImages := make([]ParsedImage, imagesCount)

	for idx, img := range rawImages {
		// convert to KL time
		imgDate, _ := time.Parse(time.RFC3339, img.CreatedAt)
		location, _ := time.LoadLocation("Asia/Kuala_Lumpur")
		imgDate = imgDate.In(location)

		thumbnail := strings.ReplaceAll(img.URL, img.ID, img.ID+"b")
		thumbnail = strings.ReplaceAll(thumbnail, ".gif", ".jpg")
		thumbnail = strings.ReplaceAll(thumbnail, ".jpeg", ".jpg")

		imgTags := []string{}
		if _, ok := oldTagsMap[img.ID]; ok {
			imgTags = oldTagsMap[img.ID]
		} else {
			fmt.Printf("(i) unable to find tags for %s\n", img.ID)
		}

		parsedImage := ParsedImage{
			Index:       -1,
			ID:          img.ID,
			Datetime:    imgDate.Unix(),
			Type:        img.MimeType,
			Width:       img.Width,
			Height:      img.Height,
			Size:        img.Size,
			Link:        strings.ReplaceAll(img.URL, ".jpeg", ".jpg"),
			Thumbnail:   thumbnail,
			TimeDisplay: imgDate.Format("02 Jan, 03:04 PM"),
			GroupBy:     imgDate.Format("January 2006"),
			Year:        imgDate.Format("2006"),
			Tags:        imgTags,
		}

		parsedImages[idx] = parsedImage
	}

	sort.Slice(parsedImages, func(i, j int) bool {
		return parsedImages[i].Datetime > parsedImages[j].Datetime
	})
	fmt.Println("(i) finish parsing images")

	// ========================================
	// update the tags file
	// ========================================

	tagsFile, _ = os.Create(tagsDataPath)

	// iterate in reverse
	for i := len(parsedImages) - 1; i >= 0; i-- {
		currentImg := parsedImages[i]
		linePattern := fmt.Sprintf("%s: [ %s ]\n", currentImg.ID, strings.Join(currentImg.Tags, ", "))

		tagsFile.WriteString(linePattern)
	}

	tagsFile.Close()
	fmt.Println("(i) finish updating the tags")

	// ========================================
	// grouping the images by month and year
	// ========================================

	// use a variable to keep track of the ordered keys
	orderedGroupKey := []string{}
	groupedData := make(map[string][]ParsedImage)
	uniqueYears := make(map[string]int)

	for idx, img := range parsedImages {
		img.Index = imagesCount - idx

		if _, ok := uniqueYears[img.Year]; !ok {
			uniqueYears[img.Year] = 0
		}

		uniqueYears[img.Year] += 1

		if _, ok := groupedData[img.GroupBy]; ok {
			groupedData[img.GroupBy] = append(groupedData[img.GroupBy], img)
		} else {
			orderedGroupKey = append(orderedGroupKey, img.GroupBy)
			groupedData[img.GroupBy] = []ParsedImage{img}
		}
	}

	orderedGroupData := make([]GroupParsedImage, len(orderedGroupKey))
	for idx, key := range orderedGroupKey {
		orderedGroupData[idx] = GroupParsedImage{
			Name:  key,
			Items: groupedData[key],
		}
	}

	fmt.Printf("(i) grouped %d images into %d months\n", imagesCount, len(orderedGroupKey))

	// ========================================
	// write the data to files
	// ========================================

	jekyllDataPath := filepath.Join(baseDir, "./_data/imgur-parsed.json")
	assetDataPath := filepath.Join(baseDir, "./assets/imgur.json")

	jsonData, _ := json.MarshalIndent(orderedGroupData, "", "  ")
	os.WriteFile(jekyllDataPath, jsonData, 0644)
	os.WriteFile(assetDataPath, jsonData, 0644)

	fmt.Println("(i) finish writing json data to files")

	// ========================================
	// writing the jekyll md files
	// ========================================

	collectionsDir := filepath.Join(baseDir, "./_imgur/")
	mdFiles, _ := os.ReadDir(collectionsDir)

	for _, file := range mdFiles {
		if strings.HasSuffix(file.Name(), ".md") && !file.IsDir() {
			os.Remove(filepath.Join(collectionsDir, file.Name()))
		}
	}

	for year, count := range uniqueYears {
		content := []string{
			"---",
			fmt.Sprintf("title: imgur %s", year),
			fmt.Sprintf("data_year: %s", year),
			fmt.Sprintf("data_count: %d", count),
			"---",
		}

		strContent := strings.Join(content, "\n")
		byteContent := []byte(strContent)

		imgurRecordPath := filepath.Join(collectionsDir, fmt.Sprintf("%s.md", year))
		os.WriteFile(imgurRecordPath, byteContent, 0644)
	}

	fmt.Println("(i) finish writing imgur collections")

}
