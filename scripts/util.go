package main

type GroupParsedData[G any] struct {
	Name  string `json:"name"`
	Items []G    `json:"items"`
}

func ConvertToGroupedData[T any](keys []string, groupedData map[string][]T) []GroupParsedData[T] {
	orderedGroupData := make([]GroupParsedData[T], len(keys))
	for idx, key := range keys {
		orderedGroupData[idx] = GroupParsedData[T]{
			Name:  key,
			Items: groupedData[key],
		}
	}
	return orderedGroupData
}
