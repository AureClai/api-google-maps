package main

import (
	"database/sql"
	"fmt"
	"time"

	"googlemaps.github.io/maps"

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
		origin_lat TEXT,
		origin_long TEXT,
		destination_lat TEXT,
		destination_long TEXT,
		polyline TEXT
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
		origin_lat,
		origin_long,
		destination_lat,
		destination_long,
		polyline
		) VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		fmt.Println("COnstruction Statement impossible")
		fmt.Println(err)
	}
	fmt.Println("New Database created")
	return &DBController{
		Database:           database,
		AddRecordStatement: statement,
		AddPathStatement:   addPathStatement,
	}
}

func AddRecord(record *maps.DistanceMatrixResponse) {
	if theModel.DBController.Database != nil &&
		theModel.DBController.AddRecordStatement != nil {
		theModel.mu.Lock()
		_, err := theModel.DBController.AddRecordStatement.Exec(time.Now().Format("2006-01-02 15:04:05"), record.Rows[0].Elements[0].DurationInTraffic.Seconds())
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Correctly stored in DB")
		theModel.mu.Unlock()
	}
}

func AddPath(infos *PathInfos, record []maps.Route) {
	fmt.Printf("%v#\n", theModel.DBController)
	if theModel.DBController.Database != nil &&
		theModel.DBController.AddPathStatement != nil {

		theModel.mu.Lock()
		fmt.Println("Creating new path in DB")
		result, err := theModel.DBController.AddPathStatement.Exec(
			infos.Name,
			infos.Coordinates.Origin.Lat,
			infos.Coordinates.Origin.Long,
			infos.Coordinates.Destination.Lat,
			infos.Coordinates.Destination.Long,
			record[0].OverviewPolyline.Points,
		)
		if err != nil {
			fmt.Println(err)
		}
		theModel.mu.Unlock()
		idLast, err := result.LastInsertId()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Correctly stored in DB - ID : %v", idLast)
		infos.ID = int(idLast)
		sendMessage(&Message{
			Name: "add path",
			Data: infos,
		})
		sendMessage(&Message{
			Name: "test callback",
			Data: map[string]string{
				"type":    "SUCCESS",
				"message": "Ajout de l'itinéraire résussi",
			},
		})
	}
}

func sendAllPaths() {
	allPaths := make([]PathInfos, 0)
	theModel.mu.Lock()
	rows, err := theModel.DBController.Database.Query(`SELECT
	id,
	name,
	origin_lat,
	origin_long,
	destination_lat,
	destination_long FROM paths`)
	if err != nil {
		fmt.Println(err)
	}
	var id int
	var name, origin_lat, origin_long, destination_lat, destination_long string
	theModel.mu.Unlock()

	for rows.Next() {
		rows.Scan(&id, &name, &origin_lat, &origin_long, &destination_lat, &destination_long)
		pathInfos := PathInfos{ID: id, Name: name}
		pathInfos.Coordinates.Origin.Lat = origin_lat
		pathInfos.Coordinates.Origin.Long = origin_long
		pathInfos.Coordinates.Destination.Lat = destination_lat
		pathInfos.Coordinates.Destination.Long = destination_long
		allPaths = append(allPaths, pathInfos)
	}
	sendMessage(&Message{
		Name: "init paths",
		Data: allPaths,
	})
}

func removePath(id int) {
	sqlStatement := `
	DELETE FROM paths
	WHERE id = $1;
	`
	theModel.mu.Lock()
	_, err := theModel.DBController.Database.Exec(sqlStatement, id)
	if err != nil {
		fmt.Println(err)
		theModel.mu.Unlock()
		return
	}
	theModel.mu.Unlock()
	sendMessage(&Message{
		Name: "remove path",
		Data: map[string]string{
			"id": fmt.Sprintf("%v", id),
		},
	})
	sendMessage(&Message{
		Name: "test callback",
		Data: map[string]string{
			"type":    "SUCCESS",
			"message": "Suppression de l'itinéraire résussie",
		},
	})

}

func sendPolyline(id int) {
	sqlStatement := `
	SELECT polyline
	FROM paths
	WHERE id = $1;
	`
	theModel.mu.Lock()
	rows, err := theModel.DBController.Database.Query(sqlStatement, id)
	if err != nil {
		fmt.Println(err)
	}
	theModel.mu.Unlock()
	if rows.Next() {
		var encodedPolyline string
		rows.Scan(&encodedPolyline)
		sendMessage(&Message{
			Name: "path map",
			Data: convertToGeoJSON(encodedPolyline),
		})
	} else {
		fmt.Println("ERROR ! no row ! 001")
	}

}
