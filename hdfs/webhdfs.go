package hdfs

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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
		log.Println(err)
	} else {
		fmt.Fprintf(w, "cp %s %s failed \n", src, dst)
	}
}

func GetContentSummary(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	summary, err := getContentSummary(&path)
	if err != nil {
		fmt.Fprintf(w, "couldn't get content summary of file %s \n", path)
		log.Println(err)
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
		log.Println(err)
	}
}

func WriteFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	rc, _, err := r.FormFile("file")
	defer rc.Close()
	if err != nil {
		fmt.Fprintf(w, "couldn't read file from request \n")
		log.Println(err)
		return
	}

	_, err = write(rc, &path)
	if err != nil {
		fmt.Fprintf(w, "couldn't write file %s \n", path)
		log.Println(err)
		return
	}
}

func Write(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	rc := r.Body
	defer rc.Close()

	_, err := write(rc, &path)
	if err != nil {
		fmt.Fprintf(w, "couldn't write file %s \n", path)
		log.Println(err)
		return
	}
}

func Rm(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	recursive := false
	removed, err := rm(&path, &recursive)
	if err != nil {
		fmt.Fprintf(w, "couldn't remove %s \n", path)
		log.Println(err)
	} else if removed == false {
		fmt.Fprintf(w, "couldn't remove %s \n", path)
	} else {
		fmt.Fprintf(w, "removed %s \n", path)
	}
}

func Ls(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	response, err := ls(&path)
	if err != nil {
		fmt.Fprintf(w, "couldn't ls %s \n", path)
		log.Println(err)
	} else {
		fmt.Fprintf(w, "%s", response)
	}
}

func RmAll(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	recursive := true
	removed, err := rm(&path, &recursive)
	if err != nil {
		fmt.Fprintf(w, "couldn't remove %s \n", path)
		log.Println(err)
	} else if removed == false {
		fmt.Fprintf(w, "couldn't remove %s \n", path)
	} else {
		fmt.Fprintf(w, "removed %s \n", path)
	}
}

func Chown(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	user := r.URL.Query().Get("user")
	group := r.URL.Query().Get("group")

	err := chown(&path, &user, &group)
	if err != nil {
		fmt.Fprintf(w, "couldn't chown user %s group %s %s \n", user, group, path)
		log.Println(err)
	} else {
		fmt.Fprintf(w, "chown user %s group %s %s \n", user, group, path)
	}
}

func Chmod(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	perm := r.URL.Query().Get("mask")
	applied := false
	if len(perm) > 0 {
		mask, err := strconv.Atoi(perm)
		log.Println(err)
		fmt.Fprintf(w, "couldn't parse mask %s \n", perm)
		if err != nil {
			return
		}
		_mask := os.FileMode(mask)
		applied, _ = chmod(&path, &_mask)
	} else {
		fmt.Fprintf(w, "no mask provided \n")
		return
	}
	if applied {
		fmt.Fprintf(w, "chmod %s %s \n", perm, path)
	} else {
		fmt.Fprintf(w, "couldn't chmod %s \n", path)
	}
}
