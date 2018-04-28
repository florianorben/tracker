package cmd

import (
	"fmt"

	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"tracker/helpers"
	"tracker/tracker"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Get the frames from the server and push the new ones.",
	Long: `Get the frames from the server and push the new ones.

  The URL of the server and the User Token must be defined via the 'tracker
	config' command.

  Example:

  $ tracker config backend.url http://localhost:4242
  $ tracker config backend.token 7e329263e329
  $ tracker upload
  Received 42 frames from the server. Added: 0 Updated: 3
  Pushed 23 frames to the server`,
	Run: upload,
}

func init() {
	RootCmd.AddCommand(uploadCmd)
}

func upload(cmd *cobra.Command, args []string) {
	url := viper.GetString("backend.url")
	token := viper.GetString("backend.token")

	if url == "" || token == "" {
		fmt.Printf("Error: %s\n", helpers.PrintRed("You need to set backend url and token before being able to upload."))
		fmt.Println("        tracker config backend.url http://some.url")
		fmt.Println("        tracker config backend.token mytoken")
		return
	}

	frames := tracker.GetFrames()
	b, err := json.Marshal(frames)

	if err != nil {
		fmt.Printf("Error: %s. %s\n", helpers.PrintRed("Unabled to marshal frames"), err.Error())
		return
	}

	err = send(url, token, b)
	if err != nil {
		fmt.Printf("Error: %s\n", helpers.PrintRed(err.Error()))
		return
	}

	fmt.Println("Sync successul")
}

func send(url, token string, data []byte) error {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("X-Token", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return errors.New(resp.Status)
	}

	var body tracker.Frames
	if err := json.Unmarshal(b, &body); err != nil {
		return err
	}

	frames := tracker.GetFrames()
	cnt := len(frames)
	added, updated := frames.Merge(body)

	fmt.Printf(
		"Received %s frames from server. Added: %s Updated: %s\n",
		helpers.PrintBold(fmt.Sprintf("%d", len(body))),
		helpers.PrintBold(fmt.Sprintf("%d", added)),
		helpers.PrintBold(fmt.Sprintf("%d", updated)),
	)
	fmt.Printf("Pushed %s frames to server.\n", helpers.PrintBold(fmt.Sprintf("%d", cnt)))

	return nil
}
