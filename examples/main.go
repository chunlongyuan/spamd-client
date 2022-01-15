// Copyright (C) 2018-2021 Andrew Colin Kissa <andrew@datopdog.io>
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.

/*
Package main
spamd-client - Golang Spamd SpamAssassin Client
*/
package main

// StatusCode StatusCode
// StatusMsg  string
// Version    string
// Score      float64
// BaseScore  float64
// IsSpam     bool
// Headers    textproto.MIMEHeader
// Msg        *Msg
// Rules      []map[string]string

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	spamdclient "github.com/baruwa-enterprise/spamd-client/pkg"
	"github.com/baruwa-enterprise/spamd-client/pkg/response"
	flag "github.com/spf13/pflag"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

var (
	cfg     *Config
	cmdName string
)

// Config holds the configuration
type Config struct {
	Address        string
	Port           int
	UseTLS         bool
	User           string
	UseCompression bool
	RootCA         string
}

func d(r *response.Response) {
	// log.Println("===================================")
	log.Printf("RequestMethod => %v\n", r.RequestMethod)
	log.Printf("StatusCode => %v\n", r.StatusCode)
	log.Printf("StatusMsg => %v\n", r.StatusMsg)
	log.Printf("Version => %v\n", r.Version)
	log.Printf("Score => %v\n", r.Score)
	log.Printf("BaseScore => %v\n", r.BaseScore)
	log.Printf("IsSpam => %v\n", r.IsSpam)
	log.Printf("Headers => %v\n", r.Headers)
	log.Printf("Msg => %v\n", r.Msg)
	log.Printf("Msg.Header => %v\n", r.Msg.Header)
	log.Printf("Msg.Body => %s", r.Msg.Body)
	log.Printf("Msg.Raw => %s", r.Raw)
	log.Printf("Rules => %v\n", r.Rules)
	log.Println("===================================")
}

func init() {
	cfg = &Config{}
	cmdName = path.Base(os.Args[0])
	flag.StringVarP(&cfg.Address, "host", "H", "192.168.15.185",
		`Specify Spamd host to connect to.`)
	flag.IntVarP(&cfg.Port, "port", "p", 783,
		`In TCP/IP mode, connect to spamd server listening on given port`)
	flag.BoolVarP(&cfg.UseTLS, "use-tls", "S", false,
		`Use TLS.`)
	flag.StringVarP(&cfg.User, "user", "u", "exim",
		`User for spamd to process this message under.`)
	flag.BoolVarP(&cfg.UseCompression, "use-compression", "z", false,
		`Compress mail message sent to spamd.`)
	flag.StringVarP(&cfg.RootCA, "root-ca", "r", "/Users/andrew/tmp/frontend-ca.pem",
		`The CA certificate for verifying the TLS connection.`)
}

func parseAddr(a string, p int) (n string, h string) {
	if strings.HasPrefix(a, "/") {
		n = "unix"
		h = a
	} else {
		n = "tcp"
		if strings.Contains(a, ":") {
			h = fmt.Sprintf("[%s]:%d", a, p)
		} else {
			h = fmt.Sprintf("%s:%d", a, p)
		}
	}
	return
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", cmdName)
	fmt.Fprint(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.ErrHelp = errors.New("")
	flag.CommandLine.SortFlags = false
	flag.Parse()
	ctx := context.Background()
	network, address := parseAddr(cfg.Address, cfg.Port)
	m := []byte(`Date: Mon, 23 Jun 2021 11:40:36 -0400
From: Gopher <test@gmail.com>
To: Another Gopher <spamtest@gmail.com>
Subject: Gophers test spam at Gophercon
Message-Id: <v0421010eb70653b14e06@[192.168.15.185]>

Message body
James

My Workd

++++++++++++++
`)

	var wg sync.WaitGroup

	wg.Add(1)
	go func(m []byte) {
		defer wg.Done()
		c, err := spamdclient.NewClient(network, address, cfg.User, cfg.UseCompression)
		if err != nil {
			log.Println(err)
			return
		}
		c.SetCmdTimeout(30 * time.Second)
		if cfg.UseTLS {
			err = c.SetRootCA(cfg.RootCA)
			if err != nil {
				log.Println("ERROR:", err)
				return
			}
			c.EnableTLS()
		}
		ir := bytes.NewReader(m)
		r, e := c.Check(ctx, ir)
		if e != nil {
			log.Println(e)
			return
		}
		d(r)
	}(m)

	wg.Add(1)
	go func(m []byte) {
		defer wg.Done()
		c, err := spamdclient.NewClient(network, address, cfg.User, cfg.UseCompression)
		if err != nil {
			log.Println("ERROR:", err)
			return
		}

		if cfg.UseTLS {
			err = c.SetRootCA(cfg.RootCA)
			if err != nil {
				log.Println("ERROR:", err)
				return
			}
			c.EnableTLS()
		}
		c.EnableRawBody()
		ir := bytes.NewReader(m)
		r, e := c.Headers(context.Background(), ir)
		if e != nil {
			log.Println("ERROR:", e)
			return
		}
		d(r)
	}(m)

	wg.Add(1)
	go func(m []byte) {
		defer wg.Done()
		c, err := spamdclient.NewClient(network, address, cfg.User, cfg.UseCompression)
		if err != nil {
			log.Println(err)
			return
		}
		if cfg.UseTLS {
			err = c.SetRootCA(cfg.RootCA)
			if err != nil {
				log.Println("ERROR:", err)
				return
			}
			c.EnableTLS()
		}
		c.EnableRawBody()
		ir := bytes.NewReader(m)
		r, e := c.Process(ctx, ir)
		if e != nil {
			log.Println(e)
			return
		}
		d(r)
	}(m)

	wg.Add(1)
	go func(m []byte) {
		defer wg.Done()
		c, err := spamdclient.NewClient(network, address, cfg.User, cfg.UseCompression)
		if err != nil {
			log.Println(err)
			return
		}
		if cfg.UseTLS {
			err = c.SetRootCA(cfg.RootCA)
			if err != nil {
				log.Println("ERROR:", err)
				return
			}
			c.EnableTLS()
		}
		c.EnableRawBody()
		ir := bytes.NewReader(m)
		r, e := c.Report(ctx, ir)
		if e != nil {
			log.Println(e)
			return
		}
		d(r)
	}(m)

	wg.Add(1)
	go func(m []byte) {
		defer wg.Done()
		c, err := spamdclient.NewClient(network, address, cfg.User, cfg.UseCompression)
		if err != nil {
			log.Println(err)
			return
		}
		if cfg.UseTLS {
			err = c.SetRootCA(cfg.RootCA)
			if err != nil {
				log.Println("ERROR:", err)
				return
			}
			c.EnableTLS()
		}
		c.EnableRawBody()
		ir := bytes.NewReader(m)
		r, e := c.ReportIfSpam(ctx, ir)
		if e != nil {
			log.Println(e)
			return
		}
		d(r)
	}(m)

	wg.Add(1)
	go func(m []byte) {
		defer wg.Done()
		c, err := spamdclient.NewClient(network, address, cfg.User, cfg.UseCompression)
		if err != nil {
			log.Println(err)
			return
		}
		if cfg.UseTLS {
			err = c.SetRootCA(cfg.RootCA)
			if err != nil {
				log.Println("ERROR:", err)
				return
			}
			c.EnableTLS()
		}
		c.EnableRawBody()
		ir := bytes.NewReader(m)
		r, e := c.Symbols(ctx, ir)
		if e != nil {
			log.Println(e)
			return
		}
		d(r)
	}(m)
	wg.Wait()
	//c, err := spamdclient.NewClient(network, address, cfg.User, cfg.UseCompression)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//if cfg.UseTLS {
	//	err = c.SetRootCA(cfg.RootCA)
	//	if err != nil {
	//		log.Println("ERROR:", err)
	//		return
	//	}
	//	c.EnableTLS()
	//}
	// c.SetConnTimeout(2 * time.Second)
	// c.SetCmdTimeout(2 * time.Second)
	// c.SetConnRetries(5)
	//ir := bytes.NewReader(m)
	//r, e := c.Tell(ctx, ir, request.Ham, request.LearnAction)
	//if e != nil {
	//	log.Println(e)
	//	return
	//}
	//d(r)
	//ir.Reset(m)
	//r, e = c.Tell(ctx, ir, request.Ham, request.ForgetAction)
	//if e != nil {
	//	log.Println(e)
	//	return
	//}
	//d(r)
	//ir.Reset(m)
	//r, e = c.Tell(ctx, ir, request.Spam, request.LearnAction)
	//if e != nil {
	//	log.Println(e)
	//	return
	//}
	//d(r)
	//ir.Reset(m)
	//r, e = c.Tell(ctx, ir, request.Spam, request.ForgetAction)
	//if e != nil {
	//	log.Println(e)
	//	return
	//}
	//d(r)
}
