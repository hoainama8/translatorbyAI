package main

import (
	"sync/atomic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

type translationActivity struct {
	inProgress atomic.Bool
}

func (a *translationActivity) setInProgress(inProgress bool) {
	a.inProgress.Store(inProgress)
}

func installTranslationCloseGuard(window fyne.Window, activity *translationActivity) {
	window.SetCloseIntercept(func() {
		if activity.inProgress.Load() {
			dialog.ShowInformation(
				"Đang dịch",
				"Vui lòng chờ quá trình dịch hoàn tất trước khi đóng ứng dụng.",
				window,
			)
			return
		}

		window.Close()
	})
}
