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
		rest.Get("/rgba", GetRgba),
		rest.Post("/rgba", PostRgba),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}

type Rgba struct {
	R float32
	G float32
	B float32
	A float32
}

var rgbaStore = map[int]*Rgba{}

var rgbaLock = sync.RWMutex{}

func GetRgba(w rest.ResponseWriter, r *rest.Request) {
	rgbaLock.RLock()
	if rgbaStore[0] == nil {
		rgba := Rgba{}
		rgbaStore[0] = &rgba
	}
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
