package main

import (
	"encoding/csv"
	"fmt"
	"image/color"
	"math/rand"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// --- データ構造 ---

type WordEntry struct {
	ID         string
	Word       string
	Pos        string
	Meaning    string
	Example    string
	Bookmarked bool
}

// --- グローバル状態 ---

var (
	allEntries     []WordEntry
	displayEntries []WordEntry
	currentIdx     int
	isShuffle      bool
)

// --- UI コンポーネント (グローバル変数として定義) ---

var (
	// ヘッダー
	progressLabel *widget.Label
	progressBar   *widget.ProgressBar
	checkShuffle  *widget.Check

	// メインカード
	idLabel    *widget.Label
	wordLabel  *widget.RichText
	starButton *widget.Button

	// 解説エリア
	detailContainer *fyne.Container
	posText         *widget.RichText
	meaningText     *widget.RichText
	exampleText     *widget.RichText

	// コントロールボタン (ここをグローバルにしないと updateUI などから参照できない)
	btnShow *widget.Button
	btnPrev *widget.Button
	btnNext *widget.Button
)

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	myApp := app.New()
	myApp.Settings().SetTheme(&myTheme{})
	myWindow := myApp.NewWindow("はば単 - Smart Vocabulary")

	// CSV読み込み
	var err error
	allEntries, err = loadCSV("habatan.csv")
	if err != nil {
		myWindow.SetContent(widget.NewLabel(fmt.Sprintf("エラー: %v", err)))
		myWindow.ShowAndRun()
		return
	}

	displayEntries = make([]WordEntry, len(allEntries))
	copy(displayEntries, allEntries)

	if len(displayEntries) == 0 {
		myWindow.SetContent(widget.NewLabel("データがありません。"))
		myWindow.ShowAndRun()
		return
	}

	// --- UI構築 ---

	// 1. ヘッダー
	progressLabel = widget.NewLabel("0 / 0")
	progressBar = widget.NewProgressBar()
	progressBar.TextFormatter = func() string { return "" }

	checkShuffle = widget.NewCheck("シャッフル", func(on bool) {
		isShuffle = on
		toggleShuffle(on)
	})

	headerBar := container.NewBorder(
		nil, nil,
		checkShuffle,
		progressLabel,
		progressBar,
	)

	// 2. カードエリア
	cardBackground := canvas.NewRectangle(color.NRGBA{R: 40, G: 40, B: 40, A: 255})
	cardBackground.CornerRadius = 16

	idLabel = widget.NewLabel("No.1")
	idLabel.Alignment = fyne.TextAlignCenter

	// アイコン修正: OutlineStarIconがない場合があるため、確実に存在するアイコンを使用
	// ここでは「保存」アイコンなどで代用
	starButton = widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), func() {
		toggleBookmark()
	})
	starButton.Importance = widget.LowImportance

	wordLabel = widget.NewRichTextFromMarkdown("")
	if len(wordLabel.Segments) > 0 {
		if seg, ok := wordLabel.Segments[0].(*widget.TextSegment); ok {
			seg.Style.Alignment = fyne.TextAlignCenter
		}
	}

	posText = widget.NewRichTextFromMarkdown("")
	meaningText = widget.NewRichTextFromMarkdown("")
	exampleText = widget.NewRichTextFromMarkdown("")

	detailContent := container.NewVBox(
		widget.NewSeparator(),
		posText,
		meaningText,
		layout.NewSpacer(),
		widget.NewLabelWithStyle("【例文】", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true}),
		exampleText,
	)
	detailContainer = container.NewPadded(detailContent)
	detailContainer.Hide()

	cardTop := container.NewBorder(nil, nil, nil, starButton, idLabel)

	wordArea := container.NewVBox(
		layout.NewSpacer(),
		wordLabel,
		layout.NewSpacer(),
	)

	detailScroll := container.NewVScroll(detailContainer)
	detailScroll.SetMinSize(fyne.NewSize(0, 150))

	cardInner := container.NewBorder(
		cardTop,
		detailScroll,
		nil, nil,
		wordArea,
	)

	cardStack := container.NewStack(
		cardBackground,
		container.NewPadded(cardInner),
	)
	mainCardArea := container.NewPadded(cardStack)

	// 3. フッターボタン (グローバル変数へ代入)
	btnShow = widget.NewButton("答えを見る", func() {
		showAnswer()
	})
	btnShow.Importance = widget.HighImportance

	btnPrev = widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		moveIndex(-1)
	})
	btnNext = widget.NewButtonWithIcon("次へ", theme.NavigateNextIcon(), func() {
		moveIndex(1)
	})

	footerButtons := container.NewBorder(
		nil, nil,
		btnPrev, btnNext,
		btnShow,
	)

	finalLayout := container.NewBorder(
		container.NewPadded(headerBar),
		container.NewPadded(footerButtons),
		nil, nil,
		mainCardArea,
	)

	updateUI()

	myWindow.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		switch k.Name {
		case fyne.KeySpace, fyne.KeyEnter, fyne.KeyDown:
			if detailContainer.Hidden {
				showAnswer()
			} else {
				moveIndex(1)
			}
		case fyne.KeyRight:
			moveIndex(1)
		case fyne.KeyLeft:
			moveIndex(-1)
		}
	})

	myWindow.SetContent(finalLayout)
	myWindow.Resize(fyne.NewSize(450, 700))
	myWindow.SetFixedSize(true) // サイズ変更を禁止
	myWindow.ShowAndRun()
}

func loadCSV(filename string) ([]WordEntry, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.LazyQuotes = true

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var result []WordEntry
	for i, row := range rows {
		if i == 0 || len(row) < 5 {
			continue
		}
		entry := WordEntry{
			ID:         row[0],
			Word:       row[1],
			Pos:        row[2],
			Meaning:    row[3],
			Example:    row[4],
			Bookmarked: false,
		}
		result = append(result, entry)
	}
	return result, nil
}

func moveIndex(delta int) {
	currentIdx += delta
	if currentIdx < 0 {
		currentIdx = 0
	}
	if currentIdx >= len(displayEntries) {
		currentIdx = len(displayEntries) - 1
	}
	detailContainer.Hide()
	btnShow.Enable()
	btnShow.SetText("答えを見る")
	updateUI()
}

func toggleShuffle(on bool) {
	if on {
		shuffled := make([]WordEntry, len(allEntries))
		copy(shuffled, allEntries)
		rand.Shuffle(len(shuffled), func(i, j int) {
			shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
		})
		displayEntries = shuffled
	} else {
		displayEntries = make([]WordEntry, len(allEntries))
		copy(displayEntries, allEntries)
	}
	currentIdx = 0
	detailContainer.Hide()
	btnShow.Enable()
	updateUI()
}

func toggleBookmark() {
	entry := &displayEntries[currentIdx]
	entry.Bookmarked = !entry.Bookmarked
	updateBookmarkIcon(entry.Bookmarked)
}

func updateBookmarkIcon(isBookmarked bool) {
	if isBookmarked {
		// ConfirmIcon(チェックマーク)で代用し、HighImportanceで色をつける
		starButton.SetIcon(theme.ConfirmIcon())
		starButton.Importance = widget.HighImportance
	} else {
		// 保存アイコンで未保存状態を表す
		starButton.SetIcon(theme.DocumentSaveIcon())
		starButton.Importance = widget.LowImportance
	}
	starButton.Refresh()
}

func showAnswer() {
	detailContainer.Show()
	btnShow.Disable()
	btnShow.SetText("Next ->")
}

func updateUI() {
	if len(displayEntries) == 0 {
		return
	}
	entry := displayEntries[currentIdx]

	progressLabel.SetText(fmt.Sprintf("%d / %d", currentIdx+1, len(displayEntries)))
	progressBar.SetValue(float64(currentIdx+1) / float64(len(displayEntries)))

	idLabel.SetText(fmt.Sprintf("No. %s", entry.ID))

	wordLabel.ParseMarkdown("# " + entry.Word)
	// Markdownパースでスタイルがリセットされる可能性があるため、中央揃えを再適用
	if len(wordLabel.Segments) > 0 {
		if seg, ok := wordLabel.Segments[0].(*widget.TextSegment); ok {
			seg.Style.Alignment = fyne.TextAlignCenter
		}
	}
	wordLabel.Refresh()

	posStr := strings.ReplaceAll(entry.Pos, "\n", " ")
	posText.ParseMarkdown(fmt.Sprintf("**%s**", posStr))

	meanStr := strings.ReplaceAll(entry.Meaning, "\n", "\n")
	meaningText.ParseMarkdown(meanStr)

	exStr := strings.ReplaceAll(entry.Example, "\n", "\n")
	exampleText.ParseMarkdown(fmt.Sprintf("> %s", strings.ReplaceAll(exStr, "\n", "\n> ")))

	updateBookmarkIcon(entry.Bookmarked)
}
