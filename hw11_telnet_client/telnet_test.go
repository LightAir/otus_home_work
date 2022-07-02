package main

import (
	"bytes"
	"io/ioutil"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out)

			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})

	t.Run("Send. Not init", func(t *testing.T) {
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		client := NewTelnetClient("localhost", time.Second*10, ioutil.NopCloser(in), out)
		err := client.Send()

		require.ErrorIs(t, err, errConnectionIsNotInit)
	})

	t.Run("Receive. Not init", func(t *testing.T) {
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		client := NewTelnetClient("localhost", time.Second*10, ioutil.NopCloser(in), out)
		err := client.Receive()

		require.ErrorIs(t, err, errConnectionIsNotInit)
	})

	t.Run("Close not init", func(t *testing.T) {
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		client := NewTelnetClient("localhost", time.Second*10, ioutil.NopCloser(in), out)
		err := client.Close()

		require.Nil(t, err)
	})

	t.Run("Close not init", func(t *testing.T) {
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		client := NewTelnetClient("localhost", time.Second*10, ioutil.NopCloser(in), out)
		err := client.Close()

		require.Nil(t, err)
	})

	t.Run("Invalid port", func(t *testing.T) {
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		client := NewTelnetClient("127.0.0.1:123456789", time.Second*10, ioutil.NopCloser(in), out)
		err := client.Connect()

		require.Equal(t, err.Error(), "dial tcp: address 123456789: invalid port")
	})

	t.Run("lose", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		client := NewTelnetClient(l.Addr().String(), time.Second, ioutil.NopCloser(in), out)
		require.NoError(t, client.Connect())

		err = client.Close()
		require.NoError(t, err)

		err = client.Close()
		require.Error(t, err)

		in.WriteString("hello\n")
		err = client.Send()
		require.Error(t, err)

		err = client.Receive()
		require.Error(t, err)
	})
}
