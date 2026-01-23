package slink

import "github.com/charmbracelet/log"

func SetLogLevel(level string) error {
	var err error
	ll, err := log.ParseLevel(level)
	if err != nil {
		return err
	}
	log.SetLevel(ll)
	return nil
}
