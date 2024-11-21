package getaway

type Metrics interface {
	IncPayProductSumCounter(status string)
	AddAmountPayProduct(amount float64)
}
