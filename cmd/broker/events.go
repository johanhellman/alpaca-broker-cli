package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	client "github.com/johanhellman/alpaca-broker-cli/pkg/brokerclient"
	"github.com/johanhellman/alpaca-broker-cli/pkg/brokerclient/api"
	"github.com/spf13/cobra"
)

var (
	eventsSince   string
	eventsUntil   string
	eventsSinceID int
	eventsUntilID int
)

var eventsCmd = &cobra.Command{
	Use:   "events",
	Short: "Stream SSE events from the broker API",
}

var eventsAccountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Stream account status events",
	RunE: func(cmd *cobra.Command, args []string) error {
		return streamEvents("accounts")
	},
}

var eventsJournalsCmd = &cobra.Command{
	Use:   "journals",
	Short: "Stream journal status events",
	RunE: func(cmd *cobra.Command, args []string) error {
		return streamEvents("journals")
	},
}

var eventsTradesCmd = &cobra.Command{
	Use:   "trades",
	Short: "Stream trade updates events",
	RunE: func(cmd *cobra.Command, args []string) error {
		return streamEvents("trades")
	},
}

func streamEvents(eventType string) error {
	c, err := api.NewClient()
	if err != nil {
		return err
	}

	// Wait for interrupt signal to gracefully shutdown the stream
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println("\nClosing stream...")
		cancel()
	}()

	var sinceTime, untilTime *time.Time
	if eventsSince != "" {
		t, err := time.Parse(time.RFC3339, eventsSince)
		if err != nil {
			return fmt.Errorf("invalid since format: %w", err)
		}
		sinceTime = &t
	}
	if eventsUntil != "" {
		t, err := time.Parse(time.RFC3339, eventsUntil)
		if err != nil {
			return fmt.Errorf("invalid until format: %w", err)
		}
		untilTime = &t
	}

	var sinceIDPtr, untilIDPtr *int
	if eventsSinceID != 0 {
		sinceIDPtr = &eventsSinceID
	}
	if eventsUntilID != 0 {
		untilIDPtr = &eventsUntilID
	}

	switch eventType {
	case "accounts":
		params := &client.GetEventsAccountsStatusParams{
			Since:   sinceTime,
			Until:   untilTime,
			SinceId: sinceIDPtr,
			UntilId: untilIDPtr,
		}
		resp, err := c.GetEventsAccountsStatus(ctx, params)
		if err != nil {
			return fmt.Errorf("failed to connect to stream: %w", err)
		}
		defer func() { _ = resp.Body.Close() }() //nolint:errcheck
		_, err = io.Copy(os.Stdout, resp.Body)
		return err

	case "journals":
		params := &client.GetEventsJournalsStatusParams{
			Since:   sinceTime,
			Until:   untilTime,
			SinceId: sinceIDPtr,
			UntilId: untilIDPtr,
		}
		resp, err := c.GetEventsJournalsStatus(ctx, params)
		if err != nil {
			return fmt.Errorf("failed to connect to stream: %w", err)
		}
		defer func() { _ = resp.Body.Close() }() //nolint:errcheck
		_, err = io.Copy(os.Stdout, resp.Body)
		return err

	case "trades":
		params := &client.GetEventsTradesParams{
			Since:   sinceTime,
			Until:   untilTime,
			SinceId: sinceIDPtr,
			UntilId: untilIDPtr,
		}
		resp, err := c.GetEventsTrades(ctx, params)
		if err != nil {
			return fmt.Errorf("failed to connect to stream: %w", err)
		}
		defer func() { _ = resp.Body.Close() }() //nolint:errcheck
		_, err = io.Copy(os.Stdout, resp.Body)
		return err

	default:
		return fmt.Errorf("unknown event type: %s", eventType)
	}
}

func init() {
	rootCmd.AddCommand(eventsCmd)

	for _, subCmd := range []*cobra.Command{eventsAccountsCmd, eventsJournalsCmd, eventsTradesCmd} {
		subCmd.Flags().StringVar(&eventsSince, "since", "", "Filter events after this time (RFC3339)")
		subCmd.Flags().StringVar(&eventsUntil, "until", "", "Filter events before this time (RFC3339)")
		subCmd.Flags().IntVar(&eventsSinceID, "since-id", 0, "Filter events after this ID")
		subCmd.Flags().IntVar(&eventsUntilID, "until-id", 0, "Filter events before this ID")
		eventsCmd.AddCommand(subCmd)
	}
}
