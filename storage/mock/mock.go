package mock

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/coral/mt-420/storage"
)

type Mock struct {
	basepath string
}

func New(basepath string) *Mock {
	return &Mock{
		basepath: basepath,
	}
}

func (m *Mock) Init() {

}

func (m *Mock) ListFiles() []os.FileInfo {
	files, err := ioutil.ReadDir(m.basepath)
	if err != nil {
		log.Fatal(err)
	}

	var t []os.FileInfo

	for _, file := range files {
		if !file.IsDir() {
			if filepath.Ext(file.Name()) == ".mid" || filepath.Ext(file.Name()) == ".midi" {
				t = append(t, file)
			}
		}
	}

	return t
}

func (m *Mock) LoadFile(f os.FileInfo) ([]byte, error) {
	dat, err := ioutil.ReadFile(filepath.Join(m.basepath, f.Name()))
	if err != nil {
		return nil, err
	}
	return dat, nil
}

func (m *Mock) Watch() <-chan storage.Events {

	return make(<-chan storage.Events)
}
