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
The program is initially set to listen on port 9000. If you wish to modify this setting, please refer to the "main.go" file. Additionally, ensure that the path to the program correctly points to your executable file, and that your current username is reflected in the "supervisor-pastezone.conf" file.
