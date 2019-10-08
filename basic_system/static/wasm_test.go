package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"io/ioutil"
	"net/http"
	"syscall/js"
	"time"
)

type PatchBody struct {
	Patch string `json:"patch"`
	Text  string `json:"text"`
	Err   string `json:"err"`
}

var (
	docStr string
	client *http.Client
	dmp    *diffmatchpatch.DiffMatchPatch
)

func init() {
	docStr = ""
	client = &http.Client{
		Timeout: time.Second * 10,
	}
	dmp = diffmatchpatch.New()
}

func sendPatch(text string) {
	patches := dmp.PatchMake(docStr, text)
	patchStr := dmp.PatchToText(patches)
	patchBody := PatchBody{Patch: patchStr}
	patchJSON, err := json.Marshal(patchBody)
	if err != nil {
		fmt.Println("Error marshalling request body:", err)
		return
	}

	reqBody := bytes.NewReader(patchJSON)

	resp, err := http.Post("/patch", "application/json", reqBody)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}
		patchResp := PatchBody{}
		err = json.Unmarshal(respBody, &patchResp)
		if err != nil {
			fmt.Println("Error parsing response:", err)
			return
		}
		fmt.Printf("Patch was unsuccessful:", patchResp.Err)
	}
}

func getPatch(text string) {
	patchBody := PatchBody{Text: text}
	patchJSON, err := json.Marshal(patchBody)
	if err != nil {
		fmt.Println(err)
		return
	}

	reqBody := bytes.NewReader(patchJSON)

	req, err := http.NewRequest("GET", "/patch", reqBody)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	respData := PatchBody{}

	err = json.Unmarshal(respBody, &respData)
	if err != nil {
		fmt.Println(err)
		return
	}

	patches, err := dmp.PatchFromText(respData.Patch)
	if err != nil {
		fmt.Println(err)
		return
	}

	newText, _ := dmp.PatchApply(patches, text)

	doc := js.Global().Get("document").Call("querySelector", "#doc")
	doc.Set("value", newText)
	docStr = newText
}

func inputHandler(this js.Value, args []js.Value) interface{} {
	go func() {
		event := args[0]
		if inputType := event.Get("inputType"); inputType.String() == "insertText" ||
			inputType.String() == "insertLineBreak" {
			doc := js.Global().Get("document").Call("querySelector", "#doc")
			sendPatch(doc.Get("value").String())
		}
	}()

	return nil
}

func main() {
	doc := js.Global().Get("document").Call("querySelector", "#doc")
	doc.Call(
		"addEventListener",
		"input",
		js.FuncOf(inputHandler),
	)
	c := make(chan bool)
	<-c
}
