package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"web-app-test/entity"
	"web-app-test/repository/mongorepo"
	"web-app-test/repository/proxyrepo"
	"web-app-test/service"
)

func IntMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestIntMinBasic(t *testing.T) {
	ans := IntMin(2, -2)
	if ans != -2 {
		t.Errorf("IntMin(2, -2) = %d; want -2", ans)
	}
}

func TestProxyRepoGetOk(t *testing.T) {
	proxyRepository := proxyrepo.New("https://httpstat.us/")
	resp, err := proxyRepository.Get("/200")
	if err != nil {
		t.Error(err)
	}
	if resp.Request.URL.Host != "httpstat.us" {
		t.Error("Invalid host: ", resp.Request.URL.Host, " expected: ", "httpstat.us")
	}
	if resp.StatusCode != http.StatusOK {
		t.Error("Invalid resonse status: ", resp.StatusCode, " expected: ", http.StatusOK)
	}
}

func TestProxyRepoGetNotFound(t *testing.T) {
	proxyRepository := proxyrepo.New("https://httpstat.us/")
	resp, err := proxyRepository.Get("/unexisting")
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Error("Invalid resonse status: ", resp.StatusCode, " expected: ", http.StatusNotFound)
	}
}

func TestMongoRepositoryInvalidUrl(t *testing.T) {
	mongoDbUrl := "https://httpstat.us/"
	defer func() {
		if r := recover(); r == nil {
			t.Error("The code did not panic, thus provided db URL is invalid")
		}
	}()
	mongoRepository := mongorepo.NewDB(mongoDbUrl)
	defer mongoRepository.Client.Disconnect(context.TODO())
}

type MockDBrepo struct {
	Calls int
}

func (repo *MockDBrepo) Insert(entity.RequestObj) {
	repo.Calls++
}

func TestProxyServiceGetModifiedEmbededOk(t *testing.T) {
	proxyHost := "https://jsonplaceholder.typicode.com/"
	mockedDbRepo := &MockDBrepo{}

	proxyRepository := proxyrepo.New(proxyHost)
	proxyService := service.New(proxyRepository, mockedDbRepo)
	resp, err := proxyService.GetModified("/todos/1")
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Error("Invalid resonse status: ", resp.StatusCode, " expected: ", http.StatusOK)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var results entity.EmbededJson
	err = json.Unmarshal(body, &results)
	_, exists := results["custom_key"]
	if !exists {
		t.Error("Custom key-value was not attached to response")
	}
	if mockedDbRepo.Calls == 0 {
		t.Error("Proxy service didn't trigger data saving")
	}
}

func TestProxyServiceGetModifiedListOk(t *testing.T) {
	proxyHost := "https://jsonplaceholder.typicode.com/"
	mockedDbRepo := &MockDBrepo{}

	proxyRepository := proxyrepo.New(proxyHost)
	proxyService := service.New(proxyRepository, mockedDbRepo)
	resp, err := proxyService.GetModified("/todos")
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Error("Invalid resonse status: ", resp.StatusCode, " expected: ", http.StatusOK)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var results entity.ListJson
	err = json.Unmarshal(body, &results)
	_, exists := results[0]["custom_key"]
	if !exists {
		t.Error("Custom key-value was not attached to response")
	}
	if mockedDbRepo.Calls == 0 {
		t.Error("Proxy service didn't trigger data saving")
	}
}
