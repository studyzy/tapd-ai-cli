// Package cmd 中的 markdown.go 提供 Markdown 与 HTML 的双向转换辅助函数
package cmd

import (
	"strings"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// markdownToHTML 将 Markdown 文本转换为 HTML。
// 仅当转换结果包含块级 HTML 元素时才返回 HTML，否则返回原始文本。
// 这样可以避免对已经是纯文本的内容进行不必要的转换。
func markdownToHTML(md string) string {
	if md == "" {
		return md
	}

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs |
		parser.FencedCode | parser.Tables
	p := parser.NewWithExtensions(extensions)

	opts := html.RendererOptions{Flags: html.CommonFlags}
	renderer := html.NewRenderer(opts)

	htmlBytes := markdown.ToHTML([]byte(md), p, renderer)
	result := strings.TrimSpace(string(htmlBytes))

	// 安全检查：仅当输出包含块级 HTML 元素时才使用 HTML
	if containsBlockHTML(result) {
		return result
	}
	return md
}

// containsBlockHTML 检查 HTML 字符串是否包含块级元素标签
func containsBlockHTML(s string) bool {
	blockTags := []string{"<p>", "<p ", "<h1", "<h2", "<h3", "<h4", "<h5", "<h6",
		"<ul", "<ol", "<li", "<pre", "<blockquote", "<table", "<div"}
	lower := strings.ToLower(s)
	for _, tag := range blockTags {
		if strings.Contains(lower, tag) {
			return true
		}
	}
	return false
}

// htmlToMarkdown 将 HTML 文本转换为 Markdown。
// 空字符串直接返回；转换失败时返回原始 HTML。
func htmlToMarkdown(h string) string {
	if h == "" {
		return h
	}
	md, err := htmltomarkdown.ConvertString(h)
	if err != nil {
		return h
	}
	return md
}
