/*
A simple test for web server status. This package is intended for use with Nagios.
*/
package main

import (
	"fmt"
	"net/http"
	"flag"
	"os"
)

type NagiosStatusVal int

const (
	NAGIOS_OK NagiosStatusVal = iota
	NAGIOS_WARNING
	NAGIOS_CRITICAL
	NAGIOS_UNKNOWN
)

var (
	valMessages = []string{
		"OK:",
		"WARNING:",
		"CRITICAL:",
		"UNKNOWN:",
	}
)

// Take a bunch of NagiosStatus pointers and find the highest value, then
// combine all the messages. Things win in the order of highest to lowest.
type NagiosStatus struct {
	Message string
	Value   NagiosStatusVal
}


func (status *NagiosStatus) Aggregate(otherStatuses []*NagiosStatus) {
	for _, s := range otherStatuses {
		if status.Value < s.Value {
			status.Value = s.Value
		}

		status.Message += " - " + s.Message
	}
}

// Exit with an UNKNOWN status and appropriate message
func Unknown(output string) {
	ExitWithStatus(&NagiosStatus{output, NAGIOS_UNKNOWN})
}

// Exit with an CRITICAL status and appropriate message
func Critical(output string) {
	ExitWithStatus(&NagiosStatus{output, NAGIOS_CRITICAL})
}

// Exit with an WARNING status and appropriate message
func Warning(output string) {
	ExitWithStatus(&NagiosStatus{output, NAGIOS_WARNING})
}

// Exit with an OK status and appropriate message
func Ok(output string) {
	ExitWithStatus(&NagiosStatus{output, NAGIOS_OK})
}

// Exit with a particular NagiosStatus
func ExitWithStatus(status *NagiosStatus) {
	fmt.Fprintln(os.Stdout, valMessages[status.Value], status.Message)
	os.Exit(int(status.Value))
}


// main expects one parameter on the command line: a valid website name.
// This host is called using https, and returns OK and the status if the status is 200, or
// Critical and the status if it's anything else.
func main() {

	hostPtr := flag.String("host", "somedomain.com", "A valid internet site without http:// or https://")
	protocolPtr := flag.String("protocol", "https", "Protocol - either https or http")

	flag.Parse()

	url := *protocolPtr + "://" +  *hostPtr

	resp, err := http.Get(url)

	if err != nil {
		msg := "CRITICAL- host did not respond"
		Critical(msg)

	} else {
		if resp.StatusCode != 200 {
			msg := "CRITICAL- " + *hostPtr + " " + resp.Proto + " " + resp.Status
			Critical(msg)
		} else {
			msg := "OK- " + *hostPtr + " responded with " + resp.Proto + " " + resp.Status
			Ok(msg)
		}
	}

}
