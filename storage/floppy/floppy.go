package floppy

import (
	"io/ioutil"
	"os"
	"path/filepath"
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
	go f.watchFloppy()
}

func (f *Floppy) ListFiles() []os.FileInfo {
	return f.FileIndex
}

func (m *Floppy) LoadFile(f os.FileInfo) ([]byte, error) {
	dat, err := ioutil.ReadFile(filepath.Join(m.mountpoint, f.Name()))
	if err != nil {
		return nil, err
	}
	return dat, nil
}

func (m *Floppy) checkMountStatus() bool {

}

// func (d *Disk) Watch() {
// 	cmd := exec.Command("diskd", "-d", d.device).Run()
// 	fmt.Println(cmd)
// 	d.mount()
// }

// func (d *Disk) mount() {
// 	cmd := exec.Command("mount", d.mountpoint).Run()
// 	fmt.Println(cmd)

// 	d.indexDisk()
// }

// func (d *Disk) indexDisk() {
// 	files, err := ioutil.ReadDir(d.mountpoint)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	d.FileIndex = files

// 	for _, file := range d.FileIndex {
// 		fmt.Println(file.Name())
// 	}
// }
