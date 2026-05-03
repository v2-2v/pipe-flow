echo "Building input.go..."
go build -o ./build/ input.go

echo "Building data-bank.go..."
go build -o ./build/ data-bank.go

echo "Building data-print.go..."
go build -o ./build/ data-print.go

echo "Building data-send.go..."
go build -o ./build/ data-send.go

echo "Building data-clear.go..."
go build -o ./build/ data-clear.go