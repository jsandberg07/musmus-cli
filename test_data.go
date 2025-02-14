package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jsandberg07/clitest/internal/database"
)

func (cfg *Config) testData() error {
	// add admin
	fmt.Println("* Creating admin...")
	err := cfg.createAdmin()
	if err != nil {
		return err
	}

	// settings are done
	// add default positions
	err = addTestPositions(cfg)
	if err != nil {
		return err
	}

	// add bunk investigators
	err = addTestInvestigators(cfg)
	if err != nil {
		return err
	}

	// add protocols
	err = addTestProtocols(cfg)
	if err != nil {
		return err
	}

	// add investigators to protocols
	err = addTestInvestigatorToProtocol(cfg)
	if err != nil {
		return err
	}

	// add strains
	err = addTestStrains(cfg)
	if err != nil {
		return err
	}

	// add cage cards
	err = addTestCageCards(cfg)
	if err != nil {
		return err
	}

	// now activate them
	err = activateTestCageCards(cfg)
	if err != nil {
		return err
	}

	// deactivate a few
	err = deactivateTestCageCards(cfg)
	if err != nil {
		return err
	}

	// get active cage cards
	err = getTestActiveCageCards(cfg)
	if err != nil {
		return err
	}

	// add a reminder for today
	err = addTestReminder(cfg)
	if err != nil {
		return err
	}

	// add an order expected today
	err = addTestOrder(cfg)
	if err != nil {
		return err
	}

	return nil
}

func addTestPositions(cfg *Config) error {
	fmt.Println("* Creating test positions...")

	posPI := database.CreatePositionParams{
		Title:             "PI",
		CanActivate:       true,
		CanDeactivate:     true,
		CanAddOrders:      true,
		CanReceiveOrders:  true,
		CanQuery:          true,
		CanChangeProtocol: true,
		CanAddStaff:       true,
		CanAddReminders:   true,
	}
	posRes := database.CreatePositionParams{
		Title:         "Researcher",
		CanActivate:   true,
		CanDeactivate: true,
		CanAddOrders:  true,
		CanQuery:      true,
	}
	posAss := database.CreatePositionParams{
		Title:    "Assistant",
		CanQuery: true,
	}
	deniedPos := database.CreatePositionParams{
		Title: "Denied",
	}
	positions := []database.CreatePositionParams{posPI, posRes, posAss, deniedPos}

	for i, position := range positions {
		cPos, err := cfg.db.CreatePosition(context.Background(), position)
		if err != nil {
			fmt.Printf("Error creating position %v -- %s\n", i, position.Title)
			return err
		}
		if verbose {
			fmt.Print(cPos)
		}
	}

	return nil
}

func addTestInvestigators(cfg *Config) error {
	fmt.Println("* Creating test investigators...")
	// fake names
	// josh england
	// sharon thornton
	// sonya ball

	PIpos, err := cfg.db.GetPositionByTitle(context.Background(), "PI")
	if err != nil {
		return err
	}
	invPI := database.CreateInvestigatorParams{
		IName:    "Josh England",
		Nickname: sql.NullString{Valid: true, String: "Johnny Boi"},
		Email:    sql.NullString{Valid: true, String: "je@test.com"},
		Position: PIpos.ID,
	}

	ResPos, err := cfg.db.GetPositionByTitle(context.Background(), "Researcher")
	if err != nil {
		return err
	}
	invRes := database.CreateInvestigatorParams{
		IName:    "Sharon Thornton",
		Email:    sql.NullString{Valid: true, String: "st@test.com"},
		Position: ResPos.ID,
	}

	AssPos, err := cfg.db.GetPositionByTitle(context.Background(), "Assistant")
	if err != nil {
		return err
	}
	invAss := database.CreateInvestigatorParams{
		IName:    "Sonya Ball",
		Nickname: sql.NullString{Valid: true, String: "Coco"},
		Position: AssPos.ID,
	}

	DenPos, err := cfg.db.GetPositionByTitle(context.Background(), "Denied")
	if err != nil {
		return err
	}
	denInv := database.CreateInvestigatorParams{
		IName:    "Daniel Denial",
		Nickname: sql.NullString{Valid: true, String: "Dan"},
		Position: DenPos.ID,
	}
	investigators := []database.CreateInvestigatorParams{invPI, invRes, invAss, denInv}
	for i, investigator := range investigators {
		ci, err := cfg.db.CreateInvestigator(context.Background(), investigator)
		if err != nil {
			fmt.Printf("Error creating investigator %v -- %s\n", i, investigator.IName)
			return err
		}
		if verbose {
			fmt.Println(ci)
		}
	}

	return nil
}

func addTestProtocols(cfg *Config) error {
	fmt.Println("* Creating test protocols...")

	PI, err := cfg.db.GetInvestigatorByName(context.Background(), "Josh England")
	if err != nil {
		return err
	}
	if len(PI) > 1 {
		return errors.New("vague PI name")
	}

	prot1 := database.CreateProtocolParams{
		PNumber:             "12-24-32",
		PrimaryInvestigator: PI[0].ID,
		Title:               "IRS-3 and metabolism",
		Allocated:           200,
		Balance:             50,
		ExpirationDate:      time.Now().AddDate(3, 1, 1),
	}
	prot2 := database.CreateProtocolParams{
		PNumber:             "18-12-16",
		PrimaryInvestigator: PI[0].ID,
		Title:               "Extended release coagulants",
		Allocated:           115,
		Balance:             110,
		ExpirationDate:      time.Now().AddDate(0, 2, 0),
	}
	protocols := []database.CreateProtocolParams{prot1, prot2}
	for i, protocol := range protocols {
		cp, err := cfg.db.CreateProtocol(context.Background(), protocol)
		if err != nil {
			fmt.Printf("Error creating protocol %v -- %s\n", i, cp.Title)
			return err
		}
		if verbose {
			fmt.Println(cp)
		}
	}

	return nil

}

func addTestInvestigatorToProtocol(cfg *Config) error {
	// get everybody, add them at once
	fmt.Println("* Adding investigators to protocols...")

	investigatorNames := []string{"Johnny Boi", "Sharon Thornton", "Sonya Ball"}
	investigators := []database.Investigator{}
	for i, name := range investigatorNames {
		tID, err := cfg.db.GetInvestigatorByName(context.Background(), name)
		if err != nil {
			fmt.Printf("Error getting investigator #%v -- %v\n", i+1, name)
			return err
		}
		if len(tID) == 0 {
			fmt.Println("Investigator not found...")
			continue
		}
		if len(tID) > 1 {
			fmt.Printf("Error getting investigator #%v -- %v\n", i+1, name)
			return errors.New("vague investigator name")
		}

		investigators = append(investigators, tID[0])
	}

	protocols, err := cfg.db.GetProtocols(context.Background())
	if err != nil {
		fmt.Println("Error getting protocols")
		return err
	}
	if verbose {
		fmt.Println(protocols)
	}

	for _, protocol := range protocols {
		for _, investigator := range investigators {
			addInvToProt := database.AddInvestigatorToProtocolParams{
				InvestigatorID: investigator.ID,
				ProtocolID:     protocol.ID,
			}
			_, err := cfg.db.AddInvestigatorToProtocol(context.Background(), addInvToProt)
			if err != nil {
				fmt.Printf("Error adding test %s to test %s\n", investigator.IName, protocol.PNumber)
				return err
			}
		}
	}

	// now remove somebody from one >:^3
	// Sonya isn't on uhh something it might not be consistent
	rmvInvFromProt := database.RemoveInvestigatorFromProtocolParams{
		InvestigatorID: investigators[2].ID,
		ProtocolID:     protocols[0].ID,
	}
	err = cfg.db.RemoveInvestigatorFromProtocol(context.Background(), rmvInvFromProt)
	if err != nil {
		fmt.Println("Error removing test investigator from test protocol")
		return err
	}
	return nil

}

func addTestStrains(cfg *Config) error {
	fmt.Println("* Adding test strains...")
	asC57 := database.AddStrainParams{
		SName:      "C57BL6/J",
		Vendor:     "Jax",
		VendorCode: "000664",
	}
	asBALB := database.AddStrainParams{
		SName:      "BALB/cJ",
		Vendor:     "Jax",
		VendorCode: "000651",
	}
	asCD1 := database.AddStrainParams{
		SName:      "CD-1",
		Vendor:     "CRL",
		VendorCode: "022",
	}
	strains := []database.AddStrainParams{asC57, asBALB, asCD1}
	for i, strain := range strains {
		ts, err := cfg.db.AddStrain(context.Background(), strain)
		if err != nil {
			fmt.Printf("Error adding strain %v -- %s\n", i, strain.SName)
			return err
		}
		if verbose {
			fmt.Println(ts)
		}
	}

	return nil
}

func addTestCageCards(cfg *Config) error {
	fmt.Println("* Adding test cage cards...")
	invest, err := cfg.db.GetInvestigatorByName(context.Background(), "Sharon Thornton")
	if err != nil {
		fmt.Println("Error getting investigator for cage card")
	}
	if len(invest) > 1 {
		fmt.Println("Error getting investigator for cage card")
		return errors.New("vague investigator name")
	}
	prot1, err := cfg.db.GetProtocolByNumber(context.Background(), "12-24-32")
	if err != nil {
		fmt.Println("Error getting protocol 1")
		return err
	}
	prot2, err := cfg.db.GetProtocolByNumber(context.Background(), "18-12-16")
	if err != nil {
		fmt.Println("Error getting protocol 2")
		return err
	}

	ccStart := int32(100)
	ccEnd := int32(120)
	for i := ccStart; i < ccEnd; i++ {
		aCC := database.AddCageCardParams{
			CcID:           i,
			InvestigatorID: invest[0].ID,
			ProtocolID:     prot1.ID,
		}
		cCC, err := cfg.db.AddCageCard(context.Background(), aCC)
		if err != nil {
			fmt.Printf("Error adding cage card %v", i)
			return err
		}
		if verbose {
			fmt.Println(cCC)
		}
	}

	ccStart = 121
	ccEnd = 140
	for i := ccStart; i < ccEnd; i++ {
		aCC := database.AddCageCardParams{
			CcID:           i,
			InvestigatorID: invest[0].ID,
			ProtocolID:     prot2.ID,
		}
		cCC, err := cfg.db.AddCageCard(context.Background(), aCC)
		if err != nil {
			fmt.Printf("Error adding cage card %v\n", i)
			return err
		}
		if verbose {
			fmt.Println(cCC)
		}
	}

	return nil

}

func activateTestCageCards(cfg *Config) error {
	fmt.Println("* Activating test cage cards")
	lastWeek := normalizeDate(time.Now().AddDate(0, 0, -7))
	activatedBy, err := cfg.db.GetInvestigatorByName(context.Background(), "Sonya Ball")
	if err != nil {
		fmt.Println("Error getting investigator for activation")
		return err
	}
	if len(activatedBy) > 1 {
		fmt.Println("Error getting investigator for activation")
		return errors.New("vague investigator name")
	}

	strain1, err := cfg.db.GetStrainByName(context.Background(), "CD-1")
	if err != nil {
		fmt.Println("Error gettting strain 1")
		return err
	}

	/* just for when i want more varried test data
	strain2, err := cfg.db.GetStrainByName(context.Background(), "000664")
	if err != nil {
		fmt.Println("Error getting strain 2")
		return err
	}
	*/

	cardsToActivate := []int32{100, 102, 103, 104, 108, 109, 111, 121, 123, 134, 139}
	for i, cardID := range cardsToActivate {
		aCC := database.TrueActivateCageCardParams{
			CcID:        cardID,
			ActivatedOn: sql.NullTime{Valid: true, Time: lastWeek},
			ActivatedBy: uuid.NullUUID{Valid: true, UUID: activatedBy[0].ID},
			Strain:      uuid.NullUUID{Valid: true, UUID: strain1.ID},
		}
		cc, err := cfg.db.TrueActivateCageCard(context.Background(), aCC)
		if err != nil {
			fmt.Printf("Error activating cage card %v -- %v", i, cardID)
			return err
		}
		if verbose {
			fmt.Printf("CC %v activated by %s\n", cc.CcID, activatedBy[0].ID)
		}
	}

	return nil
}

func deactivateTestCageCards(cfg *Config) error {
	fmt.Println("* Deactivating test cage cards")
	yesterday := normalizeDate(time.Now().AddDate(0, 0, -1))
	deactivatedBy, err := cfg.db.GetInvestigatorByName(context.Background(), "Sonya Ball")
	if err != nil {
		fmt.Println("Error getting investigator for activation")
		return err
	}
	if len(deactivatedBy) > 1 {
		fmt.Println("Error getting investigator for activation")
		return errors.New("vague investigator name")
	}

	cardsToDeactivate := []int32{108, 109, 111}
	for i, cardID := range cardsToDeactivate {
		dCC := database.DeactivateCageCardParams{
			CcID:          cardID,
			DeactivatedOn: sql.NullTime{Valid: true, Time: yesterday},
			DeactivatedBy: uuid.NullUUID{Valid: true, UUID: deactivatedBy[0].ID},
		}
		cc, err := cfg.db.DeactivateCageCard(context.Background(), dCC)
		if err != nil {
			fmt.Printf("Error deactivating cage card %v -- %v", i, cardID)
			return err
		}
		if verbose {
			fmt.Printf("CC %v activated by %s\n", cc.CcID, deactivatedBy[0].ID)
		}
	}

	return nil
}

func getTestActiveCageCards(cfg *Config) error {
	fmt.Println("* Getting all active test cage cards")
	activeCards, err := cfg.db.GetAllActiveCageCards(context.Background())
	if err != nil {
		fmt.Println("Error getting all active test cage cards")
		return err
	}
	for _, cc := range activeCards {
		if verbose {
			fmt.Println(cc)
		}
	}

	return nil
}

func addTestReminder(cfg *Config) error {
	// actually add two and make sure they dont both show up
	// better for testing login later
	names := []string{"Sharon Thornton", "Johnny Boi"}
	ccID := 100

	for i, name := range names {
		inv, err := cfg.db.GetInvestigatorByName(context.Background(), name)
		if err != nil {
			return err
		}
		if len(inv) == 0 {
			return errors.New("no investigators found")
		}
		if len(inv) > 1 {
			return errors.New("vague investigator while adding reminder")
		}
		arp := database.AddReminderParams{
			RDate:          normalizeDate(time.Now()),
			RCcID:          int32(ccID),
			InvestigatorID: inv[0].ID,
			Note:           "Test Reminder",
		}
		reminder, err := cfg.db.AddReminder(context.Background(), arp)
		if err != nil {
			fmt.Printf("Error adding reminder %v -- %s", i, name)
			return err
		}
		if verbose {
			fmt.Println(reminder)
		}
		ccID++

	}
	return nil
}

func addTestOrder(cfg *Config) error {
	inv, err := cfg.db.GetInvestigatorByName(context.Background(), "Johnny Boi")
	if err != nil {
		return err
	}
	if len(inv) == 0 {
		return errors.New("no investigators found")
	}
	if len(inv) > 1 {
		return errors.New("vague investigator while adding reminder")
	}
	pro, err := cfg.db.GetProtocolByNumber(context.Background(), "12-24-32")
	if err != nil {
		return err
	}
	strain, err := cfg.db.GetStrainByName(context.Background(), "000664")
	if err != nil {
		return err
	}

	cnop := database.CreateNewOrderParams{
		OrderNumber:    "B1234",
		ExpectedDate:   normalizeDate(time.Now()),
		ProtocolID:     pro.ID,
		InvestigatorID: inv[0].ID,
		StrainID:       strain.ID,
		Note:           sql.NullString{Valid: true, String: "Test Order"},
	}

	order, err := cfg.db.CreateNewOrder(context.Background(), cnop)
	if err != nil {
		return err
	}
	if verbose {
		fmt.Println(order)
	}

	return nil
}
