#!/bin/bash

# -------------------------------------------------------------
# Simple diameter send command using NokiaAAA
# --------------------------------------------------------------
export _THIS_FILE_DIRNAME=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
source $_THIS_FILE_DIRNAME/env.rc

ORIGIN_HOST=cc.client
ORIGIN_REALM=nokia
APPLICATION_ID=3GPP-Gx
COMMAND=Credit-Control
DESTINATION_HOST=cc.server
DESTINATION_REALM=minsait
DESTINATION_ADDRESS=127.0.0.1:3868

# Test parameters
REQUESTFILE=$_THIS_FILE_DIRNAME/CCRequest.txt

# May be overriden by command line using -count <x>
COUNT=1
LOGLEVEL=info

# Delete Garbage
rm $_THIS_FILE_DIRNAME/out/*

# Diameter CCR -------------------------------------------------------------
echo 
echo Gx Credit Control request
echo

echo Session-Id = \"session-id-1\" > $REQUESTFILE
echo Auth-Application-Id = 16777238 >> $REQUESTFILE
echo CC-Request-Type = 1 >> $REQUESTFILE
echo CC-Request-Number = 1 >> $REQUESTFILE
echo Subscription-Id = \"Subscription-Id-Type=1, Subscription-Id-Data=913374871\" >> $REQUESTFILE


# Send the packet
# -overlap <number of simultaneous requests>
$DIAMETER -debug $LOGLEVEL -count $COUNT -oh $ORIGIN_HOST -or $ORIGIN_REALM -dh $DESTINATION_HOST -dr $DESTINATION_REALM -destinationAddress $DESTINATION_ADDRESS -Application $APPLICATION_ID -command $COMMAND -request "@$REQUESTFILE" $*