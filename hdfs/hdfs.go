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
		fmt.Printf("hadoop address not provided, using default 0.0.0.0:9000\n")
		hadoopAddress = "0.0.0.0:9000"
	}
	var err error
	hadoopClient, err = hdfs.New(hadoopAddress)
	if err != nil {
		fmt.Printf("couldn't connect to hadoop on %s\n", hadoopAddress)
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Printf("connected to hadoop on %s as user %s\n", hadoopAddress, hadoopClient.User())
	}
}

// gets the 'src' and 'dst' params from http.Request and moves src to dst on hdfs
func Mv(w http.ResponseWriter, r *http.Request) {
	src := r.URL.Query().Get("src")
	dst := r.URL.Query().Get("dst")
	err := hadoopClient.Rename(src, dst)
	if err != nil {
		fmt.Fprintf(w, "mv %s %s failed\n", src, dst)
		log.Println(err)
	} else {
		fmt.Fprintf(w, "mv %s %s\n", src, dst)
		log.Printf("mv %s %s\n", src, dst)
	}
}
