package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
)

const (
	NumOfChairSearchData                 = 100
	NumOfEstateSearchData                = 100
	NumOfRecommendedEstatesWithChairData = 100
	NumOfEstatesNazotteData              = 100
)

func init() {
	rand.Seed(19700101)
}

func writeSnapshotDataToFile(path string, snapshot Snapshot) {
	bytes, err := json.Marshal(snapshot)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(path, bytes, os.FileMode(0777))
	if err != nil {
		panic(err)
	}
}

func main() {
	flags := flag.NewFlagSet("isucon10-qualify", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)

	var TargetServer string
	var DestDirectoryPath string

	flags.StringVar(&TargetServer, "target-url", "http://127.0.0.1:1323", "target url")
	flags.StringVar(&DestDirectoryPath, "dest-dir", "./result/verification_data", "destination directory")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}

	MkdirIfNotExists(DestDirectoryPath)

	// chair search
	MkdirIfNotExists(filepath.Join(DestDirectoryPath, "chair_search"))
	for i := 0; i < NumOfChairSearchData; i++ {
		wg.Add(1)
		go func(id int) {
			req := Request{
				Method:   "GET",
				Resource: "/api/chair/search",
				Query:    createRandomChairSearchQuery().Encode(),
				Body:     "",
			}

			snapshot := getSnapshotFromRequest(TargetServer, req)

			filename := fmt.Sprintf("%d.json", id)
			writeSnapshotDataToFile(filepath.Join(DestDirectoryPath, "chair_search", filename), snapshot)
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Println("Done generating verification data of /api/chair/search")

	// estate search
	MkdirIfNotExists(filepath.Join(DestDirectoryPath, "estate_search"))
	for i := 0; i < NumOfEstateSearchData; i++ {
		wg.Add(1)
		go func(id int) {
			req := Request{
				Method:   "GET",
				Resource: "/api/estate/search",
				Query:    createRandomEstateSearchQuery().Encode(),
				Body:     "",
			}

			snapshot := getSnapshotFromRequest(TargetServer, req)
			filename := fmt.Sprintf("%d.json", id)
			writeSnapshotDataToFile(filepath.Join(DestDirectoryPath, "estate_search", filename), snapshot)
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Println("Done generating verification data of /api/estate/search")

	// recommended_chair
	MkdirIfNotExists(filepath.Join(DestDirectoryPath, "recommended_chair"))
	wg.Add(1)
	go func() {
		req := Request{
			Method:   "GET",
			Resource: "/api/recommended_chair",
			Query:    "",
			Body:     "",
		}

		snapshot := getSnapshotFromRequest(TargetServer, req)
		writeSnapshotDataToFile(filepath.Join(DestDirectoryPath, "recommended_chair", "0.json"), snapshot)
		wg.Done()
	}()
	wg.Wait()
	log.Println("Done generating verification data of /api/recommended_chair")

	// recommended_estate
	MkdirIfNotExists(filepath.Join(DestDirectoryPath, "recommended_estate"))
	wg.Add(1)
	go func() {
		req := Request{
			Method:   "GET",
			Resource: "/api/recommended_estate",
			Query:    "",
			Body:     "",
		}

		snapshot := getSnapshotFromRequest(TargetServer, req)
		writeSnapshotDataToFile(filepath.Join(DestDirectoryPath, "recommended_estate", "0.json"), snapshot)
		wg.Done()
	}()
	wg.Wait()
	log.Println("Done generating verification data of /api/recommended_estate")

	// recommended_estate/:id
	MkdirIfNotExists(filepath.Join(DestDirectoryPath, "recommended_estate_with_chair"))
	for i := 0; i < NumOfRecommendedEstatesWithChairData; i++ {
		wg.Add(1)
		go func(id int) {
			req := Request{
				Method:   "GET",
				Resource: fmt.Sprintf("/api/recommended_estate/%d", id),
				Query:    "",
				Body:     "",
			}

			snapshot := getSnapshotFromRequest(TargetServer, req)
			fileName := fmt.Sprintf("%d.json", id)
			writeSnapshotDataToFile(filepath.Join(DestDirectoryPath, "recommended_estate_with_chair", fileName), snapshot)
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Println("Done generating verification data of /api/recommended_estate/:id")

	// estate nazotte
	MkdirIfNotExists(filepath.Join(DestDirectoryPath, "estate_nazotte"))
	for i := 0; i < NumOfEstatesNazotteData; i++ {
		wg.Add(1)
		go func(id int) {
			req := Request{
				Method:   "POST",
				Resource: "/api/estate/nazotte",
				Query:    "",
				Body:     createRandomConvexhull(),
			}

			snapshot := getSnapshotFromRequest(TargetServer, req)
			fileName := fmt.Sprintf("%d.json", id)
			writeSnapshotDataToFile(filepath.Join(DestDirectoryPath, "estate_nazotte", fileName), snapshot)
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Println("Done generating verification data of /api/estate/nazotte")
}
