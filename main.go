package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// SensorValues structure that contains all the individual measured values
type SensorValues struct {
	ID          string `json:"id"`
	Temperature string `json:"temperature"`
	Humidity    string `json:"humidity"`
	Co2         string `json:"co2"`
	TIme        string `json:"time"`
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "testuser"
	dbPass := "testpassword"
	dbName := "SENSORDATA"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		fmt.Println(err)
	}
	return db
}

var db = dbConn()

func getTime() string {
	time := fmt.Sprint(time.Now().Format("15:04:05"))
	return time
}

func getReadings(w http.ResponseWriter, req *http.Request) {
	var sensVal SensorValues
	var readingSlice []SensorValues

	rows, err := db.Query("SELECT id, Temperature,Humidity,CO2,Time FROM READINGS ")
	defer rows.Close()
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		err := rows.Scan(&sensVal.ID, &sensVal.Temperature, &sensVal.Humidity, &sensVal.Co2, &sensVal.TIme)
		if err != nil {
			log.Print(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		readingSlice = append(readingSlice, sensVal)
	}
	bytes, _ := json.MarshalIndent(readingSlice, "", " ")
	fmt.Fprintf(w, string(bytes))
}

func getReading(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	if params["id"] == "" {
		log.Println("No ID")
		http.Error(w, "No ID", http.StatusNoContent)
		return
	}
	var sensVal SensorValues
	var readingSlice []SensorValues
	rows, err := db.Query("SELECT id, Temperature,Humidity,CO2,Time FROM READINGS WHERE id = ? ", params["id"])
	defer rows.Close()
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		err := rows.Scan(&sensVal.ID, &sensVal.Temperature, &sensVal.Humidity, &sensVal.Co2, &sensVal.TIme)
		if err != nil {
			log.Print(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		readingSlice = append(readingSlice, sensVal)
	}
	bytes, _ := json.MarshalIndent(readingSlice, "", " ")
	fmt.Fprintf(w, string(bytes))
}

func postReading(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	var sensVal SensorValues
	err = json.Unmarshal(body, &sensVal)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if sensVal.ID == "" {
		log.Println("No ID")
		http.Error(w, "No ID", http.StatusNoContent)
		return
	}
	sensVal.TIme = getTime()
	stmt, err := db.Prepare("INSERT INTO READINGS(id, Temperature, Humidity, Co2, Time) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(sensVal.ID, sensVal.Temperature, sensVal.Humidity, sensVal.Co2, sensVal.TIme)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
func deleteReading(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	if params["id"] == "" {
		log.Println("No ID")
		http.Error(w, "No ID", http.StatusNoContent)
		return
	}
	_, err := db.Query("DELETE FROM READINGS WHERE id=?", params["id"])
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func updateReading(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	if params["id"] == "" {
		log.Println("No ID")
		http.Error(w, "No ID", http.StatusNoContent)
		return
	}
	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	var sensVal SensorValues
	err = json.Unmarshal(body, &sensVal)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_, err = db.Exec("UPDATE READINGS SET Temperature = ?, Humidity = ?, CO2 = ?, Time = ? where id = ?", sensVal.Temperature, sensVal.Humidity, sensVal.Co2, getTime(), params["id"])
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

}

// AccessControl middleware function inserts access control parameters
func AccessControl(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "3.5") //firefox
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		if req.Method == "OPTIONS" {
			return
		}
		handler.ServeHTTP(w, req)
	}

}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/getReadings", AccessControl(getReadings))
	router.HandleFunc("/getReading/{id}", AccessControl(getReading))
	router.HandleFunc("/postReading", AccessControl(postReading))
	router.HandleFunc("/deleteReading/{id}", AccessControl(deleteReading))
	router.HandleFunc("/updateReading/{id}", AccessControl(updateReading))
	http.ListenAndServe(":8090", router)

}
