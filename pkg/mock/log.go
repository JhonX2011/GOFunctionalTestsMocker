package mock

import (
	"fmt"
	"time"
)

func LogInfo(message string, args ...any) {
	log("INFO", message, args)
}

func LogError(message string, args ...any) {
	log("ERROR", message, args)
}

func log(verb string, message string, args []any) {
	template := "[%s] %s - %s"
	actualDate := time.Now().String()
	if len(args) > 0 {
		message = fmt.Sprintf(message, args)
	}
	fmt.Println(fmt.Sprintf(template, verb, actualDate, message))
}
