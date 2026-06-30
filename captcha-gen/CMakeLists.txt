cmake_minimum_required(VERSION 3.14)
project(CoreVerCap3)
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)
include(FetchContent)
# ── Static linking (single statically linked binary) ──────────────────────
set(BUILD_SHARED_LIBS OFF CACHE BOOL "" FORCE)
if(MSVC)
    set(CMAKE_MSVC_RUNTIME_LIBRARY "MultiThreaded$<$<CONFIG:Debug>:Debug>")
else()
    set(CMAKE_EXE_LINKER_FLAGS "${CMAKE_EXE_LINKER_FLAGS} -static -static-libgcc -static-libstdc++")
endif()
# ── Threads (pthread on Linux/macOS) ──────────────────────────────────────
set(THREADS_PREFER_PTHREAD_FLAG ON)
find_package(Threads REQUIRED)
# ── nlohmann/json ─────────────────────────────────────────────────────────
FetchContent_Declare(
    nlohmann_json
    GIT_REPOSITORY https://github.com/nlohmann/json.git
    GIT_TAG        v3.11.3
)
FetchContent_MakeAvailable(nlohmann_json)
# ── ZLIB (Auto-downloaded) ────────────────────────────────────────────────
set(ZLIB_BUILD_EXAMPLES OFF CACHE BOOL "" FORCE)
FetchContent_Declare(
    zlib
    GIT_REPOSITORY https://github.com/madler/zlib.git
    GIT_TAG        v1.3.1
)
FetchContent_MakeAvailable(zlib)
# ── SQLite3 amalgamation (self-contained, no system dependency) ───────────
set(SQLITE3_DIR "${CMAKE_BINARY_DIR}/sqlite3-src")
if(NOT DEFINED SQLITE3_URL)
    set(SQLITE3_URL "https://www.sqlite.org/2024/sqlite-amalgamation-3460100.zip")
endif()
if(NOT EXISTS "${SQLITE3_DIR}/sqlite3.c")
    message(STATUS "Downloading SQLite3 amalgamation...")
    file(DOWNLOAD ${SQLITE3_URL}
         "${CMAKE_BINARY_DIR}/sqlite3.zip"
         STATUS DL_STATUS)
    list(GET DL_STATUS 0 DL_ERROR)
    if(DL_ERROR)
        message(FATAL_ERROR "Failed to download SQLite3: ${DL_STATUS}")
    endif()
    file(ARCHIVE_EXTRACT INPUT "${CMAKE_BINARY_DIR}/sqlite3.zip"
         DESTINATION "${SQLITE3_DIR}")
    file(REMOVE "${CMAKE_BINARY_DIR}/sqlite3.zip")
endif()
file(GLOB SQLITE3_SUBDIR "${SQLITE3_DIR}/sqlite-amalgamation-*")
if(NOT SQLITE3_SUBDIR)
    message(FATAL_ERROR "SQLite3 amalgamation directory not found")
endif()
add_library(sqlite3 STATIC "${SQLITE3_SUBDIR}/sqlite3.c")
target_include_directories(sqlite3 PUBLIC "${SQLITE3_SUBDIR}")
if(UNIX)
    target_link_libraries(sqlite3 PUBLIC dl)
endif()
# ── System libraries: curl, openssl (statically linked) ───────────────────
# A statically-built system libcurl pulls in a chain of transitive static
# deps (nghttp2, libssh, OpenLDAP, Kerberos/GSSAPI, librtmp, libpsl, zstd,
# brotli, etc.) that find_package(CURL) does NOT resolve — it only links
# libcurl.a itself, leaving those symbols undefined at link time. pkg-config
# --static is what correctly expands the full dependency closure, so it is
# used here instead on Unix. Requires `pkg-config` and the matching -dev
# packages (libcurl4-openssl-dev pulls these in via build-deps) to be
# installed on the build machine.
if(UNIX AND NOT APPLE)
    find_package(PkgConfig REQUIRED)
    pkg_check_modules(CURL_STATIC REQUIRED IMPORTED_TARGET libcurl)
    pkg_check_modules(OPENSSL_STATIC REQUIRED IMPORTED_TARGET openssl)

    # Ubuntu/Debian never ship static (.a) builds of Kerberos/GSSAPI
    # libraries — only shared (.so) ones — even though the system libcurl
    # was built with GSSAPI support and pkg-config --static asks for them.
    # Strip those specific flags from the static link line; they're added
    # back below as an explicit dynamic link instead.
    list(FILTER CURL_STATIC_STATIC_LDFLAGS EXCLUDE REGEX "gssapi_krb5|krb5|k5crypto|com_err|gssrpc|kdb5")
else()
    find_package(CURL REQUIRED)
    find_package(OpenSSL REQUIRED)
endif()
# ── Main executable ───────────────────────────────────────────────────────
add_executable(main main.cpp)
# Expose the fetched ZLIB headers to main.cpp
target_include_directories(main PRIVATE ${zlib_SOURCE_DIR} ${zlib_BINARY_DIR})
target_link_libraries(main
    nlohmann_json::nlohmann_json
    sqlite3
    zlibstatic          # The statically built target from FetchContent
    Threads::Threads
)
if(UNIX AND NOT APPLE)
    # pkg-config's --static output already includes every transitive static
    # dependency curl/openssl need (nghttp2, ssh, ldap, rtmp, psl, zstd,
    # brotli, etc.) in the correct link order. Kerberos/GSSAPI was filtered
    # out above since no static archive exists on Ubuntu — it's linked
    # dynamically here instead, the only non-static part of the binary.
    target_compile_options(main PRIVATE ${CURL_STATIC_CFLAGS} ${OPENSSL_STATIC_CFLAGS})
    target_link_libraries(main
        ${CURL_STATIC_STATIC_LDFLAGS}
        ${OPENSSL_STATIC_STATIC_LDFLAGS}
        dl
        -Wl,-Bdynamic
        gssapi_krb5
        krb5
        k5crypto
        com_err
        -Wl,-Bstatic
    )
else()
    target_link_libraries(main CURL::libcurl OpenSSL::SSL OpenSSL::Crypto)
endif()
# Windows needs ws2_32/crypt32/bcrypt for Winsock and TLS/crypto support
# pulled in transitively by a static libcurl/OpenSSL.
if(WIN32)
    target_link_libraries(main ws2_32 crypt32 bcrypt)
endif()
