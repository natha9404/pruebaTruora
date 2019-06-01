package main

import (
	"database/sql"
	"net/http"
	"log"
	"encoding/json"
	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"strings"
)

var db *sql.DB

//Domain structure
type Domain struct {
	ID 		string `json:"id"`
	Name	string `json:"domain"`
	Updated	string `json:"updated_at"`
}

type Server struct {
	Address *string `json:"address",omitempty`
	Ssl_grade *string `json:"ssl_grade",omitempty`
	Country *string `json:"country",omitempty`
	Owner *string `json:"owner",omitempty`
 }
 
 type DomainInfo struct{
	 Servers []Server `json:"servers"`
	 Logo *string `json:"logo",omitempty`
	 Title *string `json:"title",omitempty`
 
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

	//Funciont that saves the search history
	saveSearchHistory(objectDomain.Name)
	searchInfoServer(objectDomain.Name)
	//Function that searches information of servers of a domain.


	
	log.Printf(objectDomain.Name)

	w.Write([]byte("welcome get Domain"))
	log.Printf("getDomain")
}

func getDomains(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome getDomains"))
	log.Printf("getDomains")
}

func saveSearchHistory(domain string){
	db, err := sql.Open("postgres", "postgresql://root@localhost:26257/company_db?sslmode=disable")
	if err != nil {
		log.Println("error connecting to the database: ", err)
		return
	}

	var sql = "INSERT INTO domainregister (domain, updated_at) VALUES ($1, NOW())"
	log.Println(sql)	   
	
	if _, err := db.Exec(sql, domain); err != nil {
		panic(err)
		return
	}

}

func searchInfoServer(domain string){

	url := "https://api.ssllabs.com/api/v3/analyze?host="+domain
	log.Println("url")

	log.Println(url)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	var infoSSL map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&infoSSL)
	if err != nil {
		log.Fatalln(err)
	}
	
	var endPoints = infoSSL["endpoints"].([]interface{})

	var whoIs = searchWhoIs(domain)

	var favicon = getFavicon(domain)

	var title = getTitle(domain)

	var title2 = getTitle2(domain)


	/* Country : whoIs["WhoisRecord"].(map[string]interface{})["registrant"].(map[string]interface{})["country"].(string),
	Owner : whoIs["WhoisRecord"].(map[string]interface{})["registrant"].(map[string]interface{})["organization"].(string)}
	 */
	
	log.Println("endpoints")
	log.Println(endPoints)
	log.Println("whois")
	log.Println(whoIs)
	log.Println("favicon")
	log.Println(favicon)
	log.Println("title")
	log.Println(title)
	log.Println("title2")
	log.Println(title2)
	log.Println("terminoWhois")

}

func searchWhoIs(domain string) (interface{}){

	url := "https://www.whoisxmlapi.com/whoisserver/WhoisService?apiKey=at_5UhpXqA9prtTSlHrPE2UJiUyASacC&domainName="+domain+"&outputFormat=json"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	var whoIs map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&whoIs)
	if err != nil {
		log.Fatalln(err)
	}
	
	return whoIs
	
}

func getFavicon(domain string) (interface{}){

	url := "https://besticon-demo.herokuapp.com/allicons.json?url="+domain

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	var favicon map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&favicon)
	if err != nil {
		log.Fatalln(err)
	}
	
	return favicon

}

func getTitle(domain string) (string){

	resp, err := http.Get("https://"+domain)
    if err != nil{
        log.Fatal(err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200{
        log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
    }

    doc, err := goquery.NewDocumentFromReader(resp.Body)
  	if err != nil {
        log.Fatal(err)
    }
	title := (doc.Find("title").Text())


	
	return title
}

func getTitle2(domain string) (string){

	url := "http://"+domain
	resp, err := http.Get(url)
	// handle the error if there is one
	if err != nil {
		panic(err)
	}
	// do this now so it won't be forgotten
	defer resp.Body.Close()
	// reads html as a slice of bytes

	dataInBytes, err := ioutil.ReadAll(resp.Body)
    pageContent := string(dataInBytes)
	evaluatetitle:= false
    // Find a substr
    titleStartIndex := strings.Index(pageContent, "<title>")
    if titleStartIndex == -1 {
        log.Println("No title element found")
		evaluatetitle = true
        
    }

    // The start index of the title is the index of the first
    // character, the < symbol. We don't want to include
    // <title> as part of the final value, so let's offset
    // the index by the number of characers in <title>
    titleStartIndex += 7
	
    // Find the index of the closing tag
    titleEndIndex := strings.Index(pageContent, "</title>")
    if titleEndIndex == -1 {
        log.Println("No closing tag for title found.")
        evaluatetitle = true
    }
	
    // (Optional)
    // Copy the substring in to a separate variable so the
    // variables with the full document data can be garbage collected
	pageTitle := " "
	if (evaluatetitle){
		pageTitle = " " 
	}else{
		pageTitle = string([]byte(pageContent[titleStartIndex:titleEndIndex]))
	}
	
	return pageTitle
}

// func connectDB()(*sql.DB, err){
// 	// Connect to the "pruebaTruora" database.
// 	/*db, err := sql.Open("postgres", "postgresql://root@localhost:26257/pruebaTruora?sslmode=disable")
// 	if err != nil {
// 		log.Println("error connecting to the database: ", err)
// 	}*/
// }

func main() {
	route := chi.NewRouter()
	route.Get("/", getDomain)
	route.Get("/getDomains", getDomains)
	log.Printf("starting server")
	log.Fatal(http.ListenAndServe(":3000", route))

}
