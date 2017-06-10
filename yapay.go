package yapay

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Проверка корректности оплаты
func PayCorrectCheck(o InitObj, d Data) (ans AnsObj) {

	ans = AnsObj{
		RequestDatetime: o.RequestDatetime,
		InvoiceId:       o.InvoiceId,
		ShopId:          o.ShopId,
	}

	// Проверяем shopId
	shopId, err := strconv.ParseInt(o.ShopId, 10, 32)
	if err != nil {
		ans.Message = "Указан неверный магазин"
		ans.TechMessage = "Указан некоректный shopId"
		ans.Code = 100
		return
	} else if int(shopId) != d.ShopId {
		ans.Message = "Указан неверный магазин"
		ans.TechMessage = "Указан неверный shopId"
		ans.Code = 100
		return
	}

	// Проверяем scid
	scid, err := strconv.ParseInt(o.Scid, 10, 32)
	if err != nil {
		ans.Message = "Указан неверный магазин"
		ans.TechMessage = "Указан некоректный shopId"
		ans.Code = 100
		return
	} else if int(scid) != d.Scid {
		ans.Message = "Указан неверный магазин"
		ans.TechMessage = "Указан неверный scid"
		ans.Code = 100
		return
	}

	// Формируем подпись
	crcstr := strings.Join([]string{o.Type, o.OrderSumAmount, o.OrderSumCurrencyPaycash, o.OrderSumBankPaycash,
		o.ShopId, o.InvoiceId, o.CustomerNumber, d.Secret}, ";")

	crc := strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(crcstr))))
	if crc != o.Md5 {
		log.Println("[error]", crc, o.Md5, crcstr)
		ans.Message = "Ошибка обработки запроса."
		ans.TechMessage = "Несовпадение значения параметра md5 с результатом расчета хэш-функции."
		ans.Code = 100
		return
	}

	return
}

// Возвращаем ответ кассе
func SendAnswer(w http.ResponseWriter, ans AnsObj) {
	// Определяем тип ответа
	var t string
	if ans.Aviso {
		t = "paymentAvisoResponse"
	} else {
		t = "checkOrderResponse"
	}

	// Формируем строку ответа
	str := `<?xml version="1.0" encoding="utf-8"?>
`

	if ans.Message == "" {
		str += fmt.Sprintf(`<%s performedDatetime="%s" code="%d" invoiceId="%s" shopId="%s"/>`,
			t, ans.RequestDatetime, ans.Code, ans.InvoiceId, ans.ShopId)
	} else {
		str += fmt.Sprintf(`<%s performedDatetime="%s" code="%d" invoiceId="%s" shopId="%s" message="%s" techMessage="%s"/>`,
			t, ans.RequestDatetime, ans.Code, ans.InvoiceId, ans.ShopId, ans.Message, ans.TechMessage)
	}

	// Пишем ответ
	log.Println("[info]", str)
	//fmt.Fprint(ro.W, str)
	_, err := w.Write([]byte(str))
	if err != nil {
		log.Println("[error]", err)
	}
}
