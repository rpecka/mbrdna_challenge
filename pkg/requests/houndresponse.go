package requests

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type HoundResponse struct {
	WrittenResponse string
	Intent          string
}

// Adapted from houndify.ParseWrittenResponse
func ParseHoundResponse(serverResponseJSON string) (*HoundResponse, error) {
	result := make(map[string]interface{})
	err := json.Unmarshal([]byte(serverResponseJSON), &result)
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("failed to decode json")
	}
	if !strings.EqualFold(result["Status"].(string), "OK") {
		return nil, errors.New(result["ErrorMessage"].(string))
	}
	if result["NumToReturn"].(float64) < 1 {
		return nil, errors.New("no results to return")
	}
	firstResult := result["AllResults"].([]interface{})[0].(map[string]interface{})
	writtenResponse := firstResult["WrittenResponseLong"].(string)
	intent := extractIntent(firstResult)
	return &HoundResponse{
		WrittenResponse: writtenResponse,
		Intent:          intent,
	}, nil
}

func extractIntent(result map[string]interface{}) string {
	r, ok := result["Result"].(map[string]interface{})
	if !ok {
		return ""
	}
	intent, ok := r["intent"].(string)
	if !ok {
		return ""
	}
	return intent
}
