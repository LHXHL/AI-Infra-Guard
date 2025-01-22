package parser

import "testing"

func TestParseTokens(t *testing.T) {
	s := `body="href=\"http://www.thinkphp.cn\">thinkphp</a>" || body="thinkphp_show_page_trace" || icon="f49c4a4bde1eec6c0b80c2277c76e3dbs"`
	tokens, err := ParseTokens(s)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tokens)
}

func TestParseTokens2(t *testing.T) {
	s := "body~=\"(<center><strong>EZCMS ([\\d\\.]+) )\""
	tokens, err := ParseTokens(s)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tokens)
}

func TestParseAdvisorTokens2(t *testing.T) {
	s := "version >= \"1.0.0\" || version < \"2.0.0\" || version == \"3.0.0\""
	tokens, err := ParseAdvisorTokens(s)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tokens)
}
