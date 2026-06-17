package trading

func CalculateFee(amount, price, quantity float64, feeBPS int) float64 {
	return amount * float64(feeBPS) / 10000
}

func CalculateMakerFee(amount float64, makerFeeBPS int) float64 {
	return CalculateFee(amount, 0, 0, makerFeeBPS)
}

func CalculateTakerFee(amount float64, takerFeeBPS int) float64 {
	return CalculateFee(amount, 0, 0, takerFeeBPS)
}

func IsMakerOrder(timeInForce string) bool {
	return timeInForce == "GTC"
}
