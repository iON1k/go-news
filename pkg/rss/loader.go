package rss

import (
	"io"
	"net/http"
	"os"
)

// Интерфейс загрузчика данных для RSS
type Loader interface {
	// Загрузка сырых данных для RSS
	LoadFeed() ([]byte, error)
}

// Загрузчик данных из HTTP
type HttpLoader struct {
	url string
}

// Конструктор загрузчика данных из HTTP
func NewHttpLoader(url string) *HttpLoader {
	return &HttpLoader{url: url}
}

// Загрузка сырых данных из HTTP
func (l HttpLoader) LoadFeed() ([]byte, error) {
	resp, err := http.Get(l.url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Загрузчик данных для RSS из файла
type FileLoader struct {
	file string
}

// Конструктор загрузчика данных из файла
func NewFileLoader(file string) *FileLoader {
	return &FileLoader{file: file}
}

// Загрузка сырых данных из файла
func (l FileLoader) LoadFeed() ([]byte, error) {
	bytes, err := os.ReadFile(l.file)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
