package irmagician

import (
	"fmt"
	"strconv"
	"time"

	"github.com/tarm/serial"
)

const (
	DefaultTimeout  = 1 * time.Second
	DefaultBaudRate = 9600
	DefaultWait     = 1 * time.Second
	BufferSize      = 640
)

type IrMagician struct {
	c *serial.Config
	s *serial.Port
}

func NewIrMagician(name string, rate int, timeout time.Duration) (*IrMagician, error) {
	var c *serial.Config
	if timeout == 0 {
		c = &serial.Config{
			Name: name,
			Baud: rate,
		}
	} else {
		c = &serial.Config{
			Name:        name,
			Baud:        rate,
			ReadTimeout: timeout,
		}
	}
	s, err := serial.OpenPort(c)
	if err != nil {
		return nil, err
	}
	return &IrMagician{c, s}, nil
}

func (ir *IrMagician) writeread(command string, waitsec int) ([]byte, error) {
	_, err := ir.s.Write([]byte(command))
	if err != nil {
		return nil, err
	}
	if waitsec != 0 {
		time.Sleep(time.Duration(waitsec) * time.Second)
	}
	buf := make([]byte, BufferSize)
	n, err := ir.s.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

func (ir *IrMagician) BankSet(n int) error {
	if n > 9 || n < 0 {
		return fmt.Errorf("BankSet: %v is out of Bank range (0-9)", n)
	}
	_, err := ir.s.Write([]byte("b," + strconv.Itoa(n) + "\r\n"))
	if err != nil {
		return err
	}
	return nil
}

func (ir *IrMagician) Capture() ([]byte, error) {
	_, err := ir.s.Write([]byte("c\r\n"))
	if err != nil {
		return nil, err
	}
	time.Sleep(DefaultWait)
	buf := make([]byte, BufferSize)
	n, err := ir.s.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

func (ir *IrMagician) Dump(n int) ([]byte, error) {
	if n > 63 || n < 0 {
		return nil, fmt.Errorf("Dump: %v is out of memory number range (0-63)", n)
	}
	resp, err := ir.writeread("d,"+strconv.Itoa(n)+"\r\n", 0)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (ir *IrMagician) Information(param int) ([]byte, error) {
	if param > 7 || param < 0 {
		return nil, fmt.Errorf("Information: %v is out of parameter range (0-7)", param)
	}
	resp, err := ir.writeread("I,"+strconv.Itoa(param)+"\r\n", 1)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (ir *IrMagician) SetPostScaler(v int) ([]byte, error) {
	if v > 255 || v < 1 {
		return nil, fmt.Errorf("SetPostScaler: %v is out of value range (1-255)", v)
	}
	resp, err := ir.writeread("k,"+strconv.Itoa(v)+"\r\n", 0)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (ir *IrMagician) LED(on bool) ([]byte, error) {
	toggle := "0"
	if on {
		toggle = "1"
	}
	resp, err := ir.writeread("L,"+toggle+"\r\n", 0)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (ir *IrMagician) Modulation(param int) ([]byte, error) {
	if param > 2 || param < 0 {
		return nil, fmt.Errorf("Modulation: %v is out of paramter options (0,1,2)")
	}
	resp, err := ir.writeread("m,"+strconv.Itoa(param)+"\r\n", 0)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (ir *IrMagician) SetRecordPointer(point int) ([]byte, error) {
	// TODO: Implement me.
	return nil, nil
}

func (ir *IrMagician) Play() ([]byte, error) {
	resp, err := ir.writeread("P\r\n", 0)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (ir *IrMagician) Reset(n int) ([]byte, error) {
	if n > 1 || n < 0 {
		return nil, fmt.Errorf("Reset: %v is not in options (0, 1)", n)
	}
	resp, err := ir.writeread("R,"+strconv.Itoa(n)+"\r\n", 1)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (ir *IrMagician) Version() ([]byte, error) {
	resp, err := ir.writeread("V\r\n", 1)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (ir *IrMagician) Write(pos int, data byte) error {
	if pos > 63 || pos < 0 {
		return fmt.Errorf("Write: %v is out of memory position range (0-63)", pos)
	}
	_, err := ir.s.Write([]byte("W," + strconv.Itoa(pos) + "," + string(data)))
	if err != nil {
		return err
	}
	return nil
}
