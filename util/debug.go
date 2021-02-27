package util

import "log"

//Debug writting debug message if b is true
func Debug(b bool, format string, args ...interface{}) {
	if b {
		log.Printf(format, args...)
	}
}
