#!/usr/bin/env bash

# Move the service file to the correct spot
sudo mv /tmp/aats/aats.service /etc/systemd/system/aats.service

# Stop the service, if it's running, so you can replace the files
sudo systemctl stop aats.service

# Make the app destination folder if needed
mkdir -p ~/aats

# Move the files to the appropriate place
mv /tmp/aats/build.gz ~/aats/build.gz
gunzip ~/aats/build.gz
mv /tmp/aats/.env ~/aats/.env

# Cycle systemctl, especially in the event the service didn't exist before
sudo systemctl daemon-reload
sudo systemctl enable aats.service
sudo systemctl start aats.service

# Exit the SSH session
exit
