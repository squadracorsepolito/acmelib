# acmelib

[![Go Reference](https://pkg.go.dev/badge/github.com/squadracorsepolito/acmelib.svg)](https://pkg.go.dev/github.com/squadracorsepolito/acmelib)
[![Go Report Card](https://goreportcard.com/badge/github.com/squadracorsepolito/acmelib)](https://goreportcard.com/report/github.com/squadracorsepolito/acmelib)

A [Golang](https://go.dev/) package for modelling complex CAN networks.

The package documentation can be found [here](https://pkg.go.dev/github.com/FerroO2000/acmelib).

## Model

```mermaid
flowchart
    subgraph Network
        bus(Bus)

        subgraph Node
            node-int(Node Interface)
        end

        message(Message)

        subgraph Signal
            std-sig(Standard Signal)
            enum-sig(Enum Signal)
            mux-sig(Multiplexer Signal)
        end

        sig-type(Signal Type)
        sig-unit(Signal Unit)
        sig-enum(Signal Enum)

        attribute(Attribute)
    end

    bus --Attaches--o node-int

    node-int --Sends--o message
    message -.Receives.-o node-int

    message --Contains--o Signal

    std-sig --o sig-type
    std-sig --o sig-unit

    enum-sig --o sig-enum

    mux-sig --o std-sig
    mux-sig --o enum-sig
    mux-sig --o mux-sig

    attribute -.- bus
    attribute -.- Node
    attribute -.- message
    attribute -.- Signal
```

## Getting started

### Prerequisites

[Golang](https://go.dev/) 1.22

### Installation

Run the following Go command to install the `acmelib` package:

```sh
go get -u github.com/squadracorsepolito/acmelib
```

## TODOs

-   Adding examples
