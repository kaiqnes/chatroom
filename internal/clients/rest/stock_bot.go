package rest

import (
	"chatroom/internal/config"
	"chatroom/internal/domain"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
)

type stockBot struct {
	UrlTemplate string
}

func NewStockBot(cfg *config.Config) domain.StockBotClient {
	return &stockBot{
		UrlTemplate: cfg.StockBotTemplateURL,
	}
}

func (c *stockBot) Call(req domain.StockBotRequest) (*domain.StockBotResponseDto, error) {
	var response *domain.StockBotResponseDto

	resp, err := http.Get(c.ParseURL(req.StockCode))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	isHeader := true
	for {
		row, err := reader.Read()
		if err == io.EOF || err != nil {
			return nil, err
		}
		if isHeader {
			isHeader = false
			continue
		}
		response = &domain.StockBotResponseDto{
			Symbol: row[0],
			Date:   row[1],
			Time:   row[2],
			Open:   row[3],
			High:   row[4],
			Low:    row[5],
			Close:  row[6],
			Volume: row[7],
		}
		break
	}
	return response, nil
}

func (c *stockBot) ParseURL(stock string) string {
	return fmt.Sprintf(c.UrlTemplate, stock)
}
