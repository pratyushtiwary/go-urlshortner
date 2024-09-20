package shortner

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

/* TYPES */
type ShortenRequest struct {
	Url  string
	Name string
}

// Key = name, value = url
type StoreIndex = map[string]string

/* FUNCTIONS */

func createDatastore() error {
	err := os.MkdirAll(Config.DatastoreDir, os.ModePerm)

	if err != nil {
		return err
	}
	return nil
}

func closeFile(file *os.File) {
	file.Close()
}

func createIndex() (*os.File, error) {
	err := createDatastore()

	if err != nil {
		return nil, err
	}

	_, statErr := os.Stat(Config.IndexFile)

	var file *os.File
	var opErr error

	if errors.Is(statErr, os.ErrNotExist) {
		// create file
		file, opErr = os.Create(Config.IndexFile)

		if opErr != nil {
			return nil, opErr
		}
	} else {
		file, opErr = os.OpenFile(Config.IndexFile, os.O_CREATE|os.O_WRONLY, 0666)

		if opErr != nil {
			return nil, opErr
		}
	}

	return file, nil
}

func readIndex() (*StoreIndex, error) {
	file, err := createIndex()

	defer closeFile(file)

	if err != nil {
		return nil, err
	}

	indexFileBytes, readErr := os.ReadFile(Config.IndexFile)

	if readErr != nil {
		return nil, readErr
	}

	indexContents := StoreIndex{}
	json.Unmarshal(indexFileBytes, &indexContents)

	return &indexContents, nil
}

func writeIndex(data *StoreIndex) error {
	indexFile, err := createIndex()

	defer closeFile(indexFile)

	if err != nil {
		return err
	}

	indexContentsEncoded, encodingErr := json.Marshal(data)

	if encodingErr != nil {
		return encodingErr
	}

	_, writeErr := indexFile.Write(indexContentsEncoded)

	if writeErr != nil {
		return writeErr
	}

	return nil
}

func saveData(data *ShortenRequest) error {
	dsErr := createDatastore()

	if dsErr != nil {
		return dsErr
	}

	indexContents, indexReadErr := readIndex()

	if indexReadErr != nil {
		return indexReadErr
	}

	(*indexContents)[data.Name] = data.Url

	indexWriteError := writeIndex(indexContents)

	if indexWriteError != nil {
		return indexWriteError
	}

	return nil
}

func getUrl(name string) (*string, error) {
	dsErr := createDatastore()

	if dsErr != nil {
		return nil, dsErr
	}

	index, readIndexErr := readIndex()

	if readIndexErr != nil {
		return nil, readIndexErr
	}

	value, exists := (*index)[name]

	if !exists {
		return nil, errors.New("provided name doesn't exists")
	}

	return &value, nil
}

func shorten(request *ShortenRequest) string {
	err := saveData(request)

	if err != nil {
		log.Fatal("Failed to save data ", err)
	}
	return "go/" + request.Name
}

func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	data := &ShortenRequest{}
	decoder.Decode(data)

	fmt.Fprint(w, shorten(data))
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name, exists := params["name"]

	if !exists {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Invalid name provided")
	}

	createDatastore()
	index, indexReadError := readIndex()

	if indexReadError != nil {
		log.Fatal(indexReadError)
	}

	url, urlExists := (*index)[name]

	if !urlExists {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Invalid name provided")
	}

	http.Redirect(w, r, url, http.StatusSeeOther)
}
