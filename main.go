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
)

// SensorValues structure that contains all the individual measured values
type SensorValues struct {
	ID          int     `json:"id"`
	Temperature int     `json:"temperature"`
	Humidity    float32 `json:"humidity"`
	Co2         int     `json:"co2"`
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "testuser"
	dbPass := "testpassword"
	dbName := "SENSORDATA"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}
func getTime() string {
	time := fmt.Sprint(time.Now().Format("15:04:05"))
	return time
}

func getReadingsDB(w http.ResponseWriter, req *http.Request) {

	var sensVal SensorValues
	db := dbConn()
	rows, err := db.Query("SELECT id, Temperature,Humidity,CO2 FROM READINGS ")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&sensVal.ID, &sensVal.Temperature, &sensVal.Humidity, &sensVal.Co2)
		if err != nil {
			log.Fatal(err)
		}
		bytes, _ := json.MarshalIndent(sensVal, "", " ")
		fmt.Fprintf(w, string(bytes))
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

}

func getReadingDB(w http.ResponseWriter, req *http.Request) {
	ID := req.URL.Query().Get("ID")
	var sensVal SensorValues
	db := dbConn()
	rows, err := db.Query("SELECT id, Temperature,Humidity,CO2 FROM READINGS WHERE id = ? ", ID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&sensVal.ID, &sensVal.Temperature, &sensVal.Humidity, &sensVal.Co2)
		if err != nil {
			log.Fatal(err)
		}
		bytes, _ := json.MarshalIndent(sensVal, "", " ")
		fmt.Fprintf(w, string(bytes))
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

}

func postReadingDB(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	var sensVal SensorValues
	err = json.Unmarshal(body, &sensVal)
	if err != nil {
		println(err)
	}
	db := dbConn()
	stmt, err := db.Prepare("INSERT INTO READINGS(id, Temperature, Humidity, Co2, Time) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		log.Print(err)
	}

	_, err = stmt.Exec(sensVal.ID, sensVal.Temperature, sensVal.Humidity, sensVal.Co2, getTime())
	if err != nil {
		log.Print(err)
	}
}

func main() {

	http.HandleFunc("/getReadingsDB", getReadingsDB)
	http.HandleFunc("/getReadingDB", getReadingDB)
	http.HandleFunc("/postReadingDB", postReadingDB)

	http.ListenAndServe(":8090", nil)
	// koristit /configuration/:id/ a ne /configuration?id=1
}
