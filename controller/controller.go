package main

import (
	"fmt"
	"log"
	"net/http"
	"webhdfs/hdfs"
)

func main() {
	http.HandleFunc("/v1/hdfs/mv", hdfs.Mv)
	http.HandleFunc("/v1/hdfs/cp", hdfs.Cp)
	http.HandleFunc("/v1/hdfs/createfile", hdfs.CreateFile)
	http.HandleFunc("/v1/hdfs/readfile", hdfs.ReadFile)
	http.HandleFunc("/v1/hdfs/getcontentsummary", hdfs.GetContentSummary)
	http.HandleFunc("/v1/hdfs/mkdir", hdfs.Mkdir)
	http.HandleFunc("/v1/hdfs/rm", hdfs.Rm)
	http.HandleFunc("/v1/hdfs/rmall", hdfs.RmAll)
	http.HandleFunc("/v1/hdfs/chmod", hdfs.Chmod)
	http.HandleFunc("/v1/hdfs/chown", hdfs.Chown)
	http.HandleFunc("/v1/hdfs/writefile", hdfs.WriteFile)

	http.HandleFunc("/v1/up/", up)
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", nil))
}

func up(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Service Up! \n")
}
