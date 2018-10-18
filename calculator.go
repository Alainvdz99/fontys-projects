package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:]
	fmt.Println(calculate(args))
}

// calculates the value
func calculate(args []string) string {
	return addOrSubtract(multiplyOrDivide(args))[0]
}

// multiplies or divides values
func multiplyOrDivide(args []string) []string {
	var temp1 []string
	var temp2 []string
	var temp []string

	var noMultipleOrDivide = true

	for i, element := range args  {
		if element == "x" && (i + 1) < len(args) {
			var int1, _ = strconv.Atoi(args[i - 1])
			var int2, _ = strconv.Atoi(args[i + 1])

			temp1 = args[:i - 1]
			temp2 = args[i + 1:]
			temp2[0] = strconv.Itoa(int1 * int2)

			temp = append(temp1, temp2...)

			noMultipleOrDivide = false
			break
		}
		if element == "/" && (i + 1) < len(args) {
			var int1, _ = strconv.Atoi(args[i - 1])
			var int2, _ = strconv.Atoi(args[i + 1])

			temp1 = args[:i - 1]
			temp2 = args[i + 1:]
			temp2[0] = strconv.Itoa(int1 / int2)
			noMultipleOrDivide = false

			temp = append(temp1, temp2...)

			break
		}
	}

	if noMultipleOrDivide {
		return args
	}
	if len(temp2) == 1 {
		return temp
	} else {
		return multiplyOrDivide(temp)
	}
}

// adds or subtracts values
func addOrSubtract(args []string) []string {
	var temp1 []string
	var temp2 []string
	var temp []string

	var noPlusOrMinus = true

	for i, element := range args  {
		if element == "+" && (i + 1) < len(args) {
			var int1, _ = strconv.Atoi(args[i - 1])
			var int2, _ = strconv.Atoi(args[i + 1])

			temp1 = args[:i - 1]
			temp2 = args[i + 1:]
			temp2[0] = strconv.Itoa(int1 + int2)

			temp = append(temp1, temp2...)

			noPlusOrMinus = false
			break
		}
		if element == "-" && (i + 1) < len(args) {
			var int1, _ = strconv.Atoi(args[i - 1])
			var int2, _ = strconv.Atoi(args[i + 1])

			temp1 = args[:i - 1]
			temp2 = args[i + 1:]
			temp2[0] = strconv.Itoa(int1 - int2)
			noPlusOrMinus = false

			temp = append(temp1, temp2...)

			break
		}
	}

	if noPlusOrMinus {
		return args
	}
	if len(temp2) == 1 {
		return temp
	} else {
		return addOrSubtract(temp)
	}
}
