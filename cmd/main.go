package main

import (
	"flag"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	path   = flag.String("path", "D:/Infogenesis/Log_dir/log_pms", "relative path to directory")
	outDir = flag.String("out", "D:/Infogenesis/Log_dir/log_pms/BinBackup", "relative path to output directory")
)

const maxGoroutines = 100 // Adjust as needed

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

	sem := make(chan struct{}, maxGoroutines)
	var wg sync.WaitGroup
	errCh := make(chan error)

	// Filter and move .BIN files
	for _, entry := range entries {
		if strings.HasSuffix(strings.ToUpper(entry.Name()), ".BIN") {
			wg.Add(1)
			go func(entry fs.DirEntry) {
				defer wg.Done()

				sem <- struct{}{} // Acquire semaphore
				err := moveFile(filepath.Join(*path, entry.Name()), binBackupDir)
				<-sem // Release semaphore

				if err != nil {
					errCh <- err
				}
			}(entry)
		}
	}

	// Close the error channel after all goroutines have completed
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// Handle errors
	for err := range errCh {
		log.Println("Error:", err)
	}
}

func moveFile(srcPath, destDir string) error {
	destPath := filepath.Join(destDir, filepath.Base(srcPath))
	return os.Rename(srcPath, destPath)
}
