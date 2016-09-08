package main

import (
	"image"
	"image/color"
	"testing"
)

func TestCreateImage(t *testing.T) {
	width, height := 200, 100
	img := *createImage(width, height)

	size := image.Rectangle{image.Point{0, 0}, image.Point{width, height}}
	if img.Bounds() != size {
		t.Error(img.Bounds(), " != ", size)
	}

	black_color := color.RGBA{0, 0, 0, 255}
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if img.At(x, y) != black_color {
				t.Error("image [", x, ",", y, "] != ", black_color)
			}
		}
	}
}

func TestGetStats(t *testing.T) {
	if getStats() != "{\"num_images\":0,\"average_width_px\":0,\"average_height_px\":0}" {
		t.Error("{\"num_images\":0,\"average_width_px\":0,\"average_height_px\":0} != ", getStats())
	}
	stats_json.CountImages = 5
	stats_json.total_width = 5000
	stats_json.total_height = 1000
	if getStats() != "{\"num_images\":5,\"average_width_px\":1000,\"average_height_px\":200}" {
		t.Error("{\"num_images\":5,\"average_width_px\":1000,\"average_height_px\":200} != ", getStats())
	}
}
