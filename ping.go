package main

import (
	//"github.com/tatsushid/go-fastping"
	"bufio"
	"io"
	"strings"
	"os/exec"
	"github.com/tatsushid/go-fastping"
	"net"
	"time"
)

func Ping(){
	//output := make(map[string]string)
	//ipmap := GetHostMap()
	pinger := fastping.NewPinger()
	for k,_ := range GetHostMap() {
		pinger.AddIP(k)
	}
	pinger.OnRecv = func(addr *net.IPAddr, rtt time.Duration){
		Logger().Println("success:",addr.String())
	}
	pinger.OnIdle = func() {
		Logger().Println("Idle")
	}
	err := pinger.Run()
	if err != nil {
		Logger().Println(err.Error())
	}
}

//由于未知原因，部分机器在数据库无ip
func nslookup(hostname string) (string, error) {
	var ip string

	command := "nslookup"
	param := []string{hostname}

	cmd,reader,err := ExecCommand(command, param)
	if err != nil {
		Logger().Println(err.Error())
	}

	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		line = strings.TrimSpace(line)
		if line == ""{
			continue
		}

		info := strings.Split(line, " ")
		if strings.Contains(line, "Address") && len(info) == 2 {
			ip = info[1]
			break
		}
	}
	cmd.Wait()
	return ip ,nil
}

func ExecCommand(commandName string, params []string) (*exec.Cmd, *bufio.Reader,error) {
	cmd := exec.Command(commandName, params...)

	//显示运行的命令
	//fmt.Println(cmd.Args)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		Logger().Println(err)
		return nil,nil,err
	}

	cmd.Start()

	reader := bufio.NewReader(stdout)

	return cmd,reader,nil


}


