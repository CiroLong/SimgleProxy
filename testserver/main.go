package main

import "net/http"

func main() {

	http.HandleFunc("/test", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("hello simgle proxy v1!"))
	})

	http.ListenAndServe(":8081", nil)
}
