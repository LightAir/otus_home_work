package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	timeout := flag.Duration("timeout", time.Second*10, "timeout")
	flag.Parse()

	host := flag.Arg(0)
	port := flag.Arg(1)

	address := host + ":" + port

	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)

	err := client.Connect()
	if err != nil {
		logrus.WithError(err).Error("connection error")
		return
	}

	fmt.Fprintf(os.Stderr, "...Connected to %s \n", address)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		err := client.Send()
		if err != nil {
			logrus.WithError(err).Error("sending error")
		}

		fmt.Fprintln(os.Stderr, "...EOF")
		cancel()
	}()

	go func() {
		err := client.Receive()
		if err != nil {
			logrus.WithError(err).Error("receiving error")
		}

		fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
		cancel()
	}()

	<-ctx.Done()

	err = client.Close()
	if err != nil {
		logrus.WithError(err).Error("connection closing error")
	}
}
