package tui

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// FilePicker は単純なファイル選択コンポーネントです
type FilePicker struct {
	dir          string        // 現在のディレクトリ
	files        []fs.DirEntry // ディレクトリ内のファイルとディレクトリ
	cursor       int           // 選択中のアイテムのインデックス
	selected     string        // 選択されたファイルパス
	err          error         // エラー発生時のエラー情報
	allowDirs    bool          // ディレクトリも選択可能かのフラグ
	showHidden   bool          // 隠しファイルを表示するかのフラグ
	height       int           // 表示する行の高さ
	filterSuffix string        // ファイル拡張子によるフィルター
	title        string        // ファイルピッカーのタイトル
}

// NewFilePicker は新しいFilePickerを生成します
func NewFilePicker(initialDir string, options ...FilePickerOption) (*FilePicker, error) {
	if initialDir == "" {
		// 初期ディレクトリが指定されていない場合は現在のディレクトリを使用
		currentDir, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		initialDir = currentDir
	}

	abs, err := filepath.Abs(initialDir)
	if err != nil {
		return nil, err
	}

	fp := &FilePicker{
		dir:          abs,
		allowDirs:    false,
		showHidden:   false,
		height:       10,
		filterSuffix: "",
		title:        "ファイル選択",
	}

	// オプションの適用
	for _, opt := range options {
		opt(fp)
	}

	// 初期ディレクトリのファイル一覧を読み込む
	err = fp.loadFiles()
	if err != nil {
		return nil, err
	}

	return fp, nil
}

// FilePickerOption はFilePickerのオプションを設定する関数型です
type FilePickerOption func(*FilePicker)

// WithAllowDirs はディレクトリ選択を許可するオプションです
func WithAllowDirs(allow bool) FilePickerOption {
	return func(fp *FilePicker) {
		fp.allowDirs = allow
	}
}

// WithShowHidden は隠しファイルを表示するオプションです
func WithShowHidden(show bool) FilePickerOption {
	return func(fp *FilePicker) {
		fp.showHidden = show
	}
}

// WithHeight は表示する行数を設定するオプションです
func WithHeight(height int) FilePickerOption {
	return func(fp *FilePicker) {
		fp.height = height
	}
}

// WithFilterSuffix はファイル拡張子でフィルタリングするオプションです
func WithFilterSuffix(suffix string) FilePickerOption {
	return func(fp *FilePicker) {
		fp.filterSuffix = suffix
	}
}

// WithTitle はファイルピッカーのタイトルを設定するオプションです
func WithTitle(title string) FilePickerOption {
	return func(fp *FilePicker) {
		fp.title = title
	}
}

// Selected は選択されたファイルパスを返します
func (fp *FilePicker) Selected() string {
	return fp.selected
}

// loadFiles は現在のディレクトリ内のファイル一覧を読み込みます
func (fp *FilePicker) loadFiles() error {
	entries, err := os.ReadDir(fp.dir)
	if err != nil {
		return err
	}

	// フィルタリング
	var filteredEntries []fs.DirEntry
	for _, entry := range entries {
		name := entry.Name()

		// 隠しファイルのフィルタリング
		if !fp.showHidden && strings.HasPrefix(name, ".") {
			continue
		}

		// 拡張子によるフィルタリング（ディレクトリの場合は無視）
		if !entry.IsDir() && fp.filterSuffix != "" {
			if !strings.HasSuffix(strings.ToLower(name), strings.ToLower(fp.filterSuffix)) {
				continue
			}
		}

		filteredEntries = append(filteredEntries, entry)
	}

	// ディレクトリとファイルを分けてソート
	var dirs, files []fs.DirEntry
	for _, entry := range filteredEntries {
		if entry.IsDir() {
			dirs = append(dirs, entry)
		} else {
			files = append(files, entry)
		}
	}

	// それぞれをソート
	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].Name() < dirs[j].Name()
	})
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	// 親ディレクトリへのリンクを追加（ルートディレクトリでない場合）
	fp.files = []fs.DirEntry{}
	if fp.dir != "/" {
		// 親ディレクトリへのエントリを作成（ダミーエントリ）
		fp.files = append(fp.files, parentDirEntry{})
	}

	// ディレクトリを先に、その後にファイルを追加
	fp.files = append(fp.files, append(dirs, files...)...)
	fp.cursor = 0

	return nil
}

// parentDirEntry は親ディレクトリへのエントリを表すダミー構造体です
type parentDirEntry struct{}

func (p parentDirEntry) Name() string               { return ".." }
func (p parentDirEntry) IsDir() bool                { return true }
func (p parentDirEntry) Type() fs.FileMode          { return fs.ModeDir }
func (p parentDirEntry) Info() (fs.FileInfo, error) { return nil, nil }

// Init はBubbleteaのイニシャライザです
func (fp *FilePicker) Init() tea.Cmd {
	return nil
}

// Update はBubbleteaの更新関数です
func (fp *FilePicker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			// 選択をキャンセル
			fp.selected = ""
			return fp, tea.Quit

		case "up", "k":
			if fp.cursor > 0 {
				fp.cursor--
			}

		case "down", "j":
			if fp.cursor < len(fp.files)-1 {
				fp.cursor++
			}

		case "enter":
			if len(fp.files) == 0 {
				return fp, nil
			}

			selected := fp.files[fp.cursor]

			// 親ディレクトリへの移動
			if selected.Name() == ".." {
				parentDir := filepath.Dir(fp.dir)
				fp.dir = parentDir
				err := fp.loadFiles()
				if err != nil {
					fp.err = err
				}
				return fp, nil
			}

			// ディレクトリを選択した場合
			if selected.IsDir() {
				if fp.allowDirs {
					// ディレクトリ選択が許可されている場合、選択として確定
					fp.selected = filepath.Join(fp.dir, selected.Name())
					return fp, nil
				} else {
					// ディレクトリ内に移動
					newPath := filepath.Join(fp.dir, selected.Name())
					fp.dir = newPath
					err := fp.loadFiles()
					if err != nil {
						fp.err = err
					}
				}
				return fp, nil
			}

			// ファイルを選択した場合
			fp.selected = filepath.Join(fp.dir, selected.Name())
			return fp, nil
		}
	}

	return fp, nil
}

// View はBubbleteaの表示関数です
func (fp *FilePicker) View() string {
	if fp.err != nil {
		return ErrorStyle.Render(fmt.Sprintf("エラー: %s", fp.err))
	}

	var s strings.Builder

	// タイトルと現在のディレクトリを表示
	s.WriteString(TitleStyle.Render(fp.title) + "\n")
	s.WriteString(SubtextStyle.Render(fmt.Sprintf("現在のディレクトリ: %s", fp.dir)) + "\n\n")

	if len(fp.files) == 0 {
		s.WriteString("ファイルが見つかりません\n")
		return s.String()
	}

	// 表示範囲の計算
	start := 0
	numFiles := len(fp.files)

	if fp.height > 0 && numFiles > fp.height {
		halfHeight := fp.height / 2
		start = fp.cursor - halfHeight
		if start < 0 {
			start = 0
		} else if start+fp.height > numFiles {
			start = numFiles - fp.height
		}
	}

	end := numFiles
	if fp.height > 0 && start+fp.height < numFiles {
		end = start + fp.height
	}

	// ファイル一覧の表示（ファイル数とフィルターの情報を追加）
	s.WriteString(SubtextStyle.Render(fmt.Sprintf("ファイル数: %d (フィルター: %s)\n\n", numFiles, fp.filterSuffix)))

	// ファイル一覧の表示
	for i := start; i < end; i++ {
		file := fp.files[i]
		cursor := " "
		if i == fp.cursor {
			cursor = ">"
		}

		name := file.Name()
		if file.IsDir() {
			name += "/"
		}

		line := fmt.Sprintf("%s %s", cursor, name)

		if i == fp.cursor {
			s.WriteString(SelectedButtonStyle.Render(line))
		} else {
			if file.IsDir() {
				// ディレクトリは青色で表示
				s.WriteString(ButtonStyle.Render(line))
			} else {
				s.WriteString(TextStyle.Render(line))
			}
		}
		s.WriteString("\n")
	}

	// 操作ガイドを表示
	s.WriteString("\n")
	s.WriteString(SubtextStyle.Render("↑/↓: 移動   Enter: 選択   Esc: キャンセル"))

	return s.String()
}

// RunFilePicker はファイル選択ダイアログを実行して選択されたファイルパスを返します
func RunFilePicker(initialDir string, options ...FilePickerOption) (string, error) {
	fp, err := NewFilePicker(initialDir, options...)
	if err != nil {
		return "", err
	}

	p := tea.NewProgram(fp)
	_, err = p.Run()
	if err != nil {
		return "", err
	}

	return fp.selected, nil
}
