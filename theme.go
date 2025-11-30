package main

import (
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type myTheme struct{}

var _ fyne.Theme = (*myTheme)(nil)

func (m *myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	// Force Light Theme
	return theme.DefaultTheme().Color(name, theme.VariantLight)
}

func (m *myTheme) Font(style fyne.TextStyle) fyne.Resource {
	// 埋め込みではなく、ローカルファイルを読み込む簡易実装
	// 本番では bundle.go などを使うのが一般的だが、今回はシンプルに
	fontData, err := os.ReadFile("DroidSansFallbackFull.ttf")
	if err != nil {
		// フォント読み込み失敗時はデフォルトに戻す（文字化けするがクラッシュは防ぐ）
		return theme.DefaultTheme().Font(style)
	}
	return fyne.NewStaticResource("DroidSansFallbackFull.ttf", fontData)
}

func (m *myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m *myTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
