package main

import (
	log "github.com/Sirupsen/logrus"
	"net/http"
	"io"
	"io/ioutil"
	"strings"
	"math/rand"
	"time"
	"fmt"
	"os"
)

var (
	HitTotal = 0
	UseType = HitTypeB
	UseMethod = "POST"
	UseURL = "http://1cssvepl0m60.runscope.net/shittp"
	UsePayloadLength = 1024
)

func main() {
	setup()
	done := make(chan Hit)
	count := 0
	for i := 0; i < HitTotal; i++ {
		h := NewHit(UseType, i, UseMethod, UseURL, randStringBytesMaskImpr(UsePayloadLength))
		log.WithFields(log.Fields{
			"type": "a",
			"num": h.num,
		}).Info("Begin HTTP")
		go h.Do(done)
	}
	for count < HitTotal {
		h := <-done
		if h.err != nil {
			log.WithFields(log.Fields{
				"type": "a",
				"num": h.num,
				"error": h.err,
				"duration": time.Since(h.startAt),
			}).Warn("Done HTTP (Failed)")
		} else {
			log.WithFields(log.Fields{
				"type": "a",
				"num": h.num,
				"duration": time.Since(h.startAt),
			}).Info("Done HTTP")
		}
		count++
	}
	teardown()
}

func setup() {
	var err error
	fmt.Printf("\nWelcome!\npid: %v\n\nEnter total # of HTTP to hit, then <ENTER>\n", os.Getpid())
	_, err = fmt.Scanf("%d\n", &HitTotal)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nEnter type of hit (A or B) then <ENTER>\n")
	_, err = fmt.Scanf("%s\n", &UseType)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Okay, hitting %d times, type %s..\n\n", HitTotal, UseType)
}

func teardown() {
	log.WithFields(log.Fields{
		"total": HitTotal,
	}).Info("Done")
	var exit string
	fmt.Printf("\nDone!\nWhen ready to exit pid %v, press <ENTER>\n", os.Getpid())
	fmt.Scanf("%s\n", &exit)
}

func NewHit(t HitType, num int, meth string, url string, payload string) Hit {
	return Hit {
		t:t,
		num:num,
		meth:meth,
		url:url,
		payload:payload,
		startAt: time.Now(),
		err:nil,
	}
}

type Hit struct {
	t       HitType
	num     int
	meth    string
	url     string
	payload string
	startAt time.Time
	err     error
}

func (h *Hit) Do(done chan<- Hit) {
	if h.t == HitTypeA {
		client := &http.Client{}
		req, _ := http.NewRequest(h.meth, h.url, strings.NewReader(h.payload))
		req.Header.Add("Connection", "close")
		req.Header.Add("Content-Type", "application/json")
		_, err := client.Do(req)
		if err != nil {
			h.err = err
		}
	} else if h.t == HitTypeB {
		client := &http.Client{}
		req, _ := http.NewRequest(h.meth, h.url, strings.NewReader(h.payload))
		req.Header.Add("Connection", "close")
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			h.err = err
		} else {
			io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		}
	}
	done <- *h
}

type HitType string
var (
	HitTypeA = HitType("A")
	HitTypeB = HitType("B")
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)
func randStringBytesMaskImpr(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

