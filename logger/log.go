package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

var err error

var infoFile *os.File
var warnFile *os.File
var errorFile *os.File

var InfoLogger *log.Logger
var WarnLogger *log.Logger
var ErrorLogger *log.Logger
var ConsoleLogger *log.Logger

var logFiles = map[string]*os.File{"infolog": infoFile, "warnlog": warnFile, "errorlog": errorFile}

func init() {

	if _, err = os.Stat("logs"); os.IsNotExist(err) {

		err = os.Mkdir("logs", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	timeFormat := time.Now().Format("02-01-2006")

	for filename, file := range logFiles {
		createLogFile(file, filename, timeFormat)
	}

	// Initialize logs
	ConsoleLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger = log.New(infoFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLogger = log.New(warnFile, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(errorFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func createLogFile(file *os.File, fileName string, timeFormat string) {
	file, err = os.OpenFile(fmt.Sprintf("logs/%s_%s.log", fileName, timeFormat),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
}

// Close closes the log files
func Close() {

	for _, file := range logFiles {
		err = file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
}
