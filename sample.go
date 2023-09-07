package main

type Employees struct {
	Name   string  `json:"Name"`
	Salary float64 `json:"Salary"`
}

type Person struct {
	Epoch    int         `json:"Epoch"`
	Watch    string      `json:"Watch"`
	Xyz      int         `json:"Xyz"`
	Salary   float64     `json:"Salary"`
	Employee []Employees `json:"Employee"`
}
