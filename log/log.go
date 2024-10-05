package log

import (
	"Peregrine/stru"
	"encoding/json"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var (

	// 前台Logger
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger

	// 文件Logger
	debugFileLogger *log.Logger
	infoFileLogger  *log.Logger
	warnFileLogger  *log.Logger
	errorFileLogger *log.Logger
	fatalFileLogger *log.Logger

	// JsonLogger
	debugJsonLogger *log.Logger
	infoJsonLogger  *log.Logger
	warnJsonLogger  *log.Logger
	errorJsonLogger *log.Logger
	fatalJsonLogger *log.Logger

	logOut      *os.File
	logJsonOut  *os.File
	logLevel    int
	currentDay  int
	logFile     string
	logJsonFile string
	fileLock    sync.RWMutex

	jsonStatus bool = false
	fileStatus bool = false
)

const (
	DebugLevel = iota // 0
	InfoLevel         // 1
	WarnLevel         // 2
	ErrorLevel        // 3
	FatalLevel        // 4
)

const (
	// 颜色重置
	colorReset = "\033[0m"
	// 红色
	colorRed = "\033[31m"
	// 黄色
	colorYellow = "\033[33m"
	// 青色
	colorCyan = "\033[36m"
	// 灰色
	colorGray = "\033[90m"
)

type jsonMsg struct {
	Level   string `json:"level"`
	Msg     string `json:"msg"`
	Time    int64  `json:"time"`
	Details any    `json:"details,omitempty"` // 可选字段
}

func InitLogOut(cfg stru.Log) {
	fileLock = sync.RWMutex{}
	setLevel(cfg.Level)
	fileStatus = cfg.File
	jsonStatus = cfg.Json
	if cfg.File {
		setFile(cfg.Path + "/peregrine.log")
	}
	if cfg.Json {
		setJsonFile(cfg.Path + "/peregrine.jlog")
	}
	// 获取今天是当年的第几天
	currentDay = time.Now().YearDay()
	initLog(logOut, logJsonOut)

}

func setLevel(level int) {
	logLevel = level
}

func setFile(file string) {
	logFile = file
	var err error
	logOut, err = os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		panic(err)
	}
}
func setJsonFile(file string) {
	logJsonFile = file
	var err error
	logJsonOut, err = os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		panic(err)
	}
}

func Debug(format string, v ...any) {
	checkIfDayChange()
	if logLevel <= DebugLevel {
		debugLogger.Printf(format, v...)
		if fileStatus {
			debugFileLogger.Printf(format, v...)
		}
		if jsonStatus {
			str := getJsonDataStr("debug", format, v...)
			debugJsonLogger.Printf(str)
		}
	}
}
func Info(format string, v ...any) {
	checkIfDayChange()
	if logLevel <= InfoLevel {
		infoLogger.Printf(format, v...)
		if fileStatus {
			infoFileLogger.Printf(format, v...)
		}
		if jsonStatus {
			str := getJsonDataStr("info", format, v...)
			infoJsonLogger.Printf(str)
		}
	}
}
func Warn(format string, v ...any) {
	checkIfDayChange()
	if logLevel <= WarnLevel {
		warnLogger.Printf(format, v...)
		if fileStatus {
			warnFileLogger.Printf(format, v...)
		}
		if jsonStatus {
			str := getJsonDataStr("warn", format, v...)
			warnJsonLogger.Printf(str)
		}
	}
}
func Error(format string, v ...any) {
	checkIfDayChange()
	if logLevel <= ErrorLevel {
		errorLogger.Printf(getPrefix()+format, v...)
		if fileStatus {
			errorFileLogger.Printf(format, v...)
		}
		if jsonStatus {
			str := getJsonDataStr("error", format, v...)
			errorJsonLogger.Printf(str)
		}
	}
}
func Fatal(format string, v ...any) {
	checkIfDayChange()
	if logLevel <= FatalLevel {
		fatalLogger.Printf(getPrefix()+format, v...)
		if fileStatus {
			fatalFileLogger.Printf(format, v...)
		}
		if jsonStatus {
			str := getJsonDataStr("fatal", format, v...)
			fatalJsonLogger.Printf(str)
		}
	}
	os.Exit(1)
}

func getCallTrace() (string, int) {
	// 函数名 文件 行号 是否出现异常
	//pc, file, line, ok := runtime.Caller(0)

	_, file, line, ok := runtime.Caller(3)
	if ok {
		return file, line
	}
	return "", 0
}

func getPrefix() string {
	file, line := getCallTrace()

	return file + ":" + strconv.Itoa(line) + " "
}

func checkIfDayChange() {
	// 锁
	fileLock.Lock()
	defer fileLock.Unlock()
	day := time.Now().YearDay()
	if day == currentDay {
		return
	} else {
		currentDay = day
		logOut.Close()
		postFix := time.Now().Add(-24 * time.Hour).Format("20060102")
		if fileStatus {
			os.Rename(logFile, logFile+"."+postFix)
			logOut, _ = os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
		}
		if jsonStatus {
			os.Rename(logJsonFile, logJsonFile+"."+postFix)
			logJsonOut, _ = os.OpenFile(logJsonFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
		}
		initLog(logOut, logJsonOut)
	}
}

func initLog(logOut, logJsonOut *os.File) {
	// 前台Logger
	infoLogger = log.New(os.Stdout, colorCyan+"[INFO] "+colorReset, log.LstdFlags)
	debugLogger = log.New(os.Stdout, colorGray+"[DEBUG] "+colorReset, log.LstdFlags)
	warnLogger = log.New(os.Stdout, colorYellow+"[WARN] "+colorReset, log.LstdFlags)
	errorLogger = log.New(os.Stdout, colorRed+"[ERROR] "+colorReset, log.LstdFlags)
	fatalLogger = log.New(os.Stdout, colorRed+"[FATAL] "+colorReset, log.LstdFlags)

	// 文件Logger
	if fileStatus {
		infoFileLogger = log.New(logOut, "[INFO] ", log.LstdFlags)
		debugFileLogger = log.New(logOut, "[DEBUG] ", log.LstdFlags)
		warnFileLogger = log.New(logOut, "[WARN] ", log.LstdFlags)
		errorFileLogger = log.New(logOut, "[ERROR] ", log.LstdFlags)
		fatalFileLogger = log.New(logOut, "[FATAL] ", log.LstdFlags)
	}

	// json格式Logger
	if jsonStatus {
		infoJsonLogger = log.New(logJsonOut, "", 0)
		debugJsonLogger = log.New(logJsonOut, "", 0)
		warnJsonLogger = log.New(logJsonOut, "", 0)
		errorJsonLogger = log.New(logJsonOut, "", 0)
		fatalJsonLogger = log.New(logJsonOut, "", 0)
	}

}

func getJsonDataStr(level, msg string, details ...any) string {
	data := jsonMsg{
		Level:   level,
		Msg:     msg,
		Details: details,
		Time:    time.Now().Unix(),
	}
	byt, err := json.Marshal(data)
	if err != nil {
		log.Println("出大问题，json数据格式化失败，日志系统出现问题，请联系管理员....")
	}
	return string(byt)
}
