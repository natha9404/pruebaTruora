package main

import (
	"database/sql"
	"net/http"
	"log"
	"encoding/json"
	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

//Domain structure
type Domain struct {
	ID 		string `json:"id"`
	Name	string `json:"domain"`
	Updated	string `json:"updated_at"`
}

type Domains []Domain

/* Function that receive the domain to search and return a data json with the next information:
servers, server_changed, ssl_grade, previous_ssl_grade, logo, tittle and is_down */
func getDomain(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var objectDomain Domain
	err := decoder.Decode(&objectDomain)

	if err != nil {
		log.Println(err)
		http.Error(w, "Failed Error", 500)
		return
	}

	db, err := sql.Open("postgres", "postgresql://root@localhost:26257/company_db?sslmode=disable")
	if err != nil {
		log.Println("error connecting to the database: ", err)
		return
	}

	var sql = "INSERT INTO domainregister (domain, updated_at) VALUES ($1, NOW())"
	log.Println(sql)

	//db.Exec(sql)
	   
	
	if _, err := db.Exec(sql, objectDomain.Name); err != nil {
		log.Println(err)
		http.Error(w, "Insertion Error", 500)
		return
	}	

	log.Printf(objectDomain.Name)

	w.Write([]byte("welcome get Domain"))
	log.Printf("getDomain")
}

func getDomains(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome getDomains"))
	log.Printf("getDomains")
}

func connectDB(){
	// Connect to the "pruebaTruora" database.
	/*db, err := sql.Open("postgres", "postgresql://root@localhost:26257/pruebaTruora?sslmode=disable")
	if err != nil {
		log.Println("error connecting to the database: ", err)
	}*/
}

func main() {
	route := chi.NewRouter()
	route.Get("/", getDomain)
	route.Get("/getDomains", getDomains)
	log.Printf("starting server")
	connectDB()
	log.Fatal(http.ListenAndServe(":3000", route))

}
