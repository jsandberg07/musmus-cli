![Image](https://drive.google.com/file/d/1_VdXTD9r81iy57pmL3gaZbmxONBJuomM/view?usp=sharing)
# musmus-cli
# ᘛ⁐̤ᕐᐷ

musmus is a mouse colony tracker cli written in golang. It tracks and organizes cages by card id, IACUC protocol, staff member, and individual cage notes. musmus allows for personalized reminders and tracking of incoming orders. All cage information is exportalable into CSV files. musmus also monitors protocol allotment and balance, as is requird by IACUC. For keeping track of expenses, musmus calculates the number of care days for a given date range.

## Quick Start
musmus uses postgres as a database.
1. Install postgres, create a database, and find your connection string.
2. Create a .env file and past the connection string in the DB_URL field (as seen in the .env.example)
3. musmus supports using pressly/goose to set up the db. Copying and pasting the schema into postgres works as well.

## Usage
Musmus has demo data that can be loaded on first time start up for testing features. Follow the prompts, then log in as 'admin' using password 'admin.'
New users will be prompted for a password on first sign in.
Navigate the cli using 'goto' from the main menu, and 'help' for a list of the available menus.
Every menu and command allows for the 'help' command to list available options.
The database can be reset by the admin account under the settings menu. This allows for the test data to be cleared.