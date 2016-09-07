package tlit

import (
	"testing"
)

var tranTest = map[System][]struct {
	in   string
	want string
}{
	Default: {
		{"Йух", "Yuh"},
		{"Йух ® йух™", "Yuh ® yuh™"},
		{"", ""},
	},
	UkrainianStd: {
		{"Згорани", "Zghorany"},
		{"Розгін", "Rozghin"},
		{"Сява Сянтович", "Siava Siantovych"},
		{"Сява Сiaнтович", "Siava Siantovych"},
		{"Зайцев", "Zaitsev"},
		{"Заіцев", "Zaitsev"},
		{"Заїцев", "Zaitsev"},
		{"Зайтсев", "Zaitsev"},
		{"Заітсев", "Zaitsev"},
		{"Заїтсев", "Zaitsev"},
		{"Рудко", "Rudko"},
		{"Рудько", "Rudko"},
		{"Булкін", "Bulkin"},
		{"Булькін", "Bulkin"},
		{"Вишневскі", "Vyshnevski"},
		{"Вишневски", "Vyshnevsky"},
		{"Вишневський", "Vyshnevskyi"},
		{"Вишневські", "Vyshnevski"},
		{"Україна", "Ukraina"},
	},
}

func AssertEquals(t *testing.T, got, want string) {
	if got != want {
		t.Errorf("Got %q but want %q", got, want)
	}
}

func TestMarshalString(t *testing.T) {
	for sys := range tableTransliteration {
		for _, tt := range tranTest[sys] {
			got, err := MarshalString(tt.in, sys)
			if err != nil {
				AssertEquals(t, err.Error(), tt.want)
			}
			AssertEquals(t, got, tt.want)
		}
	}
}

func TestMarshalStringURLru(t *testing.T) {
	in, want := "Путин -® IC Хуйло™", "putin-ic-huylo"
	AssertEquals(t, MarshalStringURLru(in), want)
}

func TestMarshalStringURLua(t *testing.T) {
	in, want := "Путін -® IC Хуйло™", "putin-ic-huylo"
	AssertEquals(t, MarshalStringURLua(in), want)
}

func BenchmarkMarshalString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = MarshalString(tranTest[Default][0].in, Default)
	}
}
