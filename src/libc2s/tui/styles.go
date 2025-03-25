package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// 色の定義
var (
	primaryColor   = lipgloss.Color("#1E88E5") // 青色
	secondaryColor = lipgloss.Color("#FFC107") // 黄色
	accentColor    = lipgloss.Color("#E53935") // 赤色
	textColor      = lipgloss.Color("#212121") // 濃いグレー
	subtextColor   = lipgloss.Color("#757575") // 薄いグレー
	bgColor        = lipgloss.Color("#FFFFFF") // 白色
)

// スタイルの定義
var (
	// ヘッダースタイル
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(bgColor).
			Background(primaryColor).
			Padding(0, 1).
			Width(80)

	// タイトルスタイル
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			MarginBottom(1)

	// 通常テキストスタイル
	TextStyle = lipgloss.NewStyle().
			Foreground(textColor)

	// サブテキストスタイル
	SubtextStyle = lipgloss.NewStyle().
			Foreground(subtextColor).
			Italic(true)

	// 入力フィールドスタイル
	InputFieldStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(0, 1)

	// ボタンスタイル
	ButtonStyle = lipgloss.NewStyle().
			Foreground(bgColor).
			Background(primaryColor).
			Padding(0, 3).
			MarginRight(1)

	// 選択ボタンスタイル
	SelectedButtonStyle = lipgloss.NewStyle().
				Foreground(bgColor).
				Background(accentColor).
				Padding(0, 3).
				MarginRight(1)

	// エラーメッセージスタイル
	ErrorStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true)

	// 成功メッセージスタイル
	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4CAF50")). // 緑色
			Bold(true)

	// セクションスタイル
	SectionStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(primaryColor).
			Padding(1).
			MarginTop(1).
			MarginBottom(1)

	// フッタースタイル
	FooterStyle = lipgloss.NewStyle().
			Foreground(bgColor).
			Background(primaryColor).
			Padding(0, 1).
			Width(80)
)
