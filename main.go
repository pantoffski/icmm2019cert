package main

import (
	"strconv"
	"log"
	"icmm2019cert/cfg"
	"icmm2019cert/cert"
	"net/http"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)
 
// InitRoutes the main route config
func InitRoutes() *chi.Mux {
	router := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	router.Use(cors.Handler,
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	) 
	router.Get("/{bibNO}",genCert)
	return router
}
func genCert(w http.ResponseWriter, r *http.Request) {
	bibNO,_:=strconv.Atoi(chi.URLParam(r,"bibNO"))
	cert.Image(bibNO,w)
}
func main() {
	r := InitRoutes()
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s -> %s\n", method, route)
		return nil
	}
	if err := chi.Walk(r, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error())
	}
	log.Fatal(http.ListenAndServe(":"+cfg.Getenv("CERT2019_PORT"), r))
}
