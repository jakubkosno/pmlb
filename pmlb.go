package pmlb

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func ReadAllSummaryStats() ([]byte, error) {
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

	return content, nil
}