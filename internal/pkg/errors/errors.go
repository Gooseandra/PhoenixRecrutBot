package errors

import (
	"log"
	"strconv"
)

func Warn(id int64, message string) {
	log.Println("WARNING"+strconv.Itoa(int(id))+": %s", message)
}
