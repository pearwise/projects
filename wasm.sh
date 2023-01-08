mkdir $1

cd $1

mkdir .vscode

echo -e '{
    "go.toolsEnvVars": {
        "GOARCH": "wasm",
        "GOOS": "js"
    },
    "go.testEnvVars": {
        "GOARCH": "wasm",
        "GOOS": "js"
    },
    "go.installDependenciesWhenBuilding": false
}' > .vscode/settings.json

echo -e 'package main

import "syscall/js"

var (
    document = js.Global().Get("document")
)

func main() {

}' > main.go

go mod init

echo 'GOOS=js GOARCH=wasm go build -o wasm/main.wasm

cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" wasm' > run.sh

chmod 777 run.sh

mkdir cmd

cd cmd

echo -e 'package main

func main() {

}' > main.go

code ..