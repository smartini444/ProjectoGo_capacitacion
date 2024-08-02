package cotizadores

import (
	"encoding/json"
	"fmt"
	"net/http"
	criptomonedas "primerProjecto/internal/entities/criptomonedas"
	"time"
)

type CoinPaprikaCotizador struct{}

type CoinpaprikaResponse struct {
	Name   string `json:"name"`
	Quotes map[string]struct {
		Price float64 `json:"price"`
	} `json:"quotes"`
}

func (s *CoinPaprikaCotizador) GetCotizacionExterna(moneda, codigo, fiat string) (criptomonedas.Cotizacion, error) {
	// Paso 1: Buscar el ID de la criptomoneda en CoinPaprika
	coinListURL := "https://api.coinpaprika.com/v1/coins"
	var cotizacion criptomonedas.Cotizacion
	resp, err := http.Get(coinListURL)
	if err != nil {
		return cotizacion, fmt.Errorf("error al obtener la lista de monedas: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return cotizacion, fmt.Errorf("error en la solicitud de lista de monedas: %s", resp.Status)
	}

	var coins []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&coins); err != nil {
		return cotizacion, fmt.Errorf("error al decodificar la lista de monedas: %v", err)
	}

	var coinID string
	for _, coin := range coins {
		if coin.Name == moneda {
			coinID = coin.ID
			break
		}
	}

	if coinID == "" {
		return cotizacion, fmt.Errorf("no se encontró la criptomoneda %s en CoinPaprika", moneda)
	}

	// Paso 2: Usar el ID para obtener la cotización más reciente
	tickerURL := fmt.Sprintf("https://api.coinpaprika.com/v1/tickers/%s", coinID)
	resp, err = http.Get(tickerURL)
	if err != nil {
		return cotizacion, fmt.Errorf("error al obtener la cotización: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return cotizacion, fmt.Errorf("error en la solicitud de cotización: %s", resp.Status)
	}

	var result CoinpaprikaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return cotizacion, fmt.Errorf("error al decodificar la respuesta de cotización: %v", err)
	}

	quote, ok := result.Quotes[fiat]
	if !ok {
		return cotizacion, fmt.Errorf("no se encontró la cotización para la moneda fiat %s", fiat)
	}

	/* ESTO LO TIENE QUE HACER EL SERVICE AHORA.
	// Paso 3: Buscar la criptomoneda por nombre en tu base de datos
	cripto, err := s.FindCryptoByName(moneda)
	if err != nil {
		return cotizacion, err
	}

	if cripto == nil {
		return cotizacion, fmt.Errorf("la criptomoneda %s no está registrada en la base de datos", moneda)
	}*/

	cotizacion = criptomonedas.Cotizacion{ //service tiene que buscar el id de la cripto para esto
		Cotizacion: quote.Price,
		Fecha:      time.Now(),
	}

	return cotizacion, nil
}
