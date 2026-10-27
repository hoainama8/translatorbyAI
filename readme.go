package main

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func TabGuide(tabs *container.AppTabs) {
	content := container.NewVBox(
		guideHeading("Hướng dẫn sử dụng"),
		guideText("Translator by AI sử dụng các mô hình Gemini và Gemma để dịch tài liệu. Làm theo các bước dưới đây để bắt đầu."),
		guideHeading("1. Lấy API key"),
		guideText("Ứng dụng cần API key của Google AI để dịch. Hãy giữ khóa bí mật và không chia sẻ khóa với người khác."),
		widget.NewHyperlink("Xem video hướng dẫn lấy API key", &url.URL{
			Scheme:   "https",
			Host:     "youtu.be",
			Path:     "/EswE96hRt9U",
			RawQuery: "si=HyeNFT2g1GjwgihJ",
		}),
		widget.NewHyperlink("Mở Google AI Studio", &url.URL{
			Scheme: "https",
			Host:   "aistudio.google.com",
		}),
		guideHeading("2. Dịch tài liệu"),
		guideText("1. Đính kèm tài liệu cần dịch ở khu vực điều khiển.\n2. Nhập API key, chọn ngôn ngữ nguồn (hoặc để tự động phát hiện) và ngôn ngữ đích. Ngôn ngữ nguồn và đích không được trùng nhau.\n3. Chọn mô hình AI. Tham khảo tab đơn giá để xem chi phí của từng mô hình.\n4. Nhấn “DỊCH == AI” và chờ xử lý. Kết quả sẽ xuất hiện trong tab “Nội dung dịch”."),
		guideHeading("3. Xem và lưu kết quả"),
		guideText("Dùng các nút xuất tệp để lưu bản dịch dưới dạng PDF, Word, TXT hoặc ODT. Các nút mở thư mục giúp truy cập tài liệu gốc và tài liệu đã dịch. Tab “Lịch sử dịch” xem các lần dịch; tab “Thống kê chi phí” xem lượng token và chi phí đã sử dụng."),
		guideHeading("4. Video hướng dẫn cài đặt và sử dụng"),
		widget.NewHyperlink("Hướng dẫn trên Windows", &url.URL{
			Scheme:   "https",
			Host:     "youtu.be",
			Path:     "/JvDTyStEiWA",
			RawQuery: "si=KL_yRCLZdYjlFSL6",
		}),
		widget.NewHyperlink("Hướng dẫn trên Linux", &url.URL{
			Scheme:   "https",
			Host:     "youtu.be",
			Path:     "/ougJP7lzzA4",
			RawQuery: "si=t7wZ1tvSig-1021Y1",
		}),
	)
	guideContent := container.NewScroll(container.NewPadded(content))
	tabs.Append(container.NewTabItem("Hướng dẫn sử dụng", container.NewThemeOverride(
		guideContent,
		outputEditorTheme{Theme: theme.DefaultTheme()},
	)))
}

func TabCopyright(tabs *container.AppTabs) {
	content := container.NewVBox(
		guideHeading("Bản quyền"),
		guideText("Phần mềm do Ngô Hoài Nam phát triển bằng ngôn ngữ Go và được cung cấp miễn phí cho mục đích phi thương mại. Bạn có thể dùng phần mềm để tạo tài liệu hoặc sách đã dịch mà không phải trả phí bản quyền cho phần mềm; quyền sử dụng các tài liệu do bạn dịch do bạn tự quyết định."),
		guideHeading("Tài trợ và tip"),
		guideText("Nhà tài trợ có thể giới thiệu thông tin sản phẩm và thương hiệu của mình tại phần Tài trợ trong phần mềm."),
		guideText("Tip/ủng hộ là hoàn toàn tự nguyện, không bắt buộc và không phải phí sử dụng phần mềm."),
		guideText("Ngân hàng: Vietcombank, chi nhánh Đồng Nai\nSố tài khoản: 0121002153944\nChủ tài khoản: Ngô Hoài Nam"),
	)
	copyrightContent := container.NewScroll(container.NewPadded(content))
	tabs.Append(container.NewTabItem("Bản quyền", container.NewThemeOverride(
		copyrightContent,
		outputEditorTheme{Theme: theme.DefaultTheme()},
	)))
}

func guideHeading(text string) *widget.Label {
	label := widget.NewLabel(text)
	label.TextStyle = fyne.TextStyle{Bold: true}
	label.Wrapping = fyne.TextWrapWord
	return label
}

func guideText(text string) *widget.Label {
	label := widget.NewLabel(text)
	label.Wrapping = fyne.TextWrapWord
	return label
}
