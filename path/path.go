package path

import (
	"math"

	"../common"
	"../io/postgres"
)

func createMap(artists []common.Artist) map[string]int {
	artistsMap := make(map[string]int)
	for i := range artists {
		artistsMap[artists[i].ID] = i
	}
	return artistsMap
}

func setupMatrices(length int) ([][]int, [][]int) {
	var distanceMatrix [][]int = make([][]int, length)
	var pathMatrix [][]int = make([][]int, length)
	for i := 0; i < length; i++ {
		distanceRow := make([]int, length)
		for j := range distanceRow {
			distanceRow[j] = math.MaxInt16
		}
		distanceRow[i] = 0
		distanceMatrix[i] = distanceRow
		pathMatrix[i] = make([]int, length)
	}
	return distanceMatrix, pathMatrix
}

// FloydWarshall runs the floyd warshall algorithm and finds the
// shortest path between all pairs
func FloydWarshall() (map[string]int, [][]int, [][]int) {
	artists := postgres.GetAllArtists()
	keyMap := createMap(artists)
	distance, path := setupMatrices(len(artists))
	for ind := range artists {
		artist := artists[ind]
		recommended := artist.Recommended
		for j := range recommended {
			index := keyMap[recommended[j]]
			distance[ind][index] = 1
			distance[index][ind] = 1
			path[ind][index] = ind
			path[index][ind] = index
		}
	}
	for k := 0; k < len(artists); k++ {
		for i := 0; i < len(artists); i++ {
			for j := 0; j < len(artists); j++ {
				if distance[i][j] > (distance[i][k] + distance[k][j]) {
					distance[i][j] = distance[i][k] + distance[k][j]
					path[i][j] = path[k][j]
				}
			}
		}
	}
	return keyMap, distance, path
}

func getPath(i int, j int, pathMatrix [][]int, path []int) []int {
	if i == j {
		return append(path, i)
	}
	newPath := append(path, j)
	path = getPath(i, pathMatrix[i][j], pathMatrix, newPath)
	return path
}
