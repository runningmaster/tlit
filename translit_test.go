package tlit_test

import (
	"testing"

	"github.com/runningmaster/tlit"
)

var tranTest = map[tlit.System][]struct {
	in   string
	want string
}{
	tlit.Default: {
		{"Йух", "Yuh"},
		{"Йух ® йух™", "Yuh ® yuh™"},
		{"", ""},
	},
	tlit.UkrainianStd: {
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
	for _, sys := range []tlit.System{
		tlit.Default,
		tlit.DriversLicense,
		tlit.GOST1971,
		tlit.GOST2002B,
		tlit.GOST2006,
		tlit.Passport1997,
		tlit.Passport2010,
		tlit.Passport2013ICAO,
		tlit.Telegram,
		tlit.UkrainianStd,
		tlit.UkrainianWeb,
	} {
		for _, tt := range tranTest[sys] {
			got, err := tlit.MarshalString(tt.in, sys)
			if err != nil {
				AssertEquals(t, err.Error(), tt.want)
			}
			AssertEquals(t, got, tt.want)
		}
	}
}

func TestMarshalStringURLru(t *testing.T) {
	in, want := "Путин -® IC Хуйло™", "putin-ic-huylo"
	AssertEquals(t, tlit.MarshalStringURLru(in), want)
}

func TestMarshalStringURLua(t *testing.T) {
	in, want := "Путін -® IC Хуйло™", "putin-ic-huylo"
	AssertEquals(t, tlit.MarshalStringURLua(in), want)
}

func BenchmarkMarshalString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = tlit.MarshalString(tranTest[tlit.Default][0].in, tlit.Default)
	}
}
