package yapay

type Data struct {
	ShopId int
	Scid   int
	Secret string
}

type InitObj struct {
	ShopId                  string
	Scid                    string
	OrderSumAmount          string
	OrderSumCurrencyPaycash string
	OrderSumBankPaycash     string
	InvoiceId               string
	CustomerNumber          string
	Md5                     string
	RequestDatetime         string
	Type                    string
}

type AnsObj struct {
	Message         string
	TechMessage     string
	Code            int
	InvoiceId       string
	ShopId          string
	RequestDatetime string
	Aviso           bool
}
