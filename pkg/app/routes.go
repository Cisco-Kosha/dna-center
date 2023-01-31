package app

import (
	"encoding/json"
	"github.com/kosha/dna-center/pkg/httpclient"
	"github.com/kosha/dna-center/pkg/logger"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"strings"
)

const (
	ApiKey    = "API_KEY"
	BasicAuth = "BASIC_AUTH"
)

func (a *App) commonMiddleware(log logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//Allow CORS here By * or specific origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")

		if (*r).Method == "OPTIONS" {
			w.WriteHeader(200)
			return
		}

		serverUrl := a.Cfg.GetServerURL()
		requestUri := r.RequestURI
		method := r.Method
		queryParams := r.URL.Query().Encode()

		var contentTypeHeaderFound bool

		url := serverUrl + requestUri
		if queryParams != "" && !strings.Contains(requestUri, "?") {
			url += "?" + queryParams
		}

		var c interface{}
		decoder := json.NewDecoder(r.Body)
		_ = decoder.Decode(&c)
		defer r.Body.Close()

		headers := make(map[string]string)
		// Loop over header names
		if len(r.Header) > 0 {
			for name, values := range r.Header {
				// Loop over all values for the name.
				if strings.ToLower(name) == "content-type" {
					contentTypeHeaderFound = true
				}
				for _, value := range values {
					if name != "" && value != "" {
						headers[name] = value
					}
				}
			}
		}
		// use application/json as default content type
		if !contentTypeHeaderFound {
			headers["Content-Type"] = "application/json; charset=utf-8"
		}

		username, password := a.Cfg.GetUsernameAndPassword()

		res, statusCode, err := httpclient.MakeHttpCall(headers, username, password, method, a.Cfg.GetServerURL(), url, c, log, a.Cfg.IsSecure())
		if err != nil {
			a.Log.Errorf("Encountered an error while making a call: %v\n", err)
			respondWithError(w, statusCode, err.Error())
			return
		}
		if res == nil {
			respondWithJSON(w, statusCode, res)
		}
		respondWithJSON(w, statusCode, res)
		return

	})
}

func (a *App) InitializeRoutes(log logger.Logger) {
	a.Router.PathPrefix("/").Handler(a.commonMiddleware(log)).Methods("GET", "POST", "PUT", "DELETE", "OPTIONS")

	// Swagger
	a.Router.PathPrefix("/docs").Handler(httpSwagger.WrapHandler)
}
