package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	dbUser     = "joshua468"
	dbPassword = ""
	dbName     = "health_tracker"
)

var db *sql.DB

type HealthData struct {
	ID       int       `json:"id"`
	Date     time.Time `json:"date"`
	Weight   float64   `json:"weight"`
	Steps    int       `json:"steps"`
	Sleep    float64   `json:"sleep"`
	Calories int       `json:"calories"`
	Water    float64   `json:"water"`
}

func main() {
	dbURI := dbUser + ":" + dbPassword + "@tcp(127.0.0.1:3306)/" + dbName
	var err error
	db, err = sql.Open("mysql", dbURI)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/healthdata", getHealthData).Methods("GET")
	r.HandleFunc("/healthdata/{id}", getHealthDataByID).Methods("GET")
	r.HandleFunc("/healthdata", createHealthData).Methods("POST")
	r.HandleFunc("/healthdata/{id}", updateHealthData).Methods("PUT")
	r.HandleFunc("/healthdata/{id}", deleteHealthData).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func getHealthData(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM health_data")
	if err != nil {
		log.Println("Error retrieving health data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var healthData []HealthData
	for rows.Next() {
		var data HealthData
		err := rows.Scan(&data.ID, &data.Date, &data.Weight, &data.Steps, &data.Sleep, &data.Calories, &data.Water)
		if err != nil {
			log.Println("Error scanning health data row:", err)
			continue
		}
		healthData = append(healthData, data)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over health data rows:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(healthData)
}

func getHealthDataByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var data HealthData
	err := db.QueryRow("SELECT * FROM health_data WHERE id = ?", params["id"]).Scan(&data.ID, &data.Date, &data.Weight, &data.Steps, &data.Sleep, &data.Calories, &data.Water)
	if err != nil {
		log.Println("Error retrieving health data by ID:", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(data)
}

func createHealthData(w http.ResponseWriter, r *http.Request) {
	var data HealthData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Println("Error decoding JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err := db.Exec("INSERT INTO health_data (date, weight, steps, sleep, calories, water) VALUES (?, ?, ?, ?, ?, ?)", data.Date, data.Weight, data.Steps, data.Sleep, data.Calories, data.Water)
	if err != nil {
		log.Println("Error inserting health data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func updateHealthData(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var data HealthData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Println("Error decoding JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err := db.Exec("UPDATE health_data SET date=?, weight=?, steps=?, sleep=?, calories=?, water=? WHERE id=?", data.Date, data.Weight, data.Steps, data.Sleep, data.Calories, data.Water, params["id"])
	if err != nil {
		log.Println("Error updating health data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func deleteHealthData(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	_, err := db.Exec("DELETE FROM health_data WHERE id=?", params["id"])
	if err != nil {
		log.Println("Error deleting health data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
