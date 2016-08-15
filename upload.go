package gifs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

const (
	baseURL  = "https://api.gifs.com"
	endpoint = "/media/upload"
)

func Upload(name string, r io.Reader) (*wrapperResponse, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormFile("file", name)
	if err != nil {
		return &wrapperResponse{}, err
	}
	if _, err = io.Copy(fw, r); err != nil {
		return &wrapperResponse{}, err
	}

	w.Close()

	req, err := http.NewRequest("POST", baseURL+endpoint, &b)
	if err != nil {
		return &wrapperResponse{}, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &wrapperResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("returned status: %s", res.Status)
		return &wrapperResponse{}, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &wrapperResponse{}, err
	}

	uRes := new(wrapperResponse)
	err = json.Unmarshal(body, uRes)
	if err != nil {
		return &wrapperResponse{}, err
	}

	return uRes, nil
}
