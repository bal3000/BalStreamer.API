package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/bal3000/BalStreamer.API/models"
	// adapator for postgres
	_ "github.com/lib/pq"
)

// CastController - controller for casting to chromecast
type CastController struct {
	Database *sql.DB
}

// NewCastController - constructor to return new controller while passing in dependacies
func NewCastController(connectionString string) *CastController {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}

	return &CastController{Database: db}
}

// CastStream - streams given data to given chromecast
func (controller *CastController) CastStream(res http.ResponseWriter, req *http.Request) {
	castCommand := new(models.StreamToCast)
	if err := convertJSON(req, castCommand); err != nil {
		log.Println(err)
		respondWithError(res, http.StatusInternalServerError, err.Error())
	}
	respondWithJSON(res, http.StatusNoContent, "WORKING")
}

func convertJSON(r *http.Request, toObj interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(toObj)
	defer r.Body.Close()
	return err
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if payload != nil {
		response, err := json.Marshal(payload)
		if err != nil {
			log.Println(err)
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		w.Write(response)
	}
}
