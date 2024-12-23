#!/bin/bash

# Path to your .env file
ENV_FILE="../../.env"

# Check if the .env file exists
if [ -f ".env" ]; then
    # Export the variables
    export $(grep -v '^#' ".env" | xargs)
    echo "Environment variables have been exported."
else
    echo "No .env file found at .env."
fi

