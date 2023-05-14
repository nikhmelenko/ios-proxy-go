package service

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"web-app-test/entity"
	"web-app-test/ports"
)

type service struct {
	dbRepository    ports.DBRepository
	proxyRepository ports.ProxyRepository
}

func New(proxyRepository ports.ProxyRepository, dbRepository ports.DBRepository) *service {
	return &service{
		dbRepository:    dbRepository,
		proxyRepository: proxyRepository,
	}
}

func (srv *service) GetModified(path string) (*http.Response, error) {
	resp, err := srv.proxyRepository.Get(path)
	if err != nil {
		return nil, err
	}
	modifiedBody, err := srv.modifyJson(resp)
	resp.Body = ioutil.NopCloser(bytes.NewReader(modifiedBody))
	resp.ContentLength = int64(len(modifiedBody))
	srv.Save(resp, err)
	return resp, err
}

func (srv *service) Save(resp *http.Response, err error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	// Replace the body with a new reader after reading from the original
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	obj := entity.RequestObj{
		Time:         time.Now(),
		Path:         resp.Request.URL.Path,
		Method:       resp.Request.Method,
		Status:       resp.StatusCode,
		ResponseBody: string(body),
	}
	srv.dbRepository.Insert(obj)
}

func (srv *service) modifyJson(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	var results entity.EmbededJson
	err = json.Unmarshal(body, &results)
	// error happens if JSON is not Embeded
	if err != nil {
		var listResults entity.ListJson
		json.Unmarshal(body, &listResults)
		listResults.AppendCustomKey()
		return json.Marshal(listResults)
	}
	results.AppendCustomKey()
	return json.Marshal(results)
}
