// Пакет handlers
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
	// Получаем текущую рабочую директорию
	wd, err := os.Getwd()
	if err != nil {
		log.Printf("Ошибка получения рабочей директории: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Формируем абсолютный путь к index.html
	indexPath := filepath.Join(wd, "index.html")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		log.Printf("Ошибка чтения %s: %v", indexPath, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}

func ConvertStr(w http.ResponseWriter, r *http.Request) {
	const maxMemory = 32 << 20 // 32MB

	// Принимаем только POST-запросы
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Парсинг формы с увеличенным лимитом памяти
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		log.Printf("Ошибка парсинга формы: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Получаем файл из поля "myFile" (как в автотесте)
	file, header, err := r.FormFile("myFile")
	if err != nil {
		log.Printf("Файл не найден: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Чтение содержимого файла
	data, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Ошибка чтения файла: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	content := strings.TrimSpace(string(data))
	if content == "" {
		http.Error(w, "Empty File", http.StatusBadRequest)
		return
	}

	// Конвертация
	converted, err := service.AutoConvert(content)
	if err != nil {
		log.Printf("Ошибка конвертации: %v", err)
		http.Error(w, "Conversion Error", http.StatusInternalServerError)
		return
	}

	// Создание  результата
	timestamp := strings.ReplaceAll(
		time.Now().UTC().Format("2006-01-02_15-04-05"),
		":", "-",
	)
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".txt"
	}
	outputFilename := "result_" + timestamp + ext

	outputFile, err := os.Create(outputFilename)
	if err != nil {
		log.Printf("Ошибка создания файла: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer outputFile.Close()

	if _, err := outputFile.WriteString(converted); err != nil {
		log.Printf("Ошибка записи файла: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Ответ клиенту
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(converted))
}
