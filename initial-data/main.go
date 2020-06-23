package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
)

func main() {
	c := new(http.Client)

	estateDirs := []string{
		"./request_content/estate_search",
		"./request_content/recommend_estate",
		"./request_content/recommend_estate_with_chair",
	}
	chairDirs := []string{
		"./request_content/chair_search",
		"./request_content/recommend_chair",
	}
	nazotteDir := "./request_content/estate_nazotte"

	for _, dir := range estateDirs {
		estateFilePaths, err := ioutil.ReadDir(dir)
		if err != nil {
			fmt.Errorf("%v", err)
		}
		for _, fp := range estateFilePaths {
			SaveEstateResponseFile(fp.Name(), c, dir)
		}
	}

	for _, dir := range chairDirs {
		chairFilePaths, err := ioutil.ReadDir(dir)
		if err != nil {
			fmt.Errorf("%v", err)
		}
		for _, fp := range chairFilePaths {
			SaveChairResponseFile(fp.Name(), c, dir)
		}
	}
	nazotteFilePaths, err := ioutil.ReadDir(nazotteDir)
	if err != nil {
		fmt.Errorf("%v", err)
	}
	for _, fp := range nazotteFilePaths {
		SaveNazotteResponseFile(fp.Name(), c, nazotteDir)
	}
}
func GetQueryParams(request RequestFrame) url.Values {
	// 読み込んだ構造体を元にリクエストパラメーターを作成する
	queryParam := url.Values{}
	val := reflect.ValueOf(&(request.RequectContent.Query)).Elem()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		if fmt.Sprintf("%v", valueField.Interface()) == "" {
		} else {
			queryParam.Set(tag.Get("json"), fmt.Sprintf("%v", valueField.Interface()))
		}
	}
	return queryParam
}
func RequestResponseWithFilePath(filepath string, c *http.Client) ([]byte, Request) {
	raw, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Errorf("Error: on opening filePath [%v] because of %v", raw, err)
	}
	var request RequestFrame
	json.Unmarshal(raw, &request)
	path := request.RequectContent.Uri
	method := request.RequectContent.Method
	u := "http://localhost:1323" + path
	if request.RequectContent.Id != "" {
		u = u + "/" + request.RequectContent.Id
	}
	q := GetQueryParams(request)
	req, err := http.NewRequest(method, u, nil)
	if err != nil {
		fmt.Errorf("Error: on making httpRequest [%v %v] because of %v", u, q, err)
	}
	req.URL.RawQuery = q.Encode()
	resp, err := c.Do(req)
	if err != nil {
		fmt.Errorf("Error: on doing httpRequest [%v %v] because of %v", u, q, err)
		panic("")
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return body, request.RequectContent
}

func SaveEstateResponseFile(filePath string, c *http.Client, srcDir string) error {
	fmt.Printf("Start make %s \n", filePath)
	ef := EstatesAnswerJson{}
	body, req := RequestResponseWithFilePath(srcDir+"/"+filePath, c)
	ef.Req = req
	_ = json.Unmarshal(body, &ef.Res.Body)
	bytes, _ := json.Marshal(ef)
	_ = ioutil.WriteFile("./generate_verification/"+filePath, bytes, os.FileMode(0777))
	return nil
}
func SaveChairResponseFile(filePath string, c *http.Client, srcDir string) error {
	fmt.Printf("Start make %s \n", filePath)
	cf := ChairsAnswerJson{}
	body, req := RequestResponseWithFilePath(srcDir+"/"+filePath, c)
	cf.Req = req
	_ = json.Unmarshal(body, &cf.Res.Body)
	bytes, _ := json.Marshal(cf)
	_ = ioutil.WriteFile("./generate_verification/"+filePath, bytes, os.FileMode(0777))
	return nil
}
func SaveNazotteResponseFile(filePath string, c *http.Client, srcDir string) error {
	fmt.Printf("Start make %s \n", filePath)
	ef := EstatesAnswerJson{}
	body, req := RequestNazotteFilePath(srcDir+"/"+filePath, c)
	ef.Req = req
	_ = json.Unmarshal(body, &ef.Res.Body)
	bytes, _ := json.Marshal(ef)
	_ = ioutil.WriteFile("./generate_verification/"+filePath, bytes, os.FileMode(0777))
	return nil
}
func RequestNazotteFilePath(filepath string, c *http.Client) ([]byte, Request) {
	raw, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Errorf("Error: on opening filePath [%v] because of %v", raw, err)
	}
	var request RequestFrame
	json.Unmarshal(raw, &request)
	path := request.RequectContent.Uri
	method := request.RequectContent.Method
	u := "http://localhost:1323" + path
	rb, _ := json.Marshal(request.RequectContent.Body)
	req, err := http.NewRequest(method, u, bytes.NewBuffer(rb))
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		fmt.Errorf("Error: on doing httpRequest [%v %v] because of %v", u, rb, err)
		panic("")
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return body, request.RequectContent
}
