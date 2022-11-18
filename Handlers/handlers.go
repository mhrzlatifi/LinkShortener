package Handlers

import (
	"encoding/json"
	"linkShortener/DB"
	"math/rand"
	"net/http"
	"time"
)

const BaseLink string = "link-shortener/"

func InsertLink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content_Type", "application/json")

	link := r.URL.Query().Get("link")
	uniqueHash := generateHash(6)

	short, err := insertLinkInDB(link, uniqueHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error_message": err.Error()})
		return
	} else {
		setInRedis(link, uniqueHash)
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"short_link": BaseLink + short})
}

func setInRedis(link, hash string) {
	DB.RDB.Set(DB.CTX, hash, link, 5*time.Minute)
	return
}

func GetGeneralLink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content_Type", "application/json")

	short := r.URL.Query().Get("short")

	// check Redis
	link, _ := DB.RDB.Get(DB.CTX, short).Result()

	var err error
	if link == "" {
		// get from MySQL
		link, err = getOriginalLinkFromDB(short)
	}
	if link == "" {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error_message": "link doesn't exist"})
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error_message": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"original_link": link})
}

func getOriginalLinkFromDB(h string) (link string, err error) {
	row := DB.MYSQL.QueryRow("SELECT original_link FROM links WHERE hash = ?", h)

	err = row.Scan(&link)
	return
}

func generateHash(n int) string {
	rand.Seed(time.Now().UnixNano())

	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func insertLinkInDB(link, hash string) (short string, err error) {
	// check if exists
	row := DB.MYSQL.QueryRow("SELECT hash FROM links WHERE original_link = ?", link)
	_ = row.Scan(&short)

	if short == "" {
		// insert new record
		row = DB.MYSQL.QueryRow("INSERT INTO links (original_link, hash) VALUES (?,?)", link, hash)
		err = row.Scan()
		short = hash
		if err.Error() == "sql: no rows in result set" {
			err = nil
		}
	}
	return
}
