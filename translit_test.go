package tlit_test

import (
	"bytes"
	"testing"

	"github.com/runningmaster/tlit"
)

var allSystems = []tlit.System{
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
}

var tranTest = map[tlit.System][]struct {
	in   string
	want string
}{
	tlit.Default: {
		{"Йух", "Yuh"},
		{"Йух ® йух™", "Yuh ® yuh™"},
		{"", ""},
	},
	tlit.DriversLicense: {
		{"Ёж", "Yozh"},
		{"ЁЖ", "YoZh"},
		{"Европа", "Yevropa"},
		{"ЕВРОПА", "YeVROPA"},
		{"Артём", "Artyem"},
		{"семьи", "sem'yi"},
		// е and ё become ye/yo at word boundaries, not just start of string
		{"Привет Европа", "Privet Yevropa"},
	},
	tlit.GOST1971: {
		{"Ёж", "Jozh"},
		{"Йод", "Jjod"},
		{"Москва", "Moskva"},
	},
	tlit.GOST2002B: {
		{"Цирк", "Cirk"},
		{"Цена", "Cena"},
		{"Цукор", "Czukor"},
	},
	tlit.GOST2006: {
		{"Юг", "Iug"},
		{"Москва", "Moskva"},
	},
	tlit.Passport1997: {
		{"Ольега", "Olyega"},
		{"Москва", "Moskva"},
	},
	tlit.Passport2010: {
		{"Юг", "Iug"},
		{"Москва", "Moskva"},
	},
	tlit.Passport2013ICAO: {
		{"Юг", "Iug"},
		{"Москва", "Moskva"},
	},
	tlit.Telegram: {
		{"Живот", "Jivot"},
		{"Юг", "Iug"},
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
		{"Єгипет", "Yehypet"},
		{"Їжак", "Yizhak"},
		// є/ї/й become ye/yi/y at word boundaries, not just start of string
		{"Юрій Євгенов", "Iurii Yevhenov"},
		{"Їжак Йорк", "Yizhak York"},
	},
	tlit.UkrainianWeb: {
		{"Привіт", "Privit"},
		{"Гора", "Gora"},
		// Ґ maps to G (same as Г in UkrainianWeb, unlike UkrainianStd where Г→H)
		{"Ґанок", "Ganok"},
		// Є maps to E (unlike UkrainianStd where Є→Ie)
		{"Єнот", "Enot"},
		// Ї maps to I (unlike UkrainianStd where initial Ї→Yi)
		{"Їжак", "Izhak"},
		// apostrophe is consumed silently
		{"м'яч", "miach"},
	},
}

func assertEquals(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMarshalString(t *testing.T) {
	for _, sys := range allSystems {
		for _, tt := range tranTest[sys] {
			t.Run(sys.String()+"/"+tt.in, func(t *testing.T) {
				assertEquals(t, tlit.MarshalString(tt.in, sys), tt.want)
			})
		}
	}
}

func TestMarshal(t *testing.T) {
	for _, sys := range allSystems {
		for _, tt := range tranTest[sys] {
			t.Run(sys.String()+"/"+tt.in, func(t *testing.T) {
				assertEquals(t, string(tlit.Marshal([]byte(tt.in), sys)), tt.want)
			})
		}
	}
}

func TestEncoder(t *testing.T) {
	for _, sys := range allSystems {
		for _, tt := range tranTest[sys] {
			t.Run(sys.String()+"/"+tt.in, func(t *testing.T) {
				var buf bytes.Buffer
				enc := tlit.NewEncoder(&buf, sys)
				if err := enc.EncodeString(tt.in); err != nil {
					t.Fatal(err)
				}
				assertEquals(t, buf.String(), tt.want)
			})
		}
	}
}

func TestEncoderClose(t *testing.T) {
	var buf bytes.Buffer
	enc := tlit.NewEncoder(&buf, tlit.Default)
	if err := enc.Encode([]byte("Йух")); err != nil {
		t.Fatal(err)
	}
	if err := enc.Close(); err != nil {
		t.Fatal(err)
	}
	assertEquals(t, buf.String(), "Yuh")
}

func TestMarshalStringURLru(t *testing.T) {
	in, want := "Путин -® IC Хуйло™", "putin-ic-huylo"
	assertEquals(t, tlit.MarshalStringURLru(in), want)
}

func TestMarshalStringURLua(t *testing.T) {
	in, want := "Путін -® IC Хуйло™", "putin-ic-huylo"
	assertEquals(t, tlit.MarshalStringURLua(in), want)
}

func TestMarshalStringURLuaUkrainian(t *testing.T) {
	in, want := "Привіт Світ", "privit-svit"
	assertEquals(t, tlit.MarshalStringURLua(in), want)
}

const benchInput = "Съешь же ещё этих мягких французских булок, да выпей чаю"

func BenchmarkMarshalString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tlit.MarshalString(tranTest[tlit.Default][0].in, tlit.Default)
	}
}

func BenchmarkMarshalStringLong(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tlit.MarshalString(benchInput, tlit.Default)
	}
}

func BenchmarkMarshalStringURL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tlit.MarshalStringURLru(benchInput)
	}
}
