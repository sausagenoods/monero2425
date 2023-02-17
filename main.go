package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"net/http"

	"github.com/mb-14/gomarkov"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	train := flag.Bool("train", false, "Train the markov chain")
	bind := flag.String("bind", ":3000", "Bind address")
	flag.Parse()
	if *train {
		chain, err := buildModel()
		if err != nil {
			fmt.Println(err)
			return
		}
		saveModel(chain)
	} else {
		chain, err := loadModel()
		if err != nil {
			fmt.Println(err)
			return
		}
		router(bind, chain)
	}
}

func router(bind *string, chain *gomarkov.Chain) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, generateLedgerSpam(chain))
	})
	http.ListenAndServe(*bind, r)
}

func buildModel() (*gomarkov.Chain, error) {
	spam, err := readLedgerSpam()
	if err != nil {
		return nil, err
	}
	chain := gomarkov.NewChain(1)
	fmt.Println("Adding 24/25 ledger spam to markov chain...")
	for _, s := range spam {
		chain.Add(strings.Split(s, " "))
	}
	return chain, nil
}

func loadModel() (*gomarkov.Chain, error) {
	var chain gomarkov.Chain
	data, err := os.ReadFile("model.json")
	if err != nil {
		return &chain, err
	}
	err = json.Unmarshal(data, &chain)
	if err != nil {
		return &chain, err
	}
	return &chain, nil
}

func readLedgerSpam() ([]string, error) {
	var spam []string
	data, err := os.ReadFile("2425.json")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &spam)
	return spam, err
}

func saveModel(chain *gomarkov.Chain) {
	jsonObj, _ := json.Marshal(chain)
	err := os.WriteFile("model.json", jsonObj, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func generateLedgerSpam(chain *gomarkov.Chain) string {
	tokens := []string{gomarkov.StartToken}
	for tokens[len(tokens)-1] != gomarkov.EndToken {
		next, _ := chain.Generate(tokens[(len(tokens) - 1):])
		tokens = append(tokens, next)
	}
	return strings.Join(tokens[1:len(tokens)-1], " ")
}
