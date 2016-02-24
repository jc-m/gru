package main

import (
	"fmt"
	aurora "github.com/medallia/gru/aurora"
	"flag"
)

func main() {
	var sched *aurora.AuroraScheduler
	var proxy string
	var address string


	flag.StringVar(&proxy, "socks", "localhost:8888", "socks address:port")
	flag.StringVar(&address, "address", "localhost:8081", "Aurora Scheduler address:port")

	flag.Parse()

	uri:= fmt.Sprintf("http://%s/api",address)
	sched = aurora.NewSchedulerFactory(uri, proxy)


	summaries := sched.GetRoleSummary()
	fmt.Println(summaries)
	for summary, _ := range summaries {
		fmt.Println(sched.GetJobSummary(summary.GetRole()))
	}
	fmt.Println("--------------")
	fmt.Println(sched.GetJobs(""))
}