package cmd

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jroimartin/gocui"
	"github.com/walletConsole/pancakeswap-console/config"
	"github.com/walletConsole/pancakeswap-console/utils"
	"log"
	"strings"
	"time"
)

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	return nil
}

func layout(g *gocui.Gui) error {
	maxX, _ := g.Size()

	if v, err := g.SetView("help", maxX-43, 0, maxX-1, 7); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "KEYBINDINGS")
		fmt.Fprintln(v, "Mouse Left On Token: Swap")
		fmt.Fprintln(v, "Enter: Confirm Input")
		fmt.Fprintln(v, "^C: Exit")
		fmt.Fprintln(v, "^E: Close Window/Cancel Input")
	}

	if v, err := g.SetView("account", 1, 0, maxX-44, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Account"
		account := Account{
			address: config.CF.FromAddress,
			Refreshable: Refreshable{
				client:  client,
				g:       g,
				view:    "account",
				content: make(chan string),
			},
		}
		account.start()

	}

	stepX := 20
	stepY := 5
	startY := 3
	startX := 1
	//stopY := startY + stepY
	for ind, symbol := range config.CF.PriceToken {
		i := strings.LastIndex(symbol, "-")
		toToken := symbol[i+1:]
		fromToken := symbol[0:i]
		viewName := fmt.Sprintf("%s", symbol)
		title := symbol

		stopX := startX + stepX
		if v, err := g.SetView(viewName, startX, startY, stopX, startY+stepY); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Title = title
			//go updatePrice(g, viewName, symbol)
			tradePairs := TradePairs{
				fromSymbol: fromToken,
				toSymbol:   toToken,
				Refreshable: Refreshable{
					client:  client,
					g:       g,
					view:    viewName,
					content: make(chan string),
				},
			}
			tradePairs.start()

			if err := g.SetKeybinding(viewName, gocui.MouseLeft, gocui.ModNone, buyTokenInputFunc(fromToken, toToken)); err != nil {
				return err
			}
		}
		if (ind+1)%8 == 0 {
			startX = 1
			startY = startY + stepY + 1
		} else {
			startX = stopX + 1
		}
	}
	return nil
}

func buyTokenInputFunc(fromSymbol, toSymbol string) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		_, maxY := g.Size()
		_, cy := v.Cursor()
		l, err := v.Line(cy)
		if err != nil {
			log.Panicln(err)
		}
		//selectText := l[0:strings.Index(l, ":")]
		selectText := strings.TrimSpace(l[0:strings.Index(l, ":")])
		if selectText == "price" {
			return nil
		}
		amountText := strings.TrimSpace(l[strings.Index(l, ":")+1:])
		float := utils.Str2Float(amountText)
		if float <= 0 {
			showMsg(g, "balance too low")
			return nil
		}
		var balance float64
		if strings.ToLower(selectText) == "bnb" {
			ethBalance := utils.EthBalance(config.CF.FromAddress, client)
			balance = utils.Str2Float(ethBalance)
		} else {
			erc20Token := utils.Erc20Token(config.CF.BscToken[strings.ToLower(selectText)], client)
			balanceOf, err := erc20Token.BalanceOf(utils.DefaultCallOpts, common.HexToAddress(config.CF.FromAddress))
			if err != nil {
				showMsg(g, err.Error())
				return nil
			}
			decimals, _ := erc20Token.Decimals(utils.DefaultCallOpts)
			ethBalance := utils.ToDecimal(balanceOf.String(), int(decimals))
			balance = utils.Str2Float(ethBalance.String())

		}
		if float > balance {
			showMsg(g, "balance too low")
			return nil
		}

		if v, err := g.SetView("buybox", 1, maxY-10, 50, maxY-4); err != nil {

			if preView != nil {
				closeBox(g, preView)
			}
			preView = v
			if err != gocui.ErrUnknownView {
				return err
			}
			if _, err := g.SetCurrentView("buybox"); err != nil {
				return err
			}
			v.Title = "Input " + selectText + " Amount  <Enter>"
			v.Editable = true
			if strings.ToLower(selectText) == strings.ToLower(fromSymbol) {
				if err := g.SetKeybinding("buybox", gocui.KeyEnter, gocui.ModNone, buyShowConfirmFunc(fromSymbol, toSymbol, amountText)); err != nil {
					log.Panicln(err)
				}
			} else if strings.ToLower(selectText) == strings.ToLower(toSymbol) {
				if err := g.SetKeybinding("buybox", gocui.KeyEnter, gocui.ModNone, buyShowConfirmFunc(toSymbol, fromSymbol, amountText)); err != nil {
					log.Panicln(err)
				}
			}

			if err := g.SetKeybinding("buybox", gocui.KeyCtrlE, gocui.ModNone, closeBox); err != nil {
				log.Panicln(err)
			}

		}
		return nil
	}
}

func buyShowConfirmFunc(fromSymbol, toSymbol, amount string) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		buffer := v.ViewBuffer()
		//_, cy := v.Cursor()
		buffer = strings.Trim(buffer, "\n ")
		float := utils.Str2Float(buffer)
		if float <= 0 {
			showMsg(g, "balance too low")
			return nil
		}

		closeBox(g, v)
		_, maxY := g.Size()
		if buyboxOut, err := g.SetView("buyShowConfirmFunc", 1, maxY-10, 50, maxY-4); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			if _, err := g.SetCurrentView("buyShowConfirmFunc"); err != nil {
				return err
			}
			if preView != nil {
				closeBox(g, preView)
			}
			preView = buyboxOut

			buyboxOut.Title = "Confirm <Enter>"
			amountsOut, path := GetAmountsOut(fromSymbol, toSymbol, buffer)
			buyboxOut.Write([]byte(fmt.Sprintf("%s: %s\n", fromSymbol, buffer)))
			buyboxOut.Write([]byte(fmt.Sprintf("%s: %f\n", toSymbol, amountsOut)))
			buyboxOut.Write([]byte(fmt.Sprintf("%s", path)))

			if err := g.SetKeybinding("buyShowConfirmFunc", gocui.KeyEnter, gocui.ModNone, buyConfirmFunc(fromSymbol, toSymbol, buffer)); err != nil {
				log.Panicln(err)
			}
			if err := g.SetKeybinding("buyShowConfirmFunc", gocui.KeyCtrlE, gocui.ModNone, closeBox); err != nil {
				log.Panicln(err)
			}

		}

		return nil
	}
}

func buyConfirmFunc(fromSymbol, toSymbol, amount string) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		//_, cy := v.Cursor()

		//g.DeleteKeybindings(v.Name())
		//if err := g.DeleteView(v.Name()); err != nil {
		//	return err
		//}
		closeBox(g, v)
		_, maxY := g.Size()
		if buyboxOut, err := g.SetView("buyConfirmFunc", 1, maxY-10, 100, maxY-4); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			if _, err := g.SetCurrentView("buyConfirmFunc"); err != nil {
				return err
			}
			if preView != nil {
				closeBox(g, preView)
			}
			preView = buyboxOut
			buyboxOut.Title = "Sending"
			buyboxOut.Wrap = true
			hash := Swap(fromSymbol, toSymbol, amount)

			buyboxOut.Write([]byte(fmt.Sprintf("%s\n", hash)))
			if err := g.SetKeybinding("buyConfirmFunc", gocui.KeyCtrlE, gocui.ModNone, closeBox); err != nil {
				log.Panicln(err)
			}
			go func(hash string, v *gocui.View) {
				for {
					select {
					case <-time.After(1 * time.Second):
						var sprintf string
						tx, err := client.TransactionReceipt(context.Background(), common.HexToHash(hash))
						if err != nil {
							sprintf = "pending......"
							g.Update(func(g *gocui.Gui) error {
								v.Write([]byte(fmt.Sprintf("%s\t", sprintf)))
								return nil
							})
						} else {
							status := int(tx.Status)
							if status != 1 {
								sprintf = "failed....... closing window"
							} else {
								sprintf = "success...... closing window"
							}
							g.Update(func(g *gocui.Gui) error {
								v.Write([]byte(fmt.Sprintf("%s\n", sprintf)))
								return nil
							})
							go func() {
								time.Sleep(2 * time.Second)
								closeBox(g, v)
							}()
							return
						}
					}
				}
			}(hash, buyboxOut)
		}

		return nil
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func showMsg(g *gocui.Gui, msg string) error {

	if preView != nil {
		closeBox(g, preView)
	}

	maxX, maxY := g.Size()
	if v, err := g.SetView("msg", 1, maxY-3, maxX-10, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintf(v, "\u001B[31m%s!\u001B[0m\n", msg)
		//if err := g.SetKeybinding(v.Name(), gocui.MouseLeft, gocui.ModNone, closeBox); err != nil {
		//	log.Panicln(err)
		//}
		go func() {
			time.Sleep(1 * time.Second)
			closeBox(g, v)
		}()

	}
	return nil
}
func closeBox(g *gocui.Gui, v *gocui.View) error {
	g.DeleteKeybindings(v.Name())

	if err := g.DeleteView(v.Name()); err != nil {
		return err
	}

	return nil
}
