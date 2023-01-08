package utils

func MethodToChar(method string) string {
	switch method {
	case "GET":
		return "0"
	case "POST":
		return "1"
	case "PUT":
		return "2"
	case "DELETE":
		return "3"
	case "HEAD":
		return "4"
	case "OPTIONS":
		return "5"
	case "PATCH":
		return "6"
	case "TRACE":
		return "7"
	case "CONNECT":
		return "8"
	default:
		return "9"
	}
}

func MethodToByte(method string) byte {
	switch method {
	case "GET":
		return '0'
	case "POST":
		return '1'
	case "PUT":
		return '2'
	case "DELETE":
		return '3'
	case "HEAD":
		return '4'
	case "OPTIONS":
		return '5'
	case "PATCH":
		return '6'
	case "TRACE":
		return '7'
	case "CONNECT":
		return '8'
	default:
		return '9'
	}
}


func CheckMethodPath(routerPath, requestPath, method string) bool {
	if routerPath[0] == MethodToByte(method)&&routerPath[1:] == requestPath {
		return true
	}
	return false
}