package src

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

// region " JSON structs "
type FileIO_Response struct {
	Success bool   `json:"success"`
	Key     string `json:"key"`
	Link    string `json:"link"`
	Expiry  string `json:"expiry"`
}
// endregion

func UploadFile(fileLocation string, expiration string) (string, error) {
	data, err := os.Open(fileLocation)
	if err != nil {
		return "", err
	}
	defer data.Close()

	if expiration == "" {
		expiration = "1w"
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("https://file.io?expires=%s", expiration), data)
	if err != nil {
		return "", err
	}

	contentType, err := GetFileContentType(data)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	jsonResp := FileIO_Response{}
	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		return "", err
	}

	return jsonResp.Link, nil
}

func DownloadFile(downloadLink string) ([]byte, error) {
	fileName := path.Base(downloadLink)
	resp, err := http.Get(downloadLink)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	f, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	return body, nil
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