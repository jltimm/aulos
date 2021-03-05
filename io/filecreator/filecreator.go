package filecreator

import (
	"encoding/json"
	"os"

	"../../common"
	"../postgres"
)

// CreateFile creates a JSON file from the key map, distance matrix, and path
func CreateFile(keyMap map[string]int, distance [][]int, path [][]int) {
	output := &common.Output{
		AllArtists:          postgres.GetAllArtists(),
		FloydWarshallMatrix: distance,
		ShortestPathsMatrix: path,
		IndexIDMap:          createIndexIDList(keyMap),
	}
	jsonObject, _ := json.Marshal(output)
	f, err := os.Create("data.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(string(jsonObject))
}

func createIndexIDList(keyMap map[string]int) []common.IndexID {
	var indexIDList []common.IndexID
	for key, value := range keyMap {
		indexIDList = append(indexIDList, common.IndexID{
			Index: value,
			ID:    key,
		})
	}
	return indexIDList
}
