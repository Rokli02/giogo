package minesweeper

import (
	"fmt"
	"giogo/ui"
	routerModule "giogo/ui/router"
	"giogo/ui/styles"
	"giogo/utils"
	"image"
	"math/rand"
	"time"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

const (
	heightOfHeader = 42
)

type MineField struct {
	BtnSize       image.Point
	BtnMatrix     [][]MineButton
	State         MinesweeperState
	MaxNumOfMines uint16
	RevealedCells uint16
	MarkedCells   uint16

	w                      *app.Window
	router                 *routerModule.Router[ui.ApplicationCycles, string]
	animationDuration      time.Duration
	plantedMines           uint16
	countOfHorizontalCells uint16
	countOfVerticalCells   uint16
	returnHomeClickable    widget.Clickable
}

/*
 * Constructor
 */

func NewMinefield(w *app.Window, router *routerModule.Router[ui.ApplicationCycles, string], horizontalCells, verticalCells, numberOfMines uint16) *MineField {
	mineField := &MineField{
		State:         UNDEFINED,
		BtnSize:       image.Pt(32, 32),
		MaxNumOfMines: numberOfMines,

		w:                      w,
		router:                 router,
		animationDuration:      time.Millisecond * 40,
		plantedMines:           numberOfMines,
		countOfHorizontalCells: horizontalCells,
		countOfVerticalCells:   verticalCells,
		returnHomeClickable:    widget.Clickable{},
	}

	mineField.Initialize()

	return mineField
}

/*
 * Public methods
 */

func (mf *MineField) Initialize() {
	sizeOfMinesweeper := image.Pt(mf.BtnSize.X*int(mf.countOfHorizontalCells), mf.BtnSize.Y*int(mf.countOfVerticalCells)+heightOfHeader)
	mf.w.Option(func(_ unit.Metric, c *app.Config) {
		c.Title = "Minesweeper COPY"
		c.MinSize = sizeOfMinesweeper
		c.MaxSize = sizeOfMinesweeper
		c.Size = sizeOfMinesweeper
		c.Decorated = true
	})

	shouldInitMatrix := len(mf.BtnMatrix) == 0 || len(mf.BtnMatrix) != int(mf.countOfVerticalCells)

	if shouldInitMatrix {
		mf.BtnMatrix = make([][]MineButton, mf.countOfVerticalCells)
	}

	if shouldInitMatrix || len(mf.BtnMatrix[0]) == 0 || len(mf.BtnMatrix[0]) != int(mf.countOfHorizontalCells) {
		for index := range mf.BtnMatrix {
			mf.BtnMatrix[index] = make([]MineButton, mf.countOfHorizontalCells)
		}
	}

	if shouldInitMatrix {
		for rowIndex := range mf.BtnMatrix {
			btnList := &mf.BtnMatrix[rowIndex]

			for colIndex := range *btnList {
				btn := &(*btnList)[colIndex]

				btn.onClick = mf.onButtonClick
				btn.Pos = image.Point{colIndex, rowIndex}
				btn.State = mf.State
				btn.Size = mf.BtnSize
				btn.Hidden = true
				btn.Value = 0
				btn.Marked = false
			}
		}

		mf.Restart()
	}
}

func (mf *MineField) Restart() {
	mf.State = START
	mf.RevealedCells = 0
	mf.MarkedCells = 0

	for rowIndex := range mf.BtnMatrix {
		btnList := &mf.BtnMatrix[rowIndex]

		for colIndex := range *btnList {
			btn := &(*btnList)[colIndex]

			btn.State = mf.State
			btn.Hidden = true
			btn.Value = 0
			btn.Marked = false
		}
	}
}

func (mf *MineField) Close() {
}

func (mf *MineField) Layout(gtx layout.Context) layout.Dimensions {
	for _, event := range gtx.Events(mf) {
		switch evt := event.(type) {
		case key.Event:
			switch evt.Name {
			case "R":
				if evt.Modifiers == key.ModCtrl && evt.State == key.Release {
					mf.Restart()
				}
			}
		case pointer.Event:
			switch evt.Kind {
			case pointer.Press:
				mf.Restart()
			}
		}
	}

	defer utils.SetBackgroundColor(&gtx, styles.BACKGROUND_COLOR).Pop()

	if mf.State == LOADING {
		op.InvalidateOp{At: time.Now().Add(mf.animationDuration >> 1)}.Add(gtx.Ops)
	} else {
		key.InputOp{
			Tag:  mf,
			Keys: key.Set("Ctrl-R"),
		}.Add(gtx.Ops)
	}

	return layout.N.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		flexD := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(mf.headerComponent),
			layout.Rigid(mf.bodyComponent),
		)

		flexD.Size.X = mf.BtnSize.X * int(mf.countOfHorizontalCells)

		return flexD
	})
}

/*
 * Private methods
 */

func (mf *MineField) headerComponent(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min = image.Pt(0, 0)
	gtx.Constraints.Max = image.Pt(int(mf.countOfHorizontalCells)*mf.BtnSize.X, heightOfHeader)

	defer utils.SetBackgroundColor(&gtx, styles.HEADER_BACKGROUND).Pop()

	headerD := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
		layout.Flexed(3, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min = image.Point{gtx.Constraints.Max.Y - 4, gtx.Constraints.Max.Y - 4}
			icon, _ := widget.NewIcon(icons.ActionHome)
			iconColor := styles.BLOOD_ORANGE

			if mf.returnHomeClickable.Clicked(gtx) {
				mf.router.GoTo(routerModule.MenuPage)
			}

			if mf.returnHomeClickable.Pressed() {
				iconColor.R -= 10
				iconColor.G -= 10
				iconColor.B -= 10
			} else if mf.returnHomeClickable.Hovered() {
				iconColor.R += 20
				iconColor.G += 20
				iconColor.B += 20
			}

			mf.returnHomeClickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				pointer.Cursor.Add(pointer.CursorPointer, gtx.Ops)

				return icon.Layout(gtx, iconColor)
			})

			return layout.Dimensions{Size: gtx.Constraints.Max}
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(4).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min = image.Pt(0, 0)
				gtx.Constraints.Max = image.Pt(heightOfHeader-8, heightOfHeader-8)

				paint.FillShape(gtx.Ops, styles.NIGHT_BLACK, clip.Stroke{Width: 2, Path: clip.Ellipse{Min: gtx.Constraints.Min, Max: gtx.Constraints.Max}.Path(gtx.Ops)}.Op())
				defer clip.Ellipse{Min: gtx.Constraints.Min, Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()

				switch mf.State {
				case START:
					fallthrough
				case END:
					fallthrough
				case LOADING:
					fallthrough
				case RUNNING:
					// Sárga neutrális
					paint.ColorOp{Color: styles.YELLOW}.Add(gtx.Ops)
					paint.PaintOp{}.Add(gtx.Ops)

					drawNeutralFace(&gtx)
				case WIN:
					// Zöld mosolygó
					paint.ColorOp{Color: styles.GREEN}.Add(gtx.Ops)
					paint.PaintOp{}.Add(gtx.Ops)

					drawHappyFace(&gtx)
				case LOSE:
					// Piros mérges/szomorú
					paint.ColorOp{Color: styles.RED}.Add(gtx.Ops)
					paint.PaintOp{}.Add(gtx.Ops)

					drawMadFace(&gtx)
				}

				pointer.InputOp{
					Tag:   mf,
					Kinds: pointer.Press,
				}.Add(gtx.Ops)

				return layout.Dimensions{Size: gtx.Constraints.Max}
			})
		}),
		layout.Flexed(3, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min = image.Point{0, 0}

			alignToRight := gtx.Constraints.Max.X>>2*3 - 12
			gtx.Constraints.Max.X -= alignToRight
			defer op.Offset(image.Pt(alignToRight, 0)).Push(gtx.Ops).Pop()

			return layout.UniformInset(4).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				clr := styles.YELLOWISH_GREEN
				clr.A = 0x6F

				paint.FillShape(gtx.Ops, styles.NIGHT_BLACK, clip.Stroke{Width: 2, Path: clip.Rect{Min: gtx.Constraints.Min, Max: gtx.Constraints.Max}.Path()}.Op())
				defer utils.SetBackgroundColor(&gtx, clr).Pop()

				macro := op.Record(gtx.Ops)
				label := material.Label(styles.MaterialTheme, unit.Sp(16), fmt.Sprintf("%v/%v", mf.MarkedCells, mf.plantedMines))
				label.MaxLines = 1
				label.Font.Weight = font.Medium
				label.Alignment = text.Middle
				labelDim := label.Layout(gtx)
				labelMacro := macro.Stop()

				labelOffset := image.Point{X: (gtx.Constraints.Max.X - labelDim.Size.X) >> 1, Y: (gtx.Constraints.Max.Y - labelDim.Size.Y) >> 1}
				defer op.Offset(labelOffset).Push(gtx.Ops).Pop()

				labelMacro.Add(gtx.Ops)

				return layout.Dimensions{Size: utils.GetSize(&gtx)}
			})
		}),
	)

	headerD.Size = gtx.Constraints.Max

	return headerD
}

func (mf *MineField) bodyComponent(gtx layout.Context) layout.Dimensions {
	rowList := &layout.List{Axis: layout.Vertical, Alignment: layout.Start}

	return rowList.Layout(gtx, len(mf.BtnMatrix), func(gtx layout.Context, rowIndex int) layout.Dimensions {
		list := &layout.List{Axis: layout.Horizontal, Alignment: layout.Start}

		return list.Layout(gtx, len(mf.BtnMatrix[rowIndex]), func(gtx layout.Context, colIndex int) layout.Dimensions {
			btnDimension := mf.BtnMatrix[rowIndex][colIndex].Layout(gtx)

			return btnDimension
		})
	})
}

func (mf *MineField) floodNeutralCells(startingPoint image.Point) {
	mf.State = LOADING
	floodCells := make([]*MineButton, 0, 8)
	floodCells = append(floodCells, &mf.BtnMatrix[startingPoint.Y][startingPoint.X])

	countOfFloodedCells := uint16(0)

	for iterator := 0; iterator < len(floodCells); iterator++ {
		// Venni a jelenlegi elem rejtett környezetét és azokat hozzáadni egy listához
		for i := -1; i <= 1; i++ {
			rowIndex := floodCells[iterator].Pos.Y + i

			// Kilóg felül, vagy alul
			if rowIndex < 0 || rowIndex > int(mf.countOfVerticalCells-1) {
				continue
			}

			// Column loop
			for j := -1; j <= 1; j++ {
				if mf.State != LOADING {
					return
				}

				if j == 0 && i == 0 {
					continue
				}

				colIndex := floodCells[iterator].Pos.X + j

				// Kilóg bal, vagy jobb oldalt
				if colIndex < 0 || colIndex > int(mf.countOfHorizontalCells-1) {
					continue
				}

				// Az adott elem értékét megvizsgálni
				btn := &mf.BtnMatrix[rowIndex][colIndex]
				if !btn.Hidden || btn.Marked {
					continue
				}

				// felfedni és növelni a felfedettek számát
				btn.Hidden = false
				countOfFloodedCells++

				// Ha 0, akkor feldeni, listához hozzáadni
				if btn.Value == 0 {
					floodCells = append(floodCells, btn)
				}

				// Lassú animáció
				if mf.State != LOADING {
					return
				}
				time.Sleep(mf.animationDuration)
			}
		}
	}

	mf.State = RUNNING
	mf.RevealedCells += countOfFloodedCells
}

func (mf *MineField) revealFields() {
	for rowIndex := range mf.BtnMatrix {
		btnList := &mf.BtnMatrix[rowIndex]

		for colIndex := range *btnList {
			btn := &(*btnList)[colIndex]

			btn.Hidden = false
		}
	}
}

func (mf *MineField) generateMines(clickPos image.Point) {
	minePositions := make([]image.Point, 0, mf.MaxNumOfMines)
	// Calculate them

	for i, calcTries := 0, 0; i < cap(minePositions) || calcTries > 5; calcTries++ {
		validPos := true

		// Random num
		minePos := image.Point{int(rand.Int31n(int32(mf.countOfHorizontalCells))), int(rand.Int31n(int32(mf.countOfVerticalCells)))}

		// Check if is it stored already or at 'clickPos'
		if minePos == clickPos {
			continue
		}

		for i := 0; i < len(minePositions); i++ {
			if minePositions[i] == minePos {
				validPos = false
				break
			}
		}

		if validPos {
			minePositions = append(minePositions, minePos)
			i++
			calcTries = 0
		}
	}

	// Clear minefield
	for rowIndex := range mf.BtnMatrix {
		btnList := &mf.BtnMatrix[rowIndex]

		for colIndex := range *btnList {
			btn := &(*btnList)[colIndex]

			btn.State = mf.State
			btn.Value = 0
		}
	}

	// Plant mines
	for i := range minePositions {
		mf.BtnMatrix[minePositions[i].Y][minePositions[i].X].Value = -1
	}
	mf.plantedMines = uint16(len(minePositions))

	// Find mines in the neighborhood, MAAAN!
	for rowIndex := range mf.BtnMatrix {
		btnList := &mf.BtnMatrix[rowIndex]

		for colIndex := range *btnList {
			btn := &(*btnList)[colIndex]

			if btn.Value == -1 {
				continue
			}

			btn.Value = mf.neighboringMines(rowIndex, colIndex)
		}
	}
}

func (mf *MineField) neighboringMines(rowIndexParam, colIndexParam int) int8 {
	sum := int8(0)
	// Row loop
	for i := -1; i <= 1; i++ {
		rowIndex := rowIndexParam + i
		if rowIndex < 0 || rowIndex > int(mf.countOfVerticalCells-1) {
			continue
		}

		// Column loop
		for j := -1; j <= 1; j++ {
			if j == 0 && i == 0 {
				continue
			}

			colIndex := colIndexParam + j
			if colIndex < 0 || colIndex > int(mf.countOfHorizontalCells-1) {
				continue
			}

			if mf.BtnMatrix[rowIndex][colIndex].Value == -1 {
				sum++
			}
		}
	}

	return sum
}

func (mf *MineField) onButtonClick(mb *MineButton) {
	switch mf.State {
	case START:
		fmt.Println("--- Game Start ---")
		mf.State = RUNNING

		// Legenerálni a bombákat
		mf.generateMines(mb.Pos)
		mb.Hidden = false
		mf.RevealedCells++

		if mb.Value == 0 {
			go mf.floodNeutralCells(mb.Pos)
		}
	case RUNNING:
		if mb.Hidden && mb.clickType == pointer.ButtonSecondary {
			mb.Marked = !mb.Marked

			if mb.Marked {
				mf.MarkedCells++
			} else {
				mf.MarkedCells--
			}

			break
		}

		mf.RevealedCells++

		switch mb.Value {
		case -1:
			mf.State = LOSE

			mf.revealFields()
		case 0:
			mb.Hidden = false

			go mf.floodNeutralCells(mb.Pos)
		default:
			mb.Hidden = false
		}
	}

	if mf.State == RUNNING && mf.RevealedCells == mf.countOfHorizontalCells*mf.countOfVerticalCells-mf.plantedMines {
		mf.State = WIN
		mf.MarkedCells = mf.plantedMines
		fmt.Println("--- GG, WIN ---")
	}
}
