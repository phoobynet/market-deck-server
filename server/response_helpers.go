package server

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

func writeErr(w http.ResponseWriter, statusCode int, err error) error {
	logrus.Error(err)
	return writeJSON(
		w, statusCode, map[string]string{
			"error": err.Error(),
		},
	)
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) error {
	j, err := json.Marshal(data)

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(j)

	return err
}
