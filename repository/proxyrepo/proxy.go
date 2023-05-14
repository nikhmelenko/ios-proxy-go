package proxyrepo

import (
	"net/http"
	"net/url"
	"time"
)

type proxyRepo struct {
	client  *http.Client
	hostUrl string
}

func New(hostUrl string) *proxyRepo {
	return &proxyRepo{hostUrl: hostUrl, client: &http.Client{Timeout: 10 * time.Second}}
}

func (repo *proxyRepo) Get(path string) (resp *http.Response, err error) {
	redirectUrl, _ := url.JoinPath(repo.hostUrl, path)
	return repo.client.Get(redirectUrl)
}
