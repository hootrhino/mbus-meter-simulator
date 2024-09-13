package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const (
	slaveID = 1       // 水表从机的ID
	port    = ":3000" // 监听的端口号
)

// 模拟的水表数据
type WaterMeterData struct {
	Consumption float64 // 消耗量（立方米）
	Temperature float64 // 温度（摄氏度）
}

// 发送模拟的水表数据
func sendWaterMeterData(conn net.Conn) {
	data := WaterMeterData{
		Consumption: 123.45,
		Temperature: 22.5,
	}

	// 将数据转换为字符串格式，实际应用中可能需要转换为特定的M-Bus帧格式
	dataStr := fmt.Sprintf("ID: %d, Consumption: %.2f m³, Temperature: %.1f°C\n", slaveID, data.Consumption, data.Temperature)
	conn.Write([]byte(dataStr))
}

func main() {
	// 监听指定端口
	l, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on port", port)
	go func() {
		for {
			// 接受主设备连接
			conn, err := l.Accept()
			if err != nil {
				fmt.Println("Error accepting: ", err.Error())
				os.Exit(1)
			}

			// 模拟处理连接
			go handleRequest(conn)
		}
	}()
	go func() {
		for {
			slaveAddress := "localhost:3000" // 水表从机的地址和端口
			response, err := RequestMeter(slaveAddress)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Received data from slave:", response)
			}
			time.Sleep(5 * time.Second)
		}
	}()
	select {}
}

// 处理来自主设备的请求
func handleRequest(conn net.Conn) {
	defer conn.Close()

	// 假设我们接收到一个简单的请求，然后发送水表数据
	buffer := make([]byte, 1024)
	N, _ := conn.Read(buffer)
	fmt.Println(buffer[:N])
	// 模拟从机的响应延迟
	time.Sleep(2 * time.Second)

	// 发送模拟的水表数据
	sendWaterMeterData(conn)
}

// RequestMeter 尝试连接到水表从机，发送查询请求，并返回响应。
// 参数 address 是从机的TCP地址，例如 "localhost:3000"。
// 返回值是字符串形式的响应数据，如果发生错误则返回错误信息。
func RequestMeter(address string) (string, error) {
	// 连接到水表从机
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return "", fmt.Errorf("error connecting to slave: %v", err)
	}
	defer conn.Close()

	// 发送查询请求
	query := fmt.Sprintf("Request data from slave ID %d\n", 1)
	_, err = conn.Write([]byte(query))
	if err != nil {
		return "", fmt.Errorf("error sending query: %v", err)
	}

	// 读取从机的响应
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	// 去除响应字符串末尾的换行符
	response = strings.TrimSpace(response)
	return response, nil
}
