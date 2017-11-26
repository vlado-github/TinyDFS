package logging

import (
	"fmt"
	"log"
	"os"
	"path"
	"sync"
)

const pathToLogDir = "../tinydfs/log"

type tinylogger struct {
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger

	TraceLogFile   *os.File
	InfoLogFile    *os.File
	WarningLogFile *os.File
	ErrorLogFile   *os.File

	Verbose bool
}

var instance *tinylogger
var once sync.Once

func getInstance() *tinylogger {
	once.Do(func() {
		traceDest := createDestinationFile("log.trace")
		infoDest := createDestinationFile("log.info")
		warnDest := createDestinationFile("log.warn")
		errorDest := createDestinationFile("log.error")

		instance = &tinylogger{
			Trace: log.New(traceDest,
				"TRACE: ",
				log.Ldate|log.Ltime|log.Lshortfile),
			Info: log.New(infoDest,
				"INFO: ",
				log.Ldate|log.Ltime|log.Lshortfile),
			Warning: log.New(warnDest,
				"WARNING: ",
				log.Ldate|log.Ltime|log.Lshortfile),
			Error: log.New(errorDest,
				"ERROR: ",
				log.Ldate|log.Ltime|log.Lshortfile),
			TraceLogFile:   traceDest,
			InfoLogFile:    infoDest,
			WarningLogFile: warnDest,
			ErrorLogFile:   errorDest,
		}
	})
	return instance
}

func AddTrace(v ...interface{}) {
	tl := getInstance()
	tl.Trace.Println(v)
	postLog("TRACE: ", v)
}

func AddInfo(v ...interface{}) {
	tl := getInstance()
	tl.Info.Println(v)
	postLog("INFO:", v)
}

func AddWarning(v ...interface{}) {
	tl := getInstance()
	tl.Warning.Println(v)
	postLog("WARN: ", v)
}

func AddError(v ...interface{}) {
	tl := getInstance()
	tl.Error.Println(v)
	postLog("ERROR: ", v)
}

func SetVerbose(verbose bool) {
	tl := getInstance()
	tl.Verbose = verbose
}

func Close() {
	tl := getInstance()
	tl.TraceLogFile.Close()
	tl.InfoLogFile.Close()
	tl.WarningLogFile.Close()
	tl.ErrorLogFile.Close()
}

func createDestinationFile(filename string) *os.File {
	err := os.MkdirAll(pathToLogDir, os.ModePerm)
	if err != nil {
		fmt.Println("Persistance: Can not create a directory.", err.Error())
		log.Fatal(err)
	}
	pathToFile := path.Clean(path.Join(pathToLogDir, filename))
	logFile, err := os.OpenFile(pathToFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("error creating log file:", err)
	}
	return logFile
}

func postLog(v ...interface{}) {
	tl := getInstance()
	if tl.Verbose {
		fmt.Println(v)
	}
}
