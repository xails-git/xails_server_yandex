package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

func ReturnHTML(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("index.html")
	if err != nil {
		http.Error(w, "ошибка чтения файла index.html", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(data)
	if err != nil {
		log.Printf("Ошибка при отправке данных: %v", err)
	}
}

func ConvertStr(w http.ResponseWriter, r *http.Request) {
	const maxMemory = 32 << 20

	err := r.ParseMultipartForm(maxMemory)
	if err != nil {
		http.Error(w, "ошибка парсинга", http.StatusInternalServerError)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "ошибка файла", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// обработка файла

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "ошибка чтения файла", http.StatusInternalServerError)
		return
	}

	content := string(data)
	if content == "" {
		http.Error(w, "файл пустой", http.StatusBadRequest)
		return
	}

	converted, err := service.AutoConvert(content)
	if err != nil {
		http.Error(w, "ошибка конвертации", http.StatusInternalServerError)
		return
	}

	// создаем файл для результата

	timestamp := time.Now().UTC().Format("2006-01-02_15-04-05")
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".txt"
	}
	outputFilename := "result_" + timestamp + ext

	outputFile, err := os.Create(outputFilename)
	if err != nil {
		http.Error(w, "Ошибка создания файла", http.StatusInternalServerError)
		return
	}
	defer outputFile.Close()

	_, err = outputFile.WriteString(converted)
	if err != nil {
		http.Error(w, "ошибка записи в файл", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(converted))
	if err != nil {
		log.Printf("ошибка при отправке ответа %v", err)
	}
}
