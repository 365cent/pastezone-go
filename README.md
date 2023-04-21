# Pastezone Go
Pastezone Go is a simple pastezone written in Golang and store paste entries in a json file

## Dependencies
Golang installed and added to path

Supervisor installed

## Usage
```bash
mkdir paste && cd paste
git clone https://github.com/365cent/pastezone-go.git
touch pastes.json
cd pastezone-go
go build -o pastezone main.go
cp supervisor-pastezone.conf /etc/supervisor/conf.d
supervisorctl reread
supervisorctl reload
```

To test this program without running as service:
```bash
go run main.go
```

## Note
The program is initially set to listen on port **9000**. If you wish to modify this setting, please refer to the "main.go" file. Additionally, ensure that the path to the program correctly points to your executable file, and that your current username is reflected in the "supervisor-pastezone.conf" file.

Storing pastes in a JSON file is not a secure approach. Simply moving the file to an upper directory will only conceal entries from other visitors.

By default, each paste will expire after 7 days, and the entry will be reused for new pastes once it has expired. If desired, the expiry time can be adjusted by modifying the "expires" variable in the "main.go" file.
