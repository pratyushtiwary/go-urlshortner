package shortner

import "path/filepath"

var datastoreDir string = filepath.Join(".", "datastore")

type config struct {
	DatastoreDir string
	IndexFile    string
}

var Config *config = &config{
	DatastoreDir: datastoreDir,
	IndexFile:    filepath.Join(datastoreDir, "index.json"),
}
