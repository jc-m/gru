package aurora

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/jc-m/gru/aurora/gen/api"
	"os"
	"net/url"
	"net/http"
	"golang.org/x/net/proxy"
)

type Scheduler interface {
	GetRoleSummary() map[*api.RoleSummary]bool
	getJobSummary(role string) map[*api.JobSummary]bool
	GetJobs(role string) map[*api.JobConfiguration]bool
}

type AuroraScheduler struct {
	schedulerUrl *url.URL
	client *api.AuroraSchedulerManagerClient
}

func NewSchedulerFactory(urlString string, proxyString string) *AuroraScheduler {

	var protocolFactory thrift.TProtocolFactory
	var sched AuroraScheduler

	parsedUrl, err := url.Parse(urlString)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing URL: ", err)
		return nil
	}

	sched.schedulerUrl = parsedUrl

	// setup a http client
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}

	if len(proxyString) > 0 {
		dialer, err := proxy.SOCKS5("tcp", "localhost:9999", nil, proxy.Direct)
		if err != nil {
			fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
			os.Exit(1)
		}

		// set our socks5 as the dialer
		httpTransport.Dial = dialer.Dial
	}

	transport, err := thrift.NewTHttpPostClientWithOptions(parsedUrl.String(), thrift.THttpClientOptions{Client:httpClient})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error Creating client: ", err)
		return nil
	}
	protocolFactory = thrift.NewTJSONProtocolFactory()

	sched.client = api.NewAuroraSchedulerManagerClientFactory(transport, protocolFactory)

	return &sched
}

func (a *AuroraScheduler) GetRoleSummary() map[*api.RoleSummary]bool {
	res, err := a.client.GetRoleSummary()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting Role Summary ", err)
		return make(map[*api.RoleSummary]bool)
	}
	summaries := res.GetResult_().GetRoleSummaryResult_().GetSummaries()
	return summaries
}

func (a *AuroraScheduler) GetJobSummary(role string) map[*api.JobSummary]bool {
	res, err := a.client.GetJobSummary(role)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting Job Summary ", err)
		return make(map[*api.JobSummary]bool)
	}
	summaries := res.GetResult_().GetJobSummaryResult_().GetSummaries()
	return summaries
}

func (a *AuroraScheduler) GetJobs(role string) map[*api.JobConfiguration]bool {
	res, err := a.client.GetJobs(role)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting Jobs ", err)
		return make(map[*api.JobConfiguration]bool)
	}
	configs := res.GetResult_().GetGetJobsResult_().GetConfigs()
	return configs
}