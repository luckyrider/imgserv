package main // import "github.com/ridfrustum/imgserv"

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"strings"
	"time"
)

type StatsJson struct {
	CountImages   int `json:"num_images"`
	AverageWidth  int `json:"average_width_px"`
	AverageHeight int `json:"average_height_px"`
	total_width   int
	total_height  int
}

var stats_json = &StatsJson{}

func updateStats(width, height int) {
	stats_json.total_width += width
	stats_json.total_height += height
	stats_json.CountImages++
}

func getStats() string {
	if stats_json.CountImages > 0 {
		stats_json.AverageWidth = stats_json.total_width / stats_json.CountImages
		stats_json.AverageHeight = stats_json.total_height / stats_json.CountImages
	}
	result, _ := json.Marshal(stats_json)
	return string(result)
}

func generateParams(r *http.Request) (string, int, int, error) {
	if r.Method != http.MethodGet {
		return "", 0, 0, errors.New("Bad request, GET method expected!")
	}
	url_info := "Check url /generate/{type:png or jpg}/{width_px}/{height_px}."

	params := strings.Split(r.URL.Path[1:], "/")
	if len(params) < 4 {
		return "", 0, 0, errors.New(url_info)
	}

	img_type := params[1]
	if img_type != "png" && img_type != "jpg" {
		return "", 0, 0, errors.New(url_info + " {img_type} need 'png' || 'jpg'")
	}

	width, err := strconv.Atoi(params[2])
	if err != nil || width <= 0 {
		return "", 0, 0, errors.New(url_info + " {width_px} need 'int' and > 0")
	}
	height, err := strconv.Atoi(params[3])
	if err != nil || height <= 0 {
		return "", 0, 0, errors.New(url_info + " {height_px} need 'int' and > 0")
	}
	return img_type, width, height, nil
}

func createImage(width, height int) *image.RGBA {
	myimage := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})
	for i := 0; i < len(myimage.Pix); i += 4 {
		myimage.Pix[i+3] = 255
	}
	return myimage
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	img_type, width, height, err := generateParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Generate [%s] %dx%d", img_type, width, height)

	c := make(chan bytes.Buffer)
	go func(img_type string, width, height int) {
		var buf bytes.Buffer
		myimage := *createImage(width, height)
		switch img_type {
		case "png":
			png.Encode(&buf, &myimage)
		case "jpg":
			jpeg_options := jpeg.Options{Quality: jpeg.DefaultQuality}
			jpeg.Encode(&buf, &myimage, &jpeg_options)
		}
		c <- buf
	}(img_type, width, height)
	buf := <-c

	switch img_type {
	case "png":
		w.Header().Set("Content-Type", "image/png")
	case "jpg":
		w.Header().Set("Content-Type", "image/jpeg")
	}
	fmt.Fprint(w, buf.String())
	updateStats(width, height)
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, getStats())
}

func main() {
	Port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		Port = 9000
	}
	addr := fmt.Sprintf(":%d", Port)
	log.Printf("Starting IMGSERV on port %d...", Port)

	http.HandleFunc("/generate/", generateHandler)
	http.HandleFunc("/stats", statsHandler)
	srv := &http.Server{
		Addr:           addr,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		// Handler:        http.DefaultServeMux,
	}
	log.Fatal(srv.ListenAndServe())
}
