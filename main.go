package main

import (
	// imageutils "picmask/imageUtils"
	"picmask/server"
)

func main() {
	// path := "./images/test1.jpeg"
	// imageutils.ProcessImage(path, true)
	// path = "output.jpeg"
	// imageutils.ProcessImage(path, false)
	r := server.Router()
	r.Run("0.0.0.0:8888")
	// test()
}

// func test() {
// 	path := "./images/enc125.jpeg"
// 	// path := "C:/Users/ADMINI~1\\AppData\\Local\\Temp\\go-build3862783964\\121_WMAXzsxg6H8518fc5b2ddcd44328247c7ed6e1a22095.jpeg"
// 	out := "./images/dec125.jpeg"
// 	imageutils.ProcessImage(path, 125, false, out)
// }
