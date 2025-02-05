package digest

import (
	"net/http"
	"time"

	cdigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-payment/types"
	"github.com/ProtoconNet/mitum2/base"
)

func (hd *Handlers) handlePaymentDesign(w http.ResponseWriter, r *http.Request) {
	cacheKey := cdigest.CacheKeyPath(r)
	if err := cdigest.LoadFromCache(hd.cache, cacheKey, w); err == nil {
		return
	}

	contract, err, status := cdigest.ParseRequest(w, r, "contract")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handlePaymentDesignInGroup(contract)
	}); err != nil {
		cdigest.HTTP2HandleError(w, err)
	} else {
		cdigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)

		if !shared {
			cdigest.HTTP2WriteCache(w, cacheKey, time.Second*3)
		}
	}
}

func (hd *Handlers) handlePaymentDesignInGroup(contract string) ([]byte, error) {
	var de *types.Design
	var st base.State

	de, st, err := PaymentDesign(hd.database, contract)
	if err != nil {
		return nil, err
	}

	i, err := hd.buildPaymentDesign(contract, *de, st)
	if err != nil {
		return nil, err
	}
	return hd.encoder.Marshal(i)
}

func (hd *Handlers) buildPaymentDesign(contract string, de types.Design, st base.State) (cdigest.Hal, error) {
	h, err := hd.combineURL(HandlerPathPaymentDesign, "contract", contract)
	if err != nil {
		return nil, err
	}

	var hal cdigest.Hal
	hal = cdigest.NewBaseHal(de, cdigest.NewHalLink(h, nil))

	h, err = hd.combineURL(cdigest.HandlerPathBlockByHeight, "height", st.Height().String())
	if err != nil {
		return nil, err
	}
	hal = hal.AddLink("block", cdigest.NewHalLink(h, nil))

	for i := range st.Operations() {
		h, err := hd.combineURL(cdigest.HandlerPathOperation, "hash", st.Operations()[i].String())
		if err != nil {
			return nil, err
		}
		hal = hal.AddLink("operations", cdigest.NewHalLink(h, nil))
	}

	return hal, nil
}

func (hd *Handlers) handlePaymentAccountInfo(w http.ResponseWriter, r *http.Request) {
	cachekey := cdigest.CacheKeyPath(r)
	if err := cdigest.LoadFromCache(hd.cache, cachekey, w); err == nil {
		return
	}

	contract, err, status := cdigest.ParseRequest(w, r, "contract")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	account, err, status := cdigest.ParseRequest(w, r, "address")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	if v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		return hd.handlePaymentAccountInfoInGroup(contract, account)
	}); err != nil {
		cdigest.HTTP2HandleError(w, err)
	} else {
		cdigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)

		if !shared {
			cdigest.HTTP2WriteCache(w, cachekey, time.Second*3)
		}
	}
}

func (hd *Handlers) handlePaymentAccountInfoInGroup(contract, account string) ([]byte, error) {
	var accountInfoValue *AccountInfoValue

	accountInfoValue, err := AccountInfo(hd.database, contract, account)
	if err != nil {
		return nil, err
	}

	i, err := hd.buildAccountInfoValue(contract, *accountInfoValue)
	if err != nil {
		return nil, err
	}
	return hd.encoder.Marshal(i)
}

func (hd *Handlers) buildAccountInfoValue(contract string, it AccountInfoValue) (cdigest.Hal, error) {
	h, err := hd.combineURL(
		HandlerPathPaymentAccountInfo,
		"contract", contract, "address", it.setting.Address().String(),
	)
	if err != nil {
		return nil, err
	}

	var hal cdigest.Hal
	hal = cdigest.NewBaseHal(it, cdigest.NewHalLink(h, nil))

	return hal, nil
}
