// Copyright 2018 Yoshi Yamaguchi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"
)

func main() {
	http.HandleFunc("/grep/", fineAllLinesHandler)
	http.ListenAndServe("0.0.0.0:8080", nil)
}

// findAllLinesHandler calls findAllLines with URL parameter q.
// eg. http://localhost:8080/grep/?q=タイ
func fineAllLinesHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	q := values.Get("q")
	results := findAllLines(q, dataFiles)
	total := 0
	results.Range(func(k, v interface{}) bool {
		total += len(v.([]string))
		return true
	})
	fmt.Fprintf(w, "number of results: %v", total)
}

// findAllLines find lines in all files in filenames fs that
// matches regexp pattern. Results is the map of filename and all matched lines.
func findAllLines(pattern string, fs []string) sync.Map {
	var sm sync.Map
	var wg sync.WaitGroup
	for _, f := range fs {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			r, err := findLines(pattern, f)
			if err != nil {
				log.Println(err)
			}
			sm.Store(f, r)
		}(f)
	}
	wg.Wait()
	return sm
}

// findLines find lines from filename f that matches regexp pattern.
func findLines(pattern string, f string) ([]string, error) {
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		if ok, _ := regexp.MatchString(pattern, s); ok {
			result = append(result, s)
		}
	}
	return result, nil
}
