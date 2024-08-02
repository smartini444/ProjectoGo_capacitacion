package cotizadores

import (
	"encoding/json"
	"fmt"
	"net/http"
	criptomonedas "primerProjecto/internal/entities/criptomonedas"
	"time"
)

type Exchange struct { //
	Ask      float64 `json:"ask"`
	TotalAsk float64 `json:"totalAsk"`
	Bid      float64 `json:"bid"`
	TotalBid float64 `json:"totalBid"`
	Time     int64   `json:"time"`
}

type CryptoYaQueryResponse struct { //
	SatoshiTango Exchange `json:"satoshitango"`
	LetsBit      Exchange `json:"letsbit"`
	BinanceP2P   Exchange `json:"binancep2p"`
	FiWind       Exchange `json:"fiwind"`
	TiendaCrypto Exchange `json:"tiendacrypto"`
	Calypso      Exchange `json:"calypso"`
	BanexCoin    Exchange `json:"banexcoin"`
	BitsoAlpha   Exchange `json:"bitsoalpha"`
	X4T          Exchange `json:"x4t"`
}

type CryptoYaCotizador struct{}

func (s *CryptoYaCotizador) GetCotizacionExterna(moneda, codigo, fiat string) (criptomonedas.Cotizacion, error) {

	volumen := 0.1
	var cotizacion criptomonedas.Cotizacion

	// Construir la URL del endpoint
	url := fmt.Sprintf("https://criptoya.com/api/%s/%s/%.2f", codigo, fiat, volumen)

	resp, err := http.Get(url)
	if err != nil {
		return cotizacion, fmt.Errorf("error al realizar la solicitud HTTP: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return cotizacion, fmt.Errorf("error en la solicitud: %s", resp.Status)
	}

	var apiResponse CryptoYaQueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return cotizacion, fmt.Errorf("error al decodificar la respuesta JSON: %v", err)
	}

	// usando el primer exchange
	price := apiResponse.SatoshiTango.Ask

	/* ESTO LO TIENE QUE HACER EL SERVICE AHORA.
	// Buscar la criptomoneda por nombre
	cripto, err := s.repo.FindCryptoByName(moneda)
	if err != nil {
		return cotizacion, err
	}

	if cripto == nil {
		return cotizacion, fmt.Errorf("la criptomoneda %s no está registrada en la base de datos", moneda)
	}
	*/
	cotizacion = criptomonedas.Cotizacion{
		Cotizacion: price,
		Fecha:      time.Now(),
	}

	return cotizacion, nil
}
