WELCOME TO GATOR - This is a CLI that can aggregate RSS feeds for mutiple users.

To run gator, postgres and go are necessary.  Install both.  If you are unsure about the status of Go or Postgres on your machine, try the command go version.  For postgres: psql -version.

To install gator, run the following command:  go install github.com/CoupDeGrace92/gator@latest.

In the .gatorconfig.json - we have a link to a locally hosted postgres database - change this to the location of your local database of choice.  If one is not set up:
    1. create a postgres db with the following command:
    createdb gator
    2. Then you can form a connection to that database:
    postgres://USERNAME:PASSWORD@localhost:5432/gator?sslmode=disable

If you have not set up postgres, open up the Postgres shell as an admin:
    psql -U postgres
Then take the following steps:
    1. Create a user with a password:
    CREATE USER gator_user WITH PASSWORD 'some_strong_password';

    2. Create the database (if it is not already created):
    CREATE DATABASE gator;
    OR IF THE DATABASE IS CREATED
    \c gator

    3. Give the user access to the database:
    GRANT ALL PRIVILAGES ON DATABASE gator TO gator_user

Once installed, gator has the following commands with arguments listed in parenthesis:
    register (name) - sets (name) as the current user and adds name into the users db
    login (name) - sets (name) as the current user in the config
    reset - resets the database to empty
    users - lists all users in the database and indicates the current user
    agg (delay)- meant to be run in a second terminal, in the background will populate the posts db based on feeds the current user follows.  Checks the next feed in line given a delay (eg 10s, 1m, 10m, 1h etc.)  The command might look like: gator agg 10s
    addfeed (feed_name url) - given a url and a custom feed name, will add that feed to the db.  The url should point to the feed but the feed_name can be custom to the user
    feeds - gives a list of all feeds in the feeds database
    follow (url) - allows a user to follow a feed already added by another user
    following - gives a list of feeds the active user follows
    unfollow (url) - allows a user to unfollow the feed at the given url
    browse (num)- prints a list of posts from the feeds the user follows starting with the most recent.  Only returns the num most recent.  If no value is specified, gator will default to 2
