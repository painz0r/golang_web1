package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

//var d []interface{}

//const p = "death"
//
//func zero(z *int) {
//	fmt.Println(z)
//	fmt.Println(*z)
//	*z = 0
//}
//
type Human struct {
	Age          int
	Name, Gender string
}

type origin struct {
	Human
	Skin        string
	NotExported bool
}

func (p *origin) askName(name ...string) []string {
	fmt.Println("Enter your name ")
	fmt.Scanln(&name)
	fmt.Print("Your name is ")
	p.Name = "Ira"
	//fmt.Println(p.Name)
	return name
}

func (p *origin) defaultName() string {
	p.Name = "Ira2"
	return p.Name
}

type people []string

func (p people) Len() int {
	return len(p)
}
func (p people) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p people) Less(i, j int) bool {
	return p[i] < p[j]
}

type Sorting interface {
	Len() int
	Swap(i, j int)
	Less(i, j int) bool
}

type Reverse interface {
}

var v interface{}

func incrementor() <-chan int {
	out := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			out <- i
		}
		close(out)
	}()
	return out
}

func puller(c <-chan int) chan int {
	out := make(chan int)
	go func() {
		var sum int
		for n := range c {
			sum += n
		}
		out <- sum
		close(out)
	}()
	return out
}

func main() {

	c := make(chan int)
	done := make(chan bool)

	go func() {
		for i := 0; i <= 10; i++ {
			c <- i
		}
		done <- true
	}()

	go func() {
		for i := 0; i <= 10; i++ {
			c <- i
		}
		done <- true
	}()
	go func() {
		<-done
		<-done
		close(c)
	}()

	for n := range c {
		fmt.Println(n)
	}

	/*channels as arguments in functions*/
	ch := incrementor()
	cSum := puller(ch)

	for n := range cSum {
		fmt.Println("Total Value is ", n)
	}

	v = 7
	//var f interface{} = 7
	v, ok := v.(int)
	if ok {
		fmt.Printf("%v is of type %T\n", v, v)
	}
	if ok := "equal"; true {
		fmt.Println(ok)
	}

	studyGroup := people{"Zeno", "John", "Al", "Jenny"}
	s := []string{"Zeno", "John", "Al", "Jenny"}
	n := []int{7, 4, 8, 2, 9, 19, 12, 32, 3}
	fmt.Println(studyGroup)
	sort.Sort(Sorting(studyGroup))
	fmt.Println(studyGroup)
	sort.Sort(Sorting(sort.Reverse(studyGroup)))
	fmt.Println(studyGroup)
	fmt.Println("----------------------")
	fmt.Println(s)
	//sort.Strings(s)
	//sort.StringSlice(s).Sort()
	sort.Sort(sort.StringSlice(s))
	fmt.Println(s)
	sort.Sort(sort.Reverse(sort.StringSlice(s)))
	fmt.Println(s)
	fmt.Println("----------------------")
	fmt.Println(n)
	sort.Ints(n)
	fmt.Println(n)
	sort.Sort(sort.Reverse(sort.IntSlice(n)))
	fmt.Println(n)

	person1 := origin{
		Human: Human{Age: 32,
			Name:   "Ross",
			Gender: "Male",
		},
		Skin: "white",
	}
	fmt.Println(person1.Name, person1.Age, person1.Gender, person1.Skin)
	fmt.Println(person1.askName())
	fmt.Println(person1.Name)
	fmt.Println(person1.defaultName())
	fmt.Println(person1.Name)
	//bs, _ := json.Marshal(person1)
	//fmt.Printf("%T\n", bs)
	//fmt.Println(bs)
	//fmt.Println(string(bs))
	json.NewEncoder(os.Stdout).Encode(person1)
	////x := 5
	////fmt.Println(&x)
	////zero(&x)
	////fmt.Println(x) // x is 0
	////
	//a := "test"
	//b := 4.44
	//c := 12
	//d := []interface{}{}
	//fmt.Printf("%T\n", d)
	//d = append(d, a, b, c)
	//for i:=0; i < len(d); i++{
	//	fmt.Printf("value %v ",d[i])
	//	fmt.Printf("type %T\n",d[i])
	//}
	////const q = 42
	////fmt.Println(q)
	////fmt.Println(p)
	//g := 41
	//fmt.Println(g)
	//fmt.Println(&g)
	//var f *int = &g
	//fmt.Println(f)
	//fmt.Println(*f)
	//*f = 42
	//fmt.Println(*f)
	//for i := 0; i <= 10; i++ {
	//	for j := 0; j <= 10; j++ {
	//		fmt.Println(i, " - ", j)
	//	}
	//}
	//i := 0
	//
	//for {
	//	fmt.Println(i)
	//	if i >= 10 {
	//		break
	//	}
	//	i++
	//}
	//l := 10000
	//
	//for {
	//	l++
	//	if l % 2 == 0 {
	//		continue
	//	}
	//	fmt.Println(l, " - ", string(l), " - ", []byte(string(l)))
	//	if l >= 10050 {
	//		break
	//	}
	//
	//}
	//for i := 1; i <= 100; i++ {
	//	if i % 15 == 0 {
	//		fmt.Println(i, "--FIZZBUZZ")
	//	} else if i % 3 == 0 {
	//		fmt.Println(i, "--FIZZ")
	//
	//	} else if i % 5 == 0 {
	//		fmt.Println(i, "--BUZZ")
	//	} else {
	//		fmt.Println(i)
	//	}
	//}
	//
	//n := []float64{23, 242, 5, 45, rand.Float64()}
	//
	//variadic_Expression := func(data ...float64) float64 {
	//	var sums float64
	//	for _, i := range data {
	//		sums += i
	//	}
	//	return sums / float64(len(data))
	//}
	//res := variadic_Expression(n...)
	//fmt.Println(res)
	//
	//callback_expression := func(numbers []int, callback func(single_number int)) {
	//	for _, n := range numbers {
	//		callback(n)
	//	}
	//}
	//
	//callback_expression([]int{1, 2, 3, 4, 5}, func(n int) {
	//	fmt.Println(n)
	//})
	//xs := []int{}
	//
	//fmt.Printf("%T\n", xs)
	//
	//fmt.Println(factorial(5))
	//
	//factorial_2 := func(n int) int {
	//	num := 1
	//	for n >= 1 {
	//		num *= n
	//		n -= 1
	//	}
	//	return num
	//}
	//fmt.Println(factorial_2(5))
	//
	//fmt.Println(func() string {
	//	return "Hey i'm annonymous self-executing function"
	//}())
	//
	//my_1 := func(n int) (float64, bool) {
	//
	//	return float64(n) / 2, n % 2 == 0
	//}
	//h, even := my_1(5)
	//
	//fmt.Println(h, even)
	//
	//fmt.Println(my_2([]int{1,4,5,6,10,3,15}...))
	//
	//foo(1, 2)
	//foo(1, 2, 3)
	//aSlice := []int{1, 2, 3, 4}
	//foo(aSlice...)
	//foo()
	//
	//mySlices := make([][]int, 0, 5)
	//for i :=1; i <= 5; i++ {
	//	mySlice := make([]int, 0,6)
	//	for j :=1; j <= 6; j++ {
	//		mySlice = append(mySlice, j)
	//	}
	//	mySlices = append(mySlices, mySlice)
	//}
	//for _, v := range mySlices {
	//	fmt.Println(v)
	//}
	//
	//myString := make([]int, 58)
	//for i := 65; i <= 122; i++ {
	//	myString[i-65] = i
	//}
	//for k,v := range myString {
	//	fmt.Printf("key %v - dec %X - value %v - " +
	//		"byte %v\n", k, v, string(v), []byte(string(v)))
	//}
	//for i := 65; i <= 122; i++ {
	//	fmt.Println(i, "- ", string(i), "- ", i % 12)
	//}
	//res, err := http.Get("http://www.gutenberg.org/files/2701/old/moby10b.txt")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//scanner, err := ioutil.ReadAll(res.Body)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer res.Body.Close()
	//fmt.Printf("%s", scanner)
}

func factorial(n int) int {
	if n == 0 {
		return 1
	}
	return n * factorial(n-1)
}

func my_2(numbers ...int) int {
	fmt.Printf("%T\n", numbers)
	var greatest int
	for _, i := range numbers {
		if i > greatest {
			greatest = i
		}
	}
	return greatest
}

func foo(n ...int) {
	fmt.Println(n)
}
