package types

import (
	"encoding/json"

	"github.com/imfact-labs/currency-model/common"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/mitum2/util/hint"
	"github.com/imfact-labs/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var DesignHint = hint.MustNewHint("mitum-payment-design-v0.0.1")

var maxAccounts = 1000

type Design struct {
	hint.BaseHinter
	settings map[string]Setting
}

func NewDesign() Design {
	settings := make(map[string]Setting)
	return Design{
		BaseHinter: hint.NewBaseHinter(DesignHint),
		settings:   settings,
	}
}

func (de Design) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		de.BaseHinter,
	); err != nil {
		return err
	}

	return nil
}

func (de Design) Bytes() []byte {
	var bac []byte
	if de.settings != nil {
		ac, _ := json.Marshal(de.settings)
		bac = valuehash.NewSHA256(ac).Bytes()
	} else {
		bac = []byte{}
	}

	return util.ConcatBytesSlice(bac)
}

func (de Design) Hash() util.Hash {
	return de.GenerateHash()
}

func (de Design) GenerateHash() util.Hash {
	return valuehash.NewSHA256(de.Bytes())
}

func (de Design) AccountSettings() map[string]Setting {
	return de.settings
}

func (de Design) AccountSetting(account string) *Setting {
	v, found := de.settings[account]

	if !found {
		return nil
	}

	return &v
}

func (de *Design) AddAccountSetting(setting Setting) error {
	de.settings[setting.Address().String()] = setting

	if len(de.settings) > maxAccounts {
		return common.ErrValOOR.Wrap(
			errors.Errorf("accounts over allowed, %d > %d", len(de.settings), maxAccounts))
	}

	return nil
}

func (de *Design) UpdateAccountSetting(account Setting) error {
	_, found := de.settings[account.Address().String()]
	if !found {
		return common.ErrValueInvalid.Wrap(
			errors.Errorf("account, %v not registered in service", account.Address()))
	}
	de.settings[account.Address().String()] = account

	return nil
}

func (de *Design) RemoveAccountSetting(account base.Address) error {
	_, found := de.settings[account.String()]
	if !found {
		return errors.Errorf("account, %v not registered in service", account)
	}

	delete(de.settings, account.String())

	return nil
}
