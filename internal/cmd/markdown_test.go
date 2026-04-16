package cmd

import (
	"strings"
	"testing"
)

func TestMarkdownToHTML_EmptyString(t *testing.T) {
	got := markdownToHTML("")
	if got != "" {
		t.Errorf("markdownToHTML(\"\") = %q, want \"\"", got)
	}
}

func TestMarkdownToHTML_Heading(t *testing.T) {
	got := markdownToHTML("## 需求背景")
	if !strings.Contains(got, "<h2") {
		t.Errorf("markdownToHTML heading result should contain <h2>, got: %s", got)
	}
	if !strings.Contains(got, "需求背景") {
		t.Errorf("markdownToHTML heading result should contain '需求背景', got: %s", got)
	}
}

func TestMarkdownToHTML_Paragraph(t *testing.T) {
	got := markdownToHTML("这是一段普通文本。")
	if !strings.Contains(got, "<p>") {
		t.Errorf("markdownToHTML paragraph result should contain <p>, got: %s", got)
	}
	if !strings.Contains(got, "这是一段普通文本。") {
		t.Errorf("markdownToHTML paragraph result should contain original text, got: %s", got)
	}
}

func TestMarkdownToHTML_BoldAndItalic(t *testing.T) {
	got := markdownToHTML("支持 **粗体** 和 *斜体* 文本")
	if !strings.Contains(got, "<strong>粗体</strong>") {
		t.Errorf("expected <strong> tag, got: %s", got)
	}
	if !strings.Contains(got, "<em>斜体</em>") {
		t.Errorf("expected <em> tag, got: %s", got)
	}
}

func TestMarkdownToHTML_FencedCodeBlock(t *testing.T) {
	md := "```go\nfunc main() {\n    fmt.Println(\"Hello\")\n}\n```"
	got := markdownToHTML(md)
	if !strings.Contains(got, "<pre>") || !strings.Contains(got, "<code") {
		t.Errorf("markdownToHTML code block should contain <pre> and <code>, got: %s", got)
	}
	if !strings.Contains(got, "func main()") {
		t.Errorf("markdownToHTML code block should contain code content, got: %s", got)
	}
}

func TestMarkdownToHTML_UnorderedList(t *testing.T) {
	md := "- 项目一\n- 项目二\n- 项目三"
	got := markdownToHTML(md)
	if !strings.Contains(got, "<ul>") {
		t.Errorf("expected <ul> tag, got: %s", got)
	}
	if !strings.Contains(got, "<li>") {
		t.Errorf("expected <li> tag, got: %s", got)
	}
}

func TestMarkdownToHTML_OrderedList(t *testing.T) {
	md := "1. 第一步\n2. 第二步\n3. 第三步"
	got := markdownToHTML(md)
	if !strings.Contains(got, "<ol>") {
		t.Errorf("expected <ol> tag, got: %s", got)
	}
}

func TestMarkdownToHTML_Table(t *testing.T) {
	md := "| 名称 | 描述 |\n|------|------|\n| A | B |"
	got := markdownToHTML(md)
	if !strings.Contains(got, "<table>") {
		t.Errorf("expected <table> tag, got: %s", got)
	}
}

func TestMarkdownToHTML_Blockquote(t *testing.T) {
	md := "> 这是一段引用"
	got := markdownToHTML(md)
	if !strings.Contains(got, "<blockquote>") {
		t.Errorf("expected <blockquote> tag, got: %s", got)
	}
}

func TestMarkdownToHTML_Link(t *testing.T) {
	md := "参考 [TAPD](https://www.tapd.cn)"
	got := markdownToHTML(md)
	if !strings.Contains(got, "<a href=\"https://www.tapd.cn\"") {
		t.Errorf("expected link, got: %s", got)
	}
}

func TestMarkdownToHTML_ComplexDocument(t *testing.T) {
	md := `## 需求背景

这是一个测试 **Markdown** 的需求。

### 功能要点

1. 支持**粗体**和*斜体*
2. 支持` + "`行内代码`" + `

> 引用文字

- 列表项 A
- 列表项 B`

	got := markdownToHTML(md)
	if !strings.Contains(got, "<h2") {
		t.Error("expected h2 tag")
	}
	if !strings.Contains(got, "<h3") {
		t.Error("expected h3 tag")
	}
	if !strings.Contains(got, "<ol>") {
		t.Error("expected ol tag")
	}
	if !strings.Contains(got, "<ul>") {
		t.Error("expected ul tag")
	}
	if !strings.Contains(got, "<blockquote>") {
		t.Error("expected blockquote tag")
	}
	if !strings.Contains(got, "<strong>") {
		t.Error("expected strong tag")
	}
}

func TestContainsBlockHTML_WithBlockTags(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"<p>text</p>", true},
		{"<h1>heading</h1>", true},
		{"<ul><li>item</li></ul>", true},
		{"<ol><li>item</li></ol>", true},
		{"<pre><code>code</code></pre>", true},
		{"<blockquote>quote</blockquote>", true},
		{"<table><tr><td>cell</td></tr></table>", true},
		{"<div>content</div>", true},
		{"plain text without tags", false},
		{"<span>inline only</span>", false},
		{"<em>emphasis</em>", false},
		{"", false},
	}

	for _, tt := range tests {
		got := containsBlockHTML(tt.input)
		if got != tt.want {
			t.Errorf("containsBlockHTML(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
