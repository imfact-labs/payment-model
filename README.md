### mitum-deposit

*mitum-deposit* is a payment contract model based on the second version of mitum(aka [mitum2](https://github.com/ProtoconNet/mitum2)).

#### Installation

Before you build `mitum-deposit`, make sure to run mongodb for digest api.

```sh
$ git clone https://github.com/ProtoconNet/mitum-deposit

$ cd mitum-deposit

$ go build -o ./mitum-deposit
```

#### Run

```sh
$ ./mitum-deposit init --design=<config file> <genesis file>

$ ./mitum-deposit run <config file> --dev.allow-consensus
```

[standalong.yml](standalone.yml) is a sample of `config file`.

[genesis-design.yml](genesis-design.yml) is a sample of `genesis design file`.
