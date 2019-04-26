package main

import (
	"time"
)

//messageは1つのメッセージを指す

type message struct {
	Name string
	Message string
	When time.Time
}
