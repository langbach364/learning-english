#!/bin/bash

FILE_TO_WATCH="./translate/trans.txt"
TRANSLATED_FILE="./translate/trans_ed.txt"
TRANSLATION_MODEL="google"
SOCKET_PATH="./tmp/translation_complete.sock"

translate_file() {
    local input_file=$1
    local output_file=$2

    trans -b -i "$input_file" -o "$output_file" -s en -t vi -e "$TRANSLATION_MODEL"
}

if [ -f "$FILE_TO_WATCH" ]; then
    echo "Bắt đầu dịch file..."
    translate_file "$FILE_TO_WATCH" "$TRANSLATED_FILE"
    echo "Dịch xong. Kết quả được lưu trong $TRANSLATED_FILE" | nc -U $SOCKET_PATH
    exit 0
else
    echo "File cần dịch không tồn tại: $FILE_TO_WATCH"
    exit 1
fi
