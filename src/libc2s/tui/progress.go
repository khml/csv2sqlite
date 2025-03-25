package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ProgressModel は進捗表示コンポーネントを表す構造体です
type ProgressModel struct {
	Title      string    // 処理のタイトル
	Total      int       // 処理する総アイテム数
	Current    int       // 現在処理済みのアイテム数
	Logs       []string  // ログメッセージのリスト
	MaxLogs    int       // 保持する最大ログ数
	Status     string    // 現在のステータス（実行中、完了、失敗など）
	StartTime  time.Time // 処理開始時間
	FinishTime time.Time // 処理終了時間
	Error      error     // エラー情報
	Width      int       // 進捗バーの幅
	Running    bool      // 処理中フラグ
	Quitting   bool      // 終了中フラグ
	Result     string    // 完了時の結果サマリー
}

// 進捗状態を表す定数
const (
	StatusRunning   = "実行中"
	StatusPaused    = "一時停止"
	StatusCompleted = "完了"
	StatusFailed    = "失敗"
)

// NewProgressModel は新しい進捗表示モデルを作成します
func NewProgressModel(title string, total int) *ProgressModel {
	return &ProgressModel{
		Title:     title,
		Total:     total,
		Current:   0,
		Logs:      []string{},
		MaxLogs:   10,
		Status:    StatusRunning,
		StartTime: time.Now(),
		Width:     50,
		Running:   true,
	}
}

// SetTotal は処理する総アイテム数を設定します
func (p *ProgressModel) SetTotal(total int) {
	p.Total = total
}

// Increment は進捗を1つ進めます
func (p *ProgressModel) Increment() {
	p.Current++
	if p.Current >= p.Total && p.Total > 0 {
		p.Complete()
	}
}

// SetCurrent は現在の進捗を設定します
func (p *ProgressModel) SetCurrent(current int) {
	p.Current = current
	if p.Current >= p.Total && p.Total > 0 {
		p.Complete()
	}
}

// AddLog はログメッセージを追加します
func (p *ProgressModel) AddLog(message string) {
	p.Logs = append(p.Logs, message)
	if len(p.Logs) > p.MaxLogs {
		// 古いログを削除
		p.Logs = p.Logs[len(p.Logs)-p.MaxLogs:]
	}
}

// Complete は処理を完了状態にします
func (p *ProgressModel) Complete() {
	p.Status = StatusCompleted
	p.FinishTime = time.Now()
	p.Running = false
}

// Failed は処理を失敗状態にします
func (p *ProgressModel) Failed(err error) {
	p.Status = StatusFailed
	p.Error = err
	p.FinishTime = time.Now()
	p.Running = false
}

// SetResult は完了時の結果サマリーを設定します
func (p *ProgressModel) SetResult(result string) {
	p.Result = result
}

// ProgressUpdate は進捗更新メッセージです
type ProgressUpdate struct {
	Current int
	Log     string
}

// ProgressComplete は処理完了メッセージです
type ProgressComplete struct {
	Result string
}

// ProgressError はエラーメッセージです
type ProgressError struct {
	Error error
}

// Init はBubbleteaのイニシャライザです
func (p *ProgressModel) Init() tea.Cmd {
	return nil
}

// Update はBubbleteaの更新関数です
func (p *ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			p.Quitting = true
			return p, tea.Quit
		}

	case ProgressUpdate:
		if msg.Current > 0 {
			p.SetCurrent(msg.Current)
		}
		if msg.Log != "" {
			p.AddLog(msg.Log)
		}

	case ProgressComplete:
		p.Complete()
		p.SetResult(msg.Result)

	case ProgressError:
		p.Failed(msg.Error)
	}

	return p, nil
}

// View はBubbleteaの表示関数です
func (p *ProgressModel) View() string {
	var s strings.Builder

	// タイトルを表示
	s.WriteString(TitleStyle.Render(p.Title) + "\n\n")

	// ステータス表示
	statusText := fmt.Sprintf("ステータス: %s", p.Status)
	switch p.Status {
	case StatusRunning:
		s.WriteString(TextStyle.Render(statusText) + "\n")
	case StatusPaused:
		s.WriteString(SubtextStyle.Render(statusText) + "\n")
	case StatusCompleted:
		s.WriteString(SuccessStyle.Render(statusText) + "\n")
	case StatusFailed:
		s.WriteString(ErrorStyle.Render(statusText) + "\n")
	}

	// 進捗バーの表示
	if p.Total > 0 {
		percent := float64(p.Current) / float64(p.Total)
		percentStr := fmt.Sprintf(" %3.0f%% ", percent*100)

		// 進捗バーの実体部分の長さを計算
		barWidth := p.Width - len(percentStr)
		completedWidth := int(float64(barWidth) * percent)

		// 進捗バーを構築
		bar := ""
		if completedWidth > 0 {
			bar += strings.Repeat("█", completedWidth)
		}
		if completedWidth < barWidth {
			bar += strings.Repeat("░", barWidth-completedWidth)
		}

		// スタイルを適用して表示
		completed := lipgloss.NewStyle().Foreground(lipgloss.Color("#2196F3")).Render(bar[:completedWidth])
		remaining := lipgloss.NewStyle().Foreground(lipgloss.Color("#B0BEC5")).Render(bar[completedWidth:])
		s.WriteString(completed + remaining + "\n\n")

		// カウンター表示
		counterText := fmt.Sprintf("%d / %d 完了", p.Current, p.Total)
		s.WriteString(TextStyle.Render(counterText) + "\n\n")
	}

	// 経過時間の表示
	elapsed := time.Since(p.StartTime)
	if !p.Running && !p.FinishTime.IsZero() {
		elapsed = p.FinishTime.Sub(p.StartTime)
	}
	elapsedText := fmt.Sprintf("経過時間: %s", formatDuration(elapsed))
	s.WriteString(TextStyle.Render(elapsedText) + "\n\n")

	// ログメッセージ表示
	s.WriteString(TitleStyle.Render("ログ") + "\n")

	// ログが空の場合のメッセージ
	if len(p.Logs) == 0 {
		s.WriteString(SubtextStyle.Render("ログはまだありません") + "\n")
	} else {
		// 直近のログを表示
		for _, log := range p.Logs {
			s.WriteString(TextStyle.Render(log) + "\n")
		}
	}

	// エラー表示
	if p.Error != nil {
		s.WriteString("\n" + ErrorStyle.Render(fmt.Sprintf("エラー: %v", p.Error)) + "\n")
	}

	// 結果サマリー表示
	if p.Status == StatusCompleted && p.Result != "" {
		s.WriteString("\n" + SectionStyle.Render(p.Result) + "\n")
	}

	// 操作ガイド
	s.WriteString("\n")
	if p.Running {
		s.WriteString(SubtextStyle.Render("処理中... q: 中断"))
	} else {
		s.WriteString(SubtextStyle.Render("q: 終了"))
	}

	return s.String()
}

// formatDuration は時間を見やすい形式にフォーマットします
func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}

// RunProgress は進捗処理を開始して更新コマンドを返します
func RunProgress(model *ProgressModel, updateCh <-chan ProgressUpdate, errCh <-chan error) tea.Cmd {
	return func() tea.Msg {
		for {
			select {
			case update, ok := <-updateCh:
				if !ok {
					// チャネルが閉じられた場合は完了
					return ProgressComplete{Result: "処理が完了しました"}
				}
				return update

			case err, ok := <-errCh:
				if !ok {
					continue
				}
				return ProgressError{Error: err}
			}
		}
	}
}
