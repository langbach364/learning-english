package main

import ()

type WordData struct {
    Word string
    Types []WordType
}

type WordType struct {
    Type string
    Definitions []Definition
}

type Definition struct {
    Meaning string
    Examples []string
}