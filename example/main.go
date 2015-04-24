package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	_ "time"

	"bitbucket.org/ymotongpoo/irmagician"
)

func main() {
	ir, err := irmagician.NewIrMagician("/dev/ttyACM0", 9600, 0)
	if err != nil {
		log.Fatal(err)
	}

	_, err = ir.Capture()
	if err != nil {
		log.Fatal(err)
	}

	resps, err := ir.Information(1)
	if err != nil {
		log.Fatal(err)
	}
	resp := bytes.Split(resps, []byte("\r\n"))[0]
	rec, err := irmagician.ParseRawInt(resp, 16)
	if err != nil {
		log.Fatal("rec num", err)
	}
	resps, err = ir.Information(6)
	if err != nil {
		log.Fatal(err)
	}
	resp = bytes.Split(resps, []byte("\r\n"))[0]
	scale, err := irmagician.ParseRawInt(resp, 10)
	if err != nil {
		log.Fatal("scale num", err)
	}
	log.Printf("rec: %v", rec)

	raw := make([]string, rec)
	for i := 0; i < rec; i++ {
		bank := i / 64
		pos := i % 64
		if pos == 0 {
			ir.BankSet(bank)
		}
		resp, err := ir.Dump(pos)
		if err != nil {
			log.Fatal(err)
		}
		/*
			data, err := irmagician.ParseRawInt(resp)
			if err != nil {
				log.Fatalf("dumping: %v (bank: %v, pos: %v)", err, bank, pos)
			}
		*/
		raw[i] = string(bytes.TrimSpace(resp))
	}
	jsonData := fmt.Sprintf(`{'format': 'raw', 'freq': 38, 'data': %v, 'postscale': scale}\n`, strings.Join(raw, ","), scale)
	err = ioutil.WriteFile("dump.txt", []byte(jsonData), 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("dumped to dump.txt")
}
