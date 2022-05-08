package util

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// reads an .env file and sets an environment variable for each correctly formatted line.
//
// Key and value must be separated by =, it only splits by the first =,
//
// so you can have = in the value but not in the key name
//
// Both key and value will be trimmed of leading and trailing white spaces,
//
// Lines starting with '#' are ignored, so are the incorrectly formatted ones./n
//
// Param: path - the url to the file where the env variables are declared.
func init() {
	var file *os.File
	var err error
	if file, err = os.Open(".env"); err != nil {
		log.Fatalf("Error opening environment variable file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if len(scanner.Text()) > 0 && !strings.HasPrefix(scanner.Text(), "#") {
			pairKeyValue := strings.SplitN(scanner.Text(), "=", 2)
			if len(pairKeyValue) == 2 {
				if err = os.Setenv(strings.Trim(pairKeyValue[0], " "),
					strings.Trim(pairKeyValue[1], " ")); err != nil {
					log.Fatalf("Error reading environment variable file: %s", err)
				}
			}
		}
	}

	if err = file.Close(); err != nil {
		log.Fatalf("Error closing environment variable file: %s", err)
	}

}
