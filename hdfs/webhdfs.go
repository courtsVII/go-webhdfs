package hdfs

import (
	"fmt"
	"net/http"
	"os"
)

func Mv(w http.ResponseWriter, r *http.Request) {
	src := r.URL.Query().Get("src")
	dst := r.URL.Query().Get("dst")
	moved, _ := mv(&src, &dst)
	if moved {
		fmt.Fprintf(w, "mv %s %s \n", src, dst)
	} else {
		fmt.Fprintf(w, "mv %s %s failed \n", src, dst)
	}
}

func Cp(w http.ResponseWriter, r *http.Request) {
	src := r.URL.Query().Get("src")
	dst := r.URL.Query().Get("dst")
	_, err := cp(&src, &dst)
	if err == nil {
		fmt.Fprintf(w, "cp %s %s \n", src, dst)
	} else {
		fmt.Fprintf(w, "cp %s %s failed \n", src, dst)
	}
}

func GetContentSummary(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	summary, err := getContentSummary(&path)
	if err != nil {
		fmt.Fprintf(w, "couldn't get content summary of file %s \n", path)
	} else {
		fmt.Fprintf(w, summary)
	}
}

func CreateFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	created, _ := createEmptyFile(&path)
	if created {
		fmt.Fprintf(w, "created file %s \n", path)
	} else {
		fmt.Fprintf(w, "couldn't create file %s \n", path)
	}
}

func Mkdir(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	mask := os.FileMode(0777)
	created, _ := mkdir(&path, &mask)
	if created {
		fmt.Fprintf(w, "made directory %s \n", path)
	} else {
		fmt.Fprintf(w, "couldn't make directory %s \n", path)
	}
}

func ReadFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	_, err := readFile(w, &path)
	if err != nil {
		fmt.Fprintf(w, "couldn't read file %s \n", path)
	}
}

func Rm(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	recursive := false
	removed, err := rm(&path, &recursive)
	if err != nil || removed == false {
		fmt.Fprintf(w, "couldn't remove %s \n", path)
	} else {
		fmt.Fprintf(w, "removed %s \n", path)
	}
}

func RmAll(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	recursive := true
	removed, err := rm(&path, &recursive)
	if err != nil || removed == false {
		fmt.Fprintf(w, "couldn't remove %s \n", path)
	} else {
		fmt.Fprintf(w, "removed %s \n", path)
	}
}