package byteSample

type StudentClass struct {
	Name string `json:"Name"`
}

type Employees struct {
	Name    string         `json:"Name"`
	Salary  float64        `json:"Salary"`
	Student []StudentClass `json:"Student"`
}

type Person struct {
	Epoch    int               `json:"Epoch"`
	Watch    [][]int           `json:"Watch"`
	Xyz      int               `json:"Xyz"`
	Salary   float64           `json:"Salary"`
	Employee [][][][]Employees `json:"Employee"`
}
