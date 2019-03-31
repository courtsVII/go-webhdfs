package hdfs

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Mv moves file from src to dst
// *http.Request requires src and dst parameters
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

// Cp copies file from src to dst
// *http.Request requires src and dst parameters
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

// GetContentSummary gets HDFS content summary for specified path
// *http.Request requires path parameter
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

// CreateFile creates empty file at specified path
// *http.Request requires path parameter
func CreateFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	created, _ := createEmptyFile(&path)
	if created {
		fmt.Fprintf(w, "created file %s \n", path)
	} else {
		fmt.Fprintf(w, "couldn't create file %s \n", path)
	}
}

// Mkdir creates directory at specified path
// *http.Request requires path parameter
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

// ReadFile writes contents of file to http.ResponseWriter
// *http.Request requires path parameter
func ReadFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	_, err := readFile(w, &path)
	if err != nil {
		fmt.Fprintf(w, "couldn't read file %s \n", path)
		log.Println(err)
	}
}

// WriteFile writes request FormFile to a file in HDFS
// *http.Request requires path parameter and a FormFile parameter named file
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

// Write writes request body to a file in HDFS
// *http.Request requires path parameter
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

// Rm removes file at specified path from HDFS
// *http.Request requires path parameter
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

// Ls lists files at specified HDFS path
// *http.Request requires path parameter. Allows for an optional recursive parameter to be set to true. This will make the ls recursive
func Ls(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	recursive := r.URL.Query().Get("recursive")

	var response = []string{}
	var err error
	if recursive == "true" {
		response, err = ls(&path)
	} else {
		var fullList = []string{}
		fullList, err = ls(&path)
		regex := regexp.MustCompile("/")
		matches := regex.FindAllStringIndex(path, -1)
		slashes := len(matches)

		for _, entry := range fullList {
			matches = regex.FindAllStringIndex(entry, -1)
			entrySlashes := len(matches)
			if slashes == entrySlashes {
				response = append(response, entry)
			}
		}
	}

	if err != nil {
		fmt.Fprintf(w, "couldn't ls %s \n", path)
		log.Println(err)
	} else {
		fmt.Fprintf(w, "%s\n", strings.Join(response, "\n"))
	}
}

// RmAll removes filea at specified path from HDFS recursively
// *http.Request requires path parameter
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

// Chown changes ownership of file/directory in HDFS
// *http.Request requires path, user and group parameters. Path is the path to file/folder to chown. User is the user to set file ownership as. Group is the group to set file ownership as.
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

// Chmod chmod file at specified path from HDFS
// *http.Request requires path parameter and mask parameter. Mask is specified in UNIX numeric permisson notation
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
