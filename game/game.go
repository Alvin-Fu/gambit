package game

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	dt "github.com/dylhunn/dragontoothmg"

	"github.com/maaslalani/gambit/board"
	"github.com/maaslalani/gambit/border"
	"github.com/maaslalani/gambit/fen"
	"github.com/maaslalani/gambit/moves"
	"github.com/maaslalani/gambit/pieces"
	"github.com/maaslalani/gambit/position"
	. "github.com/maaslalani/gambit/style"
)

// MoveMsg is a message that controls the board from outside the model.
type MoveMsg struct {
	From string
	To   string
}

// NotifyMsg is a message that gets emitted when the user makes a move.
type NotifyMsg struct {
	From      string
	To        string
	Turn      bool
	Check     bool
	Checkmate bool
}

// Game stores the state of the chess game.
//
// It tracks the board, legal moves, and the selected piece. It also keeps
// track of the subset of legal moves for the currently selected piece
type Game struct {
	board      *dt.Board
	moves      []dt.Move
	pieceMoves []dt.Move
	selected   string
	buffer     string
	flipped    bool
}

// NewGame returns an initial model of the game board.
func NewGame() *Game {
	return NewGameWithPosition(dt.Startpos)
}

// NewGameWithPosition returns an initial model of the game board with the
// specified FEN position.
func NewGameWithPosition(position string) *Game {
	m := &Game{}

	if !fen.IsValid(position) {
		position = dt.Startpos
	}
	board := dt.ParseFen(position)
	m.board = &board
	m.moves = m.board.GenerateLegalMoves()

	return m
}

// Init Initializes the model
func (gm *Game) Init() tea.Cmd {
	return nil
}

// View converts a FEN string into a human readable chess board. All pieces and
// empty squares are arranged in a grid-like pattern. The selected piece is
// highlighted and the legal moves for the selected piece are indicated by a
// dot (.) for empty squares. Pieces that may be captured by the selected piece
// are highlighted.
//
// For example, if the user selects the white pawn on E2 we indicate that they
// can move to E3 and E4 legally.
//
//    ┌───┬───┬───┬───┬───┬───┬───┬───┐
//  8 │ ♖ │ ♘ │ ♗ │ ♕ │ ♔ │ ♗ │ ♘ │ ♖ │
//    ├───┼───┼───┼───┼───┼───┼───┼───┤
//  7 │ ♙ │ ♙ │ ♙ │ ♙ │ ♙ │ ♙ │ ♙ │ ♙ │
//    ├───┼───┼───┼───┼───┼───┼───┼───┤
//  6 │   │   │   │   │   │   │   │   │
//    ├───┼───┼───┼───┼───┼───┼───┼───┤
//  5 │   │   │   │   │   │   │   │   │
//    ├───┼───┼───┼───┼───┼───┼───┼───┤
//  4 │   │   │   │   │ . │   │   │   │
//    ├───┼───┼───┼───┼───┼───┼───┼───┤
//  3 │   │   │   │   │ . │   │   │   │
//    ├───┼───┼───┼───┼───┼───┼───┼───┤
//  2 │ ♟ │ ♟ │ ♟ │ ♟ │ ♟ │ ♟ │ ♟ │ ♟ │
//    ├───┼───┼───┼───┼───┼───┼───┼───┤
//  1 │ ♜ │ ♞ │ ♝ │ ♛ │ ♚ │ ♝ │ ♞ │ ♜ │
//    └───┴───┴───┴───┴───┴───┴───┴───┘
//      A   B   C   D   E   F   G   H
//
func (gm *Game) View() string {
	var s strings.Builder
	s.WriteString(border.Top())

	// Traverse through the rows and columns of the board and print out the
	// pieces and empty squares. Once a piece is selected, highlight the legal
	// moves and pieces that may be captured by the selected piece.
	var rows = fen.Grid(gm.board.ToFen())
	for r := board.FirstRow; r < board.Rows; r++ {
		row := pieces.ToPieces(rows[r])
		rr := board.LastRow - r

		// reverse the row if the board is flipped
		if gm.flipped {
			row = pieces.ToPieces(rows[board.LastRow-r])
			for i, j := 0, len(row)-1; i < j; i, j = i+1, j-1 {
				row[i], row[j] = row[j], row[i]
			}
			rr = r
		}

		s.WriteString(Faint(fmt.Sprintf(" %d ", rr+1)) + border.Vertical)

		for c, piece := range row {
			whiteTurn := gm.board.Wtomove
			display := piece.Display()
			check := gm.board.OurKingInCheck()
			selected := position.ToSquare(r, c, gm.flipped)

			// The user selected the current cell, highlight it so they know it is
			// selected. If it is a check, highlight the king in red.
			if gm.selected == selected {
				display = Cyan(display)
			} else if check && piece.IsKing() {
				if (whiteTurn && piece.IsWhite()) || (!whiteTurn && piece.IsBlack()) {
					display = Red(display)
				}
			}

			// Show all the cells to which the piece may move. If it is an empty cell
			// we present a coloured dot, otherwise color the capturable piece.
			if moves.IsLegal(gm.pieceMoves, selected) {
				if piece.IsEmpty() {
					display = "."
				}
				display = Magenta(display)
			}

			s.WriteString(fmt.Sprintf(" %s %s", display, border.Vertical))
		}
		s.WriteRune('\n')

		if r != board.LastRow {
			s.WriteString(border.Middle())
		}
	}

	s.WriteString(border.Bottom() + Faint(border.BottomLabels(gm.flipped)))
	return s.String()
}

func (gm *Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		if msg.Type != tea.MouseLeft {
			return gm, nil
		}

		// Find the square the user clicked on, this will either be our square
		// square for our piece or the destination square for a move if a piece is
		// already square and that destination square completes a legal move
		square := border.Cell(msg.X, msg.Y, gm.flipped)
		return gm.Select(square)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return gm, tea.Quit
		case "ctrl+f":
			gm.flipped = !gm.flipped
		case "a", "b", "c", "d", "e", "f", "g", "h":
			gm.buffer = msg.String()
		case "1", "2", "3", "4", "5", "6", "7", "8":
			var move string
			if gm.buffer != "" {
				move = gm.buffer + msg.String()
				gm.buffer = ""
			}
			return gm.Select(move)
		case "esc":
			return gm.Deselect()
		}
	case MoveMsg:
		gm.selected = msg.From
		gm.pieceMoves = moves.LegalSelected(gm.moves, gm.selected)
		return gm.Select(msg.To)
	}

	return gm, nil
}

func (gm *Game) Notify(from, to string, turn, check, checkmate bool) tea.Cmd {
	return func() tea.Msg {
		return NotifyMsg{
			From: from, To: to, Turn: turn,
			Check: check, Checkmate: checkmate,
		}
	}
}

func (gm *Game) Deselect() (tea.Model, tea.Cmd) {
	gm.selected = ""
	gm.pieceMoves = []dt.Move{}
	return gm, nil
}

func (gm *Game) Select(square string) (tea.Model, tea.Cmd) {
	// If the user has already selected a piece, check see if the square that
	// the user clicked on is a legal move for that piece. If so, make the move.
	if gm.selected != "" {
		from := gm.selected
		to := square

		for _, move := range gm.pieceMoves {
			if move.String() == from+to {
				var cmds []tea.Cmd
				gm.board.Apply(move)

				// We have applied a new move and the chess board is in a new state.
				// We must generate the new legal moves for the new state.
				gm.moves = gm.board.GenerateLegalMoves()
				check := gm.board.OurKingInCheck()
				checkmate := check && len(gm.moves) == 0

				// We have made a move, so we no longer have a selected piece or
				// legal moves for any selected pieces.
				g, cmd := gm.Deselect()
				cmds = append(cmds, cmd, gm.Notify(from, to, gm.board.Wtomove, check, checkmate))
				return g, tea.Batch(cmds...)
			}
		}

		// The user clicked on a square that wasn't a legal move for the selected
		// piece, so we select the piece that was clicked on instead
		gm.selected = to
	} else {
		gm.selected = square
	}

	// After a mouse click, we must generate the legal moves for the selected
	// piece, if there is a newly selected piece
	gm.pieceMoves = moves.LegalSelected(gm.moves, gm.selected)

	return gm, nil
}

// SetFlipped sets the board to be flipped or not.
func (gm *Game) SetFlipped(flip bool) {
	gm.flipped = flip
}

// Position returns the current FEN position of the board.
func (gm *Game) Position() string {
	return gm.board.ToFen()
}
