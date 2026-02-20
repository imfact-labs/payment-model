### payment-model

*payment-model* is a payment contract model based on [mitum](https://github.com/imfact-labs/mitum2)).

#### Installation

Before you build `payment-model`, make sure to run mongodb for digest api.

```sh
$ git clone https://github.com/imfact-labs/payment-model

$ cd payment-model

$ go build -o ./imfact ./main.go
```

#### Run

```sh
$ ./imfact init --design=<config file> <genesis file>

$ ./imfact run <config file> --dev.allow-consensus
```

[standalong.yml](standalone.yml) is a sample of `config file`.

[genesis-design.yml](genesis-design.yml) is a sample of `genesis design file`.
