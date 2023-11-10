package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	// Connect to your PostgreSQL database
	db, err := sql.Open("postgres", "host=localhost port=5432 user=atif password=atif dbname=aqaryint sslmode=disable")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Execute the seat swapping SQL query
	query := `
		UPDATE Seat AS s1
		SET student = s2.student
		FROM Seat AS s2
		WHERE s1.id = s2.id - 1 AND s1.id % 2 = 1;
	`

	_, err = db.Exec(query)
	if err != nil {
		panic(err.Error())
	}

	// Fetch the result after swapping
	rows, err := db.Query("SELECT * FROM Seat ORDER BY id ASC")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	// Print the result
	fmt.Println("Output:")
	fmt.Println("+----+---------+")
	fmt.Println("| id | student |")
	fmt.Println("+----+---------+")
	for rows.Next() {
		var id int
		var student string
		err := rows.Scan(&id, &student)
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("| %d | %s |\n", id, student)
	}
	fmt.Println("+----+---------+")
}
