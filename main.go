package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func getNSISubjectsLinks() Response {
	/*
		Function to get NSI subjects from the official API of the French Ministry of Education.
		Subjects are returned as a Response struct.
	*/

	req, err := http.NewRequest("GET", "https://cyclades.education.gouv.fr/delos/api/public/sujets/ece?sort=libelle&order=ASC&page=0&itemsPerPage=99&globalFilter=", nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response: ", err)
	}

	var responseData Response
	err = json.Unmarshal(responseBody, &responseData)
	if err != nil {
		log.Fatal("Error parsing json: ", err)
	}

	return responseData
}

func saveNSISubjects(subjectsAPIResponse Response) {
	/*
		Function to save NSI subjects in the export folder.
	*/

	var fileLen int8
	fileLen = 0

	for _, subject := range subjectsAPIResponse.Content {
		fileLen++

		for _, file := range subject.Fichiers {
			req, err := http.NewRequest("GET", "https://cyclades.education.gouv.fr/delos/api/file/public/"+file.ID, nil)
			if err != nil {
				log.Fatal(err)
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			responseBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal("Error reading response: ", err)
			}

			var fileExtension string
			// isPDF because NSI Subjects are two files types: .pdf and .txt (one for the subject and one for python code)
			// todo: Replace with ternary operator

			if !isPDF(responseBody) {
				fileExtension = ".txt"
			} else {
				fileExtension = ".pdf"
			}

			err = ioutil.WriteFile("./export/"+strconv.FormatInt(int64(fileLen), 10)+"_"+file.ID+fileExtension, responseBody, 0644)
			if err != nil {
				log.Fatal("Error writing file: ", err)
			}
		}
	}
}

func main() {
	/*
		Main function to download NSI subjects from the official website of the French Ministry of Education.
		Subjects are downloaded in the export folder.
	*/

	timeStart := time.Now()
	println("Starting download...")

	subjetsLinks := getNSISubjectsLinks()
	saveNSISubjects(subjetsLinks)

	println("Download finished in ", time.Since(timeStart).Milliseconds(), "ms")
}
