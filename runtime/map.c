#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#include <stdio.h>

// Key kinds:
// 0: int64 (signed)
// 1: uint64 / bool
// 2: string (struct { i8* data; i64 len })
// 3: pointer

typedef struct {
    int64_t key_size;
    int64_t val_size;
    int32_t key_kind;
    int64_t len;
    int64_t cap;
    uint8_t* keys;
    uint8_t* vals;
} gominic_map;

static int key_equal(const gominic_map* m, const uint8_t* a, const uint8_t* b) {
    switch (m->key_kind) {
    case 0: {
        int64_t va = *(const int64_t*)a;
        int64_t vb = *(const int64_t*)b;
        return va == vb;
    }
    case 1: {
        uint64_t va = *(const uint64_t*)a;
        uint64_t vb = *(const uint64_t*)b;
        return va == vb;
    }
    case 2: {
        const char* sa = *(char* const*)a;
        int64_t la = *(const int64_t*)(a + sizeof(void*));
        const char* sb = *(char* const*)b;
        int64_t lb = *(const int64_t*)(b + sizeof(void*));
        if (la != lb) return 0;
        if (sa == sb) return 1;
        if (sa == NULL || sb == NULL) return 0;
        return memcmp(sa, sb, (size_t)la) == 0;
    }
    case 3: {
        const void* pa = *(void* const*)a;
        const void* pb = *(void* const*)b;
        return pa == pb;
    }
    default:
        return 0;
    }
}

static void ensure_capacity(gominic_map* m, int64_t need) {
    if (need <= m->cap) return;
    int64_t new_cap = m->cap == 0 ? 8 : m->cap * 2;
    while (new_cap < need) new_cap *= 2;
    uint8_t* new_keys = (uint8_t*)realloc(m->keys, (size_t)(new_cap * m->key_size));
    uint8_t* new_vals = (uint8_t*)realloc(m->vals, (size_t)(new_cap * m->val_size));
    if (!new_keys || !new_vals) {
        abort();
    }
    // zero newly allocated tails
    if (new_cap > m->cap) {
        int64_t old_cap = m->cap;
        memset(new_keys + old_cap * m->key_size, 0, (size_t)((new_cap - old_cap) * m->key_size));
        memset(new_vals + old_cap * m->val_size, 0, (size_t)((new_cap - old_cap) * m->val_size));
    }
    m->keys = new_keys;
    m->vals = new_vals;
    m->cap = new_cap;
}

// gominic_map_new allocates a simple linear map.
// keyKind: 0 int64, 1 uint64, 2 string, 3 pointer.
void* gominic_map_new(int64_t keySize, int64_t valSize, int32_t keyKind) {
    gominic_map* m = (gominic_map*)calloc(1, sizeof(gominic_map));
    if (!m) return NULL;
    m->key_size = keySize;
    m->val_size = valSize;
    m->key_kind = keyKind;
    m->len = 0;
    m->cap = 0;
    m->keys = NULL;
    m->vals = NULL;
    return m;
}

// gominic_map_set inserts or replaces.
void gominic_map_set(void* mapPtr, void* keyPtr, void* valPtr) {
    if (!mapPtr) return;
    gominic_map* m = (gominic_map*)mapPtr;
    // search existing
    for (int64_t i = 0; i < m->len; i++) {
        uint8_t* k = m->keys + i * m->key_size;
        if (key_equal(m, k, (const uint8_t*)keyPtr)) {
            memcpy(m->vals + i * m->val_size, valPtr, (size_t)m->val_size);
            return;
        }
    }
    ensure_capacity(m, m->len + 1);
    memcpy(m->keys + m->len * m->key_size, keyPtr, (size_t)m->key_size);
    memcpy(m->vals + m->len * m->val_size, valPtr, (size_t)m->val_size);
    m->len++;
}

// gominic_map_get returns 1 if found; writes value to valPtr.
int gominic_map_get(void* mapPtr, void* keyPtr, void* valPtr) {
    if (!mapPtr) return 0;
    gominic_map* m = (gominic_map*)mapPtr;
    for (int64_t i = 0; i < m->len; i++) {
        uint8_t* k = m->keys + i * m->key_size;
        if (key_equal(m, k, (const uint8_t*)keyPtr)) {
            memcpy(valPtr, m->vals + i * m->val_size, (size_t)m->val_size);
            return 1;
        }
    }
    // not found: zero val
    memset(valPtr, 0, (size_t)m->val_size);
    return 0;
}

int64_t gominic_map_len(void* mapPtr) {
    if (!mapPtr) return 0;
    gominic_map* m = (gominic_map*)mapPtr;
    return m->len;
}
