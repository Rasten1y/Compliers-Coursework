#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

typedef struct {
	const char *data;
	int64_t len;
} gominic_string;

// Minimal runtime helpers callable from generated code.
void gominic_print(const char *s, int64_t len) {
	fwrite(s, 1, (size_t)len, stdout);
}

void gominic_printInt(int64_t v) {
	printf("%lld", (long long)v);
}

void gominic_println(void) {
	fputc('\n', stdout);
}

// Simple heap allocation for slices: returns void*.
void *gominic_makeSlice(int64_t len, int64_t cap, int64_t elemSize) {
	// allocate cap elements; if cap < len, fall back to len to avoid UB
	int64_t count = cap;
	if (cap < len) {
		count = len;
	}
	if (count <= 0 || elemSize <= 0) {
		count = 1;
		elemSize = 1;
	}
	size_t bytes = (size_t)(count * elemSize);
	void *p = calloc(1, bytes);
	return p;
}

// memcpy wrapper for aggregate copies.
void gominic_memcpy(void *dst, const void *src, int64_t n) {
	if (n <= 0) {
		return;
	}
	memcpy(dst, src, (size_t)n);
}

// Abort helper for runtime checks.
void gominic_abort(void) {
	fprintf(stderr, "gominic runtime abort\n");
	abort();
}

// String equality: compares lengths then bytes.
int gominic_str_eq(const char *a, int64_t alen, const char *b, int64_t blen) {
	if (alen != blen) {
		return 0;
	}
	if (alen == 0) {
		return 1;
	}
	return memcmp(a, b, (size_t)alen) == 0;
}

gominic_string gominic_str_concat(gominic_string a, gominic_string b) {
	int64_t total = a.len + b.len;
	char *buf = malloc((size_t)(total + 1));
	if (buf == NULL) {
		gominic_string empty = { "", 0 };
		return empty;
	}
	if (a.len > 0) {
		memcpy(buf, a.data, (size_t)a.len);
	}
	if (b.len > 0) {
		memcpy(buf + a.len, b.data, (size_t)b.len);
	}
	buf[total] = 0;
	gominic_string res;
	res.data = buf;
	res.len = total;
	return res;
}
