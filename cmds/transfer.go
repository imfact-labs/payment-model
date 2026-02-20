package cmds

import (
	"context"

	ccmds "github.com/imfact-labs/currency-model/app/cmds"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/payment-model/operation/payment"
	"github.com/pkg/errors"
)

type TransferCommand struct {
	BaseCommand
	ccmds.OperationFlags
	Sender   ccmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract ccmds.AddressFlag    `arg:"" name:"contract" help:"contract address" required:"true"`
	Receiver ccmds.AddressFlag    `arg:"" name:"receiver" help:"receiver address" required:"true"`
	Amount   ccmds.BigFlag        `arg:"" name:"amount" help:"amount" required:"true"`
	Currency ccmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	sender   base.Address
	contract base.Address
	receiver base.Address
}

func (cmd *TransferCommand) Run(pctx context.Context) error { // nolint:dupl
	if _, err := cmd.prepare(pctx); err != nil {
		return err
	}

	if err := cmd.parseFlags(); err != nil {
		return err
	}

	op, err := cmd.createOperation()
	if err != nil {
		return err
	}

	ccmds.PrettyPrint(cmd.Out, op)

	return nil
}

func (cmd *TransferCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	a, err := cmd.Sender.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid sender format, %q", cmd.Sender)
	} else {
		cmd.sender = a
	}

	a, err = cmd.Contract.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid contract format, %q", cmd.Contract)
	} else {
		cmd.contract = a
	}

	a, err = cmd.Receiver.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid receiver format, %q", cmd.Receiver)
	} else {
		cmd.receiver = a
	}

	return nil
}

func (cmd *TransferCommand) createOperation() (base.Operation, error) { // nolint:dupl
	e := util.StringError("failed to create transfer operation")

	fact := payment.NewTransferFact([]byte(cmd.Token), cmd.sender, cmd.contract, cmd.receiver, cmd.Amount.Big, cmd.Currency.CID)

	op, err := payment.NewTransfer(fact)
	if err != nil {
		return nil, e.Wrap(err)
	}
	err = op.Sign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
