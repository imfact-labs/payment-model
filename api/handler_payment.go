package api

import (
	"github.com/imfact-labs/payment-model/digest"
	"net/http"

	apic "github.com/imfact-labs/currency-model/api"
	ctypes "github.com/imfact-labs/currency-model/types"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/payment-model/types"
)

var (
	HandlerPathPaymentDesign      = `/payment/{contract:(?i)` + ctypes.REStringAddressString + `}`
	HandlerPathPaymentAccountInfo = `/payment/{contract:(?i)` + ctypes.REStringAddressString + `}/account/{address:(?i)` + ctypes.REStringAddressString + `}`
)

func SetHandlers(hd *apic.Handlers) {
	get := 1000
	_ = hd.SetHandler(HandlerPathPaymentAccountInfo, HandlePaymentAccountInfo, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.SetHandler(HandlerPathPaymentDesign, HandlePaymentDesign, true, get, get).
		Methods(http.MethodOptions, "GET")
}

func HandlePaymentDesign(hd *apic.Handlers, w http.ResponseWriter, r *http.Request) {
	cacheKey := apic.CacheKeyPath(r)
	if err := apic.LoadFromCache(hd.Cache(), cacheKey, w); err == nil {
		return
	}

	contract, err, status := apic.ParseRequest(w, r, "contract")
	if err != nil {
		apic.HTTP2ProblemWithError(w, err, status)

		return
	}

	if v, err, shared := hd.RG().Do(cacheKey, func() (interface{}, error) {
		return handlePaymentDesignInGroup(hd, contract)
	}); err != nil {
		apic.HTTP2HandleError(w, err)
	} else {
		apic.HTTP2WriteHalBytes(hd.Encoder(), w, v.([]byte), http.StatusOK)

		if !shared {
			apic.HTTP2WriteCache(w, cacheKey, hd.ExpireShortLived())
		}
	}
}

func handlePaymentDesignInGroup(hd *apic.Handlers, contract string) ([]byte, error) {
	var de *types.Design
	var st base.State

	de, st, err := digest.PaymentDesign(hd.Database(), contract)
	if err != nil {
		return nil, err
	}

	i, err := buildPaymentDesign(hd, contract, *de, st)
	if err != nil {
		return nil, err
	}
	return hd.Encoder().Marshal(i)
}

func buildPaymentDesign(hd *apic.Handlers, contract string, de types.Design, st base.State) (apic.Hal, error) {
	h, err := hd.CombineURL(HandlerPathPaymentDesign, "contract", contract)
	if err != nil {
		return nil, err
	}

	var hal apic.Hal
	hal = apic.NewBaseHal(de, apic.NewHalLink(h, nil))

	h, err = hd.CombineURL(apic.HandlerPathBlockByHeight, "height", st.Height().String())
	if err != nil {
		return nil, err
	}
	hal = hal.AddLink("block", apic.NewHalLink(h, nil))

	for i := range st.Operations() {
		h, err := hd.CombineURL(apic.HandlerPathOperation, "hash", st.Operations()[i].String())
		if err != nil {
			return nil, err
		}
		hal = hal.AddLink("operations", apic.NewHalLink(h, nil))
	}

	return hal, nil
}

func HandlePaymentAccountInfo(hd *apic.Handlers, w http.ResponseWriter, r *http.Request) {
	cachekey := apic.CacheKeyPath(r)
	if err := apic.LoadFromCache(hd.Cache(), cachekey, w); err == nil {
		return
	}

	contract, err, status := apic.ParseRequest(w, r, "contract")
	if err != nil {
		apic.HTTP2ProblemWithError(w, err, status)

		return
	}

	account, err, status := apic.ParseRequest(w, r, "address")
	if err != nil {
		apic.HTTP2ProblemWithError(w, err, status)

		return
	}

	if v, err, shared := hd.RG().Do(cachekey, func() (interface{}, error) {
		return handlePaymentAccountInfoInGroup(hd, contract, account)
	}); err != nil {
		apic.HTTP2HandleError(w, err)
	} else {
		apic.HTTP2WriteHalBytes(hd.Encoder(), w, v.([]byte), http.StatusOK)

		if !shared {
			apic.HTTP2WriteCache(w, cachekey, hd.ExpireShortLived())
		}
	}
}

func handlePaymentAccountInfoInGroup(hd *apic.Handlers, contract, account string) ([]byte, error) {
	var accountInfoValue *digest.AccountInfoValue

	accountInfoValue, err := digest.AccountInfo(hd.Database(), contract, account)
	if err != nil {
		return nil, err
	}

	i, err := buildAccountInfoValue(hd, contract, *accountInfoValue)
	if err != nil {
		return nil, err
	}
	return hd.Encoder().Marshal(i)
}

func buildAccountInfoValue(hd *apic.Handlers, contract string, it digest.AccountInfoValue) (apic.Hal, error) {
	h, err := hd.CombineURL(
		HandlerPathPaymentAccountInfo,
		"contract", contract, "address", it.AccountInfo().Address().String(),
	)
	if err != nil {
		return nil, err
	}

	var hal apic.Hal
	hal = apic.NewBaseHal(it, apic.NewHalLink(h, nil))

	return hal, nil
}
