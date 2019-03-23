package hdfs

import (
	"fmt"
	"net/http"
)

func Mv(w http.ResponseWriter, r *http.Request) {
	src := r.URL.Query().Get("src")
	dst := r.URL.Query().Get("dst")
	moved, _ := mv(src, dst)
	if moved {
		fmt.Fprintf(w, "mv %s %s \n", src, dst)
	} else {
		fmt.Fprintf(w, "mv %s %s failed \n", src, dst)
	}
}

func GetContentSummary(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	summary, err := getContentSummary(path)
	if err != nil {
		fmt.Fprintf(w, "couldn't get content summary of file %s \n", path)
	} else {
		fmt.Fprintf(w, summary)
	}
}

func CreateFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	created, _ := createEmptyFile(path)
	if created {
		fmt.Fprintf(w, "created file %s \n", path)
	} else {
		fmt.Fprintf(w, "couldn't create file %s \n", path)
	}
}