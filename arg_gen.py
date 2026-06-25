# The key area is PC 1243-1389.

# State at PC 1243 (start of output byte computation):
# Locals: n=idx(loop ctr), e=p-ptr, a=q-ptr (just updated + maybe swapped), m=undefined
# r = permutation array

# PC 1243: LOAD m push=1       → stack: [m_prev]  (m is uninit/0 at loop start)
# PC 1246: LOAD e push=1       → stack: [m_prev, e]
# PC 1249-1259: r[e] push=1    → stack: [m_prev, e, r[e]]
# PC 1261: ADD push=1          → stack: [m_prev, e+r[e]]
# PC 1263: ADD push=1          → stack: [m_prev + e + r[e]]
# SET m                        → m = (old_m) + e + r[e]
#
# BUT old_m was set JUST BEFORE at PC 1225-1241:
# PC 1225: LOAD n push=1      → [idx]
# PC 1228-1231: o push=1      → [idx, o_string]
# PC 1233: PUSH_C 'charCodeAt'
# PC 1235: CALL_M 1 push=1    → [o.charCodeAt(idx)]
# PC 1238-1241: SET m          → m = o.charCodeAt(idx)
#
# So after PC 1268: m = o.charCodeAt(idx) + e + r[e]
#
# PC 1270-1295: m = m - (a + r[a])
# PC 1270: LOAD m push=1      → [m]
# PC 1273: LOAD a push=1      → [m, a]
# PC 1276-1286: r[a] push=1   → [m, a, r[a]]
# PC 1288: ADD push=1         → [m, a+r[a]]
# PC 1290: SUB push=1         → [m - (a+r[a])]
# SET m                       → m = m - a - r[a]
#                              = charCode + e + r[e] - a - r[a]
#
# PC 1297-1331: m = m XOR (r[e]+r[a])
# PC 1297: LOAD m             → [m]
# PC 1300-1310: r[e]          → [m, r[e]]
# PC 1312-1322: r[a]          → [m, r[e], r[a]]
# PC 1324: ADD                → [m, r[e]+r[a]]
# PC 1326: XOR                → [m ^ (r[e]+r[a])]
# SET m
#
# PC 1333-1389: m = m XOR r[(r[e]+r[a]) & 63]
# PC 1333: LOAD m             → [m]
# PC 1336: LOAD r (the array) → [m, r_arr]   ← LOAD_VAR V[2]='r'
# PC 1342-1351: r[e]          → [m, r_arr, r[e]]    ← r[scope.e]
# PC 1353-1363: r[a]          → [m, r_arr, r[e], r[a]]
# PC 1365: ADD                → [m, r_arr, r[e]+r[a]]
# PC 1368-1380: (r[e]+r[a]) & (r.length-1) → [m, r_arr, (r[e]+r[a])&63]
# PC 1382: GET_PROP           → [m, r_arr[(r[e]+r[a])&63]]
# PC 1384: XOR                → [m ^ r[(r[e]+r[a])&63]]
# SET m
#
# PC 1391-1401: m = m & 255
# PC 1403-1419: m = String.fromCharCode(m)  [h.fromCharCode(m)]
# PC 1421-1432: t = t + m
# PC 1434-1459: e = (e+1) & 63

print("Final byte formula confirmed:")
print("  m = charCode(o[idx]) + e + r[e] - a - r[a]")
print("  m = m ^ (r[e] + r[a])")
print("  m = m ^ r[(r[e]+r[a]) & 63]")
print("  m = m & 255")
print()
print("NOTE: r[e] and r[a] used here are POST-swap values")
print("(the XOR-swap happens before this computation)")


# XOR-swap of r[e] and r[a]:
#   r[e] = r[e] ^ r[a]
#   r[a] = r[e] ^ r[a]   (= old_r[e])
#   r[e] = r[e] ^ r[a]   (= old_r[a])
# This is correct standard XOR swap.

# Now let me verify my computation independently in Python:
import math

def generateArg(certifyId, constant):
    # UTF-8 encode
    o = bytes(certifyId, 'utf-8').decode('latin-1')  # = unescape(encodeURIComponent(...))
    n = constant

    # 64-element permutation table
    r = [32,50,10,51,6,44,37,16,46,11,62,19,43,25,23,30,
         60,33,53,34,7,26,12,48,5,2,20,4,61,13,47,49,
         18,29,27,22,1,17,39,56,41,38,55,31,15,58,52,40,
         8,57,45,35,59,36,42,54,63,3,24,28,14,9,0,21]

    L = len(r)  # 64

    # KSA
    i, j = 0, 0
    while i < L:
        j = (((i + j + r[i] + r[j]) >> 1) + ord(n[i % len(n)])) & (L - 1)
        if i != j:
            r[i] ^= r[j]; r[j] ^= r[i]; r[i] ^= r[j]
        i += 1

    print(f"After KSA: r[:8] = {r[:8]}")

    # PRGA
    t = ''
    e, a = 0, 0
    for idx in range(len(o)):
        # Update a (q)
        a = ((e ^ a) + (r[e] ^ r[a])) & (L - 1)
        # XOR-swap if different
        if e != a:
            r[e] ^= r[a]; r[a] ^= r[e]; r[e] ^= r[a]

        m = ord(o[idx])
        m = m + e + r[e] - a - r[a]
        m = m ^ (r[e] + r[a])
        m = m ^ r[(r[e] + r[a]) & (L - 1)]
        m = m & 255
        t += chr(m)
        e = (e + 1) & (L - 1)

    import base64
    return base64.b64encode(t.encode('latin-1')).decode()

result = generateArg("", "4xrihv8zb8tf1mfj")
print(f"Result:   {result}")
