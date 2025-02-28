package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

// Order - структура для заказа
type Order struct {
	CustomerName  string `json:"customer_name"`
	City          string `json:"city"`
	PhoneNumber   string `json:"phone_number"`
	Email         string `json:"email"`
	ArticleNumber string `json:"article_number"`
}

// Подключение к базе данных
func initDB(host string, port int, user, password, dbname string) {
	var err error
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к БД: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("❌ Ошибка соединения с БД: %v", err)
	}

	log.Println("✅ Подключение к базе данных установлено!")
}

// Функция для обработки создания заказа
func createOrder(w http.ResponseWriter, r *http.Request) {
	var order Order

	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "❌ Ошибка обработки JSON", http.StatusBadRequest)
		return
	}

	sqlStatement := `INSERT INTO orders (customer_name, city, phone_number, email, article_number) 
                   VALUES ($1, $2, $3, $4, $5) RETURNING id;`
	var orderID int

	err = db.QueryRow(sqlStatement, order.CustomerName, order.City, order.PhoneNumber, order.Email, order.ArticleNumber).Scan(&orderID)
	if err != nil {
		log.Printf("❌ Ошибка при создании заказа: %v\n", err)
		http.Error(w, "❌ Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Заказ успешно создан! ID: %d\n", orderID)

	response := map[string]interface{}{
		"message":  "✅ Заказ успешно создан!\n Номер вашего заказа: ",
		"order_id": orderID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Получаем порт из флага командной строки
	dbHost := flag.String("db-host", "127.0.0.1", "IP-адрес базы данных")
	dbPort := flag.Int("db-port", 5432, "Порт для подключения к БД")
	serverPort := flag.Int("server-port", 8080, "Порт для веб-сервера")
	flag.Parse()

	// Логируем порты
	log.Printf("📌 IP-адрес БД: %s", *dbHost)
	log.Printf("📌 Используется порт для БД: %d\n", *dbPort)
	log.Printf("📌 Сервер запустится на: %d\n", *serverPort)

	// Подключаемся к БД
	initDB(*dbHost, *dbPort, "postgres", "postgres", "postgres")
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/order", createOrder).Methods("POST")

	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/").Handler(fs)

	log.Printf("🚀 Сервер запущен на порту %d...", *serverPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *serverPort), r))
}
