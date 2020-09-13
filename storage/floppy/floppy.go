package floppy

import (
	"os"
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

}

func (f *Floppy) ListFiles() []os.FileInfo {
	return f.FileIndex
}

func ()

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
