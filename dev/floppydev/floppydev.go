package main

import "github.com/coral/mt-420/storage/floppy"

func main() {
	//Floppy
	//fl := floppy.New("/dev/fd0", "/media/floppy")
	storage := floppy.New("/dev/sdb", "/media/floppy")
	storage.Init()
}
