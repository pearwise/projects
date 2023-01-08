mkdir $1

cd $1

echo -e 'package main

func main() {

}' > main.go

go mod init

code .