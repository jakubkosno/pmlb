package pmlb

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type DatasetInfo struct {
	Dataset              string  `json:"dataset"`
	NInstances           int     `json:"n_instances"`
	NFeatures            int     `json:"n_features"`
	NBinaryFeatures      int     `json:"n_binary_features"`
	NCategoricalFeatures int     `json:"n_categorical_features"`
	NContinuousFeatures  int     `json:"n_continuous_features"`
	EndpointType         string  `json:"endpoint_type"`
	NClasses             int     `json:"n_classes"`
	Imbalance            float64 `json:"imbalance"`
	Task                 string  `json:"task"`
}

func FetchData(datasetName string) ([][]string, error) {
	url := "https://github.com/EpistasisLab/pmlb/raw/master/datasets/" + datasetName + "/" + datasetName + ".tsv.gz"

	// Send GET request
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while reading file: ", err)
		return nil, err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		fmt.Printf("Error: HTTP status: %d\n", response.StatusCode)
		return nil, err
	}

	// Unpack gzip
	gzipReader, err := gzip.NewReader(response.Body)
	if err != nil {
		fmt.Println("Error unpacking gzip file: ", err)
		return nil, err
	}
	defer gzipReader.Close()

	// Read .tsv file
	var contentBuilder strings.Builder
	_, err = io.Copy(&contentBuilder, gzipReader)
	if err != nil {
		fmt.Println("Error reading file content: ", err)
		return nil, err
	}

	var tsvData [][]string
	lines := strings.Split(contentBuilder.String(), "\n")
	for _, line := range lines {
		fields := strings.Split(line, "\t")
		tsvData = append(tsvData, fields)
	}

	return tsvData, nil
}

func FetchXYData(datasetName string) ([][]string, []string, error) {
	url := "https://github.com/EpistasisLab/pmlb/raw/master/datasets/" + datasetName + "/" + datasetName + ".tsv.gz"

	// Send GET request
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while reading file: ", err)
		return nil, nil, err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		fmt.Printf("Error: HTTP status: %d\n", response.StatusCode)
		return nil, nil, err
	}

	// Unpack gzip
	gzipReader, err := gzip.NewReader(response.Body)
	if err != nil {
		fmt.Println("Error unpacking gzip file: ", err)
		return nil, nil, err
	}
	defer gzipReader.Close()

	// Read .tsv file
	var contentBuilder strings.Builder
	_, err = io.Copy(&contentBuilder, gzipReader)
	if err != nil {
		fmt.Println("Error reading file content: ", err)
		return nil, nil, err
	}

	var tsvXData [][]string
	var tsvYData []string
	lines := strings.Split(contentBuilder.String(), "\n")
	for _, line := range lines {
		fields := strings.Split(line, "\t")
		tsvXData = append(tsvXData, fields[:len(fields) - 1])
		tsvYData = append(tsvYData, fields[len(fields) - 1])
	}

	return tsvXData, tsvYData, nil
}

func FindDatasets(task string) ([]string, error) {
	var desiredDatasets []string
	allDatasets, err := readAllSummaryStats()
	if err != nil {
		return nil, err
	}

	for _, datasetInfo := range allDatasets {
		if datasetInfo.Task == task {
			desiredDatasets = append(desiredDatasets, datasetInfo.Dataset)
		}
	}

	return desiredDatasets, nil
}

func readAllSummaryStats() ([]DatasetInfo, error) {
	url := "https://raw.githubusercontent.com/EpistasisLab/pmlb/master/pmlb/all_summary_stats.tsv"
	// Send GET request
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error: HTTP status: %d", response.StatusCode)
	}

	// Read HTTP response
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return parseAllSummaryStats(string(content)), nil
}

func parseAllSummaryStats(content string) ([]DatasetInfo){
	lines := strings.Split(content, "\n")
	var allDatasets []DatasetInfo
	for i, line := range lines {
		if i == 0 {
			continue // Skip headers
		}

		fields := strings.Split(line, "\t")
		// Only create DatasetInfo if row has all information
		if len(fields) == 10 {
			datasetInfo := DatasetInfo{
				Dataset:              fields[0],
				NInstances:           parseInt(fields[1]),
				NFeatures:            parseInt(fields[2]),
				NBinaryFeatures:      parseInt(fields[3]),
				NCategoricalFeatures: parseInt(fields[4]),
				NContinuousFeatures:  parseInt(fields[5]),
				EndpointType:         fields[6],
				NClasses:             parseInt(fields[7]),
				Imbalance:            parseFloat(fields[8]),
				Task:                 fields[9],
			}

			allDatasets = append(allDatasets, datasetInfo)
		}
	}

	return allDatasets
}

func parseInt(s string) int {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	if err != nil {
		return 0
	}
	return result
}

func parseFloat(s string) float64 {
	var result float64
	_, err := fmt.Sscanf(s, "%f", &result)
	if err != nil {
		return 0.0
	}
	return result
}
