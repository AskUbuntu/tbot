package util

import (
	"encoding/json"
	"os"
)

// Load a JSON file from disk if it exists. The first value returned is true
// if the file existed and false if it did not (default values should be set).
func LoadJSON(name string, v interface{}) (bool, error) {
	r, err := os.Open(name)
	if err != nil {
		if !os.IsNotExist(err) {
			return false, err
		}
		return false, nil
	}
	defer r.Close()
	if err := json.NewDecoder(r).Decode(v); err != nil {
		return true, err
	}
	return true, nil
}

// Save writes the JSON file to disk.
func SaveJSON(name string, v interface{}) error {
	w, err := os.Create(name)
	if err != nil {
		return err
	}
	defer w.Close()
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return err
	}
	return nil
}
