package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DBController struct {
	Database           *sql.DB
	AddRecordStatement *sql.Stmt
	AddPathStatement   *sql.Stmt
}

func initDBController() *DBController {
	database, err := sql.Open("sqlite3", "./requests.db")
	if err != nil {
		fmt.Println("COnstruction DB impossible")
		fmt.Println(err)
	}
	// Création des chemins
	statement, err := database.Prepare(`CREATE TABLE IF NOT EXISTS paths (
		id INTEGER PRIMARY KEY,
		name TEXT,
		origin TEXT,
		destination TEXT,
		free_flow_travel_time INTEGER
	)`)
	if err != nil {
		fmt.Println("Création de table chemins impossible")
		fmt.Println(err)
	}
	statement.Exec()

	// Création de la table des résultats
	statement, err = database.Prepare(`CREATE TABLE IF NOT EXISTS requests (
		id INTEGER PRIMARY KEY,
		datetime TEXT,
		path_0 INTEGER
		)`)
	if err != nil {
		fmt.Println("Création table impossible")
		fmt.Println(err)
	}
	statement.Exec()

	// Create add record statement
	statement, err = database.Prepare(`INSERT INTO requests (
		datetime,
		path_0
		) VALUES (?, ?)`)
	if err != nil {
		fmt.Println("COnstruction Statement impossible")
		fmt.Println(err)
	}

	addPathStatement, err := database.Prepare(`INSERT INTO paths (
		name,
		origin,
		destination,
		free_flow_travel_time
		) VALUES (?, ?, ?, ?)`)
	fmt.Println("New Database created")
	return &DBController{
		Database:           database,
		AddRecordStatement: statement,
		AddPathStatement:   addPathStatement,
	}
}

func AddRecord(record *Result) {
	if theModel.DBController.Database != nil &&
		theModel.DBController.AddRecordStatement != nil &&
		theModel.DBController.AddPathStatement != nil {
		go AddRow(record)
		theModel.mu.Lock()
		_, err := theModel.DBController.AddRecordStatement.Exec(time.Now().Format("2006-01-02 15:04:05"), record.Rows[0].Elements[0].DistanceInTraffic.Value)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Correctly stored in DB")
		theModel.mu.Unlock()
	}
}

func AddRow(record *Result) {
	theModel.mu.Lock()
	// check if path exists in db
	rows, _ := theModel.DBController.Database.Query(`SELECT * FROM paths WHERE paths.name=="path_0"`)
	if !rows.Next() {
		fmt.Println("Creating new path in DB")
		theModel.DBController.AddPathStatement.Exec("path_0", record.Origines[0], record.Destinations[0], record.Rows[0].Elements[0].Duration.Value)
	}
	theModel.mu.Unlock()
}
