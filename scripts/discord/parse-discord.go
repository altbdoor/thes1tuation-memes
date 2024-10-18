package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

type RawDiscordData struct {
	Filename  string   `json:"filename"`
	Timestamp string   `json:"timestamp"`
	Tags      []string `json:"tags"`
}

type ParsedDiscord struct {
	Index       int    `json:"index"`
	ID          string `json:"id"`
	Datetime    int64  `json:"datetime"`
	Type        string `json:"type"`
	TimeDisplay string `json:"timeDisplay"`
	GroupBy     string `json:"groupBy"`
}

type GroupParsedDiscord struct {
	Name  string          `json:"name"`
	Items []ParsedDiscord `json:"items"`
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
	// read json data
	// ========================================

	rawDataPath := filepath.Join(baseDir, "./_data/discord.json")
	rawBytes, _ := os.ReadFile(rawDataPath)

	rawData := []RawDiscordData{}
	json.Unmarshal(rawBytes, &rawData)

	// ========================================
	// parse json data
	// ========================================

	parsedItems := make([]ParsedDiscord, len(rawData))

	for idx, item := range rawData {
		itemDate, _ := time.Parse("01/02/2006 3:04 PM", item.Timestamp)

		itemType := "audio"
		if strings.HasSuffix(item.Filename, ".mp4") {
			itemType = "video"
		}

		parsedItems[idx] = ParsedDiscord{
			Index:       -1,
			ID:          item.Filename,
			Datetime:    itemDate.Unix(),
			Type:        itemType,
			TimeDisplay: itemDate.Format("02 Jan, 03:04 PM"),
			GroupBy:     itemDate.Format("January 2006"),
		}
	}

	sort.Slice(parsedItems, func(i, j int) bool {
		return parsedItems[i].Datetime > parsedItems[j].Datetime
	})
	fmt.Println("(i) finish parsing discord data")

	// ========================================
	// parse json data
	// ========================================

	itemsCount := len(parsedItems)
	orderedGroupKey := []string{}
	groupedData := make(map[string][]ParsedDiscord)

	for idx, item := range parsedItems {
		item.Index = itemsCount - idx

		if _, ok := groupedData[item.GroupBy]; ok {
			groupedData[item.GroupBy] = append(groupedData[item.GroupBy], item)
		} else {
			orderedGroupKey = append(orderedGroupKey, item.GroupBy)
			groupedData[item.GroupBy] = []ParsedDiscord{item}
		}
	}

	orderedGroupData := make([]GroupParsedDiscord, len(orderedGroupKey))
	for idx, key := range orderedGroupKey {
		orderedGroupData[idx] = GroupParsedDiscord{
			Name:  key,
			Items: groupedData[key],
		}
	}

	fmt.Printf("(i) grouped %d media into %d months\n", itemsCount, len(orderedGroupKey))

	// ========================================
	// write the data to files
	// ========================================

	jekyllDataPath := filepath.Join(baseDir, "./_data/discord-parsed.json")
	jsonData, _ := json.MarshalIndent(orderedGroupData, "", "  ")
	os.WriteFile(jekyllDataPath, jsonData, 0644)

	fmt.Println("(i) finish writing json data to files")

}
