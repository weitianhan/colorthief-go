package colorthief

import (
	"errors"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"

	"github.com/weitianhan/colorthief-go/mediancut"
)

var DefaultMaxCubes = 6

// GetColorFromFile return the base color from the image file
func GetColorFromFile(imgPath string) (color.Color, error) {
	colors, _, err := GetPaletteFromFile(imgPath, DefaultMaxCubes)
	if err != nil {
		return color.RGBA{}, nil
	}
	return colors[0], nil
}

// GetColor return the base color from the image
func GetColor(img image.Image) (color.Color, error) {
	colors, _, err := GetPalette(img, DefaultMaxCubes)
	if err != nil {
		return color.RGBA{}, nil
	}
	return colors[0], nil
}

// GetPaletteFromFile return cluster similar colors from the image file
func GetPaletteFromFile(imgPath string, maxCubes int) ([]color.Color, []float64, error) {
	f, err := os.Open(imgPath)
	if err != nil {
		return nil, nil, err
	}

	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, nil, err
	}

	bounds := img.Bounds()
	totalPixels := (bounds.Max.X - bounds.Min.X) * (bounds.Max.Y - bounds.Min.Y)
	colors, pixNum, err := GetPalette(img, maxCubes)
	ratio := make([]float64, len(pixNum))
	for i := range pixNum {
		ratio[i] = float64(pixNum[i]) / float64(totalPixels)
	}
	return colors, ratio, err
}

// GetPalette return cluster similar colors by the median cut algorithm
func GetPalette(img image.Image, maxCubes int) ([]color.Color, []int, error) {
	return mediancut.GetPalette(img, maxCubes)
}

func PrintColor(colors []color.Color, filename string) error {
	imgWidth := 100 * len(colors)
	imgHeight := 200
	if imgWidth == 0 {
		return errors.New("colors empty")
	}

	paletted := image.NewPaletted(image.Rect(0, 0, imgWidth, imgHeight), colors)

	for x := 0; x < imgWidth; x++ {
		idx := x / 100
		for y := 0; y < imgHeight; y++ {
			paletted.SetColorIndex(x, y, uint8(idx))
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, paletted)
}
