package hdfs

import (
	"fmt"
	"io"
	"log"
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

func mv(src *string, dst *string) (bool, error) {
	err := hadoopClient.Rename(*src, *dst)
	if err != nil {
		log.Println(err)
		return false, err
	} else {
		log.Printf("mv %s %s \n", *src, *dst)
		return true, nil
	}
}

func rm(path *string, recursive *bool) (bool, error) {
	var err error = nil
	if *recursive {
		err = hadoopClient.RemoveAll(*path)
	} else {
		err = hadoopClient.Remove(*path)
	}

	if err != nil {
		log.Println(err)
		return false, err
	} else {
		log.Printf("rm %s \n", *path)
		return true, nil
	}
}

func cp(src *string, dst *string) (int64, error) {
	r, err := hadoopClient.Open(*src)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	w, err := hadoopClient.Create(*dst)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	bytesCopied, err := io.Copy(w, r)
	if err != nil {
		log.Println(err)
		return 0, err
	} else {
		log.Printf("cp %s %s \n", *src, *dst)
	}
	err = r.Close()
	if err != nil {
		log.Println(err)
	}

	err = w.Close()
	if err != nil {
		log.Println(err)
	}

	//not a lot we can do if closes fail
	return bytesCopied, nil
}

func mkdir(path *string, perm *os.FileMode) (bool, error) {
	err := hadoopClient.MkdirAll(*path, *perm)
	if err != nil {
		log.Println(err)
		return false, err
	}
	log.Printf("mkdir %s \n", *path)
	return true, nil
}

func createEmptyFile(path *string) (bool, error) {
	err := hadoopClient.CreateEmptyFile(*path)
	if err != nil {
		log.Println(err)
		return false, err
	} else {
		log.Printf("created file %s \n", *path)
		return true, nil
	}
}

func getContentSummary(path *string) (string, error) {
	s, err := hadoopClient.GetContentSummary(*path)
	if err != nil {
		log.Println(err)
		return "", err
	} else {
		log.Printf("got content summary for %s \n", *path)
		return fmt.Sprintf("%s: \nsize %d \nsize after replication: %d \nspace quota: %d \ndirectory count %d \nfile count %d \nname quota %d", *path, s.Size(), s.SizeAfterReplication(), s.SpaceQuota(), s.DirectoryCount(), s.FileCount(), s.NameQuota()), nil
	}
}

func chmod(path *string, mask *os.FileMode) (bool, error) {
	err := hadoopClient.Chmod(*path, *mask)
	if err != nil {
		log.Println(err)
		return false, err
	} else {
		log.Printf("chmod applied to %s \n", *path)
		return true, nil
	}
}

func readFile(w io.Writer, path *string) (int64, error) {
	r, err := hadoopClient.Open(*path)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	bytesCopied, err := io.Copy(w, r)
	if err != nil {
		log.Println(err)
		return 0, err
	} else {
		log.Printf("reading from %s \n", *path)
	}
	err = r.Close()
	if err != nil {
		log.Println(err)
	}
	return bytesCopied, nil
}
