package cotizadores

//go:generate echo $GOPACKAGE/$GOFILE
//go:generate mockgen -source=./$GOFILE -destination=./mock/$GOFILE -package mock

import (
	"fmt"
	criptomonedas "primerProjecto/internal/entities/criptomonedas"
)

type Cotizador interface {
	GetCotizacionExterna(moneda, codigo, fiat string) (criptomonedas.Cotizacion, error)
}

var CotizadoresMap = map[string]Cotizador{
	"coinpaprika": &CoinPaprikaCotizador{},
	"criptoya":    &CryptoYaCotizador{},
	// Agrega otros cotizadores aquí, como Cryptoya.
}

func GetCotizador(name string) (Cotizador, error) {
	cotizador, exists := CotizadoresMap[name]
	if !exists {
		return nil, fmt.Errorf("cotizador %s no soportado", name)
	}
	return cotizador, nil
}
