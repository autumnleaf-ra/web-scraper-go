package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type PokemonProduct struct {
	url, image, name, price string
}

var pokemonProducts []PokemonProduct

func main() {

	c := colly.NewCollector()

	// Request Page
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	// Error
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Something went wrong: ", err)
	})

	// On-Response
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Status Code: ", r.StatusCode)
	})

	// Request HTML - Product
	c.OnHTML("li.product", func(e *colly.HTMLElement) {
		//  initializing a new PokemonProduct instance
		pokemonProduct := PokemonProduct{}

		// scraping the data
		pokemonProduct.url = e.ChildAttr("a", "href")
		pokemonProduct.image = e.ChildAttr("img", "src")
		pokemonProduct.name = e.ChildText("h2")
		pokemonProduct.price = e.ChildText(".price")

		// adding the product instance with scraped data to the list of products
		pokemonProducts = append(pokemonProducts, pokemonProduct)
	})

	// Visit
	c.Visit("https://scrapeme.live/shop/")

	// Convert To CSV File

	file, err := os.Create("products.csv")
	if err != nil {
		log.Fatalln("Failed to create output CSV file", err)
	}
	defer file.Close()

	// initializing a file writer
	wr := csv.NewWriter(file)

	// define CSV headers
	headers := []string{
		"url",
		"image",
		"name",
		"price",
	}

	// writing the column headers
	wr.Write(headers)

	// adding Pokemon product to the CSV output file
	for _, pokemonProduct := range pokemonProducts {
		// convert a PokemonProduct to an array of strings
		record := []string{
			pokemonProduct.url,
			pokemonProduct.image,
			pokemonProduct.name,
			pokemonProduct.price,
		}

		// writing a new CSV record
		wr.Write(record)
	}
	wr.Flush()
}
