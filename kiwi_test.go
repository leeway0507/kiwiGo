package kiwigo

import (
	"fmt"
	"testing"
)

func TestKiwi(t *testing.T) {

	t.Run("version", func(t *testing.T) {
		Version()
	})

	t.Run("Analyze", func(t *testing.T) {
		func() {
			kb := NewBuilder("./models/base", 1 /*=numThread*/, KIWI_BUILD_INTEGRATE_ALLOMORPH /*=options*/)
			kb.AddWord("코딩냄비", "NNP", 0)

			k := kb.Build()
			defer k.Close() // don't forget to Close()!

			results, _ := k.Analyze("안녕하세요 코딩냄비입니다. 부글부글.", 1 /*=topN*/, KIWI_MATCH_ALL)
			fmt.Println(results)
			// Output:
			// [{[{0 NNG 안녕} {2 XSA 하} {4 EP 시} {3 EC 어요} {6 NNP 코딩냄비} {10 VCP 이} {11 EF ᆸ니다} {13 SF .} {15 NNP 부글부} {18 NNG 글} {19 SF .}] -69.74997}]
		}()

	})

}
