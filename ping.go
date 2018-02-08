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
	//"fmt"
	//"fmt"
	"github.com/open-falcon/common/model"
	"fmt"
	"os"
	"encoding/json"
	"bytes"
	"net/http"
)

func ping(){
	output := make(map[string]int)
	ipmap := GetHostMap()
	pinger := fastping.NewPinger()

	for k,_ := range GetHostMap() {
		output[k] = 0
		pinger.AddIP(k)
	}
	pinger.OnRecv = func(addr *net.IPAddr, rtt time.Duration){
		Logger().Println("success:",addr.String())
		output[addr.String()] = 1
	}
	pinger.OnIdle = func() {
		Logger().Println("finish loop")
	}

	pinger.MaxRTT = time.Second * 5

	for i:=0;i<2;i++ {
		err := pinger.Run()
		if err != nil {
			Logger().Println(err.Error())
		}
	}

	/*for k,v := range output{
		if v == 0{
			fmt.Println(k)
		}
	}*/

	metrics, err := formatMetric(output,ipmap)
	if err != nil {
		Logger().Println(err.Error())
	}

	PostToAgent(metrics)

}

func formatMetric(output map[string]int,ipmap map[string]string)(metrics []*model.MetricValue, err error){
	hostname,err := os.Hostname()
	if err != nil {
		return metrics,err
	}
	for k,v := range output {
		tags := fmt.Sprintf("ip=%s,hostname=%s",k,ipmap[k])
		now := time.Now().Unix()
		singleMetric := &model.MetricValue{
			Endpoint:  hostname,
			Metric:    "host.ping",
			Value:     v,
			Timestamp: now,
			Step:      60,
			Type:      "GAUGE",
			Tags:      tags,
		}
		metrics = append(metrics, singleMetric)
	}

	return metrics,nil
}

func PostToAgent(metrics []*model.MetricValue) {
	if len(metrics) == 0 {
		return
	}

	contentJson, err := json.Marshal(metrics)
	if err != nil {
		Logger().Println("Error for PostToAgent json Marshal: ", err)
		return
	}
	contentReader := bytes.NewReader(contentJson)
	req, err := http.NewRequest("POST", GetConfig().Agent_path, contentReader)
	if err != nil {
		Logger().Println("Error for PostToAgent in NewRequest: ", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Logger().Println("Error for PostToAgent in http client Do: %v", err)
		return
	}
	defer resp.Body.Close()

	Logger().Println("<= ", resp.Body)
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


