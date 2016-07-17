package base91

import (
	"testing"
)

func testEq(a, b []byte) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

var str = `https://dl.google.com/tag/s/appguid%3D%7B8A69D345-D564-463C-AFF1-A69D9E530F96%7D%26iid%3D%7B1286A4E6-C3EC-E4C5-608C-0318FEE0C519%7D%26lang%3Dzh-CN%26browser%3D4%26usagestats%3D0%26appname%3DGoogle%2520Chrome%26needsadmin%3Dprefers%26ap%3Dx64-stable%26brand%3DCHWL%26installdataindex%3Ddefaultbrowser/update2/installers/ChromeStandaloneSetup64.exe`

func Test(t *testing.T) {
	raw := []byte(str)
	en := StdEncoding.EncodeToString(raw)
	de, _ := StdEncoding.DecodeString(en)

	if !testEq(raw, de) {
		t.Errorf("Decoded string not match raw string while \n\nraw:%v\n\ndec:%v", raw, de)
	}
}
