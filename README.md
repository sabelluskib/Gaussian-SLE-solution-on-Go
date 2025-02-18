# Gaussian-SLE-solution-on-Go
Solution of systems of linear algebraic equations by Gauss method with choice of the principal element

# How to run the program
## single_thread

1. `sudo apt-get install golang`
2. Generate the right amount of input using generate.go by fixing the source code
3. Put the generated file into the code of the program `gauss_single_thread.go`
4. `go run gauss_single_thread.go`

## multi_thread

1. Put the generated file into the code of the program `gauss_multi_thread.go`
2. `go run gauss_multi_thread.go`

# Results

On Intel(R) Xeon(R) CPU E5-2660 v2 @ 2.20GHz in single-threaded mode a system of 3000 unknowns is solved in 1 minute 42 seconds and on 20 threads in 43 seconds
