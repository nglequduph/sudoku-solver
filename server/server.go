package server

import (
	"fmt"
	"html/template"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"strconv"

	"example.com/sdk/ocr"
	"example.com/sdk/solver"
)

type PageData struct {
	Original   [9][9]int
	Grid       [9][9]int
	Solved     bool
	Done       bool
	Steps      int
	Error      string
	IsOriginal func(r, c int) bool
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		return
	}
	t.Execute(w, PageData{}) // Pass empty data to render initial grid
}

func SolveHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "ParseForm error: "+err.Error(), http.StatusBadRequest)
		return
	}

	var inputGrid [9][9]int
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			fieldName := fmt.Sprintf("cell-%d-%d", row, col)
			valStr := r.FormValue(fieldName)
			if valStr != "" {
				if val, err := strconv.Atoi(valStr); err == nil && val >= 1 && val <= 9 {
					inputGrid[row][col] = val
				}
			}
		}
	}

	log.Println("Solving Sudoku from manual input...")
	solvedGrid, success, steps := solver.SolveSudoku(inputGrid)

	data := PageData{
		Original: inputGrid,
		Grid:     solvedGrid,
		Solved:   success,
		Done:     true,
		Steps:    steps,
		IsOriginal: func(r, c int) bool {
			return inputGrid[r][c] != 0
		},
	}

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

func OCRHandler(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		renderError(w, "Could not decode image: "+err.Error())
		return
	}

	log.Println("Starting OCR extraction...")
	inputGrid, err := ocr.ExtractGridFromImage(img)
	if err != nil {
		renderError(w, "OCR Failed: "+err.Error())
		return
	}
	log.Printf("OCR Complete. Extracted Grid: %+v\n", inputGrid)

	// Pre-fill the grid with OCR results for user review
	data := PageData{
		Grid: inputGrid, // Populate the inputs
		// We could mark them as "Original" if we want style,
		// but since the user needs to edit/correct them, plain inputs are fine.
		// If we want bold style for detected numbers:
		IsOriginal: func(r, c int) bool {
			return inputGrid[r][c] != 0
		},
	}

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

func renderError(w http.ResponseWriter, msg string) {
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	t.Execute(w, PageData{Error: msg})
}

func Start(port string) {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/solve", SolveHandler)
	http.HandleFunc("/ocr", OCRHandler)

	log.Printf("Server running on http://localhost:%s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
