package utils

func CalculateScore(whiteNum, blackNum, whiteChessIsConnect, blackChessIsConnect uint8) (score uint32) {
	// 计算白棋得分
	if whiteChessIsConnect < 2 {
		switch whiteNum {
		case 0:
		case 1:
			if whiteChessIsConnect==0 {
				score += 10
			}
		case 2:
			if whiteChessIsConnect==0 {
				score += 50
			} else {
				score += 25
			}
		case 3:
			if whiteChessIsConnect==0 {
				score += 10000
			} else {
				score += 55
			}
		default:
			score += 30000
		}
	}

	// 计算黑棋得分
	if blackChessIsConnect < 2 {
		switch blackNum {
		case 0:
		case 1:
			if blackChessIsConnect==0 {
				score += 10
			}
		case 2:
			if blackChessIsConnect==0 {
				score += 40
			} else {
				score += 30
			}
		case 3:
			if blackChessIsConnect==0 {
				score += 200
			} else {
				score += 60
			}
		default:
			score += 20000
		}
	}
	return
}

func GetChessColor(chess uint8) string {
	switch chess {
	case 1:
		return "white"
	case 2:
		return "black"
	default:
		return ""
	}
}