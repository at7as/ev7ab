package main

func main() {

	// l := lab.New(&lib.Example1{})
	// l.AddProject([][]lab.Node{})
	// l.Run()

}

// type data struct {
// 	Rows   []*dataRow
// 	Entity *dataEntity
// 	S      dataSimple
// }

// type dataRow struct {
// 	E   *dataEntity
// 	S   dataSimple
// 	I   int
// 	Str string
// }

// type dataEntity struct {
// 	I1 int
// 	I2 *int
// }

// type dataSimple struct {
// 	S1 string
// 	S2 *string
// }

// func main() {

// 	rows := make([]*dataRow, 100000)
// 	for i := range rows {
// 		i2 := 2
// 		e := &dataEntity{
// 			I1: 1,
// 			I2: &i2,
// 		}
// 		s2 := "2"
// 		s := dataSimple{
// 			S1: "1",
// 			S2: &s2,
// 		}
// 		rows[i] = &dataRow{e, s, i, "str"}
// 	}
// 	i2 := 2
// 	entity := &dataEntity{
// 		I1: 1,
// 		I2: &i2,
// 	}
// 	s2 := "2"
// 	s := dataSimple{
// 		S1: "1",
// 		S2: &s2,
// 	}
// 	d := &data{
// 		Rows:   rows,
// 		Entity: entity,
// 		S:      s,
// 	}

// 	// fmt.Println(d)

// 	filesave, err := os.Create("./data.gob")
// 	writer := bufio.NewWriter(filesave)

// 	if err != nil {
// 		fmt.Println(err, 0)
// 		os.Exit(1)
// 	}

// 	enc := gob.NewEncoder(writer)
// 	err = enc.Encode(d)
// 	if err != nil {
// 		fmt.Println(err, 1)
// 		os.Exit(1)
// 	}

// 	err = writer.Flush()
// 	if err != nil {
// 		fmt.Println(err, 2)
// 		os.Exit(1)
// 	}

// 	filesave.Close()

// 	fileopen, err := os.Open("./data.gob")
// 	if err != nil {
// 		fmt.Println(err, 3)
// 		os.Exit(1)
// 	}

// 	var decD data
// 	dec := gob.NewDecoder(fileopen)
// 	err = dec.Decode(&decD)
// 	if err != nil {
// 		fmt.Println(err, 4)
// 		os.Exit(1)
// 	}

// 	// fmt.Println(decD)
// 	// fmt.Printf("%+v\n", decD)
// 	// fmt.Printf("%+v\n", decD.Rows)
// 	// for _, row := range decD.Rows {
// 	// 	fmt.Printf("%+v\n", row)
// 	// }
// 	// fmt.Printf("%+v\n", decD.Entity)
// 	// fmt.Printf("%+v\n", decD.S)

// 	fileopen.Close()

// }
