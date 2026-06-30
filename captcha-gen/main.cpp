/*
 * CAPTCHA Token Service - Silent TCP Server
 * 
 * BUILD REQUIREMENTS:
 * - Download SQLite3 amalgamation from https://www.sqlite.org/download.html
 *   (sqlite3.c, sqlite3.h, sqlite3ext.h) and place in project directory
 * - Uncomment the sqlite3.c include line below for static embedding
 * - OR: install libsqlite3-dev and add "sqlite3" to CMakeLists.txt link libraries
 * 
 * USAGE:
 *   ./main &                    # Linux/Mac - run in background
 *   start /B main.exe           # Windows - run in background
 * 
 * PROTOCOL:
 *   TCP connect to port 7777
 *   Send any data (e.g., newline)
 *   Receive: base64 payload + newline on success
 *            "ERROR:..." + newline on failure
 *   Connection closes
 */

// ============================================================================
// SQLite3 - amalgamation files must be in project directory
// ============================================================================
#include "sqlite3.h"
// For static embedding without system sqlite3, uncomment:
// #define SQLITE_OMIT_LOAD_EXTENSION
// #define SQLITE_THREADSAFE 0
// #include "sqlite3.c"

// ============================================================================
// Cross-platform socket setup
// ============================================================================
#ifdef _WIN32
    #ifndef WIN32_LEAN_AND_MEAN
    #define WIN32_LEAN_AND_MEAN
    #endif
    #include <winsock2.h>
    #include <ws2tcpip.h>
    #pragma comment(lib, "ws2_32.lib")
    using socket_t = SOCKET;
    constexpr socket_t INVALID_SOCK = INVALID_SOCKET;
    inline void close_socket(socket_t s) { closesocket(s); }
    inline void socket_init() {
        WSADATA wsa;
        WSAStartup(MAKEWORD(2, 2), &wsa);
    }
    inline void socket_cleanup() { WSACleanup(); }
#else
    #include <sys/socket.h>
    #include <netinet/in.h>
    #include <arpa/inet.h>
    #include <unistd.h>
    #include <signal.h>
    using socket_t = int;
    constexpr socket_t INVALID_SOCK = -1;
    inline void close_socket(socket_t s) { close(s); }
    inline void socket_init() {
        signal(SIGPIPE, SIG_IGN);
    }
    inline void socket_cleanup() {}
#endif

#include <openssl/hmac.h>
#include <openssl/evp.h>
#include <openssl/bio.h>
#include <openssl/buffer.h>

#include <curl/curl.h>
#include <nlohmann/json.hpp>
#include <zlib.h>

#include <algorithm>
#include <array>
#include <cstdint>
#include <chrono>
#include <cstring>
#include <ctime>
#include <iomanip>
#include <iostream>
#include <map>
#include <mutex>
#include <random>
#include <set>
#include <sstream>
#include <string>
#include <vector>

using json = nlohmann::json;

// ============================================================================
// Configuration
// ============================================================================
static const std::string ACCESS_KEY = "LTAI5tSEBwYMwVKAQGpxmvTd";
static const std::string SECRET_KEY = "YSKfst7GaVkXwZYvVihJsKF9r89koz";
static const std::string SCENE_ID = "didk33e0";
static const int TCP_PORT = 7777;
static const int MAX_RETRIES = 3;
static const std::string DB_PATH = "tokens.sqlite";

// ============================================================================
// Utility: URL encoding
// ============================================================================
std::string url_encode(const std::string& s, const std::string& safe = "") {
    static const char hex[] = "0123456789ABCDEF";
    std::set<char> safe_set(safe.begin(), safe.end());
    std::string out;
    out.reserve(s.size() * 3);
    for (unsigned char c : s) {
        if (std::isalnum(c) || c == '-' || c == '_' || c == '.' || c == '~'
            || safe_set.count(static_cast<char>(c))) {
            out.push_back(static_cast<char>(c));
        } else {
            out.push_back('%');
            out.push_back(hex[c >> 4]);
            out.push_back(hex[c & 0x0F]);
        }
    }
    return out;
}

// ============================================================================
// Utility: Base64 encode
// ============================================================================
std::string base64_encode(const std::string& data) {
    BIO* b64 = BIO_new(BIO_f_base64());
    BIO* mem = BIO_new(BIO_s_mem());
    b64 = BIO_push(b64, mem);
    BIO_set_flags(b64, BIO_FLAGS_BASE64_NO_NL);
    BIO_write(b64, data.data(), static_cast<int>(data.size()));
    BIO_flush(b64);
    BUF_MEM* bptr = nullptr;
    BIO_get_mem_ptr(b64, &bptr);
    std::string result(bptr->data, bptr->length);
    BIO_free_all(b64);
    return result;
}

std::string base64_encode(const std::vector<uint8_t>& data) {
    std::string s(data.begin(), data.end());
    return base64_encode(s);
}

// ============================================================================
// Utility: Hex encode
// ============================================================================
std::string hex_encode(const uint8_t* data, size_t len) {
    std::ostringstream oss;
    oss << std::hex << std::setfill('0');
    for (size_t i = 0; i < len; ++i) {
        oss << std::setw(2) << static_cast<int>(data[i]);
    }
    return oss.str();
}

// ============================================================================
// HMAC-SHA1
// ============================================================================
std::string hmac_sha1(const std::string& key, const std::string& msg) {
    unsigned char digest[EVP_MAX_MD_SIZE];
    unsigned int digest_len = 0;
    HMAC(EVP_sha1(),
         key.data(), static_cast<int>(key.size()),
         reinterpret_cast<const unsigned char*>(msg.data()),
         msg.size(),
         digest, &digest_len);
    return std::string(reinterpret_cast<char*>(digest), digest_len);
}

// ============================================================================
// UUID v4 generation
// ============================================================================
std::string generate_uuid() {
    static std::random_device rd;
    static std::mt19937_64 gen(rd());
    std::uniform_int_distribution<uint64_t> dis(0, UINT64_MAX);
    uint64_t a = dis(gen);
    uint64_t b = dis(gen);
    a = (a & 0x0FFFFFFFFFFFFFFFULL) | 0x4000000000000000ULL;
    b = (b & 0x3FFFFFFFFFFFFFFFULL) | 0x8000000000000000ULL;

    std::ostringstream oss;
    oss << std::hex << std::setfill('0');
    oss << std::setw(8) << (a >> 32) << "-";
    oss << std::setw(4) << ((a >> 16) & 0xFFFF) << "-";
    oss << std::setw(4) << (a & 0xFFFF) << "-";
    oss << std::setw(4) << (b >> 48) << "-";
    oss << std::setw(12) << (b & 0x0000FFFFFFFFFFFFULL);
    return oss.str();
}

// ============================================================================
// UTC timestamp in ISO 8601 format
// ============================================================================
std::string get_timestamp_utc() {
    auto now = std::chrono::system_clock::now();
    std::time_t now_t = std::chrono::system_clock::to_time_t(now);
    std::tm tm_utc{};
#ifdef _WIN32
    gmtime_s(&tm_utc, &now_t);
#else
    gmtime_r(&now_t, &tm_utc);
#endif
    char buf[32];
    std::strftime(buf, sizeof(buf), "%Y-%m-%dT%H:%M:%SZ", &tm_utc);
    return std::string(buf);
}

// ============================================================================
// Current time in milliseconds since epoch
// ============================================================================
int64_t current_time_millis() {
    auto now = std::chrono::system_clock::now();
    auto ms = std::chrono::duration_cast<std::chrono::milliseconds>(
                  now.time_since_epoch());
    return static_cast<int64_t>(ms.count());
}

// ============================================================================
// Generate Aliyun signature
// ============================================================================
std::string generate_signature(const std::map<std::string, std::string>& params,
                               const std::string& secret_key) {
    std::vector<std::pair<std::string, std::string>> sorted_params(
        params.begin(), params.end());
    std::sort(sorted_params.begin(), sorted_params.end(),
              [](const auto& a, const auto& b) { return a.first < b.first; });

    std::string canonical;
    for (size_t i = 0; i < sorted_params.size(); ++i) {
        if (i > 0) canonical += "&";
        canonical += url_encode(sorted_params[i].first);
        canonical += "=";
        canonical += url_encode(sorted_params[i].second);
    }

    std::string string_to_sign = "POST&" + url_encode("/") + "&" +
                                 url_encode(canonical);

    std::string signing_key = secret_key + "&";
    std::string digest = hmac_sha1(signing_key, string_to_sign);

    return base64_encode(digest);
}

// ============================================================================
// libcurl write callback
// ============================================================================
static size_t write_callback(void* contents, size_t size, size_t nmemb,
                             void* userp) {
    size_t total_size = size * nmemb;
    std::string* str = static_cast<std::string*>(userp);
    str->append(static_cast<char*>(contents), total_size);
    return total_size;
}

// ============================================================================
// HTTP POST helper
// ============================================================================
std::string http_post(const std::string& url,
                      const std::string& body,
                      const std::vector<std::pair<std::string, std::string>>&
                          extra_headers = {}) {
    CURL* curl = curl_easy_init();
    if (!curl) {
        throw std::runtime_error("Failed to init curl");
    }
    std::string response_data;

    struct curl_slist* headers = nullptr;
    headers = curl_slist_append(headers,
        "Content-Type: application/x-www-form-urlencoded; charset=UTF-8");
    for (const auto& h : extra_headers) {
        headers = curl_slist_append(headers, (h.first + ": " + h.second).c_str());
    }

    curl_easy_setopt(curl, CURLOPT_URL, url.c_str());
    curl_easy_setopt(curl, CURLOPT_POST, 1L);
    curl_easy_setopt(curl, CURLOPT_POSTFIELDS, body.c_str());
    curl_easy_setopt(curl, CURLOPT_POSTFIELDSIZE, body.size());
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_callback);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &response_data);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT, 30L);
    curl_easy_setopt(curl, CURLOPT_SSL_VERIFYPEER, 1L);
    curl_easy_setopt(curl, CURLOPT_NOSIGNAL, 1L);

    CURLcode res = curl_easy_perform(curl);
    long http_code = 0;
    curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &http_code);

    curl_slist_free_all(headers);
    curl_easy_cleanup(curl);

    if (res != CURLE_OK) {
        throw std::runtime_error(std::string("curl error: ") + curl_easy_strerror(res));
    }
    if (http_code >= 400) {
        throw std::runtime_error("HTTP error: " + std::to_string(http_code));
    }
    return response_data;
}

// ============================================================================
// TokenManager - SQLite token pool management
// ============================================================================
class TokenManager {
    sqlite3* db_ = nullptr;
    std::mutex mtx_;

public:
    explicit TokenManager(const std::string& path) {
        int rc = sqlite3_open(path.c_str(), &db_);
        if (rc != SQLITE_OK) {
            std::cerr << "[ERROR] SQLite open failed: " 
                      << (db_ ? sqlite3_errmsg(db_) : "unknown") << "\n";
            if (db_) { sqlite3_close(db_); db_ = nullptr; }
            return;
        }
        // Enable WAL mode for better concurrency
        sqlite3_exec(db_, "PRAGMA journal_mode=WAL;", nullptr, nullptr, nullptr);
        // Verify table exists
        sqlite3_stmt* check = nullptr;
        const char* check_sql = "SELECT name FROM sqlite_master WHERE type='table' AND name='tokens';";
        if (sqlite3_prepare_v2(db_, check_sql, -1, &check, nullptr) == SQLITE_OK) {
            if (sqlite3_step(check) != SQLITE_ROW) {
                std::cerr << "[ERROR] Table 'tokens' not found in " << path << "\n";
            }
            sqlite3_finalize(check);
        }
    }

    ~TokenManager() {
        if (db_) sqlite3_close(db_);
    }

    bool is_open() const { return db_ != nullptr; }

    // Atomically get and remove one token
    std::string acquire() {
        std::lock_guard<std::mutex> lock(mtx_);
        if (!db_) return "";

        // Begin immediate transaction to prevent concurrent reads
        if (sqlite3_exec(db_, "BEGIN IMMEDIATE;", nullptr, nullptr, nullptr) != SQLITE_OK) {
            std::cerr << "[ERROR] SQLite begin transaction failed: " 
                      << sqlite3_errmsg(db_) << "\n";
            return "";
        }

        std::string token;
        sqlite3_stmt* sel = nullptr;
        const char* sel_sql = "SELECT token FROM tokens ORDER BY id LIMIT 1;";
        if (sqlite3_prepare_v2(db_, sel_sql, -1, &sel, nullptr) == SQLITE_OK) {
            if (sqlite3_step(sel) == SQLITE_ROW) {
                const char* t = reinterpret_cast<const char*>(sqlite3_column_text(sel, 0));
                if (t) token = t;
            }
            sqlite3_finalize(sel);
        }

        if (!token.empty()) {
            sqlite3_stmt* del = nullptr;
            const char* del_sql = "DELETE FROM tokens WHERE token = ?;";
            if (sqlite3_prepare_v2(db_, del_sql, -1, &del, nullptr) == SQLITE_OK) {
                sqlite3_bind_text(del, 1, token.c_str(), -1, SQLITE_TRANSIENT);
                sqlite3_step(del);
                sqlite3_finalize(del);
            }
        }

        sqlite3_exec(db_, "COMMIT;", nullptr, nullptr, nullptr);
        return token;
    }

    // Return count of remaining tokens
    int count() {
        std::lock_guard<std::mutex> lock(mtx_);
        if (!db_) return 0;
        int cnt = 0;
        sqlite3_stmt* stmt = nullptr;
        if (sqlite3_prepare_v2(db_, "SELECT COUNT(*) FROM tokens;", -1, &stmt, nullptr) == SQLITE_OK) {
            if (sqlite3_step(stmt) == SQLITE_ROW) {
                cnt = sqlite3_column_int(stmt, 0);
            }
            sqlite3_finalize(stmt);
        }
        return cnt;
    }
};

// ============================================================================
// InitCaptchaV3 - silent version, returns certify_id
// ============================================================================
std::string init_captcha_silent(const std::string& access_key,
                                const std::string& secret_key,
                                const std::string& scene_id) {
    std::map<std::string, std::string> params = {
        {"AccessKeyId", access_key},
        {"Action", "InitCaptchaV3"},
        {"Format", "JSON"},
        {"Language", "en"},
        {"Mode", "popup"},
        {"SceneId", scene_id},
        {"SignatureMethod", "HMAC-SHA1"},
        {"SignatureNonce", generate_uuid()},
        {"SignatureVersion", "1.0"},
        {"Timestamp", get_timestamp_utc()},
        {"UpLang", "true"},
        {"Version", "2023-03-05"},
    };

    std::string signature = generate_signature(params, secret_key);
    params["Signature"] = signature;

    std::string body;
    bool first = true;
    for (const auto& kv : params) {
        if (!first) body += "&";
        first = false;
        body += url_encode(kv.first) + "=" + url_encode(kv.second);
    }

    std::string response = http_post(
        "https://no8xfe.captcha-open-southeast.aliyuncs.com/", body);

    auto resp_json = json::parse(response);
    return resp_json.value("CertifyId", "");
}

// ============================================================================
// Generate arg field
// ============================================================================
std::string generate_arg(const std::string& certify_id,
                         const std::string& constant = "4xrihv8zb8tf1mfj") {
    std::string encoded = url_encode(certify_id);
    std::string o;
    for (size_t i = 0; i < encoded.size();) {
        if (encoded[i] == '%' && i + 2 < encoded.size()) {
            int hi = std::stoi(encoded.substr(i + 1, 2), nullptr, 16);
            o.push_back(static_cast<char>(hi));
            i += 3;
        } else {
            o.push_back(encoded[i]);
            i += 1;
        }
    }

    std::string n = constant;
    std::vector<int> r = {
        32,50,10,51,6,44,37,16,46,11,62,19,43,25,23,30,
        60,33,53,34,7,26,12,48,5,2,20,4,61,13,47,49,
        18,29,27,22,1,17,39,56,41,38,55,31,15,58,52,40,
        8,57,45,35,59,36,42,54,63,3,24,28,14,9,0,21
    };
    int rlen = static_cast<int>(r.size());

    int i = 0, j = 0;
    while (i < rlen) {
        j = (((i + j + r[i] + r[j]) >> 1) + n[i % n.size()]) & (rlen - 1);
        if (i != j) {
            r[i] ^= r[j];
            r[j] ^= r[i];
            r[i] ^= r[j];
        }
        i += 1;
    }

    std::string t;
    int e = 0, a = 0;
    for (size_t idx = 0; idx < o.size(); ++idx) {
        a = ((e ^ a) + (r[e] ^ r[a])) & (rlen - 1);
        if (e != a) {
            r[e] ^= r[a];
            r[a] ^= r[e];
            r[e] ^= r[a];
        }
        int m = static_cast<unsigned char>(o[idx]);
        m = m + e + r[e] - a - r[a];
        m = m ^ (r[e] + r[a]);
        m = m ^ r[(r[e] + r[a]) & (rlen - 1)];
        m = m & 255;
        t.push_back(static_cast<char>(m));
        e = (e + 1) & (rlen - 1);
    }

    return base64_encode(t);
}

// ============================================================================
// ali_hash
// ============================================================================
std::string ali_hash(const std::string& input_str, const std::string& salt_str) {
    std::string o = input_str;
    std::string r = salt_str;
    int a_len = static_cast<int>(o.size());
    int m = static_cast<int>(r.size());

    std::vector<int> e;
    for (int i = 0; i < 16; ++i) {
        e.push_back((i << 4) + (i % 16));
    }
    int f = static_cast<int>(e.size());

    int i = 0, j = 0;
    while (i < f) {
        j = (((i + j + e[i] + e[j]) >> 1) + r[i % m]) & (f - 1);
        std::swap(e[i], e[j]);
        i += 1;
    }

    int idx = 0, p = 0, q = 0;
    while (idx < a_len) {
        q = ((p ^ q) + (e[p] ^ e[q])) & (f - 1);
        std::swap(e[p], e[q]);
        int C = static_cast<unsigned char>(o[idx]);
        C = (C + p + q) ^ e[p] ^ e[q];
        C = C & 255;
        e[p] = C;
        p = (p + 1) & (f - 1);
        idx += 1;
    }

    for (int step = 0; step < 2 * f; ++step) {
        int pos = step % f;
        if (pos != 0) {
            e[pos] ^= e[pos - 1];
        } else {
            e[0] ^= e[f - 1];
        }
    }

    std::ostringstream result;
    for (int b : e) {
        result << std::hex << std::setw(2) << std::setfill('0') << (b & 0xFF);
    }
    return result.str();
}

// ============================================================================
// encrypt
// ============================================================================
std::string encrypt(const std::vector<uint8_t>& plaintext_bytes) {
    std::string o(plaintext_bytes.begin(), plaintext_bytes.end());

    std::string n = "3e627e1b4c63f913";
    std::vector<int> r = {
        32,50,10,51,6,44,37,16,46,11,62,19,43,25,23,30,
        60,33,53,34,7,26,12,48,5,2,20,4,61,13,47,49,
        18,29,27,22,1,17,39,56,41,38,55,31,15,58,52,40,
        8,57,45,35,59,36,42,54,63,3,24,28,14,9,0,21
    };
    int rlen = static_cast<int>(r.size());

    int o_ksa = 0, t_ksa = 0;
    while (o_ksa < rlen) {
        t_ksa = (((o_ksa + t_ksa + r[o_ksa] + r[t_ksa]) >> 1)
                 + n[o_ksa % n.size()]) & (rlen - 1);
        if (o_ksa != t_ksa) {
            r[o_ksa] ^= r[t_ksa];
            r[t_ksa] ^= r[o_ksa];
            r[o_ksa] ^= r[t_ksa];
        }
        o_ksa += 1;
    }

    std::string t;
    int n_prga = 0, e = 0, a = 0;
    while (n_prga < static_cast<int>(o.size())) {
        a = ((e ^ a) + (r[e] ^ r[a])) & (rlen - 1);
        if (e != a) {
            r[e] ^= r[a];
            r[a] ^= r[e];
            r[e] ^= r[a];
        }
        int m = static_cast<unsigned char>(o[n_prga]);
        m = m + e + r[e];
        m = m - a - r[a];
        m = m ^ (r[e] + r[a]);
        m = m ^ r[(r[e] + r[a]) & (rlen - 1)];
        m = m & 255;
        t.push_back(static_cast<char>(m));
        e = (e + 1) & (rlen - 1);
        n_prga += 1;
    }

    return base64_encode(t);
}

// ============================================================================
// zlib compress
// ============================================================================
std::vector<uint8_t> zlib_compress(const std::string& data) {
    z_stream zs{};
    if (deflateInit(&zs, Z_DEFAULT_COMPRESSION) != Z_OK) {
        throw std::runtime_error("deflateInit failed");
    }
    zs.next_in = reinterpret_cast<Bytef*>(const_cast<char*>(data.data()));
    zs.avail_in = static_cast<uInt>(data.size());

    std::vector<uint8_t> out;
    out.resize(data.size() + 1024);
    int ret = Z_OK;
    do {
        if (zs.total_out >= out.size()) {
            out.resize(out.size() * 2);
        }
        zs.next_out = reinterpret_cast<Bytef*>(out.data() + zs.total_out);
        zs.avail_out = static_cast<uInt>(out.size() - zs.total_out);
        ret = deflate(&zs, Z_FINISH);
    } while (ret == Z_OK);

    out.resize(zs.total_out);
    deflateEnd(&zs);

    if (ret != Z_STREAM_END) {
        throw std::runtime_error("deflate failed");
    }
    return out;
}

// ============================================================================
// Build captcha data value (parts 2-7 combined)
// ============================================================================
std::string build_captcha_data(const std::string& certify_id) {
    std::string arg_value = generate_arg(certify_id);

    int64_t ct = current_time_millis();
    json track_json = {
        {"TrackList", {
            {"mc", ""}, {"tc", ""}, {"mu", ""}, {"te", ""},
            {"mp", ""}, {"tmv", ""}, {"ks", ""}, {"fi", ""},
            {"startTime", ct}
        }},
        {"TrackStartTime", ct},
        {"VerifyTime", ct + 300},
        {"arg", arg_value}
    };
    std::string json_str = track_json.dump();

    std::string hash_result = ali_hash(json_str, "0000");
    std::string combined = hash_result + json_str;

    std::vector<uint8_t> compressed = zlib_compress(combined);
    std::string final_data_b64 = base64_encode(compressed);
    std::string final_value = encrypt(
        std::vector<uint8_t>(final_data_b64.begin(), final_data_b64.end()));

    return final_value;
}

// ============================================================================
// VerifyCaptchaV3 - silent version, returns base64 payload or empty on fail
// ============================================================================
std::string verify_captcha_silent(const std::string& access_key,
                                  const std::string& secret_key,
                                  const std::string& scene_id,
                                  const std::string& certify_id,
                                  const std::string& data_value,
                                  const std::string& device_token) {
    json captcha_verify_param = {
        {"sceneId", scene_id},
        {"certifyId", certify_id},
        {"deviceToken", device_token},
        {"data", data_value}
    };
    std::string captcha_verify_param_str = captcha_verify_param.dump();

    std::map<std::string, std::string> params = {
        {"AccessKeyId", access_key},
        {"Action", "VerifyCaptchaV3"},
        {"Format", "JSON"},
        {"SignatureMethod", "HMAC-SHA1"},
        {"SignatureVersion", "1.0"},
        {"Timestamp", get_timestamp_utc()},
        {"Version", "2023-03-05"},
        {"SceneId", scene_id},
        {"CertifyId", certify_id},
        {"CaptchaVerifyParam", captcha_verify_param_str},
        {"SignatureNonce", generate_uuid()},
    };

    std::string signature = generate_signature(params, secret_key);
    params["Signature"] = signature;

    std::string body;
    bool first = true;
    for (const auto& kv : params) {
        if (!first) body += "&";
        first = false;
        body += url_encode(kv.first) + "=" + url_encode(kv.second);
    }

    std::string response = http_post(
        "https://no8xfe-verify.captcha-open-southeast.aliyuncs.com/", body,
        {{"Referer", ""}});

    auto resp_json = json::parse(response);

    if (resp_json.value("Success", false)) {
        auto result = resp_json.value("Result", json::object());
        if (result.value("VerifyResult", false)) {
            std::string security_token = result.value("securityToken", "");
            std::string resp_certify_id = result.value("certifyId", "");

            if (!security_token.empty() && !resp_certify_id.empty()) {
                json final_payload = {
                    {"certifyId", resp_certify_id},
                    {"sceneId", scene_id},
                    {"isSign", true},
                    {"securityToken", security_token}
                };
                return base64_encode(final_payload.dump());
            }
        }
    }

    return "";
}

// ============================================================================
// Full captcha flow with single token - returns payload or empty
// ============================================================================
std::string run_captcha_flow(const std::string& device_token) {
    std::string certify_id = init_captcha_silent(ACCESS_KEY, SECRET_KEY, SCENE_ID);
    if (certify_id.empty()) {
        std::cerr << "[ERROR] InitCaptcha returned empty CertifyId\n";
        return "";
    }

    std::string data_value = build_captcha_data(certify_id);
    return verify_captcha_silent(ACCESS_KEY, SECRET_KEY, SCENE_ID,
                                 certify_id, data_value, device_token);
}

// ============================================================================
// TCP Server
// ============================================================================
class TCPServer {
    socket_t listen_fd_ = INVALID_SOCK;
    int port_;

public:
    explicit TCPServer(int port) : port_(port) {}

    bool start() {
        listen_fd_ = socket(AF_INET, SOCK_STREAM, 0);
        if (listen_fd_ == INVALID_SOCK) {
            std::cerr << "[ERROR] socket() failed\n";
            return false;
        }

        int opt = 1;
        setsockopt(listen_fd_, SOL_SOCKET, SO_REUSEADDR,
                   reinterpret_cast<const char*>(&opt), sizeof(opt));

        struct sockaddr_in addr{};
        std::memset(&addr, 0, sizeof(addr));
        addr.sin_family = AF_INET;
        addr.sin_addr.s_addr = INADDR_ANY;
        addr.sin_port = htons(static_cast<uint16_t>(port_));

        if (bind(listen_fd_, reinterpret_cast<struct sockaddr*>(&addr),
                 sizeof(addr)) < 0) {
            std::cerr << "[ERROR] bind() failed on port " << port_ << "\n";
            close_socket(listen_fd_);
            listen_fd_ = INVALID_SOCK;
            return false;
        }

        if (listen(listen_fd_, 16) < 0) {
            std::cerr << "[ERROR] listen() failed\n";
            close_socket(listen_fd_);
            listen_fd_ = INVALID_SOCK;
            return false;
        }

        return true;
    }

    socket_t accept_one() {
        struct sockaddr_in client_addr{};
        socklen_t client_len = sizeof(client_addr);
        return accept(listen_fd_,
                      reinterpret_cast<struct sockaddr*>(&client_addr),
                      &client_len);
    }

    void stop() {
        if (listen_fd_ != INVALID_SOCK) {
            close_socket(listen_fd_);
            listen_fd_ = INVALID_SOCK;
        }
    }

    ~TCPServer() { stop(); }
};

// ============================================================================
// Handle a single TCP client connection
// ============================================================================
void handle_client(socket_t client_fd, TokenManager& token_mgr) {
    // Read client request (drain any incoming data with short timeout)
    char buf[256];
#ifdef _WIN32
    DWORD timeout_ms = 1000;
    setsockopt(client_fd, SOL_SOCKET, SO_RCVTIMEO,
               reinterpret_cast<const char*>(&timeout_ms), sizeof(timeout_ms));
#else
    struct timeval tv;
    tv.tv_sec = 1;
    tv.tv_usec = 0;
    setsockopt(client_fd, SOL_SOCKET, SO_RCVTIMEO, &tv, sizeof(tv));
#endif
    recv(client_fd, buf, sizeof(buf), 0);

    // Attempt captcha with retries
    std::string payload;
    for (int attempt = 0; attempt < MAX_RETRIES; ++attempt) {
        std::string token = token_mgr.acquire();
        if (token.empty()) {
            payload = "ERROR:NO_TOKENS";
            break;
        }

        try {
            payload = run_captcha_flow(token);
            if (!payload.empty()) {
                break;  // Success
            }
            // Token failed, will retry with next token
        } catch (const std::exception& e) {
            std::cerr << "[ERROR] Attempt " << (attempt + 1) 
                      << " exception: " << e.what() << "\n";
        }
    }

    if (payload.empty()) {
        payload = "ERROR:ALL_RETRIES_FAILED";
    }

    // Send response with newline terminator
    payload += "\n";
    send(client_fd, payload.c_str(), static_cast<int>(payload.size()), 0);
    close_socket(client_fd);
}

// ============================================================================
// Main entry point
// ============================================================================
int main() {
    // Suppress stdout completely for silent background operation
    std::cout.setstate(std::ios_base::failbit);

    socket_init();
    curl_global_init(CURL_GLOBAL_ALL);

    // Open token database
    TokenManager token_mgr(DB_PATH);
    if (!token_mgr.is_open()) {
        std::cerr << "[ERROR] Cannot open " << DB_PATH << "\n";
        curl_global_cleanup();
        socket_cleanup();
        return 1;
    }

    int available = token_mgr.count();
    if (available == 0) {
        std::cerr << "[ERROR] No tokens available in " << DB_PATH << "\n";
        curl_global_cleanup();
        socket_cleanup();
        return 1;
    }

    // Start TCP server
    TCPServer server(TCP_PORT);
    if (!server.start()) {
        curl_global_cleanup();
        socket_cleanup();
        return 1;
    }

    // Main event loop - silent, handles one client at a time
    while (true) {
        socket_t client = server.accept_one();
        if (client == INVALID_SOCK) {
            continue;
        }
        handle_client(client, token_mgr);
    }

    // Unreachable in normal operation
    curl_global_cleanup();
    socket_cleanup();
    return 0;
}
