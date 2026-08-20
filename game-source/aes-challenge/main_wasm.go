//go:build js && wasm

package main

import (
	"fmt"
	"html"
	"math/rand"
	"strconv"
	"strings"
	"syscall/js"
	"time"
)

type game struct {
	root        js.Value
	deck        []Round
	roundIndex  int
	correct     bool
	wrong       bool
	wrongChoice int
	attempts    int
	click       js.Func
}

func newGame(root js.Value) *game {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	g := &game{root: root, deck: newDeck(random), wrongChoice: -1}
	g.click = js.FuncOf(g.handleClick)
	root.Call("addEventListener", "click", g.click)
	return g
}

func (g *game) currentRound() Round {
	return g.deck[g.roundIndex]
}

func (g *game) resetRound() {
	g.correct = false
	g.wrong = false
	g.wrongChoice = -1
}

func (g *game) restart() {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	g.deck = newDeck(random)
	g.roundIndex = 0
	g.attempts = 0
	g.resetRound()
	g.render()
}

func (g *game) handleClick(_ js.Value, arguments []js.Value) any {
	if len(arguments) == 0 {
		return nil
	}
	target := arguments[0].Get("target")
	if target.IsNull() || target.IsUndefined() {
		return nil
	}
	control := target.Call("closest", "[data-action]")
	if control.IsNull() || control.IsUndefined() {
		return nil
	}

	switch control.Get("dataset").Get("action").String() {
	case "choose":
		if g.correct {
			return nil
		}
		choice, err := strconv.Atoi(control.Get("dataset").Get("choice").String())
		if err != nil {
			return nil
		}
		g.attempts++
		if choice == g.currentRound().KeyPosition {
			g.correct = true
			g.wrong = false
			g.wrongChoice = -1
		} else {
			g.wrong = true
			g.wrongChoice = choice
		}
		g.render()
	case "retry":
		g.wrong = false
		g.wrongChoice = -1
		g.render()
	case "next":
		if !g.correct {
			return nil
		}
		g.roundIndex++
		if g.roundIndex < len(g.deck) {
			g.resetRound()
		}
		g.render()
	case "restart":
		g.restart()
	}
	return nil
}

func groupHex(value string) string {
	groups := make([]string, 0, (len(value)+3)/4)
	for start := 0; start < len(value); start += 4 {
		end := start + 4
		if end > len(value) {
			end = len(value)
		}
		groups = append(groups, value[start:end])
	}
	return strings.Join(groups, " ")
}

func letter(position int) string {
	if position == 0 {
		return "A"
	}
	return "B"
}

func (g *game) card(position int, value string) string {
	classes := []string{"choice-card"}
	badge := "Choose as key"
	accessibleStatus := "Choose as the key string."
	disabled := ""

	if g.correct {
		disabled = " disabled"
		if position == g.currentRound().KeyPosition {
			classes = append(classes, "is-correct")
			badge = "Key string"
			accessibleStatus = "Key string."
		} else {
			classes = append(classes, "is-ciphertext")
			badge = "Ciphertext"
			accessibleStatus = "Ciphertext."
		}
	} else if g.wrong && position == g.wrongChoice {
		classes = append(classes, "is-wrong")
		badge = "Not the key"
		accessibleStatus = "Not the key."
	}

	label := letter(position)
	return fmt.Sprintf(`
      <button class="%s" type="button" data-action="choose" data-choice="%d" aria-label="String %s, %s. %s"%s>
        <span class="choice-heading">
          <span class="choice-label">String %s</span>
          <span class="choice-badge">%s</span>
        </span>
        <code>%s</code>
      </button>`,
		strings.Join(classes, " "), position, label, value, accessibleStatus, disabled, label, badge, groupHex(value))
}

func (g *game) feedback() string {
	if g.correct {
		challenge := challenges[g.currentRound().ChallengeIndex]
		key, err := challenge.actualKey()
		if err != nil {
			key = "Unavailable"
		}
		message, err := challenge.decrypt()
		if err != nil {
			message = "Unable to decrypt"
		}
		start := challenge.KeyOffset + 1
		end := challenge.KeyOffset + aesKeyLength
		nextLabel := "Next round"
		if g.roundIndex == len(g.deck)-1 {
			nextLabel = "See results"
		}
		return fmt.Sprintf(`
      <section class="feedback is-success" aria-live="polite">
        <div>
          <p class="feedback-kicker">Correct</p>
          <h2>String %s contains the key.</h2>
		  <p>The other string is the ciphertext. Using characters %d–%d (counting from 1) as the AES-128 key decrypts its message.</p>
        </div>
        <div class="reveal" aria-label="Round result">
          <span><small>Hidden key</small><code>%s</code></span>
          <span><small>Decoded message</small><strong>%s</strong></span>
        </div>
        <button class="primary-action" type="button" data-action="next">%s <span aria-hidden="true">→</span></button>
      </section>`, letter(g.currentRound().KeyPosition), start, end, html.EscapeString(key), html.EscapeString(message), nextLabel)
	}

	if g.wrong {
		return fmt.Sprintf(`
      <section class="feedback is-error" aria-live="assertive">
        <div>
          <p class="feedback-kicker">Try again</p>
          <h2>String %s is the ciphertext.</h2>
          <p>Reset this pair, then choose the other string as the key container.</p>
        </div>
        <button class="secondary-action" type="button" data-action="retry">Retry this round</button>
      </section>`, letter(g.wrongChoice))
	}

	return `
      <section class="hint" aria-live="polite">
        <span class="hint-icon" aria-hidden="true">i</span>
        <p><strong>Same appearance, different jobs.</strong> One 32-character string contains a consecutive 16-character AES key. The other represents one 16-byte ciphertext block in hexadecimal.</p>
      </section>`
}

func (g *game) results() string {
	return fmt.Sprintf(`
    <main class="results-card">
      <p class="eyebrow">Challenge complete</p>
      <h1>You classified all eight pairs.</h1>
      <p>Every assignment string was used once. The hidden key windows decrypted eight fruit messages, just as in the original classroom challenge.</p>
      <dl>
        <div><dt>Pairs solved</dt><dd>8 / 8</dd></div>
        <div><dt>Total choices</dt><dd>%d</dd></div>
      </dl>
      <button class="primary-action" type="button" data-action="restart">Play a new order <span aria-hidden="true">↻</span></button>
    </main>`, g.attempts)
}

func (g *game) render() {
	if g.roundIndex >= len(g.deck) {
		g.root.Set("innerHTML", g.results())
		return
	}

	round := g.currentRound()
	items := round.strings()
	status := "Choose the key string"
	if g.wrong {
		status = "Retry this pair"
	} else if g.correct {
		status = "Pair classified"
	}
	solved := g.roundIndex
	if g.correct {
		solved++
	}

	markup := fmt.Sprintf(`
    <header class="lab-header">
      <div>
        <p class="eyebrow">Project I · AES Challenge</p>
        <h1>Key or ciphertext?</h1>
        <p class="lab-intro">The assignment mixed eight key strings with eight encrypted strings. Classify one matched pair at a time.</p>
      </div>
      <div class="progress" aria-label="Round %d of %d">
        <span>Round</span>
        <strong>%02d <small>/ %02d</small></strong>
      </div>
    </header>
    <main>
      <div class="round-heading">
        <div>
          <p class="round-label">Your task</p>
          <h2>%s</h2>
        </div>
        <p class="score"><span>%d</span> pairs solved</p>
      </div>
      <div class="choices" aria-label="AES string choices">
        %s
        <span class="versus" aria-hidden="true">or</span>
        %s
      </div>
      %s
    </main>`,
		g.roundIndex+1, len(g.deck), g.roundIndex+1, len(g.deck), status, solved, g.card(0, items[0]), g.card(1, items[1]), g.feedback())

	g.root.Set("innerHTML", markup)
}

func main() {
	document := js.Global().Get("document")
	root := document.Call("getElementById", "aes-challenge")
	if root.IsNull() || root.IsUndefined() {
		return
	}

	game := newGame(root)
	game.render()
	document.Get("body").Get("classList").Call("add", "is-ready")
	js.Global().Set("__aesChallengeReady", true)
	select {}
}
