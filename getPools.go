package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Pool struct {
	ID string `json:"id"`
}

type PoolResponse struct {
	Pools []Pool `json:"pools"`
}

type Data struct {
	Data PoolResponse `json:"data"`
}

type Pools struct {
	Pools []string
}

func getUniswapV3Pools() (Pools, error) {

	s := struct {
		Query string `json:"query"`
	}{"{pools(where:{txCount_gt:50000}){id}}"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(s)

	resp, err := http.Post("https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v3", "application/json", b)

	if err != nil {
		return Pools{}, err
	}

	defer resp.Body.Close()

	tmp := Data{}
	err = json.NewDecoder(resp.Body).Decode(&tmp)

	if err != nil {
		return Pools{}, err
	}

	rtnPools := Pools{}
	for _, pool := range tmp.Data.Pools {
		rtnPools.Pools = append(rtnPools.Pools, pool.ID)
	}
	return rtnPools, nil
}
