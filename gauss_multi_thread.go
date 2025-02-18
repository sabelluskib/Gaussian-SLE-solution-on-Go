package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
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
	fmt.Println("Количество нулевых строк:", c1)
	fmt.Println("Количество почти нулевых строк:", c2)
	for i := 0; i < a1; i++ {
		fmt.Println("Строка матрицы:", mat1[i])
	}
	if c1 == c2 {
		return 1, nil
	} else {
		return 0, fmt.Errorf("система несовместна")
	}
}

func main() {
	file, err := os.Open("gen_3000.txt") // открываем файл
	if err != nil {                      //проверяем открытие на ошибки
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	data := make([]byte, 10000000000)     //создаем буффер для данных
	a1 = 3000                             //записываем количество строк
	a2 = 3001                             //записываем число столбцов
	var max_el float64 = -math.MaxFloat64 // переменная для максимального элемента
	var s1 float64                        //буфер для переноса строк
	var s2 float64                        //буфер для вычитания
	var x [][]float64                     //инициализация массива коррекции
	k := make([]float64, a1)
	ans := make([]float64, a1)

	_, errr := file.Read(data) //пытаемся прочитать файл
	if errr != nil && errr != io.EOF {
		fmt.Println(errr)
		os.Exit(1)
	}

	r := strings.NewReader(string(data[:]))
	s := bufio.NewScanner(r)

	var matrix [][]float64

	for s.Scan() {
		records := strings.Fields(s.Text())
		line := make([]float64, len(records))
		matrix = append(matrix, line)
		for i := range records {
			val, _ := strconv.ParseFloat(records[i], 64)
			line[i] = val
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

	numCPU := runtime.NumCPU() // Получаем количество ядер процессора

	for o := 0; o < a1-1; o++ {
		max_el = -math.MaxFloat64 // reset max_el for each iteration 'o'

		// Поиск максимального элемента - можно параллелизировать, но для простоты оставим последовательно
		for i := o; i < a1; i++ {
			for j := o; j < a2-1; j++ {
				if math.Abs(matrix[i][j]) > max_el {
					max_el = matrix[i][j]
					mi = i
					mj = j
				}
			}
		}

		// Перенос строк и столбцов - последовательные операции
		for j := 0; j < a2; j++ {
			if mi == o && mi != o { //условие никогда не выполнится mi == o == o
				s1 = matrix[o][j] //перенос строк исходного массива
				matrix[o][j] = matrix[mi][j]
				matrix[mi][j] = s1
				s1 = x[o][j] //перенос строк массива коррекции
				x[o][j] = x[mi][j]
				x[mi][j] = s1
			}
		}

		for i := 0; i < a1; i++ {
			if mj == o && mj != o { //условие никогда не выполнится mj == o == o
				s1 = matrix[i][o] //перенос столбцов исходного массива
				matrix[i][o] = matrix[i][mj]
				matrix[i][mj] = s1
				s1 = x[i][o] //перенос столбцов массива коррекции
				x[i][o] = x[i][mj]
				x[i][mj] = s1
			}
		}

		// Нормировка строк - можно распараллелить по столбцам
		var wgNorm sync.WaitGroup
		for j := o; j < a2; j++ {
			wgNorm.Add(1)
			go func(col int) {
				defer wgNorm.Done()
				matrix[o][col] = matrix[o][col] / max_el
			}(j)
		}
		wgNorm.Wait()

		// Приведение к треугольному виду - распараллелим по строкам
		var wgTriang sync.WaitGroup
		s2Chan := make(chan float64)

		go func() { // Goroutine для расчета s2 чтобы избежать data race
			s2Chan <- matrix[o][o]
		}()
		s2 = <-s2Chan

		for i := o + 1; i < a1; i++ { // Начинаем с o+1, чтобы избежать деления на 0 в первом цикле
			wgTriang.Add(1)
			go func(row int) {
				defer wgTriang.Done()
				factor := matrix[row][o] / s2 // Вычисляем фактор для текущей строки
				for j := o; j < a2; j++ {
					matrix[row][j] = matrix[row][j] - factor*matrix[o][j]
				}
			}(i)
		}
		wgTriang.Wait()

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
	// Нормировка последней строки - можно распараллелить по столбцам
	var wgLastNorm sync.WaitGroup
	for j := a1 - 1; j < a2; j++ {
		wgLastNorm.Add(1)
		go func(col int) {
			defer wgLastNorm.Done()
			matrix[a1-1][col] = matrix[a1-1][col] / max_el
		}(j)
	}
	wgLastNorm.Wait()

	// Обратный ход, получение корней - распараллелим внутренний цикл
	var wgBackSub sync.WaitGroup
	for i := a1 - 1; i >= 0; i-- {
		wgBackSub.Add(1)
		go func(row int) {
			defer wgBackSub.Done()
			k[row] = matrix[row][a1]
			var sumChan = make(chan float64)

			go func() { // Горутина для расчета суммы чтобы избежать data race
				sum := float64(0)
				for j := 0; j < a1; j++ {
					if j != row {
						sum -= matrix[row][j] * k[j]
					}
				}
				sumChan <- sum
			}()
			sumVal := <-sumChan
			k[row] += sumVal // Добавляем расчитанную сумму к k[row]
			k[row] = k[row] / matrix[row][row]
		}(i)
	}
	wgBackSub.Wait()

	for i := 0; i < a1; i++ {
		for j := 0; j < a1; j++ {
			if i == j {
				ans[int(x[i][j])] = k[i]
			}
		}
	}

	elapsed := time.Since(start)
	fmt.Println("Количество ядер процессора:", numCPU)
	fmt.Println("Решение системы:")
	for i := 0; i < a1; i++ {
		fmt.Println(ans[i])
	}
	fmt.Println("Время выполнения:", elapsed)
}
