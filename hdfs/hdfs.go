package hdfs

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/colinmarc/hdfs"
)

var hadoopClient *hdfs.Client

func init() {
	hadoopAddress := os.Getenv("HADOOP_ADDRESS")
	if len(hadoopAddress) == 0 {
		fmt.Printf("hadoop address not provided, using default 0.0.0.0:9000 \n")
		hadoopAddress = "0.0.0.0:9000"
	}
	var err error
	hadoopClient, err = hdfs.New(hadoopAddress)
	if err != nil {
		fmt.Printf("couldn't connect to hadoop on %s \n", hadoopAddress)
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Printf("connected to hadoop on %s as user %s \n", hadoopAddress, hadoopClient.User())
	}
}

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

func mv(src string, dst string) (bool, error) {
	err := hadoopClient.Rename(src, dst)
	if err != nil {
		log.Println(err)
		return false, err
	} else {
		log.Printf("mv %s %s \n", src, dst)
		return true, nil
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

func createEmptyFile(path string) (bool, error) {
	err := hadoopClient.CreateEmptyFile(path)
	if err != nil {
		log.Println(err)
		return false, err
	} else {
		log.Printf("created file %s \n", path)
		return true, nil
	}
}

func GetContentSummary(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	summary, err := getContentSummary(path)
	if err != nil {
		fmt.Fprintf(w, "couldn't get content summary of file %s \n", path)
	} else {
		fmt.Fprintf(w, summary)
		log.Printf("got content summary for %s \n", path)
	}
}

func getContentSummary(path string) (string, error) {
	s, err := hadoopClient.GetContentSummary(path)
	if err != nil {
		log.Println(err)
		return "", err
	} else {
		log.Printf("got content summary for %s \n", path)
		return fmt.Sprintf("%s: \nsize %d \nsize after replication: %d \nspace quota: %d \ndirectory count %d \nfile count %d \nname quota %d", path, s.Size(), s.SizeAfterReplication(), s.SpaceQuota(), s.DirectoryCount(), s.FileCount(), s.NameQuota()), nil
	}
}
