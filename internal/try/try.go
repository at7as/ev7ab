package try

import (
	"bufio"
	"math"
	"math/rand"
	"os"

	"gonum.org/v1/gonum/stat"
)

// var dlen int = 100
// var mlen1 int = 10000
// var mlen2 int = 10000
// var stages1 []int = []int{4, 2}
// var stages2 []int = []int{16, 2}
var outmod [8]uint8 = [8]uint8{1, 2, 4, 8, 16, 32, 64, 128}

// var data1 []dataSample1
// var m1 []*model1
// var data2 []dataSample2
// var m2 []*model2

// func main() {

// 	fmt.Printf("dlen: %v\n", dlen)
// 	fmt.Printf("mlen1: %v\n", mlen1)
// 	fmt.Printf("stages1: %v\n", stages1)

// 	data1 = make([]dataSample1, dlen)
// 	for i := range data1 {
// 		data1[i] = newDataSample1()
// 	}

// 	m1 = make([]*model1, mlen1)
// 	for i := range m1 {
// 		m := &model1{in: 4}
// 		for _, v := range stages1 {
// 			m.add(v)
// 		}
// 		m1[i] = m
// 	}

// 	tr0 := math.MaxFloat64
// 	tr1 := math.MaxFloat64

// 	for _, m := range m1 {

// 		mrs0 := make([]float64, dlen)
// 		mrs1 := make([]float64, dlen)

// 		for i, ds := range data1 {
// 			r := m.exec(ds.in)
// 			mrs0[i] = math.Abs(r[0]/ds.out[0] - 1.0)
// 			mrs1[i] = math.Abs(r[1]/ds.out[1] - 1.0)
// 			m.clear()
// 		}

// 		sort.Float64s(mrs0)
// 		sort.Float64s(mrs1)
// 		mr0 := min(stat.Mean(mrs0, nil), stat.Quantile(0.5, stat.Empirical, mrs0, nil))
// 		mr1 := min(stat.Mean(mrs1, nil), stat.Quantile(0.5, stat.Empirical, mrs1, nil))

// 		if mr0 < tr0 && mr1 < tr1 {
// 			fmt.Printf("%.3f %.3f\n", mr0, mr1)
// 			tr0 = mr0
// 			tr1 = mr1
// 		}

// 	}

// 	fmt.Printf("mlen2: %v\n", mlen2)
// 	fmt.Printf("stages2: %v\n", stages2)

// 	data2 = make([]dataSample2, dlen)
// 	for i, v := range data1 {
// 		data2[i] = convertDataSample1to2(v)
// 	}

// 	m2 = make([]*model2, mlen2)
// 	for i := range m2 {
// 		m := &model2{in: 4}
// 		m.add(16)
// 		m.out(2)
// 		m2[i] = m
// 	}

// 	tr0 = math.MaxFloat64
// 	tr1 = math.MaxFloat64

// 	for _, m := range m2 {

// 		mrs0 := make([]float64, dlen)
// 		mrs1 := make([]float64, dlen)

// 		for i, ds := range data2 {
// 			r := m.exec(ds.in)
// 			mrs0[i] = math.Abs(float64(r[0])/256.0/ds.out[0] - 1.0)
// 			mrs1[i] = math.Abs(float64(r[1])/256.0/ds.out[1] - 1.0)
// 			m.clear()
// 		}

// 		sort.Float64s(mrs0)
// 		sort.Float64s(mrs1)
// 		mr0 := min(stat.Mean(mrs0, nil), stat.Quantile(0.5, stat.Empirical, mrs0, nil))
// 		mr1 := min(stat.Mean(mrs1, nil), stat.Quantile(0.5, stat.Empirical, mrs1, nil))

// 		if mr0 < tr0 && mr1 < tr1 {
// 			fmt.Printf("%.3f %.3f\n", mr0, mr1)
// 			tr0 = mr0
// 			tr1 = mr1
// 		}

// 	}

// }

type DataSample1 struct {
	In  []float64
	Out []float64
}

func NewDataSample1() DataSample1 {
	a := rand.Float64()
	b := rand.Float64()
	c := rand.Float64()
	d := rand.Float64()
	x := (a + b*b + c*d) / 3.0
	y := (d + c*c + b*a) / 3.0
	return DataSample1{
		In:  []float64{a, b, c, d},
		Out: []float64{x, y},
	}
}

type DataSample2 struct {
	In  []uint8
	Out []float64
}

func ConvertDataSample1to2(ds1 DataSample1) DataSample2 {
	ds2 := DataSample2{
		In:  make([]uint8, len(ds1.In)),
		Out: ds1.Out,
	}
	for i, v := range ds1.In {
		ds2.In[i] = uint8(math.Floor(v * 256))
	}
	return ds2
}

type Model1 struct {
	In   int
	link [][]float64
	node [][]atom1
}

type atom1 struct {
	in    []float64
	value float64
}

func (m *Model1) Add(n int) {

	var in int
	if len(m.node) == 0 {
		in = m.In
	} else {
		in = len(m.node[len(m.node)-1])
	}
	link := make([]float64, in*n)
	for i := range link {
		link[i] = rand.Float64()
	}
	m.link = append(m.link, link)
	node := make([]atom1, n)
	for i := range node {
		node[i] = atom1{
			in: make([]float64, in),
		}
	}
	m.node = append(m.node, node)

}

func (m *Model1) Exec(in []float64) []float64 {

	for i := range m.node {

		n := len(m.node[i])
		if i == 0 {

			for iin, vin := range in {
				for ia := range m.node[i] {
					m.node[i][ia].in[iin] = m.link[i][iin*n+ia] * vin
				}
			}

		} else {

			for iin, vin := range m.node[i-1] {
				for ia := range m.node[i] {
					m.node[i][ia].in[iin] = m.link[i][iin*n+ia] * vin.value
				}
			}

		}

		for ia, a := range m.node[i] {
			m.node[i][ia].value = stat.Mean(a.in, nil)
		}

	}

	res := []float64{}
	for _, a := range m.node[len(m.node)-1] {
		res = append(res, a.value)
	}

	return res
}

func (m *Model1) Clear() {

	for i := range m.node {
		for ii := range m.node[i] {
			m.node[i][ii].value = 0
			m.node[i][ii].in = make([]float64, len(m.node[i][ii].in))
		}
	}

}

type Model2 struct {
	In   int
	link [][]uint8
	node [][]atom2
}

type atom2 struct {
	in    []uint8
	value uint8
}

func (m *Model2) Add(n int) {

	var in int
	if len(m.node) == 0 {
		in = m.In
	} else {
		in = len(m.node[len(m.node)-1])
	}
	link := make([]uint8, in*n)
	for i := range link {
		link[i] = uint8(rand.Intn(math.MaxUint8 + 1))
	}
	m.link = append(m.link, link)
	node := make([]atom2, n)
	for i := range node {
		node[i] = atom2{
			in: make([]uint8, in),
		}
	}
	m.node = append(m.node, node)

}

func (m *Model2) Out(n int) {

	// link := make([]uint8, 8*n)
	// for i := 0; i < n; i++ {
	// 	for ii, v := range outmod {
	// 		link[i*8+ii] = v
	// 	}
	// }
	// m.link = append(m.link, link)
	node := make([]atom2, n)
	for i := range node {
		node[i] = atom2{
			in: make([]uint8, 8),
		}
	}
	m.node = append(m.node, node)

}

func (m *Model2) Exec(in []uint8) []uint8 {

	for i := range m.node {

		n := len(m.node[i])

		if i == 0 {

			for iin, vin := range in {
				for ia := range m.node[i] {
					m.node[i][ia].in[iin] = vin ^ m.link[i][iin*n+ia]
				}
			}

			for ia, a := range m.node[i] {
				m.node[i][ia].value = bitSumAND(a.in)
			}

		} else {

			last := m.node[len(m.node)-2]
			for ii := range m.node[len(m.node)-1] {
				for iii := range outmod {
					if last[ii*8+iii].value > 0 {
						m.node[len(m.node)-1][ii].in[iii] = outmod[iii]
					}
				}
				m.node[len(m.node)-1][ii].value = sum(m.node[len(m.node)-1][ii].in)
			}

		}

	}

	res := []uint8{}
	for _, a := range m.node[len(m.node)-1] {
		res = append(res, a.value)
	}

	return res
}

func (m *Model2) Clear() {

	for i := range m.node {
		for ii := range m.node[i] {
			m.node[i][ii].value = 0
			m.node[i][ii].in = make([]uint8, len(m.node[i][ii].in))
		}
	}

}

func bitSumAND(v []uint8) uint8 {
	sum := v[0]
	for i := 1; i < len(v); i++ {
		sum &= v[i]
	}
	return sum
}

func sum(v []uint8) uint8 {
	var sum uint8 = 0
	for _, vv := range v {
		sum += vv
	}
	return sum
}

func wait() {
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
