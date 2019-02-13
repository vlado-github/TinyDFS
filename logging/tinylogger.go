package logging

import (
	"fmt"
	"log"
	"os"
	"path"
	"sync"
)

// VerbosityLevelType how many details on a console
type VerbosityLevelType int

const (
	ALL     VerbosityLevelType = iota
	TRACE   VerbosityLevelType = iota
	INFO    VerbosityLevelType = iota
	ERROR   VerbosityLevelType = iota
	WARNING VerbosityLevelType = iota
	NONE    VerbosityLevelType = iota
)

// Configuration represents the variables
// used by logger
type Configuration struct {
	pathToLogDir string
	verbose      VerbosityLevelType
}

// Config contains info on logger configuration
var config = &Configuration{verbose: ALL, pathToLogDir: "."}

type tinylogger struct {
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger

	TraceLogFile   *os.File
	InfoLogFile    *os.File
	WarningLogFile *os.File
	ErrorLogFile   *os.File
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
	postLog(TRACE, "TRACE: ", v)
}

func AddInfo(v ...interface{}) {
	tl := getInstance()
	tl.Info.Println(v)
	postLog(INFO, "INFO:", v)
}

func AddWarning(v ...interface{}) {
	tl := getInstance()
	tl.Warning.Println(v)
	postLog(WARNING, "WARN: ", v)
}

func AddError(v ...interface{}) {
	tl := getInstance()
	tl.Error.Println(v)
	postLog(ERROR, "ERROR: ", v)
}

// SetConfiguration is used to setup config variables
func SetConfiguration(verboseType VerbosityLevelType, pathToDir string) {
	config = &Configuration{verbose: verboseType, pathToLogDir: pathToDir}
}

func Close() {
	tl := getInstance()
	tl.TraceLogFile.Close()
	tl.InfoLogFile.Close()
	tl.WarningLogFile.Close()
	tl.ErrorLogFile.Close()
}

func createDestinationFile(filename string) *os.File {
	err := os.MkdirAll(config.pathToLogDir, os.ModePerm)
	if err != nil {
		fmt.Println("Persistance: Can not create a directory.", err.Error())
		log.Fatal(err)
	}
	pathToFile := path.Clean(path.Join(config.pathToLogDir, filename))
	logFile, err := os.OpenFile(pathToFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("error creating log file:", err)
	}
	return logFile
}

func postLog(verbosityLevel VerbosityLevelType, logText ...interface{}) {
	if config.verbose == ALL {
		fmt.Println(logText)
	} else if verbosityLevel == config.verbose || verbosityLevel == ERROR {
		fmt.Println(logText)
	}
}
