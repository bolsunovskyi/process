package main

import (
	"fmt"
	"path"
)

func main() {
	//proc, err := process.CreateProcess("./example_process/example", "-ct=3600")
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//if err := proc.Start(); err != nil {
	//	log.Fatalln(err)
	//}
	//
	//if err := proc.Wait(); err != nil {
	//	log.Println(err)
	//}
	fmt.Println(path.Clean("../foo/bar/"))
	fmt.Println(path.Base("../foo/bar"))
}
