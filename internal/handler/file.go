package handler

// import (
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"os"
// )

// func upload(w http.ResponseWriter, r *http.Request) {
// 	r.ParseMultipartForm(32 << 20)
// 	file, handler, err := r.FormFile("file")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer file.Close()
// 	// fmt.Fprintf(w, "%v", handler.Header)
// 	f, err := os.Create(handler.Filename)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer f.Close()
// 	io.Copy(f, file)
// }
