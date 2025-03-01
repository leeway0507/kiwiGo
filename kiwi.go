package kiwigo

/*
#cgo CFLAGS: -I ./kiwi/include
#cgo LDFLAGS: -L ./kiwi/lib -lkiwi
#include <stdlib.h>
#include <string.h>
#include <stdint.h> // for uintptr_t

#include <kiwi/capi.h>
*/
import "C"
import (
	"log"
	"time"
)

func Version() {
	start := time.Now()
	version := C.GoString(C.kiwi_version())
	log.Println("Kiwi version:", version)
	log.Println(time.Since(start).String())
}

type AnalyzeOption int

const (
	KIWI_MATCH_URL                  AnalyzeOption = C.KIWI_MATCH_URL
	KIWI_MATCH_EMAIL                AnalyzeOption = C.KIWI_MATCH_EMAIL
	KIWI_MATCH_HASHTAG              AnalyzeOption = C.KIWI_MATCH_HASHTAG
	KIWI_MATCH_MENTION              AnalyzeOption = C.KIWI_MATCH_MENTION
	KIWI_MATCH_ALL                  AnalyzeOption = C.KIWI_MATCH_ALL
	KIWI_MATCH_NORMALIZE_CODA       AnalyzeOption = C.KIWI_MATCH_NORMALIZE_CODA
	KIWI_MATCH_ALL_WITH_NORMALIZING AnalyzeOption = C.KIWI_MATCH_ALL_WITH_NORMALIZING
)

// BuildOption is a bitwise OR of the KiwiBuildOption values.
type BuildOption int

const (
	KIWI_BUILD_LOAD_DEFAULT_DICT   BuildOption = C.KIWI_BUILD_LOAD_DEFAULT_DICT
	KIWI_BUILD_INTEGRATE_ALLOMORPH BuildOption = C.KIWI_BUILD_INTEGRATE_ALLOMORPH
	KIWI_BUILD_DEFAULT             BuildOption = C.KIWI_BUILD_DEFAULT
)

// Kiwi is a wrapper for the kiwi C library.
type Kiwi struct {
	handler C.kiwi_h
}

type TokenResult struct {
	Tokens []TokenInfo
	Score  float32
}

// TokenInfo returns the token info for the given token(Str).
type TokenInfo struct {
	// Position is the index of this token appears in the original text.
	Position int

	// Tag represents a type of this token (e.g. VV, NNG, ...).
	Tag POSType

	// Form is the actual string of this token.
	Form string
}

func (k *Kiwi) Analyze(text string, topN int, options AnalyzeOption) ([]TokenResult, error) {
	kiwiResH := C.kiwi_analyze(k.handler, C.CString(text), C.int(topN), C.int(options), nil, nil)

	defer C.kiwi_res_close(kiwiResH)

	resSize := int(C.kiwi_res_size(kiwiResH))
	res := make([]TokenResult, resSize)

	for i := 0; i < resSize; i++ {
		tokens := make([]TokenInfo, int(C.kiwi_res_word_num(kiwiResH, C.int(i))))

		for j := 0; j < len(tokens); j++ {
			pos, err := ParsePOSType(C.GoString(C.kiwi_res_tag(kiwiResH, C.int(i), C.int(j))))
			if err != nil {
				return nil, err
			}
			tokens[j] = TokenInfo{
				Form:     C.GoString(C.kiwi_res_form(kiwiResH, C.int(i), C.int(j))),
				Tag:      pos,
				Position: int(C.kiwi_res_position(kiwiResH, C.int(i), C.int(j))),
			}
		}

		res[i] = TokenResult{
			Tokens: tokens,
			Score:  float32(C.kiwi_res_prob(kiwiResH, C.int(i))),
		}
	}

	return res, nil
}
func (k *Kiwi) Close() int {
	if k.handler != nil {
		out := int(C.kiwi_close(k.handler))
		k.handler = nil
		return out
	}
	return 0
}

// KiwiBuilder is a wrapper for the kiwi C library.
type KiwiBuilder struct {
	handler C.kiwi_builder_h
}

// NewBuilder returns a new KiwiBuilder instance.
// Don't forget to call Close after this.
func NewBuilder(modelPath string, numThread int, options BuildOption) *KiwiBuilder {
	return &KiwiBuilder{
		handler: C.kiwi_builder_init(C.CString(modelPath), C.int(numThread), C.int(options)),
	}
}

// AddWord set custom word with word, pos, score.
func (kb *KiwiBuilder) AddWord(word string, pos POSType, score float32) int {
	return int(C.kiwi_builder_add_word(kb.handler, C.CString(word), C.CString(string(pos)), C.float(score)))
}

// LoadDict loads user dict with dict file path.
func (kb *KiwiBuilder) LoadDict(dictPath string) int {
	return int(C.kiwi_builder_load_dict(kb.handler, C.CString(dictPath)))
}

// Build creates kiwi instance with user word etc.
func (kb *KiwiBuilder) Build() *Kiwi {
	h := C.kiwi_builder_build(kb.handler, nil, C.float(0))
	defer kb.Close()
	return &Kiwi{
		handler: h,
	}
}

func (kb *KiwiBuilder) Close() int {
	if kb.handler != nil {
		out := int(C.kiwi_builder_close(kb.handler))
		kb.handler = nil
		return out
	}
	return 0
}
