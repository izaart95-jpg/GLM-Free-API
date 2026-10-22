// cmd/token-collector/feilin.go
//
// FeiLin / Aliyun-SAF device-token primitives, ported to Go.
//
// The device token that the Z.AI bridge consumes is NOT minted in the browser
// any more: the browser is used exactly once, briefly, to obtain a *trusted*
// device session plus the genuine collector payload template; every token is
// then produced in-process by assembling the payload and AES-sealing it.
//
//	token = base64( tF # Q # w # tC # p )
//	  tF = region prefix        ("SG_WEB")
//	  Q  = device sessionId     ("<appKey>-h-<serverTs-1>-<32 hex>")
//	  w  = base64( AES-128-CBC(sessionKey, IV, PKCS7("#".join(payload))) )
//	  tC = GatherCost (ms)
//	  p  = MD5( tF # Q # w # tC # "daye,raolewoba!" )
//
// Everything below is derived from live captures; see
// .assets/reports/aliyun-captcha-report.md and
// .assets/resources/live_captures.md.

package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ---------- protocol constants ----------

const (
	// aesIV is the ONE IV used everywhere: the UTF-8 bytes of the literal
	// "0123456789ABCDEF" (NOT hex-decoded).
	aesIV = "0123456789ABCDEF"

	// resKey decrypts the DeviceConfig blob returned by the device service.
	resKey = "87f879f135f27da7"

	// md5Secret seals the token (field p).
	md5Secret = "daye,raolewoba!"

	// regionPrefix is the tF field for the ap-southeast endpoints.
	regionPrefix = "SG_WEB"

	// sceneID is the Z.AI captcha scene.
	sceneID = "didk33e0"
)

// ---------- AES-128-CBC (PKCS7) ----------

func aesCBCEncrypt(key, plaintext []byte) ([]byte, error) {
	if len(key) > aes.BlockSize {
		key = key[:aes.BlockSize]
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	padLen := aes.BlockSize - len(plaintext)%aes.BlockSize
	padded := make([]byte, 0, len(plaintext)+padLen)
	padded = append(padded, plaintext...)
	padded = append(padded, bytes.Repeat([]byte{byte(padLen)}, padLen)...)
	out := make([]byte, len(padded))
	cipher.NewCBCEncrypter(block, []byte(aesIV)).CryptBlocks(out, padded)
	return out, nil
}

func aesCBCDecrypt(key, ciphertext []byte) ([]byte, error) {
	if len(key) > aes.BlockSize {
		key = key[:aes.BlockSize]
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) == 0 || len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a whole number of AES blocks")
	}
	out := make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(block, []byte(aesIV)).CryptBlocks(out, ciphertext)
	// strip PKCS7
	padLen := int(out[len(out)-1])
	if padLen < 1 || padLen > aes.BlockSize || padLen > len(out) {
		return nil, errors.New("bad PKCS7 padding")
	}
	for _, b := range out[len(out)-padLen:] {
		if int(b) != padLen {
			return nil, errors.New("bad PKCS7 padding")
		}
	}
	return out[:len(out)-padLen], nil
}

// encryptB64 encrypts with AES-128-CBC and returns base64.
func encryptB64(key, plaintext string) (string, error) {
	ct, err := aesCBCEncrypt([]byte(key), []byte(plaintext))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ct), nil
}

// decryptB64 reverses encryptB64.
func decryptB64(key, b64 string) (string, error) {
	ct, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", err
	}
	pt, err := aesCBCDecrypt([]byte(key), ct)
	if err != nil {
		return "", err
	}
	return string(pt), nil
}

// ---------- the device session ----------

// deviceSession is the "trusted" device registration handed out by the device
// service (cloudauth-device-dualstack) in its Log1 response.
type deviceSession struct {
	SessionKey string // AES key that seals the token payload
	SessionID  string // Q field of the token
	ServerTs   int64  // server timestamp (ms)
	InitTime   int64  // session start (ms), serverTs-6606 in live traffic
	IP         string // the server's view of our IP
	Region     string // tF field ("SG_WEB")
	GatherCost string // tC field (per-session constant)
}

// decryptDeviceConfig parses the DeviceConfig blob returned by the device
// service.  The blob is AES-128-CBC(RES key, IV) over a '#'-joined record:
//
//	f[0] base64(sessionKey)  f[1] "1"          f[2] sessionId
//	f[3] feilin version       f[7] serverTs     f[8] ip
func decryptDeviceConfig(blob string) (deviceSession, error) {
	var s deviceSession
	plain, err := decryptB64(resKey, blob)
	if err != nil {
		return s, fmt.Errorf("decrypt DeviceConfig: %w", err)
	}
	f := strings.Split(plain, "#")
	if len(f) < 9 {
		return s, fmt.Errorf("DeviceConfig has %d fields, expected >= 9", len(f))
	}
	sk, err := base64.StdEncoding.DecodeString(f[0])
	if err != nil {
		return s, fmt.Errorf("DeviceConfig sessionKey: %w", err)
	}
	s.SessionKey = string(sk)
	s.SessionID = f[2]
	s.ServerTs, _ = strconv.ParseInt(f[7], 10, 64)
	s.IP = f[8]
	s.InitTime = s.ServerTs - 6606 // live gap; the SDK derives it the same way
	s.Region = regionPrefix
	s.GatherCost = "0"
	return s, nil
}

// timeNowMillis is the wall clock in milliseconds (the token's per-token
// timestamp source).
func timeNowMillis() int64 { return time.Now().UnixMilli() }

// ---------- the payload template ----------

// payloadTemplate is the genuine collector payload captured from the browser
// (the '#'-joined field list).  It is device- and session-bound; only the
// per-session timestamps and the trailing probe-log entries are refreshed per
// token.
type payloadTemplate struct {
	fields []string
}

// field indices of the per-session / per-token values inside the payload.
//
// The index map is stable across the SDK builds we have captured (140, 111 and
// 155 field variants): 42=ip, 43=probe log, 72=initTime, 74=now, 87=serverTs.
// Index 21 is NOT portable (older builds carried a per-session opaque probe
// there, newer ones a capability vector), so it is copied verbatim - the risk
// layer does not check it.
const (
	idxIP        = 42
	idxLog       = 43
	idxInitTime  = 72
	idxNow       = 74
	idxServerTs  = 87
	logEntry93   = "93"
	logEntry94   = "94"
	idxEmpty93   = 93
	idxEmpty94   = 94
	maxLogJitter = 15
)

// newPayloadTemplate stores the field list captured from the browser.
func newPayloadTemplate(fields []string) (*payloadTemplate, error) {
	if len(fields) < 90 {
		return nil, fmt.Errorf("payload template has %d fields, expected >= 90", len(fields))
	}
	return &payloadTemplate{fields: fields}, nil
}

// clone returns a deep copy of the template's field list.
func (t *payloadTemplate) clone() []string {
	out := make([]string, len(t.fields))
	copy(out, t.fields)
	return out
}

// build refreshes the per-session and per-token values for one token.
//
//	nowMs      wall clock of the mint
//	jitterMs   offset of probe-log entry 94 over entry 93 (live: 9-15 ms)
func (t *payloadTemplate) build(s deviceSession, nowMs int64, jitterMs int) ([]string, error) {
	f := t.clone()
	if len(f) <= idxServerTs {
		return nil, fmt.Errorf("payload template too short (%d fields)", len(f))
	}

	// --- per-session ---
	f[idxIP] = s.IP
	f[idxServerTs] = strconv.FormatInt(s.ServerTs, 10)
	f[idxInitTime] = strconv.FormatInt(s.InitTime, 10)

	// --- per-token ---
	f[idxNow] = strconv.FormatInt(nowMs, 10)

	// --- probe log: the trailing 93/94 entries are defined by (now - initTime) ---
	entries := strings.Split(f[idxLog], "|")
	if len(entries) < 2 {
		return nil, fmt.Errorf("payload field 43 has no log entries")
	}
	last, prev := entries[len(entries)-1], entries[len(entries)-2]
	if !strings.HasPrefix(prev, logEntry93+"-") || !strings.HasPrefix(last, logEntry94+"-") {
		return nil, fmt.Errorf("payload field 43 must end with log entries 93/94, got %q/%q", prev, last)
	}
	if jitterMs <= 0 {
		// keep the template's observed delta (entry94 cost - entry93 cost)
		d94, err1 := strconv.Atoi(strings.SplitN(last, "-", 2)[1])
		d93, err2 := strconv.Atoi(strings.SplitN(prev, "-", 2)[1])
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("payload field 43: bad 93/94 entries: %q/%q", prev, last)
		}
		jitterMs = d94 - d93
	}
	if jitterMs < 1 || jitterMs > maxLogJitter {
		jitterMs = 9
	}
	cost := nowMs - s.InitTime
	if cost < 0 {
		return nil, fmt.Errorf("mint time %d is before initTime %d", nowMs, s.InitTime)
	}
	entries[len(entries)-2] = logEntry93 + "-" + strconv.FormatInt(cost, 10)
	entries[len(entries)-1] = logEntry94 + "-" + strconv.FormatInt(cost+int64(jitterMs), 10)
	f[idxLog] = strings.Join(entries, "|")

	// --- the payload's own 93/94 slots stay empty (the log entries live in 43) ---
	if len(f) > idxEmpty94 {
		f[idxEmpty93] = ""
		f[idxEmpty94] = ""
	}
	return f, nil
}

// ---------- token assembly ----------

// assembleToken builds the final base64 device token.
func assembleToken(s deviceSession, fields []string) (string, error) {
	w, err := encryptB64(s.SessionKey, strings.Join(fields, "#"))
	if err != nil {
		return "", err
	}
	tC := s.GatherCost
	seal := md5.Sum([]byte(strings.Join([]string{s.Region, s.SessionID, w, tC, md5Secret}, "#")))
	p := hex.EncodeToString(seal[:])
	return base64.StdEncoding.EncodeToString([]byte(strings.Join([]string{s.Region, s.SessionID, w, tC, p}, "#"))), nil
}

// mintTokens produces n distinct device tokens bound to one device session.
//
// The tokens differ only in the per-token timestamp (field 74) and the two
// trailing probe-log entries, which is exactly what the SDK does when it is
// asked for a second token for the same session.
func mintTokens(s deviceSession, tpl *payloadTemplate, n int, jitterMs int) ([]string, error) {
	if n <= 0 {
		return nil, errors.New("token count must be > 0")
	}
	out := make([]string, 0, n)
	// The probe log is monotonic: the first token must not be minted before the
	// collector's last fixed probe cost has elapsed.
	base := timeNowMillis()
	fixed := maxFixedLogCost(tpl.fields)
	start := base
	if s.InitTime > 0 {
		if minStart := s.InitTime + fixed; start < minStart {
			start = minStart
		}
	}
	for i := 0; i < n; i++ {
		now := start + int64(i)
		fields, err := tpl.build(s, now, jitterMs)
		if err != nil {
			return nil, err
		}
		tok, err := assembleToken(s, fields)
		if err != nil {
			return nil, err
		}
		out = append(out, tok)
	}
	return out, nil
}

// maxFixedLogCost returns the largest probe-log cost of the template's fixed
// entries (everything before the trailing 93/94 pair).
func maxFixedLogCost(fields []string) int64 {
	if len(fields) <= idxLog {
		return 0
	}
	entries := strings.Split(fields[idxLog], "|")
	if len(entries) < 2 {
		return 0
	}
	fixed := entries[:len(entries)-2]
	var max int64
	for _, e := range fixed {
		parts := strings.SplitN(e, "-", 2)
		if len(parts) != 2 {
			continue
		}
		if c, err := strconv.ParseInt(parts[1], 10, 64); err == nil && c > max {
			max = c
		}
	}
	return max
}

// ---------- verification (self-test helper) ----------

// verifyToken decodes a device token and checks the MD5 seal; it returns the
// payload fields when the seal is good.
func verifyToken(token, sessionKey string) ([]string, bool, error) {
	raw, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, false, fmt.Errorf("token is not base64: %w", err)
	}
	parts := strings.Split(string(raw), "#")
	if len(parts) != 5 {
		return nil, false, fmt.Errorf("token has %d fields, expected 5", len(parts))
	}
	tF, Q, w, tC, p := parts[0], parts[1], parts[2], parts[3], parts[4]
	seal := md5.Sum([]byte(strings.Join([]string{tF, Q, w, tC, md5Secret}, "#")))
	if hex.EncodeToString(seal[:]) != p {
		return nil, false, errors.New("md5 seal mismatch")
	}
	plain, err := decryptB64(sessionKey, w)
	if err != nil {
		return nil, false, fmt.Errorf("decrypt payload: %w", err)
	}
	return strings.Split(plain, "#"), true, nil
}
