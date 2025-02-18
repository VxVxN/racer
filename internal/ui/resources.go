package ui

import (
	"image/color"
	"strconv"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	backgroundColor = "131a22"

	textIdleColor     = "dff4ff"
	textDisabledColor = "5a7a91"

	labelIdleColor     = textIdleColor
	labelDisabledColor = textDisabledColor

	buttonIdleColor     = textIdleColor
	buttonDisabledColor = labelDisabledColor

	listSelectedBackground         = "4b687a"
	listDisabledSelectedBackground = "2a3944"

	listFocusedBackground = "2a3944"

	headerColor = textIdleColor

	textInputCaretColor         = "e7c34b"
	textInputDisabledCaretColor = "766326"

	toolTipColor = backgroundColor

	separatorColor = listDisabledSelectedBackground
)

type UiResources struct {
	Fonts *Fonts

	Background *image.NineSlice

	SeparatorColor color.Color

	Text        *textResources
	Button      *buttonResources
	Label       *labelResources
	Checkbox    *checkboxResources
	comboButton *comboButtonResources
	list        *listResources
	Slider      *sliderResources
	progressBar *progressBarResources
	panel       *panelResources
	tabBook     *tabBookResources
	Header      *headerResources
	TextInput   *textInputResources
	textArea    *textAreaResources
	toolTip     *toolTipResources
}

type textResources struct {
	IdleColor     color.Color
	DisabledColor color.Color
	Face          text.Face
	TitleFace     text.Face
	BigTitleFace  text.Face
	SmallFace     text.Face
}

type buttonResources struct {
	Image   *widget.ButtonImage
	Text    *widget.ButtonTextColor
	Face    text.Face
	Padding widget.Insets
}

type checkboxResources struct {
	image   *widget.ButtonImage
	graphic *widget.CheckboxGraphicImage
	spacing int
}

type labelResources struct {
	Text *widget.LabelColor
	Face text.Face
}

type comboButtonResources struct {
	image   *widget.ButtonImage
	text    *widget.ButtonTextColor
	face    text.Face
	graphic *widget.ButtonImageImage
	padding widget.Insets
}

type listResources struct {
	image        *widget.ScrollContainerImage
	track        *widget.SliderTrackImage
	trackPadding widget.Insets
	handle       *widget.ButtonImage
	handleSize   int
	face         text.Face
	entry        *widget.ListEntryColor
	entryPadding widget.Insets
}

type sliderResources struct {
	TrackImage *widget.SliderTrackImage
	Handle     *widget.ButtonImage
	HandleSize int
}

type progressBarResources struct {
	trackImage *widget.ProgressBarImage
	fillImage  *widget.ProgressBarImage
}

type panelResources struct {
	image    *image.NineSlice
	titleBar *image.NineSlice
	padding  widget.Insets
}

type tabBookResources struct {
	buttonFace    text.Face
	buttonText    *widget.ButtonTextColor
	buttonPadding widget.Insets
}

type headerResources struct {
	background *image.NineSlice
	padding    widget.Insets
	face       text.Face
	color      color.Color
}

type textInputResources struct {
	Image   *widget.TextInputImage
	Padding widget.Insets
	Face    text.Face
	Color   *widget.TextInputColor
}

type textAreaResources struct {
	image        *widget.ScrollContainerImage
	track        *widget.SliderTrackImage
	trackPadding widget.Insets
	handle       *widget.ButtonImage
	handleSize   int
	face         text.Face
	entryPadding widget.Insets
}

type toolTipResources struct {
	background *image.NineSlice
	padding    widget.Insets
	face       text.Face
	color      color.Color
}

func NewUIResources() (*UiResources, error) {
	background := image.NewNineSliceColor(hexToColor(backgroundColor))

	fonts, err := loadFonts()
	if err != nil {
		return nil, err
	}

	button, err := newButtonResources(fonts)
	if err != nil {
		return nil, err
	}

	checkbox, err := newCheckboxResources()
	if err != nil {
		return nil, err
	}

	comboButton, err := newComboButtonResources(fonts)
	if err != nil {
		return nil, err
	}

	list, err := newListResources(fonts)
	if err != nil {
		return nil, err
	}

	slider, err := newSliderResources()
	if err != nil {
		return nil, err
	}

	progressBar, err := newProgressBarResources()
	if err != nil {
		return nil, err
	}

	panel, err := newPanelResources()
	if err != nil {
		return nil, err
	}

	tabBook, err := newTabBookResources(fonts)
	if err != nil {
		return nil, err
	}

	header, err := newHeaderResources(fonts)
	if err != nil {
		return nil, err
	}

	textInput, err := newTextInputResources(fonts)
	if err != nil {
		return nil, err
	}
	textArea, err := newTextAreaResources(fonts)
	if err != nil {
		return nil, err
	}
	toolTip, err := newToolTipResources(fonts)
	if err != nil {
		return nil, err
	}

	return &UiResources{
		Fonts: fonts,

		Background: background,

		SeparatorColor: hexToColor(separatorColor),

		Text: &textResources{
			IdleColor:     hexToColor(textIdleColor),
			DisabledColor: hexToColor(textDisabledColor),
			Face:          fonts.face,
			TitleFace:     fonts.titleFace,
			BigTitleFace:  fonts.bigTitleFace,
			SmallFace:     fonts.toolTipFace,
		},

		Button:      button,
		Label:       newLabelResources(fonts),
		Checkbox:    checkbox,
		comboButton: comboButton,
		list:        list,
		Slider:      slider,
		panel:       panel,
		tabBook:     tabBook,
		Header:      header,
		TextInput:   textInput,
		toolTip:     toolTip,
		textArea:    textArea,
		progressBar: progressBar,
	}, nil
}

func newButtonResources(fonts *Fonts) (*buttonResources, error) {
	idle, err := loadImageNineSlice("assets/graphics/button-idle.png", 12, 0)
	if err != nil {
		return nil, err
	}

	hover, err := loadImageNineSlice("assets/graphics/button-hover.png", 12, 0)
	if err != nil {
		return nil, err
	}
	pressed_hover, err := loadImageNineSlice("assets/graphics/button-selected-hover.png", 12, 0)
	if err != nil {
		return nil, err
	}
	pressed, err := loadImageNineSlice("assets/graphics/button-pressed.png", 12, 0)
	if err != nil {
		return nil, err
	}

	disabled, err := loadImageNineSlice("assets/graphics/button-disabled.png", 12, 0)
	if err != nil {
		return nil, err
	}

	i := &widget.ButtonImage{
		Idle:         idle,
		Hover:        hover,
		Pressed:      pressed,
		PressedHover: pressed_hover,
		Disabled:     disabled,
	}

	return &buttonResources{
		Image: i,

		Text: &widget.ButtonTextColor{
			Idle:     hexToColor(buttonIdleColor),
			Disabled: hexToColor(buttonDisabledColor),
		},

		Face: fonts.face,

		Padding: widget.Insets{
			Left:  30,
			Right: 30,
		},
	}, nil
}

func newCheckboxResources() (*checkboxResources, error) {
	idle, err := loadImageNineSlice("assets/graphics/checkbox-idle.png", 20, 0)
	if err != nil {
		return nil, err
	}

	hover, err := loadImageNineSlice("assets/graphics/checkbox-hover.png", 20, 0)
	if err != nil {
		return nil, err
	}

	disabled, err := loadImageNineSlice("assets/graphics/checkbox-disabled.png", 20, 0)
	if err != nil {
		return nil, err
	}

	checked, err := loadGraphicImages("assets/graphics/checkbox-checked-idle.png", "assets/graphics/checkbox-checked-disabled.png")
	if err != nil {
		return nil, err
	}

	unchecked, err := loadGraphicImages("assets/graphics/checkbox-unchecked-idle.png", "assets/graphics/checkbox-unchecked-disabled.png")
	if err != nil {
		return nil, err
	}

	greyed, err := loadGraphicImages("assets/graphics/checkbox-greyed-idle.png", "assets/graphics/checkbox-greyed-disabled.png")
	if err != nil {
		return nil, err
	}

	return &checkboxResources{
		image: &widget.ButtonImage{
			Idle:     idle,
			Hover:    hover,
			Pressed:  hover,
			Disabled: disabled,
		},

		graphic: &widget.CheckboxGraphicImage{
			Checked:   checked,
			Unchecked: unchecked,
			Greyed:    greyed,
		},

		spacing: 10,
	}, nil
}

func newLabelResources(fonts *Fonts) *labelResources {
	return &labelResources{
		Text: &widget.LabelColor{
			Idle:     hexToColor(labelIdleColor),
			Disabled: hexToColor(labelDisabledColor),
		},

		Face: fonts.face,
	}
}

func newComboButtonResources(fonts *Fonts) (*comboButtonResources, error) {
	idle, err := loadImageNineSlice("assets/graphics/combo-button-idle.png", 12, 0)
	if err != nil {
		return nil, err
	}

	hover, err := loadImageNineSlice("assets/graphics/combo-button-hover.png", 12, 0)
	if err != nil {
		return nil, err
	}

	pressed, err := loadImageNineSlice("assets/graphics/combo-button-pressed.png", 12, 0)
	if err != nil {
		return nil, err
	}

	disabled, err := loadImageNineSlice("assets/graphics/combo-button-disabled.png", 12, 0)
	if err != nil {
		return nil, err
	}

	i := &widget.ButtonImage{
		Idle:     idle,
		Hover:    hover,
		Pressed:  pressed,
		Disabled: disabled,
	}

	arrowDown, err := loadGraphicImages("assets/graphics/arrow-down-idle.png", "assets/graphics/arrow-down-disabled.png")
	if err != nil {
		return nil, err
	}

	return &comboButtonResources{
		image: i,

		text: &widget.ButtonTextColor{
			Idle:     hexToColor(buttonIdleColor),
			Disabled: hexToColor(buttonDisabledColor),
		},

		face:    fonts.face,
		graphic: arrowDown,

		padding: widget.Insets{
			Left:  30,
			Right: 30,
		},
	}, nil
}

func newListResources(fonts *Fonts) (*listResources, error) {
	idle, err := newImageFromFile("assets/graphics/list-idle.png")
	if err != nil {
		return nil, err
	}

	disabled, err := newImageFromFile("assets/graphics/list-disabled.png")
	if err != nil {
		return nil, err
	}

	mask, err := newImageFromFile("assets/graphics/list-mask.png")
	if err != nil {
		return nil, err
	}

	trackIdle, err := newImageFromFile("assets/graphics/list-track-idle.png")
	if err != nil {
		return nil, err
	}

	trackDisabled, err := newImageFromFile("assets/graphics/list-track-disabled.png")
	if err != nil {
		return nil, err
	}

	handleIdle, err := newImageFromFile("assets/graphics/slider-handle-idle.png")
	if err != nil {
		return nil, err
	}

	handleHover, err := newImageFromFile("assets/graphics/slider-handle-hover.png")
	if err != nil {
		return nil, err
	}

	return &listResources{
		image: &widget.ScrollContainerImage{
			Idle:     image.NewNineSlice(idle, [3]int{25, 12, 22}, [3]int{25, 12, 25}),
			Disabled: image.NewNineSlice(disabled, [3]int{25, 12, 22}, [3]int{25, 12, 25}),
			Mask:     image.NewNineSlice(mask, [3]int{26, 10, 23}, [3]int{26, 10, 26}),
		},

		track: &widget.SliderTrackImage{
			Idle:     image.NewNineSlice(trackIdle, [3]int{5, 0, 0}, [3]int{25, 12, 25}),
			Hover:    image.NewNineSlice(trackIdle, [3]int{5, 0, 0}, [3]int{25, 12, 25}),
			Disabled: image.NewNineSlice(trackDisabled, [3]int{0, 5, 0}, [3]int{25, 12, 25}),
		},

		trackPadding: widget.Insets{
			Top:    5,
			Bottom: 24,
		},

		handle: &widget.ButtonImage{
			Idle:     image.NewNineSliceSimple(handleIdle, 0, 5),
			Hover:    image.NewNineSliceSimple(handleHover, 0, 5),
			Pressed:  image.NewNineSliceSimple(handleHover, 0, 5),
			Disabled: image.NewNineSliceSimple(handleIdle, 0, 5),
		},

		handleSize: 5,
		face:       fonts.face,

		entry: &widget.ListEntryColor{
			Unselected:         hexToColor(textIdleColor),
			DisabledUnselected: hexToColor(textDisabledColor),

			Selected:         hexToColor(textIdleColor),
			DisabledSelected: hexToColor(textDisabledColor),

			SelectedBackground:         hexToColor(listSelectedBackground),
			DisabledSelectedBackground: hexToColor(listDisabledSelectedBackground),

			FocusedBackground:         hexToColor(listFocusedBackground),
			SelectedFocusedBackground: hexToColor(listSelectedBackground),
		},

		entryPadding: widget.Insets{
			Left:   30,
			Right:  30,
			Top:    2,
			Bottom: 2,
		},
	}, nil
}

func newSliderResources() (*sliderResources, error) {
	idle, err := newImageFromFile("assets/graphics/slider-track-idle.png")
	if err != nil {
		return nil, err
	}

	disabled, err := newImageFromFile("assets/graphics/slider-track-disabled.png")
	if err != nil {
		return nil, err
	}

	handleIdle, err := newImageFromFile("assets/graphics/slider-handle-idle.png")
	if err != nil {
		return nil, err
	}

	handleHover, err := newImageFromFile("assets/graphics/slider-handle-hover.png")
	if err != nil {
		return nil, err
	}

	handleDisabled, err := newImageFromFile("assets/graphics/slider-handle-disabled.png")
	if err != nil {
		return nil, err
	}

	return &sliderResources{
		TrackImage: &widget.SliderTrackImage{
			Idle:     image.NewNineSlice(idle, [3]int{0, 19, 0}, [3]int{6, 0, 0}),
			Hover:    image.NewNineSlice(idle, [3]int{0, 19, 0}, [3]int{6, 0, 0}),
			Disabled: image.NewNineSlice(disabled, [3]int{0, 19, 0}, [3]int{6, 0, 0}),
		},

		Handle: &widget.ButtonImage{
			Idle:     image.NewNineSliceSimple(handleIdle, 0, 5),
			Hover:    image.NewNineSliceSimple(handleHover, 0, 5),
			Pressed:  image.NewNineSliceSimple(handleHover, 0, 5),
			Disabled: image.NewNineSliceSimple(handleDisabled, 0, 5),
		},

		HandleSize: 6,
	}, nil
}

func newProgressBarResources() (*progressBarResources, error) {
	idle, err := newImageFromFile("assets/graphics/progressbar-track-idle.png")
	if err != nil {
		return nil, err
	}
	fill_idle, err := newImageFromFile("assets/graphics/progressbar-fill-idle.png")
	if err != nil {
		return nil, err
	}
	disabled, err := newImageFromFile("assets/graphics/slider-track-disabled.png")
	if err != nil {
		return nil, err
	}

	return &progressBarResources{
		trackImage: &widget.ProgressBarImage{
			Idle:     image.NewNineSlice(idle, [3]int{4, 11, 4}, [3]int{2, 2, 2}),
			Hover:    image.NewNineSlice(idle, [3]int{4, 11, 4}, [3]int{2, 2, 2}),
			Disabled: image.NewNineSlice(disabled, [3]int{4, 11, 4}, [3]int{2, 2, 2}),
		},

		fillImage: &widget.ProgressBarImage{
			Idle:     image.NewNineSlice(fill_idle, [3]int{4, 11, 4}, [3]int{2, 2, 2}),
			Hover:    image.NewNineSlice(fill_idle, [3]int{4, 11, 4}, [3]int{2, 2, 2}),
			Disabled: image.NewNineSlice(fill_idle, [3]int{4, 11, 4}, [3]int{2, 2, 2}),
		},
	}, nil
}
func newPanelResources() (*panelResources, error) {
	i, err := loadImageNineSlice("assets/graphics/panel-idle.png", 10, 10)
	if err != nil {
		return nil, err
	}
	t, err := loadImageNineSlice("assets/graphics/titlebar-idle.png", 10, 10)
	if err != nil {
		return nil, err
	}
	return &panelResources{
		image:    i,
		titleBar: t,
		padding: widget.Insets{
			Left:   30,
			Right:  30,
			Top:    20,
			Bottom: 20,
		},
	}, nil
}

func newTabBookResources(fonts *Fonts) (*tabBookResources, error) {

	return &tabBookResources{
		buttonFace: fonts.face,

		buttonText: &widget.ButtonTextColor{
			Idle:     hexToColor(buttonIdleColor),
			Disabled: hexToColor(buttonDisabledColor),
		},

		buttonPadding: widget.Insets{
			Left:  30,
			Right: 30,
		},
	}, nil
}

func newHeaderResources(fonts *Fonts) (*headerResources, error) {
	bg, err := loadImageNineSlice("assets/graphics/header.png", 446, 9)
	if err != nil {
		return nil, err
	}

	return &headerResources{
		background: bg,

		padding: widget.Insets{
			Left:   25,
			Right:  25,
			Top:    4,
			Bottom: 4,
		},

		face:  fonts.bigTitleFace,
		color: hexToColor(headerColor),
	}, nil
}

func newTextInputResources(fonts *Fonts) (*textInputResources, error) {
	idle, err := newImageFromFile("assets/graphics/text-input-idle.png")
	if err != nil {
		return nil, err
	}

	disabled, err := newImageFromFile("assets/graphics/text-input-disabled.png")
	if err != nil {
		return nil, err
	}

	return &textInputResources{
		Image: &widget.TextInputImage{
			Idle:     image.NewNineSlice(idle, [3]int{9, 14, 6}, [3]int{9, 14, 6}),
			Disabled: image.NewNineSlice(disabled, [3]int{9, 14, 6}, [3]int{9, 14, 6}),
		},

		Padding: widget.Insets{
			Left:   8,
			Right:  8,
			Top:    4,
			Bottom: 4,
		},

		Face: fonts.face,

		Color: &widget.TextInputColor{
			Idle:          hexToColor(textIdleColor),
			Disabled:      hexToColor(textDisabledColor),
			Caret:         hexToColor(textInputCaretColor),
			DisabledCaret: hexToColor(textInputDisabledCaretColor),
		},
	}, nil
}

func newTextAreaResources(fonts *Fonts) (*textAreaResources, error) {
	idle, err := newImageFromFile("assets/graphics/list-idle.png")
	if err != nil {
		return nil, err
	}

	disabled, err := newImageFromFile("assets/graphics/list-disabled.png")
	if err != nil {
		return nil, err
	}

	mask, err := newImageFromFile("assets/graphics/list-mask.png")
	if err != nil {
		return nil, err
	}

	trackIdle, err := newImageFromFile("assets/graphics/list-track-idle.png")
	if err != nil {
		return nil, err
	}

	trackDisabled, err := newImageFromFile("assets/graphics/list-track-disabled.png")
	if err != nil {
		return nil, err
	}

	handleIdle, err := newImageFromFile("assets/graphics/slider-handle-idle.png")
	if err != nil {
		return nil, err
	}

	handleHover, err := newImageFromFile("assets/graphics/slider-handle-hover.png")
	if err != nil {
		return nil, err
	}

	return &textAreaResources{
		image: &widget.ScrollContainerImage{
			Idle:     image.NewNineSlice(idle, [3]int{25, 12, 22}, [3]int{25, 12, 25}),
			Disabled: image.NewNineSlice(disabled, [3]int{25, 12, 22}, [3]int{25, 12, 25}),
			Mask:     image.NewNineSlice(mask, [3]int{26, 10, 23}, [3]int{26, 10, 26}),
		},

		track: &widget.SliderTrackImage{
			Idle:     image.NewNineSlice(trackIdle, [3]int{5, 0, 0}, [3]int{25, 12, 25}),
			Hover:    image.NewNineSlice(trackIdle, [3]int{5, 0, 0}, [3]int{25, 12, 25}),
			Disabled: image.NewNineSlice(trackDisabled, [3]int{0, 5, 0}, [3]int{25, 12, 25}),
		},

		trackPadding: widget.Insets{
			Top:    5,
			Bottom: 24,
		},

		handle: &widget.ButtonImage{
			Idle:     image.NewNineSliceSimple(handleIdle, 0, 5),
			Hover:    image.NewNineSliceSimple(handleHover, 0, 5),
			Pressed:  image.NewNineSliceSimple(handleHover, 0, 5),
			Disabled: image.NewNineSliceSimple(handleIdle, 0, 5),
		},

		handleSize: 5,
		face:       fonts.face,

		entryPadding: widget.Insets{
			Left:   30,
			Right:  30,
			Top:    2,
			Bottom: 2,
		},
	}, nil
}

func newToolTipResources(fonts *Fonts) (*toolTipResources, error) {
	bg, err := newImageFromFile("assets/graphics/tool-tip.png")
	if err != nil {
		return nil, err
	}

	return &toolTipResources{
		background: image.NewNineSlice(bg, [3]int{19, 6, 13}, [3]int{19, 5, 13}),

		padding: widget.Insets{
			Left:   15,
			Right:  15,
			Top:    10,
			Bottom: 10,
		},

		face:  fonts.toolTipFace,
		color: hexToColor(toolTipColor),
	}, nil
}

func hexToColor(h string) color.Color {
	u, err := strconv.ParseUint(h, 16, 0)
	if err != nil {
		panic(err)
	}

	return color.NRGBA{
		R: uint8(u & 0xff0000 >> 16),
		G: uint8(u & 0xff00 >> 8),
		B: uint8(u & 0xff),
		A: 255,
	}
}
