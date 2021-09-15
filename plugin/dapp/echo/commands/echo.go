package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/33cn/chain33/rpc/jsonclient"
	echotypes "github.com/33cn/plugin/plugin/dapp/echo/types/echo"
	"github.com/spf13/cobra"
)

func EchoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "echo",
		Short: "echo commandline interface",
		Args:  cobra.MinimumNArgs(1),
	}
	cmd.AddCommand(
		QueryCmd(), 
	)
	return cmd
}
func QueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query",
		Short: "query message history",
		Run:   queryMesage,
	}
	addPingPangFlags(cmd)
	return cmd
}

func addPingPangFlags(cmd *cobra.Command) {
	cmd.Flags().Uint32P("type", "t", 1, "message type, 1:ping  2:pang")
	//cmd.MarkFlagRequired("type")

	cmd.Flags().StringP("message", "m", "", "message content")
	cmd.MarkFlagRequired("message")
}

func queryMesage(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	echoType, _ := cmd.Flags().GetUint32("type")
	msg, _ := cmd.Flags().GetString("message")
	client, err := jsonclient.NewJSONClient(rpcLaddr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	var action = &echotypes.Query{Msg: msg}
	if echoType != 1 {
		fmt.Fprintln(os.Stderr, "not support")
		return
	}

	var result echotypes.QueryResult
	err = client.Call("echo.QueryPing", action, &result)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	data, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println(string(data))
}
