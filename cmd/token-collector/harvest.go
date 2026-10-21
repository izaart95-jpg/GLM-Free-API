// cmd/token-collector/harvest.go
//
// ONE brief browser use, then programmatic minting.
//
// The Z.AI risk layer trusts a device session only when the `Log1` request
// that created it came from a client it recognises as a real browser, and only
// once the device has uploaded its fingerprint (`Log2`/`Log3`) for that very
// session.  A pure HTTP client can therefore never mint a usable device token:
// the request bytes are not the problem, the *client* is.
//
// So the browser is used exactly once, briefly, to
//
//  1. boot the Aliyun-SAF captcha SDK on chat.z.ai (the SDK performs Log1 ->
//     InitCaptchaV3 -> Log2 -> Log3, which registers the device), and
//  2. hand us (a) the DeviceConfig that carries the trusted session and
//     (b) one genuine device token whose payload is the collector template.
//
// After that the browser is closed and every token is produced in-process by
// feilin.go.  The browser is never used again for token minting.

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mxschmitt/playwright-go"
)

// captchaSDKURL is the Aliyun captcha front-end bundle.
const captchaSDKURL = "https://o.alicdn.com/captcha-frontend/aliyunCaptcha/AliyunCaptcha.js"

// deviceServiceURL is the FeiLin device service that issues the session.
const deviceServiceURL = "https://cloudauth-device-dualstack.ap-southeast-1.aliyuncs.com"

// bootCaptchaJS boots the captcha SDK inside the page and returns one genuine
// device token.  It mirrors what the SDK does on a real Z.AI page: the SDK
// registers the device (Log1), initialises the challenge (InitCaptchaV3) and
// uploads the collector fingerprint (Log2/Log3) - that registration is what
// makes the session trusted.
//
// The mount element must be visible (an off-screen element never gets an
// instance) and the SDK has to be reachable, so the script injects it itself.
const bootCaptchaJS = `(async () => {
    const sleep = (ms) => new Promise(r => setTimeout(r, ms));
    try {
        window.AliyunCaptchaConfig = { region: 'sgp', prefix: 'no8xfe' };
        if (!window.initAliyunCaptcha) {
            await new Promise((res, rej) => {
                const s = document.createElement('script');
                s.src = '` + captchaSDKURL + `';
                s.onload = res;
                s.onerror = () => rej(new Error('captcha sdk load failed'));
                document.head.appendChild(s);
            });
        }
        const el = document.createElement('div');
        el.id = '__zai_cap_el';
        el.style.cssText = 'position:fixed;left:10px;top:200px;width:320px;height:200px;z-index:2147483647;';
        document.body.appendChild(el);
        const btn = document.createElement('button');
        btn.id = '__zai_cap_btn';
        btn.type = 'button';
        btn.textContent = 'go';
        btn.style.cssText = 'position:fixed;left:10px;top:420px;width:120px;height:40px;z-index:2147483647;';
        document.body.appendChild(btn);

        let instance = null;
        window.initAliyunCaptcha({
            SceneId: '` + sceneID + `',
            mode: 'embed',
            element: '#__zai_cap_el',
            button: '#__zai_cap_btn',
            captchaLogoImg: 'https://z-cdn.chatglm.cn/z-ai/static/logo.svg',
            upLang: true,
            language: 'en',
            timeout: 10000,
            delayBeforeSuccess: false,
            success: () => {}, fail: () => {}, onError: () => {}, onClose: () => {},
            getInstance: (i) => { instance = i; window.__zai_cap_inst = i; },
        });

        const t0 = Date.now();
        while (Date.now() - t0 < 60000) {
            if (instance && instance.captcha && window.z_um &&
                typeof window.z_um.getToken === 'function') break;
            await sleep(200);
        }
        if (!window.z_um || typeof window.z_um.getToken !== 'function') {
            return { ok: false, error: 'captcha SDK did not expose z_um.getToken' };
        }
        // let the collector finish its probes so the payload template is complete
        await sleep(3000);
        const tok = window.z_um.getToken();
        return { ok: true, token: (tok && typeof tok.then === 'function') ? await tok : tok };
    } catch (e) {
        return { ok: false, error: String((e && e.stack) || e) };
    }
})()`

// harvestResult carries everything one browser visit yields.
type harvestResult struct {
	Session  deviceSession
	Template *payloadTemplate
	Token    string
}

// harvestSession performs the one-off browser visit.
//
// It opens chat.z.ai (so the page origin, referrer and TLS session look like a
// real visit), boots the captcha SDK, waits for the device registration to
// complete and extracts
//
//   - the DeviceConfig from the device service's Log1 response, and
//   - one genuine device token, whose payload is the collector template.
//
// The caller closes the browser afterwards; nothing else is needed from it.
func harvestSession(browser playwright.Browser) (*harvestResult, error) {
	ctx, page, err := newWorkerPage(browser)
	if err != nil {
		return nil, fmt.Errorf("create page: %w", err)
	}
	defer ctx.Close()

	// Capture the device service's Log1 response - it carries the DeviceConfig
	// with the trusted session.
	var log1Mu chan struct{}
	log1Ch := make(chan string, 4)
	log1Mu = nil
	_ = log1Mu
	page.On("response", func(r playwright.Response) {
		if !strings.Contains(r.URL(), "cloudauth-device-dualstack") {
			return
		}
		body, err := r.Body()
		if err != nil {
			return
		}
		txt := string(body)
		if strings.Contains(txt, "DeviceConfig") {
			select {
			case log1Ch <- txt:
			default:
			}
		}
	})

	if _, err := page.Goto(URL, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
		Timeout:   playwright.Float(60000),
	}); err != nil {
		return nil, fmt.Errorf("goto: %w", err)
	}

	tuiSetStatus("Booting captcha SDK (one-off device registration)...")
	res, err := page.Evaluate(bootCaptchaJS)
	if err != nil {
		return nil, fmt.Errorf("boot captcha: %w", err)
	}
	obj, _ := res.(map[string]interface{})
	if obj == nil {
		return nil, fmt.Errorf("boot captcha: unexpected result %T", res)
	}
	if ok, _ := obj["ok"].(bool); !ok {
		msg, _ := obj["error"].(string)
		return nil, fmt.Errorf("boot captcha: %s", msg)
	}
	token, _ := obj["token"].(string)
	if token == "" {
		return nil, fmt.Errorf("boot captcha: SDK returned an empty device token")
	}

	// The Log1 response may have arrived before or after the SDK call returns.
	var log1 string
	select {
	case log1 = <-log1Ch:
	case <-time.After(15 * time.Second):
	}
	if log1 == "" {
		return nil, fmt.Errorf("device service Log1 response not captured")
	}

	var lr struct {
		ResultObject struct {
			DeviceConfig string `json:"DeviceConfig"`
		} `json:"ResultObject"`
	}
	if err := json.Unmarshal([]byte(log1), &lr); err != nil {
		return nil, fmt.Errorf("parse Log1 response: %w", err)
	}
	if lr.ResultObject.DeviceConfig == "" {
		return nil, fmt.Errorf("Log1 response carried no DeviceConfig")
	}
	session, err := decryptDeviceConfig(lr.ResultObject.DeviceConfig)
	if err != nil {
		return nil, err
	}

	// The genuine token's payload is the collector template for this device.
	fields, ok, err := verifyToken(token, session.SessionKey)
	if err != nil {
		return nil, fmt.Errorf("decrypt captured device token: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("captured device token failed its md5 seal")
	}
	tpl, err := newPayloadTemplate(fields)
	if err != nil {
		return nil, err
	}

	// The GatherCost (tC) is a per-SESSION constant (HC3: 500-5000 ms).  Take it
	// from the token the SDK itself produced rather than inventing one.
	if raw, derr := base64.StdEncoding.DecodeString(token); derr == nil {
		if parts := strings.Split(string(raw), "#"); len(parts) == 5 {
			if gc, cerr := strconv.Atoi(parts[3]); cerr == nil && gc > 0 {
				session.GatherCost = parts[3]
			}
		}
	}

	return &harvestResult{Session: session, Template: tpl, Token: token}, nil
}

// ---------- session minting manager ----------

// sessionMinter turns ONE brief browser visit into many device tokens.
//
// The browser is opened only when there is no usable session yet (or the
// current one has aged past maxAge - the device session has a server side
// lifetime, measured at >= 35 min in live traffic).  Between harvests the
// tokens are minted in-process.
type sessionMinter struct {
	browser      playwright.Browser
	session      deviceSession
	template     *payloadTemplate
	harvestedAt  time.Time
	harvestCount int
	minted       int
}

// ensure makes sure a usable session is available, harvesting when needed.
func (m *sessionMinter) ensure(maxAge time.Duration) error {
	if m.session.SessionID != "" && time.Since(m.harvestedAt) < maxAge {
		return nil
	}
	tuiSetStatus("Harvesting trusted device session (one-off browser visit)...")
	fmt.Println("🌐 One-off browser visit: registering this device with the FeiLin device service...")
	h, err := harvestSession(m.browser)
	if err != nil {
		return fmt.Errorf("harvest: %w", err)
	}
	m.session = h.Session
	m.template = h.Template
	m.harvestedAt = time.Now()
	m.harvestCount++
	fmt.Printf("✅ Trusted device session %s (key %s) from %d field payload template\n",
		m.session.SessionID, m.session.SessionKey, len(h.Template.fields))
	return nil
}

// mint produces n fresh device tokens bound to the current session.
func (m *sessionMinter) mint(n int) ([]string, error) {
	if m.template == nil {
		return nil, fmt.Errorf("no payload template - harvest first")
	}
	toks, err := mintTokens(m.session, m.template, n, 0)
	if err != nil {
		return nil, err
	}
	m.minted += len(toks)
	return toks, nil
}

// stats reports how much browser time the run actually cost.
func (m *sessionMinter) stats() string {
	return fmt.Sprintf("%d browser visit(s), %d tokens minted programmatically",
		m.harvestCount, m.minted)
}

// exportSession writes the harvested session + the genuine payload template as
// JSON so the Python reference tooling (.assets/scripts/live_validate.py
// --session-mint) can re-use the very same session.
func (m *sessionMinter) exportSession(path string) error {
	if m.template == nil {
		return fmt.Errorf("no session to export")
	}
	blob := map[string]any{
		"session": map[string]any{
			"sessionKey": m.session.SessionKey,
			"sessionId":  m.session.SessionID,
			"serverTs":   m.session.ServerTs,
			"initTime":   m.session.InitTime,
			"ip":         m.session.IP,
			"tC":         m.session.GatherCost,
			"region":     m.session.Region,
		},
		"template": m.template.fields,
	}
	data, err := json.MarshalIndent(blob, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
