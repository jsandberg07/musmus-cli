package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jsandberg07/clitest/internal/database"
)

func (cfg *Config) testData() error {
	// settings are done
	// add default positions
	err := addTestPositions(cfg)
	if err != nil {
		return err
	}

	// add bunk investigators
	err = addTestInvestigators(cfg)

	// add protocols

	// add investigators to protocols

	// add strains

	// add cage cards

	return nil
}

func addTestPositions(cfg *Config) error {
	fmt.Println("* Creating test positions...")

	posPI := database.CreatePositionParams{
		Title:             "PI",
		CanActivate:       true,
		CanDeactivate:     true,
		CanAddOrders:      true,
		CanQuery:          true,
		CanChangeProtocol: true,
		CanAddStaff:       true,
	}
	cPosPI, err := cfg.db.CreatePosition(context.Background(), posPI)
	if err != nil {
		fmt.Println("Error creating PI position.")
		return err
	}

	posRes := database.CreatePositionParams{
		Title:         "Researcher",
		CanActivate:   true,
		CanDeactivate: true,
		CanAddOrders:  true,
		CanQuery:      true,
	}
	cPosRes, err := cfg.db.CreatePosition(context.Background(), posRes)
	if err != nil {
		fmt.Println("Error creating Researcher position.")
		return err
	}

	posAss := database.CreatePositionParams{
		Title:    "Assistant",
		CanQuery: true,
	}
	cPosAss, err := cfg.db.CreatePosition(context.Background(), posAss)
	if err != nil {
		fmt.Println("Error creating Assistant position.")
		return err
	}

	if verbose {
		fmt.Println(cPosPI)
		fmt.Println(cPosRes)
		fmt.Println(cPosAss)
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
	cInvPI, err := cfg.db.CreateInvestigator(context.Background(), invPI)
	if err != nil {
		fmt.Println("Error creating investigator that is a PI")
		return err
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
	cInvRes, err := cfg.db.CreateInvestigator(context.Background(), invRes)
	if err != nil {
		fmt.Println("Error creating investigator that is a researcher")
		return err
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
	cInvAss, err := cfg.db.CreateInvestigator(context.Background(), invAss)
	if err != nil {
		fmt.Println("Error creating investigator that is an assistant")
		return err
	}

	if verbose {
		fmt.Println(cInvPI)
		fmt.Println(cInvRes)
		fmt.Println(cInvAss)
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
	cProt1, err := cfg.db.CreateProtocol(context.Background(), prot1)
	if err != nil {
		fmt.Println("Error adding test protocol 1")
		return err
	}

	prot2 := database.CreateProtocolParams{
		PNumber:             "18-12-16",
		PrimaryInvestigator: PI[0].ID,
		Title:               "Extended release coagulants",
		Allocated:           115,
		Balance:             110,
		ExpirationDate:      time.Now().AddDate(0, 2, 0),
	}
	cProt2, err := cfg.db.CreateProtocol(context.Background(), prot2)
	if err != nil {
		fmt.Println("Error adding test protocol 2")
		return err
	}

	if verbose {
		fmt.Println(cProt1)
		fmt.Println(cProt2)
	}

	return nil

}
