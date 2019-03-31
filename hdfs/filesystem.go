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
		fmt.Printf("HADOOP_ADDRESS not provided, using default 0.0.0.0:9000 \n")
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
	}
	log.Printf("mv %s %s \n", *src, *dst)
	return true, nil
}

func rm(path *string, recursive *bool) (bool, error) {
	var err error

	if *recursive {
		err = hadoopClient.RemoveAll(*path)
	} else {
		err = hadoopClient.Remove(*path)
	}

	if err != nil {
		log.Println(err)
		return false, err
	}
	log.Printf("rm %s \n", *path)
	return true, nil
}

func cp(src *string, dst *string) (int64, error) {
	r, err := hadoopClient.Open(*src)
	defer r.Close()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	w, err := hadoopClient.Create(*dst)
	defer w.Close()
	if err != nil {
		log.Println(err)
		return 0, err
	}

	bytesCopied, err := io.Copy(w, r)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	log.Printf("cp %s %s \n", *src, *dst)
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
	}
	log.Printf("created file %s \n", *path)
	return true, nil
}

func ls(startPath *string) ([]string, error) {
	var response = []string{}
	first := true
	err := hadoopClient.Walk(*startPath, func(path string, _ os.FileInfo, _ error) error {
		if first {
			//don't put startPath in response
			first = false
			return nil
		}
		response = append(response, path)
		return nil
	})

	if err != nil {
		log.Println(err)
		return response, err
	}
	log.Printf("ls %s \n", *startPath)
	return response, nil
}

func getContentSummary(path *string) (string, error) {
	s, err := hadoopClient.GetContentSummary(*path)
	if err != nil {
		log.Println(err)
		return "", err
	}
	log.Printf("got content summary for %s \n", *path)
	return fmt.Sprintf("%s: \nsize %d \nsize after replication: %d \nspace quota: %d \ndirectory count %d \nfile count %d \nname quota %d", *path, s.Size(), s.SizeAfterReplication(), s.SpaceQuota(), s.DirectoryCount(), s.FileCount(), s.NameQuota()), nil
}

func chmod(path *string, mask *os.FileMode) (bool, error) {
	err := hadoopClient.Chmod(*path, *mask)
	if err != nil {
		log.Println(err)
		return false, err
	}
	log.Printf("chmod applied to %s \n", *path)
	return true, nil
}

func chown(path *string, user *string, group *string) error {
	err := hadoopClient.Chown(*path, *user, *group)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("chown user %s %s group applied to %s \n", *user, *group, *path)
	return nil
}

func readFile(w io.Writer, path *string) (int64, error) {
	r, err := hadoopClient.Open(*path)
	defer r.Close()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	bytesCopied, err := io.Copy(w, r)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	log.Printf("reading from %s \n", *path)
	return bytesCopied, nil
}

func write(r io.ReadCloser, path *string) (int64, error) {
	w, err := hadoopClient.Create(*path)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	bytesCopied, writeErr := io.Copy(w, r)
	flushErr := w.Flush()
	closeErr := w.Close()
	if flushErr != nil {
		log.Println(writeErr)
	}
	if writeErr != nil {
		log.Println(writeErr)
	}
	if closeErr != nil {
		log.Println(closeErr)
	}
	if writeErr != nil || closeErr != nil {
		return bytesCopied, err
	}
	log.Printf("wrote to %s \n", *path)
	return bytesCopied, nil
}
