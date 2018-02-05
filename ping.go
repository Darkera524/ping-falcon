package main

import (
	//"github.com/tatsushid/go-fastping"
	"bufio"
	"io"
	"strings"
	"os/exec"
	"fmt"
)

func Ping(){
	//output := make(map[string]string)
	//ipmap := make(map[string]string)
	//pinger := fastping.NewPinger()
	for k,v := range GetHostMap() {
		ip := v
		//if v == "" {
		insip,err := nslookup(k)
		if err != nil {
			Logger().Println(err.Error())
		}
		ip = insip
		//}
		fmt.Println(k,":",ip)
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


