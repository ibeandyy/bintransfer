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
	path   = flag.String("path", "D:/Infogenesis/Log_dir/log_pms", "(absolute or relative) path to directory")
	outDir = flag.String("out", "D:/Infogenesis/Log_dir/log_pms/BinBackup", "(absolute or relative) path to output directory")
)

func main() {
	flag.Parse()

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

	// Filter and move .BIN files
	for _, entry := range entries {
		if strings.HasSuffix(strings.ToUpper(entry.Name()), ".BIN") {
			if err := moveFile(filepath.Join(*path, entry.Name()), binBackupDir); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func moveFile(srcPath, destDir string) error {
	destPath := filepath.Join(destDir, filepath.Base(srcPath))
	return os.Rename(srcPath, destPath)
}
