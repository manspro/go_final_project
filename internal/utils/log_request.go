package utils

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

func LogRequestBody(r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Ошибка при чтении тела запроса:", err)
		return
	}

	log.Printf("Тело запроса: %s", string(bodyBytes))

	// Восстанавливаем io.ReadCloser для последующего использования
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
}
