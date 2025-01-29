package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// CurrencyConversionRequest định nghĩa request đầu vào cho Lambda
type CurrencyConversionRequest struct {
	Amount float64 `json:"amount"`
	From   string  `json:"from"`
	To     string  `json:"to"`
}

// CurrencyConversionResponse định nghĩa response của Lambda
type CurrencyConversionResponse struct {
	ConvertedAmount float64 `json:"converted_amount"`
	From            string  `json:"from"`
	To              string  `json:"to"`
	Rate            float64 `json:"rate"`
}

// Fake exchange rates (giả định)
var exchangeRates = map[string]map[string]float64{
	"USD": {"EUR": 0.85, "JPY": 110.0, "VND": 23000.0},
	"EUR": {"USD": 1.18, "JPY": 130.0, "VND": 27000.0},
	"JPY": {"USD": 0.0091, "EUR": 0.0077, "VND": 210.0},
	"VND": {"USD": 0.000043, "EUR": 0.000037, "JPY": 0.0048},
}

// convertCurrency thực hiện chuyển đổi tiền tệ
func convertCurrency(amount float64, from string, to string) (float64, float64, error) {
	if rates, exists := exchangeRates[from]; exists {
		if rate, exists := rates[to]; exists {
			return amount * rate, rate, nil
		}
	}
	return 0, 0, fmt.Errorf("hông tìm thấy tỷ giá từ %s sang %s", from, to)
}

// Handler là hàm chính xử lý Lambda
func Handler(ctx context.Context, request CurrencyConversionRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Nhận request: %+v\n", request)

	convertedAmount, rate, err := convertCurrency(request.Amount, request.From, request.To)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	response := CurrencyConversionResponse{
		ConvertedAmount: convertedAmount,
		From:            request.From,
		To:              request.To,
		Rate:            rate,
	}
	responseBody, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"error": "Internal Server Error"}`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(responseBody),
	}, nil
}

func main() {
	lambda.Start(Handler)
}
