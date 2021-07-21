package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

func generateFields(r *http.Request) logrus.Fields {
	fields := logrus.Fields{
		"ts":          time.Now().UTC().Format(time.RFC3339),
		"http_proto":  r.Proto,
		"http_method": r.Method,
		"remote_addr": r.RemoteAddr,
		"user_agent":  r.UserAgent(),
		"uri":         r.RequestURI,
	}

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		fields["req_id"] = reqID
	}

	return fields
}

func j(w http.ResponseWriter, code int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Status", strconv.Itoa(code))
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		return fmt.Errorf("cannot write response: %w", err)
	}

	return nil
}

func jError(w http.ResponseWriter, err error) error {
	code := http.StatusInternalServerError
	localizedError := "Внутренняя ошибка!"

	switch err {

	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Status", strconv.Itoa(code))
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"error":           err.Error(),
		"localized_error": localizedError,
	}); err != nil {
		return fmt.Errorf("cannot write response: %w", err)
	}

	return nil
}
