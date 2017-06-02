package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"libraries/constant"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func PostRequest(urlstring string, params map[string]string) (output []byte, err error) {

	urlValues := url.Values{}

	for k, v := range params {
		urlValues.Add(k, v)
	}

	//gzipbyte, _ := GzipEncode([]byte(urlValues.Encode()))
	request, err := http.NewRequest("POST", urlstring, strings.NewReader(urlValues.Encode()))
	if err != nil {
		return nil, NewError(constant.ERROR_UTILS_POST, err.Error())
	}
	//request.Header.Add("Accept-Encoding", "gzip, deflate")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded; param=value")
	request.Header.Add("Connection", "keep-alive")
	//request.Header.Add("Connection", "application/x-gzip")

	//proxyUrl, _ := url.Parse("http://127.0.0.1:8888/")
	client := http.Client{
	//Timeout: 60 * time.Second,
	//Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
	}
	resp, err := client.Do(request)
	if resp != nil {
		defer resp.Body.Close()
	}

	//resp, err := client.PostForm(urlstring, urlValues)
	if err != nil {
		return nil, NewError(constant.ERROR_UTILS_POST, err.Error())
	}

	output, err2 := ioutil.ReadAll(resp.Body)

	if err2 != nil {
		err = NewError(constant.ERROR_UTILS_POST_RESPONSE_READ, err2.Error())
	}

	return
}

func WebPostRequest(urlstring string, planids interface{}) (output []byte, err error) {

	c, err := json.Marshal(planids)
	body := bytes.NewBuffer([]byte(c))

	//res, err := http.Post(urlstring, "application/json;charset=utf-8", body)
	//gzipbyte, _ := GzipEncode([]byte(urlValues.Encode()))
	request, err := http.NewRequest("POST", urlstring, body)
	if err != nil {
		return nil, NewError(constant.ERROR_UTILS_POST, err.Error())
	}
	//request.Header.Add("Accept-Encoding", "gzip, deflate")
	request.Header.Add("Content-Type", "application/json;charset=utf-8; param=value")
	request.Header.Add("Connection", "keep-alive")
	//request.Header.Add("Connection", "application/x-gzip")

	//proxyUrl, _ := url.Parse("http://127.0.0.1:8888/")
	client := http.Client{
		Timeout: 30 * time.Second,
		//Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
	}
	resp, err := client.Do(request)
	if resp != nil {
		defer resp.Body.Close()
	}

	//resp, err := client.PostForm(urlstring, urlValues)
	if err != nil {
		return nil, NewError(constant.ERROR_UTILS_POST, err.Error())
	}

	output, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		err = NewError(constant.ERROR_UTILS_POST_RESPONSE_READ, err2.Error())
	}

	return
}

func GzipEncode(in []byte) ([]byte, error) {
	var (
		buffer bytes.Buffer
		out    []byte
		err    error
	)
	writer := gzip.NewWriter(&buffer)
	_, err = writer.Write(in)
	if err != nil {
		writer.Close()
		return out, err
	}
	err = writer.Close()
	if err != nil {
		return out, err
	}

	return buffer.Bytes(), nil
}

func GzipDecode(in []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		var out []byte
		return out, err
	}
	defer reader.Close()

	return ioutil.ReadAll(reader)
}
