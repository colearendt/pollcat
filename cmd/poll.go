package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"github.com/colearendt/pollcat/internal/model"
	"github.com/colearendt/pollcat/internal/poller"
	"github.com/colearendt/pollcat/internal/store"
	"github.com/colearendt/pollcat/internal/ui"
)

var (
	tcpFlags    []string
	udpFlags    []string
	icmpFlags   []string
	dnsFlags    []string
	intervalStr string
	durationStr string
	timeoutStr  string
	outFile     string
	noTUI       bool
)

func init() {
	pollCmd.Flags().StringArrayVar(&tcpFlags, "tcp", nil, "TCP target(s) to poll, e.g. 1.2.3.4:80")
	pollCmd.Flags().StringArrayVar(&udpFlags, "udp", nil, "UDP target(s) to poll, e.g. 1.2.3.4:53")
	pollCmd.Flags().StringArrayVar(&icmpFlags, "icmp", nil, "ICMP target(s) to poll, e.g. 1.2.3.4")
	pollCmd.Flags().StringArrayVar(&dnsFlags, "dns", nil, "DNS target(s) to poll, e.g. example.com")
	pollCmd.Flags().StringVarP(&intervalStr, "interval", "i", "1s", "Polling interval")
	pollCmd.Flags().StringVarP(&durationStr, "duration", "d", "0", "Total duration to poll (0 = until interrupted)")
	pollCmd.Flags().StringVar(&timeoutStr, "timeout", "5s", "Timeout for each individual poll")
	pollCmd.Flags().StringVarP(&outFile, "out", "o", "", "Output file for raw results (JSON)")
	pollCmd.Flags().BoolVar(&noTUI, "no-tui", false, "Disable interactive TUI and print to stdout")
	rootCmd.AddCommand(pollCmd)
}

var pollCmd = &cobra.Command{
	Use:   "poll",
	Short: "Start polling TCP, UDP, ICMP and/or DNS targets.",
	RunE: func(cmd *cobra.Command, args []string) error {
		interval, err := time.ParseDuration(intervalStr)
		if err != nil {
			return fmt.Errorf("invalid interval: %w", err)
		}
		var duration time.Duration
		if durationStr != "0" {
			duration, err = time.ParseDuration(durationStr)
			if err != nil {
				return fmt.Errorf("invalid duration: %w", err)
			}
		}
		timeout, err := time.ParseDuration(timeoutStr)
		if err != nil {
			return fmt.Errorf("invalid timeout: %w", err)
		}

		if len(tcpFlags) == 0 && len(udpFlags) == 0 && len(icmpFlags) == 0 && len(dnsFlags) == 0 {
			return fmt.Errorf("at least one --tcp, --udp, --icmp or --dns target is required")
		}

		if !noTUI && !isatty.IsTerminal(os.Stdout.Fd()) {
			noTUI = true
		}

		var targets []model.Target
		for _, t := range tcpFlags {
			targets = append(targets, model.Target{Type: model.PollTypeTCP, Address: t, Interval: interval})
		}
		for _, t := range udpFlags {
			targets = append(targets, model.Target{Type: model.PollTypeUDP, Address: t, Interval: interval})
		}
		for _, t := range icmpFlags {
			targets = append(targets, model.Target{Type: model.PollTypeICMP, Address: t, Interval: interval})
		}
		for _, t := range dnsFlags {
			targets = append(targets, model.Target{Type: model.PollTypeDNS, Address: t, Interval: interval})
		}

		resultsCh := make(chan model.Result, len(targets)*2)
		st := store.New()

		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()
		if duration > 0 {
			var cancelTimeout context.CancelFunc
			ctx, cancelTimeout = context.WithTimeout(ctx, duration)
			defer cancelTimeout()
		}

		for _, t := range targets {
			t := t // capture range variable
			go func() {
				ticker := time.NewTicker(t.Interval)
				defer ticker.Stop()

				var p poller.Poller
				switch t.Type {
				case model.PollTypeTCP:
					p = &poller.TCPPoller{Dialer: &net.Dialer{Timeout: timeout}}
				case model.PollTypeUDP:
					p = &poller.UDPPoller{Dialer: &net.Dialer{Timeout: timeout}}
				case model.PollTypeICMP:
					p = &poller.ICMPPoller{Pinger: &poller.DefaultPinger{Timeout: timeout}}
				case model.PollTypeDNS:
					p = &poller.DNSPoller{Resolver: net.DefaultResolver, Timeout: timeout}
				}

				for {
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
						res := p.Poll(ctx, t)
						st.Append(res)
						resultsCh <- res
					}
				}
			}()
		}

		if noTUI {
			fmt.Println("Polling started...")
			for {
				select {
				case <-ctx.Done():
					goto done
				case res := <-resultsCh:
					status := "OK"
					if !res.Success {
						status = "FAIL"
					}
					fmt.Printf("[%s] %-6s %-30s %s %12s\n",
						res.Timestamp.Format("15:04:05"),
						res.Type,
						res.Target,
						status,
						res.Latency.Round(time.Microsecond),
					)
				}
			}
		} else {
			modelUI := ui.New(resultsCh)
			p := tea.NewProgram(modelUI, tea.WithAltScreen())

			go func() {
				<-ctx.Done()
				p.Quit()
			}()

			if _, err := p.Run(); err != nil {
				return fmt.Errorf("TUI error: %w", err)
			}
		}

	done:
		// Stop pollers and drain any remaining results.
		cancel()
		go func() {
			for range resultsCh {
			}
		}()

		if outFile != "" {
			f, err := os.Create(outFile)
			if err != nil {
				return fmt.Errorf("create output file: %w", err)
			}
			defer f.Close()
			enc := json.NewEncoder(f)
			enc.SetIndent("", "  ")
			if err := enc.Encode(st.Results()); err != nil {
				return fmt.Errorf("encode results: %w", err)
			}
		}

		return nil
	},
}
