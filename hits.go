package main

import (
	log "github.com/Sirupsen/logrus"
	"net/http"
	"io"
	"io/ioutil"
	"strings"
	"math/rand"
	"time"
)

const (
	S_HITS = 100
	S_METH = "POST"
	S_HITTY_URL = "http://1cssvepl0m60.runscope.net/shittp"
	S_PAYLOAD_LENGTH = 1024
)

func main() {
	qHits := make(chan Hit)
	doneHits := make([]Hit, 0)
	for i := 0; i < S_HITS; i++ {
		s := Hit{HitTypeA, i, S_METH, S_HITTY_URL, randStringBytesMaskImpr(S_PAYLOAD_LENGTH), time.Now()}
		go s.Hit(qHits)
	}
	for len(doneHits) < S_HITS {
		h := <-qHits
		doneHits = append(doneHits, h)
	}
	log.WithFields(log.Fields{
		"total": len(doneHits),
	}).Info("Done")
}

type Hit struct {
	t testHitType
	num int
	meth string
	url string
	payload string
	startAt time.Time
}

func (h *Hit) Hit(qHits chan<- Hit) {
	switch h.t {
		case HitTypeA:
			log.WithFields(log.Fields{
				"type": "a",
				"num": h.num,
				}).Info("HTTP Request Begin")
			client := &http.Client{}
			req, _ := http.NewRequest(h.meth, h.url, strings.NewReader(h.payload))
			req.Header.Add("Connection", "close")
			req.Header.Add("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				log.WithFields(log.Fields{
					"type": "a",
					"num": h.num,
					"error": err,
					"duration": time.Since(h.startAt),
					}).Warn("HTTP Request Failed")
			} else {
				io.Copy(ioutil.Discard, resp.Body)
				resp.Body.Close()
				log.WithFields(log.Fields{
					"type": "a",
					"num": h.num,
					"duration": time.Since(h.startAt),
					}).Info("HTTP Request OK")
			}
		case testHitTypeB:
			// get f----d; there is no type B
	}
	qHits <- *h
}

type testHitType string
var (
	HitTypeA = testHitType("A")
	testHitTypeB = testHitType("B")
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

