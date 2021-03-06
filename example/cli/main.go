package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"time"

	"github.com/ymotongpoo/irmagician"
)

var (
	play    = flag.Bool("p", false, "play stored data")
	capture = flag.Bool("c", false, "capture data")
	save    = flag.Bool("s", false, "save stored data")
	path    = flag.String("f", "", "path to data file")
)

func main() {
	flag.Parse()

	ir, err := irmagician.NewIrMagician("/dev/ttyACM0", 9600, irmagician.DefaultTimeout)
	if err != nil {
		log.Fatal(err)
	}
	defer ir.Close()

	switch {
	case *capture:
		err = captureData(ir)
	case *play:
		err = playData(ir, *path)
	case *save:
		err = saveData(ir, *path)
	default:
		log.Println("confirm options with --help.")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func captureData(ir *irmagician.IrMagician) error {
	_, err := ir.Capture()
	return err
}

func playData(ir *irmagician.IrMagician, path string) error {
	log.Println("playData")
	if path == "" {
		out, err := ir.Play()
		time.Sleep(1 * time.Second)
		log.Println(string(out))
		return err
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	out, err := ir.PlayData(data)
	if err != nil {
		return err
	}
	log.Println(string(out))
	return nil
}

func saveData(ir *irmagician.IrMagician, path string) error {
	if path == "" {
		path = time.Now().Format("20060102150405") + ".json"
	}

	resps, err := ir.Information(1)
	if err != nil {
		return err
	}
	resp := bytes.Split(resps, []byte("\r\n"))[0]
	rec, err := irmagician.ParseRawInt(resp, 16)
	if err != nil {
		return err
	}
	resps, err = ir.Information(6)
	if err != nil {
		return err
	}
	resp = bytes.Split(resps, []byte("\r\n"))[0]
	scale, err := irmagician.ParseRawInt(resp, 10)
	if err != nil {
		return err
	}
	log.Printf("rec: %v", rec)

	raw := make([]byte, rec)
	for i := 0; i < rec; i++ {
		bank := i / 64
		pos := i % 64
		if pos == 0 {
			ir.BankSet(bank)
		}
		resp, err := ir.Dump(pos)
		if err != nil {
			return err
		}
		b, err := irmagician.ParseRawInt(resp, 16)
		if err != nil {
			return err
		}
		raw[i] = byte(b)
	}
	dump := irmagician.Dump{
		Scale:  scale,
		Format: "raw",
		Freq:   38,
		Data:   raw,
	}
	data, err := json.Marshal(dump)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}
	log.Printf("dumped to %v", path)
	return nil
}
