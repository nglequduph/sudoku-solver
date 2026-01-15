package ocr

import (
	"bytes"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/otiai10/gosseract/v2"
)

func ExtractGridFromImage(img image.Image, client *gosseract.Client) ([9][9]int, error) {
	var grid [9][9]int

	// Preprocessing:
	// 1. Grayscale
	// 2. Resize to fixed 900x900
	// 3. Simple Binarize (Thresholding) to remove shadows/noise
	// Thresholding at 160 works well for clean images, but lowering it to 140
	// might help if numbers are thin/faint.
	procImg := imaging.Grayscale(img)
	procImg = imaging.Resize(procImg, 900, 900, imaging.Lanczos)
	procImg = binarize(procImg, 140) // Reduced threshold to keep lines thicker

	// Dilation: Make text bolder to help Tesseract
	// Since imaging/v1 doesn't have Dilate/Erode, we can simulate or skip.
	// Often Tesseract prefers clearer text.

	cellWidth := 100
	cellHeight := 100
	// Increase margin further to 22
	margin := 22

	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			// Crop individual cell
			rect := image.Rect(
				c*cellWidth+margin,
				r*cellHeight+margin,
				(c+1)*cellWidth-margin,
				(r+1)*cellHeight-margin,
			)
			cellImg := imaging.Crop(procImg, rect)

			// Simple check: if cell is mostly empty, skip OCR
			if isMostlyEmpty(cellImg) {
				grid[r][c] = 0
				continue
			}

			// Pass to OCR
			buf := new(bytes.Buffer)
			if err := imaging.Encode(buf, cellImg, imaging.PNG); err != nil {
				return grid, err
			}

			client.SetImageFromBytes(buf.Bytes())
			client.SetPageSegMode(gosseract.PSM_SINGLE_CHAR)
			client.SetWhitelist("123456789")

			text, err := client.Text()
			if err != nil {
				return grid, err
			}

			text = strings.TrimSpace(text)
			if text != "" {
				val, err := strconv.Atoi(text)
				if err == nil && val >= 1 && val <= 9 {
					grid[r][c] = val
				}
			}
		}
	}
	return grid, nil
}

func isMostlyEmpty(img *image.NRGBA) bool {
	darkPixels := 0
	totalPixels := img.Bounds().Dx() * img.Bounds().Dy()
	threshold := 150

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			c := img.At(x, y).(color.NRGBA)
			if int(c.R) < threshold {
				darkPixels++
			}
		}
	}

	percentage := float64(darkPixels) / float64(totalPixels)
	return percentage < 0.02 // Lower threshold for "mostly empty"
}

// Manually binarize image to Black/White
func binarize(img *image.NRGBA, threshold uint8) *image.NRGBA {
	bounds := img.Bounds()
	dst := image.NewNRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y).(color.NRGBA)
			// Since it's grayscale, R=G=B. Check R.
			if c.R < threshold {
				dst.SetNRGBA(x, y, color.NRGBA{0, 0, 0, 255}) // Black
			} else {
				dst.SetNRGBA(x, y, color.NRGBA{255, 255, 255, 255}) // White
			}
		}
	}
	return dst
}
