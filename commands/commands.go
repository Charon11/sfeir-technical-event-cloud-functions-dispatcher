package commands

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type Command interface{}

func getEntityId(r *http.Request) string {
	re := regexp.MustCompile(`\w{8}-\w{4}-\w{4}-\w{4}-\w{12}`)
	return re.FindString(r.URL.Path)
}

func getBearerToken(r *http.Request) string {
	return strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
}

func getCommand(r *http.Request, command Command) Command {
	if err := json.NewDecoder(r.Body).Decode(&command); err != nil {
		log.Fatalf("Unmarshal command: %v", err)
	}
	return command
}
