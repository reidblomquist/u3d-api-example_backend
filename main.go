package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
	"sync"
)

func main() {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	api.Use(&rest.CorsMiddleware{
		RejectNonCorsRequests: false,
		OriginValidator: func(origin string, request *rest.Request) bool {
			return origin == "http://localhost:3000" || origin == "http://localhost:8000"
		},
		AllowedMethods: []string{"GET", "POST", "PUT"},
		AllowedHeaders: []string{
			"Accept", "Content-Type", "X-Custom-Header", "Origin"},
		AccessControlAllowCredentials: true,
		AccessControlMaxAge:           3600,
	})
	router, err := rest.MakeRouter(
		rest.Get("/countries", GetAllCountries),
		rest.Post("/countries", PostCountry),
		rest.Get("/countries/:code", GetCountry),
		rest.Delete("/countries/:code", DeleteCountry),
		rest.Get("/rgba", GetRgba),
		rest.Post("/rgba", PostRgba),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}

type Country struct {
	Code string
	Name string
}

type Rgba struct {
	R float32;
	G float32;
	B float32;
	A float32;
}

var countryStore = map[string]*Country{}

var countryLock = sync.RWMutex{}

func GetCountry(w rest.ResponseWriter, r *rest.Request) {
	code := r.PathParam("code")

	countryLock.RLock()
	var country *Country
	if countryStore[code] != nil {
		country = &Country{}
		*country = *countryStore[code]
	}
	countryLock.RUnlock()

	if country == nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(country)
}

func GetAllCountries(w rest.ResponseWriter, r *rest.Request) {
	countryLock.RLock()
	countries := make([]Country, len(countryStore))
	i := 0
	for _, country := range countryStore {
		countries[i] = *country
		i++
	}
	countryLock.RUnlock()
	w.WriteJson(&countries)
}

func PostCountry(w rest.ResponseWriter, r *rest.Request) {
	country := Country{}
	err := r.DecodeJsonPayload(&country)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if country.Code == "" {
		rest.Error(w, "country code required", 400)
		return
	}
	if country.Name == "" {
		rest.Error(w, "country name required", 400)
		return
	}
	countryLock.Lock()
	countryStore[country.Code] = &country
	countryLock.Unlock()
	w.WriteJson(&country)
}

func DeleteCountry(w rest.ResponseWriter, r *rest.Request) {
	code := r.PathParam("code")
	countryLock.Lock()
	delete(countryStore, code)
	countryLock.Unlock()
	w.WriteHeader(http.StatusOK)
}

var rgbaStore = map[int]*Rgba{}

var rgbaLock = sync.RWMutex{}

func GetRgba(w rest.ResponseWriter, r *rest.Request) {
	rgbaLock.RLock()
	rgba := rgbaStore[0]
	rgbaLock.RUnlock()
	w.WriteJson(&rgba)
}

func PostRgba(w rest.ResponseWriter, r *rest.Request) {
	rgba := Rgba{}
	err := r.DecodeJsonPayload(&rgba)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rgbaLock.Lock()
	rgbaStore[0] = &rgba
	rgbaLock.Unlock()
	w.WriteJson(&rgba)
}