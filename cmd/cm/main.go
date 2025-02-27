package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

// タイムスタンプの正規表現 (例: 20250217165636.md)
var timestampRegex = regexp.MustCompile(`^(\d{4})(\d{2})(\d{2})\d{6}\.(\w+)$`)

func main() {
	path := flag.String("path", "~/memo", "memo directory path")
	ext := flag.String("ext", "md", "file extension for new memo (default: md)")
	flag.Parse()

	expandedPath, err := expandPath(*path)
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = checkDirExitsOrCreate(expandedPath)
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = archiveOldFiles(expandedPath, *ext)
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = createNewMemo(expandedPath, *ext)
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func expandPath(path string) (string, error) {
	if path[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("Failed to get home directory: %v", err)
		}
		return filepath.Join(home, path[2:]), nil
	}
	return path, nil
}

func checkDirExitsOrCreate(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return fmt.Errorf("Failed to create directory %s: %v", path, err)
	}
	return nil
}

// 1ヶ月以上前のファイルを YYYY/MM に移動（空のファイルは削除）
func archiveOldFiles(memoDir, ext string) error {
	threshold := time.Now().AddDate(0, -1, 0)

	files, err := os.ReadDir(memoDir)
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()

		matches := timestampRegex.FindStringSubmatch(name)
		if matches == nil {
			continue
		}

		year, err := strconv.Atoi(matches[1])
		if err != nil {
			return fmt.Errorf("Failed to convert year: %v", err)
		}
		month, err := strconv.Atoi(matches[2])
		if err != nil {
			return fmt.Errorf("Failed to convert month: %v", err)
		}
		day, err := strconv.Atoi(matches[3])
		if err != nil {
			return fmt.Errorf("Failed to convert day: %v", err)
		}

		fileTime := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

		oldPath := filepath.Join(memoDir, name)
		// ファイルが空なら削除
		if isEmptyFile(oldPath) {
			if err := os.Remove(oldPath); err != nil {
				log.Printf("Failed to delete empty file %s: %v\n", oldPath, err)
			} else {
				log.Printf("Deleted empty file: %s\n", oldPath)
			}
			continue // 空のファイルは削除したので、移動処理は不要
		}

		// 1ヶ月以上前なら移動
		if fileTime.Before(threshold) {
			destDir := filepath.Join(memoDir, fmt.Sprintf("%04d/%02d", year, month))
			destPath := filepath.Join(destDir, name)

			if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
				log.Printf("Failed to create directory %s: %v\n", destDir, err)
				continue
			}

			if err := os.Rename(oldPath, destPath); err != nil {
				log.Printf("Failed to move %s to %s: %v\n", oldPath, destPath, err)
				continue
			}

			log.Printf("Moved: %s -> %s\n", oldPath, destPath)
		}
	}

	return nil
}

func isEmptyFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		log.Printf("Failed to check file size: %v\n", err)
		return false
	}
	return info.Size() == 0
}

func createNewMemo(memoDir, ext string) error {
	var file *os.File
	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	timestamp := time.Now().Format("20060102150405")
	filePath := filepath.Join(memoDir, fmt.Sprintf("%s.%s", timestamp, ext))

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("Failed to create file: %v", err)
	}

	fmt.Println("Created:", filePath)
	return nil
}
