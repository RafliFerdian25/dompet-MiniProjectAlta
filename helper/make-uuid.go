package helper

import (
	"strings"

	"github.com/google/uuid"
)

func GenerateUUID() string {
	id := uuid.NewString()
	uuid := strings.ReplaceAll(id, "-", "")
	return uuid
}