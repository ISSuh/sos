#!/bin/bash

# Check if a file name is provided
if [ -z "$1" ]; then
    echo "Usage: $0 <file-to-upload>"
    exit 1
fi

# Define variables
FILE_TO_UPLOAD="$1"
CHUNK_SIZE=1m
FILENAME=$(basename "$FILE_TO_UPLOAD")
TEMP_DIR=$(mktemp -d)
TOTAL_CHUNKS=0

ADDRESS="$2"
GROUP="$3"
PARTITION="$4"
OBJECT_PATH="$5"

URL="http://localhost:8080/upload"

# Split the file into chunks
split -b $CHUNK_SIZE "$FILE_TO_UPLOAD" "$TEMP_DIR/chunk_"

# Count the total number of chunks
TOTAL_CHUNKS=$(ls "$TEMP_DIR"/chunk_* | wc -l)

# Upload each chunk
for CHUNK in "$TEMP_DIR"/chunk_*; do
    CHUNK_NUMBER=$(echo $CHUNK | sed 's/.*chunk_//')
    CHUNK_NUMBER=$(printf "%d" "'$CHUNK_NUMBER")
    CHUNK_NUMBER=$((CHUNK_NUMBER - 96)) # Convert to 1-based index

    echo "Uploading chunk $CHUNK_NUMBER of $TOTAL_CHUNKS..."

    curl -X POST $URL \
        -F "file=@$CHUNK" \
        -F "chunkNumber=$CHUNK_NUMBER" \
        -F "totalChunks=$TOTAL_CHUNKS" \
        -F "fileName=$FILENAME"

    if [ $? -ne 0 ]; then
        echo "Failed to upload chunk $CHUNK_NUMBER"
        exit 1
    fi
done

# Clean up temporary directory
rm -rf "$TEMP_DIR"

echo "File uploaded successfully"
