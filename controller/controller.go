package main

import (
	"fmt"
	"log"
	"net/http"
	"webhdfs/hdfs"
)

func main() {
	http.HandleFunc("/v1/hdfs/mv", hdfs.Mv)
	http.HandleFunc("/v1/up/", up)
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", nil))
}

func up(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Service Up! \n")
}
