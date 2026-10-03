package observer

import "log"

func logResource(action, kind, description string) {
	log.Printf("%s %s: %s", action, kind, description)
}