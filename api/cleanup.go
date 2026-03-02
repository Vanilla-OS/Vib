package api

import (
	"fmt"
	"strings"
)

func GetCleanupSuffix(cleanup []string) string {
	if len(cleanup) > 0 {
	    return fmt.Sprintf(" && rm -rf %s", strings.Join(cleanup, " "))
	}
	return ""
}
