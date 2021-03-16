package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {

	start := time.Now()

	folderPath1 := flag.String("input1", "", "Folder to be compared in left hand side")
	folderPath2 := flag.String("input2", "", "Folder to be compared in right hand side")
	flag.Parse()

	var md5Res1, md5Res2 = make(map[string]string), make(map[string]string)

	fmt.Println("Running...")

	wg.Add(2)
	go func() {
		md5Res1 = calMd5(*folderPath1, "checksum_input1.txt")
	}()
	go func() {
		md5Res2 = calMd5(*folderPath2, "checksum_input2.txt")
	}()
	wg.Wait()

	compare(md5Res1, md5Res2)

	duration := time.Since(start)
	fmt.Printf("Done. Duration: %s", duration.String())
}

func calMd5(path string, resultName string) map[string]string {
	defer wg.Done()

	fileMap := make(map[string]string)
	resultFile, resultWriteErr := os.Create(resultName)
	if resultWriteErr != nil {
		log.Fatal(resultWriteErr)
	}

	filepathWalkErr := filepath.Walk(path, func(fullPath string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.Index(info.Name(), ".") != 0 {
			file, err := os.Open(fullPath)
			if err != nil {
				log.Fatal(err)
			}

			hash := md5.New()
			if _, err := io.Copy(hash, file); err != nil {
				log.Fatal(err)
			}

			hashStr := hex.EncodeToString(hash.Sum(nil))
			tmpPath := strings.Replace(fullPath, path, "", 1)

			fileMap[tmpPath] = hashStr
			file.Close()

			resultFile.WriteString( fullPath + " > " + hashStr + "\n")
		}
		return nil
	})

	if resultWriteErr != nil {
		resultFile.Close()
	}

	if filepathWalkErr != nil {
		panic(filepathWalkErr)
	}
	return fileMap
}

func compare(fileMap1, fileMap2 map[string]string) {
	var resFm1MissInFm2, resFm2MissInFm1, resMismatch []string

	for k, v := range fileMap1 {
		if _, ok := fileMap2[k]; !ok {
			resFm1MissInFm2 = append(resFm1MissInFm2, k)
		} else if w, _ := fileMap2[k]; v != w {
			resMismatch = append(resMismatch, k)
		}
	}

	for k, _ := range fileMap2 {
		if _, ok := fileMap1[k]; !ok {
			resFm2MissInFm1 = append(resFm2MissInFm1, k)
		}
	}

	fmt.Println("==Mismatched checksum==")
	if len(resMismatch) > 0 {
		fmt.Println(strings.Join(resMismatch[:], "\n"))
	} else {
		fmt.Println("N/A")
	}

	fmt.Println("==Input 1 file(s) not found in input 2==")
	if len(resFm1MissInFm2) > 0 {
		fmt.Println(strings.Join(resFm1MissInFm2[:], "\n"))
	} else {
		fmt.Println("N/A")
	}

	fmt.Println("==Input 2 file(s) not found in input 1==")
	if len(resFm2MissInFm1) > 0 {
		fmt.Println(strings.Join(resFm2MissInFm1[:], "\n"))
	} else {
		fmt.Println("N/A")
	}
}