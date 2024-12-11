package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type DWNode struct {
	DB *sql.DB
}

func NewDWNode(connStr string) *DWNode {
	var db *sql.DB
	var err error

	// Retry logic until DB is available
	for {
		// Try to connect to the database
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Error connecting to the database: %v. Retrying...", err)
			time.Sleep(2 * time.Second)
			continue
		}

		// Try to ping the database to ensure it's reachable
		err = db.Ping()
		if err != nil {
			log.Printf("Database not ready: %v. Retrying...", err)
			time.Sleep(2 * time.Second)
			continue
		}

		log.Println("Database is ready!")
		break
	}

	// Create the table if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS data (
		key TEXT PRIMARY KEY,
		value TEXT
	)`)
	if err != nil {
		log.Fatalf("Error creating the table: %v", err)
	}

	return &DWNode{DB: db}
}
func (dw *DWNode) WriteHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Key   string `json:"key" xml:"key"`
		Value string `json:"value" xml:"value"`
	}

	if r.Header.Get("Content-Type") == "application/xml" {
		if err := xml.NewDecoder(r.Body).Decode(&data); err != nil {
			log.Printf("Error decoding XML: %v", err)
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
	} else {
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			log.Printf("Error decoding JSON: %v", err)
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
	}

	_, err := dw.DB.Exec("INSERT INTO data (key, value) VALUES ($1, $2) ON CONFLICT (key) DO UPDATE SET value = $2", data.Key, data.Value)
	if err != nil {
		log.Printf("Error inserting/updating: %v", err)
		http.Error(w, "Error writing to the database", http.StatusInternalServerError)
		return
	}

	log.Printf("Data inserted/updated: key=%s, value=%s", data.Key, data.Value)
	w.WriteHeader(http.StatusCreated)
}

func (dw *DWNode) ReadHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	log.Printf("Request to read key: %s", key)
	log.Printf("Accept header: %s", r.Header.Get("Accept"))

	row := dw.DB.QueryRow("SELECT value FROM data WHERE key = $1", key)

	var value string
	if err := row.Scan(&value); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Nonexistent key: %s", key)
			http.Error(w, "Nonexistent key", http.StatusNotFound)
		} else {
			log.Printf("Error reading from the database: %v", err)
			http.Error(w, "Error reading from the database", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("Successful read: key=%s, value=%s", key, value)

	acceptHeader := r.Header.Get("Accept")
	if containsXML(acceptHeader) {
		log.Printf("Returning XML response")
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		xml.NewEncoder(w).Encode(map[string]string{"key": key, "value": value})
	} else {
		log.Printf("Returning JSON response")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"key": key, "value": value})
	}
}

func containsXML(acceptHeader string) bool {
	return strings.Contains(strings.ToLower(acceptHeader), "application/xml")
}

func Run(port string, connStr string) {
	if connStr == "" {
		connStr = "user=replicator password=your_password dbname=dw_db host=localhost sslmode=disable"
	}

	dwNode := NewDWNode(connStr)
	defer dwNode.DB.Close()

	http.HandleFunc("/write", dwNode.WriteHandler)
	http.HandleFunc("/read", dwNode.ReadHandler)

	fmt.Println("DW node is running at http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func main() {
	Run("8082", "postgres://postgres:masterpassword@master:5432/mydatabase?sslmode=disable")
}
