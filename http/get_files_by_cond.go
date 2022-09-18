package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/d0kur0/webm-api/worker"

	"github.com/d0kur0/webm-grabber/types"
)

type getFilesCondition map[string][]string

func getFilesByCondition(w http.ResponseWriter, r *http.Request) {
	var condition getFilesCondition
	err := json.NewDecoder(r.Body).Decode(&condition)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var filteredFiles types.Output
	for _, outItem := range worker.GrabbingOutPut {
		boards, isDesiredVendor := condition[outItem.VendorName]
		if !isDesiredVendor {
			continue
		}

		isDesiredBoard := false
		for _, board := range boards {
			if board == outItem.BoardName {
				isDesiredBoard = true
				break
			}
		}

		if !isDesiredBoard {
			continue
		}

		filteredFiles = append(filteredFiles, outItem)
	}

	outputAsBytes, err := json.Marshal(filteredFiles)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	_, err = io.WriteString(w, string(outputAsBytes))
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
