package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	path   = flag.String("path", "D:/Infogenesis/Log_dir/log_pms", "path to directory")
	outDir = flag.String("out", "D:/Infogenesis/Log_dir/log_pms/BinBackup", "path to output directory")
	logDir = flag.String("log", "D:/Infogenesis/Log_dir/log_pms/", "path to log directory")
)
var logger *log.Logger

func init() {
	flag.Parse()
	initLogger()
	logger.Println("Program started at:", time.Now().Format(time.RFC3339))
}

func initLogger() {
	logFileName := "BinMover" + time.Now().Format("2006-01-02") + ".log"
	logFile, err := os.OpenFile(filepath.Join(*logDir, logFileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	logger = log.New(logFile, "", log.LstdFlags)
}

func main() {
	// Construct the bin backup directory with the current date
	binBackupDir := filepath.Join(*outDir, time.Now().Format("2006-01-02"))

	// Ensure the directory exists
	if _, err := os.Stat(binBackupDir); os.IsNotExist(err) {
		if err := os.Mkdir(binBackupDir, 0777); err != nil {
			log.Fatal(err)
		}
	}

	entries, err := os.ReadDir(*path)
	if err != nil {
		log.Fatal(err)
	}

	binCount := 0
	movedCount := 0
	for _, entry := range entries {
		if strings.HasSuffix(strings.ToUpper(entry.Name()), ".BIN") {
			binCount++
			err := moveFile(filepath.Join(*path, entry.Name()), binBackupDir)
			if err != nil {
				logger.Printf("Failed to move file %s: %v", entry.Name(), err)
			} else {
				movedCount++
			}
		}
	}

	logger.Printf("Number of .BIN files in the directory: %d", binCount)
	logger.Printf("Number of moved .BIN files: %d", movedCount)
	logger.Println("Program exited at:", time.Now().Format(time.RFC3339))
}

func moveFile(srcPath, destDir string) error {
	destPath := filepath.Join(destDir, filepath.Base(srcPath))
	return os.Rename(srcPath, destPath)
}
