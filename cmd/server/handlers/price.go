package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	apperrors "github.com/MXLange/desafio-pos-client-server-api/cmd/server/app_errors"
	"github.com/MXLange/desafio-pos-client-server-api/cmd/server/repository"

	"github.com/MXLange/desafio-pos-client-server-api/pkg/types"
)

// PriceHandler handles HTTP requests related to price data.
type PriceHandler struct {
	db *repository.PriceRepository
	callTimeout time.Duration
	dbTimeout time.Duration
	priceUrl string
}

// NewPriceHandler creates a new instance of PriceHandler.
func NewPriceHandler(db *repository.PriceRepository, callTimeout, dbTimeout time.Duration, priceUrl string) (*PriceHandler, error) {
	if db == nil {
		return nil, apperrors.ErrNilPriceRepository
	}
	return &PriceHandler{db: db, callTimeout: callTimeout, dbTimeout: dbTimeout, priceUrl: priceUrl}, nil
}

// GetPrice handles the HTTP request to fetch the current price, store it in the database, and return the bid value.
func (h *PriceHandler) GetPrice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	callCtx, cancel := context.WithTimeout(ctx, h.callTimeout)
	defer cancel()
	price, err := h.callExternalAPI(callCtx)
	if err != nil {
		log.Println("[SERVER] Error calling external API: ", err.Error())
		http.Error(w, "Failed to fetch price from external API", http.StatusInternalServerError)
		return
	}
	

	dbCtx, cancel := context.WithTimeout(ctx, h.dbTimeout)
	defer cancel()

	_, err = h.db.InsertPrice(dbCtx, price)
	if err != nil {
		log.Println("[SERVER] Error inserting price into database: ", err.Error())
		http.Error(w, "Failed to store price in database", http.StatusInternalServerError)
		return
	}

	res := types.Bid{Bid: price.USDBRL.Bid}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println("[SERVER] Error encoding response: ", err.Error())
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// callExternalAPI makes an HTTP GET request to the external price API and returns the parsed Price data.
func (h *PriceHandler) callExternalAPI(ctx context.Context) (*types.Price, error) {

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", h.priceUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var price types.Price
	if err := json.Unmarshal(respBody, &price); err != nil {
		return nil, err
	}

	return &price, nil
}
