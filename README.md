# go-rpifancon
A simple temperature based fan controller for a raspberry pi in go.

## Usage
```
  -checks int
         The number of samples of temperatures to check before state changes. (default 3)
  -debug
         Debug info.
  -help
         This info.
  -iopin int
         The GPIO pin used to control the fan. (default 17)
  -threshold int
         The temperature in celsius above which to enable the fan. (default 65)
  -wait int
         The amount of time to wait between polling temperature. Multiply this by checks to get time between pin state changes. (default 5)
```

## Systemd
I run it with a little systemd service like so: 
```
[Unit]
Description=RPi Fan control

[Service]
Type=simple
ExecStart=/usr/local/rpifancon
Restart=always

[Install]
WantedBy=multi-user.target
```

## Build
Remember to build for arm:
```
GOOS=linux GOARCH=arm go build -v main/rpifancon.go
```

## Charts
After some brief testing with the wonderful  [Stressberry](https://github.com/nschloe/stressberry) on my pi 4 it seems to work ok, preventing it throttling at high load:
![stressberry chart with fan](img/test.png)