package floppy

import (
	"fmt"
	"os/exec"
)

type Disk struct {
	device     string
	mountpoint string
}

func New(device string) *Disk {
	return &Disk{
		device: device,
	}
}

func (d *Disk) Watch() {
	cmd := exec.Command("diskd", "-d", d.device).Run()
	fmt.Println(cmd)
}
