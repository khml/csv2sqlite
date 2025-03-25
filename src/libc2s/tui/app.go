package tui

import (
	"csv2sqlite/libc2s"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Csv2SqliteApp はメインのTUIアプリケーションです
type Csv2SqliteApp struct {
	state      string
	filePicker *FilePicker
	form       *Form
	progress   *ProgressModel
	csvPath    string
	dbPath     string
	tableName  string
	updateCh   chan ProgressUpdate
	errCh      chan error
}

// 状態を表す定数
const (
	StateSelectCSV = "selectCSV"
	StateSelectDB  = "selectDB"
	StateForm      = "form"
	StateConfirm   = "confirm"
	StateProgress  = "progress"
	StateComplete  = "complete"
)

// NewCsv2SqliteApp は新しいアプリケーションを作成します
func NewCsv2SqliteApp() *Csv2SqliteApp {
	return &Csv2SqliteApp{
		state:    StateSelectCSV,
		updateCh: make(chan ProgressUpdate),
		errCh:    make(chan error),
	}
}

// Init はBubbleteaのイニシャライザです
func (app *Csv2SqliteApp) Init() tea.Cmd {
	// CSVファイル選択から開始
	picker, err := NewFilePicker(".",
		WithFilterSuffix(".csv"),
		WithTitle("CSVファイルを選択してください"),
		WithShowHidden(true), // 隠しファイルも表示
	)
	if err != nil {
		return func() tea.Msg {
			return ProgressError{Error: err}
		}
	}
	app.filePicker = picker
	return nil
}

// Update はBubbleteaの更新関数です
func (app *Csv2SqliteApp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return app, tea.Quit
		}
	}

	switch app.state {
	case StateSelectCSV:
		var cmd tea.Cmd
		newModel, cmd := app.filePicker.Update(msg)
		updatedPicker, ok := newModel.(*FilePicker)
		if ok {
			app.filePicker = updatedPicker

			// ファイルが選択された場合
			if app.filePicker.Selected() != "" {
				app.csvPath = app.filePicker.Selected()

				// 次はDBファイル選択へ
				picker, err := NewFilePicker(".",
					WithTitle("SQLiteファイルを選択または新規ファイル名を入力してください"),
					WithFilterSuffix(".sqlite"),
					WithFilterSuffix(".sqlite3"),
					WithFilterSuffix(".db"),
				)
				if err != nil {
					return app, func() tea.Msg {
						return ProgressError{Error: err}
					}
				}
				app.filePicker = picker
				app.state = StateSelectDB
			}
		}

		return app, cmd

	case StateSelectDB:
		var cmd tea.Cmd
		newModel, cmd := app.filePicker.Update(msg)
		updatedPicker, ok := newModel.(*FilePicker)
		if ok {
			app.filePicker = updatedPicker

			// ファイルが選択された場合
			selectedPath := app.filePicker.Selected()

			if selectedPath != "" {
				// 選択されたパスをDBパスとして設定
				app.dbPath = selectedPath

				// 次はフォーム入力へ
				form := NewForm("テーブル情報を入力してください")
				form.AddField("テーブル名", "例：users", true, nil)
				form.SetSubmitLabel("次へ")
				app.form = form
				app.state = StateForm

				// ここでreturnすることで確実に次の状態に移行する
				return app, cmd
			}
		}

		return app, cmd

	case StateForm:
		var cmd tea.Cmd
		newModel, cmd := app.form.Update(msg)
		updatedForm, ok := newModel.(*Form)
		if ok {
			app.form = updatedForm

			// フォームが送信された場合
			if app.form.Submitted {
				values := app.form.GetValues()
				app.tableName = values["テーブル名"]

				// 確認画面へ
				app.state = StateConfirm
			} else if app.form.Canceled {
				// キャンセルされた場合、DBファイル選択に戻る
				picker, err := NewFilePicker(".",
					WithTitle("SQLiteファイルを選択または新規ファイル名を入力してください"),
					WithFilterSuffix(".db"))
				if err != nil {
					return app, func() tea.Msg {
						return ProgressError{Error: err}
					}
				}
				app.filePicker = picker
				app.state = StateSelectDB
			}
		}

		return app, cmd

	case StateConfirm:
		// Enterキーで処理開始、Escキーで戻る
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			if keyMsg.String() == "enter" {
				// 進捗表示へ
				progress := NewProgressModel("CSV→SQLite変換", 100)
				app.progress = progress
				app.state = StateProgress

				// 実際の変換処理を開始
				go app.runConversion()

				// 進捗更新コマンドを返す
				return app, RunProgress(app.progress, app.updateCh, app.errCh)
			} else if keyMsg.String() == "esc" {
				// フォーム入力に戻る
				form := NewForm("テーブル情報を入力してください")
				form.AddField("テーブル名", app.tableName, true, nil)
				form.SetSubmitLabel("次へ")
				app.form = form
				app.state = StateForm
			}
		}

		return app, nil

	case StateProgress:
		var cmd tea.Cmd
		newModel, cmd := app.progress.Update(msg)
		updatedProgress, ok := newModel.(*ProgressModel)
		if ok {
			app.progress = updatedProgress

			// 処理が完了した場合
			if app.progress.Status == StatusCompleted || app.progress.Status == StatusFailed {
				app.state = StateComplete
			}
		}

		return app, cmd

	case StateComplete:
		// 完了画面でqキーが押されたら終了
		if msg, ok := msg.(tea.KeyMsg); ok && (msg.String() == "q" || msg.String() == "enter") {
			return app, tea.Quit
		}

		return app, nil
	}

	return app, nil
}

// View はBubbleteaの表示関数です
func (app *Csv2SqliteApp) View() string {
	switch app.state {
	case StateSelectCSV, StateSelectDB:
		return app.filePicker.View()

	case StateForm:
		return app.form.View()

	case StateConfirm:
		return confirmView(app.csvPath, app.dbPath, app.tableName)

	case StateProgress, StateComplete:
		return app.progress.View()
	}

	return "Unknown state"
}

// runConversion は実際のCSV→SQLite変換処理を実行します
func (app *Csv2SqliteApp) runConversion() {
	// 処理開始ログ
	app.updateCh <- ProgressUpdate{Current: 0, Log: "処理を開始しています..."}
	time.Sleep(500 * time.Millisecond)

	// CSVデータの読み込み
	app.updateCh <- ProgressUpdate{Current: 10, Log: fmt.Sprintf("CSVファイル '%s' を読み込んでいます...", app.csvPath)}
	csvData, err := libc2s.ReadCsv(app.csvPath)
	if err != nil {
		app.errCh <- fmt.Errorf("CSVファイルの読み込みに失敗しました: %v", err)
		close(app.updateCh)
		close(app.errCh)
		return
	}

	// テーブル作成とレコード数カウント（レコード数を取得するためにCSVを一度読み込む）
	recordCount := 0
	for {
		_, err := csvData.Reader.Read()
		if err != nil {
			break
		}
		recordCount++
	}

	// CSVファイルを再度開く（上記のカウントでファイルポインタが最後まで進んでいるため）
	app.updateCh <- ProgressUpdate{Current: 20, Log: "データベース接続を準備しています..."}
	csvData, err = libc2s.ReadCsv(app.csvPath)
	if err != nil {
		app.errCh <- fmt.Errorf("CSVファイルの再読み込みに失敗しました: %v", err)
		close(app.updateCh)
		close(app.errCh)
		return
	}

	// データベース接続
	app.updateCh <- ProgressUpdate{Current: 30, Log: fmt.Sprintf("SQLiteデータベース '%s' に接続しています...", app.dbPath)}
	db, err := libc2s.ConnectRepository(app.dbPath)
	if err != nil {
		app.errCh <- fmt.Errorf("データベース接続に失敗しました: %v", err)
		close(app.updateCh)
		close(app.errCh)
		return
	}
	defer db.Close()

	// テーブル作成
	app.updateCh <- ProgressUpdate{Current: 40, Log: fmt.Sprintf("テーブル '%s' を作成しています...", app.tableName)}
	err = db.CreateTbl(app.tableName, csvData)
	if err != nil {
		app.errCh <- fmt.Errorf("テーブル作成に失敗しました: %v", err)
		close(app.updateCh)
		close(app.errCh)
		return
	}

	// ProgressModelの合計レコード数を設定
	if app.progress != nil {
		app.progress.SetTotal(recordCount)
		app.updateCh <- ProgressUpdate{Current: 50, Log: fmt.Sprintf("合計 %d レコードの処理を開始します...", recordCount)}
	}

	// レコード挿入
	numRecords, err := db.InsertRecords(app.tableName, csvData)
	if err != nil {
		app.errCh <- fmt.Errorf("レコード挿入に失敗しました: %v", err)
		close(app.updateCh)
		close(app.errCh)
		return
	}

	// 完了メッセージ
	app.updateCh <- ProgressUpdate{
		Current: recordCount,
		Log:     fmt.Sprintf("変換完了！%d レコードをテーブル '%s' に挿入しました", numRecords, app.tableName),
	}

	// 結果サマリーを設定
	if app.progress != nil {
		summary := fmt.Sprintf("変換結果サマリー:\n"+
			"CSVファイル: %s\n"+
			"SQLiteデータベース: %s\n"+
			"テーブル名: %s\n"+
			"挿入レコード数: %d", app.csvPath, app.dbPath, app.tableName, numRecords)
		app.progress.SetResult(summary)
	}

	close(app.updateCh)
	close(app.errCh)
}

// confirmView は確認画面のビューを生成します
func confirmView(csvPath, dbPath, tableName string) string {
	title := TitleStyle.Render("設定内容の確認")
	section := SectionStyle.Render(fmt.Sprintf(`CSVファイル: %s
SQLiteファイル: %s
テーブル名: %s`, csvPath, dbPath, tableName))

	guide := SubtextStyle.Render("Enter: 処理開始   Esc: 戻る")

	return fmt.Sprintf("\n%s\n\n%s\n\n%s\n", title, section, guide)
}

// RunApp はアプリケーションを実行します
func RunApp() error {
	app := NewCsv2SqliteApp()
	p := tea.NewProgram(app)
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("アプリケーションの実行中にエラーが発生しました: %v", err)
	}
	return nil
}
