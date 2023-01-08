package main

import (
	"fmt"
	"math"
	"math/rand"
	"syscall/js"
	"time"
	"wasm/utils"
)

var (
    document = js.Global().Get("document")
    alert = js.Global().Get("alert")
    chess   = document.Call("getElementById", "chess")
    ctx = chess.Call("getContext","2d")
    curRole = document.Call("getElementById", "curRole")
    //restart = document.Call("getElementById", "restart")
    // 1:白棋,2:黑棋,other:空
    chessBoard [15][15]uint8
    //当前下棋的棋子类型1:白棋,2:黑棋
    curChess uint8
    // //当前玩家的棋子类型1:白棋,2:黑棋
    playerChess uint8
    //游戏是否结束
    isOver bool

    console = js.Global().Get("console")
)

func main() {
    // 画棋盘
    drawChessBoard(ctx)

    // 设置开局
    setStart()

    // 设置玩家点击事件
    setOnClick()

    select{}
}

// 画棋盘
func drawChessBoard(ctx js.Value) {
    ctx.Set("strokeStyle","#b9b9b9")
    for i:=0;i<15;i++ {
        //设置水平方向的起始点
        ctx.Call("moveTo", 15, 15+i*30)
        //设置水平方向的终止点
        ctx.Call("lineTo", 435, 15+i*30)
        //连接水平方向的两点
        ctx.Call("stroke")
        //设置竖直方向的起始点
        ctx.Call("moveTo", 15+i*30, 15)
        //设置竖直方向的终止点
        ctx.Call("lineTo", 15+i*30, 435)
        //连接竖直方向的两点
        ctx.Call("stroke")
    }
}

// 开局设置
func setStart() {
    rand.Seed(time.Now().Unix())
    playerChess = uint8(rand.Intn(2))+1
    curChess = 2
    curRole.Set("innerHTML", "黑棋回合")
    playerChessText := document.Call("getElementById", "playerChess")
    if playerChess == 1 {
        downChess(7, 7, "black")
        flushCurRole()
        playerChessText.Set("innerHTML", "我方执白棋")
    } else {
        playerChessText.Set("innerHTML", "我方执黑棋")
    }
}

// 设置玩家点击事件
func setOnClick() {
    chess.Set("onclick", js.FuncOf(func(this js.Value, args []js.Value) any {
        if isOver {
            console.Call("log", "game is over")
            alert.Invoke("游戏已结束")
            return nil
        }
        x := int8(args[0].Get("offsetX").Float()/30)
        y := int8(args[0].Get("offsetY").Float()/30)
        if chessBoard[x][y] != 0 {
            return nil
        }

        downChess(x, y, utils.GetChessColor(curChess))
        
        if checkOver(x, y) {
            isOver = true
            if curChess == 1 {
                alert.Invoke("白棋获胜")
                curRole.Set("innerHTML", "白棋获胜")
                return nil
            } else {
                alert.Invoke("黑棋获胜")
                curRole.Set("innerHTML", "黑棋获胜")
                return nil
            }
        }
        
        flushCurRole()

        // ai下棋
        aiDownChess()    
        
        flushCurRole()

        return nil
    }))
}

// 下棋
func downChess(x, y int8, color string) {
    console.Call("log", fmt.Sprintf("下棋的坐标:(%d, %d), 棋子的颜色: %s\n", x, y, color))
    chessBoard[x][y] = curChess

    // 画棋
    ctx.Call("beginPath")
    ctx.Call("arc", 15+uint16(x)*30, 15+uint16(y)*30, 13, 0, 2*math.Pi)
    ctx.Call("closePath")
    ctx.Set("fillStyle", color)
    ctx.Call("fill")
}

// 更新玩家回合
func flushCurRole() {
    if curChess == 1 {
        curChess = 2
    } else {
        curChess = 1
    }
    if curChess == 1 {
        curRole.Set("innerHTML", "白棋回合")
    } else {
        curRole.Set("innerHTML", "黑棋回合")
    }
}

// 检查游戏是否结束
func checkOver(x, y int8) bool {
    var count byte
    var i int8 = 1
    for i<5&&y>=i&&chessBoard[x][y-i]==curChess {
        count++
        i++
    }
    i = 1
    for i<5&&y+i<15&&chessBoard[x][y+i]==curChess {
        count++
        i++
    }
    if count > 3 {
        return true
    }
    // |
    count, i = 0, 1
    for i<5&&x>=i&&chessBoard[x-i][y]==curChess {
        count++
        i++
    }
    i = 1
    for i<5&&x+i<15&&chessBoard[x+i][y]==curChess {
        count++
        i++
    }
    if count > 3 {
        return true
    }
    // \
    count, i = 0, 1
    for i<5&&x>=i&&y>=i&&chessBoard[x-i][y-i]==curChess {
        count++
        i++
    }
    i = 1
    for i<5&&x+i<15&&y+i<15&&chessBoard[x+i][y+i]==curChess {
        count++
        i++
    }
    if count > 3 {
        return true
    }
    // /
    count, i = 0, 1
    for i<5&&x+i<15&&y>=i&&chessBoard[x+i][y-i]==curChess {
        count++
        i++
    }
    i = 1
    for i<5&&x>=i&&y+i<15&&chessBoard[x-i][y+i]==curChess {
        count++
        i++
    }
    return count > 3
}

// ai下棋
func aiDownChess() {
    x, y := calculate()
    console.Call("log", fmt.Sprintf("ai最佳坐标(%d,%d)", x, y))
    downChess(x, y, utils.GetChessColor(curChess))
    if checkOver(x, y) {
        isOver = true
        if curChess == 1 {
            alert.Invoke("白棋获胜")
            curRole.Set("innerHTML", "白棋获胜")
        } else {
            alert.Invoke("黑棋获胜")
            curRole.Set("innerHTML", "黑棋获胜")
        }
    }
}

// 计算棋盘中得分最高的点
func calculate() (int8, int8) {
    var maxScore, curScore uint32
    var i, j, x, y int8
    for i < 15 {
        j = 0
        for j < 15 {
            curScore = 0
            // 判断该点是否可以落子
            if chessBoard[i][j] != 0 {
                j++
                continue
            }

            // -

            curScore += calculateDirection(1, 0, i, j)

            // |
            
            curScore += calculateDirection(0, 1, i, j)
            
            // /
            
            curScore += calculateDirection(-1, 1, i, j)

            // \
        
            curScore += calculateDirection(1, 1, i, j)
            if curScore!=0 {
                console.Call("log", fmt.Sprintf("当前坐标(%d, %d)得分: %d", i, j, curScore))
            }
            if curScore > maxScore {
                maxScore, x, y = curScore, i, j
            }
            j++
        }
        i++
    }
    return x, y
}

// 计算某个方向的得分
func calculateDirection(i, j, x, y int8) uint32 {
    var k int8 = 1
    var playerNum, aiNum, playerChessIsLive, aiChessIsLive uint8
    // 正向计算player num
    for k < 5 {
        if x-i*k < 0 || y-j*k < 0 || x-i*k > 14 || y-j*k > 14 {
            break
        }
        switch chessBoard[x-i*k][y-j*k] {
            case playerChess:
            // 玩家
            playerNum++
            case curChess:
            // ai
            playerChessIsLive++
            goto NEXT0
            default:
            goto NEXT0
        }
        k++
    }
    NEXT0:
    k = 1
    // 反向计算player num
    for k < 5 {
        if x+i*k < 0 || y+j*k < 0 || x+i*k > 14 || y+j*k > 14 {
            break
        }
        switch chessBoard[x+i*k][y+j*k] {
            case playerChess:
            // 玩家
            playerNum++
            case curChess:
            // ai
            playerChessIsLive++
            goto NEXT1
            default:
            goto NEXT1
        }
        k++
    }
    NEXT1:
    k = 1
    // 正向计算ai num
    for k < 5 {
        if x-i*k < 0 || y-j*k < 0 || x-i*k > 14 || y-j*k > 14 {
            break
        }
        switch chessBoard[x-i*k][y-j*k] {
            case curChess:
            // ai
            aiNum++
            case playerNum:
            // 玩家
            aiChessIsLive++
            goto NEXT2
            default:
            goto NEXT2
        }
        k++
    }
    NEXT2:
    k = 1
    // 反向计算ai num
    for k < 5 {
        if x+i*k < 0 || y+j*k < 0 || x+i*k > 14 || y+j*k > 14 {
            break
        }
        switch chessBoard[x+i*k][y+j*k] {
            case curChess:
            // ai
            aiNum++
            case playerChess:
            // 玩家
            aiChessIsLive++
            goto NEXT3
            default:
            goto NEXT3
        }
        k++
    }
    NEXT3:
    if playerChess == 1 {
        // 玩家是白棋
        console.Call("log", fmt.Sprintf("坐标(%d,%d), whiteNum=%d, blackNum=%d, whiteIsLive=%d, blackIsLive=%d", x, y, playerNum, aiNum, playerChessIsLive, aiChessIsLive))
        return utils.CalculateScore(playerNum, aiNum, playerChessIsLive, aiChessIsLive)
    } else {
        // 玩家是黑棋
        console.Call("log", fmt.Sprintf("坐标(%d,%d), whiteNum=%d, blackNum=%d, whiteIsLive=%d, blackIsLive=%d", x, y, aiNum, playerNum, aiChessIsLive, playerChessIsLive))
        return utils.CalculateScore(aiNum, playerNum, aiChessIsLive, playerChessIsLive)
    }
}