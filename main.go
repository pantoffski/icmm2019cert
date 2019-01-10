package main

import (
	"icmm2019cert/cert"
	"icmm2019cert/cfg"
	"icmm2019cert/database"
	"log"
	"net/http"
	"strconv"
	"time"

	"math/rand"
	"github.com/globalsign/mgo/bson"
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
	router.Get("/{bibNO}", genCert)
	router.Get("/name/{bibNO}", getName)
	return router
}
func genCert(w http.ResponseWriter, r *http.Request) {
	bibNO, _ := strconv.Atoi(chi.URLParam(r, "bibNO"))
	// txt:=r.URL.Query().Get("txt")
	// x:=r.URL.Query().Get("x")
	// y:=r.URL.Query().Get("y")
	// size:=r.URL.Query().Get("size")
	// cert.Image(bibNO,txt,x,y,size, w)
	cert.Image(bibNO, w)
}
func getName(w http.ResponseWriter, r *http.Request) {
	bibNO, _ := strconv.Atoi(chi.URLParam(r, "bibNO"))
	db := database.GetDB()
	defer db.Session.Close()
	defer r.Body.Close()
	runner := cert.Runner{}
	err := db.C("bib_subscribers").Find(bson.M{"bibNumber": bibNO}).One(&runner)
	if err != nil {
		http.Error(w, "runner not found", 404)
		return
	}
	render.PlainText(w, r, runner.FName+" "+runner.LName)
}
func main() {

	http.DefaultClient.Timeout = time.Minute * 1
	
    rand.Seed(time.Now().Unix())
	r := InitRoutes()
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s -> %s\n", method, route)
		return nil
	}
	if err := chi.Walk(r, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error())
	}
	log.Fatal(http.ListenAndServe(":"+cfg.Getenv("PORT"), r))
}
