# はば単 - Smart Vocabulary

辞書風 UI を採用した、英単語学習アプリケーションです。

## 特徴

- 📖 **辞書風デザイン**: 見やすいライトテーマと左揃えレイアウト
- ⌨️ **片手キーボード操作**: 左手・右手どちらでも快適に操作可能
- 🔀 **シャッフル機能**: ランダムな順序で学習
- ⭐ **ブックマーク機能**: 重要な単語をマーク
- 🪟 **リサイズ対応**: ウィンドウサイズを自由に変更可能

## ダウンロード

最新のビルドは以下からダウンロードできます：

### GitHub Releases（推奨）

[Nightly Build](https://github.com/puyokura/flashcards/releases/tag/nightly)

### nightly.link（直接ダウンロード）

- [Linux](https://nightly.link/puyokura/flashcards/workflows/ci_release/main/flashcards-linux.zip)
- [Windows](https://nightly.link/puyokura/flashcards/workflows/ci_release/main/flashcards-windows.zip)
- [macOS](https://nightly.link/puyokura/flashcards/workflows/ci_release/main/flashcards-macos.zip)

## 使い方

### インストール

1. 上記のリンクから、お使いの OS に対応した ZIP ファイルをダウンロード
2. ZIP ファイルを解凍
3. 解凍したフォルダ内の実行ファイルを起動
   - **Linux/macOS**: `flashcards-linux` / `flashcards-macos`
   - **Windows**: `flashcards-windows.exe`

### キーボード操作

| 操作                  | キー                               |
| --------------------- | ---------------------------------- |
| **答えを表示 / 次へ** | `Space`, `Enter`, `Down`, `J`, `F` |
| **前へ**              | `Up`, `Left`, `K`, `H`, `A`, `D`   |
| **次へ（強制）**      | `Right`, `L`                       |
| **ブックマーク切替**  | `B`, `S`, `M`                      |

### 基本的な使い方

1. アプリを起動すると、最初の単語が表示されます
2. `Space`キーまたは`Enter`キーで答え（意味・例文）を表示
3. `Space`キーまたは矢印キーで次の単語へ
4. 重要な単語は`B`キーでブックマーク
5. 「シャッフル」にチェックを入れるとランダム順で学習できます

## 開発者向け

### ビルド方法

```bash
# 依存関係のインストール（初回のみ）
go mod download

# ビルド
make build

# ビルドして実行
make run

# クリーンアップ
make clean
```

### 必要な環境

- Go 1.21 以上
- Linux: `libgl1-mesa-dev`, `xorg-dev`
- macOS/Windows: 追加の依存関係なし

## ライセンス

このプロジェクトは個人学習用途で作成されています。
