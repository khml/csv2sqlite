package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// FormField はフォームのフィールドを表す構造体です
type FormField struct {
	Label       string              // フィールドのラベル
	Input       textinput.Model     // テキスト入力モデル
	Placeholder string              // プレースホルダーテキスト
	Required    bool                // 必須フィールドかどうか
	Validator   func(string) string // バリデーション関数、エラーメッセージを返すか空文字で成功
	ErrorMsg    string              // バリデーションエラーメッセージ
}

// Form はフォーム入力コンポーネントです
type Form struct {
	Fields      []FormField // フォームのフィールド一覧
	FocusIndex  int         // 現在フォーカスのあるフィールドのインデックス
	Validated   bool        // バリデーション済みかどうか
	SubmitLabel string      // 送信ボタンのラベル
	CancelLabel string      // キャンセルボタンのラベル
	Submitted   bool        // 送信されたかどうか
	Canceled    bool        // キャンセルされたかどうか
	Title       string      // フォームのタイトル
}

// NewForm は新しいフォームを生成します
func NewForm(title string) *Form {
	return &Form{
		Fields:      []FormField{},
		FocusIndex:  0,
		Validated:   false,
		SubmitLabel: "送信",
		CancelLabel: "キャンセル",
		Title:       title,
	}
}

// AddField はフォームにフィールドを追加します
func (f *Form) AddField(label, placeholder string, required bool, validator func(string) string) *Form {
	input := textinput.New()
	input.Placeholder = placeholder
	input.PromptStyle = InputFieldStyle
	input.TextStyle = TextStyle

	if len(f.Fields) == 0 {
		input.Focus()
	}

	f.Fields = append(f.Fields, FormField{
		Label:       label,
		Input:       input,
		Placeholder: placeholder,
		Required:    required,
		Validator:   validator,
	})

	return f
}

// SetValue はフィールドの値を設定します
func (f *Form) SetValue(index int, value string) {
	if index >= 0 && index < len(f.Fields) {
		f.Fields[index].Input.SetValue(value)
	}
}

// Value はフィールドの値を取得します
func (f *Form) Value(index int) string {
	if index >= 0 && index < len(f.Fields) {
		return f.Fields[index].Input.Value()
	}
	return ""
}

// SetSubmitLabel は送信ボタンのラベルを設定します
func (f *Form) SetSubmitLabel(label string) *Form {
	f.SubmitLabel = label
	return f
}

// SetCancelLabel はキャンセルボタンのラベルを設定します
func (f *Form) SetCancelLabel(label string) *Form {
	f.CancelLabel = label
	return f
}

// GetValues は全フィールドの値をマップで取得します
func (f *Form) GetValues() map[string]string {
	values := make(map[string]string)
	for _, field := range f.Fields {
		values[field.Label] = field.Input.Value()
	}
	return values
}

// Validate はフォームをバリデーションします
func (f *Form) Validate() bool {
	valid := true

	for i := range f.Fields {
		field := &f.Fields[i]
		field.ErrorMsg = ""

		// 必須フィールドの確認
		if field.Required && field.Input.Value() == "" {
			field.ErrorMsg = "必須項目です"
			valid = false
			continue
		}

		// カスタムバリデーション
		if field.Validator != nil && field.Input.Value() != "" {
			if errMsg := field.Validator(field.Input.Value()); errMsg != "" {
				field.ErrorMsg = errMsg
				valid = false
			}
		}
	}

	f.Validated = true
	return valid
}

// Init はBubbleteaのイニシャライザです
func (f *Form) Init() tea.Cmd {
	return textinput.Blink
}

// Update はBubbleteaの更新関数です
func (f *Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			f.Canceled = true
			return f, tea.Quit

		case "tab", "shift+tab":
			// フォーカスの移動
			direction := 1
			if msg.String() == "shift+tab" {
				direction = -1
			}

			f.FocusIndex = (f.FocusIndex + direction) % len(f.Fields)
			if f.FocusIndex < 0 {
				f.FocusIndex = len(f.Fields) - 1
			}

			cmds := make([]tea.Cmd, len(f.Fields))
			for i := range f.Fields {
				if i == f.FocusIndex {
					cmds[i] = f.Fields[i].Input.Focus()
				} else {
					f.Fields[i].Input.Blur()
				}
			}

			return f, tea.Batch(cmds...)

		case "enter":
			// 最後のフィールドでEnterを押した場合、送信アクション
			if f.FocusIndex == len(f.Fields)-1 {
				if f.Validate() {
					f.Submitted = true
					return f, tea.Quit
				}
			} else {
				// 次のフィールドにフォーカスを移動
				f.FocusIndex = (f.FocusIndex + 1) % len(f.Fields)

				cmds := make([]tea.Cmd, len(f.Fields))
				for i := range f.Fields {
					if i == f.FocusIndex {
						cmds[i] = f.Fields[i].Input.Focus()
					} else {
						f.Fields[i].Input.Blur()
					}
				}

				return f, tea.Batch(cmds...)
			}
		}
	}

	// 現在フォーカスのあるフィールドの更新
	if f.FocusIndex >= 0 && f.FocusIndex < len(f.Fields) {
		var cmd tea.Cmd
		f.Fields[f.FocusIndex].Input, cmd = f.Fields[f.FocusIndex].Input.Update(msg)
		return f, cmd
	}

	return f, nil
}

// View はBubbleteaの表示関数です
func (f *Form) View() string {
	var s strings.Builder

	// タイトルを表示
	s.WriteString(TitleStyle.Render(f.Title) + "\n\n")

	// 各フィールドを表示
	for i, field := range f.Fields {
		// ラベルを表示
		labelText := field.Label
		if field.Required {
			labelText += " *"
		}
		s.WriteString(TextStyle.Render(labelText) + "\n")

		// 入力フィールドを表示
		s.WriteString(field.Input.View() + "\n")

		// エラーメッセージがあれば表示
		if field.ErrorMsg != "" {
			s.WriteString(ErrorStyle.Render(field.ErrorMsg) + "\n")
		}

		if i < len(f.Fields)-1 {
			s.WriteString("\n")
		}
	}

	// 送信とキャンセルのボタン
	s.WriteString("\n")
	submitBtn := f.SubmitLabel
	cancelBtn := f.CancelLabel

	if f.FocusIndex == len(f.Fields) {
		s.WriteString(SelectedButtonStyle.Render(submitBtn) + " ")
		s.WriteString(ButtonStyle.Render(cancelBtn))
	} else if f.FocusIndex == len(f.Fields)+1 {
		s.WriteString(ButtonStyle.Render(submitBtn) + " ")
		s.WriteString(SelectedButtonStyle.Render(cancelBtn))
	} else {
		s.WriteString(ButtonStyle.Render(submitBtn) + " ")
		s.WriteString(ButtonStyle.Render(cancelBtn))
	}

	// 操作ガイドを表示
	s.WriteString("\n\n")
	s.WriteString(SubtextStyle.Render("Tab: 次へ  Shift+Tab: 前へ  Enter: 確定  Esc: キャンセル"))

	return s.String()
}

// RunForm はフォームを実行して結果を返します
func RunForm(form *Form) (map[string]string, bool, error) {
	p := tea.NewProgram(form)
	_, err := p.Run()
	if err != nil {
		return nil, false, err
	}

	// キャンセルされた場合
	if form.Canceled {
		return nil, false, nil
	}

	return form.GetValues(), form.Submitted, nil
}

// 一般的なバリデーション関数
// 数値のみを許可するバリデーション
func NumberValidator(input string) string {
	for _, r := range input {
		if r < '0' || r > '9' {
			return "数値のみ入力可能です"
		}
	}
	return ""
}

// 最小長バリデーション
func MinLengthValidator(minLength int) func(string) string {
	return func(input string) string {
		if len(input) < minLength {
			return fmt.Sprintf("最低%d文字必要です", minLength)
		}
		return ""
	}
}

// 最大長バリデーション
func MaxLengthValidator(maxLength int) func(string) string {
	return func(input string) string {
		if len(input) > maxLength {
			return fmt.Sprintf("最大%d文字までです", maxLength)
		}
		return ""
	}
}

// 複数のバリデーターを結合
func CombineValidators(validators ...func(string) string) func(string) string {
	return func(input string) string {
		for _, validator := range validators {
			if err := validator(input); err != "" {
				return err
			}
		}
		return ""
	}
}
