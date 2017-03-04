package main

import "fmt"

var text []interface{}
var another interface{} = true

func main(){
	text = append(text, "go")
	text = append(text, 12)
	text = append(text, 12.456)
	for _, v := range text {
		switch value := v.(type) {
		case string:
			fmt.Printf("%[1]v %[1]T\n", value)
		case float32,float64:
			fmt.Printf("%[1]v %[1]T\n", value)
		case int:
			fmt.Printf("%[1]v %[1]T\n", value)
		default:
			fmt.Println("unknown")
		}
	}
	if v, ok := another.(bool); ok {
		fmt.Printf("this is %T - %[1]v", v)
	}
}
