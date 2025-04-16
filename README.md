# musmus-cli
# ᘛ⁐̤ᕐᐷ
musmus is a mouse colony tracker cli written in Golang.

## Set Up
1. musmus uses PostgreSQL as a database. Install postgres, create a database, and find your connection string (the format of a connection string can be found in .env.example)
2. Create a .env file and past the connection string in the DB_URL field (as seen in the .env.example)
3. musmus supports using [pressly/goose](https://github.com/pressly/goose) to set up the db. Copying and pasting the schema into postgres works as well.

## Quick Start
Musmus has demo data that can be loaded on first time start up for testing features. Follow the prompts, then log in as `admin` using password `admin`.
New users will be prompted for a password on first sign in.
Navigate the cli using `goto` from the main menu, and `help` for a list of the available menus.
Every menu and command allows for the `help` command to list available options.
The database can be reset by the admin account under the settings menu. This allows for the test data to be cleared.

## Understanding the problem domain / Why?
musmus tracks and organizes cages by card id, IACUC protocol, staff member, and individual cage notes. musmus allows for personalized reminders and tracking of incoming orders. All cage information is exportalable into CSV files. musmus also monitors protocol allotment and balance, as is requird by IACUC. For keeping track of expenses, musmus calculates the number of care days for a given date range.
Given my years of experience in the administartion of an animal care faciltiy, I started with what I believed would help investigators provide quality care for their animals while performing research.

## Contributing
A wishlist of features I will continue work on is included in the repo. If you'd like to contribute, please for the repository and open a pull request to the 'main' branch. Please make sure to note if you make any changes the the DB schema.