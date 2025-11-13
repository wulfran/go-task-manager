package env

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func loadFromFile(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		key, value, ok := strings.Cut(scanner.Text(), "=")
		if !ok {
			continue
		}
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("error setting the variable %s with the value of %s, error: %v", key, value, err)
		}
	}

	return scanner.Err()
}

func load() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	return loadFromFile(filepath.Join(pwd, ".env"))
}

func init() {
	if err := load(); err != nil {
		log.Fatalf("unable to set up environment variables! %v", err)
	}
}

func Get(name string) string {
	val, ok := os.LookupEnv(name)
	if !ok {
		log.Fatalf("missing %s variable!", name)
	}

	return val
}
