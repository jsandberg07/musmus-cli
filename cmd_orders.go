package main

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jsandberg07/clitest/internal/database"
)

func getAddOrderCmd() Command {
	addOrderFlags := make(map[string]Flag)
	addOrderCmd := Command{
		name:        "add",
		description: "Used for adding orders",
		function:    addOrderFunction,
		flags:       addOrderFlags,
	}

	return addOrderCmd
}

// prompts so just save, exit, print
func getAddOrderFlags() map[string]Flag {
	addOrderFlags := make(map[string]Flag)
	saveFlag := Flag{
		symbol:      "save",
		description: "Saves the current order",
		takesValue:  false,
	}
	addOrderFlags[saveFlag.symbol] = saveFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits without saving",
		takesValue:  false,
	}
	addOrderFlags[exitFlag.symbol] = exitFlag

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints flags available for current command",
		takesValue:  false,
	}
	addOrderFlags[helpFlag.symbol] = helpFlag

	return addOrderFlags

}

// look into removing the args thing, might have to stay
func addOrderFunction(cfg *Config, args []Argument) error {
	// get flags
	flags := getAddOrderFlags()

	// set defaults
	exit := false

	// the reader
	reader := bufio.NewReader(os.Stdin)

	orderNumber, err := getStringPrompt(cfg, "Enter order number", checkIfOrderNumberUnique)
	if err != nil {
		return err
	}
	if orderNumber == "" {
		fmt.Println("Exiting...")
	}

	date, err := getDatePrompt("Enter expected date")
	if err != nil {
		return err
	}
	nilDate := time.Time{}
	if date == nilDate {
		fmt.Println("Exiting...")
		return nil
	}

	protocol, err := getStructPrompt(cfg, "Enter protocol for order", getProtocolStruct)
	if err != nil {
		return err
	}
	nilProtocol := database.Protocol{}
	if protocol == nilProtocol {
		fmt.Println("Exiting...")
		return nil
	}

	investigator, err := getStructPrompt(cfg, "Enter investigator receiving order", getInvestigatorStruct)
	if err != nil {
		return err
	}
	nilInvestigator := database.Investigator{}
	if investigator == nilInvestigator {
		fmt.Println("Exiting...")
		return nil
	}

	strain, err := getStructPrompt(cfg, "Enter strain of order", getStrainStruct)
	if err != nil {
		return err
	}
	nilStrain := database.Strain{}
	if strain == nilStrain {
		fmt.Println("Exiting...")
		return nil
	}

	note, err := getStringPrompt(cfg, "Optionally enter a note. Will be applied to all cage cards from order", checkFuncNil)
	if err != nil {
		return err
	}
	var dbNote sql.NullString
	if note == "" {
		dbNote.Valid = false
	} else {
		dbNote.Valid = true
		dbNote.String = note
	}

	cnoParam := database.CreateNewOrderParams{
		OrderNumber:    orderNumber,
		ExpectedDate:   date,
		ProtocolID:     protocol.ID,
		InvestigatorID: investigator.ID,
		StrainID:       strain.ID,
		Note:           dbNote,
	}

	// working here: create the params, set the flags, remember to check if note valid or not
	// thanks for the reminder. vacation time.
	fmt.Println("Order will be created with the following settings:")
	printNewOrder(&cnoParam, &protocol, &investigator, &strain)
	fmt.Println("Enter 'save' or 'exit'")
	// da loop
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input string: %s", err)
			os.Exit(1)
		}

		inputs, err := readSubcommandInput(text)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// do weird behavior here

		// but normal loop now
		args, err := parseArguments(flags, inputs)
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, arg := range args {
			switch arg.flag {
			case "save":
				fmt.Println("Saving...")
				order, err := cfg.db.CreateNewOrder(context.Background(), cnoParam)
				if err != nil {
					fmt.Println("Error adding order to db")
					return err
				}
				if verbose {
					fmt.Println(order)
				}
				exit = true

			case "exit":
				fmt.Println("Exiting...")
				exit = true

			case "help":
				cmdHelp(flags)
			default:
				fmt.Printf("Oops a fake flag snuck in: %s\n", arg.flag)
			}
		}

		if exit {
			break
		}

	}

	return nil
}

func checkIfOrderNumberUnique(cfg *Config, input string) error {
	_, err := cfg.db.GetOrderByNumber(context.Background(), input)
	if err != nil && err.Error() != "sql: no rows in result set" {
		// any other error
		return err
	}

	if err == nil {
		// is not unique
		return errors.New("order number is not unique. Please try again")
	}

	// is unique
	return nil

}

func printNewOrder(o *database.CreateNewOrderParams, p *database.Protocol, i *database.Investigator, s *database.Strain) {
	fmt.Printf("* Number - %s\n", o.OrderNumber)
	fmt.Printf("* Date - %v\n", o.ExpectedDate)
	fmt.Printf("* Protocol - %v\n", p.PNumber)
	fmt.Printf("* Investigator - %s\n", i.IName)
	fmt.Printf("* Strain - %s\n", s.SName)
	if o.Note.Valid {
		fmt.Printf("* Note - %s\n", o.Note.String)
	}
}

// TODO: merge with above. Find some kind of normalization, even though they both dont have an order #
func printUpdateOrder(o *database.UpdateOrderParams, p *database.Protocol, i *database.Investigator, s *database.Strain) {
	fmt.Printf("* Date - %v\n", o.ExpectedDate)
	fmt.Printf("* Protocol - %v\n", p.PNumber)
	fmt.Printf("* Investigator - %s\n", i.IName)
	fmt.Printf("* Strain - %s\n", s.SName)
	if o.Note.Valid {
		fmt.Printf("* Note - %s\n", o.Note.String)
	}
}

func getEditOrderCmd() Command {
	editOrderFlags := make(map[string]Flag)
	EditOrderCmd := Command{
		name:        "edit",
		description: "Used for editing existing orders",
		function:    editOrderFunction,
		flags:       editOrderFlags,
	}

	return EditOrderCmd
}

// can't change the number that's too much
// expected [d]ate, [i]nvestigator, [s]train, [n]ote, dont bother with unreceiving too much work for a fake program
// save print
func getEditOrderFlags() map[string]Flag {
	editOrderFlags := make(map[string]Flag)

	dFlag := Flag{
		symbol:      "d",
		description: "Sets expected date",
		takesValue:  true,
	}
	editOrderFlags["-"+dFlag.symbol] = dFlag

	nFlag := Flag{
		symbol:      "n",
		description: "Sets order note. Enter 'x' to blank out the note",
		takesValue:  true,
	}
	editOrderFlags["-"+nFlag.symbol] = nFlag

	iFlag := Flag{
		symbol:      "i",
		description: "Sets who the order is for",
		takesValue:  true,
	}
	editOrderFlags["-"+iFlag.symbol] = iFlag

	sFlag := Flag{
		symbol:      "s",
		description: "Sets order strain",
		takesValue:  true,
	}
	editOrderFlags["-"+sFlag.symbol] = sFlag

	pFlag := Flag{
		symbol:      "p",
		description: "Sets order protocol",
		takesValue:  true,
	}
	editOrderFlags["-"+pFlag.symbol] = pFlag

	// ect as needed or remove the "-"+ for longer ones

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints available flags",
		takesValue:  false,
	}
	editOrderFlags[helpFlag.symbol] = helpFlag

	saveFlag := Flag{
		symbol:      "save",
		description: "Saves the updated order",
		takesValue:  false,
	}
	editOrderFlags[saveFlag.symbol] = saveFlag

	printFlag := Flag{
		symbol:      "print",
		description: "Prints current order parameters for review",
		takesValue:  false,
	}
	editOrderFlags[printFlag.symbol] = printFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits without saving",
		takesValue:  false,
	}
	editOrderFlags[exitFlag.symbol] = exitFlag

	return editOrderFlags

}

// look into removing the args thing, might have to stay
// ask for an order number, load params, set flags
func editOrderFunction(cfg *Config, args []Argument) error {
	// get flags
	flags := getEditOrderFlags()

	order, err := getStructPrompt(cfg, "Enter order number", getOrderStruct)
	if err != nil {
		return err
	}
	nilOrder := database.Order{}
	if order == nilOrder {
		fmt.Println("Exiting...")
		return nil
	}

	// set defaults
	exit := false
	reviewed := Reviewed{
		Printed:     false,
		ChangesMade: false,
	}
	uoParams := database.UpdateOrderParams{
		ID:             order.ID,
		ExpectedDate:   order.ExpectedDate,
		InvestigatorID: order.InvestigatorID,
		StrainID:       order.StrainID,
		Note:           order.Note,
	}

	// the reader
	reader := bufio.NewReader(os.Stdin)

	// da loop
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input string: %s", err)
			os.Exit(1)
		}

		inputs, err := readSubcommandInput(text)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// do weird behavior here
		if reviewed.ChangesMade {
			reviewed.Printed = false
		}

		// but normal loop now
		args, err := parseArguments(flags, inputs)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// d, n, i, s, save, print, exit
		for _, arg := range args {
			switch arg.flag {
			case "exit":
				fmt.Println("Exiting...")
				exit = true

			case "save":
				fmt.Println("Saving...")
				if reviewed.ChangesMade || !reviewed.Printed {
					fmt.Println("Updating with the following params:")
					protocol, err := cfg.db.GetProtocolByID(context.Background(), uoParams.ProtocolID)
					if err != nil {
						return err
					}
					strain, err := cfg.db.GetStrainByID(context.Background(), uoParams.StrainID)
					if err != nil {
						return err
					}
					investigator, err := cfg.db.GetInvestigatorByID(context.Background(), uoParams.InvestigatorID)
					if err != nil {
						return err
					}
					printUpdateOrder(&uoParams, &protocol, &investigator, &strain)
				}
				err := cfg.db.UpdateOrder(context.Background(), uoParams)
				if err != nil {
					return err
				}
				exit = true

			case "print":
				protocol, err := cfg.db.GetProtocolByID(context.Background(), uoParams.ProtocolID)
				if err != nil {
					return err
				}
				strain, err := cfg.db.GetStrainByID(context.Background(), uoParams.StrainID)
				if err != nil {
					return err
				}
				investigator, err := cfg.db.GetInvestigatorByID(context.Background(), uoParams.InvestigatorID)
				if err != nil {
					return err
				}
				printUpdateOrder(&uoParams, &protocol, &investigator, &strain)
				reviewed.Printed = true
				reviewed.ChangesMade = false

			case "-n":
				if arg.value == "x" {
					uoParams.Note = sql.NullString{Valid: false}
				} else {
					uoParams.Note = sql.NullString{Valid: true, String: arg.value}
				}
				reviewed.ChangesMade = true

			case "-d":
				date, err := parseDate(arg.value)
				if err != nil {
					fmt.Println(err)
					break
				}
				uoParams.ExpectedDate = date
				reviewed.ChangesMade = true

			case "-i":
				investigator, err := getInvestigatorByFlag(cfg, arg.value)
				if err != nil {
					return err
				}
				uoParams.InvestigatorID = investigator.ID
				reviewed.ChangesMade = true

			case "-s":
				strain, err := getStrainByFlag(cfg, arg.value)
				if err != nil {
					return err
				}
				uoParams.StrainID = strain.ID
				reviewed.ChangesMade = true

			case "-p":
				protocol, err := getProtocolByFlag(cfg, arg.value)
				if err != nil {
					return err
				}
				uoParams.ProtocolID = protocol.ID
				reviewed.ChangesMade = true

			default:
				fmt.Printf("Oops a fake flag snuck in: %s\n", arg.flag)
			}
		}

		if exit {
			break
		}

	}

	return nil
}

func getOrderStruct(cfg *Config, input string) (database.Order, error) {
	order, err := cfg.db.GetOrderByNumber(context.Background(), input)
	if err != nil && err.Error() == "sql: no rows in result set" {
		// not found
		return database.Order{}, errors.New("no order found")
	}
	if err != nil {
		// any other error
		return database.Order{}, err
	}
	return order, nil
}

func getReceiveOrderCmd() Command {
	receiveOrderFlags := make(map[string]Flag)
	ReceiveOrderCmd := Command{
		name:        "receive",
		description: "Used for receiving orders",
		function:    receiveOrderFunction,
		flags:       receiveOrderFlags,
	}

	return ReceiveOrderCmd
}

func getReceiveOrderFlags() map[string]Flag {
	receiveOrderFlags := make(map[string]Flag)

	receiveFlag := Flag{
		symbol:      "receive",
		description: "Receives the order with the current parameters",
		takesValue:  false,
	}
	receiveOrderFlags[receiveFlag.symbol] = receiveFlag

	exitFlag := Flag{
		symbol:      "exit",
		description: "Exits without receiving the current order",
		takesValue:  false,
	}
	receiveOrderFlags[exitFlag.symbol] = exitFlag

	printFlag := Flag{
		symbol:      "print",
		description: "Review the order params before receiving it",
		takesValue:  false,
	}
	receiveOrderFlags[printFlag.symbol] = printFlag

	helpFlag := Flag{
		symbol:      "help",
		description: "Prints available flags",
		takesValue:  false,
	}
	receiveOrderFlags[helpFlag.symbol] = helpFlag

	return receiveOrderFlags

}

// ask for an order number
// ask for a date to receive (blank for today)
// cage card range (start to finish)
// as "this many cards will be added on this day this order ok"
// then DO IT
// look into removing the args thing, might have to stay
func receiveOrderFunction(cfg *Config, args []Argument) error {
	// flags just for saving and exiting prompt everything else
	flags := getReceiveOrderFlags()

	// set defaults
	exit := false

	// the reader
	reader := bufio.NewReader(os.Stdin)

	order, err := getStructPrompt(cfg, "Enter order number for order to receive", getOrderStruct)
	if err != nil {
		return err
	}
	nilOrder := database.Order{}
	if order == nilOrder {
		fmt.Println("Exiting...")
		return nil
	}

	start, err := getIntPrompt("Enter start of cage card range")
	if err != nil {
		return err
	}
	if start == -1 {
		fmt.Println("Exiting...")
		return nil
	}

	end, err := getIntPrompt("Enter end of cage card range")
	if err != nil {
		return err
	}
	if end == -1 {
		fmt.Println("Exiting...")
		return nil
	}

	if start > end {
		fmt.Println("Start larger than end. Swapping...")
		temp := end
		end = start
		start = temp
	}

	// da loop
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input string: %s", err)
			os.Exit(1)
		}

		inputs, err := readSubcommandInput(text)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// do weird behavior here

		// but normal loop now
		args, err := parseArguments(flags, inputs)
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, arg := range args {
			switch arg.flag {
			case "print":
				protocol, err := cfg.db.GetProtocolByID(context.Background(), order.ProtocolID)
				if err != nil {
					return err
				}
				investigator, err := cfg.db.GetInvestigatorByID(context.Background(), order.InvestigatorID)
				if err != nil {
					return err
				}
				strain, err := cfg.db.GetStrainByID(context.Background(), order.StrainID)
				if err != nil {
					return err
				}
				printReceiveOrder(start, end, &order, &protocol, &investigator, &strain)

			case "exit":
				fmt.Println("Exiting...")
				exit = true

			case "receive":
				fmt.Println("Receiving...")
				err := receiveOrder(cfg, start, end, &order)
				if err != nil {
					fmt.Println("Error receiving order")
					return err
				}

				exit = true

			case "help":
				cmdHelp(flags)

			default:
				fmt.Printf("Oops a fake flag snuck in: %s\n", arg.flag)
			}
		}

		if exit {
			break
		}

	}

	return nil
}

func printReceiveOrder(start, end int, o *database.Order, p *database.Protocol, i *database.Investigator, s *database.Strain) {
	fmt.Printf("* Order Number - %v\n", o.OrderNumber)
	fmt.Printf("* Date - %v\n", o.ExpectedDate)
	fmt.Printf("* Protocol - %v\n", p.PNumber)
	fmt.Printf("* Investigator - %v\n", i.IName)
	fmt.Printf("* Strain - %v\n", s.SName)
	fmt.Printf("* Note - %v\n", o.Note)
	fmt.Printf("* CC start - %v\n", start)
	fmt.Printf("* CC  end  - %v\n", end)
	fmt.Printf("* Total CC - %v\n", end-start+1)
}

// you should take it in steps but instead go whole hog. what i lack in it working i make
// up with being so large you wont notice edge cases kek
// loop start to end
// add cage card, activation date, order number, ect
// mark order as received at the end if it works ok
func receiveOrder(cfg *Config, start, end int, o *database.Order) error {
	// check if cage cards are already in db
	cageCards, err := cfg.db.GetCageCardsRange(context.Background(), database.GetCageCardsRangeParams{CcID: int32(start), CcID_2: int32(end)})
	if err != nil {
		fmt.Println("Could not check DB for cage cards")
		return err
	}
	if len(cageCards) != 0 {
		return errors.New("cage cards in range already added to DB. Please check start, end and try again")
	}

	// create the params
	rccParams := database.ReceiveCageCardParams{
		ProtocolID:     o.ProtocolID,
		ActivatedOn:    sql.NullTime{Valid: true, Time: o.ExpectedDate},
		InvestigatorID: o.InvestigatorID,
		Strain:         uuid.NullUUID{Valid: true, UUID: o.StrainID},
		Notes:          o.Note,
		ActivatedBy:    uuid.NullUUID{Valid: true, UUID: cfg.loggedInInvestigator.ID},
	}

	ccActivated := 0

	// <= is intentional, cage card range is inclusive
	for i := start; i <= end; i++ {
		rccParams.CcID = int32(i)
		cc, err := cfg.db.ReceiveCageCard(context.Background(), rccParams)
		if err != nil {
			fmt.Printf("Err receiving CC %v\n", i)
			fmt.Println(err)
			continue
		}
		if verbose {
			fmt.Println(cc)
		}
		ccActivated++
	}

	fmt.Printf("%v CC activated", ccActivated)
	// mark order as received
	order, err := cfg.db.MarkOrderReceived(context.Background(), o.ID)
	if err != nil {
		fmt.Println("Couldn't mark order as received")
		return err
	}
	if verbose {
		fmt.Println(order)
	}

	return nil
}

// "duplicate key value violates unique constraint"
// already exists
func getTodaysOrders(cfg *Config) error {
	gueoParams := database.GetUserExpectedOrdersParams{
		ExpectedDate:   normalizeDate(time.Now()),
		InvestigatorID: cfg.loggedInInvestigator.ID,
	}
	orders, err := cfg.db.GetUserExpectedOrders(context.Background(), gueoParams)
	if err != nil {
		return err
	}
	if len(orders) == 0 {
		fmt.Println("No orders expected today")
		return nil
	}

	fmt.Println("Orders expected today: ")
	for i, order := range orders {
		if order.Note.Valid {
			fmt.Printf("* %v -- %s -- %s\n", i+1, order.OrderNumber, order.Note.String)
		} else {
			fmt.Printf("* %v -- %s\n", i+1, order.OrderNumber)
		}

	}

	return nil
}
