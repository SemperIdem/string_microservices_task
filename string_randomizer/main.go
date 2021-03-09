package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var (
	minLen, maxLen int
	encUrl         = os.Getenv("ENCRYPTOR_URL")
)

func generateRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func randStringsList(n int) []string {
	list := make([]string, n)
	for i := 0; i < n; i++ {
		list[i] = generateRandomString(rand.Intn(maxLen-minLen) + minLen)
	}
	return list
}

func randomEncryptString(w http.ResponseWriter, r *http.Request) {
	param, ok := r.URL.Query()["size"]
	if !ok || len(param[0]) < 1 {
		fmt.Fprintf(w, "Url Param 'size' is missing")
		return
	}
	numberStrings, err := strconv.Atoi(param[0])
	if numberStrings < 1 {
		fmt.Fprintf(w, "Url Param 'size' must be positive number")
		return
	}
	randomStringList := randStringsList(numberStrings)
	jsonStringList, _ := json.Marshal(randomStringList)
	resp, err := http.Post(encUrl, "application/json", bytes.NewBuffer(jsonStringList))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Generated and sented %d random strings", numberStrings)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	list := make([]string, numberStrings)
	err = json.Unmarshal(body, &list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Random hash from the recivied list: %s ", list[rand.Intn(numberStrings)])
}

func main() {
	var errConv error
	rand.Seed(time.Now().UnixNano())
	defaultPort := os.Getenv("DEFAULT_PORT")
	minLenStr := os.Getenv("MIN_STRING_LENGTH")
	minLen, errConv = strconv.Atoi(minLenStr)
	if errConv != nil || minLen < 0 {
		minLen = 3
	}

	maxLenStr := os.Getenv("MAX_STRING_LENGTH")
	maxLen, errConv = strconv.Atoi(maxLenStr)
	if errConv != nil || minLen < 0 {
		maxLen = 15
	}

	fmt.Printf("Randomizer service service is listening on port %s.\n", defaultPort)
	http.HandleFunc("/getRandomHash", randomEncryptString)
	err := http.ListenAndServe(":"+defaultPort, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
