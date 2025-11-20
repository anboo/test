package rpc

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

var (
	validate = validator.New() //nolint:predeclared,gochecknoglobals
)

type BaseHTTPError struct {
	Message string `json:"message"`
}

func NewBaseHTTPError(message string) *BaseHTTPError {
	return &BaseHTTPError{Message: message}
}

type ValidationErrorResponse struct {
	BaseHTTPError
	Fields map[string]string `json:"fields"`
}

func ShouldBindJSON(r *http.Request, w http.ResponseWriter, obj any) bool {
	if err := json.NewDecoder(r.Body).Decode(obj); err != nil {
		writeValidationError(w, map[string]string{
			"body": "invalid_json",
		})
		return false
	}

	if err := validate.Struct(obj); err != nil {
		var verrs validator.ValidationErrors
		if errors.As(err, &verrs) {
			fields := map[string]string{}
			for _, fe := range verrs {
				fields[fe.Field()] = fe.Tag()
			}
			writeValidationError(w, fields)
			return false
		}

		writeValidationError(w, map[string]string{
			"body": err.Error(),
		})
		return false
	}

	return true
}

func writeValidationError(w http.ResponseWriter, fields map[string]string) {
	resp := ValidationErrorResponse{
		BaseHTTPError: BaseHTTPError{
			Message: "validation_failed",
		},
		Fields: fields,
	}

	WriteJSON(w, http.StatusUnprocessableEntity, resp)
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v) //nolint:errchkjson
}

func WriteBadRequest(w http.ResponseWriter, msg string) {
	WriteJSON(w, http.StatusBadRequest, NewBaseHTTPError(msg))
}

func WriteUnauthorized(w http.ResponseWriter) {
	WriteJSON(w, http.StatusUnauthorized, NewBaseHTTPError("unauthorized"))
}

func WriteNotFound(w http.ResponseWriter, msg string) {
	WriteJSON(w, http.StatusNotFound, NewBaseHTTPError(msg))
}

func WriteForbidden(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, NewBaseHTTPError("access_denied"))
}

func WriteUnexpectedError(w http.ResponseWriter, err error) {
	slog.Info("unhandled error:", "err", err)
	WriteJSON(w, http.StatusInternalServerError, NewBaseHTTPError("internal error"))
}
