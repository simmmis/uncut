package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

const (
	version = "0.1.0"
)

var (
	errMissingAPIKey   = errors.New("missing api key")
	errMissingEndpoint = errors.New("missing endpoint")
	stdinReader        = bufio.NewReader(os.Stdin)
)

type configFile struct {
	APIKey    string `json:"api_key"`
	Endpoint  string `json:"endpoint"`
	CreatedAt string `json:"created_at,omitempty"`
}

type credentials struct {
	APIKey         string
	Endpoint       string
	APIKeySource   string
	EndpointSource string
}

type client struct {
	apiKey   string
	endpoint string
	http     *http.Client
}

type envelope struct {
	Data    json.RawMessage `json:"data"`
	Error   *apiErrorBody   `json:"error"`
	Page    int             `json:"page,omitempty"`
	PerPage int             `json:"per_page,omitempty"`
	HasMore bool            `json:"has_more,omitempty"`
	Source  string          `json:"source,omitempty"`
}

type apiErrorBody struct {
	Code    string          `json:"code"`
	Message string          `json:"message"`
	Details json.RawMessage `json:"details,omitempty"`
}

type apiError struct {
	Status  int
	Code    string
	Message string
	Details json.RawMessage
	Body    string
}

func (e *apiError) Error() string {
	msg := e.Message
	if msg == "" {
		msg = strings.TrimSpace(e.Body)
	}
	if msg == "" {
		msg = http.StatusText(e.Status)
	}
	if e.Code != "" {
		return fmt.Sprintf("api error %d %s: %s", e.Status, e.Code, msg)
	}
	return fmt.Sprintf("api error %d: %s", e.Status, msg)
}

type card struct {
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	Mask           string      `json:"mask"`
	ExpireMonth    string      `json:"expire_month"`
	ExpireYear     string      `json:"expire_year"`
	Currency       string      `json:"currency"`
	Balance        json.Number `json:"balance"`
	Status         string      `json:"status"`
	Phone3DS       string      `json:"phone_3ds"`
	CreatedAt      string      `json:"created_at"`
	CardNumber     string      `json:"card_number,omitempty"`
	ExpirationDate string      `json:"expiration_date,omitempty"`
	CVV            string      `json:"cvv,omitempty"`
}

type walletBalance struct {
	Balance  json.Number `json:"balance"`
	Currency string      `json:"currency"`
}

type depositAddress struct {
	Chain   string `json:"chain"`
	Token   string `json:"token"`
	Address string `json:"address"`
}

type cardBIN struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	WalletSupport bool     `json:"wallet_support"`
	Currencies    []string `json:"currencies"`
	Rates         binRates `json:"rates"`
}

type binRates struct {
	CardIssueFee             json.Number `json:"card_issue_fee"`
	CardTopupPercent         json.Number `json:"card_topup_percent"`
	CardTopupFix             json.Number `json:"card_topup_fix"`
	CardAuthorizationPercent json.Number `json:"card_authorization_percent"`
	CardAuthorizationFix     json.Number `json:"card_authorization_fix"`
	CardWithdrawPercent      json.Number `json:"card_withdraw_percent"`
	CardWithdrawFix          json.Number `json:"card_withdraw_fix"`
}

type queuedOperation struct {
	OperationID string `json:"operation_id"`
	Status      string `json:"status"`
	Message     string `json:"message"`
}

type operation struct {
	OperationID  string      `json:"operation_id"`
	Type         string      `json:"type"`
	Status       string      `json:"status"`
	Amount       json.Number `json:"amount"`
	CreatedAt    string      `json:"created_at"`
	UpdatedAt    string      `json:"updated_at"`
	CardID       string      `json:"card_id"`
	ErrorMessage string      `json:"error_message"`
}

type cardTransaction struct {
	Type             string      `json:"type"`
	Time             string      `json:"time"`
	MerchantName     string      `json:"merchant_name"`
	MerchantCountry  string      `json:"merchant_country"`
	OriginalAmount   json.Number `json:"original_amount"`
	OriginalCurrency string      `json:"original_currency"`
	PreAuthAmount    json.Number `json:"pre_auth_amount"`
	PostedAmount     json.Number `json:"posted_amount"`
}

type walletTransaction struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Status    string      `json:"status"`
	Amount    json.Number `json:"amount"`
	Fee       json.Number `json:"fee"`
	CardID    *string     `json:"card_id"`
	Comment   *string     `json:"comment"`
	CreatedAt string      `json:"created_at"`
}

type deleteResult struct {
	ReturnedAmount   json.Number `json:"returned_amount"`
	ReturnedCurrency string      `json:"returned_currency"`
	WalletBalance    json.Number `json:"wallet_balance"`
}

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 {
		printHelp()
		return 0
	}
	if len(args) > 1 && (args[1] == "--help" || args[1] == "-h") {
		return printCommandHelp(args[0])
	}

	switch args[0] {
	case "-h", "--help":
		printHelp()
		return 0
	case "help":
		if len(args) > 1 {
			return printCommandHelp(args[1])
		}
		printHelp()
		return 0
	case "-v", "--version", "version":
		fmt.Printf("uncut %s\n", version)
		return 0
	case "login":
		return cmdLogin(args[1:])
	case "logout":
		return cmdLogout(args[1:])
	case "config":
		return cmdConfig(args[1:])
	case "balance":
		return cmdBalance(args[1:])
	case "deposit":
		return cmdDeposit(args[1:])
	case "wallet":
		return cmdWallet(args[1:])
	case "bins":
		return cmdBins(args[1:])
	case "cards":
		return cmdCards(args[1:])
	case "card":
		return cmdCard(args[1:])
	case "new":
		return cmdNew(args[1:])
	case "topup":
		return cmdTopup(args[1:])
	case "withdraw":
		return cmdWithdraw(args[1:])
	case "transactions":
		return cmdTransactions(args[1:])
	case "rename":
		return cmdRename(args[1:])
	case "phone":
		return cmdPhone(args[1:])
	case "freeze":
		return cmdFreeze(args[1:])
	case "unfreeze":
		return cmdUnfreeze(args[1:])
	case "delete":
		return cmdDelete(args[1:])
	case "operation":
		return cmdOperation(args[1:])
	case "wait":
		return cmdWait(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", args[0])
		printHelp()
		return 2
	}
}

func printHelp() {
	fmt.Println("uncut - console client for Uncutt Cards API")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  uncut login")
	fmt.Println("  uncut balance")
	fmt.Println("  uncut deposit")
	fmt.Println("  uncut wallet [--page 2]")
	fmt.Println("  uncut bins")
	fmt.Println("  uncut cards [--reveal] [--json]")
	fmt.Println("  uncut card <card_id> [--reveal] [--json]")
	fmt.Println("  uncut new <bin_id> [--name <name>] [--currency USD] [--topup 0]")
	fmt.Println("  uncut topup <card_id> --amount 50 [--wait]")
	fmt.Println("  uncut withdraw <card_id> --amount 20 [--wait]")
	fmt.Println("  uncut transactions <card_id> [--page 2]")
	fmt.Println("  uncut rename <card_id> <new_name>")
	fmt.Println("  uncut phone <card_id> --phone +10000000000")
	fmt.Println("  uncut freeze <card_id>")
	fmt.Println("  uncut unfreeze <card_id>")
	fmt.Println("  uncut delete <card_id> --yes")
	fmt.Println("  uncut operation <operation_id>")
	fmt.Println("  uncut wait <operation_id>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  login      save API key and endpoint locally")
	fmt.Println("  logout     remove saved API key")
	fmt.Println("  config     show config status without exposing secrets")
	fmt.Println("  balance    show wallet balance")
	fmt.Println("  deposit    show USDT deposit addresses")
	fmt.Println("  wallet     show wallet transaction history")
	fmt.Println("  bins       list available card BINs and fees")
	fmt.Println("  cards      list cards with safe masked details")
	fmt.Println("  card       show one card; add --reveal for PAN/CVV")
	fmt.Println("  new        create a new card")
	fmt.Println("  topup      move wallet funds to a card")
	fmt.Println("  withdraw   move card funds back to wallet")
	fmt.Println("  transactions show card transaction history")
	fmt.Println("  rename     rename a card")
	fmt.Println("  phone      update 3DS phone")
	fmt.Println("  freeze     freeze a card")
	fmt.Println("  unfreeze   unfreeze a card")
	fmt.Println("  delete     close a card and refund balance")
	fmt.Println("  operation  show async operation status")
	fmt.Println("  wait       wait for async operation completion")
	fmt.Println()
	fmt.Println("Use `uncut help <command>` for examples.")
}

func printCommandHelp(command string) int {
	switch command {
	case "login":
		fmt.Println("Usage: uncut login")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  uncut login")
		fmt.Println()
		fmt.Println("The command asks for both API key and API endpoint.")
	case "logout":
		fmt.Println("Usage: uncut logout")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  uncut logout")
	case "config":
		fmt.Println("Usage: uncut config")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  uncut config")
	case "balance":
		fmt.Println("Usage: uncut balance [--json]")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  uncut balance")
		fmt.Println("  uncut balance --json")
	case "deposit":
		fmt.Println("Usage: uncut deposit [--json]")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  uncut deposit")
		fmt.Println("  uncut deposit --json")
	case "wallet":
		fmt.Println("Usage: uncut wallet [--page <n>|--all] [--json]")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  uncut wallet")
		fmt.Println("  uncut wallet --page 2")
		fmt.Println("  uncut wallet --all")
	case "bins":
		fmt.Println("Usage: uncut bins [--json]")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  uncut bins")
		fmt.Println("  uncut bins --json")
	case "new":
		fmt.Println("Usage: uncut new <bin_id> [--name <name>] [--currency USD] [--topup 0] [--wait]")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  uncut new")
		fmt.Println("  uncut new <bin_id> --wait")
		fmt.Println("  uncut new <bin_id> --name 'Facebook Ads' --topup 25 --wait")
		printHelpBINExamples()
	case "cards":
		fmt.Println("Usage: uncut cards [--reveal|--full] [--json]")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  uncut cards")
		fmt.Println("  uncut cards --reveal")
		fmt.Println("  uncut cards --json")
	case "card":
		fmt.Println("Usage: uncut card <card_id> [--reveal|--full] [--json]")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  uncut card <card_id>")
		fmt.Println("  uncut card <card_id> --reveal")
		printHelpCardExamples("card")
	case "topup":
		fmt.Println("Usage: uncut topup <card_id> --amount <amount> [--wait] [--json]")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  uncut topup <card_id> --amount 60 --wait")
		fmt.Println("  uncut topup <card_id> 60 --wait")
		printHelpCardExamples("topup")
	case "withdraw":
		fmt.Println("Usage: uncut withdraw <card_id> --amount <amount> [--wait] [--json]")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  uncut withdraw <card_id> --amount 20 --wait")
		fmt.Println("  uncut withdraw <card_id> 20 --wait")
		printHelpCardExamples("withdraw")
	case "transactions":
		fmt.Println("Usage: uncut transactions <card_id> [--page <n>|--all] [--json]")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  uncut transactions <card_id>")
		fmt.Println("  uncut transactions <card_id> --page 2")
		printHelpCardExamples("transactions")
	case "rename":
		fmt.Println("Usage: uncut rename <card_id> <new_name> [--json]")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  uncut rename <card_id> 'Facebook Ads'")
		printHelpCardExamples("rename")
	case "phone":
		fmt.Println("Usage: uncut phone <card_id> --phone <e164_phone> [--json]")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  uncut phone <card_id> --phone +10000000000")
		printHelpCardExamples("phone")
	case "freeze":
		fmt.Println("Usage: uncut freeze <card_id> [--json]")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  uncut freeze <card_id>")
		printHelpCardExamples("freeze")
	case "unfreeze":
		fmt.Println("Usage: uncut unfreeze <card_id> [--json]")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  uncut unfreeze <card_id>")
		printHelpCardExamples("unfreeze")
	case "delete":
		fmt.Println("Usage: uncut delete <card_id> [--yes] [--json]")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  uncut delete <card_id>")
		fmt.Println("  uncut delete <card_id> --yes")
		printHelpCardExamples("delete")
	case "operation":
		fmt.Println("Usage: uncut operation <operation_id> [--json]")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  uncut operation <operation_id>")
	case "wait":
		fmt.Println("Usage: uncut wait <operation_id> [--interval <seconds>] [--timeout <seconds>] [--json]")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  uncut wait <operation_id>")
		fmt.Println("  uncut wait <operation_id> --interval 5 --timeout 180")
	default:
		fmt.Fprintf(os.Stderr, "unknown help topic: %s\n", command)
		return 2
	}
	return 0
}

func cmdLogin(args []string) int {
	fs := newFlagSet("login")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "usage: uncut login")
		return 2
	}

	key, err := readSecret("Enter API key: ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "login failed: %v\n", err)
		return 1
	}
	key = strings.TrimSpace(key)
	if key == "" {
		fmt.Fprintln(os.Stderr, "login failed: empty API key")
		return 2
	}
	endpoint, err := readLine("Enter API endpoint: ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "login failed: %v\n", err)
		return 1
	}
	endpoint, err = normalizeEndpoint(endpoint)
	if err != nil {
		fmt.Fprintf(os.Stderr, "login failed: %v\n", err)
		return 2
	}
	if err := saveConfig(key, endpoint); err != nil {
		fmt.Fprintf(os.Stderr, "login failed: %v\n", err)
		return 1
	}
	fmt.Println("login success!")
	return 0
}

func cmdLogout(args []string) int {
	fs := newFlagSet("logout")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "usage: uncut logout")
		return 2
	}

	path, err := configPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "logout failed: %v\n", err)
		return 1
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		fmt.Fprintf(os.Stderr, "logout failed: %v\n", err)
		return 1
	}
	fmt.Println("logout success!")
	return 0
}

func cmdConfig(args []string) int {
	fs := newFlagSet("config")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "usage: uncut config")
		return 2
	}

	path, err := configPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config failed: %v\n", err)
		return 1
	}
	creds, err := loadCredentials()
	if errors.Is(err, errMissingAPIKey) {
		fmt.Println("api key: not configured")
		fmt.Println("endpoint: not configured")
		fmt.Printf("config: %s\n", path)
		return 0
	}
	if errors.Is(err, errMissingEndpoint) {
		fmt.Println("api key: configured")
		fmt.Println("endpoint: not configured")
		fmt.Printf("config: %s\n", path)
		return 0
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "config failed: %v\n", err)
		return 1
	}
	fmt.Printf("api key: %s\n", maskKey(creds.APIKey))
	fmt.Printf("api key source: %s\n", creds.APIKeySource)
	fmt.Println("endpoint: configured")
	fmt.Printf("endpoint source: %s\n", creds.EndpointSource)
	fmt.Printf("config: %s\n", path)
	return 0
}

func cmdBalance(args []string) int {
	fs := newFlagSet("balance")
	jsonOut := fs.Bool("json", false, "print raw JSON")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "usage: uncut balance [--json]")
		return 2
	}

	c, code := authedClient()
	if code != 0 {
		return code
	}

	var balance walletBalance
	raw, _, err := c.get("/wallet", nil, &balance)
	if err != nil {
		return printCommandError(err)
	}
	if *jsonOut {
		printRawJSON(raw)
		return 0
	}
	fmt.Printf("Balance: %s\n", formatBalance(balance.Balance, balance.Currency))
	return 0
}

func cmdDeposit(args []string) int {
	fs := newFlagSet("deposit")
	jsonOut := fs.Bool("json", false, "print raw JSON")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "usage: uncut deposit [--json]")
		return 2
	}

	c, code := authedClient()
	if code != 0 {
		return code
	}

	var addresses []depositAddress
	raw, _, err := c.get("/wallet/deposit-addresses", nil, &addresses)
	if err != nil {
		return printCommandError(err)
	}
	if *jsonOut {
		printRawJSON(raw)
		return 0
	}
	if len(addresses) == 0 {
		fmt.Println("No deposit addresses.")
		return 0
	}
	fmt.Println("USDT deposit addresses")
	fmt.Println()
	for i, address := range addresses {
		if i > 0 {
			fmt.Println()
		}
		fmt.Printf("%s %s\n", address.Chain, address.Token)
		fmt.Println(address.Address)
	}
	return 0
}

func cmdWallet(args []string) int {
	opts, err := parseOptions(args, []string{"json", "all"}, []string{"page"})
	if err != nil {
		fmt.Fprintf(os.Stderr, "wallet failed: %v\n", err)
		return 2
	}
	if len(opts.positionals) != 0 {
		fmt.Fprintln(os.Stderr, "usage: uncut wallet [--page <n>|--all] [--json]")
		return 2
	}

	page, err := parsePage(opts.valueDefault("page", "1"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "wallet failed: %v\n", err)
		return 2
	}

	c, code := authedClient()
	if code != 0 {
		return code
	}

	if opts.bool("all") {
		return printAllWalletTransactions(c, opts.bool("json"))
	}

	items, raw, env, err := fetchWalletTransactions(c, page)
	if err != nil {
		return printCommandError(err)
	}
	if opts.bool("json") {
		printRawJSON(raw)
		return 0
	}
	printWalletTransactions(items, env)
	return 0
}

func cmdBins(args []string) int {
	fs := newFlagSet("bins")
	jsonOut := fs.Bool("json", false, "print raw JSON")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "usage: uncut bins [--json]")
		return 2
	}

	c, code := authedClient()
	if code != 0 {
		return code
	}

	bins, raw, err := fetchBins(c)
	if err != nil {
		return printCommandError(err)
	}
	if *jsonOut {
		printRawJSON(raw)
		return 0
	}
	printBins(bins)
	return 0
}

func cmdNew(args []string) int {
	opts, err := parseOptions(args, []string{"3ds", "wait", "json"}, []string{"bin", "name", "currency", "topup", "phone", "interval", "timeout"})
	if err != nil {
		fmt.Fprintf(os.Stderr, "new failed: %v\n", err)
		return 2
	}
	if len(opts.positionals) > 1 {
		fmt.Fprintln(os.Stderr, "usage: uncut new <bin_id> [--name <name>] [--currency <code>] [--topup <amount>] [--3ds --phone <number>] [--wait] [--json]")
		return 2
	}

	binID := strings.TrimSpace(opts.value("bin"))
	if binID == "" && len(opts.positionals) == 1 {
		binID = strings.TrimSpace(opts.positionals[0])
	}
	if binID != "" && len(opts.positionals) == 1 && strings.TrimSpace(opts.value("bin")) != "" {
		fmt.Fprintln(os.Stderr, "new failed: pass BIN either as <bin_id> or --bin, not both")
		return 2
	}

	var c *client
	if binID == "" {
		fmt.Fprintln(os.Stderr, "error: --bin is required")
		fmt.Fprintln(os.Stderr)
		var code int
		c, code = authedClient()
		if code != 0 {
			return code
		}
		bins, _, err := fetchBins(c)
		if err != nil {
			return printCommandError(err)
		}
		printNewMissingBinHelp(bins, defaultCardName())
		return 2
	}

	var code int
	c, code = authedClient()
	if code != 0 {
		return code
	}

	name := strings.TrimSpace(opts.value("name"))
	if name == "" {
		name = defaultCardName()
	}
	currency := strings.ToUpper(strings.TrimSpace(opts.valueDefault("currency", "USD")))
	topupRaw := strings.TrimSpace(opts.valueDefault("topup", "0"))
	if currency == "" {
		currency = "USD"
	}

	topupAmount, err := parseNonNegativeFloat(topupRaw, "--topup")
	if err != nil {
		fmt.Fprintf(os.Stderr, "new failed: %v\n", err)
		return 2
	}
	var waitOpts waitOptions
	if opts.bool("wait") {
		waitOpts, err = waitOptionsFrom(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "new failed: %v\n", err)
			return 2
		}
	}

	if existing, suggestion, err := findDuplicateCardName(c, name); err != nil {
		return printCommandError(fmt.Errorf("cannot check existing card names: %w", err))
	} else if existing != nil {
		fmt.Fprintf(os.Stderr, "new failed: card name must be unique; %q already exists\n", name)
		fmt.Fprintf(os.Stderr, "existing card: %s\n", existing.ID)
		fmt.Fprintf(os.Stderr, "try: uncut new %s --name %s --currency %s --topup %s --wait\n", binID, shellQuote(suggestion), currency, topupRaw)
		return 2
	}

	phone := strings.TrimSpace(opts.value("phone"))
	body := map[string]any{
		"bin_id":       binID,
		"name":         name,
		"currency":     currency,
		"topup_amount": topupAmount,
	}
	if opts.bool("3ds") || phone != "" {
		body["enable_3ds"] = true
	}
	if phone != "" {
		body["phone"] = phone
	}

	var queued queuedOperation
	raw, _, err := c.post("/cards", body, &queued)
	if err != nil {
		return printCommandError(err)
	}
	if opts.bool("json") && !opts.bool("wait") {
		printRawJSON(raw)
		return 0
	}
	if !opts.bool("json") {
		printQueuedOperation(queued)
	}
	if opts.bool("wait") {
		return waitForOperation(c, queued.OperationID, waitOpts, opts.bool("json"))
	}
	return 0
}

func cmdTopup(args []string) int {
	return cmdCardAmountOperation(args, "topup", "/cards/%s/topup", "uncut topup <card_id> --amount <amount> [--wait] [--json]")
}

func cmdWithdraw(args []string) int {
	return cmdCardAmountOperation(args, "withdraw", "/cards/%s/withdraw", "uncut withdraw <card_id> --amount <amount> [--wait] [--json]")
}

func cmdCardAmountOperation(args []string, command string, pathFormat string, usage string) int {
	opts, err := parseOptions(args, []string{"wait", "json"}, []string{"amount", "interval", "timeout"})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s failed: %v\n", command, err)
		if maybeAmountOptionError(err) {
			fmt.Fprintln(os.Stderr, "hint: amount must be passed as `--amount 60` or positional `60`, not `--60`")
		}
		printUsageWithCardExamples(command, usage)
		return 2
	}
	if len(opts.positionals) < 1 || len(opts.positionals) > 2 {
		printUsageWithCardExamples(command, usage)
		return 2
	}
	cardID := opts.positionals[0]
	amountRaw := opts.value("amount")
	if amountRaw == "" && len(opts.positionals) == 2 {
		amountRaw = opts.positionals[1]
	}
	amount, err := parsePositiveFloat(amountRaw, "--amount")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s failed: %v\n", command, err)
		printUsageWithCardExamples(command, usage)
		return 2
	}
	var waitOpts waitOptions
	if opts.bool("wait") {
		waitOpts, err = waitOptionsFrom(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s failed: %v\n", command, err)
			return 2
		}
	}

	c, code := authedClient()
	if code != 0 {
		return code
	}

	var queued queuedOperation
	raw, _, err := c.post(fmt.Sprintf(pathFormat, url.PathEscape(cardID)), map[string]any{"amount": amount}, &queued)
	if err != nil {
		return printCommandError(err)
	}
	if opts.bool("json") && !opts.bool("wait") {
		printRawJSON(raw)
		return 0
	}
	if !opts.bool("json") {
		printQueuedOperation(queued)
	}
	if opts.bool("wait") {
		return waitForOperation(c, queued.OperationID, waitOpts, opts.bool("json"))
	}
	return 0
}

func cmdTransactions(args []string) int {
	opts, err := parseOptions(args, []string{"json", "all"}, []string{"page"})
	if err != nil {
		fmt.Fprintf(os.Stderr, "transactions failed: %v\n", err)
		printUsageWithCardExamples("transactions", "uncut transactions <card_id> [--page <n>|--all] [--json]")
		return 2
	}
	if len(opts.positionals) != 1 {
		printUsageWithCardExamples("transactions", "uncut transactions <card_id> [--page <n>|--all] [--json]")
		return 2
	}
	cardID := opts.positionals[0]
	page, err := parsePage(opts.valueDefault("page", "1"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "transactions failed: %v\n", err)
		return 2
	}

	c, code := authedClient()
	if code != 0 {
		return code
	}

	if opts.bool("all") {
		return printAllCardTransactions(c, cardID, opts.bool("json"))
	}

	items, raw, env, err := fetchCardTransactions(c, cardID, page)
	if err != nil {
		return printCommandError(err)
	}
	if opts.bool("json") {
		printRawJSON(raw)
		return 0
	}
	printCardTransactions(cardID, items, env)
	return 0
}

func cmdRename(args []string) int {
	opts, err := parseOptions(args, []string{"json"}, []string{"name"})
	if err != nil {
		fmt.Fprintf(os.Stderr, "rename failed: %v\n", err)
		printUsageWithCardExamples("rename", "uncut rename <card_id> <new_name> [--json]")
		return 2
	}
	if len(opts.positionals) < 1 || len(opts.positionals) > 2 {
		printUsageWithCardExamples("rename", "uncut rename <card_id> <new_name> [--json]")
		return 2
	}
	cardID := opts.positionals[0]
	name := strings.TrimSpace(opts.value("name"))
	if name == "" && len(opts.positionals) == 2 {
		name = strings.TrimSpace(opts.positionals[1])
	}
	if name == "" {
		fmt.Fprintln(os.Stderr, "rename failed: new name is required")
		printUsageWithCardExamples("rename", "uncut rename <card_id> <new_name> [--json]")
		return 2
	}

	c, code := authedClient()
	if code != 0 {
		return code
	}

	var item card
	raw, _, err := c.patch("/cards/"+url.PathEscape(cardID), map[string]any{"name": name}, &item)
	if err != nil {
		return printCommandError(err)
	}
	if opts.bool("json") {
		printRawJSON(raw)
		return 0
	}
	printCard(item, false)
	return 0
}

func cmdPhone(args []string) int {
	opts, err := parseOptions(args, []string{"json"}, []string{"phone"})
	if err != nil {
		fmt.Fprintf(os.Stderr, "phone failed: %v\n", err)
		printUsageWithCardExamples("phone", "uncut phone <card_id> --phone <e164_phone> [--json]")
		return 2
	}
	if len(opts.positionals) < 1 || len(opts.positionals) > 2 {
		printUsageWithCardExamples("phone", "uncut phone <card_id> --phone <e164_phone> [--json]")
		return 2
	}
	cardID := opts.positionals[0]
	phone := strings.TrimSpace(opts.value("phone"))
	if phone == "" && len(opts.positionals) == 2 {
		phone = strings.TrimSpace(opts.positionals[1])
	}
	if phone == "" {
		fmt.Fprintln(os.Stderr, "phone failed: --phone is required")
		printUsageWithCardExamples("phone", "uncut phone <card_id> --phone <e164_phone> [--json]")
		return 2
	}

	c, code := authedClient()
	if code != 0 {
		return code
	}

	var item card
	raw, _, err := c.put("/cards/"+url.PathEscape(cardID)+"/3ds-phone", map[string]any{"phone": phone}, &item)
	if err != nil {
		return printCommandError(err)
	}
	if opts.bool("json") {
		printRawJSON(raw)
		return 0
	}
	printCard(item, false)
	return 0
}

func cmdFreeze(args []string) int {
	return cmdCardStatusOperation(args, "freeze", "/cards/%s/freeze", "uncut freeze <card_id> [--json]")
}

func cmdUnfreeze(args []string) int {
	return cmdCardStatusOperation(args, "unfreeze", "/cards/%s/unfreeze", "uncut unfreeze <card_id> [--json]")
}

func cmdCardStatusOperation(args []string, command string, pathFormat string, usage string) int {
	opts, err := parseOptions(args, []string{"json"}, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s failed: %v\n", command, err)
		printUsageWithCardExamples(command, usage)
		return 2
	}
	if len(opts.positionals) != 1 {
		printUsageWithCardExamples(command, usage)
		return 2
	}

	c, code := authedClient()
	if code != 0 {
		return code
	}

	var item card
	raw, _, err := c.post(fmt.Sprintf(pathFormat, url.PathEscape(opts.positionals[0])), nil, &item)
	if err != nil {
		return printCommandError(err)
	}
	if opts.bool("json") {
		printRawJSON(raw)
		return 0
	}
	printCard(item, false)
	return 0
}

func cmdDelete(args []string) int {
	opts, err := parseOptions(args, []string{"yes", "json"}, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "delete failed: %v\n", err)
		printUsageWithCardExamples("delete", "uncut delete <card_id> [--yes] [--json]")
		return 2
	}
	if len(opts.positionals) != 1 {
		printUsageWithCardExamples("delete", "uncut delete <card_id> [--yes] [--json]")
		return 2
	}
	cardID := opts.positionals[0]

	c, code := authedClient()
	if code != 0 {
		return code
	}

	if !opts.bool("yes") {
		confirmed, err := confirmDelete(cardID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "delete failed: %v\n", err)
			return 1
		}
		if !confirmed {
			fmt.Println("delete cancelled")
			return 1
		}
	}

	var result deleteResult
	raw, _, err := c.delete("/cards/"+url.PathEscape(cardID), &result)
	if err != nil {
		return printCommandError(err)
	}
	if opts.bool("json") {
		printRawJSON(raw)
		return 0
	}
	printDeleteResult(result)
	return 0
}

func cmdOperation(args []string) int {
	opts, err := parseOptions(args, []string{"json"}, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "operation failed: %v\n", err)
		return 2
	}
	if len(opts.positionals) != 1 {
		fmt.Fprintln(os.Stderr, "usage: uncut operation <operation_id> [--json]")
		return 2
	}

	c, code := authedClient()
	if code != 0 {
		return code
	}

	var op operation
	raw, _, err := c.get("/operations/"+url.PathEscape(opts.positionals[0]), nil, &op)
	if err != nil {
		return printCommandError(err)
	}
	if opts.bool("json") {
		printRawJSON(raw)
		return 0
	}
	printOperation(op)
	return 0
}

func cmdWait(args []string) int {
	opts, err := parseOptions(args, []string{"json"}, []string{"interval", "timeout"})
	if err != nil {
		fmt.Fprintf(os.Stderr, "wait failed: %v\n", err)
		return 2
	}
	if len(opts.positionals) != 1 {
		fmt.Fprintln(os.Stderr, "usage: uncut wait <operation_id> [--interval <seconds>] [--timeout <seconds>] [--json]")
		return 2
	}
	waitOpts, err := waitOptionsFrom(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "wait failed: %v\n", err)
		return 2
	}

	c, code := authedClient()
	if code != 0 {
		return code
	}
	return waitForOperation(c, opts.positionals[0], waitOpts, opts.bool("json"))
}

func cmdCards(args []string) int {
	fs := newFlagSet("cards")
	reveal := fs.Bool("reveal", false, "show full PAN/CVV")
	full := fs.Bool("full", false, "alias for --reveal")
	jsonOut := fs.Bool("json", false, "print JSON")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "usage: uncut cards [--reveal|--full] [--json]")
		return 2
	}

	c, code := authedClient()
	if code != 0 {
		return code
	}

	var cards []card
	raw, _, err := c.get("/cards", nil, &cards)
	if err != nil {
		return printCommandError(err)
	}
	showFull := *reveal || *full
	if !showFull {
		if *jsonOut {
			printRawJSON(raw)
			return 0
		}
		printCards(cards, false)
		return 0
	}

	details := make([]card, 0, len(cards))
	for _, item := range cards {
		var detail card
		_, _, err := c.get("/cards/"+url.PathEscape(item.ID)+"/details", nil, &detail)
		if err != nil {
			return printCommandError(fmt.Errorf("cannot reveal card %s: %w", item.ID, err))
		}
		details = append(details, detail)
	}
	if *jsonOut {
		printJSON(map[string]any{"data": details})
		return 0
	}
	printCards(details, true)
	return 0
}

func cmdCard(args []string) int {
	cardID, showFull, jsonOut, ok := parseCardArgs(args)
	if !ok {
		printUsageWithCardExamples("card", "uncut card <card_id> [--reveal|--full] [--json]")
		return 2
	}
	c, code := authedClient()
	if code != 0 {
		return code
	}

	path := "/cards/" + url.PathEscape(cardID)
	if showFull {
		path += "/details"
	}

	var item card
	raw, _, err := c.get(path, nil, &item)
	if err != nil {
		return printCommandError(err)
	}
	if jsonOut {
		printRawJSON(raw)
		return 0
	}
	printCard(item, showFull)
	return 0
}

func parseCardArgs(args []string) (cardID string, reveal bool, jsonOut bool, ok bool) {
	for _, arg := range args {
		switch arg {
		case "--reveal", "--full":
			reveal = true
		case "--json":
			jsonOut = true
		default:
			if strings.HasPrefix(arg, "-") {
				fmt.Fprintf(os.Stderr, "unknown card option: %s\n", arg)
				return "", false, false, false
			}
			if cardID != "" {
				return "", false, false, false
			}
			cardID = arg
		}
	}
	return cardID, reveal, jsonOut, cardID != ""
}

type parsedOptions struct {
	positionals []string
	bools       map[string]bool
	values      map[string]string
}

func (p parsedOptions) bool(name string) bool {
	return p.bools[name]
}

func (p parsedOptions) value(name string) string {
	return p.values[name]
}

func (p parsedOptions) valueDefault(name string, fallback string) string {
	if value := p.values[name]; value != "" {
		return value
	}
	return fallback
}

func parseOptions(args []string, boolFlags []string, valueFlags []string) (parsedOptions, error) {
	boolSet := make(map[string]bool)
	valueSet := make(map[string]bool)
	for _, name := range boolFlags {
		boolSet[name] = true
	}
	for _, name := range valueFlags {
		valueSet[name] = true
	}

	opts := parsedOptions{
		bools:  make(map[string]bool),
		values: make(map[string]string),
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if !strings.HasPrefix(arg, "--") || arg == "--" {
			if arg == "--" {
				opts.positionals = append(opts.positionals, args[i+1:]...)
				break
			}
			opts.positionals = append(opts.positionals, arg)
			continue
		}

		rawName := strings.TrimPrefix(arg, "--")
		name, value, hasValue := strings.Cut(rawName, "=")
		if boolSet[name] {
			if hasValue {
				parsed, err := parseBoolValue(value)
				if err != nil {
					return opts, fmt.Errorf("--%s expects a boolean value", name)
				}
				opts.bools[name] = parsed
			} else {
				opts.bools[name] = true
			}
			continue
		}
		if valueSet[name] {
			if !hasValue {
				if i+1 >= len(args) || strings.HasPrefix(args[i+1], "--") {
					return opts, fmt.Errorf("--%s requires a value", name)
				}
				i++
				value = args[i]
			}
			opts.values[name] = value
			continue
		}
		return opts, fmt.Errorf("unknown option --%s", name)
	}

	return opts, nil
}

func parseBoolValue(value string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "on":
		return true, nil
	case "0", "false", "no", "off":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value")
	}
}

func parsePage(raw string) (int, error) {
	page, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || page < 1 {
		return 0, fmt.Errorf("--page must be a positive integer")
	}
	return page, nil
}

func parsePositiveFloat(raw string, name string) (float64, error) {
	value, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil || value <= 0 {
		return 0, fmt.Errorf("%s must be a positive number", name)
	}
	return value, nil
}

func parseNonNegativeFloat(raw string, name string) (float64, error) {
	value, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil || value < 0 {
		return 0, fmt.Errorf("%s must be zero or a positive number", name)
	}
	return value, nil
}

func newFlagSet(name string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	return fs
}

func authedClient() (*client, int) {
	creds, err := loadCredentials()
	if errors.Is(err, errMissingAPIKey) {
		fmt.Fprintln(os.Stderr, "api key is not configured")
		fmt.Fprintln(os.Stderr, "run: uncut login")
		return nil, 3
	}
	if errors.Is(err, errMissingEndpoint) {
		fmt.Fprintln(os.Stderr, "endpoint is not configured")
		fmt.Fprintln(os.Stderr, "run: uncut login")
		return nil, 3
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "auth failed: %v\n", err)
		return nil, 1
	}
	return &client{
		apiKey:   creds.APIKey,
		endpoint: creds.Endpoint,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, 0
}

func (c *client) get(path string, query map[string]string, target any) ([]byte, *envelope, error) {
	return c.do(http.MethodGet, path, query, nil, target)
}

func (c *client) post(path string, body any, target any) ([]byte, *envelope, error) {
	return c.do(http.MethodPost, path, nil, body, target)
}

func (c *client) patch(path string, body any, target any) ([]byte, *envelope, error) {
	return c.do(http.MethodPatch, path, nil, body, target)
}

func (c *client) put(path string, body any, target any) ([]byte, *envelope, error) {
	return c.do(http.MethodPut, path, nil, body, target)
}

func (c *client) delete(path string, target any) ([]byte, *envelope, error) {
	return c.do(http.MethodDelete, path, nil, nil, target)
}

func (c *client) do(method string, path string, query map[string]string, body any, target any) ([]byte, *envelope, error) {
	endpoint := c.endpoint + path
	if len(query) > 0 {
		values := url.Values{}
		for key, value := range query {
			values.Set(key, value)
		}
		endpoint += "?" + values.Encode()
	}

	var requestBody io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return nil, nil, err
		}
		requestBody = bytes.NewReader(raw)
	}

	req, err := http.NewRequest(method, endpoint, requestBody)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := c.http.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	var env envelope
	if len(bytes.TrimSpace(raw)) > 0 {
		dec := json.NewDecoder(bytes.NewReader(raw))
		dec.UseNumber()
		if err := dec.Decode(&env); err != nil {
			if res.StatusCode >= 200 && res.StatusCode < 300 {
				return raw, nil, fmt.Errorf("non-JSON API response: %w", err)
			}
			return raw, nil, &apiError{Status: res.StatusCode, Body: string(raw)}
		}
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		apiErr := &apiError{Status: res.StatusCode, Body: string(raw)}
		if env.Error != nil {
			apiErr.Code = env.Error.Code
			apiErr.Message = env.Error.Message
			apiErr.Details = env.Error.Details
		}
		return raw, &env, apiErr
	}

	if target != nil && len(env.Data) > 0 {
		dec := json.NewDecoder(bytes.NewReader(env.Data))
		dec.UseNumber()
		if err := dec.Decode(target); err != nil {
			return raw, &env, err
		}
	}
	return raw, &env, nil
}

func fetchBins(c *client) ([]cardBIN, []byte, error) {
	var bins []cardBIN
	raw, _, err := c.get("/card-bins", nil, &bins)
	return bins, raw, err
}

func fetchCards(c *client) ([]card, []byte, error) {
	var cards []card
	raw, _, err := c.get("/cards", nil, &cards)
	return cards, raw, err
}

func findDuplicateCardName(c *client, name string) (*card, string, error) {
	cards, _, err := fetchCards(c)
	if err != nil {
		return nil, "", err
	}
	normalized := normalizeCardName(name)
	for i := range cards {
		if normalizeCardName(cards[i].Name) == normalized {
			suggestion := suggestUniqueCardName(cards, name)
			return &cards[i], suggestion, nil
		}
	}
	return nil, "", nil
}

func normalizeCardName(name string) string {
	return strings.ToLower(strings.Join(strings.Fields(name), " "))
}

func suggestUniqueCardName(cards []card, base string) string {
	base = strings.TrimSpace(base)
	if base == "" {
		base = defaultCardName()
	}
	used := make(map[string]bool, len(cards))
	for _, item := range cards {
		used[normalizeCardName(item.Name)] = true
	}
	for i := 2; i < 1000; i++ {
		candidate := fmt.Sprintf("%s-%d", base, i)
		if !used[normalizeCardName(candidate)] {
			return candidate
		}
	}
	return fmt.Sprintf("%s-%d", base, time.Now().Unix())
}

func fetchWalletTransactions(c *client, page int) ([]walletTransaction, []byte, *envelope, error) {
	var items []walletTransaction
	raw, env, err := c.get("/wallet/transactions", map[string]string{"page": strconv.Itoa(page)}, &items)
	return items, raw, env, err
}

func fetchCardTransactions(c *client, cardID string, page int) ([]cardTransaction, []byte, *envelope, error) {
	var items []cardTransaction
	raw, env, err := c.get("/cards/"+url.PathEscape(cardID)+"/transactions", map[string]string{"page": strconv.Itoa(page)}, &items)
	return items, raw, env, err
}

func printAllWalletTransactions(c *client, jsonOut bool) int {
	var all []walletTransaction
	page := 1
	var lastEnv *envelope
	for {
		items, _, env, err := fetchWalletTransactions(c, page)
		if err != nil {
			return printCommandError(err)
		}
		all = append(all, items...)
		lastEnv = env
		if env == nil || !env.HasMore {
			break
		}
		page++
	}
	if jsonOut {
		printJSON(map[string]any{"data": all})
		return 0
	}
	printWalletTransactions(all, lastEnv)
	return 0
}

func printAllCardTransactions(c *client, cardID string, jsonOut bool) int {
	var all []cardTransaction
	page := 1
	var lastEnv *envelope
	for {
		items, _, env, err := fetchCardTransactions(c, cardID, page)
		if err != nil {
			return printCommandError(err)
		}
		all = append(all, items...)
		lastEnv = env
		if env == nil || !env.HasMore {
			break
		}
		page++
	}
	if jsonOut {
		printJSON(map[string]any{"data": all})
		return 0
	}
	printCardTransactions(cardID, all, lastEnv)
	return 0
}

type waitOptions struct {
	Interval time.Duration
	Timeout  time.Duration
}

func waitOptionsFrom(opts parsedOptions) (waitOptions, error) {
	interval := 3
	timeout := 120
	if raw := opts.value("interval"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed <= 0 {
			return waitOptions{}, fmt.Errorf("--interval must be a positive integer")
		}
		interval = parsed
	}
	if raw := opts.value("timeout"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed <= 0 {
			return waitOptions{}, fmt.Errorf("--timeout must be a positive integer")
		}
		timeout = parsed
	}
	return waitOptions{
		Interval: time.Duration(interval) * time.Second,
		Timeout:  time.Duration(timeout) * time.Second,
	}, nil
}

func waitForOperation(c *client, operationID string, opts waitOptions, jsonOut bool) int {
	if strings.TrimSpace(operationID) == "" {
		fmt.Fprintln(os.Stderr, "wait failed: operation id is empty")
		return 2
	}
	deadline := time.Now().Add(opts.Timeout)
	for {
		var op operation
		raw, _, err := c.get("/operations/"+url.PathEscape(operationID), nil, &op)
		if err != nil {
			return printCommandError(err)
		}
		status := strings.ToLower(strings.TrimSpace(op.Status))
		if status == "completed" || status == "error" {
			if jsonOut {
				printRawJSON(raw)
			} else {
				printOperation(op)
			}
			if status == "error" {
				return 1
			}
			return 0
		}
		if !jsonOut {
			fmt.Printf("status: %s\n", fallback(op.Status, "pending"))
		}
		if time.Now().Add(opts.Interval).After(deadline) {
			fmt.Fprintf(os.Stderr, "wait failed: timeout after %.0f seconds\n", opts.Timeout.Seconds())
			return 1
		}
		time.Sleep(opts.Interval)
	}
}

func confirmDelete(cardID string) (bool, error) {
	fmt.Fprintf(os.Stderr, "delete card %s? type \"delete\" to confirm: ", cardID)
	value, err := stdinReader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return false, err
	}
	return strings.TrimSpace(value) == "delete", nil
}

func configPath() (string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		var err error
		home, err = os.UserHomeDir()
		if err != nil {
			return "", err
		}
	}
	return filepath.Join(home, ".config", "uncut", "config.json"), nil
}

func saveConfig(apiKey string, endpoint string) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	cfg := configFile{
		APIKey:    apiKey,
		Endpoint:  endpoint,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	raw, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	raw = append(raw, '\n')
	if err := os.WriteFile(path, raw, 0600); err != nil {
		return err
	}
	return os.Chmod(path, 0600)
}

func loadCredentials() (credentials, error) {
	var creds credentials
	path, err := configPath()
	if err != nil {
		return creds, err
	}
	raw, err := os.ReadFile(path)
	var cfg configFile
	if err == nil {
		if err := json.Unmarshal(raw, &cfg); err != nil {
			return creds, err
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return creds, err
	}

	if key := strings.TrimSpace(os.Getenv("UNCUT_API_KEY")); key != "" {
		creds.APIKey = key
		creds.APIKeySource = "env UNCUT_API_KEY"
	} else if key := strings.TrimSpace(cfg.APIKey); key != "" {
		creds.APIKey = key
		creds.APIKeySource = path
	}

	if endpoint := strings.TrimSpace(os.Getenv("UNCUT_BASE_URL")); endpoint != "" {
		normalized, err := normalizeEndpoint(endpoint)
		if err != nil {
			return creds, fmt.Errorf("invalid UNCUT_BASE_URL: %w", err)
		}
		creds.Endpoint = normalized
		creds.EndpointSource = "env UNCUT_BASE_URL"
	} else if endpoint := strings.TrimSpace(cfg.Endpoint); endpoint != "" {
		normalized, err := normalizeEndpoint(endpoint)
		if err != nil {
			return creds, fmt.Errorf("invalid saved endpoint: %w", err)
		}
		creds.Endpoint = normalized
		creds.EndpointSource = path
	}

	if creds.APIKey == "" {
		return creds, errMissingAPIKey
	}
	if creds.Endpoint == "" {
		return creds, errMissingEndpoint
	}
	return creds, nil
}

func readSecret(prompt string) (string, error) {
	fmt.Fprint(os.Stderr, prompt)
	restore, hidden := disableEcho()
	defer restore()

	value, err := stdinReader.ReadString('\n')
	if hidden {
		fmt.Fprintln(os.Stderr)
	}
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	return strings.TrimSpace(value), nil
}

func readLine(prompt string) (string, error) {
	fmt.Fprint(os.Stderr, prompt)
	value, err := stdinReader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	return strings.TrimSpace(value), nil
}

func normalizeEndpoint(raw string) (string, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", errors.New("empty endpoint")
	}
	parsed, err := url.Parse(value)
	if err != nil {
		return "", err
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", errors.New("endpoint must start with http:// or https://")
	}
	if parsed.Host == "" {
		return "", errors.New("endpoint must include a host")
	}
	parsed.Path = strings.TrimRight(parsed.Path, "/")
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return strings.TrimRight(parsed.String(), "/"), nil
}

func disableEcho() (func(), bool) {
	oldState := exec.Command("stty", "-g")
	oldState.Stdin = os.Stdin
	rawState, err := oldState.Output()
	if err != nil {
		return func() {}, false
	}

	off := exec.Command("stty", "-echo")
	off.Stdin = os.Stdin
	if err := off.Run(); err != nil {
		return func() {}, false
	}

	state := strings.TrimSpace(string(rawState))
	return func() {
		restore := exec.Command("stty", state)
		restore.Stdin = os.Stdin
		_ = restore.Run()
	}, true
}

func printCards(items []card, reveal bool) {
	if len(items) == 0 {
		fmt.Println("No cards.")
		return
	}
	for i, item := range items {
		if i > 0 {
			fmt.Println()
		}
		printCard(item, reveal)
	}
}

func printCard(item card, reveal bool) {
	name := strings.TrimSpace(item.Name)
	if name == "" {
		name = "(unnamed card)"
	}

	fmt.Println(name)
	if item.ID != "" {
		fmt.Printf("id: %s\n", item.ID)
	}
	if reveal && item.CardNumber != "" {
		fmt.Printf("💳 %s\n", groupPAN(item.CardNumber))
	} else {
		fmt.Printf("💳 %s\n", formatMask(item.Mask))
	}
	fmt.Printf("EXP:%s  CVV:%s  %s\n", formatExpiry(item, reveal), formatCVV(item, reveal), formatStatus(item.Status))
	fmt.Printf("Balance: %s\n", formatBalance(item.Balance, item.Currency))
}

func printBins(items []cardBIN) {
	if len(items) == 0 {
		fmt.Println("No BINs.")
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "id\tname\tcurrencies\tissue\ttopup\tauth\twithdraw\twallet")
	for _, item := range items {
		wallet := "no"
		if item.WalletSupport {
			wallet = "yes"
		}
		fmt.Fprintf(
			w,
			"%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			item.ID,
			fallback(item.Name, "-"),
			strings.Join(item.Currencies, ","),
			formatFixedFee(item.Rates.CardIssueFee, "USDT"),
			formatPercentFee(item.Rates.CardTopupPercent, item.Rates.CardTopupFix),
			formatPercentFee(item.Rates.CardAuthorizationPercent, item.Rates.CardAuthorizationFix),
			formatPercentFee(item.Rates.CardWithdrawPercent, item.Rates.CardWithdrawFix),
			wallet,
		)
	}
	_ = w.Flush()
}

func printNewMissingBinHelp(items []cardBIN, generatedName string) {
	fmt.Println("available bins:")
	printBins(items)
	fmt.Println()
	fmt.Println("defaults:")
	fmt.Printf("  name: %s\n", generatedName)
	fmt.Println("  currency: USD")
	fmt.Println("  topup: 0")
	fmt.Println()
	fmt.Println("copy-paste create commands:")
	if len(items) == 0 {
		fmt.Println("  no available BINs")
		return
	}
	for _, item := range items {
		fmt.Printf("  uncut new %s --name %s --wait\n", item.ID, shellQuote(generatedName))
	}
}

func printHelpBINExamples() {
	bins := helpBins()
	if len(bins) == 0 {
		fmt.Println()
		fmt.Println("Current BIN examples: run `uncut new` or `uncut bins` after login.")
		return
	}
	name := defaultCardName()
	fmt.Println()
	fmt.Println("Current BIN examples:")
	for _, item := range bins {
		fmt.Printf("  uncut new %s --name %s --wait\n", item.ID, shellQuote(name))
	}
}

func printHelpCardExamples(command string) {
	cards := helpCards()
	if len(cards) == 0 {
		fmt.Println()
		fmt.Println("Current card examples: run `uncut cards` after login.")
		return
	}
	fmt.Println()
	fmt.Println("Current card examples:")
	for _, item := range cards {
		switch command {
		case "card":
			fmt.Printf("  uncut card %s\n", item.ID)
			fmt.Printf("  uncut card %s --reveal\n", item.ID)
		case "topup":
			fmt.Printf("  uncut topup %s --amount 60 --wait\n", item.ID)
		case "withdraw":
			fmt.Printf("  uncut withdraw %s --amount 20 --wait\n", item.ID)
		case "transactions":
			fmt.Printf("  uncut transactions %s\n", item.ID)
		case "rename":
			fmt.Printf("  uncut rename %s %s\n", item.ID, shellQuote(item.Name+"-new"))
		case "phone":
			fmt.Printf("  uncut phone %s --phone +10000000000\n", item.ID)
		case "freeze":
			fmt.Printf("  uncut freeze %s\n", item.ID)
		case "unfreeze":
			fmt.Printf("  uncut unfreeze %s\n", item.ID)
		case "delete":
			fmt.Printf("  uncut delete %s\n", item.ID)
		}
	}
}

func printUsageWithCardExamples(command string, usage string) {
	fmt.Fprintf(os.Stderr, "usage: %s\n", usage)
	cards := helpCards()
	if len(cards) == 0 {
		fmt.Fprintln(os.Stderr, "examples: run `uncut cards` and copy a card id")
		return
	}
	fmt.Fprintln(os.Stderr, "examples:")
	for _, item := range cards {
		switch command {
		case "card":
			fmt.Fprintf(os.Stderr, "  uncut card %s\n", item.ID)
			fmt.Fprintf(os.Stderr, "  uncut card %s --reveal\n", item.ID)
		case "topup":
			fmt.Fprintf(os.Stderr, "  uncut topup %s --amount 60 --wait\n", item.ID)
		case "withdraw":
			fmt.Fprintf(os.Stderr, "  uncut withdraw %s --amount 20 --wait\n", item.ID)
		case "transactions":
			fmt.Fprintf(os.Stderr, "  uncut transactions %s\n", item.ID)
		case "rename":
			fmt.Fprintf(os.Stderr, "  uncut rename %s %s\n", item.ID, shellQuote(item.Name+"-new"))
		case "phone":
			fmt.Fprintf(os.Stderr, "  uncut phone %s --phone +10000000000\n", item.ID)
		case "freeze":
			fmt.Fprintf(os.Stderr, "  uncut freeze %s\n", item.ID)
		case "unfreeze":
			fmt.Fprintf(os.Stderr, "  uncut unfreeze %s\n", item.ID)
		case "delete":
			fmt.Fprintf(os.Stderr, "  uncut delete %s\n", item.ID)
		default:
			fmt.Fprintf(os.Stderr, "  %s\n", usage)
		}
	}
}

func maybeAmountOptionError(err error) bool {
	message := err.Error()
	return strings.Contains(message, "unknown option --") && strings.ContainsAny(message, "0123456789")
}

func helpCards() []card {
	c, ok := quietAuthedClient()
	if !ok {
		return nil
	}
	cards, _, err := fetchCards(c)
	if err != nil {
		return nil
	}
	return cards
}

func helpBins() []cardBIN {
	c, ok := quietAuthedClient()
	if !ok {
		return nil
	}
	bins, _, err := fetchBins(c)
	if err != nil {
		return nil
	}
	return bins
}

func quietAuthedClient() (*client, bool) {
	creds, err := loadCredentials()
	if err != nil {
		return nil, false
	}
	return &client{
		apiKey:   creds.APIKey,
		endpoint: creds.Endpoint,
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, true
}

func defaultCardName() string {
	return "card-" + time.Now().Format("20060102-1504")
}

func shellQuote(value string) string {
	if value == "" {
		return "''"
	}
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}

func printQueuedOperation(item queuedOperation) {
	fmt.Printf("operation: %s\n", item.OperationID)
	fmt.Printf("status: %s\n", fallback(item.Status, "new"))
	if item.Message != "" {
		fmt.Println(item.Message)
	}
	if item.OperationID != "" {
		fmt.Printf("next: uncut wait %s\n", item.OperationID)
	}
}

func printOperation(item operation) {
	fmt.Printf("operation: %s\n", item.OperationID)
	fmt.Printf("type: %s\n", fallback(item.Type, "-"))
	fmt.Printf("status: %s\n", fallback(item.Status, "-"))
	if item.Amount.String() != "" {
		fmt.Printf("amount: %s\n", formatNumber(item.Amount))
	}
	if item.CardID != "" {
		fmt.Printf("card: %s\n", item.CardID)
	}
	if item.CreatedAt != "" {
		fmt.Printf("created: %s\n", item.CreatedAt)
	}
	if item.UpdatedAt != "" {
		fmt.Printf("updated: %s\n", item.UpdatedAt)
	}
	if item.ErrorMessage != "" {
		fmt.Printf("error: %s\n", item.ErrorMessage)
	}
}

func printCardTransactions(cardID string, items []cardTransaction, env *envelope) {
	if len(items) == 0 {
		fmt.Println("No card transactions.")
		printPaginationHint("transactions", cardID, env)
		return
	}
	for i, item := range items {
		if i > 0 {
			fmt.Println()
		}
		name := fallback(item.MerchantName, "(unknown merchant)")
		fmt.Println(name)
		if item.Time != "" {
			fmt.Printf("time: %s\n", item.Time)
		}
		fmt.Printf("type: %s\n", fallback(item.Type, "-"))
		if item.MerchantCountry != "" {
			fmt.Printf("country: %s\n", item.MerchantCountry)
		}
		fmt.Printf("amount: %s\n", formatAmountCurrency(item.OriginalAmount, item.OriginalCurrency))
		fmt.Printf("pre-auth: %s\n", formatNumber(item.PreAuthAmount))
		fmt.Printf("posted: %s\n", formatNumber(item.PostedAmount))
	}
	printPaginationHint("transactions", cardID, env)
}

func printWalletTransactions(items []walletTransaction, env *envelope) {
	if len(items) == 0 {
		fmt.Println("No wallet transactions.")
		printPaginationHint("wallet", "", env)
		return
	}
	for i, item := range items {
		if i > 0 {
			fmt.Println()
		}
		fmt.Printf("%s  %s  %s\n", fallback(item.Type, "-"), fallback(item.Status, "-"), formatNumber(item.Amount))
		fmt.Printf("id: %s\n", item.ID)
		if item.Fee.String() != "" {
			fmt.Printf("fee: %s\n", formatNumber(item.Fee))
		}
		if item.CardID != nil && *item.CardID != "" {
			fmt.Printf("card: %s\n", *item.CardID)
		}
		if item.Comment != nil && *item.Comment != "" {
			fmt.Printf("comment: %s\n", *item.Comment)
		}
		if item.CreatedAt != "" {
			fmt.Printf("created: %s\n", item.CreatedAt)
		}
	}
	printPaginationHint("wallet", "", env)
}

func printPaginationHint(command string, target string, env *envelope) {
	if env == nil || env.Page == 0 {
		return
	}
	fmt.Println()
	if env.Source != "" {
		fmt.Printf("source: %s\n", env.Source)
	}
	fmt.Printf("page %d, %d per page\n", env.Page, env.PerPage)
	if env.HasMore {
		next := env.Page + 1
		if target != "" {
			fmt.Printf("next: uncut %s %s --page %d\n", command, target, next)
		} else {
			fmt.Printf("next: uncut %s --page %d\n", command, next)
		}
	}
}

func printDeleteResult(result deleteResult) {
	fmt.Println("card deleted")
	fmt.Printf("returned: %s\n", formatAmountCurrency(result.ReturnedAmount, result.ReturnedCurrency))
	fmt.Printf("wallet balance: %s\n", formatBalance(result.WalletBalance, "USDT"))
}

func formatMask(mask string) string {
	last4 := strings.TrimSpace(strings.TrimLeft(mask, "*"))
	if last4 == "" {
		last4 = "????"
	}
	return "**** **** **** " + last4
}

func groupPAN(pan string) string {
	digits := strings.ReplaceAll(strings.TrimSpace(pan), " ", "")
	if digits == "" {
		return ""
	}
	var parts []string
	for len(digits) > 4 {
		parts = append(parts, digits[:4])
		digits = digits[4:]
	}
	if digits != "" {
		parts = append(parts, digits)
	}
	return strings.Join(parts, " ")
}

func formatExpiry(item card, reveal bool) string {
	if !reveal {
		return "**/**"
	}
	if item.ExpirationDate != "" {
		return item.ExpirationDate
	}
	if item.ExpireMonth != "" && item.ExpireYear != "" {
		return item.ExpireMonth + "/" + twoDigitYear(item.ExpireYear)
	}
	return "**/**"
}

func twoDigitYear(year string) string {
	year = strings.TrimSpace(year)
	if len(year) >= 2 {
		return year[len(year)-2:]
	}
	return year
}

func formatCVV(item card, reveal bool) string {
	if reveal && item.CVV != "" {
		return item.CVV
	}
	return "***"
}

func formatStatus(status string) string {
	status = strings.TrimSpace(status)
	if status == "" {
		return "Unknown"
	}
	status = strings.ToLower(status)
	return strings.ToUpper(status[:1]) + status[1:]
}

func formatBalance(value json.Number, currency string) string {
	raw := strings.TrimSpace(value.String())
	if raw == "" {
		if currency == "" {
			return "-"
		}
		return "- " + currency
	}
	amount, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		if currency == "" {
			return raw
		}
		return raw + " " + currency
	}
	currency = strings.ToUpper(strings.TrimSpace(currency))
	switch currency {
	case "USD":
		return fmt.Sprintf("$%.2f", amount)
	case "USDT":
		return fmt.Sprintf("%.2f USDT", amount)
	case "":
		return fmt.Sprintf("%.2f", amount)
	default:
		return fmt.Sprintf("%.2f %s", amount, currency)
	}
}

func formatAmountCurrency(value json.Number, currency string) string {
	currency = strings.ToUpper(strings.TrimSpace(currency))
	if currency == "" {
		return formatNumber(value)
	}
	raw := strings.TrimSpace(value.String())
	if raw == "" {
		return "- " + currency
	}
	amount, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return raw + " " + currency
	}
	if currency == "USD" {
		return fmt.Sprintf("$%.2f", amount)
	}
	return fmt.Sprintf("%.2f %s", amount, currency)
}

func formatNumber(value json.Number) string {
	raw := strings.TrimSpace(value.String())
	if raw == "" {
		return "-"
	}
	amount, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return raw
	}
	return fmt.Sprintf("%.2f", amount)
}

func formatFixedFee(value json.Number, currency string) string {
	return formatAmountCurrency(value, currency)
}

func formatPercentFee(percent json.Number, fixed json.Number) string {
	percentText := formatNumber(percent) + "%"
	if isZeroNumber(fixed) {
		return percentText
	}
	return percentText + " + " + formatAmountCurrency(fixed, "USDT")
}

func isZeroNumber(value json.Number) bool {
	raw := strings.TrimSpace(value.String())
	if raw == "" {
		return true
	}
	amount, err := strconv.ParseFloat(raw, 64)
	return err == nil && amount == 0
}

func fallback(value string, replacement string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return replacement
	}
	return value
}

func maskKey(key string) string {
	key = strings.TrimSpace(key)
	if len(key) <= 12 {
		return "****"
	}
	return key[:8] + "..." + key[len(key)-4:]
}

func printCommandError(err error) int {
	var apiErr *apiError
	if errors.As(err, &apiErr) {
		printAPIError(apiErr)
		return 1
	}
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		fmt.Fprintf(os.Stderr, "network error: %v\n", urlErr.Err)
		fmt.Fprintln(os.Stderr, "hint: check internet connection, DNS, or UNCUT_BASE_URL")
		return 1
	}
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	return 1
}

func printAPIError(err *apiError) {
	title := friendlyAPIErrorTitle(err)
	fmt.Fprintf(os.Stderr, "error: %s\n", title)
	if err.Message != "" && err.Message != title {
		fmt.Fprintf(os.Stderr, "message: %s\n", err.Message)
	}
	if err.Code != "" {
		fmt.Fprintf(os.Stderr, "code: %s\n", err.Code)
	}
	if err.Status != 0 {
		fmt.Fprintf(os.Stderr, "http: %d\n", err.Status)
	}
	printAPIDetails(err.Details)
	if hint := apiErrorHint(err); hint != "" {
		fmt.Fprintf(os.Stderr, "hint: %s\n", hint)
	}
}

func friendlyAPIErrorTitle(err *apiError) string {
	switch err.Code {
	case "unauthorized":
		return "API key is missing, invalid, revoked, or account is inactive"
	case "not_found":
		return "resource was not found or does not belong to this account"
	case "validation_failed":
		return "request validation failed"
	case "insufficient_balance":
		return "wallet balance is too low"
	case "insufficient_card_balance":
		return "card balance is too low"
	case "invalid_phone":
		return "phone number is invalid"
	case "invalid_bin":
		return "BIN is unknown or inactive"
	case "enable_3ds_unsupported":
		return "3DS is not supported for this BIN"
	case "unsupported_currency":
		return "currency is not supported for this BIN"
	case "card_not_active":
		return "card is not active"
	case "card_not_frozen":
		return "card is not frozen"
	case "provider_error":
		return "card provider rejected or failed the request"
	case "provider_unavailable":
		return "card provider is temporarily unavailable"
	case "card_issue_unavailable":
		return "card issuing is temporarily unavailable"
	case "exchange_rate_unavailable":
		return "exchange rate is temporarily unavailable"
	case "deposit_address_unavailable":
		return "deposit address is temporarily unavailable"
	}
	if err.Status == http.StatusTooManyRequests {
		return "rate limit exceeded"
	}
	if err.Message != "" {
		return err.Message
	}
	return err.Error()
}

func apiErrorHint(err *apiError) string {
	switch err.Code {
	case "unauthorized":
		return "run `uncut login` again or set UNCUT_API_KEY"
	case "not_found":
		return "run `uncut cards` and copy the current card id"
	case "validation_failed":
		return "check required flags and value formats; use `uncut` for command examples"
	case "insufficient_balance":
		return "run `uncut balance`; create with `--topup 0` or fund the wallet with `uncut deposit`"
	case "insufficient_card_balance":
		return "run `uncut card <card_id>` and choose a smaller withdraw amount"
	case "invalid_phone":
		return "use E.164 format, for example +10000000000"
	case "invalid_bin":
		return "run `uncut bins` or `uncut new` and copy a current BIN id"
	case "enable_3ds_unsupported":
		return "choose a BIN where the `wallet` column is yes, or remove `--3ds`"
	case "unsupported_currency":
		return "run `uncut bins` and choose one of the currencies listed for that BIN"
	case "card_not_active":
		return "run `uncut unfreeze <card_id>` first, or choose an active card"
	case "card_not_frozen":
		return "run `uncut freeze <card_id>` first, or choose a frozen card"
	case "provider_error":
		return "retry later; if this was delete, the card was not deleted"
	case "provider_unavailable", "card_issue_unavailable", "exchange_rate_unavailable", "deposit_address_unavailable":
		return "retry later"
	}
	if err.Status == http.StatusTooManyRequests {
		return "wait a minute and retry"
	}
	return ""
}

func printAPIDetails(details json.RawMessage) {
	if len(bytes.TrimSpace(details)) == 0 {
		return
	}
	var value any
	dec := json.NewDecoder(bytes.NewReader(details))
	dec.UseNumber()
	if err := dec.Decode(&value); err != nil {
		fmt.Fprintf(os.Stderr, "details: %s\n", strings.TrimSpace(string(details)))
		return
	}
	fmt.Fprintln(os.Stderr, "details:")
	printDetailValue(value, "  ")
}

func printDetailValue(value any, indent string) {
	switch typed := value.(type) {
	case map[string]any:
		for key, val := range typed {
			switch val.(type) {
			case map[string]any, []any:
				fmt.Fprintf(os.Stderr, "%s%s:\n", indent, key)
				printDetailValue(val, indent+"  ")
			default:
				fmt.Fprintf(os.Stderr, "%s%s: %v\n", indent, key, val)
			}
		}
	case []any:
		for _, item := range typed {
			fmt.Fprintf(os.Stderr, "%s- %v\n", indent, item)
		}
	default:
		fmt.Fprintf(os.Stderr, "%s%v\n", indent, typed)
	}
}

func printRawJSON(raw []byte) {
	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 {
		fmt.Println("{}")
		return
	}
	fmt.Println(string(raw))
}

func printJSON(value any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(value)
}
