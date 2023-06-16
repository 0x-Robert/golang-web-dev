package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", index)

	http.HandleFunc("/set", set)
	http.HandleFunc("/read", read)
	http.HandleFunc("/expire", expire)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, req *http.Request) {
	// set으로 가서 쿠키 설정
	fmt.Fprintln(w, `<h1><a href="/set">set a cookie</a></h1>`)
}

// 쿠키 설정
func set(w http.ResponseWriter, req *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:  "session",
		Value: "some value",
		Path:  "/",
	})
	fmt.Fprintln(w, `<h1><a href="/read">read</a></h1>`)
}

// 쿠키를 읽기
func read(w http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie("session")
	if err != nil {
		http.Redirect(w, req, "/set", http.StatusSeeOther)
		// 리턴안해주면 계속 진행
		return
	}

	fmt.Fprintf(w, `<h1>Your Cookie:<br>%v</h1><h1><a href="/expire">expire</a></h1>`, c)
}

// 맥스에이지를 사용해서 쿠키를 관리하며 0이나 음수를 설정하면 쿠키가 삭제됨
// 이후 맥스에이지 값으로 쿠키 설정 후 루트로 다시 이동 그렇게 하면 전체 쿠키 설정, 읽기, 삭제까지 순환루프가 만들어진다.
func expire(w http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie("session")
	if err != nil {
		http.Redirect(w, req, "/set", http.StatusSeeOther)
		return
	}
	c.MaxAge = -1 // delete cookie
	http.SetCookie(w, c)
	http.Redirect(w, req, "/", http.StatusSeeOther)
}
