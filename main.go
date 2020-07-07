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
	Time        string `json:"time"`
}
type dbHandler struct {
	db *sql.DB
}

func dbConn() (* dbHandler) {
	dbDriver := "mysql"
	dbUser := "testuser"
	dbPass := "testpassword"
	dbName := "SENSORDATA"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		fmt.Println(err)
	}
	var DataBase dbHandler
	DataBase.db = db
	return &DataBase
}


func getTime() string {
	time := fmt.Sprint(time.Now().Format("15:04:05"))
	return time
}

func (dbHandler *dbHandler)getReadings(w http.ResponseWriter, req *http.Request) {
	var sensVal SensorValues
	var readingSlice []SensorValues

	rows, err := dbHandler.db.Query("SELECT id, Temperature,Humidity,CO2,Time FROM READINGS ")
	defer rows.Close()
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		err := rows.Scan(&sensVal.ID, &sensVal.Temperature, &sensVal.Humidity, &sensVal.Co2, &sensVal.Time)
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

func (dbHandler *dbHandler)getReading(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	if params["id"] == "" {
		log.Println("No ID")
		http.Error(w, "No ID", http.StatusNoContent)
		return
	}
	var sensVal SensorValues
	var readingSlice []SensorValues
	rows, err := dbHandler.db.Query("SELECT id, Temperature,Humidity,CO2,Time FROM READINGS WHERE id = ? ", params["id"])
	defer rows.Close()
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		err := rows.Scan(&sensVal.ID, &sensVal.Temperature, &sensVal.Humidity, &sensVal.Co2, &sensVal.Time)
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

func (dbHandler *dbHandler)postReading(w http.ResponseWriter, req *http.Request) {
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
	sensVal.Time = getTime()
	stmt, err := dbHandler.db.Prepare("INSERT INTO READINGS(id, Temperature, Humidity, Co2, Time) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(sensVal.ID, sensVal.Temperature, sensVal.Humidity, sensVal.Co2, sensVal.Time)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
func(dbHandler *dbHandler) deleteReading(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	if params["id"] == "" {
		log.Println("No ID")
		http.Error(w, "No ID", http.StatusNoContent)
		return
	}
	_, err := dbHandler.db.Query("DELETE FROM READINGS WHERE id=?", params["id"])
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func(dbHandler *dbHandler) updateReading(w http.ResponseWriter, req *http.Request) {
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
	sensVal.Time = getTime()
	_, err = dbHandler.db.Exec("UPDATE READINGS SET Temperature = ?, Humidity = ?, CO2 = ?, Time = ? where id = ?", sensVal.Temperature, sensVal.Humidity, sensVal.Co2, sensVal.Time, params["id"])
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
	DBconn := dbConn()
	defer DBconn.db.Close()
	router := mux.NewRouter()
	router.HandleFunc("/getReadings", AccessControl(DBconn.getReadings))
	router.HandleFunc("/getReading/{id}", AccessControl(DBconn.getReading))
	router.HandleFunc("/postReading", AccessControl(DBconn.postReading))
	router.HandleFunc("/deleteReading/{id}", AccessControl(DBconn.deleteReading))
	router.HandleFunc("/updateReading/{id}", AccessControl(DBconn.updateReading))
	log.Println("Server started...")
	http.ListenAndServe(":8090", router)

}
