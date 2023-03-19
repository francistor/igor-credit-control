package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/francistor/igor/core"
	"github.com/francistor/igor/router"
)

func main() {

	// defer profile.Start(profile.BlockProfile).Stop()

	// After ^C, signalChan will receive a message
	doneChan := make(chan struct{}, 1)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		close(doneChan)
		fmt.Println("terminating server")
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Get the command line arguments
	bootPtr := flag.String("boot", "resources/searchRules.json", "File or http URL with Configuration Search Rules")
	instancePtr := flag.String("instance", "", "Name of instance")

	flag.Parse()

	// Initialize the Config Object
	_ = core.InitPolicyConfigInstance(*bootPtr, *instancePtr, nil, true)

	// Get logger
	logger := core.GetLogger()

	// Start Diameter
	r := router.NewDiameterRouter(*instancePtr, RequestHandler)
	logger.Info("Diameter router started")

	// Start server
	r.Start()

	// Wait for termination signal
	<-doneChan

	// Close router gracefully
	r.Close()
}

// Handler for all requests
func RequestHandler(request *core.DiameterMessage) (*core.DiameterMessage, error) {

	hl := core.NewHandlerLogger()
	l := hl.L
	l.Info("")
	l.Info("----- Start Message Processing -----")

	defer func(hl *core.HandlerLogger) {
		hl.L.Info("----- Finished Message Processing -----")
		hl.WriteLog()
	}(hl)

	l.Infof("%s %s", "request", request)

	// The service to be sent as Charging-Rule-Base-Name is prepended the module 10 of the sum of the bytes of the
	// subscription data

	subscriptionDataAVP, err := request.GetAVPFromPath("Subscription-Id.Subscription-Id-Data")
	if err != nil {
		return nil, fmt.Errorf("Subscription-Id-Data not found: %s", err)
	}
	var sum int
	for b := range []byte(subscriptionDataAVP.GetString()) {
		sum += b
	}
	sum = sum % 10

	response := core.NewDiameterAnswer(request).Add("3GPP-Charging-Rule-Install", []core.DiameterAVP{
		*core.BuildDiameterAVP("3GPP-Charging-Rule-Base-Name", fmt.Sprintf("service-%d", sum)),
	})

	response.Add("Result-Code", core.DIAMETER_SUCCESS)
	l.Infof("%s %s", "response", response)

	return response, nil
}
