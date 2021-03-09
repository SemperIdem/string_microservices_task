package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var (
	worksNum int
	jobs     chan string
	results  chan string
)

func worker(id int, f func(string) string) {
	for j := range jobs {
		fmt.Printf("Work #%d processed %s -> %s \n", id, j, f(j))
		results <- f(j)
	}
}

func toSha256Hash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	sha256Hash := hex.EncodeToString(h.Sum(nil))
	return sha256Hash
}

func toHashList(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer req.Body.Close()
	var listStrings []string
	err = json.Unmarshal(body, &listStrings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	results = make(chan string, len(listStrings))
	hashList := make([]string, len(listStrings))

	for i := 0; i < len(listStrings); i++ {
		jobs <- listStrings[i]
	}
	for i := 0; i < len(listStrings); i++ {
		hashList[i] = <-results
	}
	close(results)

	resp, _ := json.Marshal(hashList)
	_, err = w.Write(resp)
	if err != nil {
		fmt.Printf("Write failed: %v", err)
	}

}

func main() {
	var errConv error
	defaultPort := os.Getenv("DEFAULT_PORT")
	worksNumStr := os.Getenv("WORKERS_NUMBER")
	worksNum, errConv = strconv.Atoi(worksNumStr)
	if errConv != nil || worksNum < 0 {
		worksNum = 4
	}

	// init jobs and workers for Work Pool
	jobs = make(chan string)
	for w := 1; w <= worksNum; w++ {
		go worker(w, toSha256Hash)
	}

	http.HandleFunc("/", toHashList)
	fmt.Printf("Encryptor service is listening on port %s.\n", defaultPort)
	err := http.ListenAndServe(":"+defaultPort, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
