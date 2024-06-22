package main

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"rclone_crypt_desktop/cryptadapt"
	"rclone_crypt_desktop/xtheme"
)

var (
	logArea = widget.NewMultiLineEntry()
)

func main() {

	theApp := app.New()
	theApp.Settings().SetTheme(&xtheme.XTheme{})

	window := theApp.NewWindow("Window")

	var configTitleLabel = canvas.NewText("配置", color.RGBA{B: 255, A: 255})
	configTitleLabel.TextSize = 20

	var configTitle = container.NewHBox(
		layout.NewSpacer(),
		configTitleLabel,
		layout.NewSpacer(),
	)

	var decryptTitleLabel = canvas.NewText("解密工具", color.RGBA{B: 255, A: 255})
	decryptTitleLabel.TextSize = 20

	var decryptTitle = container.NewHBox(
		layout.NewSpacer(),
		decryptTitleLabel,
		layout.NewSpacer(),
	)

	var encryptTitleLabel = canvas.NewText("加密工具", color.RGBA{B: 255, A: 255})
	encryptTitleLabel.TextSize = 20

	var encryptTitle = container.NewHBox(
		layout.NewSpacer(),
		encryptTitleLabel,
		layout.NewSpacer(),
	)

	var passwordEntry = widget.NewEntry()
	passwordEntry.SetPlaceHolder("（必填）")
	passwordEntry.Validator = func(text string) error {
		if text == "" {
			return errors.New("password is empty")
		} else {
			return nil
		}
	}
	var saltEntry = widget.NewEntry()
	saltEntry.SetPlaceHolder("（选填）")
	var filenameEncryptionSelect = widget.NewSelect(cryptadapt.FilenameEncryptionList, nil)
	filenameEncryptionSelect.SetSelectedIndex(0)
	var filenameEncodingSelect = widget.NewSelect(cryptadapt.FilenameEncodingList, nil)
	filenameEncodingSelect.SetSelectedIndex(0)

	var configArea = container.NewVBox(
		configTitle,
		container.NewHBox(
			widget.NewLabel("密码"),
			layout.NewSpacer(),
		),
		passwordEntry,
		container.NewHBox(
			widget.NewLabel("盐"),
			layout.NewSpacer(),
		),
		saltEntry,
		container.NewHBox(
			widget.NewLabel("文件名加密方式"),
			layout.NewSpacer(),
		),
		filenameEncryptionSelect,
		container.NewHBox(
			widget.NewLabel("文件名编码方式"),
			layout.NewSpacer(),
		),
		filenameEncodingSelect,
		widget.NewLabel(""),
	)

	var decryptInputEntry = widget.NewEntry()
	var decryptOutputEntry = widget.NewEntry()

	var decryptChooseFileButton = widget.NewButton("选择文件", func() {
		dialog.NewFileOpen(func(readCloser fyne.URIReadCloser, err error) {
			if readCloser != nil {
				decryptInputEntry.SetText(readCloser.URI().Path())
			}
		}, window).Show()
	})
	var decryptChooseDirButton = widget.NewButton("选择文件夹", func() {
		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri != nil {
				decryptInputEntry.SetText(uri.Path())
			}
		}, window).Show()
	})

	var decryptInputContainer = container.NewVBox(
		container.NewHBox(
			widget.NewLabel("源文件路径"),
			layout.NewSpacer(),
			decryptChooseFileButton,
			decryptChooseDirButton,
		),
		decryptInputEntry,
	)

	var decryptOutputChooseDirButton = widget.NewButton("选择文件夹", func() {
		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri != nil {
				decryptOutputEntry.SetText(uri.Path())
			}
		}, window).Show()
	})

	var decryptOutputContainer = container.NewVBox(
		container.NewHBox(
			widget.NewLabel("输出文件路径"),
			layout.NewSpacer(),
			decryptOutputChooseDirButton,
		),
		decryptOutputEntry,
	)

	var decryptButton = widget.NewButton("Decrypt", func() {
		var tipText, err = cryptadapt.DecryptFile(decryptInputEntry.Text, decryptOutputEntry.Text, passwordEntry.Text, saltEntry.Text)
		printLog(tipText)
		if err != nil {
			printLog(err.Error())
		}
	})

	var encryptInputEntry = widget.NewEntry()
	var encryptOutputEntry = widget.NewEntry()

	var encryptChooseFileButton = widget.NewButton("选择文件", func() {
		dialog.NewFileOpen(func(readCloser fyne.URIReadCloser, err error) {
			if readCloser != nil {
				encryptInputEntry.SetText(readCloser.URI().Path())
			}
		}, window).Show()
	})
	var encryptChooseDirButton = widget.NewButton("选择文件夹", func() {
		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri != nil {
				encryptInputEntry.SetText(uri.Path())
			}
		}, window).Show()
	})

	var encryptInputContainer = container.NewVBox(
		container.NewHBox(
			widget.NewLabel("源文件路径"),
			layout.NewSpacer(),
			encryptChooseFileButton,
			encryptChooseDirButton,
		),
		encryptInputEntry,
	)

	var encryptOutputChooseDirButton = widget.NewButton("选择文件夹", func() {
		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri != nil {
				encryptOutputEntry.SetText(uri.Path())
			}
		}, window).Show()
	})

	var encryptOutputContainer = container.NewVBox(
		container.NewHBox(
			widget.NewLabel("输出文件路径"),
			layout.NewSpacer(),
			encryptOutputChooseDirButton,
		),
		encryptOutputEntry,
	)

	var encryptButton = widget.NewButton("Encrypt", func() {
		var tipText, err = cryptadapt.EncryptFile(encryptInputEntry.Text, encryptOutputEntry.Text, passwordEntry.Text, saltEntry.Text)
		printLog(tipText)
		if err != nil {
			printLog(err.Error())
		}
	})

	var decryptArea = container.NewVBox(
		decryptTitle,
		decryptInputContainer,
		decryptOutputContainer,
		decryptButton,
		widget.NewLabel(""),
	)

	var encryptArea = container.NewVBox(
		encryptTitle,
		encryptInputContainer,
		encryptOutputContainer,
		encryptButton,
		widget.NewLabel(""),
	)

	var leftArea = container.NewVSplit(configArea, container.NewVSplit(decryptArea, encryptArea))
	leftArea.SetOffset(0.3)

	var splitView = container.NewHSplit(leftArea, logArea)
	splitView.SetOffset(0.7)

	window.SetContent(splitView)

	window.Resize(fyne.NewSize(1500, 1000))

	window.ShowAndRun()
}

func printLog(text string) {
	logArea.Append(text + "\n")
}
