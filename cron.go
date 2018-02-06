package main

import "time"
import "log"

var hostmap map[string]string

func GetHost(){
	hostmap = make(map[string]string)
	sql := "select hostname,ip from host where ip = \"\""
	rows, err := DB.Query(sql)
	if err != nil {
		Logger().Println("ERROR:", err)
		return
	}

	for rows.Next() {
		var (
			hostname string
			ip string
		)

		err = rows.Scan(&hostname, &ip)
		if err != nil {
			Logger().Println("ERROR:", err)
			continue
		}

		hostmap[ip] = hostname
		log.Println(hostname,":",ip)
	}

	defer rows.Close()
}

func CronHost(){
	for {
		time.Sleep(time.Duration(GetConfig().Interval) * time.Second)
		GetHost()
	}
}

func GetHostMap() map[string]string {
	return hostmap
}
