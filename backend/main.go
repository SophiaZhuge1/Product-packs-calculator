package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/rs/cors"
)

type Pack struct {
	Size int `json:"size"`
}

type OrderRequest struct {
	Items int `json:"items"`
}

type OrderResponse struct {
	Packs []Pack `json:"packs"`
}

var packSizes = []int{250, 500, 1000, 2000, 5000}

func getPacks(items int) []Pack {
	sort.Sort(sort.Reverse(sort.IntSlice(packSizes)))

	packs := []Pack{}
	for _, size := range packSizes {
		count := items / size
		items = items % size
		for i := 0; i < count; i++ {
			packs = append(packs, Pack{Size: size})
		}
		fmt.Printf("Remaining items: %d, Added packs: %+v\n", items, packs)
	}
	if items > 0 {
		packs = append(packs, Pack{Size: packSizes[len(packSizes)-1]})
		fmt.Printf("Remaining items added with smallest pack: %d, Added packs: %+v\n", items, packs)
	}
	return packs
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
	var orderReq OrderRequest
	err := json.NewDecoder(r.Body).Decode(&orderReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	packs := getPacks(orderReq.Items)
	orderRes := OrderResponse{Packs: packs}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderRes)
}

func main() {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://packs-calculator-frontend.vercel.app"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		Debug:            true,
	})

	handler := c.Handler(http.DefaultServeMux)

	http.HandleFunc("/order", orderHandler)
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
