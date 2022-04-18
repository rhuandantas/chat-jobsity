package client

import (
	"encoding/csv"
	"fmt"
	"github.com/rhuandantas/chat-jobsity/internal/config"
	"github.com/rhuandantas/chat-jobsity/internal/model"
	"net/http"
	"strconv"
)

type Stock struct {
	cfg config.Config
}

func NewStooqClient(cfg config.Config) *Stock {
	return &Stock{
		cfg: cfg,
	}
}

func (c *Stock) GetStock(stockCode string) (*model.Stock, error) {
	resp, err := http.Get(fmt.Sprintf(c.cfg.StooqClient.UrlTempate, stockCode))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	records, err := csv.NewReader(resp.Body).ReadAll()
	if err != nil {
		return nil, err
	}

	data := records[1]
	volume, _ := strconv.Atoi(data[7])
	stock := model.Stock{
		Symbol: data[0],
		//TODO format date and time
		Date:   data[1],
		Time:   data[2],
		Open:   parseFloat(data[3]),
		High:   parseFloat(data[4]),
		Low:    parseFloat(data[5]),
		Close:  parseFloat(data[6]),
		Volume: volume,
	}

	return &stock, nil
}

func parseFloat(s string) float64 {
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return 0.0
}
