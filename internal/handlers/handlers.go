package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

func ReturnHTML(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("C:/dev/xails_server_yandex/index.html")
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
	const maxMemory = 32 << 20 // 32MB

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Парсинг формы
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		http.Error(w, "Файл слишком большой", http.StatusBadRequest)
		return
	}

	// Получение файла
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Файл не найден в запросе", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Чтение файла
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Ошибка чтения файла", http.StatusInternalServerError)
		return
	}
	content := string(data)
	if content == "" {
		http.Error(w, "Файл пуст", http.StatusBadRequest)
		return
	}

	// Конвертация
	converted, err := service.AutoConvert(content)
	if err != nil {
		http.Error(w, "Ошибка конвертации: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Сохранение файла
	timestamp := strings.ReplaceAll(time.Now().UTC().Format("2006-01-02_15-04-05"), ":", "-")
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

	if _, err := outputFile.WriteString(converted); err != nil {
		http.Error(w, "Ошибка записи в файл", http.StatusInternalServerError)
		return
	}

	// Отправка ответа
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(converted))
}
