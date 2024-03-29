// Copyright 2017 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package compiler

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

var fileCache map[string][]byte
var infoCache map[string]interface{}
var count int64

var verboseReader = false
var fileCacheEnable = true
var infoCacheEnable = true

func initializeFileCache() {
	if fileCache == nil {
		fileCache = make(map[string][]byte, 0)
	}
}

func initializeInfoCache() {
	if infoCache == nil {
		infoCache = make(map[string]interface{}, 0)
	}
}

func EnableFileCache() {
	fileCacheEnable = true
}

func EnableInfoCache() {
	infoCacheEnable = true
}

func DisableFileCache() {
	fileCacheEnable = false
}

func DisableInfoCache() {
	infoCacheEnable = false
}

func RemoveFromFileCache(fileurl string) {
	if !fileCacheEnable {
		return
	}
	initializeFileCache()
	delete(fileCache, fileurl)
}

func RemoveFromInfoCache(filename string) {
	if !infoCacheEnable {
		return
	}
	initializeInfoCache()
	delete(infoCache, filename)
}

func GetInfoCache() map[string]interface{} {
	if infoCache == nil {
		initializeInfoCache()
	}
	return infoCache
}

func ClearFileCache() {
	fileCache = make(map[string][]byte, 0)
}

func ClearInfoCache() {
	infoCache = make(map[string]interface{})
}

func ClearCaches() {
	ClearFileCache()
	ClearInfoCache()
}

// FetchFile gets a specified file from the local filesystem or a remote location.
func FetchFile(fileurl string) ([]byte, error) {
	var bytes []byte
	initializeFileCache()
	if fileCacheEnable {
		bytes, ok := fileCache[fileurl]
		if ok {
			if verboseReader {
				log.Printf("Cache hit %s", fileurl)
			}
			return bytes, nil
		}
		if verboseReader {
			log.Printf("Fetching %s", fileurl)
		}
	}
	response, err := http.Get(fileurl)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Error downloading %s: %s", fileurl, response.Status))
	}
	bytes, err = io.ReadAll(response.Body)
	if fileCacheEnable && err == nil {
		fileCache[fileurl] = bytes
	}
	return bytes, err
}

// ReadBytesForFile reads the bytes of a file.
func ReadBytesForFile(filename string) ([]byte, error) {
	// is the filename a url?
	fileurl, _ := url.Parse(filename)
	if fileurl.Scheme != "" {
		// yes, fetch it
		bytes, err := FetchFile(filename)
		if err != nil {
			return nil, err
		}
		return bytes, nil
	}
	// no, it's a local filename
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// ReadInfoFromBytes unmarshals a file as a yaml.MapSlice.
func ReadInfoFromBytes(filename string, bytes []byte) (interface{}, error) {
	initializeInfoCache()
	if infoCacheEnable {
		cachedInfo, ok := infoCache[filename]
		if ok {
			if verboseReader {
				log.Printf("Cache hit info for file %s", filename)
			}
			return cachedInfo, nil
		}
		if verboseReader {
			log.Printf("Reading info for file %s", filename)
		}
	}
	var info yaml.MapSlice
	err := yaml.Unmarshal(bytes, &info)
	if err != nil {
		return nil, err
	}
	if infoCacheEnable && len(filename) > 0 {
		infoCache[filename] = info
	}
	return info, nil
}

// ReadInfoForRef reads a file and return the fragment needed to resolve a $ref.
func ReadInfoForRef(basefile string, ref string) (interface{}, error) {
	initializeInfoCache()
	if infoCacheEnable {
		info, ok := infoCache[ref]
		if ok {
			if verboseReader {
				log.Printf("Cache hit for ref %s#%s", basefile, ref)
			}
			return info, nil
		}
		if verboseReader {
			log.Printf("Reading info for ref %s#%s", basefile, ref)
		}
	}
	count = count + 1
	basedir, _ := filepath.Split(basefile)
	parts := strings.Split(ref, "#")
	var filename string
	if parts[0] != "" {
		filename = parts[0]
		if _, err := url.ParseRequestURI(parts[0]); err != nil {
			// It is not an URL, so the file is local
			filename = basedir + parts[0]
		}
	} else {
		filename = basefile
	}
	bytes, err := ReadBytesForFile(filename)
	if err != nil {
		return nil, err
	}
	info, err := ReadInfoFromBytes(filename, bytes)
	if err != nil {
		log.Printf("File error: %v\n", err)
	} else {
		if len(parts) > 1 {
			path := strings.Split(parts[1], "/")
			for i, key := range path {
				if i > 0 {
					m, ok := info.(yaml.MapSlice)
					if ok {
						found := false
						for _, section := range m {
							if section.Key == key {
								info = section.Value
								found = true
							}
						}
						if !found {
							infoCache[ref] = nil
							return nil, NewError(nil, fmt.Sprintf("could not resolve %s", ref))
						}
					}
				}
			}
		}
	}
	if infoCacheEnable {
		infoCache[ref] = info
	}
	return info, nil
}
