package mocks

import (
	"io"
	"net"
	"time"
)

type MockUDPConn struct {
	buffer       []byte
	readFromFunc func(b []byte) (int, *net.UDPAddr, error)
	writeToFunc  func(b []byte, addr *net.UDPAddr) (int, error)
}

func NewMockUDPConn() *MockUDPConn {
	return &MockUDPConn{
		buffer: make([]byte, 0),
		readFromFunc: func(b []byte) (int, *net.UDPAddr, error) {
			return 0, &net.UDPAddr{}, io.EOF
		},
		writeToFunc: func(b []byte, addr *net.UDPAddr) (int, error) {
			return len(b), nil
		},
	}
}

func (m *MockUDPConn) ReadFromUDP(b []byte) (int, *net.UDPAddr, error) {
	return m.readFromFunc(b)
}

func (m *MockUDPConn) WriteToUDP(b []byte, addr *net.UDPAddr) (int, error) {
	return m.writeToFunc(b, addr)
}

func (m *MockUDPConn) Close() error {
	return nil
}

func (m *MockUDPConn) SetReadBuffer(bytes int) error {
	return nil
}

func (m *MockUDPConn) SetWriteBuffer(bytes int) error {
	return nil
}

func (m *MockUDPConn) LocalAddr() net.Addr {
	return &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 2053}
}

func (m *MockUDPConn) RemoteAddr() net.Addr {
	return &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 53000}
}

func (m *MockUDPConn) SetDeadline(t time.Time) error {
	return nil
}

func (m *MockUDPConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *MockUDPConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (m *MockUDPConn) Read(b []byte) (n int, err error) {
	return 0, io.EOF
}

func (m *MockUDPConn) Write(b []byte) (n int, err error) {
	return len(b), nil
}

func (m *MockUDPConn) ReadFrom(b []byte) (int, net.Addr, error) {
	n, addr, err := m.ReadFromUDP(b)
	return n, addr, err
}

func (m *MockUDPConn) WriteTo(b []byte, addr net.Addr) (n int, err error) {
	udpAddr, ok := addr.(*net.UDPAddr)
	if !ok {
		udpAddr = &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 53000}
	}
	return m.WriteToUDP(b, udpAddr)
}

func (m *MockUDPConn) File() (f interface{}, err error) {
	return nil, nil
}
