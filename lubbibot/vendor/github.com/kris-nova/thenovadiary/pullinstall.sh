#!/bin/bash

git pull origin main
go build -o thenovadiary cmd/*
chmod +x thenovadiary
cp thenovadiary /usr/local/bin/thenovadiary
