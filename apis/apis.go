package apis

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gokultp/metar/cache"
	"github.com/gokultp/metar/metar"
	"github.com/gorilla/mux"
)

type API struct {
	Cache  cache.Cache
	Router *mux.Router
}

// NewAPI will creates a new API object
func NewAPI(redisURL, redisPassword string, redisDB int) *API {
	var api API

	getData := func(key string) (interface{}, error) {
		return metar.GetData(key)
	}

	if redisURL == "" {
		api.Cache = cache.NewInMemCache(getData)
	} else {
		api.Cache = cache.NewRedisCache(getData, redisURL, redisPassword, redisDB)
	}
	api.Router = mux.NewRouter()
	return &api
}

// InitRoutes will initialize rotes
func (api *API) InitRoutes() {
	fmt.Println("InitRoutes")
	api.Router.HandleFunc("/metar/ping", api.Ping).Methods(http.MethodGet)
	api.Router.HandleFunc("/metar/info", api.MetarInfo).Methods(http.MethodGet)
}

// Listen will start the server
func (api *API) Listen(addr string) {

	if addr == ":" {
		addr = ":8080"
	}
	fmt.Println("Listening at ", addr)

	log.Fatal(http.ListenAndServe(addr, api.Router))
}

// Ping is a test ping  call to the server
func (api *API) Ping(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"data": "pong",
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
	return
}

// MetarInfo will get the metar info of the given station
func (api *API) MetarInfo(w http.ResponseWriter, r *http.Request) {
	station := r.URL.Query().Get("scode")
	nocache := r.URL.Query().Get("nocache")

	if nocache == "1" {
		data, err := api.Cache.GetDataFromSource(station)

		if err != nil {
			jsonResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		jsonResponse(w, http.StatusOK, data)
	}

	data, err := api.Cache.GetDataFromCache(station)

	if err != nil {
		jsonResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, data)
	return

}

func jsonResponse(w http.ResponseWriter, status int, payload interface{}) {
	response := make(map[string]interface{})

	if status < 400 {
		// positive response
		response["status"] = true
		response["data"] = payload
	} else {
		// error response
		response["status"] = false
		response["errors"] = payload
	}

	jsonData, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(jsonData))
}
