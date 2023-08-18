package main

import (
	"bufio"
	"fmt"
	"github.com/sqweek/dialog"
	"os"
	"strings"
)

func readCommandsFromFile(filePath string) ([]string, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open the file: %w", err)
	}
	defer file.Close()

	// 创建一个 Scanner 读取文件内容
	scanner := bufio.NewScanner(file)

	// 创建一个切片用于存放数据
	var lines []string

	// 逐行读取文件内容并存入切片
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	// 检查是否有读取错误
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading the file: %w", err)
	}

	return lines, nil
}

func getCommandLines() ([]string, error) {
	// 声明一个 commandLines 变量用来存储命令
	var commandLines []string

	scanner := bufio.NewScanner(os.Stdin)

	// 设置一个死循环
	for {
		fmt.Println("Input txt file path, for example E:\\\\Hello\\\\draw.txt")
		fmt.Println("Hint: Press ENTER button directly can select file in a GUI window.")
		fmt.Println("If you want to exit the program, simply type the letter \"q\" and press Enter.")

		// 得到用户输入的 input 字符串
		scanner.Scan()
		input := scanner.Text()

		// 如果输入为 "q"，则退出主程序
		if strings.ToLower(input) == "q" {
			fmt.Println("Exiting the program...")
			// 跳出循环
			break
		}
		// 对应与直接输入 Enter，如果输入 ENTER 则显示对话框，不会退出主程序
		if input == "" {
			filePath, err := dialog.File().Filter("All Files", "*.*").Title("Select a File").Load()
			// 如果没有选择文件，即选择取消，则打印提示信息，并重新进入循环
			if err != nil {
				fmt.Println("Hint: You did not select a file.", err)
				fmt.Println()
				continue
			}
			// 如果选择文件，返回这个 commandLines
			lines, err := readCommandsFromFile(filePath)
			if err != nil {
				continue
			}

			fmt.Println("File path:", filePath)
			commandLines = append(commandLines, lines...)
			// 跳出循环
			break
		} else {
			// 对应直接输入 txt 文件路径，
			// 通过命令行界面得到一个装有命令的 txt 文件
			lines, err := readCommandsFromFile(input)
			if err != nil {
				continue
			}

			fmt.Println("File path:", input)
			commandLines = append(commandLines, lines...)
			// 跳出循环
			break
		}
	}

	return commandLines, nil
}

func showHead() {
	welcome := `
====================================================================
==      ========================       =============================
=  ====  =======================  ====  ============================
=  ====  =======================  ====  ============================
==  =======    ====   ====   ===  ====  ==  =   ====   ===  =   =  =
====  =====  =  ==  =  ==  =  ==  ====  ==    =  ==  =  ==  =   =  =
======  ===  =  ==     ==  =====  ====  ==  ==========  ===   =   ==
=  ====  ==    ===  =====  =====  ====  ==  ========    ===   =   ==
=  ====  ==  =====  =  ==  =  ==  ====  ==  =======  =  ==== === ===
==      ===  ======   ====   ===       ===  ========    ==== === ===
====================================================================
`
	fmt.Println(welcome)
	fmt.Println()
}

func main() {
	// 显示程序头
	showHead()
	// 实现程序头逻辑
	commandLines, err := getCommandLines()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Command lines:", commandLines)
	fmt.Println("Program finished.")
}
