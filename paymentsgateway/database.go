package paymentsgateway

import (
	"database/sql"
	"log"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/Piusdan/payments-gateway/model"

)

func connectToDatabase() *sql.DB {
	db, err := sql.Open("postgres", "postgres://{password}:{user}@{db_host}/{db_name}?sslmode=disable")
	if err != nil {
		log.Fatalln(fmt.Errorf("Unable to connect to databseL %v", err))
	}
	model.SetDatabase(db)
	return db
}
