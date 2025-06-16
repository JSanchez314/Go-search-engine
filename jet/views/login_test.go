package main

import (
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestLocalHTML(t *testing.T) {

	file, err := os.Open("login.hbs")
	if err != nil {
		t.Fatalf("Error al abrir el archivo: %v", err)
	}
	defer file.Close()

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		t.Fatalf("Error al analizar HTML: %v", err)
	}
	if doc.Find(`form`).Length() == 0 {
		t.Error("expected form to be render, but it wasn´t")
	}
}
