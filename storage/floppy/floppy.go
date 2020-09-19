package floppy

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Floppy struct {
	device     string
	mountpoint string
	FileIndex  []os.FileInfo
	mounted    bool
}

func New(device string, mountpoint string) *Floppy {
	return &Floppy{
		device:     device,
		mountpoint: mountpoint,
	}
}

func (f *Floppy) Init() {
	if f.checkMountStatus() {
		f.mounted = true
	} else {
		fmt.Println(f.mountFloppy())
		time.Sleep(1 * time.Second)
		f.Init()
	}
}

func (f *Floppy) ListFiles() []os.FileInfo {
	if f.mnt() {
		time.Sleep(100 * time.Millisecond)
	}

	files, err := ioutil.ReadDir(f.mountpoint)
	if err != nil {
		log.Fatal(err)
	}

	var t []os.FileInfo

	for _, file := range files {
		if !file.IsDir() {
			if !strings.HasPrefix(file.Name(), ".") {
				if filepath.Ext(file.Name()) == ".mid" || filepath.Ext(file.Name()) == ".midi" {

					t = append(t, file)
				}
			}
		}
	}

	return t
}

func (m *Floppy) LoadFile(f os.FileInfo) ([]byte, error) {
	dat, err := ioutil.ReadFile(filepath.Join(m.mountpoint, f.Name()))
	if err != nil {
		return nil, err
	}
	return dat, nil
}

func (m *Floppy) mnt() bool {
	if !m.checkMountStatus() {
		err := m.mountFloppy()
		if err != nil {
			m.mounted = true
			return true
		}
	}

	m.mounted = false
	return false
}

func (m *Floppy) checkMountStatus() bool {

	content, err := ioutil.ReadFile("/proc/mounts")
	if err != nil {
		fmt.Println(err)
	}
	lines := strings.Split(string(content), "\n")

	for _, mt := range lines {
		if strings.Contains(mt, "/media/floppy") {
			return true
		}
	}

	return false

}

func (m *Floppy) mountFloppy() error {

	return exec.Command("mount", m.device).Run()
}
