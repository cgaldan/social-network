package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"path/filepath"
)

const (
	avatarPoolSize    = 24
	postMediaPoolSize = 14
	avatarSize        = 256
	postMediaW        = 800
	postMediaH        = 600
	imageURLPrefix    = "/uploads/"
)

// generateSeedImages writes a small pool of placeholder PNGs into the upload
// directory and returns URL paths for them. Images already on disk are reused
// without rewriting, so this is cheap to call on every seed run.
//
// Returns (avatarURLs, postMediaURLs).
func generateSeedImages(uploadDir string) ([]string, []string, error) {
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return nil, nil, fmt.Errorf("mkdir uploads: %w", err)
	}

	avatars := make([]string, 0, avatarPoolSize)
	for i := 0; i < avatarPoolSize; i++ {
		name := fmt.Sprintf("seed_avatar_%03d.png", i)
		if err := ensurePNG(filepath.Join(uploadDir, name), func() image.Image {
			a, b := avatarPalette(i)
			return gradientWithDots(avatarSize, avatarSize, a, b, int64(i))
		}); err != nil {
			return nil, nil, err
		}
		avatars = append(avatars, imageURLPrefix+name)
	}

	media := make([]string, 0, postMediaPoolSize)
	for i := 0; i < postMediaPoolSize; i++ {
		name := fmt.Sprintf("seed_post_%03d.png", i)
		if err := ensurePNG(filepath.Join(uploadDir, name), func() image.Image {
			a, b := postPalette(i)
			return gradientWithDots(postMediaW, postMediaH, a, b, int64(i)*7)
		}); err != nil {
			return nil, nil, err
		}
		media = append(media, imageURLPrefix+name)
	}

	return avatars, media, nil
}

func ensurePNG(path string, build func() image.Image) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	defer f.Close()
	if err := png.Encode(f, build()); err != nil {
		return fmt.Errorf("encode %s: %w", path, err)
	}
	return nil
}

// gradientWithDots paints a diagonal two-color gradient and sprinkles small
// white dots over it so each image is distinguishable at a glance.
func gradientWithDots(w, h int, c1, c2 color.RGBA, seed int64) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			t := float64(x+y) / float64(w+h)
			img.Set(x, y, lerp(c1, c2, t))
		}
	}
	rng := rand.New(rand.NewSource(seed))
	dots := (w * h) / 400
	for i := 0; i < dots; i++ {
		img.Set(rng.Intn(w), rng.Intn(h), color.RGBA{255, 255, 255, 220})
	}
	return img
}

func lerp(a, b color.RGBA, t float64) color.RGBA {
	return color.RGBA{
		R: uint8(float64(a.R)*(1-t) + float64(b.R)*t),
		G: uint8(float64(a.G)*(1-t) + float64(b.G)*t),
		B: uint8(float64(a.B)*(1-t) + float64(b.B)*t),
		A: 255,
	}
}

var avatarPalettes = [][2]color.RGBA{
	{{255, 99, 132, 255}, {255, 159, 64, 255}},
	{{54, 162, 235, 255}, {153, 102, 255, 255}},
	{{75, 192, 192, 255}, {45, 100, 100, 255}},
	{{255, 205, 86, 255}, {255, 99, 132, 255}},
	{{201, 203, 207, 255}, {99, 99, 99, 255}},
	{{34, 197, 94, 255}, {14, 116, 144, 255}},
	{{236, 72, 153, 255}, {99, 102, 241, 255}},
	{{251, 191, 36, 255}, {220, 38, 38, 255}},
	{{14, 165, 233, 255}, {16, 185, 129, 255}},
	{{168, 85, 247, 255}, {236, 72, 153, 255}},
	{{249, 115, 22, 255}, {239, 68, 68, 255}},
	{{20, 184, 166, 255}, {59, 130, 246, 255}},
}

var postPalettes = [][2]color.RGBA{
	{{15, 23, 42, 255}, {30, 64, 175, 255}},
	{{120, 53, 15, 255}, {251, 191, 36, 255}},
	{{6, 78, 59, 255}, {52, 211, 153, 255}},
	{{67, 56, 202, 255}, {236, 72, 153, 255}},
	{{15, 76, 117, 255}, {50, 130, 184, 255}},
	{{124, 45, 18, 255}, {251, 146, 60, 255}},
	{{20, 83, 45, 255}, {163, 230, 53, 255}},
}

func avatarPalette(i int) (color.RGBA, color.RGBA) {
	p := avatarPalettes[i%len(avatarPalettes)]
	return p[0], p[1]
}
func postPalette(i int) (color.RGBA, color.RGBA) {
	p := postPalettes[i%len(postPalettes)]
	return p[0], p[1]
}
