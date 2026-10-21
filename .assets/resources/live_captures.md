# Live captures & session material — captain-verified

Shared reference for all team tasks. Everything here was provided by the user in the
team brief and verified by the captain in Python (AES-128-CBC, PKCS7,
IV = UTF-8 bytes of `0123456789ABCDEF`, RES key `87f879f135f27da7`,
md5 secret `daye,raolewoba!`).

---

## 1. InitCaptchaV3 request (user's live fetch example)

```
POST https://no8xfe.captcha-open-southeast.aliyuncs.com/
Content-Type: application/x-www-form-urlencoded; charset=UTF-8

body:
AccessKeyId=LTAI5tSEBwYMwVKAQGpxmvTd&SignatureMethod=HMAC-SHA1&SignatureVersion=1.0&Format=JSON&Timestamp=2026-08-29T12%3A48%3A37Z&Version=2023-03-05&Action=InitCaptchaV3&SceneId=didk33e0&Language=en&Mode=popup&UpLang=true&DeviceData=TEQYvgJq1LrMqFaBybfIzPxz2ygFyAct7X%2Fw%2BLacfXWd9rGSwE%2Fx6ZCONucD1fehS2Qpig6tUVsFK111d9wIk5pWp6rwYjzFCRgL7pNp8bzGsvOSdUXgQTopQm90YPSdCiRAlgENdODLvY7P8jrfO9eC15tPCPwLxcRIrcspVvQYqVfk9%2FyFeIlePKmTRjkM&SignatureNonce=aae91cc7-685b-4193-887a-ca47f3ee1d00&Signature=Signature
```

Raw (unencoded) DeviceData param from that body:

```
TEQYvgJq1LrMqFaBybfIzPxz2ygFyAct7X/w+LacfXWd9rGSwE/x6ZCONucD1fehS2Qpig6tUVsFK111d9wIk5pWp6rwYjzFCRgL7pNp8bzGsvOSdUXgQTopQm90YPSdCiRAlgENdODLvY7P8jrfO9eC15tPCPwLxcRIrcspVvQYqVfk9/yFeIlePKmTRjkM
```

## 2. InitCaptchaV3 response (live)

```json
{
    "CertifyId": "QOpv38o7Qs",
    "Message": "success",
    "RequestId": "04DE4994-1E5A-4387-81FB-2F8020F81B56",
    "Code": "Success",
    "LimitFlow": false,
    "Success": true,
    "StaticPath": "3.29.0/pe.064.91720ad47545c881",
    "CaptchaType": "TRACELESS",
    "DeviceConfig": "wN32kvF2iZuLdXqR9SvIIDy5WLu2p9xzJPVdtSxOaTAAFuQhfAzge1b2KK6jxvBs8dGMNR++31iG+CCS6kMRf9JYA4zuU4B5PpOiPliAnU23UQw5FN+DGQQeOwJ48T3HJy9LrGS4uz2g49qaDhLhatXK3dUOKqahUgwpj5hsLPx3qWfVeSCA8hKNrk54gtBkCoox4qBXZMbYM7zfOIU3rIskHnDS6ODyK+wX03GIGTRXDivevyfMLIsU73fFeBPic2vVwcgAhn6YQHo3ZPF5OvxCp39IYwlpND1oJqCg4jKJaCAzAy7FjHP+jXxansE8"
}
```

## 3. DeviceConfig decrypted (captain-verified, RES key `87f879f135f27da7`)

AES-128-CBC decrypt of the §2 `DeviceConfig` blob, split on `#`:

| field | value |
|---|---|
| f[0] | `NjgyZWZkMTgxNDllMzlhNg==` → base64-decode → **sessionKey `682efd18149e39a6`** |
| f[1] | `MQ==` (= 1, switch) |
| f[2] | **sessionId** `3795d28242a11619bc25f786f84e53d4-h-1788007723938-9472a56a25dc49bc8d669b22c6db0142` |
| f[3] | version `1.5.1/feilin001.874f974c24cb17ca9480aadc03f6652a9f7e071628b484381d9efc0060379b50` |
| f[4] | `` |
| f[5] | `` |
| f[6] | `` |
| f[7] | serverTs `1788007723939` |
| f[8] | ip `106.219.217.208` |
| f[9] | `0` |

## 4. Live deviceTokens (both md5-verified; `w` decrypts with `682efd18149e39a6`)

### TOKEN1

```
U0dfV0VCIzM3OTVkMjgyNDJhMTE2MTliYzI1Zjc4NmY4NGU1M2Q0LWgtMTc4ODAwNzcyMzkzOC05NDcyYTU2YTI1ZGM0OWJjOGQ2NjliMjJjNmRiMDE0MiN6a1hnQTBlZEliVDR0RkN2Ly94VHlvTmVQRC9YbEZzbWNUZTRVOWI3bEtDK1gxbmFmeG93a1oxYVRrWHFIRHhWUnpKYzZvcnZaUlorYmxtbzZzQ29pdDh1VEx6eWR6MzlPSGFBelpoR216TlMwcS9xYkc4R2pRS1JTdDdjTTBianpIZGE2UE52U1A1TmZvVHlNWGVadmV3a2N0OFNYSEQyWHVVWHRxanQweXFkTWdBWWVvc2w5M0k4NkVNNUc2NTVjS0hlamJnVFdSQ1dvb2xjTWVnOTlleU9qeEl5TXBXbE50SHQ5SkMybDcxOW9FandkZlpIVWpNUEJuOW9sT2JyYnd5UVdwdG5lQkNYQVJCaGpEMGNLeDIxN1VRTDBocGZ0MVMweDRROXZxZWlSSVVVN2FEN2kyNksxVUJ0T3VocVR2OWdPMUE3dVdpYzg5R0YwK2s3YWh2RjJMRWFWOEhCZkxCZHNSVEtmaG0rZXg3MWJVQ0hHWjF1WGlGUSszeTdPd0dOd2JMcVAzQkRSZUNVcWNuMStOdTBXT1kreldlZlZ0NnhrVldNQ3A1ZzZVYU1lL3VHd0d6M3crN3FCcUdIbG1DVTdGdlJOMVZiM2V6VEFQWlJtSU1ZUllzd3h2YXhheGJlMlE5Skx3UFdlOEg4bGxiWkxXZExYTFlocDNKWER1SDhGbUpES3AzRXJtVm1URGVxNkZkMkIvSFcrc0cxNExDU0V2UnBObTVFQ1FWREFZd0xQejFPZ3R3a0lqTS9zSFZpeXJGclBiY0krZys4OGxZdWQ4V3ZFZk1lc1lSVWxWZHhyb25MOHo3OGtHVWVtK25Sc2VVcVY1ZDlONjVqL2dPWnNtNm51dSsvMWt2aFgrbzdvdndzL2tDQTc0OGc1WFpSMUJCOHliZGZCOHBpcmlMUkVHR2lQSjEzYkw3WUZjUGRGSlNMRE1GV2srYmdhbDNpNXlvcHUzTGV5Y3hHRGJ3em5aQ3VyWVlhYXNtVkx4TURDU1VwOWpFYVdTOXcxOXd2TWdPYlpadHVWQXZ6eERscyt2dE82QmJ4RkZxNzJ2QjluOXVxSGJWMnhURlk0NUJxay9hMDhTQWpTVkFOMkJzZ1lKRFNBbGlmZWV4UHNPS2tDY1ZTU05wWTU1UlVxK0ozTkl4S1duST0jNDA1OCNhZjQzNDI5YzY4NjkzZTlhN2I1ZjIzMmM1OThlZGVjNQ==
```

### TOKEN2

```
U0dfV0VCIzM3OTVkMjgyNDJhMTE2MTliYzI1Zjc4NmY4NGU1M2Q0LWgtMTc4ODAwNzcyMzkzOC05NDcyYTU2YTI1ZGM0OWJjOGQ2NjliMjJjNmRiMDE0MiN6a1hnQTBlZEliVDR0RkN2Ly94VHlvTmVQRC9YbEZzbWNUZTRVOWI3bEtDK1gxbmFmeG93a1oxYVRrWHFIRHhWUnpKYzZvcnZaUlorYmxtbzZzQ29pdDh1VEx6eWR6MzlPSGFBelpoR216TlMwcS9xYkc4R2pRS1JTdDdjTTBianpIZGE2UE52U1A1TmZvVHlNWGVadmV3a2N0OFNYSEQyWHVVWHRxanQweXFkTWdBWWVvc2w5M0k4NkVNNUc2NTVjS0hlamJnVFdSQ1dvb2xjTWVnOTlleU9qeEl5TXBXbE50SHQ5SkMybDcxOW9FandkZlpIVWpNUEJuOW9sT2JyYnd5UVdwdG5lQkNYQVJCaGpEMGNLeDIxN1VRTDBocGZ0MVMweDRROXZxZWlSSVVVN2FEN2kyNksxVUJ0T3VocU9xYXF2cm5KVGhDbC9oRGVPMEJ0YU52cWE1V3RTSlBKTms4S1M4ZTJIcEtnSWoydGlybnI0bVJSMGpheW9IaTJLVjVHakVxRmhpZlV4a09aS21rL3pRRkRhTGZTUEJJb2pjbVRmMEQ4SE5xNXpGa2FEZXVCbUxKR1BEVFEyMW9mZ1hRcld3MmxOYWRsUzFBazNWS0JCeFlwMyt5MmtiaUtyQVRaWU5sbm9SaDJza29hUnJmTlJVREl6UXZiQ3FhUTAyZnVJamtaWDhCeVVoRi9NOXZkaEhFR216MkkxY3FqRG5hTCt4NnJhUUpJbjVjWlBKMUU3cmRXam94ZnNmbGNHVFQwazcwZFk0cE11N0YwVVpacmV3bHRERDVXbTBLb24zVkNMTzdUV2pCSWdYbWRnS0VtcTNxR1BlZ05jWm5vc1Q1eUpJV0V0U2pNUVg1TEptUEE3L2tEcjUvVW4vYW5Ub3k0ZVZqS0NZTk1sRlJXUmMvc0NaWmhnSFFVWDBDaGExWFhiMWU3MCtQRW1xelFPRUsxcG9ybEcwbDNETWtmK0NHN2xad0NRU1BmUGZGK1Nsajd3V2N0WnBNL3B4c0h2RnBZbG9YVkxzN2xnSVpYcHEwUVRSNUcvcndTQ0FoM0JUaXdoeEpxdDM1RVNHczZSbGhtR2NPSll1TlBqNXE4VVVOM2l2SlFtNnpFdUJna1NTUUFvbkl1cjYyV3d5TVFGUkoxTnBzKzRmbz0jNDA1OCM1MTI4YzM0N2VjMTQ1MmQ5MjZjYjhlZTI1YjRlMDE2MQ==
```

### Token facts (captain-verified)

| field | TOKEN1 | TOKEN2 |
|---|---|---|
| tF | `SG_WEB` | `SG_WEB` |
| Q | `3795d28242a11619bc25f786f84e53d4-h-1788007723938-9472a56a25dc49bc8d669b22c6db0142` | same |
| w | 876-char b64, 656-byte ciphertext | 876-char b64, 656-byte ciphertext |
| tC | `4058` | `4058` |
| p | `af43429c68693e9a7b5f232c598edec5` | `5128c347ec1452d926cb8ee25b4e0161` |
| md5 check | PASS | PASS |

Plaintext inside `w` is **648 bytes after PKCS7-strip** (656 = ciphertext length; 139 `#`
separators → 140 fields).

Payload (decrypted `w`) differences TOKEN1→TOKEN2 — **the ONLY diffs are payload fields 43 and 74**; everything else is byte-identical (fields at index 93/94 are empty strings — the `93-…`/`94-…` values below are log ENTRIES inside field 43, not standalone fields):

- field 43 (logs, `probeId-cost` pairs joined with `|`): `…|81-13514|93-594050|94-594065` → `…|81-13514|93-620724|94-620733`
- field 74: `1788008311383` → `1788008338057`

Verified arithmetic (token-engineer, re-confirmed): **log-entry-93 cost = field74 − field72 (initTime)**:
`1788008311383 − 594050 = 1788007717333` (= field 72 of both tokens) and
`1788008338057 − 620724 = 1788007717333`. Log-entry-94 = entry-93 + 9–15 ms jitter.

Timing relations and invariants (validator-measured across both live + §21.4 reference sessions):

- **Q's embedded ts == serverTs − 1** (live: 1788007723938 = 1788007723939−1; ref: 1787994141656 = 1787994141657−1). Q and serverTs come from the same server response — treat them as linked, never independent.
- **f72 (initTime) < serverTs** by 4–7 s (live 6606 ms; ref 4126 ms). Replays must keep a realistic few-second gap.
- **f74 is unbounded vs serverTs in real usage** — ref: f74−serverTs = 1698 ms; LIVE: 587444 ms (~9.8 min; page left open between init and token mint). So f74 and log-93/94 (= f74−f72) can legitimately be huge; the "quick interaction" constraint is weak. Fresh synthetic tokens should still stay in the 1–10 s range to look normal.
- **f87 == serverTs exactly**; **f42 == DeviceConfig f[8] (ip) exactly**; **Q == DeviceConfig f[2] exactly** — hard replay invariants in both sessions.
- **tC (GatherCost) is per-SESSION constant**: 4058 in both live tokens (minted 26 s apart); §12.1: 524; §21.4: 3496. NOT a function of token-mint time, and NOT max(log costs) (live entry `81-13514` > tC 4058) — independent collector metric frozen at collection.
- w plaintext = **648 bytes + 8 bytes PKCS7** → 656-byte ciphertext; 139 `#` separators → 140 fields.

**DeviceData generator params for the live site** (decrypt of the §1 request's DeviceData param, KEY_HE→KEY_O): sceneId `didk33e0`, prefix `no8xfe`, region `sgp`, appKey `3795d28242a11619bc25f786f84e53d4` (NOT the `ab034e…` static default in generate_device_data.py — always pass appKey explicitly).

---

## Appendix — captain's decrypt of TOKEN1 `w` (for cross-checking your own work)

Plaintext (140 `#`-fields, 656 bytes → 625 visible after trim; regenerate yourself with
`AES-128-CBC(key=b'682efd18149e39a6', iv=b'0123456789ABCDEF')` over `base64.b64decode(w)`):

```
W.10054#####Linux armv81#Chrome#149.0.0.0#############8#E2FFHTUwY2c=#4##########50a98af28ee10a81ef6f31efaae2853b##8##Linux#x86_64#####106.219.217.208#10-0|20-2847|11-3546|23-7686|30-7697|40-7759|41-12134|70-12139|71-13504|80-13504|81-13514|93-594050|94-594065#true###768*1366##5##################saf-captcha#0###4SyHGkKVW8fUJZTYIWWLKAoWeJmau7PQKrOjC8GP#1788007717333#O9l1RIX7GMGG35xDAGAL5y9Zq9vHKgI8LAOBLElIU7#1788008311383#desktop#false##9d4568c009d203ab10e33ea9953a0264#########1788007723939#MCMwIzAjMCMwIzAjMCMwIzAjMSMwIzAjMCMwIzAjMCMwIzAjMCMxIzEjMCMxMTExMTExMDExMTExMTExMTExMTExMTExMQ==#1#1#true##################0##############################
```

(§21.4 of Fielin-report.md holds a second, older reference session — key `57ad9f73260d1d46`, sessionId `…-h-1787994141656-52b528da…`, ip `223.188.28.68` — useful as a third data point for cross-session field comparison.)

---

# Session-4 captures — 2026-10-21 (device service, headful Chromium 120)

Captured live with the rewritten `cmd/token-collector` (one browser visit → a
trusted device session + the genuine collector payload) and cross-checked with
`.assets/scripts/live_validate.py --session-mint`.

Everything below was produced by the collector itself and re-verified offline
(AES-128-CBC, IV = UTF-8 bytes of `0123456789ABCDEF`).

## S4.1 The device service's `Log1` request (byte-identical to ours)

```
POST https://cloudauth-device-dualstack.ap-southeast-1.aliyuncs.com/
Content-Type: application/x-www-form-urlencoded; charset=UTF-8

body:
AaduaneId=DuaneAprqkYsF3nt1yjK29Bf&Version=2020-10-15&SignatureMethod=HMAC-SHA1&SignatureVersion=1.0&Format=JSON&Action=Log1&Data=TEQYvgJq1LrMqFaBybfIzPxz2ygFyAct7X%2Fw%2BLacfXWd9rGSwE%2Fx6ZCONucD1fehrtHiD2ADpEIEuf3%2F3x5GPEHZiSr3V6CUS5gq97FCPvfQRPQJb88%2FuidhHZLBaYMgHjb1hJo0dxcWqO7mL%2FmeqWHPlzNxKfuvTWzLV5%2Fj%2BMg%3D&SignatureNonce=36624f64-57cc-4c64-aac0-1234a3258b9d&Signature=3wr4p6WXFhff56ewS0wz2xjtCs4%3D
```

Wire headers actually sent by the browser (CDP `Network.requestWillBeSentExtraInfo`
— note: **no `Origin`, no `Referer`, no `Cookie`**):

```
Accept: */*                      Accept-Encoding: gzip, deflate, br
Accept-Language: en-US,en;q=0.9  Connection: keep-alive
Content-Length: 413              Content-Type: application/x-www-form-urlencoded; charset=UTF-8
Host: cloudauth-device-dualstack.ap-southeast-1.aliyuncs.com
Sec-Fetch-Dest: empty            Sec-Fetch-Mode: cors
Sec-Fetch-Site: cross-site       User-Agent: Mozilla/5.0 (X11; Linux x86_64) … Chrome/120.0.0.0 …
sec-ch-ua: "Not_A Brand";v="8", "Chromium";v="120"
sec-ch-ua-mobile: ?0             sec-ch-ua-platform: "Linux"
```

The `Data` param is **static** (it does not carry a timestamp, a nonce or the
fingerprint) and is byte-identical to the one `recreation.py` builds:

```
appKey#W#AES(FLAG_KEY,"W.10001.c#saf-captcha##captcha-front#no8xfe#sgp")#W20220202#CLOUD#
```

## S4.2 `blob32` is DERIVED, not a constant

The device service's upload head carries a 32-byte blob that the older analysis
treated as an opaque, rotating constant.  It is in fact:

```
blob32 = AES-128-CBC(sessionKey, IV, PKCS7("W.10055#saf-captcha#didk33e0"))
```

Verified on three independent live captures (three different session keys, three
different blobs, all reproduced exactly).  The `sceneId` is the second part, so
the blob is per-session — the hard-coded value that used to live in the code
belonged to a *different* session and was therefore always wrong.

## S4.3 The upload formats (Log2 / Log3), byte-exact

```
Log2 head : appKey#W#<blob32>#W20220202#CLOUD#<GatherCost>#501#<payload>
Log2 inner: sessionId#b64(AES(sk,RAW_DEVICE_DATA))#AES(sk,'saf-captcha')
                    #AES(sk,'W.10055')##AES(sk,ts)

Log3 head : appKey#W#<blob32>#W20220202#CLOUD#<cost>#<payload>
Log3 inner: 511#b64(A)#b64(B)
     A = sessionId##AES(sk,'saf-captcha')#AES(sk,'W.10055')##AES(sk,ts)
     B = sessionId#b64(AES(sk,BEHAVIOUR))#AES(sk,'saf-captcha')
             #AES(sk,'W.10055')##AES(sk,ts)
```

Notes that cost us time and are now fixed in the code:

* the fingerprint is base64-encoded **exactly once** — the previous
  implementation base64'd it twice, so the server never saw a fingerprint at all;
* `Log3`'s head field 5 is **empty** (the older code sent `2001`);
* `Log2`'s head field 5 is the **GatherCost** and must equal the token's `tC`;
* the behaviour log carries the **User-Agent**, so it must match the UA the
  request was sent with;
* with the old `Log3` the device service answered **HTTP 500**; with the
  correct one it answers **200**.

## S4.4 Live session used for the Session-4 measurements

```
sessionId  : 3795d28242a11619bc25f786f84e53d4-h-1792574091284-d5fea567b1404c00a0a992100bb04571
sessionKey : 3780c2097f062881
serverTs   : 1792574091286      (DeviceConfig f[7])
ip         : 223.188.38.120
tC         : 2484
payload    : 155 fields (the "warm" variant; a cold browser yields 111)
```

**HC1 relaxation (live finding):** the sessionId's embedded ts trails the
DeviceConfig's `serverTs` by **1 s** in the archived capture above and by **2 s**
in this one.  The "embedded == serverTs − 1" equality is therefore an artefact
of a single sample; the invariant is a window (≤ 3 s).  Both
`device_token_gen.build_payload` and `live_validate.build_payload_live` now
enforce the window instead of the equality.

## S4.5 One session → N tokens

`live_validate.py --session-mint 3 --session <collector export>` (same session,
one `InitCaptchaV3` + one `VerifyCaptchaV3` per token):

```
session#1: … -> T001 (accepted=True)
session#2: … -> T001 (accepted=True)
session#3: … -> T001 (accepted=True)
SESSION-MINT RESULT: 3/3 tokens accepted (same device session, own CertifyId each)
```

A long-lived session was measured separately: **18/18 mints over 35 minutes,
all `T001`, with no browser process running** (the browser had already exited
before the first mint).
