package floppy

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type Disk struct {
	device     string
	mountpoint string
	fileindex  []os.FileInfo
}

func New(device string, mountpoint string) *Disk {
	return &Disk{
		device:     device,
		mountpoint: mountpoint,
	}
}

func (d *Disk) Watch() {
	cmd := exec.Command("diskd", "-d", d.device).Run()
	fmt.Println(cmd)
	d.mount()
}

func (d *Disk) mount() {
	cmd := exec.Command("mount", d.mountpoint).Run()
	fmt.Println(cmd)

	d.indexDisk()
}

func (d *Disk) indexDisk() {
	files, err := ioutil.ReadDir(d.mountpoint)
	if err != nil {
		log.Fatal(err)
	}

	d.fileindex = files

	for _, file := range d.fileindex {
		fmt.Println(file.Name())
	}
}
