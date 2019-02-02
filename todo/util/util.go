package util

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// ReadBody -
func ReadBody(w http.ResponseWriter, r *http.Request, obj interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &obj)
	if err != nil {
		return err
	}
	return nil
}

// ParseAndWrite -
func ParseAndWrite(w http.ResponseWriter, obj interface{}) error {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	_, err = w.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}
