package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

//type Flight struct {
//	fl map[string]commentObj
//}

type commentObj struct {
	passenger string
	point     int
	comment   string
}

func main() {

	var line string
	var itr int
	var err error
	var printInputLog []string

	flighList := make(map[string][]commentObj)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	line = scanner.Text()
	//first inputs number
	itr, err = strconv.Atoi(line)
	if err != nil {
		log.Fatal("bad input numbeer for for first inputs ")
	}

	for i := 0; i < itr; i++ {
		scanner.Scan()
		line = scanner.Text()
		flighList[line] = []commentObj{}
	}
	//getting number of second inputs
	scanner.Scan()
	line = scanner.Text()
	itr, err = strconv.Atoi(line)
	if err != nil {
		log.Fatal("bad input number for for second inputs ")
	}

	for i := 0; i < itr; i++ {
		scanner.Scan()
		line = scanner.Text()
		lineSlice := strings.Split(line, " ")

		if len(lineSlice) != 2 {
			log.Fatal("bad input for input 2")

		}

		val, ok := flighList[lineSlice[1]]

		if !ok {
			printInputLog = append(printInputLog, "Invalid flight "+lineSlice[1])

		} else {

			check := true

			for _, j := range val {
				if j.passenger == lineSlice[0] {
					printInputLog = append(printInputLog, "Duplicate ticket for "+lineSlice[1]+" by "+lineSlice[0])
					check = false
				}

			}
			if check {
				flighList[lineSlice[1]] = append(flighList[lineSlice[1]], commentObj{passenger: lineSlice[0]})
			}
		}

		//	if ok && val.passenger == lineSlice[0] {
		//		printInputLog = append(printInputLog, "Duplicate ticket for "+lineSlice[1]+" by "+lineSlice[0])
		//	} else {
		//		flighList[lineSlice[1]] = commentObj{passenger: lineSlice[0]}
		//	}
	}

	//getting number of third inputs
	scanner.Scan()
	line = scanner.Text()
	itr, err = strconv.Atoi(line)
	if err != nil {
		log.Fatal("bad input number for for third inputs ")
	}

	for i := 0; i < itr; i++ {

		scanner.Scan()
		line = scanner.Text()
		lineSlice := strings.Split(line, " ")

		val, ok := flighList[lineSlice[1]]

		if !ok {
			printInputLog = append(printInputLog, "Invalid flight "+lineSlice[1])
		} else {

			check := true
			var index int
			for i, j := range val {

				if j.passenger == lineSlice[0] {
					check = false
					index = i
					break
				}
			}
			if check {
				printInputLog = append(printInputLog, "Invalid passenger for "+lineSlice[1]+" "+lineSlice[0])
			}

			if !check && val[index].comment != "" {
				printInputLog = append(printInputLog, "Duplicate comment for "+lineSlice[1]+" by "+lineSlice[0])

			}
			if !check && val[index].comment == "" {
				printInputLog = append(printInputLog, "Accepted comment "+lineSlice[1]+" by "+lineSlice[0])
				tempObj := flighList[lineSlice[1]][index]
				tempObj.comment = lineSlice[3]
				tempPoint, _ := strconv.Atoi(lineSlice[2])
				tempObj.point = tempPoint
				flighList[lineSlice[1]][index] = tempObj

			}

		}

		/*	if val.passenger != lineSlice[0] {
				printInputLog =append(printInputLog,"Invalid passenger for "+lineSlice[1]+" "+lineSlice[0])
			}

			if val.passenger == lineSlice[0] && val.comment != "" {
				printInputLog = append(printInputLog, "Invalid passenger for "+lineSlice[1]+" "+lineSlice[0])
			}else {
				printInputLog = append(printInputLog, "Accepted comment "+lineSlice[1]+" "+lineSlice[0])
				tempObj:= flighList[lineSlice[1]]
				tempObj.comment = lineSlice[3]
				tempPoint, _:=strconv.Atoi(lineSlice[2])
				tempObj.point=tempPoint
				flighList[lineSlice[1]] =tempObj

			}*/
	}

	for _, i := range printInputLog {
		fmt.Println(i)
	}

	var keys []string
	for flName, _ := range flighList {
		keys = append(keys, flName)
	}
	//	keySorted := sort.StringSlice(keys)

	sort.Strings(keys)

	for _, flName := range keys {
		var avePoint float64
		var flightLen int
		for _, j := range flighList[flName] {
			if j.point != 0 {
				avePoint += float64(j.point)
				flightLen++
			}
		}
		if flightLen > 0 {
			fmt.Printf("Average score for %s is %.2f\n", flName, avePoint/float64(flightLen))
		}

	}

}
