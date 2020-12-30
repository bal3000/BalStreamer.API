package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/bal3000/BalStreamer.API/models"
	"github.com/streadway/amqp"
)

// CastController - controller for casting to chromecast
type CastController struct {
	Database     *sql.DB
	RabbitMQ     *amqp.Channel
	ExchangeName string
}

// NewCastController - constructor to return new controller while passing in dependacies
func NewCastController(db *sql.DB, ch *amqp.Channel, en string) *CastController {
	return &CastController{Database: db, RabbitMQ: ch, ExchangeName: en}
}

// CastStream - streams given data to given chromecast
func (controller *CastController) CastStream(res http.ResponseWriter, req *http.Request) {
	castCommand := new(models.StreamToCast)
	if err := convertJSON(req, castCommand); err != nil {
		log.Println(err)
		respondWithError(res, http.StatusInternalServerError, err.Error())
	}

	// Send to chromecast
	cast := models.StreamToChromecastEvent{
		ChromeCastToStream: castCommand.Chromecast,
		Stream:             castCommand.StreamURL,
		StreamDate:         time.Now(),
	}

	b, err := json.Marshal(cast)
	if err != nil {
		log.Fatalln(err)
	}

	err = controller.RabbitMQ.Publish(
		controller.ExchangeName, // exchange
		"",                      // routing key
		false,                   // mandatory
		false,                   // immediate
		amqp.Publishing{
			ContentType: "application/vnd.masstransit+json",
			Body:        []byte(b),
		})

	if err != nil {
		log.Fatalln(err)
	}

	respondWithJSON(res, http.StatusNoContent, nil)
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
