package main

import "os"

func initPiner() error {
	file, err := os.Create("pine.yaml")
	if err != nil {
		return err
	}
}
