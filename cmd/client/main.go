package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/MXLange/desafio-pos-client-server-api/pkg/types"
)



const(
	callTimeout time.Duration = time.Duration(300) * time.Millisecond
	callTimeoutToFail time.Duration = time.Duration(2) * time.Millisecond
	apiAddress string = "http://localhost:8080/cotacao"
	apiCalls int = 5
)

func main() {

	log.Println("[CLIENT] Iniciando chamada(s), total: ", apiCalls)

	
	for i := range apiCalls {
		func(){

			ctx, cancel := context.WithTimeout(context.Background(), callTimeout)
			defer cancel()
			price, err := fetchPrice(ctx)
			if err != nil {
				log.Printf("[CLIENT] Chamada %d falhou: %s", i, err.Error())
				return
			}
			log.Printf("[CLIENT] Chamada %d preço atual é: %s", i, price.Bid)
		}()
	}

	log.Printf("[CLIENT] Iniciando chamada com timeout curto para falhar, total: %d", apiCalls)
	
	for i := range apiCalls {
		func(){
			ctx, cancel := context.WithTimeout(context.Background(), callTimeoutToFail)
			defer cancel()
			price, err := fetchPrice(ctx)
			if err != nil {
				log.Printf("[CLIENT] Chamada %d falhou como esperado: %s", i, err.Error())
				return
			}
			log.Printf("[CLIENT] Chamada %d preço atual é: %s", i, price.Bid)
		}()
	}

}


func fetchPrice(ctx context.Context) (*types.Bid, error ) {
	
	req, err := http.NewRequestWithContext(ctx, "GET", apiAddress, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var bid types.Bid
	if err := json.Unmarshal(body, &bid); err != nil {
		return nil, err
	}

	return &bid, nil	

}
