# Aliyun CaptchaJS Information Report 

# Obfuscation & Deobfuscation of literals [File:AliyunCaptcha.js] 

Starting with 
```javascript
var G = rr; // <- Deobfuscator function
window.__G = G = rr; // Added proxy to access it globally 
```

### Function rr Defintion 
```javascript
function rr(t, r) {

            var e = jt();

            return rr = function(r, n) {

                var i = e[r -= 312];

                if (void 0 === rr.avoPoV) {

                    rr.nyBjcA = function(t) {

                        for (var r, e, n = "", i = "", o = 0, c = 0; e = t.charAt(c++); ~e && (r = o % 4 ? 64 * r + e : e,

                        o++ % 4) ? n += String.fromCharCode(255 & r >> (-2 * o & 6)) : 0)

                            e = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+/=".indexOf(e);

                        for (var u = 0, a = n.length; u < a; u++)

                            i += "%" + ("00" + n.charCodeAt(u).toString(16)).slice(-2);

                        return decodeURIComponent(i)

                    }

                    ,

                    t = arguments,

                    rr.avoPoV = !0

                }

                var o = r + e[0]

                  , c = t[o];

                return c ? i = c : (i = rr.nyBjcA(i),

                t[o] = i),

                i

            }

            ,

            rr(t, r)

        }


```

### Jt  


```javascript
// Jt.prototype
Jt[G(547) + "e"] = {
            ENDPOINTS: [G(525) + G(417) + G(473) + G(405) + G(578)],
            CN_DEFAULT_ENDPOINTS: [G(525) + G(417) + G(473) + G(405) + G(578)],
            INTL_DEFAULT_ENDPOINTS: [G(525) + G(417) + G(473) + G(555) + G(362) + G(358) + G(366)],
            CN_ENDPOINTS: _t,
            INTL_ENDPOINTS: St,
            WAF_ENDPOINTS: [G(525) + G(475) + G(344) + G(377) + G(556)],
            cdnServers: [G(535) + G(366)],
            cdnDevServers: [G(501) + G(581)],
            dynamicJsPath: function(t) {
                var r = 479
                  , e = 318
                  , n = 342
                  , i = 328
                  , o = 388
                  , c = 507
                  , u = 444
                  , a = 494
                  , s = 318
                  , f = 507
                  , l = G
                  , p = {};
                p[l(494)] = function(t, r) {
                    return t + r
                }
                ,
                p[l(r)] = function(t, r) {
                    return t + r
                }
                ,
                p[l(e)] = l(n) + l(i) + l(o) + "/",
                p[l(c)] = l(u);
                var v = p;
                return v[l(a)](v[l(r)](v[l(s)], t), v[l(f)])
            },
            fallbackVersion: G(319) + G(499),
            https: G(525),
            http: G(579),
            API_VERSION: G(590) + "15",
            APP_VERSION: G(416) + "2",
            PLATFORM: G(398) + "c",
            APP_NAME: G(484) + G(508),
            DEVICE_TYPE: Vt,
            APP_KEY: G(345) + G(526) + G(424) + G(322),
            ACCESS_KEY: Yt,
            WEB_AES_SECRET_KEY: Zt,
            AES_IV: G(534) + G(427) + G(586),
            SALT: G(468) + G(549) + G(350),
            SESSION_ID_SALT: G(409) + G(533) + G(584),
            ACCESS_SEC: G(503) + G(396),
            ACTION: Xt,
            ACTION_STATE: Qt,
            WEB_REGION: $t,
            WEB_REGION_PREID: tr,
            UID_NAME_COOKIE: G(567) + "o",
            UID_NAME_LOCAL: G(497) + "s",
            initTime: Date[G(496)](),
            preCollectData: {},
            logs: [],
            _extend: function(t) {
                var r = G
                  , e = this;
                new q(t)[r(505)](function(t, r) {
                    e[t] = r
                })
            }
        };
```        


#### Deobfuscation examples
```js
var G = rr;
window.__G = G;
```
```js
var G = window.__G;
```
```js
PLATFORM: G(398) + "c"
'W.10001.c'

fallbackVersion: G(319) + G(499)
'0.0.0/feilin'

https: G(525)
'https://'

http: G(579)
'http://'

API_VERSION: G(590) + "15"
'2020-10-15'

APP_VERSION: G(416) + "2"
'W20220202'

APP_NAME: G(484) + G(508)
'saf-aliyun-com'

APP_KEY: G(345) + G(526) + G(424) + G(322)
'ab034ec0643f91399eb33e062dc7fae1'

AES_IV: G(534) + G(427) + G(586)
'd35db7e39ebbf3d001083105'

SALT: G(468) + G(549) + G(350)
'NLAoqT6K03oLbQXW2VS3zA=='

SESSION_ID_SALT: G(409) + G(533) + G(584)
'X1y5VstbB5zqghyJ9g8a0A=='

ACCESS_SEC: G(503) + G(396)
'FqJB6iRNVYdEGpwb'

UID_NAME_COOKIE: G(567) + "o"
'_c_WBKFRo'

var Zt = {};

        Zt[G(510)] = G(374) + G(410) + G(333) + G(546) + G(346) + G(399),
        Zt[G(407)] = G(321) + G(537) + G(426) + G(415) + G(593) + G(521),
        Zt[G(348)] = G(464) + G(516) + G(437) + G(325) + G(488) + G(449),
        Zt[G(564)] = G(553) + G(480) + G(577) + G(365) + G(340) + G(561),
        Zt[G(500)] = G(457) + G(403) + G(433) + G(571) + G(400) + G(341);
var WEB_AES_SECRET_KEY = Zt;
undefined
console.log(WEB_AES_SECRET_KEY)
{

    "REQ": "8KmHIQsc5+LZJA7uYex3WaHdkjgCtS6epbG/bc9xss0=",

    "RES": "9NhnQQ+LRrKCkAuxwZaUWGBtSzaFtFlNb/ksJCrCgrM=",

    "FLAG": "k+1RW0cz3iDi2RAbC/c3QKzTPiVwNmNO1910DXe6Gas=",

    "UPLOAD": "+fR9tYzlKFr07pEbumd7+KnO3xLOkphCS+qKUbJiMfA=",

    "PREID": "xLLw/t15vkI7QQBX/1scBbcb9fKx+ymxF0tJ3ds42B0="

}


var Vt = {};
Vt[G(441)] = "W";
'W'
var DEVICE_TYPE = Vt;
undefined
console.log(DEVICE_TYPE);
{

    "WEB": "W"

}

```

---

## Functions , Definitions & Variables [AliyunCaptcha.js]

 ```javascript  
   function me(t, r) {
            var e = 499
              , n = 531
              , i = 368
              , o = 487
              , c = 423
              , u = 373
              , a = 468
              , s = 487
              , f = 487
              , l = 436
              , p = 373
              , v = 429
              , h = 465
              , d = 376
              , y = re
              , g = {};
            g[y(e)] = y(n) + y(i),
            g[y(o)] = function(t, r) {
                return t === r
            }
            ,
            g[y(c)] = function(t, r) {
                return t !== r
            }
            ,
            g[y(u)] = function(t, r) {
                return t <= r
            }
            ;
            for (var m = g, x = m[y(e)][y(a)]("|"), w = 0; ; ) {
                switch (x[w++]) {
                case "0":
                    if (m[y(s)](r, void 0) || m[y(f)](r, null))
                        return null;
                    continue;
                case "1":
                    if (m[y(o)](t, void 0) || m[y(c)](t[y(l)], 16) || m[y(p)](r[y(l)], 0))
                        return null;
                    continue;
                case "2":
                    var b = r;
                    continue;
                case "3":
                    var S = ue[y(v)](t);
                    continue;
                case "4":
                    return C[y(h)](ue);
                case "5":
                    var C = ce[y(d)](b, S, pe);
                    continue
                }
                break
            }
        }
```
```javascript
var ce = te()[re(454)]
          , ue = te()[re(439)][re(383)]
          , ae = te()[re(439)][re(505)]
          , se = te()[re(439)][re(517)]
          , fe = te()[re(504)][re(494)]
          , le = ae[re(397) + "y"](se[re(429)](dt))
          , pe = {
            iv: ue[re(429)](le),
            padding: fe
        }
          , ve = nr[re(523) + re(521) + "EY"]
          , he = me(nr[re(541) + "EC"], ve[re(435)])
          , de = me(nr[re(541) + "EC"], ve[re(483)]);
```
```console
window.__value_of_ve
{
    "REQ": "8KmHIQsc5+LZJA7uYex3WaHdkjgCtS6epbG/bc9xss0=",
    "RES": "9NhnQQ+LRrKCkAuxwZaUWGBtSzaFtFlNb/ksJCrCgrM=",
    "FLAG": "k+1RW0cz3iDi2RAbC/c3QKzTPiVwNmNO1910DXe6Gas=",
    "UPLOAD": "+fR9tYzlKFr07pEbumd7+KnO3xLOkphCS+qKUbJiMfA=",
    "PREID": "xLLw/t15vkI7QQBX/1scBbcb9fKx+ymxF0tJ3ds42B0="
}


window.__value_he
'45f8ac1e1de14397'

window.__value_de
'87f879f135f27da7'

window.__se
{stringify: ƒ, parse: ƒ}
window.__se.stringify
, u = n.enc = {}
                  , a = u.Hex = {
                    stringify: function(t) {
                        for (var r = t.words, e = t.sigBytes, n = [], i = 0; i < e; i++) {
                            var o = r[i >>> 2] >>> 24 - i % 4 * 8 & 255;
                            n.push((o >>> 4).toString(16)),
                            n.push((15 & o).toString(16))
                        }
                        return n.join("")
                    },
		    
window.__se_parse
parse: function(t) {
                        for (var r = t.length, e = [], n = 0; n < r; n += 2)
                            e[n >>> 3] |= parseInt(t.substr(n, 2), 16) << 24 - n % 8 * 4;
                        return new c.init(e,r / 2)
                    }
		    

window.__ce
{encrypt: ƒ, decrypt: ƒ}


context of ce:
        7165: function(t, r, e) {
            var n;
            t.exports = (n = e(9021),
            e(9506),
            void (n.lib.Cipher || function(t) {
                var r = n
                  , e = r.lib
                  , i = e.Base
                  , o = e.WordArray
                  , c = e.BufferedBlockAlgorithm
                  , u = r.enc
                  , a = (u.Utf8,
                u.Base64)
                  , s = r.algo.EvpKDF
                  , f = e.Cipher = c.extend({
                    cfg: i.extend(),
                    createEncryptor: function(t, r) {
                        return this.create(this._ENC_XFORM_MODE, t, r)
                    },
                    createDecryptor: function(t, r) {
                        return this.create(this._DEC_XFORM_MODE, t, r)
                    },
                    init: function(t, r, e) {
                        this.cfg = this.cfg.extend(e),
                        this._xformMode = t,
                        this._key = r,
                        this.reset()
                    },
                    reset: function() {
                        c.reset.call(this),
                        this._doReset()
                    },
                    process: function(t) {
                        return this._append(t),
                        this._process()
                    },
                    finalize: function(t) {
                        return t && this._append(t),
                        this._doFinalize()
                    },
                    keySize: 4,
                    ivSize: 4,
                    _ENC_XFORM_MODE: 1,
                    _DEC_XFORM_MODE: 2,
                    _createHelper: function() {
                        function t(t) {
                            return "string" == typeof t ? x : g
                        }
                        return function(r) {
                            return {
                                encrypt: function(e, n, i) {
                                    return t(n).encrypt(r, e, n, i)
                                },
                                decrypt: function(e, n, i) {
                                    return t(n).decrypt(r, e, n, i)
                                }
                            }
                        }
                    }()
                })




window.__fe
{pad: ƒ, unpad: ƒ}
window.__fe.pad and window.__fe.unpad

                  , h = (r.pad = {}).Pkcs7 = {
                    pad: function(t, r) {
                        for (var e = 4 * r, n = e - t.sigBytes % e, i = n << 24 | n << 16 | n << 8 | n, c = [], u = 0; u < n; u += 4)
                            c.push(i);
                        var a = o.create(c, n);
                        t.concat(a)
                    },
                    unpad: function(t) {
                        var r = 255 & t.words[t.sigBytes - 1 >>> 2];
                        t.sigBytes -= r
                    }
                }

window.__le
'0123456789ABCDEF'

window.__ie
function ie(t, r) {
            var e = 469
              , n = 436
              , i = 421
              , o = 414
              , c = ge
              , u = {
                csLej: function(t, r) {
                    return t == r
                },
                tyWOd: function(t, r) {
                    return t > r
                },
                dUZcv: function(t, r) {
                    return t(r)
                },
                Tuisd: function(t, r) {
                    return t < r
                }
            };
            (u[c(360)](null, r) || u[c(e)](r, t[c(n)])) && (r = t[c(n)]);
            for (var a = 0, s = u[c(i)](Array, r); u[c(o)](a, r); a++)
                s[a] = t[a];
            return s
        }


window.__pe

{
    "iv": {
        "words": [
            808530483,
            875902519,
            943276354,
            1128547654
        ],
        "sigBytes": 16
    },
    "padding": {}
}

window.__ue
{stringify: ƒ, parse: ƒ}

window.__ue.stringify
, f = u.Utf8 = {
stringify: function(t) {
                        try {
                            return decodeURIComponent(escape(s.stringify(t)))
                        } catch (t) {
                            throw new Error("Malformed UTF-8 data")
                        }
                    },
window.__ue.parse
parse: function(t) {
                        return s.parse(unescape(encodeURIComponent(t)))
                    }


window.__ae // Base64 encoder/decoder?
{
    "_map": "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=",
    "_reverseMap": [
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        62,
        null,
        null,
        null,
        63,
        52,
        53,
        54,
        55,
        56,
        57,
        58,
        59,
        60,
        61,
        null,
        null,
        null,
        64,
        null,
        null,
        null,
        0,
        1,
        2,
        3,
        4,
        5,
        6,
        7,
        8,
        9,
        10,
        11,
        12,
        13,
        14,
        15,
        16,
        17,
        18,
        19,
        20,
        21,
        22,
        23,
        24,
        25,
        null,
        null,
        null,
        null,
        null,
        null,
        26,
        27,
        28,
        29,
        30,
        31,
        32,
        33,
        34,
        35,
        36,
        37,
        38,
        39,
        40,
        41,
        42,
        43,
        44,
        45,
        46,
        47,
        48,
        49,
        50,
        51
    ]
}


```

```javascript
function ye(t, r) {
            var e = 442
              , n = 406
              , i = 535
              , o = 385
              , c = 480
              , u = 492
              , a = 468
              , s = 479
              , f = 385
              , l = 436
              , p = 436
              , v = 429
              , h = 385
              , d = 465
              , y = re
              , g = {};
            g[y(e)] = y(n) + y(i),
            g[y(o)] = function(t, r) {
                return t === r
            }
            ,
            g[y(c)] = function(t, r) {
                return t !== r
            }
            ,
            g[y(u)] = function(t, r) {
                return t <= r
            }
            ;
            for (var m = g, x = m[y(e)][y(a)]("|"), w = 0; ; ) {
                switch (x[w++]) {
                case "0":
                    var b = ce[y(s)](C, S, pe);
                    continue;
                case "1":
                    if (m[y(f)](t, void 0) || m[y(c)](t[y(l)], 16) || m[y(u)](r[y(p)], 0))
                        return null;
                    continue;
                case "2":
                    var S = ue[y(v)](t);
                    continue;
                case "3":
                    if (m[y(f)](r, void 0) || m[y(h)](r, null))
                        return null;
                    continue;
                case "4":
                    var C = r;
                    continue;
                case "5":
                    return b[y(d)]()
                }
                break
            }
        }

```
```javascript
function be(t) {
            var r = 431
              , e = re;
            return {
                PMZDa: function(t, r, e) {
                    return t(r, e)
                }
            }[e(453)](ye, he, t[e(r)]("#"))
        }
```
```javascript

        var Ee = {
            ACTION: gt,
            ACTION_STATE: xt,
            KEY_ID: me(pt, ht.ID),
            KEY_SECRET: me(pt, ht[re(371)])
        }
```
```console
console.log(Ee)
{
    "ACTION": {
        "INIT": "InitCaptcha",
        "INITV2": "InitCaptchaV2",
        "INITV3": "InitCaptchaV3",
        "VERIFY": "VerifyCaptchaV2",
        "VERIFYV3": "VerifyCaptchaV3",
        "LOG": "UploadLog"
    },
    "ACTION_STATE": {
        "SUCCESS": "success",
        "FAIL": "fail"
    },
    "KEY_ID": "LTAI5tSEBwYMwVKAQGpxmvTd",
    "KEY_SECRET": "YSKfst7GaVkXwZYvVihJsKF9r89koz"
}
```

```javascript
  , ke = {
            ACTION: nr[re(476)],
            ACTION_STATE: nr[re(407) + re(509)],
            DEVICE_TYPE: nr[re(477) + re(388)],
            WEB_AES_SECRET_KEY: nr[re(523) + re(521) + "EY"],
            KEY_ID: me(nr[re(541) + "EC"], nr[re(485) + "EY"].ID),
            KEY_SECRET: me(nr[re(541) + "EC"], nr[re(485) + "EY"][re(371)]),
            WEB_AES_FLAG_SECRET_KEY: me(nr[re(541) + "EC"], nr[re(523) + re(521) + "EY"][re(524)])
        };

```
```console

{
    "ACTION": {
        "INIT": "Log1",
        "DEVICE_UPLOAD": "Log2",
        "COMBAT_UPLOAD": "Log3"
    },
    "ACTION_STATE": {
        "SUCCESS": "success",
        "FAIL": "fail"
    },
    "DEVICE_TYPE": {
        "WEB": "W"
    },
    "WEB_AES_SECRET_KEY": {
        "REQ": "8KmHIQsc5+LZJA7uYex3WaHdkjgCtS6epbG/bc9xss0=",
        "RES": "9NhnQQ+LRrKCkAuxwZaUWGBtSzaFtFlNb/ksJCrCgrM=",
        "FLAG": "k+1RW0cz3iDi2RAbC/c3QKzTPiVwNmNO1910DXe6Gas=",
        "UPLOAD": "+fR9tYzlKFr07pEbumd7+KnO3xLOkphCS+qKUbJiMfA=",
        "PREID": "xLLw/t15vkI7QQBX/1scBbcb9fKx+ymxF0tJ3ds42B0="
    },
    "KEY_ID": "LTAI5tGjnK9uu9GbT9GQw72p",
    "KEY_SECRET": "fpOKzILEajkqgSpr9VvU98FwAgIRcX",
    "WEB_AES_FLAG_SECRET_KEY": "c175a358550d02e2"
}
```

---

### Function Re internals [AliyunCaptcha.js]

```javascript
var nr = new Jt({});
// nr is related to Jt
```

In context of function Re i added 
```javascript
window.__value_of_r = r;
```

so the full function becomes
```javascript
function Re(t, r, e) {
            r._extend(Be({}, t));
            var n = t.appKey || r.APP_KEY
              , i = t.appName || r.APP_NAME
              , o = me(r.ACCESS_SEC, r.secretKey) || e.WEB_AES_FLAG_SECRET_KEY
              , c = r.PLATFORM + "#" + i + "#" + (r.sceneId || "") + "#captcha-normal#" + De.prefix + "#" + De.region;
            window.__value_of_r = r;
            return c = ye(o, c),
            window.__value_of_r_secretKey = r.secretKey;
	    console.log('Value of o', o)
            be([n, e.DEVICE_TYPE.WEB, c, r.APP_VERSION, "CLOUD", ""])
        }
```
      
Upon calling this function i get value of o = e.WEB_AES_FLAG_SECRET_KEY which means me(r.ACCESS_SEC, r.secretKey) is failing when i ran `console.log(window.__value_of_r_secretKey)` i got undefined and when i ran `console.log(window.__value_of_r.ACCESS_SEC)`
it worked  and logged 'FqJB6iRNVYdEGpwb'  which means `window.__value_of_r_secretKey` is failing when i logged  r.ACCESS_KEY i got 
```json
{
    "ID": "tBwmiXXwEdGaRtkcLvIoA/F6r4qktDu7lFi23FmDCbo=",
    "SECRET": "9U6Kg6ljgm9PhHvSsihjy1wJVo71uzWkRDfGiyiMVBg="
} 
```
This is just for information lets just assume  using e.WEB_AES_FLAG_SECRET_KEY since it works
For information here is the  full log of r: 

```json
{
    "preCollectData": {
        "fontsNum": 4
    },
    "ENDPOINTS": [
        "https://cloudauth-device-dualstack.cn-shanghai.aliyuncs.com",
        "https://cn-shanghai.device.saf.aliyuncs.com"
    ],
    "sceneId": "didk33e0",
    "appName": "saf-captcha",
    "appKey": "ab034ec0643f91399eb33e062dc7fae1",
    "endpoints": [
        "https://cloudauth-device-dualstack.cn-shanghai.aliyuncs.com",
        "https://cn-shanghai.device.saf.aliyuncs.com"
    ]
}

```
Some other important values related to r:
```javascript
r._extend
// Definition
_extend: function(t) {
                var r = G
                  , e = this;
                new q(t)[r(505)](function(t, r) {
                    e[t] = r
                })
            }
```
	    
#### Decryption Script

```python3
from Crypto.Cipher import AES
from Crypto.Util.Padding import unpad
import base64

def decrypt_ye(encrypted_b64, key_str, iv_words):
    """
    Decrypt data encrypted by the ye function
    """                                                                                              # Convert key string to bytes (AES-128 = 16 bytes)
    key = key_str.encode('utf-8')

    # Convert iv words to bytes
    iv_bytes = []
    for word in iv_words:
        # Each word is 32-bit integer, big-endian
        iv_bytes.extend([
            (word >> 24) & 0xFF,
            (word >> 16) & 0xFF,
            (word >> 8) & 0xFF,
            word & 0xFF
        ])
    iv = bytes(iv_bytes[:16])  # Take first 16 bytes

    # Base64 decode the encrypted data
    encrypted = base64.b64decode(encrypted_b64)

    # Decrypt using AES-CBC
    cipher = AES.new(key, AES.MODE_CBC, iv)
    decrypted_padded = cipher.decrypt(encrypted)

    # Remove PKCS7 padding
    decrypted = unpad(decrypted_padded, AES.block_size)

    # Return as string
    return decrypted.decode('utf-8')

# Your data
encrypted = "V0z3lM83c/TjjZu6HWeg95/xMYKFKoBmiixMa2THPUZkDnFK7GAZFGfreQXa7zu4FGY0zm1H4bU9uybpP2Fv0w=="
key = 'c175a358550d02e2'  # 16 bytes
iv_words = [808530483, 875902519, 943276354, 1128547654]

# Decrypt
result = decrypt_ye(encrypted, key, iv_words)
print(result)
```

Script is working

---

## Information of Rt in context of Rt.appkey [AliyunCaptcha.js]
```json
{
    "appName": {
        "1.0": "saf-captcha-waf",
        "2.0": "saf-captcha",
        "3.0": "saf-captcha"
    },
    "appKey": {
        "1.0": {
            "cn": "sh87bd1512hsb03cb405803a307dbe32",
            "sgp": "sg63c0a034gsf3f3ed381ac9c8a0bc51",
            "ga": "sg63c0a034gsf3f3ed381ac9c8a0bc51"
        },
        "2.0": {
            "cn": "ab034ec0643f91399eb33e062dc7fae1",
            "sgp": "3795d28242a11619bc25f786f84e53d4",
            "ga": "3795d28242a11619bc25f786f84e53d4"
        },
        "3.0": {
            "cn": "ab034ec0643f91399eb33e062dc7fae1",
            "sgp": "3795d28242a11619bc25f786f84e53d4",
            "ga": "3795d28242a11619bc25f786f84e53d4"
        }
    },
    "endpoints": {
        "1.0": {
            "cn": [
                "https://device.captcha-open.aliyuncs.com"
            ],
            "sgp": [
                "https://cloudauth-device-dualstack.ap-southeast-1.aliyuncs.com",
                "https://ap-southeast-1.device.saf.aliyuncs.com"
            ],
            "ga": [
                "https://cloudauth-device-dualstack.ap-southeast-1.aliyuncs.com",
                "https://ap-southeast-1.device.saf.aliyuncs.com"
            ]
        },
        "2.0": {
            "cn": [
                "https://cloudauth-device-dualstack.cn-shanghai.aliyuncs.com",
                "https://cn-shanghai.device.saf.aliyuncs.com"
            ],
            "sgp": [
                "https://cloudauth-device-dualstack.ap-southeast-1.aliyuncs.com",
                "https://ap-southeast-1.device.saf.aliyuncs.com"
            ],
            "ga": [
                "https://ap-southeast-1-ga.device.saf.aliyuncs.com",
                "https://cloudauth-device-dualstack.ap-southeast-1.aliyuncs.com"
            ]
        },
        "3.0": {
            "cn": [
                "https://cloudauth-device-dualstack.cn-shanghai.aliyuncs.com",
                "https://cn-shanghai.device.saf.aliyuncs.com"
            ],
            "sgp": [
                "https://cloudauth-device-dualstack.ap-southeast-1.aliyuncs.com",
                "https://ap-southeast-1.device.saf.aliyuncs.com"
            ],
            "ga": [
                "https://ap-southeast-1-ga.device.saf.aliyuncs.com",
                "https://cloudauth-device-dualstack.ap-southeast-1.aliyuncs.com"
            ]
        }
    },
    "shplEndpoints": {
        "1.0": {
            "cn": [
                "https://device.captcha-open.aliyuncs.com"
            ],
            "sgp": [
                "https://cloudauth-device-dualstack.ap-southeast-1.aliyuncs.com",
                "https://ap-southeast-1.device.saf.aliyuncs.com"
            ],
            "ga": [
                "https://cloudauth-device-dualstack.ap-southeast-1.aliyuncs.com",
                "https://ap-southeast-1.device.saf.aliyuncs.com"
            ]
        },
        "2.0": {
            "cn": [
                "https://cn-shanghai.device.saf.aliyuncs.com"
            ],
            "sgp": [
                "https://ap-southeast-1.device.saf.aliyuncs.com"
            ],
            "ga": [
                "https://ap-southeast-1-ga.device.saf.aliyuncs.com"
            ]
        },
        "3.0": {
            "cn": [
                "https://cn-shanghai.device.saf.aliyuncs.com"
            ],
            "sgp": [
                "https://ap-southeast-1.device.saf.aliyuncs.com"
            ],
            "ga": [
                "https://ap-southeast-1-ga.device.saf.aliyuncs.com"
            ]
        }
    }
}

```


// Simplified version of what the code does in real:

```javascript
var Rt = {
    appKey: {
        "1.0": {
            "cn": G(569) + G(481) + G(363) + G(565),      // sh87bd15 + 12hsb03c + b405803a + 307dbe32
            "sgp": G(434) + G(490) + G(523) + G(498),     // sg63c0a0 + 34gsf3f3 + 7ef9e8a2 + c8a0bc51
            "ga": G(434) + G(490) + G(523) + G(498)       // Same as sgp
        },
        "2.0": {
            "cn": G(345) + G(526) + G(424) + G(322),      // ab034ec0 + 643f9139 + 9eb33e06 + 2dc7fae1
            "sgp": G(337) + G(589) + G(588) + G(455),     // 3795d282 + 42a11619 + bc25f786 + f84e53d4
            "ga": G(337) + G(589) + G(588) + G(455)       // Same as sgp
        },
        "3.0": {
            "cn": G(345) + G(526) + G(424) + G(322),      // Same as 2.0
            "sgp": G(337) + G(589) + G(588) + G(455),     // Same as 2.0
            "ga": G(337) + G(589) + G(588) + G(455)       // Same as 2.0
        }
    }
}
```


### Calculation  of literals using Deobfuscation function rr:

#### Step 1:  The jt() array contains the raw string fragments
```javascript
function jt() {
    var t = ["Ac1KzxzPy2u", "DxrOzwfZDc0", ... "ab034ec0", "643f9139", ...];
    return (jt = function() { return t })();
}
```
#### Step 2: The decoder function rr (assigned to G)

Min Index is 312 of function rr
Max Index is 595 of function rr

```javascript
var G = rr;  // At the top of the file

function rr(t, r) {
    var e = jt();  // Gets the string array
    return rr = function(r, n) {
        var i = e[r -= 312];  // Subtract 312 to get array index
        if (void 0 === rr.avoPoV) {
            rr.nyBjcA = function(t) {
                // Base64-like decoding function
                for (var r, e, n = "", i = "", o = 0, c = 0; e = t.charAt(c++); ~e && (r = o % 4 ? 64 * r + e : e, o++ % 4) ? n += String.fromCharCode(255 & r >> (-2 * o & 6)) : 0)
                    e = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+/=".indexOf(e);
                for (var u = 0, a = n.length; u < a; u++)
                    i += "%" + ("00" + n.charCodeAt(u).toString(16)).slice(-2);
                return decodeURIComponent(i)
            },
            t = arguments,
            rr.avoPoV = !0
        }
        var o = r + e[0],
            c = t[o];
        return c ? i = c : (i = rr.nyBjcA(i), t[o] = i),
        i
    }, rr(t, r)
}
```

#### Step 3: The Ft (which becomes Rt) object definition

```javascript
// This appears in the obfuscated code (around line ~9700+ in your file)
var Ft = {};
Ft[G(548)] = Pt,        // G(548) = "appName"
Ft[G(519)] = Nt,        // G(519) = "appKey"  
Ft[G(430) + "s"] = Ht,  // G(430) = "endpoint" → "endpoints"
Ft[G(421) + G(411)] = Wt // G(421) = "shplEndp", G(411) = "oints" → "shplEndpoints"
```

#### Step 4: The Nt object (which is appKey) definition
```javascript
var Nt = {};
// For 1.0 version
Nt[G(532)] = {};  // G(532) = "1.0"
Nt[G(532)][G(441)] = G(569) + G(481) + G(363) + G(565);  // G(441) = "WEB" or "cn"?
// Actual: "sh87bd15" + "12hsb03c" + "b405803a" + "307dbe32"

Nt[G(532)][G(520)] = G(434) + G(490) + G(523) + G(498);  // G(520) = "sgp"
// Actual: "sg63c0a0" + "34gsf3f3" + "7ef9e8a2" + "c8a0bc51"

// For 2.0 version
Nt[G(487)] = {};  // G(487) = "2.0"
Nt[G(487)][G(441)] = G(345) + G(526) + G(424) + G(322);  // "ab034ec0" + "643f9139" + "9eb33e06" + "2dc7fae1"
Nt[G(487)][G(520)] = G(337) + G(589) + G(588) + G(455);  // "3795d282" + "42a11619" + "bc25f786" + "f84e53d4"

// For 3.0 version (same as 2.0)
Nt[G(423)] = {};  // G(423) = "3.0"
Nt[G(423)][G(441)] = G(345) + G(526) + G(424) + G(322);
Nt[G(423)][G(520)] = G(337) + G(589) + G(588) + G(455);
```

---

## Function rr index deobfuscated literals json log [AliyunCaptcha.js]

```json
{
  "312": "1662178hogNeO",
  "313": "theast-d",
  "314": "1NfQv95E",
  "315": "SG_WEB_P",
  "316": "dZJMm",
  "317": "re-b.ali",
  "318": "HaTtw",
  "319": "0.0.0/fe",
  "320": "theast.a",
  "321": "9NhnQQ+L",
  "322": "2dc7fae1",
  "323": "ual.aliy",
  "324": "z2k=",
  "325": "PiVwNmNO",
  "326": "a.device",
  "327": "ptchaV3",
  "328": "-fronten",
  "329": "SG_WEB",
  "330": "INIT",
  "331": "hai.devi",
  "332": "theast-b",
  "333": "Yex3WaHd",
  "334": "AIN_FAIL",
  "335": "W8YrgOBc",
  "336": "split",
  "337": "3795d282",
  "338": "liyuncs.",
  "339": "DNsKKPRH",
  "340": "S+qKUbJi",
  "341": "2B0=",
  "342": "/captcha",
  "343": "ap-south",
  "344": "aptcha-o",
  "345": "ab034ec0",
  "346": "pbG/bc9x",
  "347": ".saf.ali",
  "348": "FLAG",
  "349": "Log3",
  "350": "2VS3zA==",
  "351": "7JLsB18M",
  "352": "success",
  "353": "theast-p",
  "354": "cha",
  "355": "ONUTG",
  "356": "fail",
  "357": "l.aliyun",
  "358": "aliyuncs",
  "359": "EdGaRtkc",
  "360": "gp.aliyu",
  "361": "web-pre.",
  "362": "heast-1.",
  "363": "b405803a",
  "364": "evice.sa",
  "365": "3xLOkphC",
  "366": ".com",
  "367": "2415840gMPfPM",
  "368": "hanghai.",
  "369": "captcha-",
  "370": "INITV3",
  "371": "VERIFY",
  "372": "UploadLo",
  "373": "com",
  "374": "8KmHIQsc",
  "375": "M0v7u45+",
  "376": "re.aliyu",
  "377": "pen.aliy",
  "378": "open-b.a",
  "379": "DEVICE_U",
  "380": "af.aliyu",
  "381": "-1.devic",
  "382": "COMBAT_U",
  "383": "east-1.d",
  "384": "outheast",
  "385": "web-b.al",
  "386": "ck.ap-so",
  "387": "ncs.com/",
  "388": "d/FeiLin",
  "389": "sg50c495",
  "390": "open-pre",
  "391": "iyuncs.c",
  "392": "ptchaV2",
  "393": "pre-ap-s",
  "394": "Log1",
  "395": "open.ali",
  "396": "VYdEGpwb",
  "397": "-b.aliyu",
  "398": "W.10001.",
  "399": "ss0=",
  "400": "F0tJ3ds4",
  "401": "ncs.com",
  "402": "95bc895c",
  "403": "vkI7QQBX",
  "404": "d/aliyun",
  "405": ".aliyunc",
  "406": "ed381ac9",
  "407": "RES",
  "408": "Vo71uzWk",
  "409": "X1y5Vstb",
  "410": "5+LZJA7u",
  "411": "oints",
  "412": "NIT_FAIL",
  "413": "FAIL",
  "414": "DEVICE_M",
  "415": "SzaFtFlN",
  "416": "W2022020",
  "417": "cloudaut",
  "418": "c0ad7983",
  "419": "chaV2",
  "420": "tBwmiXXw",
  "421": "shplEndp",
  "422": "Cbo=",
  "423": "3.0",
  "424": "9eb33e06",
  "425": "REFRESH_",
  "426": "wZaUWGBt",
  "427": "9ebbf3d0",
  "428": "east-1-g",
  "429": "r4qktDu7",
  "430": "endpoint",
  "431": "5169504KEaUDm",
  "432": "open-sou",
  "433": "/1scBbcb",
  "434": "sg63c0a0",
  "435": "/main.cs",
  "436": "dev.o.al",
  "437": "C/c3QKzT",
  "438": "aptcha-p",
  "439": "VzY=",
  "440": "sh3c47a8",
  "441": "WEB",
  "442": "yuncs.co",
  "443": "com/",
  "444": ".js",
  "445": "853792CgZYEZ",
  "446": "ual-b.al",
  "447": "9U6Kg6lj",
  "448": "chaV3",
  "449": "Gas=",
  "450": "SUCCESS",
  "451": "aptcha-s",
  "452": "ddhs0305",
  "453": "INITV2",
  "454": "807878SVsivw",
  "455": "f84e53d4",
  "456": "PLOAD",
  "457": "xLLw/t15",
  "458": "anghai.a",
  "459": "NiaNh",
  "460": "518MUXIlt",
  "461": "e.saf.al",
  "462": "static-c",
  "463": "83f5e541",
  "464": "k+1RW0cz",
  "465": "75184fAvUYY",
  "466": "InitCapt",
  "467": "open-ga-",
  "468": "NLAoqT6K",
  "469": "LTXPs",
  "470": "aptcha.a",
  "471": "/bfozcSz",
  "472": "fOTuFph6",
  "473": "h-device",
  "474": "utheast-",
  "475": "device.c",
  "476": "SECRET",
  "477": "southeas",
  "478": "LOG",
  "479": "XiMPo",
  "480": "KFr07pEb",
  "481": "12hsb03c",
  "482": "l-b.aliy",
  "483": "open-dua",
  "484": "saf-aliy",
  "485": "www.aliy",
  "486": "WEB_PREI",
  "487": "2.0",
  "488": "1910DXe6",
  "489": "t-1.aliy",
  "490": "34gsf3f3",
  "491": "mosql",
  "492": "n.js",
  "493": "un.com/",
  "494": "bBLLx",
  "495": "cJS/",
  "496": "now",
  "497": "_l_KPLiP",
  "498": "c8a0bc51",
  "499": "ilin",
  "500": "PREID",
  "501": "dev.g.al",
  "502": "-pre.ali",
  "503": "FqJB6iRN",
  "504": "o.alicdn",
  "505": "_each",
  "506": "S_FAIL",
  "507": "GbSWQ",
  "508": "un-com",
  "509": "VBg=",
  "510": "REQ",
  "511": "f.aliyun",
  "512": "saf-capt",
  "513": "REID",
  "514": "202130TawWqM",
  "515": "taPHkC+T",
  "516": "3iDi2RAb",
  "517": "pre-cn-s",
  "518": "cn_dual",
  "519": "appKey",
  "520": "sgp",
  "521": "grM=",
  "522": "BdGKm",
  "523": "7ef9e8a2",
  "524": "cha-waf",
  "525": "https://",
  "526": "643f9139",
  "527": "g3E=",
  "528": "8fgs168b",
  "529": "-dualsta",
  "530": "nA7GX3d6",
  "531": "sgp_dual",
  "532": "1.0",
  "533": "B5zqghyJ",
  "534": "d35db7e3",
  "535": "g.alicdn",
  "536": "Plkdn",
  "537": "RrKCkAux",
  "538": "2023-03-",
  "539": "LvIoA/F6",
  "540": "pro-open",
  "541": "OhpiO",
  "542": "n9jH0yAC",
  "543": "sihjy1wJ",
  "544": "web.aliy",
  "545": "bC6YUaXi",
  "546": "kjgCtS6e",
  "547": "prototyp",
  "548": "appName",
  "549": "03oLbQXW",
  "550": "lFi23FmD",
  "551": "-pre-b.a",
  "552": "1.aliyun",
  "553": "+fR9tYzl",
  "554": "-pre.ap-",
  "555": ".ap-sout",
  "556": "uncs.com",
  "557": "LIMIT_FL",
  "558": "DEVICE_I",
  "559": "VERIFYV3",
  "560": "T68xcVuO",
  "561": "MfA=",
  "562": "INIT_FAI",
  "563": "OTHER",
  "564": "UPLOAD",
  "565": "307dbe32",
  "566": "PIC_FAIL",
  "567": "_c_WBKFR",
  "568": "cs.com",
  "569": "sh87bd15",
  "570": "FP/fp.mi",
  "571": "9fKx+ymx",
  "572": "cn-shang",
  "573": "Aoxz0b7v",
  "574": "d/dynami",
  "575": "gp-pre.a",
  "576": "device.s",
  "577": "umd7+KnO",
  "578": "s.com",
  "579": "http://",
  "580": "RDfGiyiM",
  "581": "icdn.com",
  "582": "VerifyCa",
  "583": "8ZpvzGBX",
  "584": "9g8a0A==",
  "585": "LxErT1sG",
  "586": "01083105",
  "587": "ck.cn-sh",
  "588": "bc25f786",
  "589": "42a11619",
  "590": "2020-10-",
  "591": "Log2",
  "592": "ce.saf.a",
  "593": "b/ksJCrC",
  "594": "gm9PhHvS",
  "595": "DYNAMICJ"
}
```
---

## Function Re reverse engineered DeviceData creation [AliyunCaptcha.js]

Function Re generates DeviceData which is sent in function Pe

#### Function Pe Definition
```javascript
function Pe() {
            return Pe = k()(P().mark(function t(r, e, n, i) {
                var o, c, u, a, s, f, l, p, v, h, d, y, g, m, x;
                return P().wrap(function(t) {
                    for (; ; )
                        switch (t.prev = t.next) {
                        case 0:
                            return De._extend({
                                initBeginTime: Date.now(),
                                logUploaded: !1,
                                logInfo: {}
                            }),
                            xr("sId", r.SceneId),
                            o = e.https,
                            c = e.initPath,
                            u = e.isDev,
                            a = e.verifyType,
                            s = o,
                            f = Qe(e),
                            l = tn(r, e),
                            p = l.action,
                            xr("pfx", v = l._prefix),
                            f = N()(f).call(f, function(t) {
                                return v + "." + t
                            }),
                            h = N()(f).call(f, function(t) {
                                return lr(s, t, c)
                            }),
                            De._extend({
                                urls: h
                            }),
                            d = i.deviceConfig,
                            y = i.deviceCallback,
                            "1.0" === a ? (delete r.DeviceToken,
                            Ie = new Jt) : e.userId && e.userUserId && (De._extend({
                                userId: void 0,
                                userUserId: void 0
                            }),
                            Ie = new Jt),
                            $e(d.endpoints, d.appName),
                            g = Re(d, Ie, ke),
                            e.isFromTraceless || void 0 !== Ie.DeviceConfig || (r.DeviceData = g),
                            t.next = 1,
                            Ne(p, r, h, e, Ee);
                        case 1:
                            !(m = t.sent).Success || m.LimitFlow || m.LimitedFlowToken ? (m.LimitedFlowToken ? m.CertifyId = m.LimitedFlowToken : m.CertifyId || (m.CertifyId = dr().substring(0, 5)),
                            xr("cId", m.CertifyId),
                            n(Ee.ACTION_STATE.FAIL, m)) : (e._extend({
                                log: nn
                            }),
                            xr("cId", m.CertifyId),
                            !e.isFromTraceless && De._extend({
                                initialRequestTime: Date.now(),
                                overTime: !1
                            }),
                            m.DeviceConfig && void 0 === Ie.DeviceConfig && Ie._extend({
                                DeviceConfig: m.DeviceConfig
                            }),
                            rn(m.DeviceConfig, y, u, "captcha"),
                            x = Se(m, e),
                            n(Ee.ACTION_STATE.SUCCESS, x));
                        case 2:
                        case "end":
                            return t.stop()
                        }
                }, t)
            })),
            Pe.apply(this, arguments)
        }
```

#### Function Re Definition
```javascript
        function Re(t, r, e) {
            r._extend(Be({}, t));
            var n = t.appKey || r.APP_KEY
              , i = t.appName || r.APP_NAME
              , o = me(r.ACCESS_SEC, r.secretKey) || e.WEB_AES_FLAG_SECRET_KEY
              , c = r.PLATFORM + "#" + i + "#" + (r.sceneId || "") + "#captcha-normal#" + De.prefix + "#" + De.region;
            return c = ye(o, c),
            be([n, e.DEVICE_TYPE.WEB, c, r.APP_VERSION, "CLOUD", ""])
        }
```
Notes 
- 1. `var n` in context of `var n = t.appKey || r.APP_KEY can have multiple values since it likely gets it values from Rt.appKey which have different values based on region and version. This is same for `var i` im not sure but it seems obviuous 
- 2. in context of `var o` the key used is e.WEB_AES_FLAG_SECRET_KEY since r.secretKey is undefined and we should stick to it since it works 



### Function Re reconstruction in python

`generate_device_data.py`

```python

"""
AliyunCaptcha - DeviceData Generator & Inspector
Replicates: Re() -> ye() -> be() -> ye() chain from AliyunCaptcha.js

Static values extracted from AliyunCaptcha.js via browser deobfuscation.
Dynamic values (sceneId, prefix, region) are per-captcha-instance.

Public API:
  generate_device_data(scene_id, prefix, region)  → full DeviceData (ye + be + ye)
  generate_only_ye(scene_id, prefix, region)       → first ye() pass only, no be()
  decrypt_ye(ciphertext_b64, key_str)              → decrypt one ye() layer
  decrypt_full(device_data)                        → decrypt both layers, return dict
"""

from Crypto.Cipher import AES
from Crypto.Util.Padding import pad, unpad
import base64


# ─────────────────────────────────────────────
#  STATIC CONSTANTS (never change, baked into JS)
# ─────────────────────────────────────────────

APP_KEY     = "3795d28242a11619bc25f786f84e53d4"    # Note: Default is G(345)+G(526)+G(424)+G(322) but derived from t.appKey or r.APP_KEY which likely points to Rt.appKey and is different per region and version
PLATFORM    = "W.10001.c"                           # G(398)+"c"
APP_NAME    = "saf-captcha"                         # appName from deviceConfig
APP_VERSION = "W20220202"                           # G(416)+"2"
DEVICE_TYPE = "W"                                   # Vt["WEB"]

# Key used in first AES pass (ye inside Re)
# = ke.WEB_AES_FLAG_SECRET_KEY = me(ACCESS_SEC, WEB_AES_SECRET_KEY["FLAG"])
KEY_O = "c175a358550d02e2"

# Key used in second AES pass (be -> ye)
# = he = me(ACCESS_SEC, WEB_AES_SECRET_KEY["REQ"])
KEY_HE = "45f8ac1e1de14397"

# Fixed IV — from pe.iv.words in the browser
# words: [808530483, 875902519, 943276354, 1128547654]
# hex:   d35db7e3 9ebbf3d0 01083105 43434346  → but use exact bytes from words
_IV_WORDS = [808530483, 875902519, 943276354, 1128547654]

def _build_iv(words: list) -> bytes:
    """Convert CryptoJS WordArray words to 16-byte IV."""
    out = bytearray()
    for w in words:
        out += w.to_bytes(4, 'big')
    return bytes(out[:16])

IV = _build_iv(_IV_WORDS)


# ─────────────────────────────────────────────
#  CORE CRYPTO — replicates ye(key, plaintext)
# ─────────────────────────────────────────────

def ye(key_str: str, plaintext_str: str) -> str:
    """
    Replicates ye(t, r) from AliyunCaptcha.js:

      case "2": S = ue.parse(t)        → UTF-8 encode key → WordArray
      case "4": C = r                  → plaintext string
      case "0": b = ce.encrypt(C, S, pe) → AES-CBC encrypt with fixed IV + PKCS7
      case "5": return b.toString()    → CipherParams.toString() → Base64 ciphertext

    Key must be exactly 16 bytes (AES-128).
    Returns: Base64 string of raw AES-CBC ciphertext (no OpenSSL "Salted__" prefix,
             because key is passed as WordArray → SerializableCipher path).
    """
    if key_str is None or len(key_str) < 16:
        return None
    if plaintext_str is None or len(plaintext_str) == 0:
        return None

    key = key_str.encode('utf-8')[:16]          # AES-128: first 16 bytes
    plaintext = plaintext_str.encode('utf-8')
    cipher = AES.new(key, AES.MODE_CBC, IV)
    encrypted = cipher.encrypt(pad(plaintext, AES.block_size))
    return base64.b64encode(encrypted).decode('utf-8')


def decrypt_ye(ciphertext_b64: str, key_str: str = KEY_O) -> str:
    """
    Decrypt one ye() layer.

    Args:
        ciphertext_b64: Base64 string produced by ye()
        key_str:        16-char AES key. Defaults to KEY_O (first-pass key).
                        Pass KEY_HE to decrypt be() outer layer.
    Returns:
        Plaintext string.
    """
    key = key_str.encode('utf-8')[:16]
    data = base64.b64decode(ciphertext_b64)
    cipher = AES.new(key, AES.MODE_CBC, IV)
    return unpad(cipher.decrypt(data), AES.block_size).decode('utf-8')


def decrypt_full(device_data: str) -> dict:
    """
    Fully decrypt a DeviceData string — reverses both encryption layers.

    Layer 1 (outer): be() → ye(KEY_HE, joined_array)
    Layer 2 (inner): ye(KEY_O, plaintext_c)

    Args:
        device_data: Base64 DeviceData string from generate_device_data()

    Returns dict with keys:
        joined_array  → the "#"-joined array before outer encryption
        array_parts   → list of the 6 array elements
        plaintext_c   → the inner plaintext before first encryption
        app_key       → extracted appKey
        device_type   → "W"
        encrypted_c   → the inner ciphertext sitting in position [2]
        app_version   → extracted APP_VERSION
        scene_id      → extracted sceneId
        prefix        → extracted prefix
        region        → extracted region
        app_name      → extracted appName
    """
    # Outer layer: decrypt be() → ye(KEY_HE, ...)
    joined_array = decrypt_ye(device_data, KEY_HE)
    array_parts = joined_array.split("#")
    # array_parts = [appKey, "W", encrypted_c, APP_VERSION, "CLOUD", ""]

    encrypted_c = array_parts[2]

    # Inner layer: decrypt ye(KEY_O, plaintext_c)
    plaintext_c = decrypt_ye(encrypted_c, KEY_O)
    # plaintext_c = "PLATFORM#appName#sceneId#captcha-normal#prefix#region"
    c_parts = plaintext_c.split("#")

    return {
        "joined_array" : joined_array,
        "array_parts"  : array_parts,
        "plaintext_c"  : plaintext_c,
        "app_key"      : array_parts[0],
        "device_type"  : array_parts[1],
        "encrypted_c"  : encrypted_c,
        "app_version"  : array_parts[3],
        "platform"     : c_parts[0],
        "app_name"     : c_parts[1],
        "scene_id"     : c_parts[2],
        "prefix"       : c_parts[4],
        "region"       : c_parts[5],
    }


# ─────────────────────────────────────────────
#  be() — replicates be(array) from JS
# ─────────────────────────────────────────────

def be(array: list) -> str:
    """
    Replicates be(t) from AliyunCaptcha.js:

      return ye(he, t.join("#"))

    Joins the array with "#" and encrypts with KEY_HE.
    """
    joined = "#".join(array)
    return ye(KEY_HE, joined)


# ─────────────────────────────────────────────
#  Re() — the main DeviceData generator
# ─────────────────────────────────────────────

def generate_device_data(
    scene_id: str,
    prefix: str = "",
    region: str = "cn",
    app_name: str = APP_NAME,
    app_key: str = APP_KEY,
) -> str:
    """
    Replicates Re(t, r, e) from AliyunCaptcha.js.

    Re() does:
      n = appKey
      i = appName
      o = WEB_AES_FLAG_SECRET_KEY          (me() fails → fallback to static key)
      c = PLATFORM + "#" + i + "#" + sceneId + "#captcha-normal#" + prefix + "#" + region
      c = ye(o, c)                          ← first AES pass
      return be([n, "W", c, APP_VERSION, "CLOUD", ""])  ← join + second AES pass

    Args:
        scene_id: Your captcha SceneId (e.g. "didk33e0")
        prefix:   er.prefix from captcha config (usually empty string "")
        region:   er.region from captcha config ("cn" or "sg")
        app_name: appName from deviceConfig (default "saf-captcha")
        app_key:  appKey from deviceConfig (default static value)

    Returns:
        DeviceData string (Base64) — attach as r.DeviceData in the init payload.
    """

    # Build the plaintext that gets first-pass encrypted
    # r.PLATFORM + "#" + i + "#" + (r.sceneId||"") + "#captcha-normal#" + De.prefix + "#" + De.region
    plaintext_c = f"{PLATFORM}#{app_name}#{scene_id}#captcha-normal#{prefix}#{region}"

    # First AES pass: ye(o, c)
    encrypted_c = ye(KEY_O, plaintext_c)

    # be([appKey, "W", encrypted_c, APP_VERSION, "CLOUD", ""])
    # → joins with "#" → ye(he, joined)
    device_data = be([app_key, DEVICE_TYPE, encrypted_c, APP_VERSION, "CLOUD", ""])

    return device_data


# ─────────────────────────────────────────────
#  generate_only_ye() — first pass only, no be()
# ─────────────────────────────────────────────

def generate_only_ye(
    scene_id: str,
    prefix: str = "",
    region: str = "cn",
    app_name: str = APP_NAME,
) -> str:
    """
    Runs only the first ye() pass from Re() — skips be().

    Produces the encrypted_c value that sits at array_parts[2] inside
    a full DeviceData. Useful if you need to inspect or inject just the
    inner ciphertext, or call be() manually yourself.

    Returns: Base64 string of ye(KEY_O, plaintext_c)
    """
    plaintext_c = f"{PLATFORM}#{app_name}#{scene_id}#captcha-normal#{prefix}#{region}"
    return ye(KEY_O, plaintext_c)


# ─────────────────────────────────────────────
#  VERIFICATION (internal helper)
# ─────────────────────────────────────────────

def _verify(device_data: str, scene_id: str, prefix: str, region: str):
    """Decrypt DeviceData back and print each layer."""
    print("\n── Verification ──")
    result = decrypt_full(device_data)
    print(f"Layer 2 (joined array): {result['joined_array']}")
    print(f"Layer 1 (plaintext c):  {result['plaintext_c']}")
    expected = f"{PLATFORM}#{APP_NAME}#{scene_id}#captcha-normal#{prefix}#{region}"
    match = "✅ MATCH" if result['plaintext_c'] == expected else "❌ MISMATCH"
    print(f"Result: {match}")


# ─────────────────────────────────────────────
#  MAIN
# ─────────────────────────────────────────────

if __name__ == "__main__":

    scene_id = "didk33e0"
    prefix   = ""
    region   = "cn"

    # 1. Full DeviceData (ye + be + ye)
    print("── 1. generate_device_data() ──")
    device_data = generate_device_data(scene_id, prefix, region)
    print(f"DeviceData: {device_data}")
    _verify(device_data, scene_id, prefix, region)

    # 2. First ye() pass only — no be()
    print("\n── 2. generate_only_ye() ──")
    only_ye = generate_only_ye(scene_id, prefix, region)
    print(f"encrypted_c: {only_ye}")
    print(f"Decrypted back: {decrypt_ye(only_ye, KEY_O)}")

    # 3. decrypt_ye() — single layer
    print("\n── 3. decrypt_ye() ──")
    outer_plaintext = decrypt_ye(device_data, KEY_HE)
    print(f"Outer layer decrypted: {outer_plaintext}")

    # 4. decrypt_full() — both layers, structured
    print("\n── 4. decrypt_full() ──")
    parsed = decrypt_full(device_data)
    for k, v in parsed.items():
        print(f"  {k:15s}: {v}")
```

## Signature and Nonce generation [AliyunCaptcha.js]

Every request is signed with a HMAC-SHA1 Signature and Signature Nonce

Shown in function je

```javascript
function je() {
            return je = k()(P().mark(function t(r, e, n, i, o) {
                var c, u;
                return P().wrap(function(t) {
                    for (; ; )
                        switch (t.prev = t.next) {
                        case 0:
                            return (c = {}).AccessKeyId = o.KEY_ID,
                            c.SignatureMethod = "HMAC-SHA1",
                            c.SignatureVersion = "1.0",
                            c.Format = "JSON",
                            c.Timestamp = hr(),
                            c.Version = lt,
                            c.Action = r,
                            ir(e) || (c = or(c, e)),
                            u = function() {
                                var t = k()(P().mark(function t(r) {
                                    var e, a, s, f, l, p, v, h;
                                    return P().wrap(function(t) {
                                        for (; ; )
                                            switch (t.prev = t.next) {
                                            case 0:
                                                return c.SignatureNonce = dr(),
                                                a = Ce(c, o.KEY_SECRET),
                                                c.Signature = a,
                                                s = Date.now(),
                                                t.next = 1,
                                                Ke(n[r], c, i);
                                            case 1:
                                                if (f = t.sent,
                                                l = Date.now(),
                                                p = f.Code,
                                                v = f.Success,
                                                h = Gr()(e = n[r]).call(e, "-b") ? "bInit" : "mInit",
                                                !("Success" === p && v || r >= n.length - 1)) {
                                                    t.next = 2;
                                                    break
                                                }
                                                return "Success" === p && v ? (xr(h, {
                                                    t: l,
                                                    s: !0,
                                                    msg: "INIT_SUCCESS",
                                                    rt: l - s
                                                }),
                                                Ge(r)) : xr(h, {
                                                    t: l,
                                                    s: !1,
                                                    msg: f.err,
                                                    rt: l - s
                                                }),
                                                t.abrupt("return", f);
                                            case 2:
                                                if (xr(h, {
                                                    t: l,
                                                    s: !1,
                                                    msg: f.err || f.Message,
                                                    rt: l - s
                                                }),
                                                !("403" === p && f.LimitedFlow || "ThrottlingByStrategy" === p)) {
                                                    t.next = 3;
                                                    break
                                                }
                                                return t.abrupt("return", f);
                                            case 3:
                                                return t.next = 4,
                                                u(r + 1);
                                            case 4:
                                                return t.abrupt("return", t.sent);
                                            case 5:
                                            case "end":
                                                return t.stop()
                                            }
                                    }, t)
                                }));
                                return function(r) {
                                    return t.apply(this, arguments)
                                }
                            }(),
                            t.next = 1,
                            u(0);
                        case 1:
                            return t.abrupt("return", t.sent);
                        case 2:
                        case "end":
                            return t.stop()
                        }
                }, t)
            })),
            je.apply(this, arguments)
        }
```

Notes:
- 1. `o.KEY_SECRET` points to `var Ee.KEY_SECRET` which is  "YSKfst7GaVkXwZYvVihJsKF9r89koz"

### Reverse Engineered Signature Generation Example in Python
`Signature.py`
```python
import hmac
import hashlib
import base64
import urllib.parse
import uuid
from datetime import datetime, timezone
import requests

def generate_nonce():
    """Generate UUID v4 as SignatureNonce"""
    return str(uuid.uuid4())

def get_timestamp():
    """Get UTC timestamp in ISO format without microseconds"""
    return datetime.now(timezone.utc).strftime('%Y-%m-%dT%H:%M:%SZ')

def generate_signature(params, secret_key):
    """
    Generate HMAC-SHA1 signature for Aliyun CAPTCHA API
    """
    # Sort parameters alphabetically
    sorted_params = sorted(params.items())

    # Build canonicalized query string
    canonicalized_query = '&'.join([
        f"{urllib.parse.quote(k, safe='')}={urllib.parse.quote(str(v), safe='')}"
        for k, v in sorted_params
    ])

    # String to sign: POST&%2F& + urlencode(canonicalized_query)
    string_to_sign = f"POST&%2F&{urllib.parse.quote(canonicalized_query, safe='')}"

    # Calculate HMAC-SHA1 (key = secret_key + "&")
    signing_key = (secret_key + "&").encode('utf-8')
    signature = hmac.new(signing_key, string_to_sign.encode('utf-8'), hashlib.sha1)
    signature_b64 = base64.b64encode(signature.digest()).decode('utf-8')

    return signature_b64

def make_captcha_request(access_key_id, secret_key, scene_id, device_token):
    """
    Make a valid CAPTCHA initialization request
    """
    # Generate new nonce and timestamp for each request
    nonce = generate_nonce()
    timestamp = get_timestamp()

    # Build parameters (without Signature)
    params = {
        'AccessKeyId': access_key_id,
        'Action': 'InitCaptchaV3',
        'Format': 'JSON',
        'Language': 'en',
        'Mode': 'popup',
        'SceneId': scene_id,
        'SignatureMethod': 'HMAC-SHA1',
        'SignatureNonce': nonce,
        'SignatureVersion': '1.0',
        'Timestamp': timestamp,
        'UpLang': 'true',
        'Version': '2023-03-05',
        'DeviceToken': device_token,
    }

    # Generate signature using all parameters
    signature = generate_signature(params, secret_key)
    params['Signature'] = signature

    # Build form data
    body = '&'.join([f"{k}={urllib.parse.quote(str(v), safe='')}" for k, v in params.items()])

    # Make request
    response = requests.post(
        'https://no8xfe.captcha-open-southeast.aliyuncs.com/',
        headers={
            'Accept': '*/*',
            'Accept-Language': 'en-US,en;q=0.9',
            'Cache-Control': 'no-cache',
            'Content-Type': 'application/x-www-form-urlencoded; charset=UTF-8',
            'Pragma': 'no-cache',
        },
        data=body,
        timeout=30
    )

    return response

# Usage
access_key = "LTAI5tSEBwYMwVKAQGpxmvTd"  # Your AccessKeyId
secret = "YSKfst7GaVkXwZYvVihJsKF9r89koz"  # You need the actual secret key
scene = "didk33e0"
device_token = "U0dfV0VCIzM3OTVkMjgyNDJhMTE2MTliYzI1Zjc4NmY4NGU1M2Q0LWgtMTc3ODc1OTg5ODA4Ny00MWYzYjlmNzAyZmI0MDBiOTkyMDQ1ODYwMWNkOTE5MyNrbDVhUnF6bjdGc0hKYnhLYXh0ZC82WWFwMHFUY1U1SUllRmZpUERMWCt3cVc5NWsrV1JQTHYxNGlRUUtWNUQ3Ui9nTXJDVy9zRzhsWlBEWW1vWFpGYTN4Y2RSMW0vc1pxNS95NmFnTzFKWHM0MkJ2VjdkVnBuZFdHTEtTQjFxL1RGRWNROFl3Qy9OM0dqTTJTY1F0L1lHbkl6ODVQTmRTdnVPRmpJMXZTRVRLWi9POHZsSjRYUzBBcy9UREJGcWlJc0l2eUI2blpYU0RDcmZ1YlUxb3lGNUQvRy9USUFFQThIWWQyUTMrR3YyT0g3U0pnN1RYbWQwSDRqSnNHTVZQdU50RU5oWDdHQXFIVGk3VEJwLzV1S0gvMlpvN2V3eVVjL0c1L0s5am4vUUh5cnlZZHRhdERzVlRQYU1JTG51Q3VrVVlrM1JpRE5VeWtkOHE4eWF5N2d1UkJvQXpQaUZ4aVdEbUtuVERwZUIzbnlXZ0ZBL21ZYWJmZ1dHU09ic2NQRENKY2xGTGRQK0dUZ3p1UWdXZG55QlpVNEVnaTVuMkQzUlNvYTlab244ckZXNFA5Sm1HZ0crQ1IwZVpWeVlYSDFIV2liYTRVVDVxZjh3SkdKclRZODZUbjJoWXNsL3dGT3ZXeTJPek8xZlVTNHd1d3V3REVra0gySTdpY2pldm1lM3doajdNNlV5OXpNUkFJa0JCTjRhUjdqWndyTnlXcFowK0JINjFWN0J4bVlXQVdqemtCcWpVSDREbEZZd3p6aVlUT1FUVlZFTDJ4dmEzV2NkK0dlOTlNVHNVOG9sRVgvWDQ0QzJBS3hpYTBDSzV0ZkZlQzB2Mkh1NUpLSW01Yk9Vd20wVjdLVGovczNHSWVNb25Vb2sxV20yMmVPd0J3aFJDbVhxWVlUZ3VqQjlPN0tPSmlzc1NyMHBOaXBHbXVVTGViYldCSXlrcnlOOUN1VTBsOWsraWsvZ1V5U1ZyRzNveStsaWk1ZVRybVZMdzFDOThGTjh6K09LOEFaU2ZrZDBudVk1UWZIdDFYSk5OYWxrRm1LOWk4ZUFIeXJzaDUydUJFek4rQUVBc1ZFQjUrdEJ6bllmOHQvVllUK05sUW9rcDkzY29jSWhVY0tjYUtMRHJ5Zz09IzY2NSNjMzEzMmRkYmZlMjVjNzZiMzkyNDYwNjgxN2ExZGYxMQ=="

response = make_captcha_request(access_key, secret, scene, device_token)
print(response.json())
```

### Reverse Engineered Signature Generation Psuedocode in Javascript
`Signature.js`
```javascript
const crypto = require('crypto');

function generateAliyunSignature(params, secretKey) {
    // 1. Remove Signature if exists
    delete params.Signature;

    // 2. Sort keys alphabetically
    const sortedKeys = Object.keys(params).sort();

    // 3. Build canonicalized query string
    const canonicalized = sortedKeys.map(key => {
        return encodeURIComponent(key) + "=" + encodeURIComponent(params[key]);
    }).join("&");

    // 4. Build string to sign
    const stringToSign = "POST&%2F&" + encodeURIComponent(canonicalized);

    // 5. Generate HMAC-SHA1 signature (IMPORTANT: add trailing &)
    const signatureKey = secretKey + "&";
    const hmac = crypto.createHmac('sha1', signatureKey);
    hmac.update(stringToSign);

    return hmac.digest('base64');                                                                                                                                         }                                                                                                                                                                                                                                                                                                                                                   // Example usage:

const params = {                                                                                                                                                              AccessKeyId: "LTAI5tSEBwYMwVKAQGpxmvTd",                                                                                                                                  Action: "InitCaptchaV3",                                                                                                                                                  Format: "JSON",
    SignatureMethod: "HMAC-SHA1",
    SignatureVersion: "1.0",
    Timestamp: new Date().toISOString().replace(/\.\d+Z$/, 'Z'),
    Version: "2023-03-05",
    SceneId: "didk33e0",
    Language: "en",
    Mode: "popup",
    UpLang: "true",
    DeviceData: "your_device_data_here"
};

// Generate nonce (UUID v4)
const signatureNonce = crypto.randomUUID();
params.SignatureNonce = signatureNonce;

// Generate signature
const signature = generateAliyunSignature(params, "YSKfst7GaVkXwZYvVihJsKF9r89koz");
params.Signature = signature;

console.log("Generated params with signature:", params);
```

---

## Traceless CaptchaFlow
### 1. AliyunCaptccha.js Loads

### 2. Loads fielin.*.js
feiling.js generates DeviceToken and exposes a call
`window.um?.getToken`
`window.z_um?.getToken`
DeviceToken is calulated by fielin js and used in client server requests

### 3. InitCaptchaV3
fetch request
```javascript
fetch("https://no8xfe.captcha-open-southeast.aliyuncs.com/", {
  "headers": {
    "accept": "*/*",
  },
  "body": "AccessKeyId=LTAI5tSEBwYMwVKAQGpxmvTd&SignatureMethod=HMAC-SHA1&SignatureVersion=1.0&Format=JSON&Timestamp=2026-05-25T02%3A42%3A52Z&Version=2023-03-05&Action=InitCaptchaV3&SceneId=didk33e0&Language=en&Mode=popup&UpLang=true&DeviceData=TEQYvgJq1LrMqFaBybfIzPxz2ygFyAct7X%2Fw%2BLacfXWd9rGSwE%2Fx6ZCONucD1fehS2Qpig6tUVsFK111d9wIk5pWp6rwYjzFCRgL7pNp8bzGsvOSdUXgQTopQm90YPSdCiRAlgENdODLvY7P8jrfO9eC15tPCPwLxcRIrcspVvQYqVfk9%2FyFeIlePKmTRjkM&SignatureNonce=9aab2319-86bf-4019-839e-82225f796296&Signature=UDiQcb82lN5kZT0ZPNOmw0hT8sw%3D",
  "method": "POST",
  "mode": "cors",
  "credentials": "omit"
});
```
DeviceData is sent on first request it is determined by if the previous response inlucded CaptchaTypre Traceless or not since before there was no request it is fired 
On second requst DeviceToken is sent if previous inclueded TRACELESS captcha type
```json
{
    "CertifyId": "gl62Aqi6e1",
    "Message": "success",
    "RequestId": "65A5902F-1573-4BA3-BB81-898D35553B23",
    "Code": "Success",
    "LimitFlow": false,
    "Success": true,
    "StaticPath": "3.25.0/pe.050.9453665072ad79a7",
    "CaptchaType": "TRACELESS",
    "DeviceConfig": "ckj8rA4/WEfd9fOcvSSuIJ6e0IDZR3X79geskaQ1tN6EJBdqOGupxhCCoHnlEZsxgowJxh2znY9HEY2Ruj07n/YhavNvqMCzc3ppJQkefLtCHXx7PUpO8NhrB+OV2e3AnSJ1oRFgSX7sWxZoUSoLNYNh8QbZOYMHO+ShVv5BADyF02OloVszuUsNbbDeDYId2moM/OpnagWMpLGNJJiVFOfHevUTtYQTsNN/X9iCN3RayeHFtcgxBjhEZcfHPzNVYn6bCGyZ2mCVE8J1swWnPFKiSJuiWiAxX4Qsn5pOtWs="
}
```

```javascript
Object.keys(AliyunCaptcha.prototype)
[
    "config",
    "deviceConfig",
    "startPOWCalculation",
    "init",
    "bindEvents",
    "show",
    "hide",
    "loading",
    "onBizSuccess",
    "onBizFail",
    "initPopup",
    "initEmbed",
    "initFloat",
    "destroyCaptcha",
    "refresh",
    "onCloseClick",
    "startTracelessVerification"
]
```

The StaticPath is randomised for eg 3.25.0/pe.050.9453665072ad79a7
It is used to load a javascript from path
for StaticPath 3.25.0/pe.050.9453665072ad79a7 it loads

https://g.alicdn.com/captcha-frontend/dynamicJS/3.25.0/pe.050.9453665072ad79a7.js

The javascript calculates captcha payload  (data) field sent in verifycaptchav3

Btw all javasripts calculates data in same way just function names are variables are renamed
using mitmproxy i tested to use a single js every time and it works


### 4. VerifyCaptchaV3
```javascript
fetch("https://no8xfe.captcha-open-southeast.aliyuncs.com/", {
  "headers": {
    "accept": "*/*",
  },
  "body": "AccessKeyId=LTAI5tSEBwYMwVKAQGpxmvTd&SignatureMethod=HMAC-SHA1&SignatureVersion=1.0&Format=JSON&Timestamp=2026-05-25T02%3A43%3A00Z&Version=2023-03-05&Action=VerifyCaptchaV3&SceneId=didk33e0&CertifyId=gl62Aqi6e1&CaptchaVerifyParam=%7B%22sceneId%22%3A%22didk33e0%22%2C%22certifyId%22%3A%22gl62Aqi6e1%22%2C%22deviceToken%22%3A%22U0dfV0VCIzM3OTVkMjgyNDJhMTE2MTliYzI1Zjc4NmY4NGU1M2Q0LWgtMTc3OTY3Njk3NzYzMC1kMjA1YTU2MjFmMTY0NDA2YjI5ZTQyMGRkNTUxMzI1MSM3dTU1SjBYc09Uc1IzaXVKS0Zhb09kY1dwVEtxcDBsN2NlRTBHd3VhVy9UVjZ4ZnpVWVhBMWtZQnNWR2pza3A1MW0wak9lMDNadFdHdDV6b3pTOElNOEFPR21ydHUwbnU2ZDdWZUlLbVIrMzNPVkZUYW1VcGtRK0VtVHBTQ1NMa3ZjcXB0NDdFTXdyOFNOMW42SnpWZ0sxMDN6QzA1Q1g1T1Y1bnBSV3hZREdzUFhRVnpTcytxWEFqczhJS3ZiU1dZb0tUNVFFM2ZVR2ZwamNlclNoaWkzSG50OU0ySmlmWmd2amxzSlowYVNqQ0JiL1kwTGlGUlRST1I2Y2FydThMbEdPTFA4TUUxem9hQ1JKWjlSeEE3ekcyVjRmdnBBa0tsU1U4aFZTVkpDMHl2N1ZCaFhuSDMrSWFTWUh1TitRR00vdndoYllqOVNPVUZhQjJDYllQU0xpandobS9HTk1QcVBBQ2xVNllNQndDTVRuYk40MDJmTVorRTVTVVU0eGNXYnBvTnJhdlg1MDVLbktUTEcwbGVmTmpmZVd4VThDL1doWDBVQ1FjZXl1S1owMjM1QjFlZys2cEI3Y2RFR0ZML2xkRTY5VzlkN1U3bzRMYnVYT1IwemZWeDlQKzhmNStBL3dySmZyZXcrSkJBelE4Wm8zYkxScEtac2RjMFZpSENza2ZQeUkwWSsyQStCaHpvd0Q4QXdMQ2tIZisyZ2tBK2xVaGdyRGhEcWlqNU5nL21YSTdwdVFEbUJkRjFKUllzY3ZCbU53OUt4T09xQVBMN2tvaFRqZGw4TjljWUVTdm0zRTQ0bEJ5d21WZEQyOVFmeWwzOWYrbmlHdlJGc0lDaDJMK3ZIMStlTFNBN3B4MFpWV2pMUVI1eUpyeHZ0eUxQVVgxQ2JpZE9FNkp1bEZpd1BpbnlYK0RvTHBEQVZTdnFqRjNWMVY5bVJtMjVqL1YwaDBhMkhkWWV1UTNQNlVOWUtDNk1wc3BPM2JqMmdSR1J1VDA2V050NVRmZWdETWNWeWVFeEZpYWxPWHZia3dhbjZ3ZWYvbEVGUHF3RUVPZ1ROWDAxdEtXUm1VNWxSOGJKQVdjaUhCTnJQOEpOY0FyMGV0Z2xjd1RxVVhOTngrZ2xrYzdkMzJ3Zk5tZ1BCUjV3Q3h3NjdIWTUzVTI3dHNQMjRVVHVOaGc3SkhjNk5kdDBmQ0JvUDRXdmxpaFdEbFBQQ3Vic1FObFZBS1VYZ0JoYUFDMUVMbmt0Ukt6eEZXNWJ1TzlTOUVDT0RBQWdTNVBrZzIzVldpRTJCOW9DbkFyRERpdVpwZDY4bUpHdnljemdZeXFhd0pWSlp6TzVoaTNoazkvNDdlcy85dExGbzAvMmlmcTVTS01kRTdhVkZlNnRtREVCdz09IzAjN2JkNzMyZDIwZTc1NmEyNmY0OTU1ZGJhYjU4MjhjOWE%3D%22%2C%22data%22%3A%22JRMlgg0gDgVASQIXRQRNMFsDRJoGeSB5eQp9YWNPH%2Fw6iT2%2B96hKSApjEY9S1AZt4rwnZQKgGGOjoeIYWpEaAJdJFhEbK15pcCcKYi0%2FNEgaXU4YLj19be54MEMWEBqDqzmjFdV5KuAcSYQ%2FYWYLNeJ3jZ8VzTdU%2FwO4QqMVb8BXhzAENTBkSQ1iXQhwvkMi01UB7yClAj59DW9NKmNkaj1eS%2BQpfSlvKJ1iF2IXEDcbqi4rbP%2B3BDaUnxMkbo0aLSs8RXxbQH8%3D%22%7D&SignatureNonce=bae00c65-c089-41fc-a704-84576ebea034&Signature=azbdVXlzzucKOIz1kICXOmafPJg%3D",
  "method": "POST",
  "mode": "cors",
  "credentials": "omit"
});
```

```json
{
    "RequestId": "F06551B1-2C26-453F-B570-AB7687D82551",
    "Message": "success",
    "HttpStatusCode": 200,
    "Code": "Success",
    "Success": true,
    "Result": {
        "securityToken": "6oOo7e72nA61uVLiZVKiLYqF1m9rOno3vEIPJKaL7KLxCJqb1UBwRpl4p7EcFTgd3yG06TCPBjR35MbCZ5lDrdjPcqaflqbQLZQdX2rYd/8bhnqhIpC7SnRlIxGPsqvX",
        "VerifyCode": "T001",
        "VerifyResult": true,
        "certifyId": "gl62Aqi6e1"
    }
}
```


### 5. captcha_verify_param generation and usage

captcha_verify_param is base64 encoded string which contains scene id certify id and token and is snet in completions requests
```base64
eyJjZXJ0aWZ5SWQiOiJnbDYyQXFpNmUxIiwic2NlbmVJZCI6ImRpZGszM2UwIiwiaXNTaWduIjp0cnVlLCJzZWN1cml0eVRva2VuIjoiNm9PbzdlNzJuQTYxdVZMaVpWS2lMWXFGMW05ck9ubzN2RUlQSkthTDdLTHhDSnFiMVVCd1JwbDRwN0VjRlRnZDN5RzA2VENQQmpSMzVNYkNaNWxEcmRqUGNxYWZscWJRTFpRZFgycllkLzhiaG5xaElwQzdTblJsSXhHUHNxdlgifQ==
```
```utf
{"certifyId":"gl62Aqi6e1","sceneId":"didk33e0","isSign":true,"securityToken":"6oOo7e72nA61uVLiZVKiLYqF1m9rOno3vEIPJKaL7KLxCJqb1UBwRpl4p7EcFTgd3yG06TCPBjR35MbCZ5lDrdjPcqaflqbQLZQdX2rYd/8bhnqhIpC7SnRlIxGPsqvX"}
```


## Verify Captcha payload generation [pe.*.*.js]

### The request

```javascript
fetch("https://no8xfe.captcha-open-southeast.aliyuncs.com/", {
  "headers": {
    "accept": "*/*",
    "accept-language": "en-US,en;q=0.9",
    "content-type": "application/x-www-form-urlencoded; charset=UTF-8",
  },
  "referrerPolicy": "no-referrer",
  "body": "AccessKeyId=LTAI5tSEBwYMwVKAQGpxmvTd&SignatureMethod=HMAC-SHA1&SignatureVersion=1.0&Format=JSON&Timestamp=2026-06-19T06%3A42%3A09Z&Version=2023-03-05&Action=VerifyCaptchaV3&SceneId=didk33e0&CertifyId=KdeAGaejyN&CaptchaVerifyParam=%7B%22sceneId%22%3A%22didk33e0%22%2C%22certifyId%22%3A%22KdeAGaejyN%22%2C%22deviceToken%22%3A%22U0dfV0VCIzM3OTVkMjgyNDJhMTE2MTliYzI1Zjc4NmY4NGU1M2Q0LWgtMTc4MTg1MTQ4NDgxOC0xZDgwY2E3MWE3ZDQ0MzA5OTI4NDEzZjAyY2FhODA0MyM4RDJWN05odUhRWHhkQzhJbWxaa0tiMjNCaUM0WjYwUWtqTlNjOEZnbFhNOFpsb2dFbzZxMVZDcjdDazhVdEM3emlGbDROa3ViUHFnK0NvVkk0RXJmYkowSUwvUlg1RnhITDFYbHd1U3E3MzZWUDhrTnFZZENCdDFtelIrTHRPMGNkRDJBRitIYVZQZkZhTFJac0pMS1FVR1QwWjd6TFFMVVlacTdCeVNYTGhKZWMvd1NkV0FtdjZEUkNLYUpZc0ZvTlkzbGJnbFdHOFVEUDBWdFVjUGt1MjJmVTlobVJmQno4UnlaWmVIeTgzRkF3dXF4bjhsY1htUGZBQ25BSUlLR3pXS2UzWTlocEFOemxwNEoxQStaaTNROFJUeE5tam1vbmxiTjdkb3h3bmozbkpPcERmQlhNeFNNaUxmdHF6MWxtRlFZbzhEbFpNQzVqaEI3UnFFbDZlellHMWJWU3UxQVA4VWFYY1VFdFhSYU43V2pEQk5YaThDWnZVTmRBY1lRRjEyZWwwS0ZZS3A4VGNpVkI1YXBKQk4vQk9Pa1BMcmYzcTYwcHA5cElZSkpZa1pDekRZbGFpVHljMWhEdTYyM2Q4WjN3UWtjOG9scEdKZHY5OElnSnp1a1Riem1jYjlaWEtuWVVVK1FaWEpEbUlBYTVzTy9iaDRYOTlJalZIZmlvc1huN1QzR0dUNU51Ny9tckRyQ2VYR2YyOWJNdzdqeVZHKzVISFVhTHltSVNwTTRZdFQ3NVFpTWVmeGlJaG1XTWNReFBIQ05WWG5RUDFUdnVHTGxieHhUa1RRZ0R1Tldrc0hvalFWbkd1QlB0Sk5zMHpTekxvSzFhQ1dZaGhnbVlXOStNWXVWR1Zkekx0TnJJY2wyRlF2U1F2N1drdVRVOXFtUFlQNEJBbCtIMWh1QVdadFVxWllENzZVLzhUWi9XV0NTM1hsU05ycEFFTTVUa1FlVVdHN0hkaHUwYmNzUHJqUU44NTE4Ym11a1M5YUZxKy9YQnJxcjRzMjFXTWNXaGRXcVpnS3NJaGVneGh3YnlMYW12ajN3YVQ1a1d3eEZzT0RRMGVycEJUdUFlRjRYbVVhNmtXQmZSNWZaVXlmIzAjMDdiOTE3ODA5YTRkNWEzMjE5YjEyZjk4NTI1ZTNjMWY%3D%22%2C%22data%22%3A%22JRMlgg03DgRAeAILRQpEEHPjZZg9Ez1fbA0qd1yjNhsriXK5OHh3ZgkJc2NaO%2BpbP7wtcQGL6Gvi5LEZqOMpPHYyKHIfAmJWGUMiBRsQD0kkUUkiCA1Jd8BAFFEFUmL%2BQMu4MyFaXRQ8WWI%2BZ3Ib%2BByNmoNjgz5I0ECwZobWT5V6ii5ADARXc1BcB3ZnpVMtJn97EAZbYV5%2BLy9SEyhzMjNYnyF8fCUVJIx0UDQ2LMIfjHwgCMu0MPWbvABhaIf6ZjkXeldLQW4%3D%22%7D&SignatureNonce=0ecd84d6-b4ce-473c-b4d3-1cb4d0dcbca9&Signature=Oyh7ojzWm1sLLw%2BdNzrs8Ulfb88%3D",
  "method": "POST",
  "mode": "cors",
  "credentials": "omit"
});
```


### Successfull Response
```json
{"RequestId":"037B581E-C73A-4CFB-A23D-AB00F498CC0D","Message":"success","HttpStatusCode":200,"Code":"Success","Success":true,"Result":{"securityToken":"6oOo7e72nA61uVLiZVKiLYqF1m9rOno3vEIPJKaL7KLxCJqb1UBwRpl4p7EcFTgdXev7OeyuHCpA5PzF5PUSndjPcqaflqbQLZQdX2rYd/8bhnqhIpC7SnRlIxGPsqvX","VerifyCode":"T001","VerifyResult":true,"certifyId":"KdeAGaejyN"}}
```

### Payload Data
```json
"data":"JRMlgg03DgRAeAILRQpEEHPjZZg9Ez1fbA0qd1yjNhsriXK5OHh3ZgkJc2NaO+pbP7wtcQGL6Gvi5LEZqOMpPHYyKHIfAmJWGUMiBRsQD0kkUUkiCA1Jd8BAFFEFUmL+QMu4MyFaXRQ8WWI+Z3Ib+ByNmoNjgz5I0ECwZobWT5V6ii5ADARXc1BcB3ZnpVMtJn97EAZbYV5+Ly9SEyhzMjNYnyF8fCUVJIx0UDQ2LMIfjHwgCMu0MPWbvABhaIf6ZjkXeldLQW4="
```
The data is base64 encoded likely encrypted zlib compressed and custom hash prepended data 

The actual data looks like : hash+json `'9333ef7396dd56dbb9d6e8f31e8f6014{"TrackList":{"mc":"","tc":"","mu":"","te":"","mp":"","tmv":"","ks":"","fi":"","startTime":1782100652835},"TrackStartTime":1782100652835,"VerifyTime":1782100652862,"arg":"JjObDGdh/ywcWQ=="}'`
mc = mouse clicks
tc = touch?
mu = mouse up?
etc.
arg is base64 encoded 10 byte value

#### Entire VerifyCaptcha is done by pe.*.*.js file returned in init captcha response

### Variables and their definitions
```javascript
var ta = ({
            0: te
        })[0](40, 62)
          , to = window[ta + (~te ? te : 5)(29, 264)]
          , tc = n[(te && te)(30, 188)]
          , ts = to[te.apply(2, [97, 51])]
          , tu = to[te(8. .valueOf(), 15. .valueOf())][(-te ? 4 : te)(19, 273)]
          , tf = to[te(8, 15)][te(~te && 90, ~te && 86)]
          , tl = to[te((te(),
        8), (te(),
        15))][(te || te)(64, 38)]
          , th = to[({
            0: te
        })[0](62, 187)][te.apply(5, [53, 149])]
          , tp = to[te.apply(1, [37, 261])]
          , td = tf[(te(),
        te)(3, 204) + "y"](tl[te(63 * (1 | te), 50 * (1 | te))](tc))
          , tA = {};
        tA.iv = tu[[te][0](63, 50)](td),
        tA[te(-te ? 5 : 82, 234 * !-te)] = th;
        var tv = tA;

```

### The values
```
ta
'__ALIYUN
to = window.__ALIYUN_CRYPT


window.__ALIYUN_CRYPT contains
AES
: 
{encrypt: Æ, decrypt: Æ}
DES
: 
{encrypt: Æ, decrypt: Æ}
EvpKDF
: 
 (t,r,e)
HmacMD5
: 
 (r,e)
HmacRIPEMD160
: 
 (r,e)
HmacSHA1
: 
 (r,e)
HmacSHA3
: 
 (r,e)
HmacSHA224
: 
 (r,e)
HmacSHA256
: 
 (r,e)
HmacSHA384
: 
 (r,e)
HmacSHA512
: 
 (r,e)
MD5
: 
 (r,e)
PBKDF2
: 
 (t,r,e)
RC4
: 
{encrypt: Æ, decrypt: Æ}
RC4Drop
: 
{encrypt: Æ, decrypt: Æ}
RIPEMD160
: 
Æ (r,e)
Rabbit
: 
{encrypt: Æ, decrypt: Æ}
RabbitLegacy
: 
{encrypt: Æ, decrypt: Æ}
SHA1
: 
Æ (r,e)
SHA3
: 
Æ (r,e)
SHA224
: 
Æ (r,e)
SHA256
: 
Æ (r,e)
SHA384
: 
Æ (r,e)
SHA512
: 
Æ (r,e)
TripleDES
: 
{encrypt: Æ, decrypt: Æ}
algo
: 
{MD5: {â€¦}, SHA1: {â€¦}, SHA256: {â€¦}, SHA224: {â€¦}, SHA512: {â€¦}, â€¦}
computeSignature
: 
Æ Ce(t,r)
enc
: 
{Hex: {â€¦}, Latin1: {â€¦}, Utf8: {â€¦}, Utf16BE: {â€¦}, Utf16: {â€¦}, â€¦}
format
: 
{OpenSSL: {â€¦}, Hex: {â€¦}}
kdf
: 
{OpenSSL: {â€¦}}
lib
: 
{Base: {â€¦}, WordArray: {â€¦}, BufferedBlockAlgorithm: {â€¦}, Hasher: {â€¦}, Cipher: {â€¦}, â€¦}
mode
: 
{CBC: {â€¦}, CFB: {â€¦}, CTR: {â€¦}, CTRGladman: {â€¦}, OFB: {â€¦}, â€¦}
pad
: 
{Pkcs7: {â€¦}, AnsiX923: {â€¦}, Iso10126: {â€¦}, Iso97971: {â€¦}, ZeroPadding: {â€¦}, â€¦}
x64
: 
{Word: {â€¦}, WordArray: {â€¦}}


tc
'd35db7e39ebbf3d001083105'
ts
{encrypt: Æ, decrypt: Æ}
tu
{stringify: Æ, parse: Æ}
tl
{stringify: Æ, parse: Æ}
tf // Base64 encode decode?
{_map: 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=', _reverseMap: Array(123), stringify: Æ, parse: Æ}
parse
: 
Æ (t)
stringify
: 
Æ (t)
_map
: 
"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
_reverseMap
: 
(123) [empty Ã 43, 62, empty Ã 3, 63, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, empty Ã 3, 64, empty Ã 3, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, empty Ã 6, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51]
[[Prototype]]
: 
Object
th
{pad: Æ, unpad: Æ}
tp
Æ (r,e){return new p.HMAC.init(t,e).finalize(r)}
td
'0123456789ABCDEF'
tA
{iv: t.init, padding: {â€¦}}
iv
: 
t.init
sigBytes
: 
16
words
: 
Array(4)
0
: 
808530483
1
: 
875902519
2
: 
943276354
3
: 
1128547654
length
: 
4

padding
: 
pad
: 
Æ (t,r)
unpad
: 
Æ (t)
[[Prototype]]
: 
Object
[[Prototype]]
: 
Object
ï»¿

w
'DNsKKPRHfOTuFph6taPHkC+TbC6YUaXi1NfQv95Ez2k='
x
'd35db7e39ebbf3d001083105'

v
{ID: '7JLsB18MnA7GX3d6LxErT1sGT68xcVuOAoxz0b7vVzY=', SECRET: 'n9jH0yACW8YrgOBcM0v7u45+/bfozcSz8ZpvzGBXg3E='}



M
VerifyCap

```

### Zlib compression
Zlib compression is done by function K

hash+json is passed to funtion K ty(t) converts string to uint8array G.deflate zlib compresses the data

```javascript
// Example
console.log(t)
'99603d2757539169328062db7568942c{"TrackList":{"mc":"","tc":"","mu":"","te":"","mp":"","tmv":"","ks":"","fi":"","startTime":1781878507084},"TrackStartTime":1781878507084,"VerifyTime":1781878507099,"arg":"VzN0TRFfHQ5OVw=="}'

console.log('Before Defalte',n);

Before Deflate
Uint8Array(221) [57, 57, 54, 48, 51, 100, 50, 55, 53, 55, 53, 51, 57, 49, 54, 57, 51, 50, 56, 48, 54, 50, 100, 98, 55, 53, 54, 56, 57, 52, 50, 99, 123, 34, 84, 114, 97, 99, 107, 76, 105, 115, 116, 34, 58, 123, 34, 109, 99, 34, 58, 34, 34, 44, 34, 116, 99, 34, 58, 34, 34, 44, 34, 109, 117, 34, 58, 34, 34, 44, 34, 116, 101, 34, 58, 34, 34, 44, 34, 109, 112, 34, 58, 34, 34, 44, 34, 116, 109, 118, 34, 58, 34, 34, 44, 34, 107, 115, 34, 58, â€¦]


G.deflate(n)
Uint8Array(149) [120, 156, 117, 141, 203, 10, 194, 48, 16, 69, 255, 101, 214, 89, 164, 73, 243, 152, 66, 183, 210, 133, 40, 106, 233, 190, 166, 173, 132, 18, 144, 36, 42, 90, 250, 239, 22, 146, 157, 184, 58, 247, 112, 47, 51, 136, 146, 242, 129, 41, 161, 4, 199, 66, 34, 103, 154, 74, 54, 92, 149, 144, 26, 75, 102, 22, 104, 125, 111, 230, 189, 13, 17, 170, 5, 156, 129, 10, 128, 64,
204, 116, 143, 236, 99, 246, 123, 118, 247, 76, 97, 14, 137, 147, 77, 12, 177, 247, 177, 181, â€¦]


G
{constants: {â€¦}, default: {â€¦}, Deflate: Æ, deflate: Æ, deflateRaw: Æ, â€¦}
Deflate
: 
Æ tW(t)
constants
: 
{Z_NO_FLUSH: 0, Z_PARTIAL_FLUSH: 1, Z_SYNC_FLUSH: 2, Z_FULL_FLUSH: 3, Z_FINISH: 4, â€¦}
default
: 
{constants: {â€¦}, Deflate: Æ, deflate: Æ, deflateRaw: Æ, gzip: Æ}
deflate
: 
Æ tQ(t,n)
length
: 
2
name
: 
"tQ"
prototype
: 
{}
: 
Æ ()
[[Scopes]]
: 
Scopes[4]
deflateRaw
: 
Æ (t,n)
gzip
: 
Æ (t,n)
__esModule
: 
true


```
```javascript
function tQ(t, n) {
                    var e = new tW(n);
                    if (e.push(t, !0),
                    e.err)
                        throw e.msg || R[e.err];
                    return e.result
                }
```

```javascript
var tU = function(t) {   if ("function" == typeof TextEncoder && TextEncoder.prototype.encode)
                        return (new TextEncoder).encode(t);
                    var n, e, r, i, a, o = t.length, c = 0;
                    for (i = 0; i < o; i++)
                        55296 == (64512 & (e = t.charCodeAt(i))) && i + 1 < o && 56320 == (64512 & (r = t.charCodeAt(i + 1))) && (e = 65536 + (e - 55296 << 10) + (r - 56320),
                        i++),
                        c += e < 128 ? 1 : e < 2048 ? 2 : e < 65536 ? 3 : 4;
                    for (n = new Uint8Array(c),
                    a = 0,
                    i = 0; a < c; i++)
                        55296 == (64512 & (e = t.charCodeAt(i))) && i + 1 < o && 56320 == (64512 & (r = t.charCodeAt(i + 1))) && (e = 65536 + (e - 55296 << 10) + (r - 56320),
                        i++),
                        e < 128 ? n[a++] = e : (e < 2048 ? n[a++] = 192 | e >>> 6 : (e < 65536 ? n[a++] = 224 | e >>> 12 : (n[a++] = 240 | e >>> 18,
                        n[a++] = 128 | e >>> 12 & 63),
                        n[a++] = 128 | e >>> 6 & 63),
                        n[a++] = 128 | 63 & e);
                    return n
                }
                  , tj = function() {
                    this.input = null,
                    this.next_in = 0,
                    this.avail_in = 0,
                    this.total_in = 0,
                    this.output = null,
                    this.next_out = 0,
                    this.avail_out = 0,
                    this.total_out = 0,
                    this.msg = "",
                    this.state = null,
                    this.data_type = 2,
                    this.adler = 0
                }
                  , tP = Object.prototype.toString
                  , tL = F.Z_NO_FLUSH
                  , tR = F.Z_SYNC_FLUSH
                  , tF = F.Z_FULL_FLUSH
                  , tV = F.Z_FINISH
                  , tG = F.Z_OK
                  , tD = F.Z_STREAM_END
                  , tq = F.Z_DEFAULT_COMPRESSION
                  , tH = F.Z_DEFAULT_STRATEGY
                  , tZ = F.Z_DEFLATED;
```




### Transformation of trackJson is done by fucntion nf as seen in this call C=nx.A(nf,n_)
```
nx.A = function(t,n) { return t(n)}
C = nf(n_) // n_ is trackJson without hash
```



```javascript

function nf(t) {
            for (e = 9; e; )
                switch (r = e >> 3,
                i = 7 & e,
                r) {
                case 0:
                    if (i < 6)
                        if (i < 2)
                            i >= 1 && (e = 10);
                        else if (i > 3)
                            if (i > 4) {
                                e -= -3;
                                try {
                                    return function(t) {
                                        var n, e, r, i;
                                        for (e = 3; e; )
                                            e < 1 || (e <= 1 ? (e += 3,
                                            r = function(t, n) {
                                                return i.B(te, n, t - -3)
                                            }
                                            ) : e < 4 ? e < 3 ? (e ^= 3,
                                            i.B = function(t, n, e) {
                                                return t(n, e)
                                            }
                                            ) : (i = {},
                                            e += -1) : (n = tr[(r && r)(52, 64)](this, 26)[i.B(r, Math.floor(237), 19)](this, arguments),
                                            e = 0));
                                        return n
                                    }(t)
                                } catch (t) {
                                    if (c.e(s, a))
                                        throw t
                                }
                            } else
                                !c * !c / 0 != 6 ? e ^= 1 : e += 3;
                        else
                            i > 2 ? (e = 14,
                            c.D = function(t, n) {
                                return t <= n
                            }
                            ) : (e += 14,
                            c.n = function(t, n) {
                                return t && n
                            }
                            );
                    else
                        i > 6 ? (e += -4,
                        c.R = function(t, n) {
                            return t !== n
                        }
                        ) : (a = c.y(c.W(c.c(arguments[o.call(1, 229, 51)], 1), 61), 58) > 58 && c.R(arguments[1], void 0) ? arguments[1] : 3,
                        e = 1);
                    break;
                case 1:
                    i > 5 ? i <= 6 ? (c.e = function(t, n) {
                        return t === n
                    }
                    ,
                    e -= 2) : (e -= 2,
                    c.W = function(t, n) {
                        return t * n
                    }
                    ) : i > 2 ? i < 5 ? i < 4 ? c.D(c.W(s - a, 50) + -6, -6) ? e ^= 15 : e = 0 : (o = function(t, n) {
                        return c.n(te, te)(n, t - 1)
                    }
                    ,
                    e ^= 10) : (e += -6,
                    c.c = function(t, n) {
                        return t - n
                    }
                    ) : i >= 1 ? i < 2 ? (e -= 7,
                    c = {}) : (s = 1,
                    e += 1) : (e = 11,
                    s++);
                    break;
                case 2:
                    e ^= 31,
                    c.y = function(t, n) {
                        return t + n
                    }
                }
            return n
        }
```

t = trackJson ( no Hash)

This gets transformed into final data with this code  n = tr[(r && r)(52, 64)](this, 26)[i.B(r, Math.floor(237), 19)](this, arguments)
console.log(n);
JRMlgg0gGwVASQITewRNMEYjU1P7AMM7DGwWaiVYKOEKajqEHFh8RA1neYRfNAdR0XUQfSWw7h2PmQTbVpMbBWVFEHMdI1sENzToeHBDa0gdW0MRcxRcTwxEEp82HSfvd8K7Ri0JKSw7UW8jV31K21N3YXwcwT937U+vWYTVSox0hsdjInZNDW8WB3FzQ0Qkw0sBAAJLKT5CKRN3xwaLCX8bYxoLkQ9oKo5XTDQAZitmn0Qhdc/1Suikdj8OaWwhRjtEPlwSHmc=

### Function nf behaviour 
```js
console.log(t)

{TrackList: {â€¦}, TrackStartTime: 1782021569701, VerifyTime: 1782021569717, arg: 'ZjyUTmpv9h8dBw=='}
```

```js
nf(t)

'JRMlgg0gGwVQeQITewRNMEbgUlclHnFEaXA2MWKm3fMsc1CfFaYZZjkfd3ASHyFF9bUHGneULV+hjAnpgZcDBJZsEmwZIGwMdCsQFi0xAFYVUUQYLwB4U9Y0MEgfFifSRy9VO85ZBx85VXMrSUtHPgKyQZJfzTF5DwWuBYXbfJwJmygJIBo0Y35rFRkTplQm0EhfERO9Kz98GhVWwX+bJjFlaxt4niZQxnEVEyYCehkVUVEwa8yXCjKUlyBtVk0KRx54Bg=='
```

```js
nf('abcdef123456789')
'JRMkbRAlaBM6LgIPamJf6WPtb2UkBx5abH4bHmS2BThZbDKC7ah/Xt0bN595ORYu8k0JPXOv6X2Vr8cJg5E7BpBSMVY='      
```
`Tests:` **whether function nf uses zlib compression or not** 
`Setup:` **In function K added logger that logs arguments**
```js
nf(t)

[FUNCTION K LOG] ARGUMENTS: 'c6d9b3718aa286fe7290fbab8e7f66e2{"TrackList":{"mc":"","tc":"","mu":"","te":"","mp":"","tmv":"","ks":"","fi":"","startTime":1782022235313},"TrackStartTime":1782022235313,"VerifyTime":1782022235354,"arg":"fMR+GVZSAyc4Ng=="}'

nf('abcdef123456789')
'JRMlZ25EVwYlLzQ+WWUF5GbgRHolfik0SDgDCUCqwMVUaUCpF3FzZgF1ZaB10PhN0EpybDWVO3Kb7fIc'

[FUNCTION K LOG] ARGUMENTS: 
'f4aaf499d633c77792ccf444f4dd0b33"abcdef123456789"'
```

### Relation
```js
tr(26, 'test')

'JRMkWWkGDigveSB2eABaM2cdclsffhBebh0gL0OrLcAuWVtA/UUFYQYbfIx0PgkT1VhpMBetN2vjjcQOioEnDQ=='
```
```js
nf('test')
'JRMkWWkGDigveSB2eABaM2cdclsffhBebh0gL0OrLcAuWVtA/UUFYQYbfIx0PgkT1VhpMBetN2vjjcQOioEnDQ=='
```

`Most likely Reason` 
**This code:**
```js
n = tr[(r && r)(52, 64)](this, 26)[i.B(r, Math.floor(237), 19)](this, arguments), e = 0
```
**My Thoughts: Correct me if im wrong**
```js
i.B = function(t, n, e) { return t(n, e) }
and r = function(t, n) { return i.B(te, n, t + 3) }
//   which simplifies to: r(t, n) = te(n, t + 3)

r && r                          // r is truthy → evaluates to r itself
r(52, 64)                       // call r with (52, 64)
= te(64, 52 + 3)                // r(t,n) = te(n, t+3)
= te(64, 55)                    // → returns a string (method name)
```
`Let's call this methodA = te(64, 55)`
```js
tr[te(64, 55)]                  // look up methodA on tr object
(this, 26)                      // call it with this=context, 26=selector

Math.floor(237)                 // = 237 (obfuscation noise, pure 237)
i.B(r, 237, 19)                 // i.B(t,n,e) = t(n,e)
= r(237, 19)                    // call r with (237, 19)
= te(19, 237 + 3)               // r(t,n) = te(n, t+3)
= te(19, 240)                   // → returns a string (method name)
```
**Now debug result shows:** 
```js
r(52, 64)
'bind'
r(237, 19)
'apply'
te(64, 55)
'bind'
te(19, 240)
'apply
```
```js
// Resolved:
n = tr.bind(this, 26).apply(this, arguments)

Once again 
n = tr(26, t) 
```

### Function tr behaviour

The first tr call when i trigger verify captcha is:
```js
console.log(t,e,r,i,a,o);  tr(t, e, r, i, a, o);
VM41108:1 35 undefined undefined undefined undefined undefined

true
```
The second call is 
```js
console.log(t,e,r,i,a,o);  tr(t, e, r, i, a, o); // decryption?
VM42148:1 9 'FqJB6iRNVYdEGpwb' '7JLsB18MnA7GX3d6LxErT1sGT68xcVuOAoxz0b7vVzY=' undefined undefined undefined

'LTAI5tSEBwYMwVKAQGpxmvTd'
```
Third call is 
```js
console.log(t,e,r,i,a,o);  tr(t, e, r, i, a, o); // decryption again?
VM42659:1 9 'FqJB6iRNVYdEGpwb' 'n9jH0yACW8YrgOBcM0v7u45+/bfozcSz8ZpvzGBXg3E=' undefined undefined undefined

'YSKfst7GaVkXwZYvVihJsKF9r89koz'
```
Fourth call is
```js
console.log(t,e,r,i,a,o);  tr(t, e, r, i, a, o); // xhr?
VM44730:1 33 rA 𝑓 {$button: button#chat-captcha-trigger, captchaVerifyCallback: undefined, onBizResultCallback: undefined, success:, fail:} {TrackList: {â€¦}, TrackStartTime: 1782100652835, VerifyTime: 1782100652862} undefined undefined undefined

Promise {<pending>}
```
At this point r doesnt have arg in it

Fifth call is
```js
console.log(t,e,r,i,a,o);  tr(t, e, r, i, a, o); // device token generation 
VM47310:1 36 '6iL4denBvY' undefined undefined undefined undefined

'U0dfV0VCIzM3OTVkMjgyNDJhMTE2MTliYzI1Zjc4NmY4NGU1M2Q0LWgtMTc4MjA5MzM1NzYyMC0wY2U5YzY2NTU2N2E0N2Q3OTBlNzNkNGE3MTgwNDMwNCNrUHNRYzl2NzB2YXBNdEl6YVVTWWZ4OUFndXBRSzd2RW5QOVpobmgxb1RsODBKWTVpYXFlQkhlWExIR2hDWUlmRHU2WkgrQXNQVk41SHNEWlB5MUNmOWVjVnExb1F6NGtwUkxJQkl1NXRhbVR4d0ZidERuY0ozT2FERTdqdDhtd2dIQ2QrYU5PQ0hlalJlU1A0L0ZoOTFUQzVFbWpCNzZMazQ1UzJvOHNWZmgrTHN0dWcvTW53emtBSUc1L0RBOXc5cGU0NWd1R0RKV21NSHBKTXJ5cm40VGpiTEsvcUxDa1hkcWlIdGkxbUpaRmxUd3N3MkVVcjAybVBxVGIzc3RZVS9INEQwSDVVZHpIUnovZ01jc1BmSnF3WEsyVmVVT2lNTjRYN3RXM3NqcWZadmxNcGtXelpiVnVMRWRLQytuY3ZqbDVZV0NGc0RuRjZrNGpYTkgyT25QdmY0bElIZCtxT25Fb2Ztem5KQ21YZ0JSTnNXK3h0WE90Wng5L0tXeVRWNVpYbHc1dmVHdk11bHhHQTQ0Z1ZZa0FqZExQdDN4ZktXTlFYUFB4K2lUNWdrbTZlRWx4QUV4TlVPanVaUzlHVW83NDcwN1RnbmQ4Y1lGelBxQUxjVDZUWXVoZEVLWEU0aXpoMThsSjR5UDgzSWpCeTFhWER3OVhnRUlLNk5HYzhxTTI1YlVlSVJpSGYwdXpLS0h5NlhGbjlWcmdJUGpUUFI0UmFFRW9XRUg3N0dwdVpkaDE0Yi9YRUQ2b3pUYWh5c3dpZVU2Sk9rZEZPWkRXdEtUbGx6MWMwZjRaWGpxaWIzN3VRRUo5NnJZUG80V2M2MGVXTDZFSVUvVDkraS9lb3J4SE8xTFhPVHoyYWhkYy92bWpIV3RPcWtFWlJXS1NBTFhvNmxuc0tuWjhoZmlxdnY3VGhJaVc0Wm9HZzJVMHJCeWtqajM3SnhBZnNmVXRQRm5oQ0JEWFVTakZhR3lMVjVKVmJjckp0UVdXMlRQRFBia0JrdXhwajI3WlNqV0NBRzBiSTlLRld1MUwzakhaVXEvQjYvYWZ6NmhFd1kyL0ZWczdzMDA0d3lsQnRHYVZtRC9zanlSMFF5STNEMTNtYlQ5d0pORUtpRFNLZ0U5dTVsTytuckswcDlxVUU5RCtVSmUyOUNXUnVWST0jNDQ4I2JiYjRlMmE5MjE5NzFmMDE2Nzk4NjRhZDk4MDEwOTlj'
```
Note:
```js
decoded = SG_WEB#3795d28242a11619bc25f786f84e53d4-h-1782093357620-0ce9c665567a47d790e73d4a71804304#kPsQc9v70vapMtIzaUSYfx9AgupQK7vEnP9Zhnh1oTl80JY5iaqeBHeXLHGhCYIfDu6ZH+AsPVN5HsDZPy1Cf9ecVq1oQz4kpRLIBIu5tamTxwFbtDncJ3OaDE7jt8mwgHCd+aNOCHejReSP4/Fh91TC5EmjB76Lk45S2o8sVfh+Lstug/MnwzkAIG5/DA9w9pe45guGDJWmMHpJMryrn4TjbLK/qLCkXdqiHti1mJZFlTwsw2EUr02mPqTb3stYU/H4D0H5UdzHRz/gMcsPfJqwXK2VeUOiMN4X7tW3sjqfZvlMpkWzZbVuLEdKC+ncvjl5YWCFsDnF6k4jXNH2OnPvf4lIHd+qOnEofmznJCmXgBRNsW+xtXOtZx9/KWyTV5ZXlw5veGvMulxGA44gVYkAjdLPt3xfKWNQXPPx+iT5gkm6eElxAExNUOjuZS9GUo74707Tgnd8cYFzPqALcT6TYuhdEKXE4izh18lJ4yP83IjBy1aXDw9XgEIK6NGc8qM25bUeIRiHf0uzKKHy6XFn9VrgIPjTPR4RaEEoWEH77GpuZdh14b/XED6ozTahyswieU6JOkdFOZDWtKTllz1c0f4ZXjqib37uQEJ96rYPo4Wc60eWL6EIU/T9+i/eorxHO1LXOTz2ahdc/vmjHWtOqkEZRWKSALXo6lnsKnZ8hfiqvv7ThIiW4ZoGg2U0rBykjj37JxAfsfUtPFnhCBDXUSjFaGyLV5JVbcrJtQWW2TPDPbkBkuxpj27ZSjWCAG0bI9KFWu1L3jHZUq/B6/afz6hEwY2/FVs7s004wylBtGaVmD/sjyR0QyI3D13mbT9wJNEKiDSKgE9u5lO+nrK0p9qUE9D+UJe29CWRuVI=#448#bbb4e2a921971f01679864ad9801099c
```
`this part:` kPsQc9v70vapMtIzaUSYfx9AgupQK7vEnP9Zhnh1oTl80JY5iaqeBHeXLHGhCYIfDu6ZH+AsPVN5HsDZPy1Cf9ecVq1oQz4kpRLIBIu5tamTxwFbtDncJ3OaDE7jt8mwgHCd+aNOCHejReSP4/Fh91TC5EmjB76Lk45S2o8sVfh+Lstug/MnwzkAIG5/DA9w9pe45guGDJWmMHpJMryrn4TjbLK/qLCkXdqiHti1mJZFlTwsw2EUr02mPqTb3stYU/H4D0H5UdzHRz/gMcsPfJqwXK2VeUOiMN4X7tW3sjqfZvlMpkWzZbVuLEdKC+ncvjl5YWCFsDnF6k4jXNH2OnPvf4lIHd+qOnEofmznJCmXgBRNsW+xtXOtZx9/KWyTV5ZXlw5veGvMulxGA44gVYkAjdLPt3xfKWNQXPPx+iT5gkm6eElxAExNUOjuZS9GUo74707Tgnd8cYFzPqALcT6TYuhdEKXE4izh18lJ4yP83IjBy1aXDw9XgEIK6NGc8qM25bUeIRiHf0uzKKHy6XFn9VrgIPjTPR4RaEEoWEH77GpuZdh14b/XED6ozTahyswieU6JOkdFOZDWtKTllz1c0f4ZXjqib37uQEJ96rYPo4Wc60eWL6EIU/T9+i/eorxHO1LXOTz2ahdc/vmjHWtOqkEZRWKSALXo6lnsKnZ8hfiqvv7ThIiW4ZoGg2U0rBykjj37JxAfsfUtPFnhCBDXUSjFaGyLV5JVbcrJtQWW2TPDPbkBkuxpj27ZSjWCAG0bI9KFWu1L3jHZUq/B6/afz6hEwY2/FVs7s004wylBtGaVmD/sjyR0QyI3D13mbT9wJNEKiDSKgE9u5lO+nrK0p9qUE9D+UJe29CWRuVI=
**is aes encryted with key 31646664393836656436633062643262 and iv words [808530483,875902519,943276354, 1128547654 ],sigBytes: 16**

Sixth call is  
```js
console.log(t,e,r,i,a,o);  tr(t, e, r, i, a, o);
VM53596:1 26 {TrackList: {â€¦}, TrackStartTime: 1782100652835, VerifyTime: 1782100652862, arg: 'JjObDGdh/ywcWQ=='} undefined undefined undefined undefined

'JRMlgg0wDgRASAITRQpNEHEZuJsRBQZDtzsmAXNCL+49eDyEFKlLbSVgFq8aNiIiME9qeyVTKaKkiAQCrb8LIk1vCHMXJG0MOWjmeHE5GGodLkQQXSRFY+BCdntpe1A1ez2zBcgfEfU6c2QvXUxD9xl6TUFapw8G+yeTcaTTe5F3mQlxBGNkcnwbeHhhRzdYKVZ3+jKsH359DR9NImN7DT1alfp/cilvJJwZemkAVjVrnyssVj6RPRt0dV0ejXAUVxYEEnZxQEA='
```
Seventh call converts into Uint8array 
```js
console.log(t,e,r,i,a,o);  tr(t, e, r, i, a, o);
VM53612:1 30 '9333ef7396dd56dbb9d6e8f31e8f6014{"TrackList":{"mc":"","tc":"","mu":"","te":"","mp":"","tmv":"","ks":"","fi":"","startTime":1782100652835},"TrackStartTime":1782100652835,"VerifyTime":1782100652862,"arg":"JjObDGdh/ywcWQ=="}' undefined undefined undefined undefined
```
This returns uint8array of string
```js
Uint8Array(221)Â [57, 51, 51, 51, 101, 102, 55, 51, 57, 54, 100, 100, 53, 54, 100, 98, 98, 57, 100, 54, 101, 56, 102, 51, 49, 101, 56, 102, 54, 48, 49, 52, 123, 34, 84, 114, 97, 99, 107, 76, 105, 115, 116, 34, 58, 123, 34, 109, 99, 34, 58, 34, 34, 44, 34, 116, 99, 34, 58, 34, 34, 44, 34, 109, 117, 34, 58, 34, 34, 44, 34, 116, 101, 34, 58, 34, 34, 44, 34, 109, 112, 34, 58, 34, 34, 44, 34, 116, 109, 118, 34, 58, 34, 34, 44, 34, 107, 115, 34, 58,Â â€¦]
```

Eigth call is 
```js
console.log(t,e,r,i,a,o);  tr(t, e, r, i, a, o);
VM54357:1 26 {TrackList: {â€¦}, TrackStartTime: 1782100652835, VerifyTime: 1782100652862, arg: 'JjObDGdh/ywcWQ=='} undefined undefined undefined undefined

'JRMlgg0wDgRASAITRQpNEHEZuJsRBQZDtzsmAXNCL+49eDyEFKlLbSVgFq8aNiIiME9qeyVTKaKkiAQCrb8LIk1vCHMXJG0MOWjmeHE5GGodLkQQXSRFY+BCdntpe1A1ez2zBcgfEfU6c2QvXUxD9xl6TUFapw8G+yeTcaTTe5F3mQlxBGNkcnwbeHhhRzdYKVZ3+jKsH359DR9NImN7DT1alfp/cilvJJwZemkAVjVrnyssVj6RPRt0dV0ejXAUVxYEEnZxQEA='
```


Ninth call is what sends the verifycaptchav3 request
```js
console.log(t,e,r,i,a,o);  tr(t, e, r, i, a, o);
VM56392:1 22 '{"sceneId":"didk33e0","certifyId":"6iL4denBvY","deviceToken":"U0dfV0VCIzM3OTVkMjgyNDJhMTE2MTliYzI1Zjc4NmY4NGU1M2Q0LWgtMTc4MjA5MzM1NzYyMC0wY2U5YzY2NTU2N2E0N2Q3OTBlNzNkNGE3MTgwNDMwNCNrUHNRYzl2NzB2YXBNdEl6YVVTWWZ4OUFndXBRSzd2RW5QOVpobmgxb1RsODBKWTVpYXFlQkhlWExIR2hDWUlmRHU2WkgrQXNQVk41SHNEWlB5MUNmOWVjVnExb1F6NGtwUkxJQkl1NXRhbVR4d0ZidERuY0ozT2FERTdqdDhtd2dIQ2QrYU5PQ0hlalJlU1A0L0ZoOTFUQzVFbWpCNzZMazQ1UzJvOHNWZmgrTHN0dWcvTW53emtBSUc1L0RBOXc5cGU0NWd1R0RKV21NSHBKTXJ5cm40VGpiTEsvcUxDa1hkcWlIdGkxbUpaRmxUd3N3MkVVcjAybVBxVGIzc3RZVS9INEQwSDVVZHpIUnovZ01jc1BmSnF3WEsyVmVVT2lNTjRYN3RXM3Nqb0wvdW4xWHpGSFh4OE5yQ1d3b093cVVJMkVPVmdKeCtBNTYvR3FETlQyTGU4WFdudFdPcDcxYjhvY0ZTQkhJNytlVUljaGJRendXN2lmcnBkNnU1WTArdEZacjYvQ1J3ajRsWUlTUkJ1WFg0VTdueTFYOGkyd3JESXdMWDZKZkQ0SHZha1c1dFdEM09QOUpTczNVelFTMjJwSXZjYkY4cHF3WHdxdFgwU3MzNC9Qa0hvRDBxL3NacFZsWDEzL1hBYTcvZ1VWUUVXRWxyTXVpZis3TS9jaEdLV0NqUXF6dnZ6WEFCaDJFY3padk5PNG9lLy8yN2lTQTk5SG5BZDRWWDlmczBOV1czUUxxd1lsR2N6R2o2NmxXVk0wYWxzbHl2YWFXMjBHcE9XVzBkdGFUVXpJZVZvVUZTUHlEZWkyTHNJaWtuS0VlWEJvM3RJOGJRdDcrSklTRHdwVlZhRVdEU1FXbFBJMXdEeG5pOUtqY3prNWZHcEgvUlpNNmlnUFU3VkZtL2xrdmFDL2lQWXVPMC9jNThnYk5JTnVqUWRvS1dZWjk3NE04WlhjOHdUYXJLb0R6MUVPSCs3L0h1VWlkRVVHQ2M3b3VXYXY4ajMvL2ZqWGdDZmFiQ01meEJFcWRoaExEZ1NXL3ozVENxQ1VNQ2pQdTNFdVRpNjRqcjBSeUJPaGl2THVSNnk5U0YydmltamdVVHovM3A2M3NzV0VSTFF4My81aEg1T1hNeENHSG5GL3p0dmdNSjRTUmNJZmlyST0jNDQ4Izg1YjgxN2M2YmM4NzNjNTlkZDQ2OWRiMzhmNTk5YjA0","data":"JRMlgg0wDgRASAITRQpNEHEZuJsRBQZDtzsmAXNCL+49eDyEFKlLbSVgFq8aNiIiME9qeyVTKaKkiAQCrb8LIk1vCHMXJG0MOWjmeHE5GGodLkQQXSRFY+BCdntpe1A1ez2zBcgfEfU6c2QvXUxD9xl6TUFapw8G+yeTcaTTe5F3mQlxBGNkcnwbeHhhRzdYKVZ3+jKsH359DR9NImN7DT1alfp/cilvJJwZemkAVjVrnyssVj6RPRt0dV0ejXAUVxYEEnZxQEA="}' YÂ {immediate: true, UserCertifyId: undefined, DeviceConfig: undefined, deviceConfig: undefined, DeviceToken: 'U0dfV0VCIzM3OTVkMjgyNDJhMTE2MTliYzI1Zjc4NmY4NGU1M2â€¦0OCMwZTI4OTcyMzdlZjcxYmRmMjkzMDA0YjRmNmQ4YTFiNg==',Â â€¦} undefined undefined undefined

PromiseÂ {<pending>}
```
Next call is likely signnature generator
```js
console.log(t,e,r,i,a,o);  tr(t, e, r, i, a, o);
VM65765:1 13 {AccessKeyId: 'LTAI5tSEBwYMwVKAQGpxmvTd', SignatureMethod: 'HMAC-SHA1', SignatureVersion: '1.0', Format: 'JSON', Timestamp: '2026-06-22T04:09:14Z',Â â€¦}AccessKeyId: "LTAI5tSEBwYMwVKAQGpxmvTd"Action: "VerifyCaptchaV3"CaptchaVerifyParam: "{\"sceneId\":\"didk33e0\",\"certifyId\":\"6iL4denBvY\",\"deviceToken\":\"U0dfV0VCIzM3OTVkMjgyNDJhMTE2MTliYzI1Zjc4NmY4NGU1M2Q0LWgtMTc4MjA5MzM1NzYyMC0wY2U5YzY2NTU2N2E0N2Q3OTBlNzNkNGE3MTgwNDMwNCNrUHNRYzl2NzB2YXBNdEl6YVVTWWZ4OUFndXBRSzd2RW5QOVpobmgxb1RsODBKWTVpYXFlQkhlWExIR2hDWUlmRHU2WkgrQXNQVk41SHNEWlB5MUNmOWVjVnExb1F6NGtwUkxJQkl1NXRhbVR4d0ZidERuY0ozT2FERTdqdDhtd2dIQ2QrYU5PQ0hlalJlU1A0L0ZoOTFUQzVFbWpCNzZMazQ1UzJvOHNWZmgrTHN0dWcvTW53emtBSUc1L0RBOXc5cGU0NWd1R0RKV21NSHBKTXJ5cm40VGpiTEsvcUxDa1hkcWlIdGkxbUpaRmxUd3N3MkVVcjAybVBxVGIzc3RZVS9INEQwSDVVZHpIUnovZ01jc1BmSnF3WEsyVmVVT2lNTjRYN3RXM3Nqb0wvdW4xWHpGSFh4OE5yQ1d3b093cVVJMkVPVmdKeCtBNTYvR3FETlQyTGU4WFdudFdPcDcxYjhvY0ZTQkhJNytlVUljaGJRendXN2lmcnBkNnU1WTArdEZacjYvQ1J3ajRsWUlTUkJ1WFg0VTdueTFYOGkyd3JESXdMWDZKZkQ0SHZha1c1dFdEM09QOUpTczNVelFTMjJwSXZjYkY4cHF3WHdxdFgwU3MzNC9Qa0hvRDBxL3NacFZsWDEzL1hBYTcvZ1VWUUVXRWxyTXVpZis3TS9jaEdLV0NqUXF6dnZ6WEFCaDJFY3padk5PNG9lLy8yN2lTQTk5SG5BZDRWWDlmczBOV1czUUxxd1lsR2N6R2o2NmxXVk0wYWxzbHl2YWFXMjBHcE9XVzBkdGFUVXpJZVZvVUZTUHlEZWkyTHNJaWtuS0VlWEJvM3RJOGJRdDcrSklTRHdwVlZhRVdEU1FXbFBJMXdEeG5pOUtqY3prNWZHcEgvUlpNNmlnUFU3VkZtL2xrdmFDL2lQWXVPMC9jNThnYk5JTnVqUWRvS1dZWjk3NE04WlhjOHdUYXJLb0R6MUVPSCs3L0h1VWlkRVVHQ2M3b3VXYXY4ajMvL2ZqWGdDZmFiQ01meEJFcWRoaExEZ1NXL3ozVENxQ1VNQ2pQdTNFdVRpNjRqcjBSeUJPaGl2THVSNnk5U0YydmltamdVVHovM3A2M3NzV0VSTFF4My81aEg1T1hNeENHSG5GL3p0dmdNSjRTUmNJZmlyST0jNDQ4Izg1YjgxN2M2YmM4NzNjNTlkZDQ2OWRiMzhmNTk5YjA0\",\"data\":\"JRMlgg0wDgRASAITRQpNEHEZuJsRBQZDtzsmAXNCL+49eDyEFKlLbSVgFq8aNiIiME9qeyVTKaKkiAQCrb8LIk1vCHMXJG0MOWjmeHE5GGodLkQQXSRFY+BCdntpe1A1ez2zBcgfEfU6c2QvXUxD9xl6TUFapw8G+yeTcaTTe5F3mQlxBGNkcnwbeHhhRzdYKVZ3+jKsH359DR9NImN7DT1alfp/cilvJJwZemkAVjVrnyssVj6RPRt0dV0ejXAUVxYEEnZxQEA=\"}"CertifyId: "6iL4denBvY"Format: "JSON"SceneId: "didk33e0"SignatureMethod: "HMAC-SHA1"SignatureNonce: "99b9c6e8-31af-4fc5-8580-7de85d63d7ae"SignatureVersion: "1.0"Timestamp: "2026-06-22T04:09:14Z"Version: "2023-03-05"[[Prototype]]: Object 'YSKfst7GaVkXwZYvVihJsKF9r89koz' undefined undefined undefined

'nrZkfG4wzTrgysSsOmq4HgEEwDE='
```

tr(11) also generates signature 

tr(14,e) url encodes e


### ZLib compression of trackjson; Function K is called by this code with trackjson( hash included)

Correct me if im wrong
Subject:
```js
nj = [nx.A(K, er), nx.F(O, (td || td)(10, 78))];
```
#### Lets deobfuscate
```js
nx.A = function(t,n) { return t(n)}

After resolving 
nj = [K(er),nx.F(O, td(10,78))];
```
```js
console.log(er)

3cf7c247ff063c6f65dadae1b49dae28{"TrackList":{"mc":"","tc":"","mu":"","te":"","mp":"","tmv":"","ks":"","fi":"","startTime":1782103567255},"TrackStartTime":1782103567255,"VerifyTime":1782103567286,"arg":"JSCTGEZbEAAgBQ=="} // MD5+JSON 
```
```js
td = function(t, n) {
                            return nx.j(ti, ti)(nx.w(t, 7), n)
                        }

nx.j = function(t, n) {
                                return t || n
                            } 

nx.w = function(t, n) {
                            return t - n
                        }


// Resolution

td = function(t, n) {
    return ti(t - 7, n)
}
```


```javascript
// Function ti string table

function ti(t, n) {
            var e = ["A=ZCC6i1", "4CavcUIg8/d", "Lk=/cq", "4S2djA0TjTE", "LkEZOv8wLQW", "IvWamq3P", "+UbStUZp", "Q6Zj53fOs3W", "+UsX", "tE2+QP2UU2E", "J5V94Pbh", "Vv8yVQe8oSe", "oq6UH5O1aqh", "13esP=M6G16", "Q37+Q3M", "B/fZMq", "MU7X8UM", "5EBQGJa16E5", "CU7xA=DGB28", "9SnuLknK9n", "aS5gaHBSOHO", "s1G7", "B3bAQ=e", "MC8F", "4YBbU/0psNO", "spDjBE1J62h", "OQEhOQOY", "lv67", "qmWwVUE6ovM", "ApIn5=bqB=h", "6PDNsNZp8Cq", "HH6d", "VSqZH08", "sCaysToRtUn", "8kp/8UhdcgW", "B1fBAG1cC=h", "oveX9vE0", "ovMeHSMXoAi", "OSeW", "UP7XxTfDCPq", "GE1C5PaBB2W", "qIe6HW", "53sGBEf1UYi", "Q27aQ1i", "VH5v", "3nEQ", "CpZ9+61H5=e", "sqOo3q", "VS=Koq", "UY2xU2I+x2i", "4HoKjHDTjPq", "m1po53qomBW", "G3GMq=imGqq", "xN1yMq", "30BM61fIOnq", "xTxfsJ1ZjU8", "UTpKMTBE", "GN1T811T8J8", "lvENO0iZcW", "Vg82cgVk", "sTGwMCfd", "MUZT4CxWjTV", "om=iLvW7Vkh", "Og3eoW", "og5kVknuLW", "9S1Sj/sdjNh", "366C3=745Ed", "5Bo1Cp1c8Ed", "xY1k", "LQWX", "VvOkVW", "8Top5P7g8Th", "9q5eLA3Y95", "tP7ptW", "A2BXG27vABO", "OS3N9AVbIve", "xNDT+m0XxPd", "6Ba3UBb4", "sT1Y8/1YGNq", "JmefLNZK", "mWOtanECq5i", "QTsP5n", "Cpbc56a6U16", "qZ3P", "HgebakEGLq", "56xGG/25", "4J5hMTZ7xCd", "ok8plSqT", "Ovn0M5", "8Nsgsn", "cAiF", "BEZ1a5", "GE1t63MV328", "tE2+U229AW", "lgWh", "IGo0", "J6E8QG6Vmq5", "53BIGB0y53O", "U=s9C=Z1U3M", "6/7hB/1X4/O", "LnMFcm=v", "8CGf+YM", "VnnRlQV2Lq", "jkxvt3bf+Yn", "cmiNLkpy", "+J7KtP7uM/d", "3Z6H3Z5AH5d", "OmEucQq", "mkEyOQd0JgM", "xCfNjn", "LAMqoS3/Vm3", "A1a6C5", "tNok4C7m4Jq", "3W6m3n66", "IIiUUnn83p3", "tJfT+Pxe8Un", "Bpo1Q=M", "LgWD", "oUOZaJ3gaQq", "cgqy9mMw95", "P5hV1GnLm5O", "m5qLI05qIn5", "QE0V52GZQEV", "AEZ9Bp5W", "QBI6Cq", "+Aqn9C=k+/8", "mS=RamOr", "4NsSt/Zw+YV", "Q=MN", "MC0XjU7dtTW", "cQOWOHE", "Up18ATfUCBW", "Q1bV62Zm", "sNbs8T7Ra2=", "G=BsC33", "tPbvtTe", "spxlG/xM6BE", "MT7T8n", "+CZh5Y20t5", "OHo5", "oki/V5", "IG56JWhIP5", "sB7UGnV", "M3bCx=BAA68", "Q=IMBn", "Lkebon", "amndLmdvVgV", "Jq=5InhMHWi", "tNxN", "t/xZtmx2xU=", "8QpgsU874Yq", "OPVgamapMY6", "G6o1UP0156=", "C17LtB7VQpM", "tNpdAUoXsUV", "C1a45=oPPpE", "tYxhMJoRxW", "+YBx8YGk4Cd", "lvnEOQeb", "+TIiQT7Dj/V", "UEG8CW", "Omd2", "lgWfaQnrLq", "66ffApoqBN8", "632o5505CB5", "53ZaG6f3", "9g5V9v=vaQW", "x67lG3h", "Q5O61G6cmEE", "j=2t567G4=E", "Q=aVB5", "VA3XVA3XVAq", "tJWTcPV2+Je", "tJaKjTGb45", "U3xj4=2+C65", "aH6hIkq/Vgq", "jUZD", "cAnZPmEKHgE", "oHhvLn", "C67s5=2A46d", "Qpbt", "C2Da5BxV", "lAnwoHdimSV", "tY2XMPxu", "CCGfx/fR", "L/D0GYI/8U8", "s3sJ51sU+=q", "41Bt6E1B623", "sYoZsYM", "61sLBW", "s/IN8Jxh8UW", "avqwaHeHOvq", "1qn13qV91n=", "6EbHU60PAEO", "3n5A3nncmZd", "4Jbj4Ybr", "x3ZX+CsZ4q", "MPxZtY2=tY6", "L0eZl5", "avZh+C0=OAO", "6CZ7tUBu6T3", "5p1JCBoa5SM", "8BIM8BxGtCe", "IWVcmJVMIWd", "cviy3Qhk", "xUa2xUbft3h", "xCINM5", "HZsL52aAQS=", "+P0psW", "CYGy8NZ74Ph", "A104BB7vCW", "qGO4In", "LUedLmpXLNW", "8Jpw83p=+NM", "cmiXHgqEck3", "UB7sG38k", "5=B55EBv61E", "GCfh4N1Y", "G2ej51IIBp8", "1GMC35MPHZh", "xCod4UGQ4n", "C2sxB3x+", "xJGT4PV", "G=7CQ1DaqNi", "m5ia3q", "xPBg8Ypnsq", "jYoH+NaD8/d", "VkONakeha5", "ag6=1m6YcQO", "OH=plvdgcCM", "G6pj", "I55CHqV3a5d", "VQ6ilnVfOg6", "MYakjNi", "GB0MBI2tB3n", "6YoP6=M", "MNbSt5", "B1bACpB6", "IgqYVm=0", "CE7ct1DOQBM", "4=0GxHqYLC=", "qIW53I5UJn", "9HOYon", "OCoitC0RBCn", "8J0/t5", "B3aA62GV", "tPaKxCGwt5", "oW5ravMwjgq", "JQhi9kOhVZO", "+Ps08n", "t/1dlvWR", "GBoPU2s66/8", "QNaS+TpujY=", "4UO0+A6n+SO", "+603G/sDsI5", "OQ0f", "6posQ/GsBBe", "5p7LG32o", "5=G1", "Vg5SaH6", "L55ravMwjgi", "s6aB5YBmB6V", "jUx5MUV", "QUbu4Y2Y+2n", "xTZbxn", "+EpjAPa4", "GYx2lW", "lviv", "Aq0GCBsHmW=", "MCIYjY1rs5", "M=1UBE29An", "sU0/sN0SxJE", "xT0k8YBKsT5", "j=GsB6GLG5", "VmEk", "lkOR9n", "BpDI", "4k8hov8dxk5", "oHnZ", "Q6b4t101CB=", "sN0TGNi", "+BZN+Uq"];
            return (ti = function(n, r) {
                var i = e[n -= 0];
                if (i) {
                    if (void 0 === ti.Sc) {
                        ti.Sc = !0;
                        var a = "c1f9f8dd8381fae0e6dfdc89ffd1d3fc85f7f1e186f2e5f388c3c49bfdc884dade80e3db8dc0e99fd9d6c9e8d5f4fbe2e7eac6d7f582fee4d487c2c7d8d2c5f6ca".match(/.{1,2}/g).map(function(t) {
                            return parseInt(t, 16)
                        });
                        ti.lF = function(t, n) {
                            for (var e, r, i = "", o = "", c = 0, s = 0; r = t.charAt(s++); ~r && (e = c % 4 ? 64 * e + r : r,
                            c++ % 4) && (i += String.fromCharCode(255 & e >> (-2 * c & 6) ^ n)))
                                r = a.indexOf(176 ^ r.charCodeAt(0));
                            for (var u = 0, f = i.length; u < f; u++)
                                o += "%" + ("00" + i.charCodeAt(u).toString(16)).slice(-2);
                            return decodeURIComponent(o)
                        }
                        ,
                        t = {}
                    }
                    var o = t[n];
                    return o ? i = o : (i = ti.lF(i, r),
                    t[n] = i),
                    i
                }
            }
            )(t, n)
        }


```

### NOTE:
When using function ti or wrapper of ti first use a valid xor key only then use single arguments 

#### Test and proofs
td(10,78)
'4c63f913'
ti(3,78)
'4c63f913'
ti(78)
'gqfbqfG`'
ti(3)
'4c63f913'

#### It looks like ti function focuses on first argument
```js
ti(7)
'ML_@JLdL'
ti(7,6)
'ML_@JLdL'
ti(6)
'mobile'
ti(6,7)
'mobile'
```
```js
// Value of td(10,78)
console.log(td(10,78)) 
'4c63f913'
```
```
// nx.F function
nx.F = function(t, n) {
                            return t + n
                        }
```
```js
// Resolution

nj = [K(er),O+td(10,78)];

O = td(219, 14)
'3e627e1b'

// Resolution
 nj = [K(er),'3e627e1b4c63f913'];
```
```js
console.log(nj)
[
    "eJx1jcEKwjAQRP9lzzm4aWLSQg+9iIgg2CJe05KWUCKaREVK/91AchI8vX3MsEO5LnuumGC8ZIgo+bBhapS66OUWBVugc2qYj8YHqBawA1QABEKmfWbX2e/Z7Ssds08cTaIPyoXO2NhHIWlcLLhEKVaShto/OYGLdmb8/EYUKQHlpvj9cH60++u72TXN7TTVNaxfJOs7Cg==",
    "3e627e1b4c63f913"
]
```
```bash
~ $ echo eJx1jcEKwjAQRP9lzzm4aWLSQg+9iIgg2CJe05KWUCKaREVK/91AchI8vX3MsEO5LnuumGC8ZIgo+bBhapS66OUWBVugc2qYj8YHqBawA1QABEKmfWbX2e/Z7Ssds08cTaIPyoXO2NhHIWlcLLhEKVaShto/OYGLdmb8/EYUKQHlpvj9cH60++u72TXN7TTVNaxfJOs7Cg== | base64 -d | xxd | head -1
00000000: 789c 758d c10a c230 1044 ff65 cf39 b869  x.u....0.D.e.9.i

Success: Valid zlib headers 
```

```python
>>> import zlib                                                                                                                                                          
>>> import base64
>>> b64_data = "eJx1jcEKwjAQRP9lzzm4aWLSQg+9iIgg2CJe05KWUCKaREVK/91AchI8vX3MsEO5LnuumGC8ZIgo+bBhapS66OUWBVugc2qYj8YHqBawA1QABEKmfWbX2e/Z7Ssds08cTaIPyoXO2NhHIWlcLLhEKVaSh\
to/OYGLdmb8/EYUKQHlpvj9cH60++u72TXN7TTVNaxfJOs7Cg=="                                                                                                                      ...
>>> decompressed = zlib.decompress(base64.b64decode(b64_data))                                                                                                            >>> print(decompressed.decode('utf-8'))
25e9b5a47459411185c04af8e3b86174{"TrackList":{"mc":"","tc":"","mu":"","te":"","mp":"","tmv":"","ks":"","fi":"","startTime":1782111358187},"TrackStartTime":1782111358187,"VerifyTime":1782111358212,"arg":"JRqSHXwAFAAnOg=="}

```



### This is hash calculation call
```js
er = P(0, [], F, V, ez, [er, (td && td)(77, 19)]) + er
```

#### Virtual Machine Interpreter
```javascript
// Function P
var P = function t(n, e, r, i, a, o) {
            var c, s, u, f, l, h, p, d, A;
            function v(t) {
                var n, e;
                for (e = 0,
                n = []; (e - t) * 10 + -30 < -30; e++)
                    n.push(e);
                return n
            }
            function b(t, n) {
                return t && "__proto__" != n ? t.hasOwnProperty(n) ? t : b(Object.getPrototypeOf(t), n) : null
            }
            for ((c = Object.create(a.s || {}))._ = window,
            c["*"] = a.t || this,
            a.e && (c[i[r[n + 1]]] = a.e),
            c.arguments = o,
            s = 0,
            u = void 0; (n - r.length) * 93 + 61 < 61; )
                if (10 == (A = r[n++]) ? (f = e.pop(),
                h = ((l = e.pop()) - f) * 48 + -38 <= -38,
                (p = r[n++]) && e.push(h)) : 20 == A ? (f = e.pop(),
                h = ~(~(l = e.pop()) | ~f),
                (p = r[n++]) && e.push(h)) : 28 == A ? (f = e.pop(),
                l = r[n++],
                f || (n = l)) : 39 == A ? (f = r[n++],
                l = r[n++],
                h = r[n++],
                p = r[n++],
                n = f,
                function() {
                    try {
                        ((f = t(l, e, r, i, {
                            s: c,
                            t: c["*"]
                        })) - 0) * 21 + -79 > -79 && (s = 1,
                        u = a.r ? e.pop() : f)
                    } catch (n) {
                        ((f = t(h, e, r, i, {
                            s: c,
                            t: c["*"],
                            e: n
                        })) - 0) * 50 + -57 > -57 && (s = 1,
                        u = a.r ? e.pop() : f)
                    } finally {
                        ((f = t(p, e, r, i, {
                            s: c,
                            t: c["*"]
                        })) - 0) * 30 + -22 > -22 && (s = 1,
                        u = a.r ? e.pop() : f)
                    }
                }()) : 44 == A ? e.push({}) : 48 == A ? (f = e.pop(),
                h = (l = e.pop()) >> f,
                (p = r[n++]) && e.push(h)) : 37 == A ? (f = r[n++],
                (l = t(n, e, r, i, {
                    s: c,
                    t: c["*"]
                })) ? 1 === l ? s = 1 : (s = 1,
                u = a.r ? e.pop() : l) : n = f) : 55 == A ? (f = r[n++],
                l = [],
                v(f).forEach(function() {
                    l.unshift(e.pop())
                }),
                r[n++] && e.push(l)) : 52 == A ? (l = !(f = e.pop()),
                r[n++] && e.push(l)) : 9 == A ? (f = e.pop(),
                h = ((l = e.pop()) - f) * 20 + 93 < 93,
                (p = r[n++]) && e.push(h)) : 50 == A ? (f = r[n++],
                l = e.pop(),
                h = e.pop(),
                o = [],
                v(f).forEach(function() {
                    o.unshift(e.pop())
                }),
                p = null === h ? l.apply(c, o) : h[l].apply(h, o),
                r[n++] && e.push(p)) : 19 == A ? function() {
                    throw e.pop()
                }() : 27 == A ? (f = e.pop(),
                l = e.pop(),
                h = e.pop(),
                (p = l === c && b(l, f) || l)[f] = h,
                (d = r[n++]) && e.push(h)) : 7 == A ? (f = e.pop(),
                l = r[n++],
                h = r[n++],
                p = function(n) {
                    return function() {
                        return t(n, e, r, i, {
                            s: c,
                            t: this,
                            r: 1
                        }, arguments)
                    }
                }(l),
                h ? e.push(p) : c[f] = p) : 0 == A ? (f = e.pop(),
                h = (l = e.pop()) >>> f,
                (p = r[n++]) && e.push(h)) : 47 == A ? (f = e.pop(),
                l = e.pop(),
                (h = e.pop())[l] = f,
                r[n++] && e.push(h)) : 40 == A ? (f = e.pop(),
                h = -95 * (l = e.pop()) / (-95 * f),
                (p = r[n++]) && e.push(h)) : 35 == A ? (l = typeof (f = e.pop()),
                r[n++] && e.push(l)) : 32 == A ? (f = e.pop(),
                l = r[n++],
                f && (n = l)) : 38 == A ? (f = e.pop(),
                h = ~(~(l = e.pop()) & ~f),
                (p = r[n++]) && e.push(h)) : 30 == A ? e.push(c) : 21 == A ? c[l = i[f = r[n++]]] = void 0 : 5 == A ? (l = ~(f = e.pop()),
                r[n++] && e.push(l)) : 3 == A ? (f = r[n++],
                l = e.pop(),
                p = b(h = e.pop(), l) || h,
                d = f ? --p[l] : p[l]--,
                r[n++] && e.push(d)) : 56 == A ? (f = e.pop(),
                l = e.pop(),
                r[n++] && e.push(l[f])) : 34 == A ? e.push(window) : 26 == A ? (f = e.pop(),
                h = ((l = e.pop()) + 12) * f - 12 * f,
                (p = r[n++]) && e.push(h)) : 23 == A ? (f = e.pop(),
                h = (l = e.pop()) & ~f | ~l & f,
                (p = r[n++]) && e.push(h)) : 45 == A ? (f = e.pop(),
                l = e.pop(),
                h = delete l[f],
                r[n++] && e.push(h)) : 49 == A ? (f = r[n++],
                o || (o = [].concat(e)),
                v(f + 1).forEach(function(t) {
                    c[e.pop()] = o[-((-51 * f - -51 * t) / 51)]
                })) : 46 == A ? (f = e.pop(),
                h = (l = e.pop())instanceof f,
                (p = r[n++]) && e.push(h)) : 8 == A ? (f = r[n++],
                l = e.pop(),
                p = b(h = e.pop(), l) || h,
                d = f ? ++p[l] : p[l]++,
                r[n++] && e.push(d)) : 36 == A ? (f = e.pop(),
                h = (l = e.pop()) !== f,
                (p = r[n++]) && e.push(h)) : 33 == A ? (f = i[r[n++]],
                h = (l = c)[f],
                r[n++] && e.push(h)) : 2 == A ? (f = e.pop(),
                h = (l = e.pop()) << f,
                (p = r[n++]) && e.push(h)) : 4 == A || (41 == A ? (f = e.pop(),
                l = e.pop(),
                h = r[n++],
                o = [],
                v(h).forEach(function() {
                    o.unshift(e.pop())
                }),
                p = new (v(h).reduce(function(t, n) {
                    return t.bind(l, o[n])
                }, l[f])),
                r[n++] && e.push(p)) : 11 == A ? s = 1 : 14 == A ? (f = e.pop(),
                l = e.pop(),
                h = r[n++],
                -1 == l && (l = h),
                n = l = f[l]) : 15 == A ? (f = e.pop(),
                h = (l = e.pop()) === f,
                (p = r[n++]) && e.push(h)) : 12 == A ? (f = e.pop(),
                h = ((l = e.pop()) - f) * 92 + 67 >= 67,
                (p = r[n++]) && e.push(h)) : 54 == A ? (f = e.pop(),
                l = void 0,
                r[n++] && e.push(l)) : 53 == A ? (f = r[n++],
                l = e.pop(),
                v(f).forEach(function() {
                    e.pop()
                }),
                r[n++] && e.push(l)) : 24 == A ? (f = e.pop(),
                h = ((l = e.pop()) - f) * 36 + 69 > 69,
                (p = r[n++]) && e.push(h)) : 42 == A ? (f = e.pop(),
                h = (l = e.pop())in f,
                (p = r[n++]) && e.push(h)) : 18 == A ? (f = e.pop(),
                h = (l = e.pop()) + f,
                (p = r[n++]) && e.push(h)) : 16 == A ? (f = r[n++],
                l = "",
                v(f).forEach(function() {
                    l += e.pop()
                }),
                r[n++] && e.push(l)) : 22 == A ? (f = e.pop(),
                h = -((-70 * (l = e.pop()) - -70 * f) / 70),
                (p = r[n++]) && e.push(h)) : 25 == A ? (s = 1,
                u = a.r ? e.pop() : r[n++]) : 51 == A ? (f = e.pop(),
                h = (l = e.pop()) != f,
                (p = r[n++]) && e.push(h)) : 31 == A ? n = f = r[n++] : 29 == A ? (f = e.pop(),
                h = (l = e.pop()) == f,
                (p = r[n++]) && e.push(h)) : 17 == A ? (f = e.pop(),
                h = -(-(l = e.pop()) % f),
                (p = r[n++]) && e.push(h)) : 6 == A ? (f = r[n++],
                void 0 === (l = t(n, e, r, i, {
                    s: c,
                    t: c["*"]
                })) ? n = f : (s = 1,
                u = a.r ? e.pop() : l)) : 43 == A ? (f = r[n++],
                e.push(i[f])) : 13 == A ? e.push(e[-((-41 * e.length - -41) / 41)]) : 1 == A ? e.pop() : function() {
                    throw Error()
                }()),
                s)
                    return u
        }

```
#### value of ez
```js
console.log(ez)
{"r": 1}
```
#### Value of td(77,19)
```js
ti(70,19)
'0000'
```

#### Value of V
```js
console.log(V)
[    "o",    "r",    "a",    "m",    "C",    "e",    "f",    "i",    "p",    "h",    "d",    "",    0,    147,    "n",    "fromCharCode",    252,    241,    249,    246,    240,    231,    "map",    "join",    202,    172,    191,    164,    169,    190,    163,    165,    87,    56,    53,    61,    50,    52,    35,    156,    250,    233,    242,    255,    232,    245,    243,    108,    24,    3,    63,    30,    5,    2,    11,    206,    167,    160,    170,    171,    182,    129,    168,    124,    12,    14,    19,    8,    25,    33,    84,    79,    69,    68,    71,    72,    "_",    55,    64,    94,    89,    83,    88,    1,    22,    97,    127,    120,    114,    121,    46,    117,    65,    76,    75,    77,    90,    74,    115,    54,    95,    82,    110,    29,    36,    166,    152,    159,    149,    158,    134,    78,    66,    85,    119,    20,    26,    18,    44,    67,    70,    73,    104,    113,    162,    205,    192,    200,    199,    193,    214,    130,    234,    239,    238,    230,    215,    207,    204,    151,    211,    248,    244,    226,    227,    181,    146,    132,    148,    133,    91,    93,    86,    201,    253,    247,    178,    221,    150,    153,    136,    220,    184,    137,    145,    161,    157,    131,    179,    154,    142,    135,    28,    125,    106,    123,    155,    128,    32,    45,    37,    42,    59,    111,    57,    38,    40,    216,    185,    174,    177,    183,    141,    195,    236,    251,    228,    212,    210,    218,    173,    188,    138,    139,    143,    144,    175,    187,    180,    4,    7,    34,    194,    223,    229,    15,    49,    58,    96,    10,    101,    16,    41,    107,    48,    9,    31,    27,    21,    6,    109,    140,    224,    225,    118,    122,    196,    197,    203,    186,    189,    176,    13,    17,    43,    99,    23,    47,    103,    112,    237,    80,    105,    102,    98,    213,    126,    81,    100,    92,    116,    51,    62,    219,    209,    254,    208,    39,    222,    198,    217,    235,    false,    60,    "Boolean",    "Number",    "String",    "j",    "t",    "u",    "arguments",    "s",    "random",    256,    "floor",    "push",    "length",    "0",    "toString",    "y",    "charCodeAt"]
```


### Lets deobfuscate and decompile VM
```
P(0, [], F, V, ez, [er, (td && td)(77, 19)]) + er

td && td = td

// Resolution
P(0, [], F, V, ez, [er,td(77, 19)]) + er

td(77,19) = '0000'

// resolution
P(0, [], F, V, ez, [er,'0000']) + er
```

In order to find how hash is generated we must disassemble the virtual machine first 
Using `disasm.py` we can get disassembly of the virtual machine after rversing xor encoded strings we will get 
full disassembly of Virtual Machine using which we can reconstruct it in javascript

## How i reverse engineered hash calculation

1. Disassebmle Virtual Machine with appropriate bytecode and other arguments: `disasm.py`
2. Finding xor encoded function/strings using `analyze.py`
3. Decoding all the XOR-encoded strings that get assigned to variable names using `decode_strings.py`
4. Extracting all encoded strings by finding BUILD_ARRAY + map(fn) + join patterns using `full_decode.py`
5. Extracting all the encoded byte arrays and decode them to strings using `extract_strings.py`
6. Analysing critical parts disassembly that does hashing using bash commands like `grep -n "" disasm.txt | tail -300 | head -200`, `tail -100 disasm.txt`.`grep -n "" disasm.txt | grep -E "^(7[6-9][0-9]{2}|8[0-2][0-9]{2}):" | head -100`. `grep -n "" disasm.txt | grep -E "^(7[4-6][0-9]{2}):" | head -100` etc.
7. Decoding the strings around pc 15396 and 15464 that appear to be key strings in the algorithm using:
```bash
python3 << 'PYEOF'
V = [
    "o","r","a","m","C","e","f","i","p","h","d","",0,147,"n","fromCharCode",
    252,241,249,246,240,231,"map","join",202,172,191,164,169,190,163,165,
    87,56,53,61,50,52,35,156,250,233,242,255,232,245,243,108,24,3,63,30,
    5,2,11,206,167,160,170,171,182,129,168,124,12,14,19,8,25,33,84,79,69,
    68,71,72,"_",55,64,94,89,83,88,1,22,97,127,120,114,121,46,117,65,76,
    75,77,90,74,115,54,95,82,110,29,36,166,152,159,149,158,134,78,66,85,
    119,20,26,18,44,67,70,73,104,113,162,205,192,200,199,193,214,130,234,
    239,238,230,215,207,204,151,211,248,244,226,227,181,146,132,148,133,
    91,93,86,201,253,247,178,221,150,153,136,220,184,137,145,161,157,131,
    179,154,142,135,28,125,106,123,155,128,32,45,37,42,59,111,57,38,40,
    216,185,174,177,183,141,195,236,251,228,212,210,218,173,188,138,139,
    143,144,175,187,180,4,7,34,194,223,229,15,49,58,96,10,101,16,41,107,
    48,9,31,27,21,6,109,140,224,225,118,122,196,197,203,186,189,176,13,17,
    43,99,23,47,103,112,237,80,105,102,98,213,126,81,100,92,116,51,62,219,
    209,254,208,39,222,198,217,235,False,60,"Boolean","Number","String",
    "j","t","u","arguments","s","random",256,"floor","push","length","0",
    "toString","y","charCodeAt"
]

# String at pc=15396, XOR key=75 (V[94]=75)
# nums: V[241]=176, V[105]=166, V[240]=189, V[168]=179
nums = [V[241], V[105], V[240], V[168]]
key = 75
s = ''.join(chr(n ^ key) for n in nums)
print(f"String at pc=15396 (xor=75): {repr(s)}")

# String at pc=15464, XOR key=75 but let's check the fn at 15440
# fn at 15440 has PUSH_CONST V[94]=75
nums2 = [V[90], V[180], V[186], V[104], V[247], V[90], V[51], V[68], 
         V[53], V[67], V[104], V[185], V[182], V[104], V[180], V[90], V[180], V[50]]
key2 = 75  # XOR key from fn at 15440
s2 = ''.join(chr(int(n) ^ key2) for n in nums2)
print(f"String at pc=15464 (xor=75): {repr(s2)}")

# String at 15560, fn at 15536 has PUSH_CONST V[222]=41
nums3 = [V[95], V[93], V[97], V[120], V[95], V[93], V[63], V[175], 
         V[218], V[174], V[120], V[73], V[80], V[120], V[74], V[93], V[74], V[151]]
key3 = 41
s3 = ''.join(chr(int(n) ^ key3) for n in nums3)
print(f"String at pc=15560 (xor=41): {repr(s3)}")

# What's V[289]? -> 'y'  
print(f"\nV[281]={repr(V[281])}")  # 's'
print(f"V[282]={repr(V[282])}")  # 'random'
print(f"V[283]={repr(V[283])}")  # 256
print(f"V[284]={repr(V[284])}")  # 'floor'
print(f"V[285]={repr(V[285])}")  # 'push'
print(f"V[286]={repr(V[286])}")  # 'length'
print(f"V[287]={repr(V[287])}")  # '0'
print(f"V[288]={repr(V[288])}")  # 'toString'
print(f"V[289]={repr(V[289])}")  # 'y'
print(f"V[290]={repr(V[290])}")  # 'charCodeAt'

PYEOF


Output

String at pc=15396 (xor=75): 'ûíöø'
String at pc=15464 (xor=75): 'encodeURIComponent'
String at pc=15560 (xor=41): 'decodeURIComponent'

V[281]='s'
V[282]='random'
V[283]=256
V[284]='floor'
V[285]='push'
V[286]='length'
V[287]='0'
V[288]='toString'
V[289]='y'
V[290]='charCodeAt'
```
8. Looking at the section around PC 15760-15870 to understand the main hash structure using `grep -n "" disasm.txt | sed -n '7550,7650p'`
9. Decoding the remaining key strings - especially at 15728 and find what 'h' and 'd' and 'u' are using
```bash
python3 << 'PYEOF'
V = [
    "o","r","a","m","C","e","f","i","p","h","d","",0,147,"n","fromCharCode",
    252,241,249,246,240,231,"map","join",202,172,191,164,169,190,163,165,
    87,56,53,61,50,52,35,156,250,233,242,255,232,245,243,108,24,3,63,30,
    5,2,11,206,167,160,170,171,182,129,168,124,12,14,19,8,25,33,84,79,69,
    68,71,72,"_",55,64,94,89,83,88,1,22,97,127,120,114,121,46,117,65,76,
    75,77,90,74,115,54,95,82,110,29,36,166,152,159,149,158,134,78,66,85,
    119,20,26,18,44,67,70,73,104,113,162,205,192,200,199,193,214,130,234,
    239,238,230,215,207,204,151,211,248,244,226,227,181,146,132,148,133,
    91,93,86,201,253,247,178,221,150,153,136,220,184,137,145,161,157,131,
    179,154,142,135,28,125,106,123,155,128,32,45,37,42,59,111,57,38,40,
    216,185,174,177,183,141,195,236,251,228,212,210,218,173,188,138,139,
    143,144,175,187,180,4,7,34,194,223,229,15,49,58,96,10,101,16,41,107,
    48,9,31,27,21,6,109,140,224,225,118,122,196,197,203,186,189,176,13,17,
    43,99,23,47,103,112,237,80,105,102,98,213,126,81,100,92,116,51,62,219,
    209,254,208,39,222,198,217,235,False,60,"Boolean","Number","String",
    "j","t","u","arguments","s","random",256,"floor","push","length","0",
    "toString","y","charCodeAt"
]

# String at 15728 with xor key = V[151] = 93
# nums: V[186]=40, V[261]=51, V[33]=56, V[90]=46, V[262]=62, V[273]=60, V[179]=45, V[33]=56
nums = [V[186], V[261], V[33], V[90], V[262], V[273], V[179], V[33]]
key = V[151]  # = 93
print(f"V[151] = {key}")
s = ''.join(chr(int(n) ^ int(key)) for n in nums)
print(f"String at 15728 (xor=93) -> 'u' = {repr(s)}")

# What's at string decoded earlier as h?
# At pc=15396: nums=[V[241]=176, V[105]=166, V[240]=189, V[168]=179] xor=75
# 176^75=235(?), let me recompute
nums_h = [V[241], V[105], V[240], V[168]]
key_h = 75
s_h = ''.join(chr(int(n) ^ key_h) for n in nums_h)
print(f"String for 'h' = {repr(s_h)}")

# Let me trace back further to find what 'C' function is, what 'h' value comes from
# The string result gets accessed as window[str] (._ = window)
# e['atob'] -> the fn arg[0] which is the JSON input? 
# No wait...

# Let me look at the variables more carefully
# V[279]='u', V[280]='arguments', V[281]='s'

# Argument 0 = o = the JSON input string  
# Argument 1 = r = td = '0000'

# 'o' = input string (JSON)
# 'r' = td = '0000'
# 'a' = [] (empty array, will be filled with 16 random bytes = AES key candidate?)
# 'C' = some function called on the scope (c.C())
# 'm' = result of C()
# 'f' = length of something (a lookup table size?)

# Let me find what 'f', 'i', 'e', 'p' are set to before the final section
# Looking at what builds them...

# From the disassembly list [6]-[7]:
# e = 'prototype'
# f = 'undefined'
# i = 'p'  (V[8]='p')
# p = 'h'  (V[9]='h')

# But these get overwritten. Let me look at the critical section more carefully
# From the extracted strings:
# [171] = 'Math'
# [172] = 'atob'  (V[336]=4 -> wait that doesn't match)

# Let me check string #171 - 'Math'
# pc=15268, xor=95, nums=[...] -> 'Math'
# Then it's accessed as window['Math'] 

# For 'atob': pc=15336, xor=6, -> 'atob'

# So the function accesses window.Math and window.atob etc.
# The 'C' function is some transformation
# Let me look at what h is - appears to be parseInt  based on the 'atob' access pattern

# Actually let me check what's at h - it comes from:
# PUSH_SCOPE V[76]='_'  -> window
# [XORED STRING] GET_PROP -> window[decoded_string]
# h = result

# At string [133]: pc=12204, decoded='p' -> but that's innerText looking at the assignment

# Let me check V[279]='u' and see what u gets set to
# From disasm line 7585-7586: v[279]='u' SET_PROP -> c.u = result of get_prop(window, XORED_STR)
# And the xored string at 15728 with xor=93 -> let me recalculate

print("\nRechecking special variable assignments:")

# 'h' at pc 15396: get_prop from window
nums_h2 = [V[241], V[105], V[240], V[168]]
# fn at ~15370 uses xor key V[94]=75 (from the pattern before)
# But wait - at 15396 we do: BUILD_ARRAY + map(fn) + join + GET_PROP + SET_PROP h
# The fn was the one at entry 15375 with xor=75
key_h2 = 75
s_h2 = ''.join(chr(int(n) ^ key_h2) for n in nums_h2)
print(f"h = window[{repr(s_h2)}]")  # parseInt!

# 'd' at pc 15500: window[decodeURIComponent_string] (already decoded as 'decodeURIComponent')
# But wait - the access is _.d = window['encodeURIComponent']? 
# No - let me re-read:
# c._ = window, then GET_PROP gives window['encodeURIComponent']
# That gets stored in 'd'
print(f"d = window['encodeURIComponent'] (the fn)")

# 'j' at pc 15560: window['decodeURIComponent']
print(f"j = window['decodeURIComponent'] (the fn)")

# 'u' at pc 15728 - string decoded above  
print(f"u = window[{repr(s)}]")  # let me check

# Looking at the scope: c._ = window
# So _.u means window.u? No, 'u' is a JS VM variable name (V[278]='u')
# The value of u = result of window['encodeURIComponent'](something)?
# No - at 15764: c.u = window['_'][decoded_string_at_15728]
# c._ = window so window[string]

s_u = ''.join(chr(int(n) ^ 93) for n in nums)
print(f"u = window[{repr(s_u)}]")

PYEOF
Output

V[151] = 93
String at 15728 (xor=93) -> 'u' = 'unescape'
String for 'h' = 'ûíöø'

Rechecking special variable assignments:
h = window['ûíöø']
d = window['encodeURIComponent'] (the fn)
j = window['decodeURIComponent'] (the fn)
u = window['unescape']
u = window['unescape']
```

10. Tracing back what h string should decode to - re-examine the XOR function used for h using:
```bash

python3 << 'PYEOF'
import json

with open('F_bytecode.txt') as f:
    F = json.load(f)

V = [
    "o","r","a","m","C","e","f","i","p","h","d","",0,147,"n","fromCharCode",
    252,241,249,246,240,231,"map","join",202,172,191,164,169,190,163,165,
    87,56,53,61,50,52,35,156,250,233,242,255,232,245,243,108,24,3,63,30,
    5,2,11,206,167,160,170,171,182,129,168,124,12,14,19,8,25,33,84,79,69,
    68,71,72,"_",55,64,94,89,83,88,1,22,97,127,120,114,121,46,117,65,76,
    75,77,90,74,115,54,95,82,110,29,36,166,152,159,149,158,134,78,66,85,
    119,20,26,18,44,67,70,73,104,113,162,205,192,200,199,193,214,130,234,
    239,238,230,215,207,204,151,211,248,244,226,227,181,146,132,148,133,
    91,93,86,201,253,247,178,221,150,153,136,220,184,137,145,161,157,131,
    179,154,142,135,28,125,106,123,155,128,32,45,37,42,59,111,57,38,40,
    216,185,174,177,183,141,195,236,251,228,212,210,218,173,188,138,139,
    143,144,175,187,180,4,7,34,194,223,229,15,49,58,96,10,101,16,41,107,
    48,9,31,27,21,6,109,140,224,225,118,122,196,197,203,186,189,176,13,17,
    43,99,23,47,103,112,237,80,105,102,98,213,126,81,100,92,116,51,62,219,
    209,254,208,39,222,198,217,235,False,60,"Boolean","Number","String",
    "j","t","u","arguments","s","random",256,"floor","push","length","0",
    "toString","y","charCodeAt"
]

# Let me find the exact bytes around pc 15380-15422
# Looking for the MAKE_FUNC and BUILD_ARRAY leading to h
r = F

# Scan around index 15380 in F (note F indices = bytecode positions, but we need array indices)
# The disasm showed pc=15396 as BUILD_ARRAY count=4 (push)
# That means F[15396] = 55, F[15397] = 4, F[15398] = 1

# Let me verify
print(f"F[15396] = {F[15396]} (expect 55=BUILD_ARRAY)")
print(f"F[15397] = {F[15397]} (expect 4=count)")

# The MAKE_FUNC for the mapper function of h - find it
# Scanning for the MAKE_FUNC before F[15396]
# Looking backwards from 15396 for a MAKE_FUNC (7)
search_range = range(15380, 15400)
for i in search_range:
    print(f"F[{i}] = {F[i]}")

print("...")
# The string bytes for h:
print(f"\nNums: V[241]={V[241]}, V[105]={V[105]}, V[240]={V[240]}, V[168]={V[168]}")

# Now find the MAKE_FUNC before this - go back
# The structure is: PUSH_CONST '' (43,11) PUSH_CONST 0 (43,12) MAKE_FUNC (7) entry push
# Then jump to after the function body
# Let's find MAKE_FUNC between 15370 and 15395
for i in range(15365, 15396):
    if F[i] == 7:
        print(f"MAKE_FUNC at F[{i}], entry={F[i+1]}, push={F[i+2]}")
        entry = F[i+1]
        # Check entry
        print(f"  Entry bytecode F[{entry}]={F[entry]}, F[{entry+1}]={F[entry+1]}")
        if F[entry] == 43 and F[entry+1] == 0:  # PUSH_CONST 'o'
            n = entry + 2
            if F[n] == 49:  # BIND_ARGS
                nargs = F[n+1]
                n += 2
                if F[n] == 33 and F[n+1] == 0:  # LOAD_VAR 'o'
                    n += 3
                    if F[n] == 43:  # PUSH_CONST KEY
                        key_idx = F[n+1]
                        print(f"  XOR key = V[{key_idx}] = {V[key_idx]}")

PYEOF
Output

F[15396] = 43 (expect 55=BUILD_ARRAY)
F[15397] = 241 (expect 4=count)
F[15380] = 198
F[15381] = 23
F[15382] = 1
F[15383] = 30
F[15384] = 43
F[15385] = 14
F[15386] = 56
F[15387] = 1
F[15388] = 43
F[15389] = 15
F[15390] = 50
F[15391] = 1
F[15392] = 1
F[15393] = 25
F[15394] = 2
F[15395] = 11
F[15396] = 43
F[15397] = 241
F[15398] = 43
F[15399] = 105
...

Nums: V[241]=176, V[105]=166, V[240]=189, V[168]=179
MAKE_FUNC at F[15367], entry=15372, push=1
  Entry bytecode F[15372]=43, F[15373]=0
  XOR key = V[198] = 210
```

11. Recalculate h string with correct XOR key 210 using:
```bash
python3 << 'PYEOF'
import json

with open('F_bytecode.txt') as f:
    F = json.load(f)

V = [
    "o","r","a","m","C","e","f","i","p","h","d","",0,147,"n","fromCharCode",
    252,241,249,246,240,231,"map","join",202,172,191,164,169,190,163,165,
    87,56,53,61,50,52,35,156,250,233,242,255,232,245,243,108,24,3,63,30,
    5,2,11,206,167,160,170,171,182,129,168,124,12,14,19,8,25,33,84,79,69,
    68,71,72,"_",55,64,94,89,83,88,1,22,97,127,120,114,121,46,117,65,76,
    75,77,90,74,115,54,95,82,110,29,36,166,152,159,149,158,134,78,221,
    150,153,136,220,184,137,145,161,157,131,
    179,154,142,135,28,125,106,123,155,128,32,45,37,42,59,111,57,38,40,
    216,185,174,177,183,141,195,236,251,228,212,210,218,173,188,138,139,
    143,144,175,187,180,4,7,34,194,223,229,15,49,58,96,10,101,16,41,107,
    48,9,31,27,21,6,109,140,224,225,118,122,196,197,203,186,189,176,13,17,
    43,99,23,47,103,112,237,80,105,102,98,213,126,81,100,92,116,51,62,219,
    209,254,208,39,222,198,217,235,False,60,"Boolean","Number","String",
    "j","t","u","arguments","s","random",256,"floor","push","length","0",
    "toString","y","charCodeAt"
]

# The BUILD_ARRAY for 'h' is at pc 15404 (F index 15404)
# F[15404]=55 = BUILD_ARRAY? Let me verify
print(f"F[15404] = {F[15404]}")
print(f"F[15405] = {F[15405]}")  # count

# Find exact positions
for i in range(15396, 15430):
    print(f"F[{i}] = {F[i]}")

PYEOF
Output

F[15404] = 55
F[15405] = 4
F[15396] = 43
F[15397] = 241
F[15398] = 43
F[15399] = 105
F[15400] = 43
F[15401] = 240
F[15402] = 43
F[15403] = 168
F[15404] = 55
F[15405] = 4
F[15406] = 1
F[15407] = 43
F[15408] = 22
F[15409] = 50
F[15410] = 1
F[15411] = 1
F[15412] = 43
F[15413] = 23
F[15414] = 50
F[15415] = 1
F[15416] = 1
F[15417] = 56
F[15418] = 1
F[15419] = 30
F[15420] = 43
F[15421] = 9
F[15422] = 27
F[15423] = 0
F[15424] = 21
F[15425] = 10
F[15426] = 30
F[15427] = 43
F[15428] = 76
F[15429] = 56
```

12. Tracing the actual hash algorithm - look at what happens with 'e', 'f', 'C' which are the core of the hash using:
```bash
python3 << 'PYEOF'
import json

with open('F_bytecode.txt') as f:
    F = json.load(f)

V = [
    "o","r","a","m","C","e","f","i","p","h","d","",0,147,"n","fromCharCode",
    252,241,249,246,240,231,"map","join",202,172,191,164,169,190,163,165,
    87,56,53,61,50,52,35,156,250,233,242,255,232,245,243,108,24,3,63,30,
    5,2,11,206,167,160,170,171,182,129,168,124,12,14,19,8,25,33,84,79,69,
    68,71,72,"_",55,64,94,89,83,88,1,22,97,127,120,114,121,46,117,65,76,
    75,77,90,74,115,54,95,82,110,29,36,166,152,159,149,158,134,78,66,85,
    119,20,26,18,44,67,70,73,104,113,162,205,192,200,199,193,214,130,234,
    239,238,230,215,207,204,151,211,248,244,226,227,181,146,132,148,133,
    91,93,86,201,253,247,178,221,150,153,136,220,184,137,145,161,157,131,
    179,154,142,135,28,125,106,123,155,128,32,45,37,42,59,111,57,38,40,
    216,185,174,177,183,141,195,236,251,228,212,210,218,173,188,138,139,
    143,144,175,187,180,4,7,34,194,223,229,15,49,58,96,10,101,16,41,107,
    48,9,31,27,21,6,109,140,224,225,118,122,196,197,203,186,189,176,13,17,
    43,99,23,47,103,112,237,80,105,102,98,213,126,81,100,92,116,51,62,219,
    209,254,208,39,222,198,217,235,False,60,"Boolean","Number","String",
    "j","t","u","arguments","s","random",256,"floor","push","length","0",
    "toString","y","charCodeAt"
]

# Now let's look at the main function that's being computed
# The initial part sets up:
# - f = 256 (lookup table size) -> from the final section
# - e = a lookup/state table array of 256 bytes (the S-box for RC4/custom cipher)
# - C = the key scheduling function
# - i = Math
# - h = btoa
# - etc.

# Let me find where 'f' and 'e' get their values in the final section
# Looking at disasm around 15780+

# The key insight from the loop structures:
# Starting at PC 15814:
#   s = 0; while s < 16: push random byte to 'a', s++
# So 'a' = array of 16 random bytes? No wait - it checks 'C()'
# Actually: C() is called, result in m
# if m: 
#   m = m[a[a.length-1]]  -- last element of a? Indexed by something
#   Then proceeds with the m value as the key
# If not m (i.e., if __ALIYUN_CRYPT is not defined):
#   Generate 16 random bytes into a
#   a = a.map(fn).join('') -> convert each byte to 2-hex-digit string
# Then jump to final section

# Let me read PC 15871-15980 more carefully from the raw bytecodes

print("Bytecodes at PC 15871-15980:")
for i in range(15871, 15982):
    print(f"  F[{i}] = {F[i]}")

PYEOF
Output

Bytecodes at PC 15871-15980:
  F[15871] = 30
  F[15872] = 43
  F[15873] = 4
  F[15874] = 50
  F[15875] = 0
  F[15876] = 1
  F[15877] = 30
  F[15878] = 43
  F[15879] = 3
  F[15880] = 27
  F[15881] = 0
  F[15882] = 33
  F[15883] = 3
  F[15884] = 1
  F[15885] = 28
  F[15886] = 15980
  F[15887] = 6
  F[15888] = 15978
  F[15889] = 33
  F[15890] = 3
  F[15891] = 1
  F[15892] = 30
  F[15893] = 43
  F[15894] = 2
  F[15895] = 56
  F[15896] = 1
  F[15897] = 30
  F[15898] = 43
  F[15899] = 2
  F[15900] = 56
  F[15901] = 1
  F[15902] = 43
  F[15903] = 286
  F[15904] = 56
  F[15905] = 1
  F[15906] = 43
  F[15907] = 83
  F[15908] = 22
  F[15909] = 1
  F[15910] = 27
  F[15911] = 0
  F[15912] = 43
  F[15913] = 11
  F[15914] = 43
  F[15915] = 12
  F[15916] = 7
  F[15917] = 15921
  F[15918] = 1
  F[15919] = 31
  F[15920] = 15957
  F[15921] = 43
  F[15922] = 0
  F[15923] = 49
  F[15924] = 0
  F[15925] = 33
  F[15926] = 0
  F[15927] = 1
  F[15928] = 43
  F[15929] = 221
  F[15930] = 9
  F[15931] = 1
  F[15932] = 28
  F[15933] = 15938
  F[15934] = 43
  F[15935] = 287
  F[15936] = 31
  F[15937] = 15940
  F[15938] = 43
  F[15939] = 11
  F[15940] = 43
  F[15941] = 221
  F[15942] = 30
  F[15943] = 43
  F[15944] = 0
  F[15945] = 56
  F[15946] = 1
  F[15947] = 43
  F[15948] = 288
  F[15949] = 50
  F[15950] = 1
  F[15951] = 1
  F[15952] = 18
  F[15953] = 1
  F[15954] = 25
  F[15955] = 2
  F[15956] = 11
  F[15957] = 30
  F[15958] = 43
  F[15959] = 2
  F[15960] = 56
  F[15961] = 1
  F[15962] = 43
  F[15963] = 22
  F[15964] = 50
  F[15965] = 1
  F[15966] = 1
  F[15967] = 43
  F[15968] = 23
  F[15969] = 50
  F[15970] = 1
  F[15971] = 1
  F[15972] = 30
  F[15973] = 43
  F[15974] = 2
  F[15975] = 27
  F[15976] = 0
  F[15977] = 11
  F[15978] = 31
  F[15979] = 16833
  F[15980] = 6
  F[15981] = 16833
```

13. Tracing the actual VM hash algorithm - let's look at what 'y' function does at PC 15982 using:
```bash
python3 << 'PYEOF'
import json

with open('F_bytecode.txt') as f:
    F = json.load(f)

V = [
    "o","r","a","m","C","e","f","i","p","h","d","",0,147,"n","fromCharCode",
    252,241,249,246,240,231,"map","join",202,172,191,164,169,190,163,165,
    87,56,53,61,50,52,35,156,250,233,242,255,232,245,243,108,24,3,63,30,
    5,2,11,206,167,160,170,171,182,129,168,124,12,14,19,8,25,33,84,79,69,
    68,71,72,"_",55,64,94,89,83,88,1,22,97,127,120,114,121,46,117,65,76,
    75,77,90,74,115,54,95,82,110,29,36,166,152,159,149,158,134,78,66,85,
    119,20,26,18,44,67,70,73,104,113,162,205,192,200,199,193,214,130,234,
    239,238,230,215,207,204,151,211,248,244,226,227,181,146,132,148,133,
    91,93,86,201,253,247,178,221,150,153,136,220,184,137,145,161,157,131,
    179,154,142,135,28,125,106,123,155,128,32,45,37,42,59,111,57,38,40,
    216,185,174,177,183,141,195,236,251,228,212,210,218,173,188,138,139,
    143,144,175,187,180,4,7,34,194,223,229,15,49,58,96,10,101,16,41,107,
    48,9,31,27,21,6,109,140,224,225,118,122,196,197,203,186,189,176,13,17,
    43,99,23,47,103,112,237,80,105,102,98,213,126,81,100,92,116,51,62,219,
    209,254,208,39,222,198,217,235,False,60,"Boolean","Number","String",
    "j","t","u","arguments","s","random",256,"floor","push","length","0",
    "toString","y","charCodeAt"
]

# Let's trace the bytecodes from 15982 onward (the main y function)
# Entry is at 15989 (from F[15986]=15989)
print("F[15982:16020]:")
for i in range(15982, 16020):
    print(f"  F[{i}] = {F[i]}", end="")
    if F[i] == 43 and i+1 < len(F):
        vi = F[i+1]
        if vi < len(V):
            print(f"  (PUSH_CONST V[{vi}]={repr(V[vi])})", end="")
    print()

print("\n\nTracing y function body (PC 15989+):")
# From the disasm we saw:
# 15989: PUSH_CONST V[0]='o', PUSH_CONST V[1]='r', BIND_ARGS 1
# Then: VAR_DECL o, r, a, m, C, e, f, i (local scope vars)
# Then: ...

# Let me decode the exact sequence at 15989
for i in range(15989, 16075):
    print(f"  F[{i}] = {F[i]}", end="")
    if F[i] in [43, 33, 21, 27, 56]:
        vi = F[i+1] if i+1 < len(F) else -1
        if vi < len(V) and vi >= 0:
            print(f"  -> V[{vi}]={repr(V[vi])}", end="")
    print()

PYEOF
Output

F[15982:16020]:
  F[15982] = 43  (PUSH_CONST V[289]='y')
  F[15983] = 289
  F[15984] = 7
  F[15985] = 15989
  F[15986] = 0
  F[15987] = 31
  F[15988] = 16815
  F[15989] = 43  (PUSH_CONST V[0]='o')
  F[15990] = 0
  F[15991] = 43  (PUSH_CONST V[1]='r')
  F[15992] = 1
  F[15993] = 49
  F[15994] = 1
  F[15995] = 21
  F[15996] = 2
  F[15997] = 21
  F[15998] = 3
  F[15999] = 21
  F[16000] = 4
  F[16001] = 21
  F[16002] = 5
  F[16003] = 21
  F[16004] = 6
  F[16005] = 21
  F[16006] = 14
  F[16007] = 33
  F[16008] = 0
  F[16009] = 1
  F[16010] = 30
  F[16011] = 43  (PUSH_CONST V[10]='d')
  F[16012] = 10
  F[16013] = 50
  F[16014] = 1
  F[16015] = 1
  F[16016] = 30
  F[16017] = 43  (PUSH_CONST V[279]='u')
  F[16018] = 279
  F[16019] = 50


Tracing y function body (PC 15989+):
  F[15989] = 43  -> V[0]='o'
  F[15990] = 0
  F[15991] = 43  -> V[1]='r'
  F[15992] = 1
  F[15993] = 49
  F[15994] = 1
  F[15995] = 21  -> V[2]='a'
  F[15996] = 2
  F[15997] = 21  -> V[3]='m'
  F[15998] = 3
  F[15999] = 21  -> V[4]='C'
  F[16000] = 4
  F[16001] = 21  -> V[5]='e'
  F[16002] = 5
  F[16003] = 21  -> V[6]='f'
  F[16004] = 6
  F[16005] = 21  -> V[14]='n'
  F[16006] = 14
  F[16007] = 33  -> V[0]='o'
  F[16008] = 0
  F[16009] = 1
  F[16010] = 30
  F[16011] = 43  -> V[10]='d'
  F[16012] = 10
  F[16013] = 50
  F[16014] = 1
  F[16015] = 1
  F[16016] = 30
  F[16017] = 43  -> V[279]='u'
  F[16018] = 279
  F[16019] = 50
  F[16020] = 1
  F[16021] = 1
  F[16022] = 30
  F[16023] = 43  -> V[0]='o'
  F[16024] = 0
  F[16025] = 27  -> V[0]='o'
  F[16026] = 0
  F[16027] = 30
  F[16028] = 43  -> V[0]='o'
  F[16029] = 0
  F[16030] = 56  -> V[1]='r'
  F[16031] = 1
  F[16032] = 43  -> V[286]='length'
  F[16033] = 286
  F[16034] = 56  -> V[1]='r'
  F[16035] = 1
  F[16036] = 30
  F[16037] = 43  -> V[2]='a'
  F[16038] = 2
  F[16039] = 27  -> V[0]='o'
  F[16040] = 0
  F[16041] = 30
  F[16042] = 43  -> V[1]='r'
  F[16043] = 1
  F[16044] = 56  -> V[1]='r'
  F[16045] = 1
  F[16046] = 43  -> V[286]='length'
  F[16047] = 286
  F[16048] = 56  -> V[1]='r'
  F[16049] = 1
  F[16050] = 30
  F[16051] = 43  -> V[3]='m'
  F[16052] = 3
  F[16053] = 27  -> V[0]='o'
  F[16054] = 0
  F[16055] = 55
  F[16056] = 0
  F[16057] = 1
  F[16058] = 30
  F[16059] = 43  -> V[5]='e'
  F[16060] = 5
  F[16061] = 27  -> V[0]='o'
  F[16062] = 0
  F[16063] = 6
  F[16064] = 16123
  F[16065] = 21  -> V[0]='o'
  F[16066] = 0
  F[16067] = 43  -> V[12]=0
  F[16068] = 12
  F[16069] = 30
  F[16070] = 43  -> V[0]='o'
  F[16071] = 0
  F[16072] = 27  -> V[0]='o'
  F[16073] = 0
  F[16074] = 37
```

Result: function y calculates hash with custom cipher

14. Now we need full disassembly of the y() function from PC 15989 to 16838 we can get that using:
```bash
python3 << 'PYEOF'
import json

with open('F_bytecode.txt') as f:
    F = json.load(f)

V = [
    "o","r","a","m","C","e","f","i","p","h","d","",0,147,"n","fromCharCode",
    252,241,249,246,240,231,"map","join",202,172,191,164,169,190,163,165,
    87,56,53,61,50,52,35,156,250,233,242,255,232,245,243,108,24,3,63,30,
    5,2,11,206,167,160,170,171,182,129,168,124,12,14,19,8,25,33,84,79,69,
    68,71,72,"_",55,64,94,89,83,88,1,22,97,127,120,114,121,46,117,65,76,
    75,77,90,74,115,54,95,82,110,29,36,166,152,159,149,158,134,78,66,85,
    119,20,26,18,44,67,70,73,104,113,162,205,192,200,199,193,214,130,234,
    239,238,230,215,207,204,151,211,248,244,226,227,181,146,132,148,133,
    91,93,86,201,253,247,178,221,150,153,136,220,184,137,145,161,157,131,
    179,154,142,135,28,125,106,123,155,128,32,45,37,42,59,111,57,38,40,
    216,185,174,177,183,141,195,236,251,228,212,210,218,173,188,138,139,
    143,144,175,187,180,4,7,34,194,223,229,15,49,58,96,10,101,16,41,107,
    48,9,31,27,21,6,109,140,224,225,118,122,196,197,203,186,189,176,13,17,
    43,99,23,47,103,112,237,80,105,102,98,213,126,81,100,92,116,51,62,219,
    209,254,208,39,222,198,217,235,False,60,"Boolean","Number","String",
    "j","t","u","arguments","s","random",256,"floor","push","length","0",
    "toString","y","charCodeAt"
]

r = F
MAX = len(F)

def opname(A):
    return {0:'SHR_U',2:'SHL',3:'DEC_PROP',4:'NOP',5:'BIT_NOT',6:'CALL_IF',
            7:'MAKE_FUNC',8:'INC_PROP',9:'LT',10:'LTE',11:'RETURN',
            12:'GTE',13:'DUP_TOP',14:'JMP_TBL',15:'EQ_STRICT',16:'BUILD_STR',
            17:'MOD',18:'ADD',19:'THROW',20:'AND',21:'VAR_DECL',22:'SUB',
            23:'XOR',24:'GT',25:'RET_VAL',26:'MUL',27:'SET_PROP',28:'JIF_F',
            29:'EQ_LOOSE',30:'PUSH_SCOPE',31:'JUMP',32:'JIF_T',33:'LOAD_VAR',
            34:'PUSH_WIN',35:'TYPEOF',36:'NEQ_S',37:'CALL_WH',38:'OR',
            39:'TRY',40:'DIV',41:'NEW',42:'IN',43:'PUSH_C',44:'PUSH_OBJ',
            45:'DEL',46:'INST',47:'SET_IDX',48:'SHR',49:'BIND',50:'CALL_M',
            51:'NEQ_L',52:'NOT',53:'POP_N',54:'PUSH_U',55:'BUILD_ARR',56:'GET_PROP'
            }.get(A, f'?{A}')

def dis1(pc):
    if pc >= MAX: return f"{pc}: EOF", pc+1
    A = r[pc]; n = pc+1
    op = opname(A)
    args = []
    if A == 43:
        vi = r[n]; n+=1
        v = repr(V[vi]) if vi < len(V) else f'?'
        args = [f'V[{vi}]={v}']
    elif A == 33:
        vi = r[n]; n+=1; push = r[n]; n+=1
        v = repr(V[vi]) if vi < len(V) else f'?'
        args = [f'V[{vi}]={v}', f'push={push}']
    elif A == 21:
        vi = r[n]; n+=1
        v = repr(V[vi]) if vi < len(V) else f'?'
        args = [f'V[{vi}]={v}']
    elif A in [27,56,0,2,5,9,10,12,15,17,18,20,22,23,24,26,29,36,38,40,48,51,52]:
        push = r[n]; n+=1; args=[f'push={push}']
    elif A in [28,31,32,6,37]:
        t = r[n]; n+=1; args=[f'->{t}']
    elif A == 49:
        cnt = r[n]; n+=1; args=[f'n={cnt}']
    elif A in [55,16]:
        cnt = r[n]; n+=1; push = r[n]; n+=1; args=[f'cnt={cnt}',f'push={push}']
    elif A == 50:
        na = r[n]; n+=1; push = r[n]; n+=1; args=[f'nargs={na}',f'push={push}']
    elif A == 7:
        entry = r[n]; n+=1; push = r[n]; n+=1; args=[f'entry={entry}',f'push={push}']
    elif A in [8,3]:
        pre=r[n];n+=1;push=r[n];n+=1;args=[f'pre={pre}',f'push={push}']
    elif A == 53:
        cnt=r[n];n+=1;push=r[n];n+=1;args=[f'cnt={cnt}',f'push={push}']
    elif A == 25:
        v=r[n];n+=1;args=[f'inline_val={v}']
    return f"{pc:6d}: {op:12} {' '.join(str(a) for a in args)}", n

# Trace from PC 15989 (y function entry) through 16838 - fully
pc = 15989
while pc < 16840:
    line, next_pc = dis1(pc)
    print(line)
    pc = next_pc
PYEOF
Output

 15989: PUSH_C       V[0]='o'
 15991: PUSH_C       V[1]='r'
 15993: BIND         n=1
 15995: VAR_DECL     V[2]='a'
 15997: VAR_DECL     V[3]='m'
 15999: VAR_DECL     V[4]='C'
 16001: VAR_DECL     V[5]='e'
 16003: VAR_DECL     V[6]='f'
 16005: VAR_DECL     V[14]='n'
 16007: LOAD_VAR     V[0]='o' push=1
 16010: PUSH_SCOPE   
 16011: PUSH_C       V[10]='d'
 16013: CALL_M       nargs=1 push=1
 16016: PUSH_SCOPE   
 16017: PUSH_C       V[279]='u'
 16019: CALL_M       nargs=1 push=1
 16022: PUSH_SCOPE   
 16023: PUSH_C       V[0]='o'
 16025: SET_PROP     push=0
 16027: PUSH_SCOPE   
 16028: PUSH_C       V[0]='o'
 16030: GET_PROP     push=1
 16032: PUSH_C       V[286]='length'
 16034: GET_PROP     push=1
 16036: PUSH_SCOPE   
 16037: PUSH_C       V[2]='a'
 16039: SET_PROP     push=0
 16041: PUSH_SCOPE   
 16042: PUSH_C       V[1]='r'
 16044: GET_PROP     push=1
 16046: PUSH_C       V[286]='length'
 16048: GET_PROP     push=1
 16050: PUSH_SCOPE   
 16051: PUSH_C       V[3]='m'
 16053: SET_PROP     push=0
 16055: BUILD_ARR    cnt=0 push=1
 16058: PUSH_SCOPE   
 16059: PUSH_C       V[5]='e'
 16061: SET_PROP     push=0
 16063: CALL_IF      ->16123
 16065: VAR_DECL     V[0]='o'
 16067: PUSH_C       V[12]=0
 16069: PUSH_SCOPE   
 16070: PUSH_C       V[0]='o'
 16072: SET_PROP     push=0
 16074: CALL_WH      ->16114
 16076: LOAD_VAR     V[0]='o' push=1
 16079: PUSH_C       V[221]=16
 16081: LT           push=1
 16083: JIF_F        ->16112
 16085: LOAD_VAR     V[0]='o' push=1
 16088: PUSH_C       V[209]=4
 16090: SHL          push=1
 16092: LOAD_VAR     V[0]='o' push=1
 16095: PUSH_C       V[221]=16
 16097: MOD          push=1
 16099: ADD          push=1
 16101: PUSH_SCOPE   
 16102: PUSH_C       V[5]='e'
 16104: GET_PROP     push=1
 16106: PUSH_C       V[285]='push'
 16108: CALL_M       nargs=1 push=0
 16111: RETURN       
 16112: RET_VAL      inline_val=1
 16114: PUSH_SCOPE   
 16115: PUSH_C       V[0]='o'
 16117: INC_PROP     pre=0 push=0
 16120: JUMP         ->16074
 16122: RETURN       
 16123: PUSH_SCOPE   
 16124: PUSH_C       V[5]='e'
 16126: GET_PROP     push=1
 16128: PUSH_C       V[286]='length'
 16130: GET_PROP     push=1
 16132: PUSH_SCOPE   
 16133: PUSH_C       V[6]='f'
 16135: SET_PROP     push=0
 16137: CALL_IF      ->16311
 16139: VAR_DECL     V[0]='o'
 16141: PUSH_C       V[12]=0
 16143: PUSH_SCOPE   
 16144: PUSH_C       V[0]='o'
 16146: SET_PROP     push=0
 16148: VAR_DECL     V[2]='a'
 16150: PUSH_C       V[12]=0
 16152: PUSH_SCOPE   
 16153: PUSH_C       V[2]='a'
 16155: SET_PROP     push=0
 16157: CALL_WH      ->16302
 16159: LOAD_VAR     V[0]='o' push=1
 16162: LOAD_VAR     V[6]='f' push=1
 16165: LT           push=1
 16167: JIF_F        ->16300
 16169: LOAD_VAR     V[0]='o' push=1
 16172: LOAD_VAR     V[2]='a' push=1
 16175: ADD          push=1
 16177: PUSH_SCOPE   
 16178: PUSH_C       V[5]='e'
 16180: GET_PROP     push=1
 16182: PUSH_SCOPE   
 16183: PUSH_C       V[0]='o'
 16185: GET_PROP     push=1
 16187: GET_PROP     push=1
 16189: ADD          push=1
 16191: PUSH_SCOPE   
 16192: PUSH_C       V[5]='e'
 16194: GET_PROP     push=1
 16196: PUSH_SCOPE   
 16197: PUSH_C       V[2]='a'
 16199: GET_PROP     push=1
 16201: GET_PROP     push=1
 16203: ADD          push=1
 16205: PUSH_C       V[83]=1
 16207: SHR          push=1
 16209: LOAD_VAR     V[0]='o' push=1
 16212: LOAD_VAR     V[3]='m' push=1
 16215: MOD          push=1
 16217: PUSH_SCOPE   
 16218: PUSH_C       V[1]='r'
 16220: GET_PROP     push=1
 16222: PUSH_C       V[290]='charCodeAt'
 16224: CALL_M       nargs=1 push=1
 16227: ADD          push=1
 16229: LOAD_VAR     V[6]='f' push=1
 16232: PUSH_C       V[83]=1
 16234: SUB          push=1
 16236: AND          push=1
 16238: PUSH_SCOPE   
 16239: PUSH_C       V[2]='a'
 16241: SET_PROP     push=0
 16243: PUSH_SCOPE   
 16244: PUSH_C       V[5]='e'
 16246: GET_PROP     push=1
 16248: PUSH_SCOPE   
 16249: PUSH_C       V[0]='o'
 16251: GET_PROP     push=1
 16253: GET_PROP     push=1
 16255: PUSH_SCOPE   
 16256: PUSH_C       V[4]='C'
 16258: SET_PROP     push=0
 16260: PUSH_SCOPE   
 16261: PUSH_C       V[5]='e'
 16263: GET_PROP     push=1
 16265: PUSH_SCOPE   
 16266: PUSH_C       V[2]='a'
 16268: GET_PROP     push=1
 16270: GET_PROP     push=1
 16272: PUSH_SCOPE   
 16273: PUSH_C       V[5]='e'
 16275: GET_PROP     push=1
 16277: PUSH_SCOPE   
 16278: PUSH_C       V[0]='o'
 16280: GET_PROP     push=1
 16282: SET_PROP     push=0
 16284: LOAD_VAR     V[4]='C' push=1
 16287: PUSH_SCOPE   
 16288: PUSH_C       V[5]='e'
 16290: GET_PROP     push=1
 16292: PUSH_SCOPE   
 16293: PUSH_C       V[2]='a'
 16295: GET_PROP     push=1
 16297: SET_PROP     push=0
 16299: RETURN       
 16300: RET_VAL      inline_val=1
 16302: PUSH_SCOPE   
 16303: PUSH_C       V[0]='o'
 16305: INC_PROP     pre=0 push=0
 16308: JUMP         ->16157
 16310: RETURN       
 16311: CALL_IF      ->16590
 16313: VAR_DECL     V[1]='r'
 16315: PUSH_C       V[12]=0
 16317: PUSH_SCOPE   
 16318: PUSH_C       V[1]='r'
 16320: SET_PROP     push=0
 16322: VAR_DECL     V[3]='m'
 16324: PUSH_C       V[12]=0
 16326: PUSH_SCOPE   
 16327: PUSH_C       V[3]='m'
 16329: SET_PROP     push=0
 16331: VAR_DECL     V[14]='n'
 16333: PUSH_C       V[12]=0
 16335: PUSH_SCOPE   
 16336: PUSH_C       V[14]='n'
 16338: SET_PROP     push=0
 16340: CALL_WH      ->16581
 16342: LOAD_VAR     V[1]='r' push=1
 16345: LOAD_VAR     V[2]='a' push=1
 16348: LT           push=1
 16350: JIF_F        ->16579
 16352: LOAD_VAR     V[3]='m' push=1
 16355: LOAD_VAR     V[14]='n' push=1
 16358: XOR          push=1
 16360: PUSH_SCOPE   
 16361: PUSH_C       V[5]='e'
 16363: GET_PROP     push=1
 16365: PUSH_SCOPE   
 16366: PUSH_C       V[3]='m'
 16368: GET_PROP     push=1
 16370: GET_PROP     push=1
 16372: PUSH_SCOPE   
 16373: PUSH_C       V[5]='e'
 16375: GET_PROP     push=1
 16377: PUSH_SCOPE   
 16378: PUSH_C       V[14]='n'
 16380: GET_PROP     push=1
 16382: GET_PROP     push=1
 16384: XOR          push=1
 16386: ADD          push=1
 16388: LOAD_VAR     V[6]='f' push=1
 16391: PUSH_C       V[83]=1
 16393: SUB          push=1
 16395: AND          push=1
 16397: PUSH_SCOPE   
 16398: PUSH_C       V[14]='n'
 16400: SET_PROP     push=0
 16402: PUSH_SCOPE   
 16403: PUSH_C       V[5]='e'
 16405: GET_PROP     push=1
 16407: PUSH_SCOPE   
 16408: PUSH_C       V[3]='m'
 16410: GET_PROP     push=1
 16412: GET_PROP     push=1
 16414: PUSH_SCOPE   
 16415: PUSH_C       V[4]='C'
 16417: SET_PROP     push=0
 16419: PUSH_SCOPE   
 16420: PUSH_C       V[5]='e'
 16422: GET_PROP     push=1
 16424: PUSH_SCOPE   
 16425: PUSH_C       V[14]='n'
 16427: GET_PROP     push=1
 16429: GET_PROP     push=1
 16431: PUSH_SCOPE   
 16432: PUSH_C       V[5]='e'
 16434: GET_PROP     push=1
 16436: PUSH_SCOPE   
 16437: PUSH_C       V[3]='m'
 16439: GET_PROP     push=1
 16441: SET_PROP     push=0
 16443: LOAD_VAR     V[4]='C' push=1
 16446: PUSH_SCOPE   
 16447: PUSH_C       V[5]='e'
 16449: GET_PROP     push=1
 16451: PUSH_SCOPE   
 16452: PUSH_C       V[14]='n'
 16454: GET_PROP     push=1
 16456: SET_PROP     push=0
 16458: LOAD_VAR     V[1]='r' push=1
 16461: PUSH_SCOPE   
 16462: PUSH_C       V[0]='o'
 16464: GET_PROP     push=1
 16466: PUSH_C       V[290]='charCodeAt'
 16468: CALL_M       nargs=1 push=1
 16471: PUSH_SCOPE   
 16472: PUSH_C       V[4]='C'
 16474: SET_PROP     push=0
 16476: LOAD_VAR     V[4]='C' push=1
 16479: LOAD_VAR     V[3]='m' push=1
 16482: LOAD_VAR     V[14]='n' push=1
 16485: ADD          push=1
 16487: ADD          push=1
 16489: PUSH_SCOPE   
 16490: PUSH_C       V[4]='C'
 16492: SET_PROP     push=0
 16494: LOAD_VAR     V[4]='C' push=1
 16497: PUSH_SCOPE   
 16498: PUSH_C       V[5]='e'
 16500: GET_PROP     push=1
 16502: PUSH_SCOPE   
 16503: PUSH_C       V[3]='m'
 16505: GET_PROP     push=1
 16507: GET_PROP     push=1
 16509: PUSH_SCOPE   
 16510: PUSH_C       V[5]='e'
 16512: GET_PROP     push=1
 16514: PUSH_SCOPE   
 16515: PUSH_C       V[14]='n'
 16517: GET_PROP     push=1
 16519: GET_PROP     push=1
 16521: XOR          push=1
 16523: XOR          push=1
 16525: PUSH_SCOPE   
 16526: PUSH_C       V[4]='C'
 16528: SET_PROP     push=0
 16530: LOAD_VAR     V[4]='C' push=1
 16533: PUSH_C       V[43]=255
 16535: AND          push=1
 16537: PUSH_SCOPE   
 16538: PUSH_C       V[4]='C'
 16540: SET_PROP     push=0
 16542: LOAD_VAR     V[4]='C' push=1
 16545: PUSH_SCOPE   
 16546: PUSH_C       V[5]='e'
 16548: GET_PROP     push=1
 16550: PUSH_SCOPE   
 16551: PUSH_C       V[3]='m'
 16553: GET_PROP     push=1
 16555: SET_PROP     push=0
 16557: LOAD_VAR     V[3]='m' push=1
 16560: PUSH_C       V[83]=1
 16562: ADD          push=1
 16564: LOAD_VAR     V[6]='f' push=1
 16567: PUSH_C       V[83]=1
 16569: SUB          push=1
 16571: AND          push=1
 16573: PUSH_SCOPE   
 16574: PUSH_C       V[3]='m'
 16576: SET_PROP     push=0
 16578: RETURN       
 16579: RET_VAL      inline_val=1
 16581: PUSH_SCOPE   
 16582: PUSH_C       V[1]='r'
 16584: INC_PROP     pre=0 push=0
 16587: JUMP         ->16340
 16589: RETURN       
 16590: CALL_IF      ->16744
 16592: VAR_DECL     V[0]='o'
 16594: PUSH_C       V[12]=0
 16596: PUSH_SCOPE   
 16597: PUSH_C       V[0]='o'
 16599: SET_PROP     push=0
 16601: VAR_DECL     V[1]='r'
 16603: PUSH_C       V[12]=0
 16605: PUSH_SCOPE   
 16606: PUSH_C       V[1]='r'
 16608: SET_PROP     push=0
 16610: CALL_WH      ->16735
 16612: LOAD_VAR     V[0]='o' push=1
 16615: LOAD_VAR     V[6]='f' push=1
 16618: PUSH_C       V[83]=1
 16620: SHL          push=1
 16622: LT           push=1
 16624: JIF_F        ->16733
 16626: LOAD_VAR     V[0]='o' push=1
 16629: LOAD_VAR     V[6]='f' push=1
 16632: MOD          push=1
 16634: PUSH_SCOPE   
 16635: PUSH_C       V[1]='r'
 16637: SET_PROP     push=0
 16639: LOAD_VAR     V[1]='r' push=1
 16642: JIF_F        ->16689
 16644: CALL_IF      ->16687
 16646: PUSH_SCOPE   
 16647: PUSH_C       V[5]='e'
 16649: GET_PROP     push=1
 16651: PUSH_SCOPE   
 16652: PUSH_C       V[1]='r'
 16654: GET_PROP     push=1
 16656: GET_PROP     push=1
 16658: PUSH_SCOPE   
 16659: PUSH_C       V[5]='e'
 16661: GET_PROP     push=1
 16663: LOAD_VAR     V[1]='r' push=1
 16666: PUSH_C       V[83]=1
 16668: SUB          push=1
 16670: GET_PROP     push=1
 16672: XOR          push=1
 16674: PUSH_SCOPE   
 16675: PUSH_C       V[5]='e'
 16677: GET_PROP     push=1
 16679: PUSH_SCOPE   
 16680: PUSH_C       V[1]='r'
 16682: GET_PROP     push=1
 16684: SET_PROP     push=0
 16686: RETURN       
 16687: JUMP         ->16732
 16689: CALL_IF      ->16732
 16691: PUSH_SCOPE   
 16692: PUSH_C       V[5]='e'
 16694: GET_PROP     push=1
 16696: PUSH_SCOPE   
 16697: PUSH_C       V[1]='r'
 16699: GET_PROP     push=1
 16701: GET_PROP     push=1
 16703: PUSH_SCOPE   
 16704: PUSH_C       V[5]='e'
 16706: GET_PROP     push=1
 16708: LOAD_VAR     V[6]='f' push=1
 16711: PUSH_C       V[83]=1
 16713: SUB          push=1
 16715: GET_PROP     push=1
 16717: XOR          push=1
 16719: PUSH_SCOPE   
 16720: PUSH_C       V[5]='e'
 16722: GET_PROP     push=1
 16724: PUSH_SCOPE   
 16725: PUSH_C       V[1]='r'
 16727: GET_PROP     push=1
 16729: SET_PROP     push=0
 16731: RETURN       
 16732: RETURN       
 16733: RET_VAL      inline_val=1
 16735: PUSH_SCOPE   
 16736: PUSH_C       V[0]='o'
 16738: INC_PROP     pre=0 push=0
 16741: JUMP         ->16610
 16743: RETURN       
 16744: PUSH_C       V[11]=''
 16746: PUSH_C       V[12]=0
 16748: MAKE_FUNC    entry=16753 push=1
 16751: JUMP         ->16789
 16753: PUSH_C       V[0]='o'
 16755: BIND         n=0
 16757: LOAD_VAR     V[0]='o' push=1
 16760: PUSH_C       V[221]=16
 16762: LT           push=1
 16764: JIF_F        ->16770
 16766: PUSH_C       V[287]='0'
 16768: JUMP         ->16772
 16770: PUSH_C       V[11]=''
 16772: PUSH_C       V[221]=16
 16774: PUSH_SCOPE   
 16775: PUSH_C       V[0]='o'
 16777: GET_PROP     push=1
 16779: PUSH_C       V[288]='toString'
 16781: CALL_M       nargs=1 push=1
 16784: ADD          push=1
 16786: RET_VAL      inline_val=2
 16788: RETURN       
 16789: PUSH_SCOPE   
 16790: PUSH_C       V[5]='e'
 16792: GET_PROP     push=1
 16794: PUSH_C       V[22]='map'
 16796: CALL_M       nargs=1 push=1
 16799: PUSH_C       V[23]='join'
 16801: CALL_M       nargs=1 push=1
 16804: PUSH_SCOPE   
 16805: PUSH_C       V[14]='n'
 16807: SET_PROP     push=0
 16809: LOAD_VAR     V[14]='n' push=1
 16812: RET_VAL      inline_val=2
 16814: RETURN       
 16815: LOAD_VAR     V[0]='o' push=1
 16818: LOAD_VAR     V[1]='r' push=1
 16821: PUSH_SCOPE   
 16822: PUSH_C       V[289]='y'
 16824: CALL_M       nargs=2 push=1
 16827: PUSH_SCOPE   
 16828: PUSH_C       V[2]='a'
 16830: SET_PROP     push=0
 16832: RETURN       
 16833: LOAD_VAR     V[2]='a' push=1
 16836: RET_VAL      inline_val=2
16838: EOF
16839: EOF
```

15. By analysing the disassembly we can reconstruct it into javascript using:
```js
// ============================================================
// Aliyun Captcha VM Hash — Full Reverse Engineering Notes
// ============================================================
//
// ENTRY POINT called by the page:
//   er = P(0, [], F, V, ez, [er, (td && td)(77, 19)]) + er
//   where er = input JSON string, td = the salt token (e.g. '0000')
//   ez = { r: 1 }  (means "return the pop value")
//
// P / t  is a stack-based virtual machine. Its opcodes live in F[]
// and its constant pool lives in V[].
//
// ─────────────────────────────────────────────────────────────
// IS THIS MD5?  NO.
// ─────────────────────────────────────────────────────────────
// window.__ALIYUN_CRYPT is undefined, so no external crypto lib
// is in play. The output is 32 hex chars like MD5 but the algo
// is entirely different — a custom 16-byte stateful stream cipher.
//
// ─────────────────────────────────────────────────────────────
// WHY DOES parseInt MATTER?
// ─────────────────────────────────────────────────────────────
// The call site is:  (td && td)(77, 19)
// When td is a STRING like '0000', this is  '0000'(77, 19) which
// would throw in strict contexts, but in this VM the outer wrapper
// evaluates td(77,19) differently. From the user's console tests:
//
//   t(0,[], F, V, ez, [input, '0000'])          -> '9333ef...'
//   t(0,[], F, V, ez, [input, '0000000000001']) -> '4ed342...' (salt≠0)
//   t(0,[], F, V, ez, [input, '0010'])          -> '70eb2c...' (salt≠0)
//
// The salt string is passed as-is to the y() function (argument[1]).
// The KSA uses r.charCodeAt(i % m) on the RAW STRING, NOT parseInt.
// The "salt doesn't matter when it's '0000'/'000'/'0'" observation
// is explained purely by the KSA arithmetic: when every salt byte
// is '0' (charCode 48), those 48s are added into j then masked to
// (f-1)=15, which happens to produce the same permutation as any
// other all-same-byte salt that maps to the same residue mod 16.
// '0010' breaks it because charCode('1')=49 ≠ charCode('0')=48.
//
// ─────────────────────────────────────────────────────────────
// ALGORITHM  (decompiled from VM bytecodes, exact formulae)
// ─────────────────────────────────────────────────────────────
//
// function aliHash(inputStr, saltStr):
//
//   Phase 0 — UTF-8 encode  [PC 16007-16025]
//     o = unescape(encodeURIComponent(inputStr))
//     a = o.length          // byte length of UTF-8 encoded input
//     r = saltStr
//     m = r.length          // salt string length
//
//   Phase 1 — State init  [PC 16067-16135]
//     e[i] = (i << 4) + (i % 16)   for i in [0..15]
//     → e = [0, 17, 34, 51, 68, 85, 102, 119, 136, 153, 170, 187, 204, 221, 238, 255]
//     f = 16                // state size
//
//   Phase 2 — KSA (Key Schedule)  [PC 16139-16310]
//     i = 0, j = 0
//     while i < f:
//       j = (((i + j + e[i] + e[j]) >> 1) + r.charCodeAt(i % m)) & (f-1)
//       swap(e[i], e[j])
//       i++
//
//   Phase 3 — PRGA (stream processing)  [PC 16313-16589]
//     // local loop counter shadows outer 'r' (salt); 'o' (input str) still visible
//     idx = 0, p = 0, q = 0
//     while idx < a:          // a = UTF-8 input length
//       q   = ((p ^ q) + (e[p] ^ e[q])) & (f-1)
//       swap(e[p], e[q])
//       C   = o.charCodeAt(idx)
//       C   = (C + p + q) ^ e[p] ^ e[q]    // uses POST-swap e[p] and e[q]
//       C   = C & 255
//       e[p] = C
//       p   = (p + 1) & (f-1)
//       idx++
//
//   Phase 4 — Final diffusion  [PC 16592-16743]
//     for step in [0 .. 2*f-1]:
//       pos = step % f
//       if pos != 0:  e[pos] ^= e[pos-1]
//       else:         e[0]   ^= e[f-1]
//
//   Phase 5 — Hex encode  [PC 16744-16814]
//     return e.map(b => (b < 16 ? '0' : '') + b.toString(16)).join('')
//
// ─────────────────────────────────────────────────────────────
// IMPLEMENTATION
// ─────────────────────────────────────────────────────────────

function aliHash(inputStr, saltStr) {
    // Phase 0: UTF-8 encode
    var o = unescape(encodeURIComponent(inputStr));
    var a = o.length;
    var r = saltStr;
    var m = r.length;

    // Phase 1: State init  e[i] = (i<<4) + (i%16)
    var e = [], f;
    for (var _i = 0; _i < 16; _i++) {
        e.push((_i << 4) + (_i % 16));
    }
    f = e.length; // 16

    // Phase 2: KSA
    var i = 0, j = 0, tmp;
    while (i < f) {
        j = (((i + j + e[i] + e[j]) >> 1) + r.charCodeAt(i % m)) & (f - 1);
        tmp = e[i]; e[i] = e[j]; e[j] = tmp;
        i++;
    }

    // Phase 3: PRGA
    var idx = 0, p = 0, q = 0, C;
    while (idx < a) {
        q = ((p ^ q) + (e[p] ^ e[q])) & (f - 1);
        tmp = e[p]; e[p] = e[q]; e[q] = tmp;
        C = o.charCodeAt(idx);
        C = (C + p + q) ^ e[p] ^ e[q];
        C = C & 255;
        e[p] = C;
        p = (p + 1) & (f - 1);
        idx++;
    }

    // Phase 4: Final diffusion
    for (var step = 0; step < 2 * f; step++) {
        var pos = step % f;
        if (pos !== 0) {
            e[pos] = e[pos] ^ e[pos - 1];
        } else {
            e[0] = e[0] ^ e[f - 1];
        }
    }

    // Phase 5: Hex encode
    return e.map(function(b) {
        return (b < 16 ? '0' : '') + b.toString(16);
    }).join('');
}

// ─────────────────────────────────────────────────────────────
// VERIFICATION
// ─────────────────────────────────────────────────────────────

var input = '{"TrackList":{"mc":"","tc":"","mu":"","te":"","mp":"","tmv":"","ks":"","fi":"","startTime":1782100652835},"TrackStartTime":1782100652835,"VerifyTime":1782100652862,"arg":"JjObDGdh/ywcWQ=="}';
var expected = '9333ef7396dd56dbb9d6e8f31e8f6014';

console.log('=== Verification ===');
console.log('Expected : ' + expected);

// All zero-byte salts produce the same result (charCode '0'=48, all map same)
['0000', '000', '0'].forEach(function(s) {
    var res = aliHash(input, s);
    console.log('salt=' + JSON.stringify(s).padEnd(20) + ' -> ' + res + '  PASS=' + (res === expected));
});

// Salts with non-zero digits break it (charCode differs)
['0010', '0000000000001', '1'].forEach(function(s) {
    var res = aliHash(input, s);
    console.log('salt=' + JSON.stringify(s).padEnd(20) + ' -> ' + res + '  PASS=' + (res === expected) + ' (expected different)');
});

// ─────────────────────────────────────────────────────────────
// INTERMEDIATE STATE TRACE (for the given example, salt='0000')
// ─────────────────────────────────────────────────────────────
console.log('\n=== Intermediate state trace (salt="0000") ===');

function aliHashTrace(inputStr, saltStr) {
    var o = unescape(encodeURIComponent(inputStr));
    var a = o.length;
    var r = saltStr;
    var m = r.length;

    var e = [], f;
    for (var _i = 0; _i < 16; _i++) e.push((_i << 4) + (_i % 16));
    f = e.length;

    console.log('Phase 1 - initial e: [' + e + ']');
    console.log('  f=' + f + ', a (UTF-8 len)=' + a + ', m (salt len)=' + m);

    var i = 0, j = 0, tmp;
    while (i < f) {
        j = (((i + j + e[i] + e[j]) >> 1) + r.charCodeAt(i % m)) & (f - 1);
        tmp = e[i]; e[i] = e[j]; e[j] = tmp;
        i++;
    }
    console.log('Phase 2 - after KSA:  [' + e + ']');

    var idx = 0, p = 0, q = 0, C;
    while (idx < a) {
        q = ((p ^ q) + (e[p] ^ e[q])) & (f - 1);
        tmp = e[p]; e[p] = e[q]; e[q] = tmp;
        C = o.charCodeAt(idx);
        C = (C + p + q) ^ e[p] ^ e[q];
        C = C & 255;
        e[p] = C;
        p = (p + 1) & (f - 1);
        idx++;
    }
    console.log('Phase 3 - after PRGA: [' + e + ']');

    for (var step = 0; step < 2 * f; step++) {
        var pos = step % f;
        if (pos !== 0) e[pos] ^= e[pos - 1];
        else           e[0]   ^= e[f - 1];
    }
    console.log('Phase 4 - after diff: [' + e + ']');

    return e.map(function(b) { return (b<16?'0':'') + b.toString(16); }).join('');
}

var traced = aliHashTrace(input, '0000');
console.log('Final:  ' + traced);
console.log('Match:  ' + (traced === expected));
JSEOF
node /tmp/aliyun_hash_final.js
Output

=== Verification ===
Expected : 9333ef7396dd56dbb9d6e8f31e8f6014
salt="0000"               -> 9333ef7396dd56dbb9d6e8f31e8f6014  PASS=true
salt="000"                -> 9333ef7396dd56dbb9d6e8f31e8f6014  PASS=true
salt="0"                  -> 9333ef7396dd56dbb9d6e8f31e8f6014  PASS=true
salt="0010"               -> 70eb2c0d763b5e3b78be269cc9a0fd28  PASS=false (expected different)
salt="0000000000001"      -> 4ed342e636bf55896c0244fa419b6f1f  PASS=false (expected different)
salt="1"                  -> 1002e563d3c7d303b20a36d2044d70ec  PASS=false (expected different)

=== Intermediate state trace (salt="0000") ===
Phase 1 - initial e: [0,17,34,51,68,85,102,119,136,153,170,187,204,221,238,255]
  f=16, a (UTF-8 len)=189, m (salt len)=4
Phase 2 - after KSA:  [204,136,238,51,34,170,17,68,85,187,102,119,0,153,255,221]
Phase 3 - after PRGA: [124,71,124,64,121,174,192,6,239,13,81,37,246,124,126,155]
Phase 4 - after diff: [147,51,239,115,150,221,86,219,185,214,232,243,30,143,96,20]
Final:  9333ef7396dd56dbb9d6e8f31e8f6014
Match:  true
```

### Hash function reverse engineering completed more details at [VMachine_HASH](./VMachine_Hash_report.md)

## TrackJson Arg Value generation

TrackJson arg is assembled in this code
```js

 a >= 4 ? a > 6 ? (e -= -23,
m[y] = b + ({ 0: s  })[0](127, 82)) : a <= 4 ? (e ^= 24,
P[(s && s)(242, 19)](void 0, v)) : a < 6 ? (P[(s && s)(242, 19)](void 0, v),
e += -45) : 0 > Math.abs(!v) ? e ^= 2 : e ^= 44 : a >= 3 ? (O[nx.F(c, nx.c(s, [170, s()][0], [91, s()][0]))] = eh,

// This code
e -= -28) : a >= 1 ? a <= 1 ? (n_[nx.c(s, 99 * !-s, 88 * !-s)] = I,
```

```js
Subject:
n_[nx.c(s, 99 * !-s, 88 * !-s)] = I,
```

### Deobfuscation Time
```js
 nx.c = function(t, n, e) {
                            return t(n, e)
                        }

// resoltion

n_[s(99 * !-s, 88 * !-s)] = I,

s = function(t, n) { return (td || td)(w.U(t, 2), n)}

-s = NaN
!NaN = true
!-s = !Nan = true
99 * true = 99
88 * true = 88

// Resolution
n_[s(99,88)] =I,

s(99,88) = 'arg'

// Resoltuion

n_['arg'] = I,
n_.arg = I,
```
```optional
// Explanation: Why -s returns NaN
// The unary minus operator (-) tries to convert its operand to a number
// When you apply - to a function, JavaScript:
// 1. First calls `ToNumber()` on the function
// 2. `ToNumber()` calls `ToPrimitive()` with hint "number"
// 3. For functions, `ToPrimitive()` calls `valueOf()` then `toString()`
//  4. Functions convert to strings (their source code)
// 5. The string cannot be parsed as a number
// 6. Result: `NaN`

// Example

A = '10'
Number(A) // parsed string to number succesfully
10
-A = -5

B = '10A'
Number(B) // Cant parse because it contains A which is not a number
NaN
-B
NaN


// Final thing
In JavaScript, NaN is falsy (coerces to false in boolean contexts), but it is not equal to false:

!NaN = true
!false = true
false === Nan result is  false
```

### How the variable I is defined?
Variable I in the context of `n_['arg'] = I` is most likely generated by the virtual machine more specifically
this call

```js
P[(s && s)(242, 19)](void 0, v)) : a < 6 ? (P[(s && s)(242, 19)](void 0, v)

// Resoltuion
P[s(242,19)](void 0,v) :  a < 6 ? P[s(242,19)](void 0,v)

s(242,19) = 'apply'

//  resoltuion
P['apply'](void 0,v)
// resoltuion
P['apply'](unndefined,v)

typeof v is object

Object.keys(v)

["0","1","2","3","4","5"]


// These keys are arguments passed to interpereter P using P.apply(undefined,arguments)

0 is the 1st argument and it resolves to `0`
1 is the 2nd arg resolving to `[]`
2 is the 3rd argument resolving to bytecode shown in L_bytecode.txt
3 is the 4th arg resolving to constant pool R shown below
4 is the 5th arg resolving to JSON {r:1}
5 is the 6th arg resolving to variable m from the code m[y] = .... and it contains: ['dynamic-string','constant_string]
eg ["IWZInZXnCl","4xrihv8zb8tf1mfj"]
m[y] = m[1] = constant_string in this case 4xrihv8zb8tf1mfj
m[0] = dynamic_string AKA CertifyID from init captcha response 
P.apply(undefined,arguments) resolves to P(all the args from var v)

// Important Note regarding the constant_string

It is important to note that there are 10 possible values of the constant strings but only one is used at a time
in my case that was 4xrihv8zb8tf1mfj

Heres list of all the constant string def

1. In case 0:
 m[y] = B + (~s ? s : 0)(277, 56)
B = nx.c(s, ~s ? 123 : 1, ~s ? 98
B = fxt8jzp3
s(277,56) = p0ienz71
// resolution  m[y] = fxt8jzp3p0ienz71

2. In case 1 (my case):
 m[y] = M + nx.c(s, nx.j(-s, 59), nx.j(-s, 26)),
nx.j = function(t, n) {return t || n}
M = s(286, 78)
M = 4xrihv8z
s(59,26) = b8tf1mfj
// resolution m[y] = 4xrihv8zb8tf1mfj

3. In case 3:
m[y] = k + (~s ? s : 8)(22, 44)

s(22,44) = '8g52k8hy'

```

**R constant Pool**
```
["o", "n", "r", "t", "e", "_", "Boolean", "a", "Number", "h", "String", "m", "", 0, 59, "fromCharCode", 90, 79, 84, 89, "map", "join", "C", 118, 20, 2, 25, 23, "l", 239, 138, 129, 140, 128, 139, 186, 189, 166, 172, 130, 159, 155, "c", 252, 152, 153, 147, 169, 174, 181, 191, 145, 146, 136, "f", 109, 8, 30, 14, 12, 29, "s", 220, 178, 185, 175, "arguments", 1, 32, 50, 10, 51, 6, 44, 37, 16, 46, 11, 62, 19, 43, 60, 33, 53, 34, 7, 26, 48, 5, 4, 61, 13, 47, 49, 18, 27, 22, 17, 39, 56, 41, 38, 55, 31, 15, 58, 52, 40, 57, 45, 35, 36, 42, 54, 63, 3, 24, 28, 9, 21, "length", "charCodeAt", 255, null, "call"]
          
```

