package irmagician

import (
	"encoding/json"
	"fmt"
	_ "log"
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

func (ir *IrMagician) writeread(command string, waitmsec int) ([]byte, error) {
	//log.Println(command)
	_, err := ir.s.Write([]byte(command))
	if err != nil {
		return nil, err
	}
	if waitmsec != 0 {
		time.Sleep(time.Duration(waitmsec) * time.Millisecond)
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
	resp, err := ir.writeread("I,"+strconv.Itoa(param)+"\r\n", 5)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (ir *IrMagician) SetPostScaler(v int) ([]byte, error) {
	if v > 255 || v < 1 {
		return nil, fmt.Errorf("SetPostScaler: %v is out of value range (1-255)", v)
	}
	resp, err := ir.writeread("k,"+strconv.Itoa(v)+"\r\n", 5)
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
	resp, err := ir.writeread("L,"+toggle+"\r\n", 5)
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
	if point < 0 {
		return nil, fmt.Errorf("SetRecordPointer: point must be unsigned.")
	}
	resp, err := ir.writeread("n,"+strconv.Itoa(point)+"\r\n", 2)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (ir *IrMagician) Play() ([]byte, error) {
	resp, err := ir.writeread("p\r\n", 2)
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
	buf := fmt.Sprintf("w,%d,%d\r\n", pos, data)
	_, err := ir.s.Write([]byte(buf))
	if err != nil {
		return err
	}
	return nil
}

func (ir *IrMagician) Close() {
	ir.s.Close()
}

func (ir *IrMagician) PlayData(data []byte) ([]byte, error) {
	var dump Dump
	if err := json.Unmarshal(data, &dump); err != nil {
		return nil, err
	}
	_, err := ir.SetRecordPointer(len(dump.Data))
	if err != nil {
		return nil, err
	}
	_, err = ir.SetPostScaler(dump.Scale)
	if err != nil {
		return nil, err
	}
	for i, b := range dump.Data {
		bank := i / 64
		pos := i % 64
		if pos == 0 {
			err = ir.BankSet(bank)
			if err != nil {
				return nil, err
			}
		}
		err = ir.Write(pos, b)
		if err != nil {
			return nil, err
		}
	}
	resp, err := ir.Play()
	if err != nil {
		return nil, err
	}
	return resp, nil
}
