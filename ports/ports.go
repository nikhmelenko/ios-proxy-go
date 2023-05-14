package ports

import (
	"net/http"
	"web-app-test/entity"
)

type DBRepository interface {
	Insert(entity.RequestObj)
}

type ProxyRepository interface {
	Get(string) (*http.Response, error)
}

type ProxyService interface {
	GetModified(string) (*http.Response, error)
	Save(*http.Response, error)
}
