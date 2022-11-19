package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	r := chi.NewRouter()

	/** habilitar cors */
	r.Use(cors.Handler(cors.Options{

		/**Use this to allow specific origin hosts*/
		AllowedOrigins: []string{"https://*", "http://*"},

		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// **********************

	r.Use(middleware.Logger)
	r.Get("/{aux}", func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Imprimir body")
		data := chi.URLParam(r, "aux")
		datoRespuesta := buscarContenido(data)
		w.Write([]byte(datoRespuesta))
	})
	http.ListenAndServe(":3000", r)
}

func buscarContenido(data string) string {
	var query = `{
		"search_type": "matchphrase",
		"query":{
				"term": "%s"
		},
		"from": 0,
		"_source": []
	}`

	s := fmt.Sprintf(string(query), data)

	req, err := http.NewRequest("POST", "http://localhost:4080/api/mail_dir/_search", strings.NewReader(s))
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth("admin", "Complexpass#123")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	log.Println(resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
	return string(body)
}
