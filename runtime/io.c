// Minimal I/O helpers for self-hosted builds.
// These wrappers avoid dependence on Go's os package by exposing argc/argv and
// basic file writing to the generated code.

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

typedef long long i64;

typedef struct {
    char *data;
    i64 len;
} gominic_string;

// Return command-line argument count.
i64 gominic_argc(void) {
#ifdef _WIN32
    extern int __argc;
    return (i64)__argc;
#else
    extern int __argc; // some libc expose this
    extern int argc;
    if (&__argc != NULL) {
        return (i64)__argc;
    }
    return (i64)argc;
#endif
}

// Return argv[idx] as a {data,len} pair via sret pointer. If out of range, returns {NULL, 0}.
void gominic_argv(gominic_string *out, i64 idx) {
    gominic_string res;
    res.data = NULL;
    res.len = 0;
#ifdef _WIN32
    extern char **__argv;
    extern int __argc;
    if (idx < 0 || idx >= (i64)__argc) {
        *out = res;
        return;
    }
    res.data = __argv[idx];
    res.len = (i64)strlen(__argv[idx]);
#else
    extern char **__argv; // glibc
    extern char **argv;
    extern int __argc;
    extern int argc;
    char **use_argv = NULL;
    i64 use_argc = 0;
    if (&__argv != NULL) {
        use_argv = __argv;
        use_argc = (i64)__argc;
    } else {
        use_argv = argv;
        use_argc = (i64)argc;
    }
    if (idx < 0 || idx >= use_argc) {
        *out = res;
        return;
    }
    res.data = use_argv[idx];
    res.len = (i64)strlen(use_argv[idx]);
#endif
    *out = res;
}

// Write a buffer to file path (not null-terminated, provided length). Returns 1 on success, 0 otherwise.
// Mode is always binary write/overwrite.
int gominic_write_file(const char *path, i64 path_len, const char *data, i64 data_len) {
    if (path == NULL || data == NULL || path_len <= 0) {
        return 0;
    }
    char *zpath = (char *)malloc((size_t)path_len + 1);
    if (zpath == NULL) {
        return 0;
    }
    memcpy(zpath, path, (size_t)path_len);
    zpath[path_len] = '\0';

#ifdef _WIN32
    FILE *f = fopen(zpath, "wb");
#else
    FILE *f = fopen(zpath, "wb");
#endif
    free(zpath);
    if (f == NULL) {
        return 0;
    }
    size_t written = fwrite(data, 1, (size_t)data_len, f);
    fclose(f);
    return written == (size_t)data_len;
}
