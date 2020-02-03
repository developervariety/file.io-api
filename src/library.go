package src

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// region " JSON structs "
type FileIO_Response struct {
	Success bool   `json:"success"`
	Key     string `json:"key"`
	Link    string `json:"link"`
	Expiry  string `json:"expiry"`
}
// endregion

func UploadFile(fileLocation string) FileIO_Response {
	data, err := os.Open(fileLocation)
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	req, err := http.NewRequest("PUT", "https://file.io", data)
	if err != nil {
		log.Fatal(err)
	}

	contentType, err := GetFileContentType(data)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	jsonResp := FileIO_Response{}
	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		log.Fatal(err)
	}

	return jsonResp
}

func DownloadFile(downloadLink string) {
	// download file with provided file.io link
}

// https://golangcode.com/get-the-content-type-of-file/
func GetFileContentType(out *os.File) (string, error) {
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)
	return contentType, nil
}