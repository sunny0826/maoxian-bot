package cmd

//var proPath string

//func InitLog() {
//	logFile, err := GetLogFile()
//	if err != nil {
//		return
//	}
//	output := &lumberjack.Logger{
//		Filename:   logFile,
//		MaxSize:    30, // m
//		MaxBackups: 1,
//		MaxAge:     2, // days
//		Compress:   false,
//	}
//	logrus.SetOutput(&fileWriterWithoutErr{output})
//
//	formatter := &logrus.TextFormatter{
//		FullTimestamp:   true,
//		TimestampFormat: time.RFC3339Nano,
//	}
//	logrus.SetFormatter(formatter)
//
//	if Debug {
//		logrus.SetLevel(logrus.DebugLevel)
//	}
//}
//
//// GetLogFile
//func GetLogFile() (string, error) {
//	logPath, err := GetLogPath()
//	if err != nil {
//		return "", err
//	}
//	logFile := path.Join(logPath, BotLog)
//	return logFile, nil
//}
//
//// GetProgramPath
//func GetProgramPath() string {
//	if proPath != "" {
//		return proPath
//	}
//	dir, err := exec.LookPath(os.Args[0])
//	if err != nil {
//		log.Fatal("can get the process path")
//	}
//	if p, err := os.Readlink(dir); err == nil {
//		dir = p
//	}
//	proPath, err = filepath.Abs(filepath.Dir(dir))
//	if err != nil {
//		log.Fatal("can get the full process path")
//	}
//	return proPath
//}
//
//func GetLogPath() (string, error) {
//	binDir := GetProgramPath()
//	logsPath := path.Join(binDir, "logs")
//	if IsExist(logsPath) {
//		return logsPath, nil
//	}
//	// mk dir
//	err := os.MkdirAll(logsPath, os.ModePerm)
//	if err != nil {
//		return "", err
//	}
//	return logsPath, nil
//}
//
////IsExist returns true if file exists
//func IsExist(fileName string) bool {
//	_, err := os.Stat(fileName)
//	return err == nil || os.IsExist(err)
//}
