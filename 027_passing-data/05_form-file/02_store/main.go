package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

// 다른 사람에게 구축하는 것을 부탁하지 마세요 어디에 얽매이면 그 틀에 갇히게 됩니다.
// 이게 표준라이브러리를 이용한 진짜 프로그래밍입니다. 이게 켄톰슨, 롭파이크가 생각했던 지향했던 Go 프로그래밍방식입니다.
// 간결하고 유연하게 프로그래밍하세요~~ -교수님? 강사님? ~~~
var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
	fmt.Println("tpl", tpl)
}

func main() {
	http.HandleFunc("/", foo)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}

func foo(w http.ResponseWriter, req *http.Request) {
	var s string
	if req.Method == http.MethodPost {

		// open
		f, h, err := req.FormFile("q")
		fmt.Println("f", f, "h", h)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		// for your information
		fmt.Println("\nfile:", f, "\nheader:", h, "\nerr", err)

		// read
		bs, err := ioutil.ReadAll(f)
		fmt.Println("bs", bs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s = string(bs)
		fmt.Println("s", s)

		// store on server
		dst, err := os.Create(filepath.Join("./user/", h.Filename))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = dst.Write(bs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tpl.ExecuteTemplate(w, "index.gohtml", s)
}
