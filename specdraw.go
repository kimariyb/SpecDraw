package main

/*
@Author: Kimariyb
@Institution: XiaMen University
@Data: 2023-08-26
*/

import (
	"bufio"
	"fmt"
	"github.com/sqweek/dialog"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func batchExecution(commandLines []string) {
	// 获取当前工作目录的路径
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error:", err)
	}

	// 读取当前文件夹的文件列表
	files, err := ioutil.ReadDir(currentDir)
	if err != nil {
		log.Fatal("Error:", err)
	}

	// 遍历文件列表
	for _, file := range files {
		// 检查文件是否为 TOML 文件
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".toml") {
			fmt.Println("TOML file:", file.Name())

			// 进入 KimariDraw 程序
			cmd := exec.Command("kimaridraw.exe", file.Name())
			// 设置终端的编码为 UTF-8
			cmd.Env = append(cmd.Env, "PYTHONIOENCODING=utf-8")
			// 设置 HOME 环境变量
			cmd.Env = append(os.Environ(), "HOME="+currentDir)

			// 获取标准输入管道
			stdin, err := cmd.StdinPipe()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			// 启动进程
			e := cmd.Start()
			if err != nil {
				fmt.Println("Error starting process:", e)
				return
			}

			// 执行命令
			for _, line := range commandLines {
				// 进入 KimariDraw 程序后以此执行 txt 命令中的内容
				_, err := io.WriteString(stdin, line+"\n")
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
			}

			// 关闭标准输入管道
			stdin.Close()

			// 等待进程完成
			err = cmd.Wait()
			if err != nil {
				fmt.Println("Error:", err)
			}

		}
	}

}

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

			fmt.Println("Selected file path:", filePath)
			fmt.Println()
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
	// 获取当前文件的绝对路径
	filePath, err := filepath.Abs(os.Args[0])
	if err != nil {
		fmt.Println("获取文件路径失败:", err)
		return
	}

	// 获取最后修改时间的时间戳
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fmt.Println("获取文件信息失败:", err)
		return
	}
	modTime := fileInfo.ModTime()
	timestamp := modTime.Format("2006-Jan-02")

	fmt.Println("SpecDraw -- A script that automatically calls KimariDraw to generate pictures in batches")
	fmt.Println("Version: v1.1.0, release date:", timestamp)
	fmt.Println("Developer: Kimariyb, Ryan Hsiun")
	fmt.Println("Address: XiaMen University, School of Electronic Science and Engineering")
	fmt.Println("Website: https://github.com/kimariyb/SpecDraw")

	// 获取当前日期和时间
	now := time.Now().Format("Jan-02-2006, 15:04:05")

	// 输出版权信息和问候语
	fmt.Printf("(Copyright 2023 Kimariyb. Currently timeline: %s)\n", now)
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
	batchExecution(commandLines)
	fmt.Println()
	fmt.Println("Program finished.")
}
