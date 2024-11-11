package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
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

func ParseDiscord(baseDir string) {
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
	log.Println("(i) discord: finish parsing discord data")

	// ========================================
	// grouping the images by month and year
	// ========================================

	itemsCount := len(parsedItems)
	orderedGroupKey := []string{}
	groupedData := make(map[string][]ParsedDiscord)

	for idx, item := range parsedItems {
		item.Index = itemsCount - idx
		group := item.GroupBy

		if _, ok := groupedData[group]; !ok {
			orderedGroupKey = append(orderedGroupKey, group)
		}

		groupedData[group] = append(groupedData[group], item)
	}

	orderedGroupData := ConvertToGroupedData(orderedGroupKey, groupedData)
	log.Printf("(i) discord: grouped %d media into %d months\n", itemsCount, len(orderedGroupKey))

	// ========================================
	// write the data to files
	// ========================================

	jekyllDataPath := filepath.Join(baseDir, "./_data/discord-parsed.json")
	jsonData, _ := json.MarshalIndent(orderedGroupData, "", "  ")
	os.WriteFile(jekyllDataPath, jsonData, 0644)

	log.Println("(i) discord: finish writing json data to files")

}
