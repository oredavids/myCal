# myCal

Google Calendar CLI app extending the Go quickstart for Google calendar.

![image](output.png)

## Prerequisite

- [Go](https://go.dev/doc/install) minimum version 1.19

## Setup

1. Clone this repository
1. Follow the environment instructions [here](https://developers.google.com/calendar/api/quickstart/go#set_up_your_environment) to obtain the API credentials required.
1. Rename your generated credentials as `myCalAppCredentials.json`
. Rename `.env.sample` as `.env` and replace the value of `MYCAL_CREDENTIALS_DIRECTORY` with the path to the parent directory of your generated credentials file.

> **Note**
> `MYCAL_CREDENTIALS_DIRECTORY` is also where the generated token will be stored. If this env is not set the current working directory will be used; this will have the side effect of needing to generate new tokens in each new directory where the myCal command is run for the first time.

## Run

```cli
go run .
```

### Run from any terminal

1. Build the app

    ```cli
    go build
    ```

1. Confirm the directory where the go app will be installed

    ```cli
    go list -f '{{.Target}}' // Example output: /Users/oredavids/go/bin/myCal
    ```

1. Install the app

    ```cli
    go install
    ```

1. Update your shell config file (e.g bashrc, .zshrc, etc)

    ```cli
    export PATH=$PATH:/Users/oredavids/go // Add directory, confirmed earlier, to your PATH variable
    export MYCAL_CREDENTIALS_DIRECTORY=/Users/path/to/credentials/directory // Where installed app can find your API credentials & store your token
    myCal // OPTIONAL - Run myCal automatically when new terminal window is opened
    ```

    Now that the app has been installed and configured you can run the executable anywhere, manually with:

    ```cli
    myCal
    ```
