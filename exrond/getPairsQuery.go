package exrond

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

type Pairs struct {
	Data struct {
		Pairs []struct {
			Address            string `json:"address"`
			FirstTokenPrice    string `json:"firstTokenPrice"`
			SecondTokenPrice   string `json:"secondTokenPrice"`
			FirstTokenReserve  string `json:"firstTokenReserve"`
			SecondTokenReserve string `json:"secondTokenReserve"`
			FirstToken         struct {
				Identifier string `json:"identifier"`
				Name       string `json:"name"`
				Decimals   int    `json:"decimals"`
				Supply     string `json:"supply"`
				Ticker     string `json:"ticker"`
				Balance    string `json:"balance"`
			} `json:"firstToken"`
			SecondToken struct {
				Identifier string `json:"identifier"`
				Name       string `json:"name"`
				Decimals   int    `json:"decimals"`
				Supply     string `json:"supply"`
				Ticker     string `json:"ticker"`
				Balance    string `json:"balance"`
			} `json:"secondToken"`
		} `json:"pairs"`
	} `json:"data"`
}

func GetPairsQuery() (Pairs, error) {
	//Encode the data
	var (
		url = "https://api.exrond.com/graphql"
	)

	request := `
	{"operationName":"GET_PAIRS_QUERY","variables":{},"query":"query GET_PAIRS_QUERY {\n  pairs {\n    address\n    firstTokenPrice\n    secondTokenPrice\n    firstTokenReserve\n    secondTokenReserve\n    firstToken {\n      identifier\n      name\n      decimals\n      supply\n      ticker\n      balance\n      assets {\n        svgUrl\n        pngUrl\n        __typename\n      }\n      __typename\n    }\n    secondToken {\n      identifier\n      name\n      decimals\n      supply\n      ticker\n      balance\n      assets {\n        svgUrl\n        pngUrl\n        __typename\n      }\n      __typename\n    }\n    __typename\n  }\n}"}
	`

	// Create a new request using http
	resp, err := http.Post(url, "application/json", strings.NewReader(request))
	if err != nil {
		log.Err(err).Msg("error making request")
		return Pairs{}, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Err(err).Msg("error reading response body")
		return Pairs{}, err
	}

	var p Pairs
	err = json.Unmarshal(body, &p)
	if err != nil {
		log.Err(err).Msg("error unmarshalling data")
		return Pairs{}, err
	}

	return p, nil
}
