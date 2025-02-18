package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

var a1 int // размерность строк
var a2 int // размерность столбцов
var mi int // строка главного элемента
var mj int // столбец главного элемента

func rank(mat1 [][]float64, a1 int, a2 int) (int, error) {

	k1 := 0
	k2 := 0
	c1 := 0
	c2 := 0
	for i := 0; i < a1; i++ {
		for j := 0; j < a2; j++ {
			if mat1[i][j] == 0 {
				k1++
			}
		}
		if k1 == a2 {
			c1++
		}
		k1 = 0
	}

	for i := 0; i < a1; i++ {
		for j := 0; j < a2-1; j++ {
			if mat1[i][j] == 0 {
				k2++
			}
		}
		if k2 == a2-1 {
			c2++
		}
		k2 = 0
	}
	fmt.Println(c1)
	fmt.Println(c2)
	for i := 0; i < a1; i++ {
		fmt.Println(mat1[i])
	}
	if c1 == c2 {
		return 1, nil
	} else {
		return 0, fmt.Errorf("cистема_несовместна")
	}
}

func main() {
	file, err := os.Open("gen_3000.txt")        // открываем файл
	data := make([]byte, 100000000)             //создаем буффер для данных
	a1 = 3000                                   //записываем количество строк
	a2 = 3001                                   //записываем число столбцов
	var max_el float64 = float64(math.MinInt64) // переменная для максимального элемента
	var s1 float64                              //буфер для переноса строк
	var s2 float64                              //буфер для вычитания
	var x [][]float64                           //инициализация массива коррекции
	k := make([]float64, a1)
	ans := make([]float64, a1)
	if err != nil { //проверяем открытие на ошибки
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		n, errr := file.Read(data) //пытаемся прочитать файл
		if errr == io.EOF {        // если конец файла
			break // выходим из цикла
		}
		fmt.Print(string(data[:n]))
	}

	r := strings.NewReader(string(data[:]))
	s := bufio.NewScanner(r)

	var matrix [][]float64

	for s.Scan() {
		records := strings.Fields(s.Text())
		line := make([]float64, len(records))
		matrix = append(matrix, line)
		for i := range records {
			line[i], _ = strconv.ParseFloat(records[i], 64)
		}
	}
	start := time.Now()
	for i := 0; i < a1; i++ { //заполнение массива коррекции
		x_cor := make([]float64, a2)
		x = append(x, x_cor)
		for j := 0; j < a2; j++ {
			x[i][j] = float64(j)
		}
	}

	for o := 0; o < a1-1; o++ {
		for i := o; i < a1; i++ { //поиск максимального (главного) элемента
			for j := o; j < a2-1; j++ {
				if math.Abs(matrix[i][j]) > max_el {
					max_el = matrix[i][j]
					mi = i
					mj = j
				}
			}
		}

		for i := 0; i < a1; i++ { //перенос строк
			for j := 0; j < a2; j++ {
				if mi == i && mi != o {
					s1 = matrix[o][j] //перенос строк исходного массива
					matrix[o][j] = matrix[i][j]
					matrix[i][j] = s1
					s1 = x[o][j] //перенос строк массива коррекции
					x[o][j] = x[i][j]
					x[i][j] = s1
				}
			}
		}

		for i := 0; i < a1; i++ { //перенос столбцов
			for j := 0; j < a2-1; j++ {
				if mj == j && mj != o {
					s1 = matrix[i][o] //перенос столбцов исходного массива
					matrix[i][o] = matrix[i][j]
					matrix[i][j] = s1
					s1 = x[i][o] //перенос столбцов массива коррекции
					x[i][o] = x[i][j]
					x[i][j] = s1
				}
			}
		}

		for j := o; j < a2; j++ { //нормировка строк
			matrix[o][j] = matrix[o][j] / max_el
		}

		for i := o; i < a1; i++ { //приведение к треугольному виду
			s2 = matrix[i][o]
			for j := o; j < a2; j++ {
				if i != o {
					matrix[i][j] = matrix[i][j] - s2*matrix[o][j]
				}
			}
		}

		max_el = float64(math.MinInt64)
	}

	rank_slice := matrix[:][:]

	_, err1 := rank(rank_slice, a1, a2)

	if err1 != nil {
		fmt.Println(err1)
		os.Exit(1)
	} else {
		fmt.Print("Cистема совместна, решение возможно \n")
	}

	max_el = matrix[a1-1][a1-1]
	for j := a1 - 1; j < a2; j++ { //нормировка последней строки
		matrix[a1-1][j] = matrix[a1-1][j] / max_el
	}

	for i := a1 - 1; i >= 0; i-- { //обратный ход, получение корней
		k[i] = matrix[i][a1]
		for j := 0; j < a1; j++ {
			if j != i {
				k[i] = k[i] - matrix[i][j]*k[j]
			}
		}
		k[i] = k[i] / matrix[i][i]
	}

	for i := 0; i < a1; i++ {
		for j := 0; j < a1; j++ {
			if i == j {
				ans[int(x[i][j])] = k[i]
			}
		}
	}

	elapsed := time.Since(start)
	for i := 0; i < a1; i++ {
		fmt.Println(ans[i])
	}
	fmt.Println(elapsed)
}
