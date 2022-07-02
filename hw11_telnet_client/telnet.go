package main

import (
	"errors"
	"io"
	"net"
	"time"
)

var errConnectionIsNotInit = errors.New("connection is closed")

type TelnetClient interface {
	Connect() error
	Close() error
	Send() error
	Receive() error
}

type TClient struct {
	address string
	timeout time.Duration
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

func (t *TClient) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return err
	}

	t.conn = conn

	return nil
}

func (t *TClient) Close() error {
	if t.conn == nil {
		return nil
	}

	if err := t.conn.Close(); err != nil {
		return err
	}

	return nil
}

func (t *TClient) Send() error {
	if t.conn == nil {
		return errConnectionIsNotInit
	}

	if _, err := io.Copy(t.conn, t.in); err != nil {
		return err
	}

	return nil
}

func (t *TClient) Receive() error {
	if t.conn == nil {
		return errConnectionIsNotInit
	}

	if _, err := io.Copy(t.out, t.conn); err != nil {
		return err
	}

	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}
